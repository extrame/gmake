package main

import (
<<<<<<< HEAD
=======
	"fmt"
	"os"
>>>>>>> 0e83d21c43682e10fd195d42b3b943bd0e4a94cf
	"strconv"
	"strings"
	"unicode/utf8"
	"os"

	"fmt"
)

// LexToken holds is a (type, value) array.
type LexToken [3]string

// EOF character
var EOF string = "+++EOF+++"

// lexerState represents the state of the scanner
// as a function that returns the next state.
type lexerState func(*lexer) lexerState

// Lexer creates a new scanner for the input string.
func Lexer(name, input string, initialState ...lexerState) (*lexer, []LexToken) {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make([]LexToken, 0),
		lineno: 1,
	}
	if len(initialState) > 0 {
		l.initialState = initialState[0]
	} else {
		l.initialState = initialLexerState
	}
	l.Run()
	return l, l.tokens
}

// lexer holds the state of the scanner.
type lexer struct {
	name         string     // used only for error reports.
	input        string     // the string being scanned.
	start        int        // start position of this item.
	pos          int        // current position in the input.
	width        int        // width of last rune read from input.
	tokens       []LexToken // scanned items.
	initialState lexerState

	lineno int
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// next returns the next rune in the input.
func (l *lexer) next() string {
	var r rune
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return string(r)
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() string {
	r := l.next()
	l.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.Index(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.Index(valid, l.next()) >= 0 {
	}
	l.backup()
}

// Run consumes a run of runes from the invalid set.
func (l *lexer) rejectRun(valid string) {
	for c := l.next(); strings.Index(valid, c) < 0; c = l.next() {
		fmt.Print(c)
		if c == EOF {
			os.Exit(1)
		}
	}
	l.backup()
}

// error returns an error and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) lexerState {
	l.tokens = nil

	args = append([]interface{}{l.lineno},args...)

	fmt.Fprintf(os.Stderr,"gmake:%d: "+format+"\n", args...)
	return nil
}

// emit passes an item back to the client.
func (l *lexer) emit(t string) {
	l.tokens = append(l.tokens, LexToken{t, l.input[l.start:l.pos], strconv.Itoa(l.lineno)})
	l.start = l.pos
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) Run() {
	for state := l.initialState; state != nil; {
		state = state(l)
	}
}

// isName() checks if a character is an alpha
func isName(char string) bool {
	testStr := alphavalues + classMarker + idMarker + platformMaker + splitMaker
	if strings.Index(testStr, char) >= 0 {
		return true
	} else {
		return false
	}
}

// isName() checks if a character is an alpha
func isCharacter(char string) bool {
	if strings.Index(alphavalues, char) >= 0 {
		return true
	} else {
		return false
	}
}

func itemLexerState(l *lexer) lexerState {
	for r := l.next(); r != EOF; r = l.next() {
		if r == "\r" {
			l.ignore()
		} else if r == "\n" {
			l.lineno += 1
			l.ignore()
		} else if r == "." {
			l.emit(T_CLASS_MARK)
		} else if r == "#" {
			l.emit(T_ID_MARK)
		} else if isCharacter(r) {
			l.acceptRun(alphavalues + numbers + splitMaker)
			if len(l.tokens) >= 1 {
				switch l.tokens[len(l.tokens)-1][0] {
				case T_CLASS_MARK:
					l.emit(T_LCLASS)
				case T_ID_MARK:
					l.emit(T_LID)
				default:
					l.emit(T_LITEM)
				}
			} else {
				l.emit(T_LITEM)
			}
		} else if r == "{" || r == "(" {
			l.backup()
			l.emit(T_RITEM)
			return initialLexerState
		}
	}
	l.emit(T_EOF)
	return nil
}

// initialState is the starting point for the
// scanner. It scans through each character and decides
// which state to create for the lexer. lexerState == nil
// is exit scanner.
func initialLexerState(l *lexer) lexerState {
	for r := l.next(); r != EOF; r = l.next() {
		if r == " " || r == "\t" || r == "\r" {
			l.ignore()
		} else if r == "\n" {
			l.lineno += 1
			l.ignore()
		} else if isName(r) {
			l.backup()
			return itemLexerState
		} else if r == "{" {
			l.emit(T_LCBRAC)
			return commandState
		} else if r == "(" {
			l.emit(T_LPAREN)
			return dependencyLexerState
		} else {
			return l.errorf("Illegal character '%s'.", r)
		}
	}

	l.emit(T_EOF)
	return nil
}

func dependencyLexerState(l *lexer) lexerState {
	for r := l.next(); r != ")"; r = l.next() {
		if r == " " || r == "\t" || r == "\r" {
			l.ignore()
		} else if r == "\n" {
			l.lineno += 1
			l.emit(T_COMMA)
		} else if r == EOF {
			return l.errorf("Unclosed dependancy switch...")
		} else if isName(r) {
			l.acceptRun(alphavalues + numbers + classMarker + idMarker + platformMaker + splitMaker)
			l.emit(T_CMDPART)
		} else if r == "," {
			l.emit(T_COMMA)
		} else {
<<<<<<< HEAD
=======
			fmt.Println("Hello world!", r)
>>>>>>> 0e83d21c43682e10fd195d42b3b943bd0e4a94cf
			return l.errorf("Illegal character '%s'.", r)
		}
	}
	l.emit(T_RPAREN)
	return initialLexerState
}

func commandState(l *lexer) lexerState {
	for r := l.next(); r != "}"; r = l.next() {
		if r == " " || r == "\t" || r == "\r" {
			l.ignore()
		} else if r == "\"" {
			l.ignore()
			l.rejectRun("\"")
			l.emit(T_CMDPART)
			l.next()
		} else if r == "\n" {
			l.lineno += 1
			l.emit(T_SEMI)
		} else if r == EOF {
			return l.errorf("Unclosed statement...")
		} else {
			l.acceptRun(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*_+=-|\]['":/?.>,<`)
			l.emit(T_CMDPART)
		}
	}
	l.emit(T_RCBRAC)
	return initialLexerState
}
