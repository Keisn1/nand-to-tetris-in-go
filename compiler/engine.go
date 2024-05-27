package compiler

import (
	"errors"
	"fmt"
)

type Engine struct {
	Tknzr  *Tokenizer
	Errors []error
}

func NewEngine(tknzr *Tokenizer) Engine {
	return Engine{Tknzr: tknzr}
}

func (e *Engine) CompileClass() (string, error) {
	ret := xmlStartClass()

	if err := e.eatKeyword(CLASS, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
		return "", fmt.Errorf("compileClass: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileClass: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		return "", fmt.Errorf("compileClass: %w", err)
	}

	e.Tknzr.Advance()
	for isClassVarDec(e.Tknzr.Keyword()) {
		classVarDec, err := e.CompileClassVarDec()
		if err != nil {
			return "", fmt.Errorf("compileClass: %w", err)
		}
		ret += classVarDec
		e.Tknzr.Advance()
	}

	for isSubRoutineDec(e.Tknzr.Keyword()) {
		subRoutineDec, err := e.CompileSubroutineDec()
		if err != nil {
			return "", fmt.Errorf("compileClass: %w", err)
		}
		ret += subRoutineDec
	}

	e.Tknzr.Advance()
	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		return "", fmt.Errorf("compileClass: %w", err)
	}

	return ret + xmlEndClass(), nil
}

