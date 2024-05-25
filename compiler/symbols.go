package compiler

const (
	LBRACE    = '{'
	RBRACE    = '}'
	LPAREN    = '('
	RPAREN    = ')'
	LSQUARE   = '['
	RSQUARE   = ']'
	POINT     = '.'
	KOMMA     = ','
	SEMICOLON = ';'
	PLUS      = '+'
	MINUS     = '-'
	STAR      = '*'
	SLASH     = '/'
	AND       = '&'
	PIPE      = '|'
	LT        = '<'
	GT        = '>'
	EQUAL     = '='
	TILDE     = '~'
)

var (
	allSymbols = map[byte]struct{}{
		LBRACE:    {},
		RBRACE:    {},
		LPAREN:    {},
		RPAREN:    {},
		LSQUARE:   {},
		RSQUARE:   {},
		POINT:     {},
		KOMMA:     {},
		SEMICOLON: {},
		PLUS:      {},
		MINUS:     {},
		STAR:      {},
		SLASH:     {},
		AND:       {},
		PIPE:      {},
		LT:        {},
		GT:        {},
		EQUAL:     {},
		TILDE:     {},
	}
)

var (
	xmlSymbols = map[rune]string{
		LBRACE:    "{",
		RBRACE:    "}",
		LPAREN:    "(",
		RPAREN:    ")",
		LSQUARE:   "[",
		RSQUARE:   "]",
		POINT:     ".",
		KOMMA:     ",",
		SEMICOLON: ";",
		PLUS:      "+",
		MINUS:     "-",
		STAR:      "*",
		SLASH:     "/",
		AND:       "&amp;",
		PIPE:      "|",
		LT:        "&lt;",
		GT:        "&gt;",
		EQUAL:     "=",
		TILDE:     "~",
	}
)
