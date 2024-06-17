package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullPrograms(t *testing.T) {
	type testCase struct {
		path         string
		wantOutFiles []string
	}

	dir := "test_programs/project11/"
	testCases := []testCase{
		{

			path:         dir + "Seven",
			wantOutFiles: []string{"Main.vm"},
		},
		{

			path:         dir + "ConvertToBin",
			wantOutFiles: []string{"Main.vm"},
		},
		{

			path:         dir + "Square",
			wantOutFiles: []string{"Main.vm", "Square.vm", "SquareGame.vm"},
		},
		{

			path:         dir + "Pong",
			wantOutFiles: []string{"Main.vm", "Bat.vm", "Ball.vm", "PongGame.vm"},
		},
		{
			path:         dir + "Average",
			wantOutFiles: []string{"Main.vm"},
		},
		{
			path:         dir + "ComplexArrays",
			wantOutFiles: []string{"Main.vm"},
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		log.SetOutput(&buf)

		flag.CommandLine = flag.NewFlagSet("cmdLineArgs", flag.ExitOnError)
		os.Args = []string{"flag ", "-path", tc.path}

		main()

		info, err := os.Stat(tc.path)
		if err != nil {
			log.Fatal(err)
		}

		dirName := tc.path
		if !info.IsDir() {
			dirName, _ = filepath.Split(tc.path)
		}

		for _, out := range tc.wantOutFiles {
			assert.FileExists(t, filepath.Join(dirName, out))

			got := removeWhiteSpaces(readFile(t, filepath.Join(dirName, out)))
			want := removeWhiteSpaces(readFile(t, filepath.Join(dirName, "Compare"+out)))
			assert.Equal(t, want, got)
		}
	}
}

func readFile(t *testing.T, fp string) string {
	t.Helper()
	input, err := os.ReadFile(fp)
	assert.NoError(t, err, "error reading file")
	return string(input)
}

func removeWhiteSpaces(input string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(input, "")
}
