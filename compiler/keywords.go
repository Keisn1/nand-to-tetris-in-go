package compiler

const (
	CLASS       = "class"
	METHOD      = "method"
	FUNCTION    = "function"
	CONSTRUCTOR = "constructor"
	INT         = "int"
	BOOLEAN     = "boolean"
	CHAR        = "char"
	VOID        = "void"
	VAR         = "var"
	STATIC      = "static"
	FIELD       = "field"
	LET         = "let"
	DO          = "do"
	IF          = "if"
	ELSE        = "else"
	WHILE       = "while"
	RETURN      = "return"
	TRUE        = "true"
	FALSE       = "false"
	NULL        = "null"
	THIS        = "this"
	EOF         = "eof"
)

var (
	keywords = map[string]string{
		"class":       CLASS,
		"method":      METHOD,
		"function":    FUNCTION,
		"constructor": CONSTRUCTOR,
		"int":         INT,
		"boolean":     BOOLEAN,
		"char":        CHAR,
		"void":        VOID,
		"var":         VAR,
		"static":      STATIC,
		"field":       FIELD,
		"let":         LET,
		"do":          DO,
		"if":          IF,
		"else":        ELSE,
		"while":       WHILE,
		"return":      RETURN,
		"true":        TRUE,
		"false":       FALSE,
		"null":        NULL,
		"this":        THIS,
		"eof":         EOF,
	}
)
