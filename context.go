package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Context struct {
	variables         map[string]interface{}
	wait              bool
	directivesInStack map[int]bool
}

func (ctx *Context) replaceVar(origin ...string) (dir string, res []string) {
	res = make([]string, 0)
	for i := 0; i < len(origin); i++ {
		var word = origin[i]
		switch word {
		case "cd":
			if i < len(origin)-2 {
				//dir and &&
				dir = origin[i+1]
				if origin[i+2] != "&&" {
					logrus.Fatal("need '&&' for cd dir and run cmd")
				}
				i = i + 2
			} else {
				logrus.Fatal("wrong argument for cd dir,need dirname and &&")
			}
		case "$currentDir":
			res = append(res, getCurrentDirectory())
		default:
			res = append(res, os.Expand(word, func(key string) string {
				if v, ok := ctx.variables[key]; ok {
					return fmt.Sprintf("%s", v)
				} else {
					return os.Getenv(key)
				}
			}))
		}
	}
	return
}
