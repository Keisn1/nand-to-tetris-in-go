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

	for e.tknzr.TokenType() != SYMBOL || e.tknzr.Symbol() != RBRACE {
		e.tknzr.Advance()
		switch e.tknzr.TokenType() {
		case KEYWORD:
			switch e.tknzr.Keyword() {
			case STATIC:
				x, _ := e.CompileClassVarDec()
				ret += x
			}
		}
	}

	ret += xmlSymbol(RBRACE)
	ret += xmlEndClass()
	return ret, nil
}

func (e Engine) CompileClassVarDec() (string, error) {
	var ret string

	if err := e.startClassVarDec(&ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	if err := e.eatType(&ret); err != nil {
		return "", err
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", err
	}

	// e.tknzr.Advance()
	// fmt.Println(e.tknzr.curToken)
	// e.tknzr.Advance()
	// fmt.Println(e.tknzr.curToken)
	// e.tknzr.Advance()
	// fmt.Println(e.tknzr.curToken)
	// for (e.tknzr.TokenType() != SYMBOL) || (e.tknzr.Symbol() != SEMICOLON) {
	// 	ret += xmlSymbol(KOMMA)
	// 	if err := e.eatIdentifier(&ret); err != nil {
	// 		return "", err
	// 	}

	// 	e.tknzr.Advance()
	// }

	e.eatSymbol(SEMICOLON, &ret)

	ret += xmlEndClassVarDec()
	return ret, nil
}

func (e Engine) eatType(ret *string) error {
	e.tknzr.Advance()
	*ret += xmlKeyword(e.tknzr.Keyword())
	return nil
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

func (e Engine) startClassVarDec(ret *string) error {
	if e.tknzr.TokenType() != KEYWORD || e.tknzr.Keyword() != STATIC {
		return errors.New("expected keyword STATIC or FIELD")
	}
	*ret = xmlStartClassVarDec()
	return nil
}

func xmlSymbol(symbol rune) string {
	return fmt.Sprintf("<symbol> %s </symbol>\n", string(symbol))
}

func xmlKeyword(kw Keyword) string {
	return fmt.Sprintf("<keyword> %s </keyword>\n", kw)
}

func xmlIdentifier(identifier string) string {
	return fmt.Sprintf("<identifier> %s </identifier>\n", identifier)
}

func xmlStartClass() string {
	return `<class>
	<keyword>class</keyword>
`
}

func xmlEndClass() string {
	return "</class>\n"
}

func xmlStartClassVarDec() string {
	return `<classVarDec>
	<keyword>static</keyword>
`
}

func xmlEndClassVarDec() string {
	return "</classVarDec>\n"
}
