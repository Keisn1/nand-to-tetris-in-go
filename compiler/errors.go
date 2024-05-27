package compiler

import "fmt"

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

type ErrSyntaxUnexpectedTokenType struct {
	ExpectedTokenType TokenType
	FoundTokenType    TokenType
}

func NewErrSyntaxUnexpectedTokenType(expectedTokenType, foundTokenType TokenType) ErrSyntaxUnexpectedTokenType {
	return ErrSyntaxUnexpectedTokenType{
		ExpectedTokenType: expectedTokenType,
		FoundTokenType:    foundTokenType,
	}
}

func (err ErrSyntaxUnexpectedTokenType) Error() string {
	return fmt.Sprintf("expected tokenType %s but found tokenType %s", err.ExpectedTokenType, err.FoundTokenType)
}
