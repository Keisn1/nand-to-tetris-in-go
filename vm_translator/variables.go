package vmtrans

import "errors"

var (
	ErrNotAdvanced = errors.New("you need to advance before reading")
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

var (
	templateFileNames = map[string]string{
		C_PUSH + " " + "local":     "pushGeneralX.asm",
		C_PUSH + " " + "argument":  "pushGeneralX.asm",
		C_PUSH + " " + "this":      "pushGeneralX.asm",
		C_PUSH + " " + "that":      "pushGeneralX.asm",
		C_PUSH + " " + "constant":  "pushConstantX.asm",
		C_ARITHMETIC + " " + "add": "add.asm",
	}
)

var (
	segmentRegisters = map[string]string{
		"local":    "R1",
		"argument": "R2",
		"this":     "R3",
		"that":     "R4",
	}
	segmentRegisterName = map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
	}
)
