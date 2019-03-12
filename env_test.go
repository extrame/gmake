package main

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		env.deploy#GOPATH {
			D://test
		}
        deploy#target1.main
		(
			.deploy
		){}
    `)
	ds := Parse("GMAKE", tokens)
	ds.Exec(false)
	setted := os.Getenv("GOPATH")
	if setted != "D://test" {
		t.Errorf("set env error expected 'D://test' but is '%s'", setted)
	}
}

func TestEnvAppend(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		env.deploy#GOPATH.append {
			D://test
		}
        deploy#target1.main
		(
			.deploy
		){}
    `)
	ds := Parse("GMAKE", tokens)
	ds.Exec(false)
	setted := os.Getenv("GOPATH")
	if setted == "D://test" {
		t.Error("set env error not expected 'D://test' ", setted)
	}
}

func TestEnvWithVar(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		env.deploy#GOPATH {
			$currentDir
		}
        deploy#target1.main
		(
			.deploy
		){}
    `)
	ds := Parse("GMAKE", tokens)
	ds.Exec(false)
	setted := os.Getenv("GOPATH")
	if setted != getCurrentDirectory() {
		t.Errorf("set env error expected '%s' but is '%s'", getCurrentDirectory(), setted)
	}
}
