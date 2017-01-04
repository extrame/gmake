package main

import (
	"fmt"
	"log"
	"testing"
)

func TestSelector(t *testing.T) {
	selector := Selector(".main.test#id")
	fmt.Println(selector)
}

func TestSelect(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		env.deploy#GOPATH {
			
		}
		compile.deploy {

		}
        deploy#target1
		(
			.deploy
		){}
    `)
	ds := Parse("GMAKE", tokens)
	selected := ds.Select(".deploy")
	if len(selected) != 2 {
		t.Error("wrong selected directive number")
	}
}

func TestSelect1(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		env.deploy#GOPATH {
			
		}
		compile {

		}
        deploy#target1.main
		(
			.deploy
		){}
    `)
	ds := Parse("GMAKE", tokens)
	selected := ds.Select(".main")
	for _, dir := range selected {
		log.Println(dir)
	}
	if len(selected) != 1 {
		t.Error("wrong selected directive number", len(selected))
	}
}

func TestSelect2(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
		env.deploy#GOPATH {
			
		}
		compile.deploy {

		}
        deploy#target1
		(
			compile
		){}
    `)
	ds := Parse("GMAKE", tokens)
	selected := ds.Select("compile")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number expected '%d', but is '%d'", 1, len(selected))
	}
}
