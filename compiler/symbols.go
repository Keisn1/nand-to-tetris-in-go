package compiler

const (
	LBRACE    = '{'
	RBRACE    = '}'
	LPAREN    = '('
	RPAREN    = ')'
	LBRAC     = '['
	RBRAC     = ']'
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
		LBRAC:     {},
		RBRAC:     {},
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
