package main

import (
	"flag"
	vmtrans "hack/vm_translator"
)

func main() {
	var fp string
	var out string

	flag.StringVar(&fp, "file", "SimpleAdd.vm", "Specify a file")
	flag.StringVar(&out, "out", "Alice", "Specify an output file")
	flag.Parse()

	p, err := vmtrans.NewParser(fp)
	c := vmtrans.NewCodeWriter(out)
	if err != nil {
		panic(err)
	}
	for p.HasMoreCommands() {
		p.Advance()
		cmdType, _ := p.CommandType()
		arg1, _ := p.Arg1()
		arg2, _ := p.Arg2()

		c.WriteArithmetic(cmdType, arg1, arg2)
		c.WriteNewline()
	}
	c.CloseFile()
}
