package compiler

type Symbol string

const (
	LBRACE Symbol = "{"
	RBRACE Symbol = "}"
)

var (
	symbols = map[Symbol]struct{}{
		LBRACE: {},
		RBRACE: {},
	}
)
