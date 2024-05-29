package compiler

import (
	"errors"
	"fmt"
)

var ErrEndOfFile = errors.New("no more tokens")

type ErrSyntaxUnexpectedToken struct {
	ExpectedToken string
	FoundToken    string
}

func NewErrSyntaxUnexpectedToken(expectedToken, foundToken string) ErrSyntaxUnexpectedToken {
	return ErrSyntaxUnexpectedToken{
		ExpectedToken: expectedToken,
		FoundToken:    foundToken,
	}
}

func (err ErrSyntaxUnexpectedToken) Error() string {
	return fmt.Sprintf("expected %s, found %s", err.ExpectedToken, err.FoundToken)
}

func NewErrSyntaxNotATerm(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf(
			"TOKENTYPE '%s' / '%s' / '%s' / '%s' or SYMBOL '%s' / '%s' / '%s'",
			INT_CONST, STRING_CONST, KEYWORD, IDENTIFIER, LPAREN, TILDE, MINUS,
		),
		foundToken,
	)
}

func NewErrSyntaxNotAKeywordConst(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s / %s / %s / %s", TRUE, FALSE, NULL, THIS),
		foundToken,
	)
}

func NewErrSyntaxNotAType(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s / %s / %s or className", INT, CHAR, BOOLEAN),
		foundToken,
	)
}

func NewErrSyntaxNotAClassVarDec(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s or %s ", STATIC, FIELD),
		foundToken,
	)
}

func NewErrSyntaxNotASubroutineDec(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s / %s / %s ", CONSTRUCTOR, FUNCTION, METHOD),
		foundToken,
	)
}

func NewErrSyntaxNotVoidOrType(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("expected KEYWORD %s or type", VOID),
		foundToken,
	)
}

type ErrSyntaxUnexpectedTokenType struct {
	ExpectedTokenType TokenType
	FoundToken        string
}

func NewErrSyntaxUnexpectedTokenType(expectedTokenType TokenType, foundToken string) ErrSyntaxUnexpectedTokenType {
	return ErrSyntaxUnexpectedTokenType{
		ExpectedTokenType: expectedTokenType,
		FoundToken:        foundToken,
	}
}

func (err ErrSyntaxUnexpectedTokenType) Error() string {
	return fmt.Sprintf("expected tokenType %s but found tokenType %s", err.ExpectedTokenType, err.FoundToken)
}
