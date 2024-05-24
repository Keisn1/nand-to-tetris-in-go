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
