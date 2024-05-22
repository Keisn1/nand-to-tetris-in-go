package compiler

import (
	"bytes"
	"errors"
	"regexp"
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

type Tokenizer struct {
	input        string
	curTokenType TokenType
	curToken     string
	currentPos   int
	readPos      int
	ch           byte
}

func NewTokenizer(input string) Tokenizer {
	return Tokenizer{input: input}
}

func (t Tokenizer) HasMoreTokens() bool {
	return t.currentPos < len(t.input)
}

func (t Tokenizer) TokenType() TokenType {
	return t.curTokenType
}

func (t *Tokenizer) Advance() error {
	if !t.HasMoreTokens() {
		return errors.New("advance should not be called if there are no more tokens")
	}

	t.readChar()
	switch {
	case t.isSymbol():
		t.curTokenType = SYMBOL
		t.curToken = strings.ToLower(t.input[t.currentPos:t.readPos])
	case t.isInteger():
		t.curTokenType = INT_CONST
		t.curToken = strings.ToLower(t.input[t.currentPos:t.readPos])
	case t.isLetterOrUnderScore():
		for t.isWordCharacter() && !t.eof() {
			t.readChar()
		}

		t.curToken = strings.ToLower(t.input[t.currentPos:t.readPos])
		if _, ok := keywords[t.curToken]; ok {
			t.curTokenType = KEYWORD
		} else {
			t.curTokenType = IDENTIFIER
		}
	}

	t.currentPos = t.readPos
	return nil
}

func (t *Tokenizer) readChar() {
	t.ch = t.input[t.readPos]
	t.readPos++
}

func (t Tokenizer) Symbol() string {
	return t.curToken
}

func (t *Tokenizer) eof() bool {
	return t.readPos >= len(t.input)
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
	return bytes.Contains([]byte{'{', '}'}, []byte{t.ch})
}
