package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"time"

	"sort"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Directive struct {
	Serial       int
	Name         *Item
	Dependencies []string
	Commands     []*Command
}

func (d *Directive) shouldJumpOver(ctx *Context) bool {
	if len(d.Name.Conditions) > 0 {
		for _, cond := range d.Name.Conditions {
			if cond.isNotAlready(ctx) {
				return true
			}
		}
	}
	return false
}

type Result struct {
	Serial int
	Result bool
}

// nodependencies: 0 - root and no skil 1 - root and skip 2 - n=not root and skip
func (d *Directive) Exec(doc *Doc, ctx *Context, parentChan chan Result, serialNo int, nodependencies int) bool {
	ctx.directivesInStack[d.Serial] = true
	logrus.Debugln("in", d.Name.String())
	var signalNo = 0
	//check Dependencies
	hasDependency := false
	var dependencies Doc
	for _, dependency := range d.Dependencies {
		logrus.Debugln("go to dependency", dependency)
		selected := doc.Select(dependency)
		for _, v := range selected {
			if _, ok := ctx.directivesInStack[v.Serial]; !ok {
				//not in stack already
				dependencies = append(dependencies, v)
			}
		}
	}
	var execChan = make(chan Result, len(dependencies))

	sort.Sort(dependencies)

	for _, d := range dependencies {
		hasDependency = true
		//TODO this should be implented by chan
		// if d.Name.Type == "env" || d.Name.Type == "var" {
		// 	d.Exec(doc, ctx, execChan, signalNo)
		// } else {
		if nodependencies > 0 {
			d.Exec(doc, ctx, execChan, signalNo, notRootAndSkip)
		} else {
			d.Exec(doc, ctx, execChan, signalNo, rootOrNoSkip)
		}
		// }
		signalNo++
	}
	//waiting for dependencies
	myCheckList := make(map[int]bool)

	for {
		exec := 0
		if hasDependency {
			logrus.Infoln("waiting in", d.Name.String())
			select {
			case result := <-execChan:
				logrus.Infoln("waited for one signal in", d.Name.String())
				myCheckList[result.Serial] = result.Result
				for _, v := range myCheckList {
					if v {
						exec++
					}
				}
			}
		}

		resultCode := true
		needWait := false
		if exec == signalNo && nodependencies < notRootAndSkip {
			logrus.Debugln("go to exec")
			var err error
			if d.shouldJumpOver(ctx) {
				logrus.WithField("name", d.Name).Infoln("condition not fill, jump over")
				resultCode = true
				needWait = false
				err = nil
			} else {
				resultCode, needWait, err = d.exec(ctx)
			}

			if err != nil {
				logrus.WithError(err).Fatal("exec error")
			}
		}
		logrus.Debugf("[%s] %d %d", d.Name.String(), len(myCheckList), signalNo)
		if len(myCheckList) == signalNo {
			parentChan <- Result{serialNo, resultCode}
		}
		if len(myCheckList) == signalNo && !ctx.wait {
			break
		} else if !hasDependency {
			if needWait {
				logrus.Infoln("sleep")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	//
	logrus.Debugln("finish")
	return true
}

// exec the command
// return values
// successed
// needWait
// err: error
func (d *Directive) exec(ctx *Context) (successed bool, needWait bool, err error) {
	switch d.Name.Type {
	case "env":

		if d.Name.Id != "" {
			envs := make([]string, 0)
			//set name in title
			for _, c := range d.Commands {
				_, args := ctx.replaceVar(c.Parts...)
				envs = append(envs, args...)
			}
			var newEnv string
			if d.Name.hasClass("append") {
				newEnv = appendEnv(d.Name.Id, envs...)
			} else {
				newEnv = setEnv(d.Name.Id, envs...)
			}
			logrus.
				WithField("setted env", newEnv).
				WithField("input env", envs).
				WithField("name", d.Name.Id).
				Infof("set env")
		} else {
			//set name as yaml
			envs := make(map[string]string)
			str := make([]string, len(d.Commands))
			for n, c := range d.Commands {
				_, env := ctx.replaceVar(c.Parts...)
				str[n] += strings.Join(env, " ")
			}
			var fullStr = strings.Join(str, "\n")
			err = yaml.Unmarshal([]byte(fullStr), &envs)
			if err == nil {
				var newEnv string
				for k, env := range envs {
					if d.Name.hasClass("append") {
						newEnv = appendEnv(k, env)
					} else {
						newEnv = setEnv(k, env)
					}
					logrus.
						WithField("setted env", newEnv).
						WithField("input env", envs).
						WithField("name", k).
						Infof("set env")
				}
			} else {
				err = errors.Wrap(err, "in parse:"+fullStr)
				return
			}
		}

	case "var":
		if ctx.Variables == nil {
			ctx.Variables = make(map[string]interface{})
		}
		str := make([]string, len(d.Commands))
		for n, c := range d.Commands {
			_, vars := ctx.replaceVar(c.Parts...)
			str[n] += strings.Join(vars, " ")
		}
		var fullStr = strings.Join(str, "\n")
		err = yaml.Unmarshal([]byte(fullStr), &ctx.Variables)
		if err != nil {
			err = errors.Wrap(err, "in parse:"+fullStr)
			return
		}
		ctx.markVariableStatus()
	case "import":
	//use go get to install the packages

	//find the packages in GOPATH
	//if user install a executable file in the package, trend it as bin
	//(not else)if there is gmake in the package, trend it as package
	//parse the gmake file in packages dir
	case "file":
		// needWait = ctx.wait
		// logrus.Infoln("watch file...")
		// for _, c := range d.Commands {
		// 	_, replacedParts := ctx.replaceVar(c.Parts...)
		// 	if isChanged := ctx.IsFileChanged(replacedParts[0]); isChanged {
		// 		successed = true
		// 		return
		// 	}
		// }
		// return
	default:
		logrus.Debugln("exec commands", d.Commands)
		for _, c := range d.Commands {
			dir, replacedParts := ctx.replaceVar(c.Parts...)
			cm, parts := replacedParts[0], replacedParts[1:]
			mustSuccess := true
			if replacedParts[0][0] == '-' {
				mustSuccess = false
				if replacedParts[0] == "-" {
					cm = replacedParts[1]
					parts = replacedParts[2:]
				} else {
					cm = replacedParts[0][1:]
				}
			}
			logrus.Infof("try to exec '%s'", strings.Join(replacedParts, " "))
			cmd := exec.Command(cm, parts...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if dir != "" {
				cmd.Dir = filepath.Join(".", dir)
			}
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil && mustSuccess {
				log.Fatalf("gmake: fatal: '%s'", err)
			}
			logrus.Infoln("exec success")
		}
	}
	successed = true
	return
}
