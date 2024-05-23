package compiler

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TokenType string

const (
	KEYWORD      TokenType = "KEYWORD"
	SYMBOL       TokenType = "SYMBOL"
	IDENTIFIER   TokenType = "IDENTIFIER"
	INT_CONST    TokenType = "INT_CONST"
	STRING_CONST TokenType = "STRING_CONST"
)

type Token struct {
	token     string
	tokenType TokenType
}

func NewToken(token string, tokenType TokenType) Token {
	return Token{token: token, tokenType: tokenType}
}

type Tokenizer struct {
	input      string
	curToken   Token
	currentPos int
	readPos    int
	ch         byte
}

func NewTokenizer(input string) Tokenizer {
	return Tokenizer{input: input}
}

func (t *Tokenizer) Advance() error {
	if !t.HasMoreTokens() {
		return errors.New("advance should not be called if there are no more tokens")
	}

	t.readCharNonWhiteSpace()

	switch {
	case t.isSymbol():
		tokenS := t.input[t.currentPos:t.readPos]
		t.curToken = NewToken(tokenS, SYMBOL)

	case t.isInteger():
		for (t.isInteger() || t.isWhiteSpace()) && !t.eof() {
			t.readChar()
		}

		tokenS := t.input[t.currentPos:t.readPos]
		t.curToken = NewToken(removeWhiteSpaces(tokenS), INT_CONST)

	case t.isLetterOrUnderScore():
		for t.isWordCharacter() && !t.eof() {
			t.readChar()
		}

		tokenS := strings.ToLower(t.input[t.currentPos:t.readPos])
		if _, ok := keywords[tokenS]; ok {
			t.curToken = NewToken(tokenS, KEYWORD)
		} else {
			t.curToken = NewToken(tokenS, IDENTIFIER)
		}

	case t.isDoubleQuote():
		t.readChar()
		for !t.isDoubleQuote() && !t.eof() {
			t.readChar()
		}

		tokenS := t.input[t.currentPos+1 : t.readPos-1]
		t.curToken = NewToken(tokenS, STRING_CONST)
	}

	t.currentPos = t.readPos
	return nil
}

func (t Tokenizer) isWhiteSpace() bool {
	validRegex := regexp.MustCompile(`\s`)
	return validRegex.MatchString(string(t.ch))
}

func (t Tokenizer) HasMoreTokens() bool {
	return t.readPos < len(t.input)
}

func (t Tokenizer) TokenType() TokenType {
	return t.curToken.tokenType
}

func (t *Tokenizer) readCharNonWhiteSpace() {
	t.readChar()
	t.jumpWhiteSpace()
}

func (t *Tokenizer) jumpWhiteSpace() {
	for t.isWhiteSpace() {
		t.currentPos = t.readPos
		t.readChar()
	}
}

func (t *Tokenizer) readChar() {
	t.ch = t.input[t.readPos]
	t.readPos++
}

func (t Tokenizer) Keyword() Keyword {
	return keywords[t.curToken.token]
}

func (t Tokenizer) Symbol() rune {
	return rune(t.curToken.token[0])
}

func (t Tokenizer) IntVal() int {
	nbr, err := strconv.Atoi(t.curToken.token)
	fmt.Println(err)
	return nbr
}

func (t Tokenizer) StringVal() string {
	return t.curToken.token
}

func (t Tokenizer) Identifier() string {
	return t.curToken.token
}

func (t *Tokenizer) eof() bool {
	return t.readPos >= len(t.input)
}

func (t Tokenizer) isDoubleQuote() bool {
	return t.ch == '"'
}

func (t Tokenizer) isWordCharacter() bool {
	validRegex := regexp.MustCompile(`\w`)
	return validRegex.MatchString(string(t.ch))
}

func (t Tokenizer) isLetterOrUnderScore() bool {
	validRegex := regexp.MustCompile(`[a-zA-Z_]`)
	return validRegex.MatchString(string(t.ch))
}

func (t Tokenizer) isInteger() bool {
	return t.ch >= '0' && t.ch <= '9'
}

func (t Tokenizer) isSymbol() bool {
	if _, ok := allSymbols[t.ch]; ok {
		return true
	}
	return false
}
