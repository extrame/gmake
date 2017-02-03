package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"os"

	"github.com/golang/glog"
)

type Environment struct {
	FileTail map[string]time.Time `json:"file_tail"`
}

var Env *Environment

func IsFileChanged(name string) bool {
	var err error
	defer func() {
		if bts, err := json.Marshal(Env); err != nil {
			glog.Fatalln(err)
		} else {
			var out bytes.Buffer
			json.Indent(&out, bts, "", "  ")
			ioutil.WriteFile(".gmake", out.Bytes(), 0666)
		}
	}()
	if Env == nil {
		var bts []byte
		if bts, err = ioutil.ReadFile(".gmake"); err == nil {
			err = json.Unmarshal(bts, &Env)
		}
	}
	if err != nil {
		glog.Warningln(err)
		Env = new(Environment)
		Env.FileTail = make(map[string]time.Time)
	}
	var matches []string
	var res bool
	if matches, err = filepath.Glob(name); err == nil {
		for _, v := range matches {
			if fi, err := os.Stat(v); err == nil {
				if oldMT, ok := Env.FileTail[v]; (ok && fi.ModTime().Sub(oldMT) > 0) || !ok {
					glog.Infoln(v, "is changed")
					res = true
				}
				Env.FileTail[v] = fi.ModTime()
			}
		}
	}
	return res
}
