package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// parserState represents the state of the scanner
// as a function that returns the next state.
type parserState func(*parser) parserState

// Parse creates a new parser with the recommended
// parameters.
func Parse(name string, ctx *lang.Context, tokens []LexToken) Doc {
	p := &parser{
		name:   name,
		tokens: tokens,
		pos:    -1,
	}
	p.initState = initialParserState
	p.run()
	return p.ast
}

func ParseItem(name string, tokens []LexToken) Item {
	p := &parser{
		name:   name,
		tokens: tokens,
		pos:    -1,
		serial: 0,
	}
	p.newDirective(name)
	p.initState = itemState
	p.run()
	return p.currentDirective.Name
}

// the parser type
type parser struct {
	name   string
	tokens []LexToken
	pos    int
	serial int

	ast              Doc // the ast
	initState        parserState
	currentDirective *Directive
	cmdparts         *Command
}

func (p *parser) newDirective(name string) {
	p.currentDirective = &Directive{Name: Item{Type: name}, Commands: make([]*Command, 0), Serial: p.serial}
	p.serial++
	// p.cmdparts = Command{Parts: make([]string, 0)}
}

func (p *parser) closeDirective() {
	p.ast = append(p.ast, p.currentDirective)
}

func (p *parser) addCmdPart(part string) {
	logrus.WithField("parts", part).Debugln("add cmd parts")
	if p.cmdparts == nil {
		p.cmdparts = &Command{Parts: make([]string, 0)}
	}
	p.cmdparts.Parts = append(p.cmdparts.Parts, part)
}

func (p *parser) flushCommand() {
	if p.cmdparts != nil {
		p.currentDirective.Commands = append(p.currentDirective.Commands, p.cmdparts)
		p.cmdparts = nil
	}
}

// peek returns what the next token is but does NOT
// advance the position.
func (p *parser) peek() *LexToken {
	tok := p.next()
	p.backup()
	return tok
}

// nest returns what the next token AND
// advances p.pos.
func (p *parser) next() *LexToken {
	if p.pos >= len(p.tokens)-1 {
		return nil
	}
	p.pos += 1
	return &p.tokens[p.pos]
}

// backup sets the position back one token
func (p *parser) backup() {
	p.pos -= 1
}

// run starts the statemachine
func (p *parser) run() {
	for state := p.initState; state != nil; {
		state = state(p)
	}
}

// the starting state for parsing
func initialParserState(p *parser) parserState {
	for t := p.next(); t != nil && t[0] != T_EOF; t = p.next() {
		if t[0] == T_LITEM {
			p.newDirective(t[1])
			return itemState
		} else if t[0] == T_LCBRAC {
			return commandsState
		} else if t[0] == T_LPAREN {
			return dependencyState
		} else {
			fmt.Printf("Doc:%s: unexpected '%s' expecting T_DIRECT\n", t[2], t[0])
			return nil
		}

	}

	return nil
}

func commandsState(p *parser) parserState {

	for t := p.next(); t[0] != T_RCBRAC; t = p.next() {

		if t[0] == T_CMDPART {
			p.addCmdPart(t[1])
		} else if t[0] == T_SEMI {
			p.flushCommand()

		} else {
			fmt.Printf("Doc:%s: unexpected '%s' expecting T_CMDPART or T_SEMI\n", t[2], t[0])
			return nil

		}

	}
	p.closeDirective()
	return initialParserState
}

func itemState(p *parser) parserState {

	for t := p.next(); t[0] != T_RITEM; t = p.next() {
		switch t[0] {
		case T_LCLASS:
			p.currentDirective.Name.Classes = append(p.currentDirective.Name.Classes, t[1])
		case T_LID:
			if p.currentDirective.Name.Id != "" {
				fmt.Printf("Doc:%s: repeated '%s' for directive id with %s\n", t[0], t[1], p.currentDirective.Name.Id)
				return nil
			}
			p.currentDirective.Name.Id = t[1]
		case T_EOF:
			return nil
		case T_LITEM:
			p.newDirective(t[1])
		}
	}
	return initialParserState
}

func dependencyState(p *parser) parserState {

	for t := p.next(); t[0] != T_RPAREN; t = p.next() {
		switch t[0] {
		case T_CMDPART:
			p.currentDirective.Dependencies = append(p.currentDirective.Dependencies, t[1])
		}
	}
	return initialParserState
}
