package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

type Directive struct {
	Name         Item
	Dependencies []string
	Commands     []Command
}

func (d *Directive) Exec(doc *Doc, ctx *Context) bool {
	glog.Infoln("in", d.Name.String())
	var _continue = true
	for _, dependency := range d.Dependencies {
		glog.Infoln("go to dependency", dependency)
		dependencies := doc.Select(dependency)
		for _, d := range dependencies {
			_continue = _continue && d.Exec(doc, ctx)
			if !_continue {
				break
			}
		}
	}
	if _continue {
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
			for _, c := range d.Commands {
				replacedParts := ctx.replaceVar(c.Parts...)
				if isChanged := IsFileChanged(replacedParts[0]); !isChanged {
					return false
				}
			}
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
	}
	glog.Infoln("finish")
	return true
}
