package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

type Directive struct {
	Name         Item
	Dependencies []string
	Commands     []Command
}

func (d *Directive) Exec(doc *Doc, ctx *Context) {
	log.Println("in", d.Name)
	for _, dependency := range d.Dependencies {
		log.Println("go to dependency", dependency)
		dependencies := doc.Select(dependency)
		for _, d := range dependencies {
			d.Exec(doc, ctx)
		}
	}
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
		log.Printf("set env %s to %s\n", d.Name.Id, envs)
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
			fmt.Printf("gmake: fatal: '%s'", err)
			os.Exit(1)
		}
	default:
		for _, c := range d.Commands {
			replacedParts := ctx.replaceVar(c.Parts...)
			cm, parts := replacedParts[0], replacedParts[1:]
			mustSuccess := true
			if cm[0] == '-' {
				mustSuccess = false
			}
			log.Printf("try to exec '%s'", strings.Join(replacedParts, " "))
			cmd := exec.Command(cm, parts...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil && mustSuccess {
				fmt.Printf("gmake: fatal: '%s'", err)
				os.Exit(1)
			}
			log.Println("exec success")
		}
	}
	log.Println("finish")
}
