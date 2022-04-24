package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

const (
	help_text string = `
Usage: gmake [OPTION]...

A very lightweight build tool.

	--help		display this help and exit
	--version	output version information and exit
	--watch		watch the file
	--verbose	set the log level
		2 => ERROR(default)
		3 => WARN
		4 => INFO
		5 => DEBUG
	
`
	version_text = `
    gmake (aisola/gmake) 0.1

    Copyright (C) 2014-2017 Abram C. Isola && Liu Ming
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var AST Doc

func combiner(strs []string) string {
	var fullstr string
	for i := 0; i < len(strs); i++ {
		fullstr = fullstr + strs[i] + " "
	}
	return fullstr
}

// Starts processing
func main() {
	// help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	watch := flag.Bool("watch", false, "watch for file changes")
	verbose := flag.Uint("verbose", 2, "open verbose log,2: error,3: warning, 4: info, 5: debug")
	nodepenence := flag.Bool("nd", false, "not execute dependencies")
	flag.Parse()

	logrus.SetLevel(logrus.Level(*verbose))

	if *version {
		logrus.Fatal(version_text)
	} else {
		// get contents
		buf, err := ioutil.ReadFile("GMakefile")
		if err != nil {
			fmt.Println("gmake: fatal: could not read GMakefile")
			return
		}

		// scan then parse
		_, tokens := Lexer("GMAKE", string(buf))
		AST = Parse("GMAKE", tokens)

		args := flag.Args()

		logrus.WithField("skip-dependencies", *nodepenence).Infoln("start exec")

		if len(args) == 0 {
			AST.Exec(*watch, *nodepenence)
		} else {
			AST.Exec(*watch, *nodepenence, args[0])
		}
	}
}
