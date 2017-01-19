package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		glog.Fatalln(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
