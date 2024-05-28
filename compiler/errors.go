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

type ErrSyntaxNotAType struct {
	ErrSyntaxUnexpectedToken
}

func NewErrSyntaxNotAType(foundToken string) ErrSyntaxNotAType {
	return ErrSyntaxNotAType{
		ErrSyntaxUnexpectedToken: NewErrSyntaxUnexpectedToken(
			fmt.Sprintf("KEYWORD %s / %s / %s or className", INT, CHAR, BOOLEAN),
			foundToken,
		),
	}
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
