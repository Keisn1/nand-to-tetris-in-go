package main

import (
	"flag"
	"fmt"
	assembler "hack/project_06_Assembler"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var fp string
	var out string

	// flag.StringVar(&fp, "file", "", "Specify a file")
	flag.StringVar(&out, "out", "", "Specify an output file")
	flag.Parse()

	nonFlagArgs := flag.Args()
	if len(nonFlagArgs) <= 0 {
		fmt.Println("Need to specify an input-file.")
		os.Exit(1)
	}

	if len(nonFlagArgs) > 1 {
		fmt.Println("Can only assemble one file at a time.")
		os.Exit(1)
	}

	fp = nonFlagArgs[0]

	if out == "" {
		out = strings.TrimSuffix(fp, filepath.Ext(fp)) + ".hack"
	}

	fmt.Println("Compiling", fp)
	fmt.Println("outfile:", out)

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
