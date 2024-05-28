package compiler

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

var unaryOperators = map[string]struct{}{
	TILDE: {},
	MINUS: {},
}
