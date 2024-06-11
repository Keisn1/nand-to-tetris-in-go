package engine

import (
	"fmt"
	"hack/compiler/symbolTable"
	"hack/compiler/token"
)

type Engine struct {
	Tknzr     *token.Tokenizer
	symTab    symbolTable.SymbolTable
	clAndSubR map[string]string
	Errors    []error
}

func NewEngine(tknzr *token.Tokenizer) Engine {
	return Engine{
		Tknzr:     tknzr,
		symTab:    symbolTable.NewSymbolTable(),
		clAndSubR: make(map[string]string),
	}
}

func (e *Engine) CompileClass() string {
	ret := xmlStart(CLASS_T)

	if err := e.eatKeyword(token.CLASS, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	e.clAndSubR[e.Tknzr.Identifier()] = "class"
	if err := e.eatIdentifier(&ret, true); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	for isStaticOrField(e.Tknzr.Keyword()) {
		ret += e.CompileClassVarDec()
	}

	for isSubRoutineDec(e.Tknzr.Keyword()) {
		e.symTab.StartSubroutine()
		ret += e.CompileSubroutineDec()
	}

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	return ret + xmlEnd(CLASS_T)
}

func (e *Engine) CompileClassVarDec() string {
	ret := xmlStart(CLASSVARDEC_T)

	var kind string
	switch e.Tknzr.Keyword() {
	case token.STATIC:
		kind = symbolTable.STATIC
		e.eatKeyword(token.STATIC, &ret)
	case token.FIELD:
		kind = symbolTable.FIELD
		e.eatKeyword(token.FIELD, &ret)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotAClassVarDec(e.Tknzr.GetTokenLiteral()))
	}

	varType := e.Tknzr.Keyword()
	if _, err := e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
	if err := e.eatIdentifier(&ret, true); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	for e.Tknzr.Symbol() == token.KOMMA {
		e.eatSymbol(token.KOMMA, &ret)

		e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
		if err := e.eatIdentifier(&ret, true); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
		}
	}

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	return ret + xmlEnd(CLASSVARDEC_T)
}

func (e *Engine) CompileSubroutineDec() string {
	ret := xmlStart(SUBROUTINEDEC_T)

	switch e.Tknzr.Keyword() {
	case token.CONSTRUCTOR:
		e.eatKeyword(token.CONSTRUCTOR, &ret)
	case token.FUNCTION:
		e.eatKeyword(token.FUNCTION, &ret)
	case token.METHOD:
		e.eatKeyword(token.METHOD, &ret)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotASubroutineDec(e.Tknzr.Keyword()))
	}

	switch e.Tknzr.Keyword() {
	case token.VOID:
		e.eatKeyword(token.VOID, &ret)
	default:
		if _, err := e.eatType(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
		}
	}

	e.clAndSubR[e.Tknzr.Identifier()] = "subroutine"
	if err := e.eatIdentifier(&ret, true); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	switch e.Tknzr.Symbol() {
	case token.RPAREN:
		ret += xmlStart(PLIST_T)
		ret += xmlEnd(PLIST_T)
	default:
		ret += e.CompileParameterList()
	}

	if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	ret += e.CompileSubroutineBody()

	return ret + xmlEnd(SUBROUTINEDEC_T)
}

func (e *Engine) CompileParameterList() string {
	ret := xmlStart(PLIST_T)

	kind := symbolTable.ARG

	var varType string
	var err error
	if varType, err = e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
	if err := e.eatIdentifier(&ret, true); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	for e.Tknzr.Symbol() == token.KOMMA {
		e.eatSymbol(token.KOMMA, &ret)

		if varType, err = e.eatType(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}

		e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
		if err := e.eatIdentifier(&ret, true); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}
	}

	return ret + xmlEnd(PLIST_T)
}

