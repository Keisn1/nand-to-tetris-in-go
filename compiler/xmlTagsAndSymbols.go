package compiler

const (
	CLASS_T          = "class"
	CLASSVARDEC_T    = "classVarDec"
	SUBROUTINEDEC_T  = "subroutineDec"
	PLIST_T          = "parameterList"
	SUBROUTINEBODY_T = "subroutineBody"
	VARDEC_T         = "varDec"
	STATEMENTS_T     = "statements"
	LET_T            = "letStatement"
	IF_T             = "ifStatement"
	DO_T             = "doStatement"
	WHILE_T          = "whileStatement"
	RETURN_T         = "returnStatement"
	EXPRESSION_T     = "expression"
	TERM_T           = "term"
)

var (
	xmlSymbols = map[string]string{
		LBRACE:    LBRACE,
		RBRACE:    RBRACE,
		LPAREN:    LPAREN,
		RPAREN:    RPAREN,
		LSQUARE:   LSQUARE,
		RSQUARE:   RSQUARE,
		DOT:       DOT,
		KOMMA:     KOMMA,
		SEMICOLON: SEMICOLON,
		PLUS:      PLUS,
		MINUS:     MINUS,
		STAR:      STAR,
		SLASH:     SLASH,
		AND:       "&amp;",
		PIPE:      PIPE,
		LT:        "&lt;",
		GT:        "&gt;",
		EQUAL:     EQUAL,
		TILDE:     TILDE,
	}
)
