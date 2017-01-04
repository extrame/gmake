package main

func Selector(desc string) Item {
	_, tokens := Lexer("Item", desc, itemLexerState)
	return ParseItem("", tokens)
}
