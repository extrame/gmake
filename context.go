package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Context struct {
	Variables map[string]interface{}
	FileTail  map[string]time.Time `json:"file_tail"`
	//variables got from cache file
	oldVariables      map[string]interface{}
	wait              bool
	directivesInStack map[int]bool
}

func (ctx *Context) markVariableStatus() {
	ctx.cache()
}

func (ctx *Context) replaceVar(origin ...string) (dir string, res []string) {
	res = make([]string, 0)
	var cmd []string
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
			explained := os.Expand(word, ctx.expend)
			if cmd != nil && strings.HasSuffix(explained, ")") {
				//in execute
				cmd = append(cmd, explained[:len(explained)-1])
				var c = exec.Command(cmd[0], cmd[1:]...)
				bts, err := c.Output()
				if err != nil {
					logrus.WithField("execution", "get string from sub command").Error(err)
				} else {
					logrus.WithField("result", string(bts)).Info("get result from sub command")
					res = append(res, string(bts))
				}
				cmd = nil
			} else if cmd != nil {
				cmd = append(cmd, explained)
			} else if strings.HasPrefix(explained, "$(") {
				//in execute
				cmd = []string{explained[2:]}
			} else {
				res = append(res, explained)
			}
		}
	}
	return
}

func (ctx *Context) expend(key string) string {
	if v, ok := ctx.Variables[key]; ok {
		return fmt.Sprintf("%s", v)
	} else {
		return os.Getenv(key)
	}
}
