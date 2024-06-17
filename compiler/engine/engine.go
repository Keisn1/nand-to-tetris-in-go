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

func (e *Engine) CompileClass() {
	if err := e.eatKeyword(token.CLASS); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}

	for isStaticOrField(e.Tknzr.Keyword()) {
		e.CompileClassVarDec()
	}

	for isSubRoutineDec(e.Tknzr.Keyword()) {
		e.CompileSubroutine()
	}

	if err := e.eatSymbol(token.RBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClass: %w", err))
	}
}

func (e *Engine) CompileClassVarDec() {
	kind := e.Tknzr.Keyword()
	switch kind {
	case token.STATIC, token.FIELD:
		e.Tknzr.Advance()
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotAClassVarDec(e.Tknzr.GetTokenLiteral()))
	}

	varType, err := e.eatType()
	if err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}

	for e.Tknzr.Symbol() == token.KOMMA {
		e.eatSymbol(token.KOMMA)

		e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
		if err := e.eatIdentifier(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
		}
	}

	if err := e.eatSymbol(token.SEMICOLON); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileClassVarDec: %w", err))
	}
}

func (e *Engine) CompileSubroutine() {
	e.symTab.StartSubroutine()

	// ------------------SubroutineDeclaration------------------------------------------
	subroutineType := e.Tknzr.Keyword()
	if err := e.eatKeyword(e.Tknzr.Keyword()); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	switch subroutineType {
	case token.CONSTRUCTOR, token.FUNCTION:
	case token.METHOD:
		e.symTab.Define(vmWriter.THIS, "Point", symbolTable.ARG)
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotASubroutineDec(subroutineType))
	}

	switch e.Tknzr.Keyword() {
	case token.VOID:
		e.eatKeyword(token.VOID)
	default:
		if _, err := e.eatType(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
		}
	}

	subroutineName := e.Tknzr.Identifier()
	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	if e.Tknzr.Symbol() != token.RPAREN {
		e.CompileParameterList()
	}

	if err := e.eatSymbol(token.RPAREN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineDec: %w", err))
	}

	// ------------------SubroutineBody------------------------------------------
	if err := e.eatSymbol(token.LBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}

	var count int
	for e.Tknzr.Keyword() == token.VAR {
		count += e.CompileVarDec()
	}
	e.vmWriter.WriteFunction(e.vmWriter.GetFilename()+"."+subroutineName, count)

	switch subroutineType {
	case token.CONSTRUCTOR:
		e.vmWriter.WritePush(vmWriter.CONST, e.symTab.VarCount(symbolTable.FIELD))
		e.vmWriter.WriteCall("Memory.alloc", 1)
		e.vmWriter.WritePop(vmWriter.POINTER, 0)
	case token.METHOD:
		e.vmWriter.WritePush(vmWriter.ARG, 0)
		e.vmWriter.WritePop(vmWriter.POINTER, 0)
	case token.FUNCTION:
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotASubroutineDec(subroutineType))
	}

	e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileSubroutineBody: %w", err))
	}
}

func (e *Engine) CompileVarDec() int {

	if err := e.eatKeyword(token.VAR); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	varType, err := e.eatType()
	if err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	e.symTab.Define(e.Tknzr.Identifier(), varType, symbolTable.VAR)
	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	count := 1
	for e.Tknzr.Symbol() == token.KOMMA {
		count++
		e.eatSymbol(token.KOMMA)

		e.symTab.Define(e.Tknzr.Identifier(), varType, symbolTable.VAR)
		if err := e.eatIdentifier(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
		}
	}

	if err := e.eatSymbol(token.SEMICOLON); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileVarDec: %w", err))
	}

	return count
}

func (e *Engine) CompileParameterList() {
	kind := symbolTable.ARG

	varType, err := e.eatType()
	if err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
	}

	for e.Tknzr.Symbol() == token.KOMMA {
		e.eatSymbol(token.KOMMA)

		if varType, err = e.eatType(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}

		e.symTab.Define(e.Tknzr.Identifier(), varType, kind)
		if err := e.eatIdentifier(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileParameterList: %w", err))
		}
	}

}

