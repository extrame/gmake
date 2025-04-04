package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func Test_getCurrentDirectory(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	fmt.Println("--------")
	fmt.Println(getCurrentDirectory())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCurrentDirectory(); got != tt.want {
				t.Errorf("getCurrentDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		var.deploy {
			test: git.oschina.net/mink-tech/bible
			test2: 23
		}
        deploy.main(
            .deploy
        ){
			go get $test
		}
    `)
	ds := Parse("GMAKE", tokens)
	ds.Exec(false, false, "")
}

func TestVarUpdated(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		var.deploy {
			test: git.oschina.net/mink-tech/bible
			test2: $(dir .)
		}
        deploy.main@test2:updated (
            .deploy
        ){
			go get $test
		}
    `)
	ds := Parse("GMAKE", tokens)
	ds.Exec(false, false, "")
}

func TestExecLs(t *testing.T) {
	var cmd = exec.Command("dir", ".")
	fmt.Println(cmd.Env)
	bts, err := cmd.Output()
	fmt.Println(string(bts), err)
}
