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

	t.readUpToNextToken()

	switch {
	case t.isSymbol():
		t.curToken = NewToken(t.input[t.currentPos:t.readPos], SYMBOL)

	case t.isInteger():
		tokenS := t.readInteger()
		t.curToken = NewToken(removeWhiteSpaces(tokenS), INT_CONST)

	case t.isLetterOrUnderScore():
		tokenS := t.readWord()
		switch _, ok := keywords[tokenS]; ok {
		case true:
			t.curToken = NewToken(tokenS, KEYWORD)
		case false:
			t.curToken = NewToken(tokenS, IDENTIFIER)
		}

	case t.isDoubleQuote():
		tokenS := t.readStrConst()
		t.curToken = NewToken(tokenS, STRING_CONST)
	}

	t.currentPos = t.readPos
	return nil
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

func (t *Tokenizer) readUpToNextToken() {
	t.readChar()
	for t.isWhiteSpace() || t.ch == '/' {
		if t.isWhiteSpace() {
			t.jumpWhiteSpace()
		}
		if t.ch == '/' {
			if t.readPos < len(t.input) {
				switch t.input[t.readPos] {
				case '/':
					t.nextLine()
				case '*':
					t.jumpEndOfComment()
				default:
					return
				}
			}
		}
	}
}

func (t *Tokenizer) jumpEndOfComment() {
	t.readPos++
	for t.input[t.readPos:t.readPos+2] != "*/" {
		t.readPos++
	}
	t.readPos += 2
	t.currentPos = t.readPos
	t.readChar()
}

func (t *Tokenizer) nextLine() {
	for t.ch != '\n' && !t.isEOF() {
		t.readChar()
	}
}

func (t *Tokenizer) jumpWhiteSpace() {
	for t.isWhiteSpace() && !t.isEOF() {
		t.currentPos = t.readPos
		t.readChar()
	}
}

func (t *Tokenizer) readChar() {
	t.ch = t.input[t.readPos]
	t.readPos++
}

func (t *Tokenizer) readBack() {
	t.readPos--
	t.ch = t.input[t.readPos-1]
}

func (t *Tokenizer) readStrConst() string {
	t.readChar()
	for !t.isDoubleQuote() && !t.isEOF() {
		t.readChar()
	}

	return t.input[t.currentPos+1 : t.readPos-1]
}
func (t *Tokenizer) readInteger() string {
	for (t.isInteger() || t.isWhiteSpace()) && !t.isEOF() {
		t.readChar()
	}

	if !t.isEOF() {
		t.readBack()
	}
	return t.input[t.currentPos:t.readPos]
}
func (t *Tokenizer) readWord() string {
	for t.isWordCharacter() && !t.isWhiteSpace() && !t.isEOF() {
		t.readChar()
	}

	if !t.isEOF() {
		t.readBack()
	}

	tokenS := t.input[t.currentPos:t.readPos]
	return removeWhiteSpaces(tokenS)
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
			t.writeToken(file, string(t.Keyword()), "keyword")
		case SYMBOL:
			xmlSymbol := xmlSymbols[t.Symbol()]
			t.writeToken(file, xmlSymbol, "symbol")
		case IDENTIFIER:
			t.writeToken(file, t.Identifier(), "identifier")
		case INT_CONST:
			t.writeToken(file, fmt.Sprint(t.IntVal()), "integerConstant")
		case STRING_CONST:
			t.writeToken(file, t.StringVal(), "stringConstant")
		}
	}
	file.Write([]byte("</tokens>"))
}
func (t Tokenizer) writeToken(file *os.File, val string, tt string) {
	fmt.Fprintf(file, "<%s>", tt)
	fmt.Fprintf(file, "%s", val)
	fmt.Fprintf(file, "</%s>\n", tt)
}

func (t *Tokenizer) isEOF() bool {
	return t.readPos >= len(t.input)
}

func (t Tokenizer) isWhiteSpace() bool {
	validRegex := regexp.MustCompile(`\s`)
	return validRegex.MatchString(string(t.ch))
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

func removeWhiteSpaces(input string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(input, "")
}
