package compiler

import (
	"fmt"
)

type Engine struct {
	Tknzr  *Tokenizer
	Errors []error
}

func NewEngine(tknzr *Tokenizer) Engine {
	return Engine{Tknzr: tknzr}
}

func (e *Engine) CompileClass() string {
	ret := xmlStartClass()

	if err := e.eatKeyword(CLASS, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	e.Tknzr.Advance()
	for isClassVarDec(e.Tknzr.Keyword()) {
		ret += e.CompileClassVarDec()
	}

	for isSubRoutineDec(e.Tknzr.Keyword()) {
		ret += e.CompileSubroutineDec()
	}

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	e.Tknzr.Advance()
	return ret + xmlEndClass()
}

func (e *Engine) CompileClassVarDec() string {
	ret := xmlStartClassVarDec()

	if err := e.eatStaticOrField(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	e.Tknzr.Advance()
	for e.Tknzr.Symbol() == KOMMA {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
		}

		e.Tknzr.Advance()
		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
		}

		e.Tknzr.Advance()
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	e.Tknzr.Advance()
	return ret + xmlEndClassVarDec()
}

func (e *Engine) CompileSubroutineDec() string {
	ret := xmlStartSubroutineDec()

	if err := e.eatSubRoutineDecStart(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatVoidOrType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatSymbol(LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	e.Tknzr.Advance()
	switch e.Tknzr.Symbol() {
	case RPAREN:
		ret += xmlStartParameterList()
		ret += xmlEndParameterList()
	default:
		ret += e.CompileParameterList()
	}

	if err := e.eatSymbol(RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	e.Tknzr.Advance()
	ret += e.CompileSubroutineBody()

	return ret + xmlEndSubroutineDec()
}

func (e *Engine) CompileParameterList() string {
	ret := xmlStartParameterList()

	if err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	e.Tknzr.Advance()
	for e.Tknzr.Symbol() == KOMMA {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}

		e.Tknzr.Advance()
		if err := e.eatType(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}

		e.Tknzr.Advance()
		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}

		e.Tknzr.Advance()
	}

	return ret + xmlEndParameterList()
}

func (e *Engine) CompileSubroutineBody() string {
	ret := xmlStartSubroutineBody()

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	e.Tknzr.Advance()
	for e.Tknzr.Keyword() == VAR {
		varDec, _ := e.CompileVarDec()
		ret += varDec
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		return ""
	}

	e.Tknzr.Advance()
	return ret + xmlEndSubroutineBody()
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

	e.Tknzr.Advance()
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

	e.Tknzr.Advance()
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
			return fmt.Errorf("%w: %w",
				NewErrSyntaxUnexpectedToken(fmt.Sprintf("expected KEYWORD %s or type", VOID), e.Tknzr.curToken.Literal),
				err,
			)
		}
	}
	return nil
}

func (e *Engine) eatType(ret *string) error {
	switch e.Tknzr.TokenType() {
	case IDENTIFIER:
		*ret += xmlIdentifier(e.Tknzr.Identifier())
	case KEYWORD:
		if e.Tknzr.Keyword() != INT && e.Tknzr.Keyword() != CHAR && e.Tknzr.Keyword() != BOOLEAN {
			return NewErrSyntaxUnexpectedToken(
				fmt.Sprintf("KEYWORD %s / %s / %s or className", INT, CHAR, BOOLEAN),
				e.Tknzr.Keyword())
		}
		*ret += xmlKeyword(e.Tknzr.Keyword())
	default:
		return NewErrSyntaxUnexpectedToken(
			fmt.Sprintf("KEYWORD %s / %s / %s or className", INT, CHAR, BOOLEAN),
			e.Tknzr.curToken.Literal)
	}
	return nil
}

func (e Engine) eatIdentifier(ret *string) error {
	if e.Tknzr.TokenType() != IDENTIFIER {
		return NewErrSyntaxUnexpectedTokenType(IDENTIFIER, e.Tknzr.curToken.Literal)
	}
	*ret += xmlIdentifier(e.Tknzr.Identifier())
	return nil
}

func (e Engine) eatSubRoutineDecStart(ret *string) error {
	if !isSubRoutineDec(e.Tknzr.Keyword()) {
		return NewErrSyntaxUnexpectedToken(
			fmt.Sprintf("KEYWORD %s / %s / %s ", CONSTRUCTOR, FUNCTION, METHOD),
			e.Tknzr.Keyword(),
		)
	}

	*ret += xmlKeyword(e.Tknzr.Keyword())
	return nil
}

func (e Engine) eatStaticOrField(ret *string) error {
	if !isClassVarDec(e.Tknzr.Keyword()) {
		return NewErrSyntaxUnexpectedToken(fmt.Sprintf("KEYWORD %s or %s ", STATIC, FIELD), e.Tknzr.Keyword())
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

func (e Engine) eatKeyword(expectedKeyword string, ret *string) error {
	if (e.Tknzr.TokenType() != KEYWORD) || (e.Tknzr.Keyword() != expectedKeyword) {
		return NewErrSyntaxUnexpectedToken(expectedKeyword, e.Tknzr.curToken.Literal)
	}

	*ret += xmlKeyword(expectedKeyword)
	return nil
}

func (e Engine) eatSymbol(expectedSymbol string, ret *string) error {
	if (e.Tknzr.TokenType() != SYMBOL) || (e.Tknzr.Symbol() != expectedSymbol) {
		return NewErrSyntaxUnexpectedToken(expectedSymbol, e.Tknzr.curToken.Literal)
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
