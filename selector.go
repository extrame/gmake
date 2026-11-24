package main

func Selector(desc string) *Item {
	// 选择器允许以 # 开头表示 ID，因此关闭行首 # 注释规则
	_, tokens := lex("Item", desc, false, itemLexerState)
	return ParseItem("", tokens)
}
