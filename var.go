package main

import (
	"os"
	"strings"

	"github.com/golang/glog"
)

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		glog.Fatalln(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
