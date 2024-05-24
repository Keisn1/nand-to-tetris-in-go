package compiler

type Keyword string

const (
	CLASS       Keyword = "class"
	METHOD      Keyword = "method"
	FUNCTION    Keyword = "function"
	CONSTRUCTOR Keyword = "constructor"
	INT         Keyword = "int"
	BOOLEAN     Keyword = "boolean"
	CHAR        Keyword = "char"
	VOID        Keyword = "void"
	VAR         Keyword = "var"
	STATIC      Keyword = "static"
	FIELD       Keyword = "field"
	LET         Keyword = "let"
	DO          Keyword = "do"
	IF          Keyword = "if"
	ELSE        Keyword = "else"
	WHILE       Keyword = "while"
	RETURN      Keyword = "return"
	TRUE        Keyword = "true"
	FALSE       Keyword = "false"
	NULL        Keyword = "null"
	THIS        Keyword = "this"
)

var (
	keywords = map[string]Keyword{
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
	}
)