func (e Engine) CompileClassVarDec() (string, error) {
	ret := xmlStartClassVarDec()

	if err := e.eatStaticOrField(&ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatType(&ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	e.Tknzr.Advance()
	for e.Tknzr.Symbol() == KOMMA {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			return "", fmt.Errorf("compileClassVarDec: %w", err)
		}

		e.Tknzr.Advance()
		if err := e.eatIdentifier(&ret); err != nil {
			return "", fmt.Errorf("compileClassVarDec: %w", err)
		}

		e.Tknzr.Advance()
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	return ret + xmlEndClassVarDec(), nil
}

func (e Engine) CompileSubroutineDec() (string, error) {
	ret := xmlStartSubroutineDec()

	if err := e.eatSubRoutineDecStart(&ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatVoidOrType(&ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatSymbol(LPAREN, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	e.Tknzr.Advance()
	switch e.Tknzr.Symbol() {
	case RPAREN:
		ret += xmlStartParameterList()
		ret += xmlEndParameterList()
	default:
		parameterList, err := e.CompileParameterList()
		if err != nil {
			return "", fmt.Errorf("compileSubroutineDec: %w", err)
		}
		ret += parameterList
	}

	if err := e.eatSymbol(RPAREN, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	e.Tknzr.Advance()
	subroutineBody, err := e.CompileSubroutineBody()
	if err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}
	ret += subroutineBody

	return ret + xmlEndSubroutineDec(), nil
}

func (e Engine) CompileParameterList() (string, error) {
	ret := xmlStartParameterList()

	if err := e.eatType(&ret); err != nil {
		return "", fmt.Errorf("compileParameterList: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileParameterList: %w", err)
	}

	e.Tknzr.Advance()
	for e.Tknzr.Symbol() == KOMMA {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			return "", fmt.Errorf("compileParameterList: %w", err)
		}

		e.Tknzr.Advance()
		if err := e.eatType(&ret); err != nil {
			return "", fmt.Errorf("compileParameterList: %w", err)
		}

		e.Tknzr.Advance()
		if err := e.eatIdentifier(&ret); err != nil {
			return "", fmt.Errorf("compileParameterList: %w", err)
		}

		e.Tknzr.Advance()
	}

	return ret + xmlEndParameterList(), nil
}

func (e Engine) CompileSubroutineBody() (string, error) {
	ret := xmlStartSubroutineBody()

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineBody: %w", err)
	}

	e.Tknzr.Advance()
	for e.Tknzr.Keyword() == VAR {
		varDec, err := e.CompileVarDec()
		if err != nil {
			return "", fmt.Errorf("compileSubroutineBody: %w", err)
		}
		ret += varDec

		e.Tknzr.Advance()
	}

	ret += e.CompileStatements()

	e.Tknzr.Advance()
	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineBody: %w", err)
	}

	return ret + xmlEndSubroutineBody(), nil
}

func (e Engine) CompileVarDec() (string, error) {
	ret := xmlStartVarDec()

	if err := e.eatKeyword(VAR, &ret); err != nil {
		return "", fmt.Errorf("compileVarDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatType(&ret); err != nil {
		return "", fmt.Errorf("compileVarDec: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileVarDec: %w", err)
	}

	e.Tknzr.Advance()
	for e.Tknzr.Symbol() == KOMMA {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			return "", fmt.Errorf("compileClassVarDec: %w", err)
		}

		e.Tknzr.Advance()
		if err := e.eatIdentifier(&ret); err != nil {
			return "", fmt.Errorf("compileVarDec: %w", err)
		}

		e.Tknzr.Advance()
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		return "", fmt.Errorf("compileVarDec: %w", err)
	}

	return ret + xmlEndVarDec(), nil
}

func (e Engine) CompileStatements() string {
	ret := xmlStartStatements()
	ret += e.CompileReturn()
	return ret + xmlEndStatements()
}

func (e Engine) CompileReturn() string {
	ret := xmlStartReturnStatement()
	e.eatKeyword(RETURN, &ret)

	e.Tknzr.Advance()
	e.eatSymbol(SEMICOLON, &ret)

	return ret + xmlEndReturnStatement()

}

func (e Engine) eatParameter(ret *string) error {
	if err := e.eatType(ret); err != nil {
		return fmt.Errorf("eatParameter: %w", err)
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(ret); err != nil {
		return fmt.Errorf("eatParameter: %w", err)
	}

	return nil
}

func (e Engine) eatVoidOrType(ret *string) error {
	switch e.Tknzr.Keyword() {
	case VOID:
		e.eatKeyword(VOID, ret)
	default:
		if err := e.eatType(ret); err != nil {
			return fmt.Errorf("eatVoidOrType: expected type or KEYWORD void: %w", err)
		}
	}
	return nil
}

func (e Engine) eatType(ret *string) error {
	switch e.Tknzr.TokenType() {
	case IDENTIFIER:
		*ret += xmlIdentifier(e.Tknzr.Identifier())
	case KEYWORD:
		if e.Tknzr.Keyword() != INT && e.Tknzr.Keyword() != CHAR && e.Tknzr.Keyword() != BOOLEAN {
			return fmt.Errorf("expected className or KEYWORD int / char / boolean, got: %s", e.Tknzr.Keyword())
		}
		*ret += xmlKeyword(e.Tknzr.Keyword())
	default:
		return fmt.Errorf("expected className or KEYWORD int / char / boolean, got: %s", e.Tknzr.Keyword())
	}
	return nil
}

func (e Engine) eatIdentifier(ret *string) error {
	if e.Tknzr.TokenType() != IDENTIFIER {
		return errors.New("expected tokenType IDENTIFIER")
	}
	*ret += xmlIdentifier(e.Tknzr.Identifier())
	return nil
}

func (e Engine) eatSubRoutineDecStart(ret *string) error {
	if !isSubRoutineDec(e.Tknzr.Keyword()) {
		return fmt.Errorf("unexpected token: %v : expected KEYWORD constructor / function / method", e.Tknzr.curToken)
	}

	*ret += xmlKeyword(e.Tknzr.Keyword())
	return nil
}

func (e Engine) eatStaticOrField(ret *string) error {
	if !isClassVarDec(e.Tknzr.Keyword()) {
		return fmt.Errorf("unexpected token: %v : expected KEYWORD static or field", e.Tknzr.curToken)
	}

	*ret += xmlKeyword(e.Tknzr.Keyword())
	return nil
}

func isSubRoutineDec(kw string) bool {
	return kw == CONSTRUCTOR || kw == FUNCTION || kw == METHOD
}

func isClassVarDec(kw string) bool {
	return kw == STATIC || kw == FIELD
}

type ErrSyntaxUnexpectedToken struct {
	ExpectedToken Token
	FoundToken    Token
}

func NewErrSyntaxUnexpectedToken(expectedToken, foundToken Token) ErrSyntaxUnexpectedToken {
	return ErrSyntaxUnexpectedToken{
		ExpectedToken: expectedToken,
		FoundToken:    foundToken,
	}
}

func (err ErrSyntaxUnexpectedToken) Error() string {
	return fmt.Sprintf("expected %s but found token %s", err.ExpectedToken, err.FoundToken)
}

func (e Engine) eatKeyword(expectedKeyword string, ret *string) error {
	if (e.Tknzr.TokenType() != KEYWORD) || (e.Tknzr.Keyword() != expectedKeyword) {
		// return ErrSyntaxUnexpectedToken{}
		expectedToken := NewToken(string(expectedKeyword), KEYWORD)
		return NewErrSyntaxUnexpectedToken(expectedToken, e.Tknzr.curToken)
	}

	*ret += xmlKeyword(expectedKeyword)
	return nil
}

func (e Engine) eatSymbol(expectedSymbol string, ret *string) error {
	if (e.Tknzr.TokenType() != SYMBOL) || (e.Tknzr.Symbol() != expectedSymbol) {
		return fmt.Errorf("expected SYMBOL %s", expectedSymbol)
	}

	*ret += xmlSymbol(expectedSymbol)
	return nil
}

func xmlSymbol(symbol string) string {
	return fmt.Sprintf("<symbol> %s </symbol>\n", string(symbol))
}

func xmlKeyword(kw string) string {
	return fmt.Sprintf("<keyword> %s </keyword>\n", kw)
}

func xmlIdentifier(identifier string) string {
	return fmt.Sprintf("<identifier> %s </identifier>\n", identifier)
}

func xmlStartClass() string {
	return "<class>\n"
}

func xmlEndClass() string {
	return "</class>\n"
}

func xmlStartClassVarDec() string {
	return "<classVarDec>\n"
}

func xmlEndClassVarDec() string {
	return "</classVarDec>\n"
}

func xmlStartSubroutineDec() string {
	return "<subroutineDec>\n"
}

func xmlEndSubroutineDec() string {
	return "</subroutineDec>\n"
}

func xmlStartParameterList() string {
	return "<parameterList>\n"
}
func xmlEndParameterList() string {
	return "</parameterList>\n"
}

func xmlStartSubroutineBody() string {
	return "<subroutineBody>\n"
}

func xmlEndSubroutineBody() string {
	return "</subroutineBody>\n"
}

func xmlStartVarDec() string {
	return "<varDec>\n"
}

func xmlEndVarDec() string {
	return "</varDec>\n"
}

func xmlStartStatements() string {
	return "<statements>\n"
}

func xmlEndStatements() string {
	return "</statements>\n"
}

func xmlStartReturnStatement() string {
	return "<returnStatement>\n"
}

func xmlEndReturnStatement() string {
	return "</returnStatement>\n"
}
