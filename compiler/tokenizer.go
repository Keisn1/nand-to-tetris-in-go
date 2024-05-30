package compiler

import (
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
	EOF_CONST    TokenType = "EOF_CONST"
)

type Token struct {
	Literal   string
	TokenType TokenType
}

func NewToken(token string, tokenType TokenType) Token {
	return Token{Literal: token, TokenType: tokenType}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s", t.TokenType, t.Literal)
}

type Tokenizer struct {
	input      string
	curToken   Token
	currentPos int
	readPos    int
	ch         byte
}

func NewTokenizer(input string) Tokenizer {
	// re := regexp.MustCompile(`//.*?\n`)
	// output1 := re.ReplaceAllString(input, "")

	re := regexp.MustCompile(`//.*?$`)
	output1 := re.ReplaceAllString(input, "")

	re = regexp.MustCompile(`//.*?\n`)
	output2 := re.ReplaceAllString(output1, "")

	// re = regexp.MustCompile(`/\**?\*/`)
	// output3 := re.ReplaceAllString(output2, "")

	return Tokenizer{input: output2}
}

func (t *Tokenizer) Advance() error {

	if !t.HasMoreTokens() {
		t.curToken = NewToken(EOF, EOF_CONST)
		return ErrEndOfFile
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

func (t *Tokenizer) HasMoreTokens() bool {
	if t.readPosAtEOF() {
		return false
	}

	return len(strings.TrimSpace(t.input[t.readPos:])) > 0
}

func (t Tokenizer) TokenType() TokenType {
	return t.curToken.TokenType
}

func (t Tokenizer) Keyword() string {
	return keywords[t.curToken.Literal]
}

func (t Tokenizer) Symbol() string {
	return string(t.curToken.Literal[0])
}

func (t Tokenizer) IntVal() int {
	nbr, err := strconv.Atoi(t.curToken.Literal)
	if err != nil {
		log.Fatalf("intVal: %v", err)
	}
	return nbr
}

func (t Tokenizer) StringVal() string {
	return t.curToken.Literal
}

func (t Tokenizer) Identifier() string {
	return t.curToken.Literal
}

func (t *Tokenizer) readUpToNextToken() {
	if t.readPosAtEOF() {
		return
	}
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
	for t.ch != '\n' && !t.readPosAtEOF() {
		t.readChar()
	}
}

func (t *Tokenizer) jumpWhiteSpace() {
	for t.isWhiteSpace() && !t.readPosAtEOF() {
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
	for !t.isDoubleQuote() && !t.readPosAtEOF() {
		t.readChar()
	}

	return t.input[t.currentPos+1 : t.readPos-1]
}
func (t *Tokenizer) readInteger() string {
	for t.isInteger() && !t.readPosAtEOF() {
		t.readChar()
	}

	if !t.isInteger() {
		t.readBack()
	}

	if t.readPosAtEOF() {
		return removeWhiteSpaces(t.input[t.currentPos:t.readPos])
	}

	return t.input[t.currentPos:t.readPos]
}

func (t *Tokenizer) readWord() string {
	for t.isWordCharacter() && !t.readPosAtEOF() {
		t.readChar()
	}

	if !t.isWordCharacter() {
		t.readBack()
	}

	if t.readPosAtEOF() {
		return removeWhiteSpaces(t.input[t.currentPos:t.readPos])
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

func (t *Tokenizer) readPosAtEOF() bool {
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
	if _, ok := allSymbols[string(t.ch)]; ok {
		return true
	}
	return false
}

func (t Tokenizer) isKeyword() bool {
	if _, ok := keywords[t.input[t.currentPos:t.readPos]]; ok {
		return true
	}
	return false
}

func removeWhiteSpaces(input string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(input, "")
}
