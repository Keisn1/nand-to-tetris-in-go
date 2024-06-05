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

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	for isStaticOrField(e.Tknzr.Keyword()) {
		ret += e.CompileClassVarDec()
	}

	for isSubRoutineDec(e.Tknzr.Keyword()) {
		ret += e.CompileSubroutineDec()
	}

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	return ret + xmlEnd(CLASS_T)
}

func (e *Engine) CompileClassVarDec() string {
	ret := xmlStart(CLASSVARDEC_T)

	switch e.Tknzr.Keyword() {
	case STATIC:
		e.eatKeyword(STATIC, &ret)
	case FIELD:
		e.eatKeyword(FIELD, &ret)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotAClassVarDec(e.Tknzr.curToken.Literal))
	}

	if err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	for e.Tknzr.Symbol() == KOMMA {
		e.eatSymbol(KOMMA, &ret)

		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
		}
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	return ret + xmlEnd(CLASSVARDEC_T)
}

func (e *Engine) CompileSubroutineDec() string {
	ret := xmlStart(SUBROUTINEDEC_T)

	switch e.Tknzr.Keyword() {
	case CONSTRUCTOR:
		e.eatKeyword(CONSTRUCTOR, &ret)
	case FUNCTION:
		e.eatKeyword(FUNCTION, &ret)
	case METHOD:
		e.eatKeyword(METHOD, &ret)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotASubroutineDec(e.Tknzr.Keyword()))
	}

	switch e.Tknzr.Keyword() {
	case VOID:
		e.eatKeyword(VOID, &ret)
	default:
		if err := e.eatType(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
		}
	}

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if err := e.eatSymbol(LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

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

	ret += e.CompileSubroutineBody()

	return ret + xmlEnd(SUBROUTINEDEC_T)
}

func (e *Engine) CompileParameterList() string {
	ret := xmlStart(PLIST_T)

	if err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	for e.Tknzr.Symbol() == KOMMA {
		e.eatSymbol(KOMMA, &ret)

		if err := e.eatType(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}

		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}
	}

	return ret + xmlEnd(PLIST_T)
}

