package main

import (
	vmtrans "hack/vm_translator"
)

func main() {
	p, err := vmtrans.NewParser("./vm_translator/test_programs/StackArithmetic/SimpleAdd/SimpleAdd.vm")
	c := vmtrans.NewCodeWriter("out.asm")
	if err != nil {
		panic(err)
	}
	for p.HasMoreCommands() {
		p.Advance()
		cmdType, _ := p.CommandType()
		arg1, _ := p.Arg1()
		arg2, _ := p.Arg2()

		c.WriteArithmetic(cmdType, arg1, arg2)
	}
	c.CloseFile()
}
