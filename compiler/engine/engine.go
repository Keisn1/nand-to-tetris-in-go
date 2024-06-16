package engine

import (
	"fmt"
	"hack/compiler/symbolTable"
	"hack/compiler/token"
	"hack/compiler/vmWriter"
	"strconv"
)

type Engine struct {
	Tknzr      *token.Tokenizer
	vmWriter   *vmWriter.VmWriter
	symTab     symbolTable.SymbolTable
	Errors     []error
	labelCount int
}

func NewEngine(tknzr *token.Tokenizer, vw *vmWriter.VmWriter) Engine {
	return Engine{
		Tknzr:    tknzr,
		vmWriter: vw,
		symTab:   symbolTable.NewSymbolTable(),
	}
}

func (e *Engine) CompileClass() string {
	ret := xmlStart(CLASS_T)

	if err := e.eatKeyword(token.CLASS, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

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

	var varType string
	var err error
	if varType, err = e.eatType(&ret); err != nil {
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

	subroutineType := e.Tknzr.Keyword()
	if err := e.eatKeyword(e.Tknzr.Keyword(), &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	switch subroutineType {
	case token.CONSTRUCTOR:
	case token.FUNCTION:
	case token.METHOD:
		e.symTab.Define(vmWriter.THIS, "Point", symbolTable.ARG)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotASubroutineDec(subroutineType))
	}

	switch e.Tknzr.Keyword() {
	case token.VOID:
		e.eatKeyword(token.VOID, &ret)
	default:
		if _, err := e.eatType(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
		}
	}

	name := e.Tknzr.Identifier()
	if err := e.eatIdentifier(&ret, true); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if e.Tknzr.Symbol() != token.RPAREN {
		ret += e.CompileParameterList()
	}

	if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	var count int
	for e.Tknzr.Keyword() == token.VAR {
		count += e.CompileVarDec()
	}

	e.vmWriter.WriteFunction(e.vmWriter.GetFilename()+"."+name, count)
	switch subroutineType {
	case token.CONSTRUCTOR:
		e.vmWriter.WritePush(vmWriter.CONST, e.symTab.VarCount(symbolTable.FIELD))
		e.vmWriter.WriteCall("Memory.alloc", 1)
		e.vmWriter.WritePop(vmWriter.POINTER, 0)
	case token.FUNCTION:
	case token.METHOD:
		e.symTab.Define(vmWriter.THIS, "Point", symbolTable.ARG)
		e.vmWriter.WritePush(vmWriter.ARG, 0)
		e.vmWriter.WritePop(vmWriter.POINTER, 0)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotASubroutineDec(subroutineType))
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	return ret + xmlEnd(SUBROUTINEBODY_T)
}

func (e *Engine) CompileVarDec() int {
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

	count := 1
	for e.Tknzr.Symbol() == token.KOMMA {
		count++
		e.eatSymbol(token.KOMMA, &ret)

		e.symTab.Define(e.Tknzr.Identifier(), varType, symbolTable.VAR)
		if err := e.eatIdentifier(&ret, true); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
		}
	}

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	return count
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
	labelCountStr := strconv.Itoa(e.labelCount)
	e.labelCount++

	if err := e.eatKeyword(token.IF, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileExpression()

	e.vmWriter.WriteArithmetic(vmWriter.NOT)

	if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	// statement 1
	e.vmWriter.WriteIf(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")
	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	ret += e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteGoto(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")
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

	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
	return ret + xmlEnd(IF_T)
}

func (e *Engine) CompileWhileStatement() string {
	ret := xmlStart(WHILE_T)
	labelCountStr := strconv.Itoa(e.labelCount)
	e.labelCount++

	if err := e.eatKeyword(token.WHILE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")
	ret += e.CompileExpression()

	if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteArithmetic(vmWriter.NOT)
	e.vmWriter.WriteIf(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
	ret += e.CompileStatements()
	e.vmWriter.WriteGoto(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")

	if err := e.eatSymbol(token.RBRACE, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
	return ret + xmlEnd(WHILE_T)
}

func (e *Engine) CompileDoStatement() string {
	ret := xmlStart(DO_T)

	if err := e.eatKeyword(token.DO, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	identifier := e.Tknzr.Identifier()
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
		e.vmWriter.WritePush(vmWriter.POINTER, 0)
		if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		count := e.CompileExpressionList()

		if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		e.vmWriter.WriteCall(e.vmWriter.GetFilename()+token.DOT+identifier, count+1)
		e.vmWriter.WritePop(vmWriter.TEMP, 0)

	case token.DOT:
		if err := e.eatSymbol(token.DOT, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		subroutineName := e.Tknzr.Identifier()
		if err := e.eatIdentifier(&ret, false); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if e.symTab.TypeOf(identifier) != "" {
			if e.symTab.KindOf(identifier) == token.FIELD {
				e.vmWriter.WritePush(vmWriter.THIS, e.symTab.IndexOf(identifier))
			} else {
				e.vmWriter.WritePush(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))
			}
			count := e.CompileExpressionList()
			e.vmWriter.WriteCall(e.symTab.TypeOf(identifier)+token.DOT+subroutineName, count+1)
			e.vmWriter.WritePop(vmWriter.TEMP, 0)
		} else {
			count := e.CompileExpressionList()
			e.vmWriter.WriteCall(identifier+token.DOT+subroutineName, count)
			e.vmWriter.WritePop(vmWriter.TEMP, 0)
		}

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

	identifier := e.Tknzr.Identifier()
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

	switch e.symTab.KindOf(identifier) {
	case token.FIELD:
		e.vmWriter.WritePop(vmWriter.THIS, e.symTab.IndexOf(identifier))
	default:
		e.vmWriter.WritePop(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))

	}

	return ret + xmlEnd(LET_T)
}

func (e *Engine) CompileExpression() string {
	ret := xmlStart(EXPRESSION_T)

	ret += e.CompileTerm()

	for isOperator(e.Tknzr.Symbol()) {
		op := e.Tknzr.Symbol()
		e.eatSymbol(e.Tknzr.Symbol(), &ret)

		ret += e.CompileTerm()

		// we are compiling term (op term)*
		// defer is LastIn FirstOut, so it is working like the stack machine
		switch op {
		default:
			defer e.vmWriter.WriteArithmetic(opSymbolToVmWriter[op])
		case token.STAR:
			defer e.vmWriter.WriteCall("Math.multiply", 2)
		case token.SLASH:
			defer e.vmWriter.WriteCall("Math.divide", 2)
		}

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
		e.vmWriter.WritePush(vmWriter.CONST, len(e.Tknzr.StringVal()))
		e.vmWriter.WriteCall("String.new", 1)

		for _, c := range e.Tknzr.StringVal() {
			e.vmWriter.WritePush(vmWriter.CONST, int(c))
			e.vmWriter.WriteCall("String.appendChar", 2)
		}

		if err := e.eatStringVal(&ret); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case token.KEYWORD:
		switch e.Tknzr.Keyword() {
		case token.TRUE:
			e.eatKeyword(token.TRUE, &ret)
			e.vmWriter.WritePush(vmWriter.CONST, vmWriter.TRUE)
			e.vmWriter.WriteArithmetic(vmWriter.NEG)
		case token.FALSE:
			e.eatKeyword(token.FALSE, &ret)
			e.vmWriter.WritePush(vmWriter.CONST, vmWriter.FALSE)
		case token.NULL:
			e.eatKeyword(token.NULL, &ret)
			e.vmWriter.WritePush(vmWriter.CONST, vmWriter.NULL)
		case token.THIS:
			e.eatKeyword(token.THIS, &ret)
			e.vmWriter.WritePush(vmWriter.POINTER, 0)
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
			e.vmWriter.WriteArithmetic(unaryOpToVmWriter[token.MINUS])

		case token.TILDE:
			if err := e.eatSymbol(token.TILDE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileTerm()
			e.vmWriter.WriteArithmetic(unaryOpToVmWriter[token.TILDE])
		}

	case token.IDENTIFIER:
		identifier := e.Tknzr.Identifier()
		if err := e.eatIdentifier(&ret, false); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

		switch e.Tknzr.Symbol() {
		default:
			switch e.symTab.KindOf(identifier) {
			case token.FIELD:
				e.vmWriter.WritePush(vmWriter.THIS, e.symTab.IndexOf(identifier))
			default:
				e.vmWriter.WritePush(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))
			}

		case token.LSQUARE:
			if err := e.eatSymbol(token.LSQUARE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			ret += e.CompileExpression()

			if err := e.eatSymbol(token.RSQUARE, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

		case token.LPAREN:
			e.vmWriter.WritePush(vmWriter.POINTER, 0)

			if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			count := e.CompileExpressionList()

			if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			e.vmWriter.WriteCall(e.vmWriter.GetFilename()+token.DOT+identifier, count+1)

		case token.DOT:
			if err := e.eatSymbol(token.DOT, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			subroutineName := e.Tknzr.Identifier()
			if err := e.eatIdentifier(&ret, false); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if err := e.eatSymbol(token.LPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if e.symTab.TypeOf(identifier) != "" {
				if e.symTab.KindOf(identifier) == token.FIELD {
					e.vmWriter.WritePush(vmWriter.THIS, e.symTab.IndexOf(identifier))
				} else {
					e.vmWriter.WritePush(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))
				}
				count := e.CompileExpressionList()
				e.vmWriter.WriteCall(e.symTab.TypeOf(identifier)+token.DOT+subroutineName, count+1)
			} else {
				count := e.CompileExpressionList()
				e.vmWriter.WriteCall(identifier+token.DOT+subroutineName, count)
			}

			if err := e.eatSymbol(token.RPAREN, &ret); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}
		}
	}

	return ret + xmlEnd(TERM_T)
}

func (e *Engine) CompileExpressionList() int {
	ret := xmlStart("expressionList")

	count := 0
	if e.isTerm() {
		count++
		ret += e.CompileExpression()

		for e.Tknzr.Symbol() == token.KOMMA {
			count++
			e.eatSymbol(token.KOMMA, &ret)
			ret += e.CompileExpression()
		}
	}
	return count
}

func (e *Engine) CompileReturn() string {
	ret := xmlStart(RETURN_T)

	if err := e.eatKeyword(token.RETURN, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	if e.isTerm() {
		ret += e.CompileExpression()
	} else {
		e.vmWriter.WritePush(vmWriter.CONST, 0)
	}

	if err := e.eatSymbol(token.SEMICOLON, &ret); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	e.vmWriter.WriteReturn()

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

	e.vmWriter.WritePush(vmWriter.CONST, e.Tknzr.IntVal())
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
