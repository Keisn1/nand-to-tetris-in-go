package compiler

import (
	"errors"
	"fmt"
	"log"
	"os"
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

		if !t.eof() {
			t.readBack()
		}

		tokenS := t.input[t.currentPos:t.readPos]
		t.curToken = NewToken(removeWhiteSpaces(tokenS), INT_CONST)

	case t.isLetterOrUnderScore():
		for t.isWordCharacter() && !t.isWhiteSpace() && !t.eof() {
			t.readChar()
		}

		if !t.eof() {
			t.readBack()
		}

		tokenS := t.input[t.currentPos:t.readPos]
		tokenS = removeWhiteSpaces(tokenS)
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

func (t *Tokenizer) OutputXML(out string) {
	file, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write([]byte("<tokens>\n"))
	for t.HasMoreTokens() {
		t.Advance()

		gotTokenType := t.TokenType()

		switch gotTokenType {
		case KEYWORD:
			gotKeyword := t.Keyword()
			file.Write([]byte("<keyword>"))
			file.Write([]byte(gotKeyword))
			file.Write([]byte("</keyword>\n"))
		case SYMBOL:
			gotSymbol := t.Symbol()
			file.Write([]byte("<symbol>"))
			switch gotSymbol {
			case GT:
				file.Write([]byte("&gt;"))
			case LT:
				file.Write([]byte("&lt;"))
			case AND:
				file.Write([]byte("&amp;"))
			default:
				file.Write([]byte(string(gotSymbol)))
			}
			file.Write([]byte("</symbol>\n"))

		case IDENTIFIER:
			gotIdent := t.Identifier()
			file.Write([]byte("<identifier>"))
			file.Write([]byte(gotIdent))
			file.Write([]byte("</identifier>\n"))
		case INT_CONST:
			gotIntVal := t.IntVal()
			file.Write([]byte("<integerConstant>"))
			fmt.Fprint(file, gotIntVal)
			file.Write([]byte("</integerConstant>\n"))
		case STRING_CONST:
			gotStrVal := t.StringVal()
			file.Write([]byte("<stringConstant>"))
			file.Write([]byte(gotStrVal))
			file.Write([]byte("</stringConstant>\n"))
		}
	}
	file.Write([]byte("</tokens>"))
}

func (t *Tokenizer) readBack() {
	t.readPos--
	t.ch = t.input[t.readPos-1]
}

func (t Tokenizer) isWhiteSpace() bool {
	validRegex := regexp.MustCompile(`\s`)
	return validRegex.MatchString(string(t.ch))
}

func (t Tokenizer) HasMoreTokens() bool {
	if t.readPos >= len(t.input) {
		return false
	}
	return len(strings.TrimSpace(t.input[t.readPos:])) > 0
}

func (t Tokenizer) TokenType() TokenType {
	return t.curToken.tokenType
}

func (t *Tokenizer) readCharNonWhiteSpace() {
	if !t.eof() {
		t.readChar()
	}
	t.jumpWhiteSpace()
	if t.ch == '/' {
		if t.readPos < len(t.input) {
			if t.input[t.readPos] == '/' || t.input[t.readPos] == '*' {
				t.jumpComment()
			}
		}
	}
}

func (t *Tokenizer) jumpComment() {
	for t.ch != '\n' && !t.eof() {
		t.currentPos = t.readPos
		t.readChar()
	}
}

func (t *Tokenizer) jumpWhiteSpace() {
	for t.isWhiteSpace() && !t.eof() {
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
	if err != nil {
		log.Fatalf("intVal: %v", err)
	}
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
