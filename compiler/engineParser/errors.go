package engineParser

import (
	"fmt"
	"hack/compiler/token"
)

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
	return fmt.Sprintf("expected %s, found %q", err.ExpectedToken, err.FoundToken)
}

func NewErrSyntaxNotATerm(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf(
			"TOKENTYPE '%s' / '%s' / '%s' / '%s' or SYMBOL '%s' / '%s' / '%s'",
			token.INT_CONST, token.STRING_CONST, token.KEYWORD, token.IDENTIFIER, token.LPAREN, token.TILDE, token.MINUS,
		),
		foundToken,
	)
}

func NewErrSyntaxNotAKeywordConst(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s / %s / %s / %s", token.TRUE, token.FALSE, token.NULL, token.THIS),
		foundToken,
	)
}

func NewErrSyntaxNotAType(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s / %s / %s or className", token.INT, token.CHAR, token.BOOLEAN),
		foundToken,
	)
}

func NewErrSyntaxNotAClassVarDec(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s or %s ", token.STATIC, token.FIELD),
		foundToken,
	)
}

func NewErrSyntaxNotASubroutineDec(foundToken string) ErrSyntaxUnexpectedToken {
	return NewErrSyntaxUnexpectedToken(
		fmt.Sprintf("KEYWORD %s / %s / %s ", token.CONSTRUCTOR, token.FUNCTION, token.METHOD),
		foundToken,
	)
}

type ErrSyntaxUnexpectedTokenType struct {
	ExpectedTokenType token.TokenType
	FoundToken        string
}

func NewErrSyntaxUnexpectedTokenType(expectedTokenType token.TokenType, foundToken string) ErrSyntaxUnexpectedTokenType {
	return ErrSyntaxUnexpectedTokenType{
		ExpectedTokenType: expectedTokenType,
		FoundToken:        foundToken,
	}
}

func (err ErrSyntaxUnexpectedTokenType) Error() string {
	return fmt.Sprintf("expected tokenType %s but found tokenType %s", err.ExpectedTokenType, err.FoundToken)
}
