package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
    gmake (extrame/gmake) 0.3.1

    Copyright (C) 2014-2026 Abram C. Isola && Liu Ming
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

// findGMakefile 递归向上查找GMakefile文件，返回实际命中文件的绝对路径
func findGMakefile(dir string) (string, error) {
	// 检查当前目录是否存在GMakefile
	filePath := filepath.Join(dir, "GMakefile")
	if _, err := os.Stat(filePath); err == nil {
		return actualFileName(dir, filePath), nil
	}

	// 检查父目录
	parentDir := filepath.Dir(dir)
	if parentDir == dir { // 已经到达根目录
		return "", fmt.Errorf("GMakefile not found")
	}

	// 递归查找父目录
	return findGMakefile(parentDir)
}

// actualFileName 在大小写不敏感的文件系统（如 macOS 默认的 APFS）上把
// GMakefile 还原成真实文件名。否则磁盘上的 Gmakefile 也会被命中，
// 但报出来的路径会写成 GMakefile，与真实文件对不上。
func actualFileName(dir, fallback string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fallback
	}
	for _, entry := range entries {
		if strings.EqualFold(entry.Name(), filepath.Base(fallback)) {
			return filepath.Join(dir, entry.Name())
		}
	}
	return fallback
}

// Starts processing
func main() {
	// help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	watch := flag.Bool("watch", false, "watch for file changes (not enabled until now)")
	verbose := flag.Uint("verbose", 2, "open verbose log,2: error,3: warning, 4: info, 5: debug")
	nodepenence := flag.Bool("nd", false, "not execute dependencies")
	list := flag.Bool("list", false, "list all targets without executing any command")
	listJSON := flag.Bool("json", false, "print the target list as JSON (use with --list)")
	flag.Parse()

	logrus.SetLevel(logrus.Level(*verbose))

	if *version {
		// 版本信息属于正常输出：走 stdout，退出码 0，不写 stderr
		fmt.Print(version_text)
	} else {
		// 获取当前目录
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("gmake: fatal: could not get current directory")
			return
		}

		// 查找GMakefile文件
		gmakeFile, err := findGMakefile(currentDir)
		if err != nil {
			fmt.Println("gmake: fatal: could not find GMakefile")
			return
		}

		// 报出本次实际命中的配置文件。默认开启、不受 --verbose/--nd 影响；
		// 写 stderr，避免混进任务命令可能被其他脚本消费的 stdout。
		fmt.Fprintf(os.Stderr, "gmake: using %s\n", gmakeFile)

		// 切换到找到GMakefile的目录
		if err := os.Chdir(filepath.Dir(gmakeFile)); err != nil {
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

		// --list 只读：只输出已解析的目标，不执行命令、不展开依赖
		if *list {
			out := buildList(gmakeFile, AST)
			if *listJSON {
				if err := printListJSON(os.Stdout, out); err != nil {
					fmt.Fprintf(os.Stderr, "gmake: fatal: could not output json: %v\n", err)
				}
			} else {
				printListText(os.Stdout, out)
			}
			return
		}

		args := flag.Args()

		logrus.WithField("skip-dependencies", *nodepenence).Infoln("start exec")

		//第一个参数是选择器，其余是指令参数
		AST.Exec(*watch, *nodepenence, args...)
	}
}
