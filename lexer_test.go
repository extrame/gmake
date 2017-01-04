package main

import (
	"reflect"
	"testing"
)

func TestSingleDirective(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy.main{
            env.deploy
        }
    `)
	ds := Parse("GMAKE", tokens)
	if len(ds) > 1 {
		t.Error("wrong directive number")
	}
}

func TestSingleClass(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy.main{
            env.deploy
        }
    `)
	ds := Parse("GMAKE", tokens)
	if len(ds[0].Name.Classes) != 1 {
		t.Error("wrong class number")
	} else if ds[0].Name.Classes[0] != "main" {
		t.Error("wrong class content")
	}
}

func TestDoubleClass(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy.main.test{
            env.deploy
        }
    `)
	ds := Parse("GMAKE", tokens)
	if len(ds[0].Name.Classes) != 2 {
		t.Error("wrong class number")
	} else if ds[0].Name.Classes[0] != "main" && ds[0].Name.Classes[1] != "test" {
		t.Error("wrong class content")
	}
}

func TestDirectiveWithID(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy#target1{
            env.deploy
        }
    `)
	ds := Parse("GMAKE", tokens)
	if len(ds) != 1 {
		t.Error("wrong directive number")
	}
}

func TestSingleID(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy#target1{
            env.deploy
        }
    `)
	ds := Parse("GMAKE", tokens)
	if ds[0].Name.Id != "target1" {
		t.Errorf("wrong id content,expected target1 but is '%s'", ds[0].Name.Id)
	}
}

func TestDependency(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy#target1
		(
			.deploy
		){}
    `)
	ds := Parse("GMAKE", tokens)
	if len(ds[0].Dependencies) != 1 {
		t.Errorf("wrong Dependencies number, expected 1 but %d", len(ds[0].Dependencies))
	} else if ds[0].Dependencies[0] != ".deploy" {
		t.Errorf("wrong Dependencies content,expected .deploy but is '%s'", ds[0].Dependencies[0])
	}
}

func TestMultiDependency(t *testing.T) {
	_, tokens := Lexer("GMAKE", `
        deploy#target1
		(
			.deploy
			compile
		){}
    `)
	ds := Parse("GMAKE", tokens)
	if len(ds[0].Dependencies) != 2 {
		t.Errorf("wrong Dependencies number, expected 1 but %d", len(ds[0].Dependencies))
	} else if ds[0].Dependencies[0] != ".deploy" {
		t.Errorf("wrong Dependencies content,expected .deploy but is '%s'", ds[0].Dependencies[0])
	} else if ds[0].Dependencies[1] != "compile" {
		t.Errorf("wrong Dependencies content,expected compile but is '%s'", ds[0].Dependencies[0])
	}
}

func TestLexer(t *testing.T) {
	type args struct {
		name         string
		input        string
		initialState []lexerState
	}
	tests := []struct {
		name  string
		args  args
		want  *lexer
		want1 []LexToken
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Lexer(tt.args.name, tt.args.input, tt.args.initialState...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Lexer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_lexer_ignore(t *testing.T) {
	tests := []struct {
		name string
		l    *lexer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.ignore()
		})
	}
}

func Test_lexer_backup(t *testing.T) {
	tests := []struct {
		name string
		l    *lexer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.backup()
		})
	}
}

func Test_lexer_next(t *testing.T) {
	tests := []struct {
		name string
		l    *lexer
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.next(); got != tt.want {
				t.Errorf("lexer.next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lexer_peek(t *testing.T) {
	tests := []struct {
		name string
		l    *lexer
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.peek(); got != tt.want {
				t.Errorf("lexer.peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lexer_accept(t *testing.T) {
	type args struct {
		valid string
	}
	tests := []struct {
		name string
		l    *lexer
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.accept(tt.args.valid); got != tt.want {
				t.Errorf("lexer.accept() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lexer_acceptRun(t *testing.T) {
	type args struct {
		valid string
	}
	tests := []struct {
		name string
		l    *lexer
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.acceptRun(tt.args.valid)
		})
	}
}

func Test_lexer_errorf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		l    *lexer
		args args
		want lexerState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.errorf(tt.args.format, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexer.errorf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lexer_emit(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		l    *lexer
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.emit(tt.args.t)
		})
	}
}

func Test_lexer_Run(t *testing.T) {
	tests := []struct {
		name string
		l    *lexer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Run()
		})
	}
}

func Test_isName(t *testing.T) {
	type args struct {
		char string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isName(tt.args.char); got != tt.want {
				t.Errorf("isName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isCharacter(t *testing.T) {
	type args struct {
		char string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCharacter(tt.args.char); got != tt.want {
				t.Errorf("isCharacter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_itemLexerState(t *testing.T) {
	type args struct {
		l *lexer
	}
	tests := []struct {
		name string
		args args
		want lexerState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := itemLexerState(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("itemLexerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initialLexerState(t *testing.T) {
	type args struct {
		l *lexer
	}
	tests := []struct {
		name string
		args args
		want lexerState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initialLexerState(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initialLexerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dependencyLexerState(t *testing.T) {
	type args struct {
		l *lexer
	}
	tests := []struct {
		name string
		args args
		want lexerState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dependencyLexerState(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dependencyLexerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commandState(t *testing.T) {
	type args struct {
		l *lexer
	}
	tests := []struct {
		name string
		args args
		want lexerState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandState(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commandState() = %v, want %v", got, tt.want)
			}
		})
	}
}