func (e *Engine) CompileStatements() {
	for isStatement(e.Tknzr.Keyword()) {
		switch e.Tknzr.Keyword() {
		case token.LET:
			e.CompileLetStatement()
		case token.DO:
			e.CompileDoStatement()
		case token.IF:
			e.CompileIfStatement()
		case token.WHILE:
			e.CompileWhileStatement()
		case token.RETURN:
			e.CompileReturn()
		}
	}
}

func (e *Engine) CompileIfStatement() {
	e.labelCount++
	labelCountStr := strconv.Itoa(e.labelCount)

	if err := e.eatKeyword(token.IF); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.CompileExpression()

	e.vmWriter.WriteArithmetic(vmWriter.NOT)

	if err := e.eatSymbol(token.RPAREN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	// statement 1
	e.vmWriter.WriteIf(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")
	if err := e.eatSymbol(token.LBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.CompileStatements()

	if err := e.eatSymbol(token.RBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteGoto(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")
	if e.Tknzr.Keyword() == token.ELSE {
		if err := e.eatKeyword(token.ELSE); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}

		if err := e.eatSymbol(token.LBRACE); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}

		e.CompileStatements()

		if err := e.eatSymbol(token.RBRACE); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
		}
	}

	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
}

func (e *Engine) CompileWhileStatement() {
	labelCountStr := strconv.Itoa(e.labelCount)
	e.labelCount++

	if err := e.eatKeyword(token.WHILE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LPAREN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")
	e.CompileExpression()

	if err := e.eatSymbol(token.RPAREN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	if err := e.eatSymbol(token.LBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteArithmetic(vmWriter.NOT)
	e.vmWriter.WriteIf(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
	e.CompileStatements()
	e.vmWriter.WriteGoto(e.vmWriter.GetFilename() + "." + labelCountStr + ".L1")

	if err := e.eatSymbol(token.RBRACE); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	e.vmWriter.WriteLabel(e.vmWriter.GetFilename() + "." + labelCountStr + ".L2")
}

func (e *Engine) CompileDoStatement() {
	if err := e.eatKeyword(token.DO); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	identifier := e.Tknzr.Identifier()
	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}

	switch e.Tknzr.Symbol() {
	case token.LPAREN:
		e.vmWriter.WritePush(vmWriter.POINTER, 0)
		if err := e.eatSymbol(token.LPAREN); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		count := e.CompileExpressionList()

		if err := e.eatSymbol(token.RPAREN); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		e.vmWriter.WriteCall(e.vmWriter.GetFilename()+token.DOT+identifier, count+1)
		e.vmWriter.WritePop(vmWriter.TEMP, 0)

	case token.DOT:
		if err := e.eatSymbol(token.DOT); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		subroutineName := e.Tknzr.Identifier()
		if err := e.eatIdentifier(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}

		if err := e.eatSymbol(token.LPAREN); err != nil {
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

		if err := e.eatSymbol(token.RPAREN); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
		}
	}

	if err := e.eatSymbol(token.SEMICOLON); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileDoStatement: %w", err))
	}
}

func (e *Engine) CompileLetStatement() {

	if err := e.eatKeyword(token.LET); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	identifier := e.Tknzr.Identifier()
	if err := e.eatIdentifier(); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	if e.Tknzr.Symbol() == token.LSQUARE {
		if err := e.eatSymbol(token.LSQUARE); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}

		e.vmWriter.WritePush(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))
		e.CompileExpression()

		e.vmWriter.WriteArithmetic(vmWriter.ADD)

		if err := e.eatSymbol(token.RSQUARE); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}

		if err := e.eatSymbol(token.EQUAL); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}

		e.CompileExpression()

		if err := e.eatSymbol(token.SEMICOLON); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
		}

		e.vmWriter.WritePop(vmWriter.TEMP, 0)
		e.vmWriter.WritePop(vmWriter.POINTER, 1)
		e.vmWriter.WritePush(vmWriter.TEMP, 0)
		e.vmWriter.WritePop(vmWriter.THAT, 0)

		return
	}

	if err := e.eatSymbol(token.EQUAL); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	e.CompileExpression()

	if err := e.eatSymbol(token.SEMICOLON); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileLetStatement: %w", err))
	}

	switch e.symTab.KindOf(identifier) {
	case token.FIELD:
		e.vmWriter.WritePop(vmWriter.THIS, e.symTab.IndexOf(identifier))
	default:
		e.vmWriter.WritePop(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))
	}

}

