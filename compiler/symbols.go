package compiler

const (
	LBRACE    = "{"
	RBRACE    = "}"
	LPAREN    = "("
	RPAREN    = ")"
	LSQUARE   = "["
	RSQUARE   = "]"
	POINT     = "."
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
	xmlSymbols = map[string]string{
		LBRACE:    LBRACE,
		RBRACE:    RBRACE,
		LPAREN:    LPAREN,
		RPAREN:    RPAREN,
		LSQUARE:   LSQUARE,
		RSQUARE:   RSQUARE,
		POINT:     POINT,
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
