package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"time"

	"sort"

	"github.com/BurntSushi/toml"
	glog "github.com/sirupsen/logrus"
)

type Directive struct {
	Serial       int
	Name         Item
	Dependencies []string
	Commands     []*Command
}

type Result struct {
	Serial int
	Result bool
}

// nodependencies: 0 - root and no skil 1 - root and skip 2 - n=not root and skip
func (d *Directive) Exec(doc *Doc, ctx *Context, parentChan chan Result, serialNo int, nodependencies int) bool {
	ctx.directivesInStack[d.Serial] = true
	glog.Debugln("in", d.Name.String())
	var signalNo = 0
	//check Dependencies
	hasDependency := false
	var dependencies Doc
	for _, dependency := range d.Dependencies {
		glog.Debugln("go to dependency", dependency)
		selected := doc.Select(dependency)
		for _, v := range selected {
			if _, ok := ctx.directivesInStack[v.Serial]; !ok {
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
			glog.Infoln("waiting in", d.Name.String())
			select {
			case result := <-execChan:
				glog.Infoln("waited for one signal in", d.Name.String())
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
			glog.Debugln("go to exec")
			resultCode, needWait = d.exec(ctx)
		}
		glog.Debugf("[%s] %d %d", d.Name.String(), len(myCheckList), signalNo)
		if len(myCheckList) == signalNo {
			parentChan <- Result{serialNo, resultCode}
		}
		if len(myCheckList) == signalNo && !ctx.wait {
			break
		} else if !hasDependency {
			if needWait {
				glog.Infoln("sleep")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	//
	glog.Debugln("finish")
	return true
}

func (d *Directive) exec(ctx *Context) (bool, bool) {
	switch d.Name.Type {
	case "env":
		envs := make([]string, 0)
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
		glog.
			WithField("setted env", newEnv).
			WithField("input env", envs).
			WithField("name", d.Name.Id).
			Infof("set env")
	case "var":
		if ctx.variables == nil {
			ctx.variables = make(map[string]interface{})
		}
		str := make([]string, len(d.Commands))
		for n, c := range d.Commands {
			str[n] += strings.Join(c.Parts, " ")
		}
		_, err := toml.Decode(strings.Join(str, "\n"), &ctx.variables)
		if err != nil {
			glog.Fatalf("gmake: fatal: '%s'", err)
		}
	case "import":
	//use go get to install the packages

	//find the packages in GOPATH
	//if user install a executable file in the package, trend it as bin
	//(not else)if there is gmake in the package, trend it as package
	//parse the gmake file in packages dir
	case "file":
		glog.Infoln("watch file...")
		for _, c := range d.Commands {
			_, replacedParts := ctx.replaceVar(c.Parts...)
			if isChanged := IsFileChanged(replacedParts[0]); isChanged {
				return true, ctx.wait
			}
		}
		return false, ctx.wait
	default:
		glog.Debugln("exec commands", d.Commands)
		for _, c := range d.Commands {
			dir, replacedParts := ctx.replaceVar(c.Parts...)
			cm, parts := replacedParts[0], replacedParts[1:]
			mustSuccess := true
			if replacedParts[0][0] == '-' {
				mustSuccess = false
			}
			glog.Infof("try to exec '%s'", strings.Join(replacedParts, " "))
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
			glog.Infoln("exec success")
		}
	}
	return true, false
}
