package vmtrans

import "errors"

var (
	ErrNotAdvanced = errors.New("you need to advance before reading")
)

const (
	pushLocalX = `//push local {{ .x }}
//*R13=LCL + {{ .x }}
@{{ .x }}
D=A
@R1
D=M+D
@R13
M=D

// SP--
@SP
M=M-1

// *R13=*SP
@R13`
)

const (
	C_ARITHMETIC = "C_ARITHMETIC"
	C_PUSH       = "C_PUSH"
	C_POP        = "C_POP"
	C_LABEL      = "C_LABEL"
	C_GOTO       = "C_GOTO"
	C_IF         = "C_IF"
	C_FUNCTION   = "C_FUNCTION"
	C_RETURN     = "C_RETURN"
	C_CALL       = "C_CALL"
)

var (
	cmdTable = map[string]string{
		"add":      C_ARITHMETIC,
		"sub":      C_ARITHMETIC,
		"neg":      C_ARITHMETIC,
		"eq":       C_ARITHMETIC,
		"get":      C_ARITHMETIC,
		"lt":       C_ARITHMETIC,
		"and":      C_ARITHMETIC,
		"or":       C_ARITHMETIC,
		"not":      C_ARITHMETIC,
		"push":     C_PUSH,
		"pop":      C_POP,
		"label":    C_LABEL,
		"goto":     C_GOTO,
		"if":       C_IF,
		"function": C_FUNCTION,
		"return":   C_RETURN,
		"call":     C_CALL,
	}
)
