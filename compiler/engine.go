package compiler

import (
	"errors"
	"fmt"
)

type Engine struct {
	tknzr *Tokenizer
}

func NewEngine(tknzr *Tokenizer) Engine {
	return Engine{tknzr: tknzr}
}

func xmlSymbol(symbol rune) string {
	return fmt.Sprintf("<symbol> %s </symbol>\n", string(symbol))
}

func xmlIdentifier(identifier string) string {
	return fmt.Sprintf("<identifier> %s </identifier>\n", identifier)
}

func xmlStartClass() string {
	return `<class>
	<keyword>class</keyword>
`
}

func (e Engine) eatIdentifier(ret *string) error {
	e.tknzr.Advance()
	if e.tknzr.TokenType() != IDENTIFIER {
		return errors.New("expected tokenType IDENTIFIER")
	}
	*ret += xmlIdentifier(e.tknzr.Identifier())
	return nil
}

func (e Engine) eatSymbol(expectedSymbol rune, ret *string) error {
	e.tknzr.Advance()
	if (e.tknzr.TokenType() != SYMBOL) || (e.tknzr.Symbol() != expectedSymbol) {
		return fmt.Errorf("expected symbol %c", expectedSymbol)
	}

	*ret += xmlSymbol(expectedSymbol)
	return nil
}

func (e Engine) startClass(ret *string) error {
	if e.tknzr.TokenType() != KEYWORD || e.tknzr.Keyword() != CLASS {
		return errors.New("expected keyword CLASS")
	}
	*ret = xmlStartClass()
	return nil
}

func (e Engine) CompileClass() (string, error) {
	var ret string
	if err := e.startClass(&ret); err != nil {
		return "", err
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", err
	}

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		return "", err
	}

	for e.tknzr.HasMoreTokens() {
		e.tknzr.Advance()
		switch e.tknzr.TokenType() {
		case KEYWORD:
			switch e.tknzr.Keyword() {
			case STATIC:
				ret += e.compileClassVarDec()
			}
		}
	}

	end := `
<symbol>}</symbol>
</class>`
	return ret + end, nil
}

func (e Engine) compileClassVarDec() string {
	ret := `<classVarDec>
<keyword> static </keyword>
`

	for e.tknzr.TokenType() != SYMBOL && e.tknzr.Symbol() != SEMICOLON {
		e.tknzr.Advance()
		if e.tknzr.TokenType() == KEYWORD {
			switch e.tknzr.Keyword() {
			case BOOLEAN:
				ret += "<keyword> boolean </keyword>"
			case INT:
				ret += "<keyword> int </keyword>"
			}
		}
		if e.tknzr.TokenType() == IDENTIFIER {
			staticVar := `<identifier> %s </identifier>`
			ret += fmt.Sprintf(staticVar, e.tknzr.Identifier())
		}
	}

	end := `<symbol> ; </symbol></classVarDec>`
	return ret + end
}