func (e *Engine) CompileSubroutineBody() string {
	ret := xmlStart(SUBROUTINEBODY_T)

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	for e.Tknzr.Keyword() == VAR {
		ret += e.CompileVarDec()
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	return ret + xmlEnd(SUBROUTINEBODY_T)
}

func (e *Engine) CompileVarDec() string {
	ret := xmlStart(VARDEC_T)

	if err := e.eatKeyword(VAR, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	if err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	for e.Tknzr.Symbol() == KOMMA {
		e.eatSymbol(KOMMA, &ret)

		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
		}

	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	return ret + xmlEnd(VARDEC_T)
}
func (e *Engine) CompileStatements() string {
	ret := xmlStart(STATEMENTS_T)

	for isStatement(e.Tknzr.Keyword()) {
		switch e.Tknzr.Keyword() {
		case LET:
			ret += e.CompileLetStatement()
		case DO:
			ret += e.CompileDoStatement()
		case IF:
			ret += e.CompileIfStatement()
		case WHILE:
			ret += e.CompileWhileStatement()
		case RETURN:
			ret += e.CompileReturn()
		}
	}

	return ret + xmlEnd(STATEMENTS_T)
}

func (e *Engine) CompileIfStatement() string {
	ret := xmlStart(IF_T)

	if err := e.eatKeyword(IF, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileExpression()

	if err := e.eatSymbol(RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if e.Tknzr.Keyword() == ELSE {
		if err := e.eatKeyword(ELSE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}

		if err := e.eatSymbol(LBRACE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}

		ret += e.CompileStatements()

		if err := e.eatSymbol(RBRACE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}
	}
	return ret + xmlEnd(IF_T)
}

func (e *Engine) CompileWhileStatement() string {
	ret := xmlStart(WHILE_T)

	if err := e.eatKeyword(WHILE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileExpression()

	if err := e.eatSymbol(RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	return ret + xmlEnd(WHILE_T)
}

func (e *Engine) CompileDoStatement() string {
	ret := xmlStart(DO_T)

	if err := e.eatKeyword(DO, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	switch e.Tknzr.Symbol() {
	case LSQUARE:
		if err := e.eatSymbol(LSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		ret += e.CompileExpression()

		if err := e.eatSymbol(RSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

	case LPAREN:
		if err := e.eatSymbol(LPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		ret += e.CompileExpressionList()

		if err := e.eatSymbol(RPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}
	case DOT:
		if err := e.eatSymbol(DOT, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if err := e.eatSymbol(LPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		ret += e.CompileExpressionList()

		if err := e.eatSymbol(RPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}
	return ret + xmlEnd(DO_T)
}

func (e *Engine) CompileLetStatement() string {
	ret := xmlStart(LET_T)

	if err := e.eatKeyword(LET, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	if err := e.eatIdentifier(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	if e.Tknzr.Symbol() == LSQUARE {
		if err := e.eatSymbol(LSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}

		ret += e.CompileExpression()

		if err := e.eatSymbol(RSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}
	}

	if err := e.eatSymbol(EQUAL, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	ret += e.CompileExpression()

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	return ret + xmlEnd(LET_T)
}

func (e *Engine) CompileExpression() string {
	ret := xmlStart(EXPRESSION_T)

	ret += e.CompileTerm()

	for isOperator(e.Tknzr.Symbol()) {
		e.eatSymbol(e.Tknzr.Symbol(), &ret)

		ret += e.CompileTerm()
	}

	return ret + xmlEnd(EXPRESSION_T)
}

func (e *Engine) CompileTerm() string {
	ret := xmlStart(TERM_T)

	switch e.Tknzr.TokenType() {
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotATerm(e.Tknzr.curToken.Literal))

	case INT_CONST:
		if err := e.eatIntVal(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case STRING_CONST:
		if err := e.eatStringVal(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case KEYWORD:
		switch e.Tknzr.Keyword() {
		case TRUE:
			e.eatKeyword(TRUE, &ret)
		case FALSE:
			e.eatKeyword(FALSE, &ret)
		case NULL:
			e.eatKeyword(NULL, &ret)
		case THIS:
			e.eatKeyword(THIS, &ret)
		default:
			e.Errors = append(e.Errors, NewErrSyntaxNotAKeywordConst(e.Tknzr.Keyword()))
		}

	case SYMBOL:
		switch e.Tknzr.Symbol() {

		default:
			e.Errors = append(e.Errors, NewErrSyntaxNotATerm(e.Tknzr.curToken.Literal))
		case LPAREN:
			if err := e.eatSymbol(LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf(" compileTerm: %w", err))
			}

			ret += e.CompileExpression()

			if err := e.eatSymbol(RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

		case MINUS:
			if err := e.eatSymbol(MINUS, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileTerm()

		case TILDE:
			if err := e.eatSymbol(TILDE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileTerm()
		}

	case IDENTIFIER:
		if err := e.eatIdentifier(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

		switch e.Tknzr.Symbol() {
		case LSQUARE:
			if err := e.eatSymbol(LSQUARE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpression()

			if err := e.eatSymbol(RSQUARE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

		case LPAREN:
			if err := e.eatSymbol(LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpressionList()

			if err := e.eatSymbol(RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}
		case DOT:
			if err := e.eatSymbol(DOT, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if err := e.eatIdentifier(&ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if err := e.eatSymbol(LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpressionList()

			if err := e.eatSymbol(RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}
		}
	}

	return ret + xmlEnd(TERM_T)
}

func (e *Engine) CompileExpressionList() string {
	ret := xmlStart("expressionList")

	if e.isTerm() {
		ret += e.CompileExpression()

		for e.Tknzr.Symbol() == KOMMA {
			e.eatSymbol(KOMMA, &ret)
			ret += e.CompileExpression()
		}
	}
	return ret + xmlEnd("expressionList")
}

func (e *Engine) CompileReturn() string {
	ret := xmlStart(RETURN_T)

	if err := e.eatKeyword(RETURN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	if e.isTerm() {
		ret += e.CompileExpression()
	}

	if err := e.eatSymbol(SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	return ret + xmlEnd(RETURN_T)
}

func (e Engine) eatType(ret *string) error {
	switch e.Tknzr.TokenType() {
	case IDENTIFIER:
		e.eatIdentifier(ret)

	case KEYWORD:
		switch e.Tknzr.Keyword() {
		case INT:
			e.eatKeyword(INT, ret)
		case CHAR:
			e.eatKeyword(CHAR, ret)
		case BOOLEAN:
			e.eatKeyword(BOOLEAN, ret)
		default:
			return NewErrSyntaxNotAType(e.Tknzr.Keyword())
		}
	default:
		return NewErrSyntaxNotAType(e.Tknzr.curToken.Literal)
	}
	return nil
}

func (e Engine) eatIntVal(ret *string) error {
	if e.Tknzr.TokenType() != INT_CONST {
		return NewErrSyntaxUnexpectedTokenType(INT_CONST, e.Tknzr.curToken.Literal)
	}

	*ret += xmlIntegerConst(e.Tknzr.IntVal())
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatStringVal(ret *string) error {
	if e.Tknzr.TokenType() != STRING_CONST {
		return NewErrSyntaxUnexpectedTokenType(STRING_CONST, e.Tknzr.curToken.Literal)
	}

	*ret += xmlStringConst(e.Tknzr.StringVal())
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatIdentifier(ret *string) error {
	if e.Tknzr.TokenType() != IDENTIFIER {
		return NewErrSyntaxUnexpectedTokenType(IDENTIFIER, e.Tknzr.curToken.Literal)
	}

	*ret += xmlIdentifier(e.Tknzr.Identifier())
	e.Tknzr.Advance()
	return nil
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

func (e Engine) eatKeyword(expectedKeyword string, ret *string) error {
	if (e.Tknzr.TokenType() != KEYWORD) || (e.Tknzr.Keyword() != expectedKeyword) {
		return NewErrSyntaxUnexpectedToken(expectedKeyword, e.Tknzr.curToken.Literal)
	}

	*ret += xmlKeyword(expectedKeyword)
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatSymbol(expectedSymbol string, ret *string) error {
	if (e.Tknzr.TokenType() != SYMBOL) || (e.Tknzr.Symbol() != expectedSymbol) {
		return NewErrSyntaxUnexpectedToken(expectedSymbol, e.Tknzr.curToken.Literal)
	}

	*ret += xmlSymbol(expectedSymbol)
	e.Tknzr.Advance()
	return nil
}

func isStatement(kw string) bool {
	statements := map[string]struct{}{
		LET:    {},
		IF:     {},
		WHILE:  {},
		DO:     {},
		RETURN: {},
	}
	if _, ok := statements[kw]; ok {
		return true
	}
	return false
}

func isSubRoutineDec(kw string) bool {
	return kw == CONSTRUCTOR || kw == FUNCTION || kw == METHOD
}

func isStaticOrField(kw string) bool {
	return kw == STATIC || kw == FIELD
}

func isOperator(sym string) bool {
	if _, ok := operators[sym]; ok {
		return true
	}
	return false
}

func xmlSymbol(symbol string) string {
	return fmt.Sprintf("<symbol> %s </symbol>\n", xmlSymbols[string(symbol)])
}

func xmlKeyword(kw string) string {
	return fmt.Sprintf("<keyword> %s </keyword>\n", kw)
}

func xmlIntegerConst(val int) string {
	return fmt.Sprintf("<integerConstant> %d </integerConstant>\n", val)
}

func xmlStringConst(val string) string {
	return fmt.Sprintf("<stringConstant> %s </stringConstant>\n", val)
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
