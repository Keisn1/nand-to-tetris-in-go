package main

import (
	"flag"
	"hack/compiler"
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
				out := filepath.Join(path, fileNameWithoutExt+".xml")

				f, err := os.Create(out)
				if err != nil {
					log.Fatal("Couldn't create file", err)
				}

				input, err := os.ReadFile(filepath.Join(path, dirEntry.Name()))
				if err != nil {
					log.Fatal("Couldn't read input", err)
				}

				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				s := engine.CompileClass()
				f.Write([]byte(s))

				f.Close()
			}
		}
	} else {

		dir, fileName := filepath.Split(path)
		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		out := filepath.Join(dir, fileNameWithoutExt+".xml")

		f, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		input, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		tknzr := compiler.NewTokenizer(string(input))
		engine := compiler.NewEngine(&tknzr)

		engine.Tknzr.Advance()
		s := engine.CompileClass()
		f.Write([]byte(s))
	}

}
