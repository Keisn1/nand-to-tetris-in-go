package compiler

type Engine struct {
	tknzr Tokenizer
}

func NewEngine(tknzr Tokenizer) Engine {
	return Engine{tknzr: tknzr}
}

func (e Engine) CompileClass() string {
	return "<class><keyword>class</keyword><identifier>Main</identifier><symbol>{</symbol><symbol>}</symbol></class>"
}
