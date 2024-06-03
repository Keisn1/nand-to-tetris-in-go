package compiler

const (
	LBRACE    = "{"
	RBRACE    = "}"
	LPAREN    = "("
	RPAREN    = ")"
	LSQUARE   = "["
	RSQUARE   = "]"
	DOT       = "."
	KOMMA     = ","
	SEMICOLON = ";"
	PLUS      = "+"
	MINUS     = "-"
	STAR      = "*"
	SLASH     = "/"
	AND       = "&"
	PIPE      = "|"
	LT        = "<"
	GT        = ">"
	EQUAL     = "="
	TILDE     = "~"
)

var (
	allSymbols = map[string]struct{}{
		LBRACE:    {},
		RBRACE:    {},
		LPAREN:    {},
		RPAREN:    {},
		LSQUARE:   {},
		RSQUARE:   {},
		DOT:       {},
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

var operators = map[string]struct{}{
	PLUS:  {},
	MINUS: {},
	STAR:  {},
	SLASH: {},
	AND:   {},
	PIPE:  {},
	LT:    {},
	GT:    {},
	EQUAL: {},
}
