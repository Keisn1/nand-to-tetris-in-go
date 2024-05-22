package compiler

type Keyword string

const (
	CLASS    Keyword = "CLASS"
	METHOD   Keyword = "METHOD"
	FUNCTION Keyword = "FUNCTION"
	STATIC   Keyword = "STATIC"
)

var (
	keywords = map[string]Keyword{
		"class":    CLASS,
		"method":   METHOD,
		"function": FUNCTION,
		"static":   STATIC,
	}
)