func (e *Engine) CompileExpression() {
	e.CompileTerm()

	for isOperator(e.Tknzr.Symbol()) {
		op := e.Tknzr.Symbol()
		e.eatSymbol(e.Tknzr.Symbol())

		e.CompileTerm()

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
}

func (e *Engine) CompileTerm() {
	switch e.Tknzr.TokenType() {
	default:
		e.Errors = append(e.Errors, NewErrSyntaxNotATerm(e.Tknzr.GetTokenLiteral()))

	case token.INT_CONST:
		if err := e.eatIntVal(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case token.STRING_CONST:
		e.vmWriter.WritePush(vmWriter.CONST, len(e.Tknzr.StringVal()))
		e.vmWriter.WriteCall("String.new", 1)

		for _, c := range e.Tknzr.StringVal() {
			e.vmWriter.WritePush(vmWriter.CONST, int(c))
			e.vmWriter.WriteCall("String.appendChar", 2)
		}

		if err := e.eatStringVal(); err != nil {
			e.Errors = append(e.Errors, fmt.Errorf("compileTerm: , %w", err))
		}

	case token.KEYWORD:
		switch e.Tknzr.Keyword() {
		case token.TRUE:
			e.eatKeyword(token.TRUE)
			e.vmWriter.WritePush(vmWriter.CONST, vmWriter.TRUE)
			e.vmWriter.WriteArithmetic(vmWriter.NEG)
		case token.FALSE:
			e.eatKeyword(token.FALSE)
			e.vmWriter.WritePush(vmWriter.CONST, vmWriter.FALSE)
		case token.NULL:
			e.eatKeyword(token.NULL)
			e.vmWriter.WritePush(vmWriter.CONST, vmWriter.NULL)
		case token.THIS:
			e.eatKeyword(token.THIS)
			e.vmWriter.WritePush(vmWriter.POINTER, 0)
		default:
			e.Errors = append(e.Errors, NewErrSyntaxNotAKeywordConst(e.Tknzr.Keyword()))
		}

	case token.SYMBOL:
		switch e.Tknzr.Symbol() {
		default:
			e.Errors = append(e.Errors, NewErrSyntaxNotATerm(e.Tknzr.GetTokenLiteral()))
		case token.LPAREN:
			if err := e.eatSymbol(token.LPAREN); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf(" compileTerm: %w", err))
			}

			e.CompileExpression()

			if err := e.eatSymbol(token.RPAREN); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

		case token.MINUS:
			if err := e.eatSymbol(token.MINUS); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			e.CompileTerm()
			e.vmWriter.WriteArithmetic(unaryOpToVmWriter[token.MINUS])

		case token.TILDE:
			if err := e.eatSymbol(token.TILDE); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			e.CompileTerm()
			e.vmWriter.WriteArithmetic(unaryOpToVmWriter[token.TILDE])
		}

	case token.IDENTIFIER:
		identifier := e.Tknzr.Identifier()
		if err := e.eatIdentifier(); err != nil {
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
			if err := e.eatSymbol(token.LSQUARE); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			e.vmWriter.WritePush(e.symTab.KindOf(identifier), e.symTab.IndexOf(identifier))
			e.CompileExpression()
			e.vmWriter.WriteArithmetic(vmWriter.ADD)

			if err := e.eatSymbol(token.RSQUARE); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			e.vmWriter.WritePop(vmWriter.POINTER, 1)
			e.vmWriter.WritePush(vmWriter.THAT, 0)

		case token.LPAREN:
			e.vmWriter.WritePush(vmWriter.POINTER, 0)

			if err := e.eatSymbol(token.LPAREN); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			count := e.CompileExpressionList()

			if err := e.eatSymbol(token.RPAREN); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			e.vmWriter.WriteCall(e.vmWriter.GetFilename()+token.DOT+identifier, count+1)

		case token.DOT:
			if err := e.eatSymbol(token.DOT); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			subroutineName := e.Tknzr.Identifier()
			if err := e.eatIdentifier(); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}

			if err := e.eatSymbol(token.LPAREN); err != nil {
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

			if err := e.eatSymbol(token.RPAREN); err != nil {
				e.Errors = append(e.Errors, fmt.Errorf("compileTerm: %w", err))
			}
		}
	}
}

