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
	FoundToken        string
}

func NewErrSyntaxUnexpectedTokenType(expectedTokenType TokenType, foundTokenType string) ErrSyntaxUnexpectedTokenType {
	return ErrSyntaxUnexpectedTokenType{
		ExpectedTokenType: expectedTokenType,
		FoundToken:        foundTokenType,
	}
}

func (err ErrSyntaxUnexpectedTokenType) Error() string {
	return fmt.Sprintf("expected tokenType %s but found tokenType %s", err.ExpectedTokenType, err.FoundToken)
}
