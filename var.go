package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatalln(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
