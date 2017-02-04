package main

import (
	"fmt"
	"os"
)

type Context struct {
	variables         map[string]interface{}
	wait              bool
	directivesInStack map[int]bool
}

func (ctx *Context) replaceVar(origin ...string) []string {
	res := make([]string, len(origin))
	for n, word := range origin {
		switch word {
		case "$currentDir":
			res[n] = getCurrentDirectory()
		default:
			res[n] = os.Expand(word, func(key string) string {
				return fmt.Sprintf("%s", ctx.variables[key])
			})
		}
	}
	return res
}
