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
	ret := xmlStartClass()

	if err := e.eatKeyword(CLASS, &ret); err != nil {
		return "", err
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", err
	}

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		return "", err
	}

	nextToken := e.tknzr.ReadAheadNextToken()
	for nextToken.TokenType == KEYWORD && nextToken.Token == "static" {
		x, _ := e.CompileClassVarDec()
		ret += x
		nextToken = e.tknzr.ReadAheadNextToken()
	}

	e.eatSymbol(RBRACE, &ret)
	return ret + xmlEndClass(), nil
}

func (e Engine) CompileClassVarDec() (string, error) {
	ret := xmlStartClassVarDec()

	if err := e.eatStaticOrField(&ret); err != nil {
		return "", fmt.Errorf("expected KEYWORD static or field: %w", err)
	}

	if err := e.eatType(&ret); err != nil {
		return "", err
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", err
	}

	nextToken := e.tknzr.ReadAheadNextToken()
	for (nextToken.TokenType == SYMBOL) && (nextToken.Token == ",") {
		e.eatSymbol(KOMMA, &ret)
		if err := e.eatIdentifier(&ret); err != nil {
			return "", err
		}
		nextToken = e.tknzr.ReadAheadNextToken()
	}

	e.eatSymbol(SEMICOLON, &ret)
	return ret + xmlEndClassVarDec(), nil
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

func (e Engine) eatStaticOrField(ret *string) error {
	e.tknzr.Advance()
	if (e.tknzr.TokenType() != KEYWORD) || (e.tknzr.Keyword() != STATIC && e.tknzr.Keyword() != FIELD) {
		return fmt.Errorf("unexpected token: %v", e.tknzr.curToken)
	}

	*ret += xmlKeyword(e.tknzr.Keyword())
	return nil
}
func (e Engine) eatKeyword(expectedKeyword Keyword, ret *string) error {
	e.tknzr.Advance()
	if (e.tknzr.TokenType() != KEYWORD) || (e.tknzr.Keyword() != expectedKeyword) {
		return fmt.Errorf("expected KEYWORD %s", expectedKeyword)
	}

	*ret += xmlKeyword(expectedKeyword)
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
	return "<class>\n"
}

func xmlEndClass() string {
	return "</class>\n"
}

func xmlStartClassVarDec() string {
	return "<classVarDec>\n"
}

func xmlEndClassVarDec() string {
	return "</classVarDec>\n"
}
