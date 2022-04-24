package main

import (
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Environment struct {
	FileTail map[string]time.Time `yaml:"file_tail"`
}

// var Env *Environment

func (ctx *Context) cache() {
	if bts, err := yaml.Marshal(ctx); err != nil {
		logrus.Fatalln(err)
	} else {
		ioutil.WriteFile(".gmake", bts, 0666)
	}
}

//got content from cache
func (ctx *Context) load() error {
	bts, err := ioutil.ReadFile(".gmake")
	if err == nil {
		err = yaml.Unmarshal(bts, ctx)
		if err == nil && ctx.Variables != nil {
			ctx.oldVariables = make(map[string]interface{})
			for k, v := range ctx.Variables {
				ctx.oldVariables[k] = v
			}
		}
	}
	return err
}

// func (ctx *Context) IsFileChanged(name string) bool {
// 	var err error

// 	if err != nil {
// 		logrus.Warningln(err)
// 		Env = new(Environment)
// 		Env.FileTail = make(map[string]time.Time)
// 	}
// 	var matches []string
// 	var res bool
// 	if matches, err = filepath.Glob(name); err == nil {
// 		for _, v := range matches {
// 			if fi, err := os.Stat(v); err == nil {
// 				if oldMT, ok := Env.FileTail[v]; (ok && fi.ModTime().Sub(oldMT) > 0) || !ok {
// 					logrus.Infoln(v, "is changed")
// 					res = true
// 				}
// 				Env.FileTail[v] = fi.ModTime()
// 			}
// 		}
// 	}
// 	return res
// }