func (e *Engine) CompileExpressionList() int {
	count := 0
	if e.isTerm() {
		count++
		e.CompileExpression()

		for e.Tknzr.Symbol() == token.KOMMA {
			count++
			e.eatSymbol(token.KOMMA)
			e.CompileExpression()
		}
	}
	return count
}

func (e *Engine) CompileReturn() {

	if err := e.eatKeyword(token.RETURN); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	if e.isTerm() {
		e.CompileExpression()
	} else {
		e.vmWriter.WritePush(vmWriter.CONST, 0)
	}

	if err := e.eatSymbol(token.SEMICOLON); err != nil {
		e.Errors = append(e.Errors, fmt.Errorf("compileReturn: %w", err))
	}

	e.vmWriter.WriteReturn()
}

func (e Engine) eatType() (string, error) {
	var varType string
	switch e.Tknzr.TokenType() {
	case token.IDENTIFIER:
		varType = e.Tknzr.Identifier()
		e.eatIdentifier()

	case token.KEYWORD:
		varType = e.Tknzr.Keyword()
		switch varType {
		case token.INT, token.CHAR, token.BOOLEAN:
			e.eatKeyword(varType)
		default:
			return "", NewErrSyntaxNotAType(e.Tknzr.Keyword())
		}
	default:
		return "", NewErrSyntaxNotAType(e.Tknzr.GetTokenLiteral())
	}
	return varType, nil
}

func (e Engine) eatIntVal() error {
	if e.Tknzr.TokenType() != token.INT_CONST {
		return NewErrSyntaxUnexpectedTokenType(token.INT_CONST, e.Tknzr.GetTokenLiteral())
	}

	e.vmWriter.WritePush(vmWriter.CONST, e.Tknzr.IntVal())
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatStringVal() error {
	if e.Tknzr.TokenType() != token.STRING_CONST {
		return NewErrSyntaxUnexpectedTokenType(token.STRING_CONST, e.Tknzr.GetTokenLiteral())
	}

	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatIdentifier() error {
	if e.Tknzr.TokenType() != token.IDENTIFIER {
		return NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, e.Tknzr.GetTokenLiteral())
	}

	e.Tknzr.Advance()
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

func (e Engine) eatKeyword(expectedKeyword string) error {
	if (e.Tknzr.TokenType() != token.KEYWORD) || (e.Tknzr.Keyword() != expectedKeyword) {
		return NewErrSyntaxUnexpectedToken(expectedKeyword, e.Tknzr.GetTokenLiteral())
	}
	e.Tknzr.Advance()
	return nil
}

func (e Engine) eatSymbol(expectedSymbol string) error {
	if (e.Tknzr.TokenType() != token.SYMBOL) || (e.Tknzr.Symbol() != expectedSymbol) {
		return NewErrSyntaxUnexpectedToken(expectedSymbol, e.Tknzr.GetTokenLiteral())
	}

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
