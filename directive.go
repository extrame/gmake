package main

import (
	"os"
	"os/exec"
	"strings"

	"time"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

type Directive struct {
	Serial       int
	Name         Item
	Dependencies []string
	Commands     []Command
}

type Result struct {
	Serial int
	Result bool
}

func (d *Directive) Exec(doc *Doc, ctx *Context, parentChan chan Result, serialNo int) bool {
	ctx.directivesInStack[d.Serial] = true
	glog.Infoln("in", d.Name.String())
	var signalNo = 0
	//check Dependencies
	hasDependency := false
	var dependencies Doc
	for _, dependency := range d.Dependencies {
		glog.Infoln("go to dependency", dependency)
		selected := doc.Select(dependency)
		for _, v := range selected {
			if _, ok := ctx.directivesInStack[v.Serial]; !ok {
				dependencies = append(dependencies, v)
			}
		}
	}
	var execChan = make(chan Result, len(dependencies))

	for _, d := range dependencies {
		hasDependency = true
		go d.Exec(doc, ctx, execChan, signalNo)
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
				glog.Infoln(myCheckList, exec, signalNo)
			}
		}

		resultCode := true
		needWait := false
		if exec == signalNo {
			glog.Infoln("go to exec")
			resultCode, needWait = d.exec(ctx)
		}
		parentChan <- Result{serialNo, resultCode}
		if len(myCheckList) == serialNo && !ctx.wait {
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
	glog.Infoln("finish")
	return true
}

func (d *Directive) exec(ctx *Context) (bool, bool) {
	switch d.Name.Type {
	case "env":
		envs := make([]string, 0)
		for _, c := range d.Commands {
			envs = append(envs, ctx.replaceVar(c.Parts...)...)
		}
		if d.Name.hasClass("append") {
			appendEnv(d.Name.Id, envs...)
		} else {
			setEnv(d.Name.Id, envs...)
		}
		glog.Infof("set env %s to %s\n", d.Name.Id, envs)
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
			replacedParts := ctx.replaceVar(c.Parts...)
			if isChanged := IsFileChanged(replacedParts[0]); isChanged {
				return true, ctx.wait
			}
		}
		return false, ctx.wait
	default:
		for _, c := range d.Commands {
			replacedParts := ctx.replaceVar(c.Parts...)
			cm, parts := replacedParts[0], replacedParts[1:]
			mustSuccess := true
			if cm[0] == '-' {
				mustSuccess = false
			}
			glog.Infof("try to exec '%s'", strings.Join(replacedParts, " "))
			cmd := exec.Command(cm, parts...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil && mustSuccess {
				glog.Fatalf("gmake: fatal: '%s'", err)
			}
			glog.Infoln("exec success")
		}
	}
	return true, false
}
