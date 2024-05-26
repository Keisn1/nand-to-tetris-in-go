package compiler

import (
	"errors"
	"fmt"
)

type Engine struct {
	tknzr *Tokenizer
}

func NewEngine(tknzr *Tokenizer) Engine {
	return Engine{tknzr: tknzr}
}

func (e Engine) CompileClass() (string, error) {
	ret := xmlStartClass()

	if err := e.eatKeyword(CLASS, &ret); err != nil {
		return "", fmt.Errorf("compileClass: %w", err)
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileClass: %w", err)
	}

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		return "", fmt.Errorf("compileClass: %w", err)
	}

	for isClassVarDec(Keyword(e.tknzr.NextToken().Literal)) {
		classVarDec, err := e.CompileClassVarDec()
		if err != nil {
			return "", fmt.Errorf("compileClass: %w", err)
		}
		ret += classVarDec
	}

	for isSubRoutineDec(Keyword(e.tknzr.NextToken().Literal)) {
		subRoutineDec, err := e.CompileSubroutineDec()
		if err != nil {
			return "", fmt.Errorf("compileClass: %w", err)
		}
		ret += subRoutineDec
	}

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

	if err := e.eatType(&ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileClassVarDec: %w", err)
	}

	for e.tknzr.NextToken().Literal == "," {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			return "", fmt.Errorf("compileClassVarDec: %w", err)
		}
		if err := e.eatIdentifier(&ret); err != nil {
			return "", fmt.Errorf("compileClassVarDec: %w", err)
		}
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

	if err := e.eatVoidOrType(&ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	if err := e.eatIdentifier(&ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	if err := e.eatSymbol(LPAREN, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	parameterList, err := e.CompileParameterList()
	if err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}
	ret += parameterList

	if err := e.eatSymbol(RPAREN, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}

	subroutineBody, err := e.CompileSubroutineBody()
	if err != nil {
		return "", fmt.Errorf("compileSubroutineDec: %w", err)
	}
	ret += subroutineBody

	return ret + xmlEndSubroutineDec(), nil
}

func (e Engine) eatParameter(ret *string) error {
	if err := e.eatType(ret); err != nil {
		return fmt.Errorf("eatParameter: %w", err)
	}

	if err := e.eatIdentifier(ret); err != nil {
		return fmt.Errorf("eatParameter: %w", err)
	}

	return nil
}

func (e Engine) CompileParameterList() (string, error) {
	ret := xmlStartParameterList()

	if e.tknzr.NextToken().Literal != string(RPAREN) {
		if err := e.eatParameter(&ret); err != nil {
			return "", fmt.Errorf("compileParameterList: %w", err)
		}
		for e.tknzr.NextToken().Literal == string(KOMMA) {
			if err := e.eatSymbol(KOMMA, &ret); err != nil {
				return "", fmt.Errorf("compileParameterList: %w", err)
			}
			if err := e.eatParameter(&ret); err != nil {
				return "", fmt.Errorf("compileParameterList: %w", err)
			}
		}
	}

	return ret + xmlEndParameterList(), nil
}

func (e Engine) CompileSubroutineBody() (string, error) {
	ret := xmlStartSubroutineBody()

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		return "", fmt.Errorf("compileSubroutineBody: %w", err)
	}

	for isVarDec(Keyword(e.tknzr.NextToken().Literal)) {
		varDec := e.CompileVarDec()
		ret += varDec
	}

	ret += e.CompileStatements()
	e.eatSymbol(RBRACE, &ret)

	return ret + xmlEndSubroutineBody(), nil
}

func (e Engine) CompileVarDec() string {
	ret := xmlStartVarDec()
	e.eatKeyword(VAR, &ret)
	e.eatIdentifier(&ret)
	e.eatIdentifier(&ret)
	e.eatSymbol(SEMICOLON, &ret)
	return ret + xmlEndVarDec()
}

func (e Engine) CompileStatements() string {
	ret := xmlStartStatements()
	ret += e.CompileReturn()
	return ret + xmlEndStatements()
}

func (e Engine) CompileReturn() string {
	ret := xmlStartReturnStatement()
	e.eatKeyword(RETURN, &ret)
	e.eatSymbol(SEMICOLON, &ret)
	return ret + xmlEndReturnStatement()

}

func (e Engine) eatVoidOrType(ret *string) error {
	switch Keyword(e.tknzr.NextToken().Literal) {
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
	e.tknzr.Advance()
	switch e.tknzr.TokenType() {
	case IDENTIFIER:
		*ret += xmlIdentifier(e.tknzr.Identifier())
	case KEYWORD:
		if e.tknzr.Keyword() != INT && e.tknzr.Keyword() != CHAR && e.tknzr.Keyword() != BOOLEAN {
			return fmt.Errorf("expected className or KEYWORD int / char / boolean, got: %s", e.tknzr.Keyword())
		}
		*ret += xmlKeyword(e.tknzr.Keyword())
	default:
		return fmt.Errorf("expected className or KEYWORD int / char / boolean, got: %s", e.tknzr.Keyword())
	}
	return nil
}

func (e Engine) eatIdentifier(ret *string) error {
	e.tknzr.Advance()
	if e.tknzr.TokenType() != IDENTIFIER {
		return errors.New("expected tokenType IDENTIFIER")
	}
	*ret += xmlIdentifier(e.tknzr.Identifier())
	return nil
}

func (e Engine) eatSubRoutineDecStart(ret *string) error {
	e.tknzr.Advance()
	if !isSubRoutineDec(e.tknzr.Keyword()) {
		return fmt.Errorf("unexpected token: %v : expected KEYWORD constructor / function / method", e.tknzr.curToken)
	}

	*ret += xmlKeyword(e.tknzr.Keyword())
	return nil
}

func (e Engine) eatStaticOrField(ret *string) error {
	e.tknzr.Advance()
	if !isClassVarDec(e.tknzr.Keyword()) {
		return fmt.Errorf("unexpected token: %v : expected KEYWORD static or field", e.tknzr.curToken)
	}

	*ret += xmlKeyword(e.tknzr.Keyword())
	return nil
}

func isSubRoutineDec(kw Keyword) bool {
	return kw == CONSTRUCTOR || kw == FUNCTION || kw == METHOD
}

func isClassVarDec(kw Keyword) bool {
	return kw == STATIC || kw == FIELD
}

func isVarDec(kw Keyword) bool {
	return kw == VAR
}

func (e Engine) eatKeyword(expectedKeyword Keyword, ret *string) error {
	e.tknzr.Advance()
	if (e.tknzr.TokenType() != KEYWORD) || (e.tknzr.Keyword() != expectedKeyword) {
		return fmt.Errorf("expected KEYWORD %s", expectedKeyword)
	}

	*ret += xmlKeyword(expectedKeyword)
	return nil
}

func (e Engine) eatSymbol(expectedSymbol rune, ret *string) error {
	e.tknzr.Advance()
	if (e.tknzr.TokenType() != SYMBOL) || (e.tknzr.Symbol() != expectedSymbol) {
		return fmt.Errorf("expected SYMBOL %c", expectedSymbol)
	}

	*ret += xmlSymbol(expectedSymbol)
	return nil
}

func xmlSymbol(symbol rune) string {
	return fmt.Sprintf("<symbol> %s </symbol>\n", string(symbol))
}

func xmlKeyword(kw Keyword) string {
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
