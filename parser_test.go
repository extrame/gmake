package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		name   string
		tokens []LexToken
	}
	tests := []struct {
		name string
		args args
		want Doc
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.name, tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseItem(t *testing.T) {
	type args struct {
		name   string
		tokens []LexToken
	}
	tests := []struct {
		name string
		args args
		want Item
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseItem(tt.args.name, tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_newDirective(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		p    *parser
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.newDirective(tt.args.name)
		})
	}
}

func Test_parser_closeDirective(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.closeDirective()
		})
	}
}

func Test_parser_addCmdPart(t *testing.T) {
	type args struct {
		part string
	}
	tests := []struct {
		name string
		p    *parser
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.addCmdPart(tt.args.part)
		})
	}
}

func Test_parser_flushCommand(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.flushCommand()
		})
	}
}

func Test_parser_peek(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
		want *LexToken
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.peek(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_next(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
		want *LexToken
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.next(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_backup(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.backup()
		})
	}
}

func Test_parser_run(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.run()
		})
	}
}

func Test_initialParserState(t *testing.T) {
	type args struct {
		p *parser
	}
	tests := []struct {
		name string
		args args
		want parserState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initialParserState(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initialParserState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commandsState(t *testing.T) {
	type args struct {
		p *parser
	}
	tests := []struct {
		name string
		args args
		want parserState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandsState(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commandsState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_itemState(t *testing.T) {
	type args struct {
		p *parser
	}
	tests := []struct {
		name string
		args args
		want parserState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := itemState(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("itemState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dependencyState(t *testing.T) {
	type args struct {
		p *parser
	}
	tests := []struct {
		name string
		args args
		want parserState
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dependencyState(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dependencyState() = %v, want %v", got, tt.want)
			}
		})
	}
}
