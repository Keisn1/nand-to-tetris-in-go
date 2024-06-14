package engine

import (
	"hack/compiler/token"
	"hack/compiler/vmWriter"
)

var tokenToVmWriter = map[string]string{
	token.PLUS:  vmWriter.ADD,
	token.MINUS: vmWriter.SUB,
}
