package main

import "flag"
import "fmt"
import "io/ioutil"

import "github.com/golang/glog"

const (
	help_text string = `
    Usage: gmake [OPTION]...
    
    A very lightweight build tool.

          --help     display this help and exit
          --version  output version information and exit
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
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *help {
		glog.Fatalln(help_text)
	} else if *version {
		glog.Fatalln(version_text)
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

		if len(args) == 0 {
			AST.Exec()
		} else {
			AST.Exec(args[0])
		}
	}
}
