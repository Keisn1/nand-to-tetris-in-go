package main

import (
	"flag"
	"hack/compiler/engine"
	"hack/compiler/token"
	"hack/compiler/vmWriter"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var path string

	flag.StringVar(&path, "path", "Main.jack", "Specify a file")
	flag.Parse()

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, dirEntry := range dirEntries {
			if filepath.Ext(dirEntry.Name()) == ".jack" {
				fileNameWithoutExt := strings.TrimSuffix(dirEntry.Name(), ".jack")
				out := filepath.Join(path, fileNameWithoutExt+".vm")

				input, err := os.ReadFile(filepath.Join(path, dirEntry.Name()))
				if err != nil {
					log.Fatal("Couldn't read input", err)
				}

				tknzr := token.NewTokenizer(string(input))
				vw := vmWriter.NewVmWriter(out)
				e := engine.NewEngine(&tknzr, &vw)

				e.Tknzr.Advance()
				e.CompileClass()
				vw.Close()
			}
		}
	} else {

		dir, fileName := filepath.Split(path)
		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		out := filepath.Join(dir, fileNameWithoutExt+".vm")

		f, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		input, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		tknzr := token.NewTokenizer(string(input))
		vw := vmWriter.NewVmWriter(out)
		e := engine.NewEngine(&tknzr, &vw)

		e.Tknzr.Advance()
		e.CompileClass()
		vw.Close()
	}
}
