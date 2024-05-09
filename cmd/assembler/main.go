package main

import (
	"flag"
	"hack/assembler"
	"os"
)

func main() {
	var fp string
	var out string

	flag.StringVar(&fp, "file", "SimpleAdd.vm", "Specify a file")
	flag.StringVar(&out, "out", "Alice", "Specify an output file")
	flag.Parse()

	asm, err := assembler.NewAssembler(fp)
	if err != nil {
		panic(err)
	}

	got := asm.Assemble()

	os.Remove(out)
	file, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write([]byte(got))
}