func (e *Engine) CompileSubroutineBody() string {
	ret := xmlStart(SUBROUTINEBODY_T)

	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	for e.Tknzr.Keyword() == token.VAR {
		ret += e.CompileVarDec()
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	return ret + xmlEnd(SUBROUTINEBODY_T)
}

func (e *Engine) CompileVarDec() string {
	ret := xmlStart(VARDEC_T)

	if err := e.eatKeyword(token.VAR, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	var varType string
	var err error
	if varType, err = e.eatType(&ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	e.symTab.Define(e.Tknzr.Identifier(), varType, symbolTable.VAR)
	if err := e.eatIdentifier(&ret, true); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	for e.Tknzr.Symbol() == token.KOMMA {
		e.eatSymbol(token.KOMMA, &ret)

		e.symTab.Define(e.Tknzr.Identifier(), varType, symbolTable.VAR)
		if err := e.eatIdentifier(&ret, true); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
		}

	}

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	return ret + xmlEnd(VARDEC_T)
}
func (e *Engine) CompileStatements() string {
	ret := xmlStart(STATEMENTS_T)

	for isStatement(e.Tknzr.Keyword()) {
		switch e.Tknzr.Keyword() {
		case token.LET:
			ret += e.CompileLetStatement()
		case token.DO:
			ret += e.CompileDoStatement()
		case token.IF:
			ret += e.CompileIfStatement()
		case token.WHILE:
			ret += e.CompileWhileStatement()
		case token.RETURN:
			ret += e.CompileReturn()
		}
	}

	return ret + xmlEnd(STATEMENTS_T)
}

func (e *Engine) CompileIfStatement() string {
	ret := xmlStart(IF_T)

	if err := e.eatKeyword(token.IF, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileExpression()

	if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if e.Tknzr.Keyword() == token.ELSE {
		if err := e.eatKeyword(token.ELSE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}

		if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}

		ret += e.CompileStatements()

		if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}
	}
	return ret + xmlEnd(IF_T)
}

func (e *Engine) CompileWhileStatement() string {
	ret := xmlStart(WHILE_T)

	if err := e.eatKeyword(token.WHILE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileExpression()

	if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	return ret + xmlEnd(WHILE_T)
}

func (e *Engine) CompileDoStatement() string {
	ret := xmlStart(DO_T)

	if err := e.eatKeyword(token.DO, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatIdentifier(&ret, false); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	switch e.Tknzr.Symbol() {
	case token.LSQUARE:
		if err := e.eatSymbol(token.LSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		ret += e.CompileExpression()

		if err := e.eatSymbol(token.RSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

	case token.LPAREN:
		if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		ret += e.CompileExpressionList()

		if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}
	case token.DOT:
		if err := e.eatSymbol(token.DOT, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if err := e.eatIdentifier(&ret, false); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		ret += e.CompileExpressionList()

		if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}
	}

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}
	return ret + xmlEnd(DO_T)
}

func (e *Engine) CompileLetStatement() string {
	ret := xmlStart(LET_T)

	if err := e.eatKeyword(token.LET, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	if err := e.eatIdentifier(&ret, false); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	if e.Tknzr.Symbol() == token.LSQUARE {
		if err := e.eatSymbol(token.LSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}

		ret += e.CompileExpression()

		if err := e.eatSymbol(token.RSQUARE, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}
	}

	if err := e.eatSymbol(token.EQUAL, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	ret += e.CompileExpression()

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
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
		e.Errors = append(e.Errors, NewErrSyntaxNotATerm(e.Tknzr.GetTokenLiteral()))

	case token.INT_CONST:
		if err := e.eatIntVal(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case token.STRING_CONST:
		if err := e.eatStringVal(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case token.KEYWORD:
		switch e.Tknzr.Keyword() {
		case token.TRUE:
			e.eatKeyword(token.TRUE, &ret)
		case token.FALSE:
			e.eatKeyword(token.FALSE, &ret)
		case token.NULL:
			e.eatKeyword(token.NULL, &ret)
		case token.THIS:
			e.eatKeyword(token.THIS, &ret)
		default:
			e.Errors = append(e.Errors, NewErrSyntaxNotAKeywordConst(e.Tknzr.Keyword()))
		}

	case token.SYMBOL:
		switch e.Tknzr.Symbol() {

		default:
			e.Errors = append(e.Errors, NewErrSyntaxNotATerm(e.Tknzr.GetTokenLiteral()))
		case token.LPAREN:
			if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf(" compileTerm: %w", err))
			}

			ret += e.CompileExpression()

			if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

		case token.MINUS:
			if err := e.eatSymbol(token.MINUS, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileTerm()

		case token.TILDE:
			if err := e.eatSymbol(token.TILDE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileTerm()
		}

	case token.IDENTIFIER:
		if err := e.eatIdentifier(&ret, false); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

		switch e.Tknzr.Symbol() {
		case token.LSQUARE:
			if err := e.eatSymbol(token.LSQUARE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpression()

			if err := e.eatSymbol(token.RSQUARE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

		case token.LPAREN:
			if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpressionList()

			if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}
		case token.DOT:
			if err := e.eatSymbol(token.DOT, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if err := e.eatIdentifier(&ret, false); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpressionList()

			if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
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

		for e.Tknzr.Symbol() == token.KOMMA {
			e.eatSymbol(token.KOMMA, &ret)
			ret += e.CompileExpression()
		}
	}
	return ret + xmlEnd("expressionList")
}

func (e *Engine) CompileReturn() string {
	ret := xmlStart(RETURN_T)

	if err := e.eatKeyword(token.RETURN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	if e.isTerm() {
		ret += e.CompileExpression()
	}

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	return ret + xmlEnd(RETURN_T)
}

func (e Engine) eatType(ret *string) (string, error) {
	var varType string
	switch e.Tknzr.TokenType() {
	case token.IDENTIFIER:
		varType = e.Tknzr.Identifier()
		e.eatIdentifier(ret, false)

	case token.KEYWORD:
		varType = e.Tknzr.Keyword()
		switch varType {
		case token.INT:
			e.eatKeyword(token.INT, ret)
		case token.CHAR:
			e.eatKeyword(token.CHAR, ret)
		case token.BOOLEAN:
			e.eatKeyword(token.BOOLEAN, ret)
		default:
			return "", NewErrSyntaxNotAType(e.Tknzr.Keyword())
		}
	default:
		return "", NewErrSyntaxNotAType(e.Tknzr.GetTokenLiteral())
	}
	return varType, nil
}

func (e Engine) eatIntVal(ret *string) error {
	if e.Tknzr.TokenType() != token.INT_CONST {
		return NewErrSyntaxUnexpectedTokenType(token.INT_CONST, e.Tknzr.GetTokenLiteral())
	}

	*ret += xmlIntegerConst(e.Tknzr.IntVal())
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatStringVal(ret *string) error {
	if e.Tknzr.TokenType() != token.STRING_CONST {
		return NewErrSyntaxUnexpectedTokenType(token.STRING_CONST, e.Tknzr.GetTokenLiteral())
	}

	*ret += xmlStringConst(e.Tknzr.StringVal())
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatIdentifier(ret *string, isDefinition bool) error {
	if e.Tknzr.TokenType() != token.IDENTIFIER {
		return NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, e.Tknzr.GetTokenLiteral())
	}

	name := e.Tknzr.Identifier()
	e.Tknzr.Advance()

	var info string
	if e.symTab.KindOf(name) != "" {
		info = fmt.Sprintf("_%s_%d", e.symTab.KindOf(name), e.symTab.IndexOf(name))
	} else {
		switch e.Tknzr.Symbol() {
		case token.LPAREN:
			info += "_subroutine"
		default:
			info += "_class"
		}
	}

	if isDefinition {
		info += "_def"
	} else {
		info += "_used"
	}

	*ret += xmlIdentifier(name, info)
	return nil
}

func (e Engine) isTerm() bool {
	switch e.Tknzr.TokenType() {
	case token.INT_CONST:
		return true
	case token.STRING_CONST:
		return true
	case token.KEYWORD:
		return true
	case token.IDENTIFIER:
		return true
	case token.SYMBOL:
		switch e.Tknzr.Symbol() {
		case token.LPAREN:
			return true
		case token.TILDE:
			return true
		case token.MINUS:
			return true
		}
	}
	return false
}

func (e Engine) eatKeyword(expectedKeyword string, ret *string) error {
	if (e.Tknzr.TokenType() != token.KEYWORD) || (e.Tknzr.Keyword() != expectedKeyword) {
		return NewErrSyntaxUnexpectedToken(expectedKeyword, e.Tknzr.GetTokenLiteral())
	}

	*ret += xmlKeyword(expectedKeyword)
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatSymbol(expectedSymbol string, ret *string) error {
	if (e.Tknzr.TokenType() != token.SYMBOL) || (e.Tknzr.Symbol() != expectedSymbol) {
		return NewErrSyntaxUnexpectedToken(expectedSymbol, e.Tknzr.GetTokenLiteral())
	}

	*ret += xmlSymbol(expectedSymbol)
	e.Tknzr.Advance()
	return nil
}

func isStatement(kw string) bool {
	statements := map[string]struct{}{
		token.LET:    {},
		token.IF:     {},
		token.WHILE:  {},
		token.DO:     {},
		token.RETURN: {},
	}
	if _, ok := statements[kw]; ok {
		return true
	}
	return false
}

func isSubRoutineDec(kw string) bool {
	return kw == token.CONSTRUCTOR || kw == token.FUNCTION || kw == token.METHOD
}

func isStaticOrField(kw string) bool {
	return kw == token.STATIC || kw == token.FIELD
}

func isOperator(sym string) bool {
	if _, ok := token.Operators[sym]; ok {
		return true
	}
	return false
}

func xmlSymbol(symbol string) string {
	return fmt.Sprintf("<symbol> %s </symbol>\n", token.XmlSymbols[string(symbol)])
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

func xmlIdentifier(name, info string) string {
	return fmt.Sprintf("<identifier%s> %s </identifier%s>\n", info, name, info)
}

func xmlStart(tag string) string {
	return fmt.Sprintf("<%s>\n", tag)
}
func xmlEnd(tag string) string {
	return fmt.Sprintf("</%s>\n", tag)
}
