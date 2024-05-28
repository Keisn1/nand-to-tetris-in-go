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
	ret := xmlStart(CLASS_T)

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

	return ret + xmlEnd(CLASS_T)
}

func (e *Engine) CompileClassVarDec() string {
	ret := xmlStart(CLASSVARDEC_T)

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

	return ret + xmlEnd(CLASSVARDEC_T)
}

func (e *Engine) CompileSubroutineDec() string {
	ret := xmlStart(SUBROUTINEDEC_T)

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
		ret += xmlStart(PLIST_T)
		ret += xmlEnd(PLIST_T)
	default:
		ret += e.CompileParameterList()
	}

	if err := e.eatSymbol(RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}
	e.Tknzr.Advance()

	ret += e.CompileSubroutineBody()

	return ret + xmlEnd(SUBROUTINEDEC_T)
}

func (e *Engine) CompileParameterList() string {
	ret := xmlStart(PLIST_T)

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

	return ret + xmlEnd(PLIST_T)
}

func (e *Engine) CompileSubroutineBody() string {
	ret := xmlStart(SUBROUTINEBODY_T)

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}
	e.Tknzr.Advance()

	for e.Tknzr.Keyword() == VAR {
		ret += e.CompileVarDec()
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		return ""
	}
	e.Tknzr.Advance()

	return ret + xmlEnd(SUBROUTINEBODY_T)
}

func (e *Engine) CompileVarDec() string {
	ret := xmlStart(VARDEC_T)

	if err := e.eatKeyword(VAR, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	e.Tknzr.Advance()
	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	e.Tknzr.Advance()
	for e.Tknzr.Symbol() == KOMMA {
		if err := e.eatSymbol(KOMMA, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
		}

		e.Tknzr.Advance()
		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
		}

		e.Tknzr.Advance()
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	e.Tknzr.Advance()
	return ret + xmlEnd(VARDEC_T)
}

func (e Engine) CompileStatements() string {
	ret := xmlStart(STATEMENTS_T)

	switch e.Tknzr.Keyword() {
	case LET:
		ret += e.CompileLetStatement()
	}

	ret += e.CompileReturn()
	return ret + xmlEnd(STATEMENTS_T)
}
func (e *Engine) CompileLetStatement() string {
	ret := xmlStart(LET_T)

	if err := e.eatKeyword(LET, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}
	e.Tknzr.Advance()

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}
	e.Tknzr.Advance()

	if e.Tknzr.Symbol() == LSQUARE {
		if err := e.eatSymbol(LSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}
		e.Tknzr.Advance()

		ret += e.CompileExpression()

		if err := e.eatSymbol(RSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}
		e.Tknzr.Advance()
	}

	if err := e.eatSymbol(EQUAL, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}
	e.Tknzr.Advance()

	if err := e.eatIntVal(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}
	e.Tknzr.Advance()

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}
	e.Tknzr.Advance()

	return ret + xmlEnd(LET_T)
}

func (e *Engine) CompileExpression() string {
	ret := xmlStart(EXPRESSION_T)

	ret += e.CompileTerm()

	for isOperator(e.Tknzr.Symbol()) {
		if err := e.eatSymbol(e.Tknzr.Symbol(), &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileExpression: %w", err))
		}
		e.Tknzr.Advance()

		ret += e.CompileTerm()
	}

	return ret + xmlEnd(EXPRESSION_T)
}

func (e Engine) isTerm() bool {
	switch e.Tknzr.TokenType() {
	case INT_CONST:
		return true
	case STRING_CONST:
		return true
	case KEYWORD:
		return true
	case IDENTIFIER:
		return true
	case SYMBOL:
		switch e.Tknzr.Symbol() {
		case LPAREN:
			return true
		case TILDE:
			return true
		case MINUS:
			return true
		}
	}
	return false
}

func (e *Engine) CompileTerm() string {
	ret := xmlStart(TERM_T)

	if !e.isTerm() {
		e.Errors = append(
			e.Errors,
			NewErrSyntaxUnexpectedToken(
				fmt.Sprintf(
					"TOKENTYPE '%s' / '%s' / '%s' / '%s' or SYMBOL '%s' / '%s' / '%s'",
					INT_CONST, STRING_CONST, KEYWORD, IDENTIFIER, LPAREN, TILDE, MINUS),
				e.Tknzr.curToken.Literal),
		)
	}

	if err := e.eatIntVal(&ret); err != nil {
		e.Errors = append(e.Errors,
			fmt.Errorf(
				"compileTerm: %w, %w",
				NewErrSyntaxUnexpectedToken("expression", e.Tknzr.curToken.Literal),
				err,
			),
		)
	}
	e.Tknzr.Advance()

	return ret + xmlEnd(TERM_T)
}

func (e *Engine) CompileReturn() string {
	ret := xmlStart(RETURN_T)

	if err := e.eatKeyword(RETURN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}
	e.Tknzr.Advance()

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}
	e.Tknzr.Advance()

	return ret + xmlEnd(RETURN_T)
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
			return NewErrSyntaxNotAType(e.Tknzr.Keyword())
			// return NewErrSyntaxUnexpectedToken(
			// 	fmt.Sprintf("KEYWORD %s / %s / %s or className", INT, CHAR, BOOLEAN),
			// 	e.Tknzr.Keyword())
		}
		*ret += xmlKeyword(e.Tknzr.Keyword())
	default:
		return NewErrSyntaxNotAType(e.Tknzr.curToken.Literal)
		// return NewErrSyntaxUnexpectedToken(
		// 	fmt.Sprintf("KEYWORD %s / %s / %s or className", INT, CHAR, BOOLEAN),
		// 	e.Tknzr.curToken.Literal)
	}
	return nil
}

func (e Engine) eatIntVal(ret *string) error {
	if e.Tknzr.TokenType() != INT_CONST {
		return NewErrSyntaxUnexpectedTokenType(INT_CONST, e.Tknzr.curToken.Literal)
	}
	*ret += xmlIntegerConst(e.Tknzr.IntVal())
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

func isOperator(sym string) bool {
	if _, ok := operators[sym]; ok {
		return true
	}
	return false
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

func xmlIntegerConst(val int) string {
	return fmt.Sprintf("<integerConstant> %d </integerConstant>\n", val)
}

func xmlIdentifier(identifier string) string {
	return fmt.Sprintf("<identifier> %s </identifier>\n", identifier)
}

func xmlStart(tag string) string {
	return fmt.Sprintf("<%s>\n", tag)
}
func xmlEnd(tag string) string {
	return fmt.Sprintf("</%s>\n", tag)
}
