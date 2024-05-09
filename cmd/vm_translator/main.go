package main

import (
	"flag"
	"fmt"
	vmtrans "hack/vm_translator"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var fp string
	var out string

	flag.StringVar(&fp, "file", "SimpleAdd.vm", "Specify a file")
	flag.StringVar(&out, "out", "Alice", "Specify an output file")
	flag.Parse()

	os.Remove(out)

	fileInfo, err := os.Stat(fp)
	if err != nil {
		log.Fatal("Error getting file info:", err)
	}

	if fileInfo.IsDir() {
		translateDir(fp, out)
	} else {
		translateFile(fp, out)
	}
}

func translateFile(fp, out string) {
	p, err := vmtrans.NewParser(fp)
	c := vmtrans.NewCodeWriter(out, fp)
	if err != nil {
		panic(err)
	}

	for p.HasMoreCommands() {
		p.Advance()
		cmdType, err := p.CommandType()
		if err != nil {
			panic(err)
		}
		arg1, err := p.Arg1()
		if err != nil {
			panic(err)
		}
		arg2, err := p.Arg2()
		if err != nil {
			panic(err)
		}

		c.Write(cmdType, arg1, arg2)
		c.WriteNewline()
	}
	c.CloseFile()

}

func translateDir(fp, out string) {
	files, err := filepath.Glob(fp + "/*.vm")
	if err != nil {
		fmt.Println("Error matching files:", err)
		return
	}

	c := vmtrans.NewCodeWriter(out, fp)
	c.WriteBootStrap()
	c.CloseFile()

	for _, file := range files {
		translateFile(file, out)
	}
}
