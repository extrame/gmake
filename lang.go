package main

// Token Types
const T_EOF string = "T_EOF"
const T_DIRECT string = "T_DIRECT"
const T_CMDPART string = "T_CMDPART"

const T_LPAREN string = "T_LPAREN"
const T_RPAREN string = "T_RPAREN"

const T_LCBRAC string = "T_LCBRAC"
const T_RCBRAC string = "T_RCBRAC"

const T_COMMA string = "T_COMMA"
const T_SEMI string = "T_SEMI"

const T_LITEM string = "T_LITEM"
const T_RITEM string = "T_RITEM"

const T_LCLASS string = "T_LCLASS"
const T_CLASS_MARK string = "T_CLASS_MARK"

const T_LID string = "T_LID"
const T_ID_MARK string = "T_ID_MARK"

const alphavalues = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
const numbers = `0123456789`
const classMarker = "."
const idMarker = "#"
const platformMaker = "@"
const splitMaker = "_"

// asts

type Command struct {
	Parts []string
}
