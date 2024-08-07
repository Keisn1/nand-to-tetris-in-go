package engine

import (
	"hack/compiler/token"
	"hack/compiler/vmWriter"
)

var opSymbolToVmWriter = map[string]string{
	token.PLUS:  vmWriter.ADD,
	token.MINUS: vmWriter.SUB,
	token.EQUAL: vmWriter.EQUAL,
	token.GT:    vmWriter.GT,
	token.LT:    vmWriter.LT,
	token.AND:   vmWriter.AND,
	token.PIPE:  vmWriter.OR,
}

var unaryOpToVmWriter = map[string]string{
	token.MINUS: vmWriter.NEG,
	token.TILDE: vmWriter.NOT,
}
