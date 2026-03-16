package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	help_text string = `
Usage: gmake [OPTION]...

A very lightweight build tool.

	--help		display this help and exit
	--version	output version information and exit
	--verbose	set the log level
		2 => ERROR(default)
		3 => WARN
		4 => INFO
		5 => DEBUG
	
`
	version_text = `
    gmake (extrame/gmake) 0.2.1

    Copyright (C) 2014-2022 Abram C. Isola && Liu Ming
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

// findGMakefile 递归向上查找GMakefile文件
func findGMakefile(dir string) (string, error) {
	// 检查当前目录是否存在GMakefile
	filePath := filepath.Join(dir, "GMakefile")
	if _, err := os.Stat(filePath); err == nil {
		return dir, nil
	}

	// 检查父目录
	parentDir := filepath.Dir(dir)
	if parentDir == dir { // 已经到达根目录
		return "", fmt.Errorf("GMakefile not found")
	}

	// 递归查找父目录
	return findGMakefile(parentDir)
}

// Starts processing
func main() {
	// help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	watch := flag.Bool("watch", false, "watch for file changes (not enabled until now)")
	verbose := flag.Uint("verbose", 2, "open verbose log,2: error,3: warning, 4: info, 5: debug")
	nodepenence := flag.Bool("nd", false, "not execute dependencies")
	flag.Parse()

	logrus.SetLevel(logrus.Level(*verbose))

	if *version {
		logrus.Fatal(version_text)
	} else {
		// 获取当前目录
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("gmake: fatal: could not get current directory")
			return
		}

		// 查找GMakefile文件
		gmakeDir, err := findGMakefile(currentDir)
		if err != nil {
			fmt.Println("gmake: fatal: could not find GMakefile")
			return
		}

		// 切换到找到GMakefile的目录
		if err := os.Chdir(gmakeDir); err != nil {
			fmt.Println("gmake: fatal: could not change directory")
			return
		}

		// 读取GMakefile文件
		buf, err := os.ReadFile("GMakefile")
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
