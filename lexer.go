package main

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

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
	return lex(name, input, true, initialState...)
}

// lex 与 Lexer 相同，但可以控制是否把行首的 # 当作行注释。
// 选择器本身允许以 # 开头（ID 选择器），解析选择器时需要关闭该规则。
func lex(name, input string, stripLineComments bool, initialState ...lexerState) (*lexer, []LexToken) {
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
	if !l.stripComments(stripLineComments) {
		return l, l.tokens
	}
	l.Run()
	return l, l.tokens
}

// stripComments 在扫描前移除 GMakefile 中的注释，支持两种写法：
//   - 块注释：/* ... */，可以跨行，也可以出现在选择器或命令中间
//   - 行注释：以 # 开头（允许前置空白）的整行注释
//
// 双引号字符串中的注释标记会原样保留；被移除的注释中的换行会被保留，
// 以便 token 行号仍然准确。
func (l *lexer) stripComments(stripLineComments bool) bool {
	input := l.input
	var b strings.Builder
	b.Grow(len(input))

	var (
		inBlock     bool
		inQuote     bool
		atLineStart = true
		blockLine   = 1
		line        = 1
	)

	for i := 0; i < len(input); i++ {
		c := input[i]

		if inBlock {
			if c == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlock = false
				i++
			} else if c == '\n' {
				b.WriteByte(c)
				line++
				atLineStart = true
			}
			continue
		}

		if inQuote {
			b.WriteByte(c)
			if c == '\n' {
				line++
			} else if c == '"' {
				inQuote = false
			}
			continue
		}

		if atLineStart {
			if c == ' ' || c == '\t' || c == '\r' {
				b.WriteByte(c)
				continue
			}
			if stripLineComments && c == '#' {
				// 跳过整行，保留行尾换行由后续循环处理
				for i+1 < len(input) && input[i+1] != '\n' {
					i++
				}
				continue
			}
			atLineStart = false
		}

		switch {
		case c == '"':
			inQuote = true
			b.WriteByte(c)
		case c == '/' && i+1 < len(input) && input[i+1] == '*':
			// 用空格代替注释起始标记，避免把前后内容拼接成一个词
			inBlock = true
			blockLine = line
			b.WriteByte(' ')
			i++
		default:
			b.WriteByte(c)
			if c == '\n' {
				line++
				atLineStart = true
			}
		}
	}

	if inBlock {
		l.lineno = blockLine
		l.errorf("Unclosed comment...")
		return false
	}
	l.input = b.String()
	return true
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

	args = append([]interface{}{l.lineno}, args...)

	fmt.Fprintf(os.Stderr, "gmake:%d: "+format+"\n", args...)
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
	testStr := alphavalues + classMarker + idMarker + conditionMaker + splitMaker + attributeMarker + attributeEqual + attributeEndMarker
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
		} else if r == idMarker {
			l.emit(T_ID_MARK)
		} else if r == conditionMaker {
			l.emit(T_CONDITION_MARK)
		} else if r == pseudoMaker {
			l.emit(T_PSEUDO_MARK)
		} else if r == attributeMarker { // 新增属性选择器开始标记
			l.emit(T_ATTRIBUTE_MARK)
		} else if r == attributeEqual { // 属性等号
			l.emit(T_ATTRIBUTE_EQUAL)
		} else if r == attributeEndMarker { // 属性结束标记
			l.emit(T_RATTRIBUTE)
		} else if isCharacter(r) {
			l.acceptRun(alphavalues + numbers + splitMaker + nameSeparator)
			if len(l.tokens) >= 1 {
				switch l.tokens[len(l.tokens)-1][0] {
				case T_CLASS_MARK:
					l.emit(T_LCLASS)
				case T_ID_MARK:
					l.emit(T_LID)
				case T_CONDITION_MARK:
					l.emit(T_LCONDITION)
				case T_PSEUDO_MARK:
					l.emit(T_LPSEUDO)
				case T_ATTRIBUTE_MARK: // 属性名称
					l.emit(T_LATTRIBUTE)
				case T_ATTRIBUTE_EQUAL: // 属性值
					l.emit(T_LITEM)
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
			l.acceptRun(alphavalues + numbers + classMarker + idMarker + conditionMaker + splitMaker + nameSeparator)
			l.emit(T_CMDPART)
		} else if r == "," {
			l.emit(T_COMMA)
		} else {
			return l.errorf("Illegal character '%s'.", r)
		}
	}
	l.emit(T_RPAREN)
	return initialLexerState
}

func commandState(l *lexer) lexerState {
	var lvl = 0
	var hascompart = false
	for r := l.next(); r != "}" || lvl > 0; r = l.next() {
		if r == " " || r == "\t" || r == "\r" {
			if hascompart && lvl == 0 {
				l.backup()
				l.emit(T_CMDPART)
				l.next()
			}
			hascompart = false
			l.ignore()
		} else if r == "\"" {
			l.ignore()
			l.rejectRun("\"")
			l.emit(T_CMDPART)
			l.next()
		} else if r == "\n" {
			if hascompart && lvl == 0 {
				l.backup()
				l.emit(T_CMDPART)
				l.next()
			}
			hascompart = false
			l.lineno += 1
			l.emit(T_SEMI)
		} else if r == EOF {
			return l.errorf("Unclosed statement...")
		} else if r == "{" {
			lvl++
		} else if r == "}" {
			lvl--
		} else {
			l.acceptRun(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*_+=-|\]['":/?.>,<`)
			hascompart = true
		}
	}
	if hascompart {
		l.backup()
		l.emit(T_CMDPART)
		l.next()
	}
	l.emit(T_RCBRAC)
	return initialLexerState
}
