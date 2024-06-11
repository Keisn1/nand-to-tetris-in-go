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
		path     string
		outFiles []string
	}
	testCases := []testCase{
		{

			path:     "test_programs/project10/ArrayTest/Main.jack",
			outFiles: []string{"Main.xml"},
		},
		{
			path:     "test_programs/project10/ExpressionLessSquare",
			outFiles: []string{"Main.xml", "Square.xml", "SquareGame.xml"},
		},
		{
			path:     "test_programs/project10/Square",
			outFiles: []string{"Main.xml", "Square.xml", "SquareGame.xml"},
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		log.SetOutput(&buf)

		flag.CommandLine = flag.NewFlagSet("cmdLineArgs", flag.ExitOnError)
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

		for _, out := range tc.outFiles {
			assert.FileExists(t, filepath.Join(dirName, out))

			got := removeWhiteSpaces(readFile(t, filepath.Join(dirName, out)))
			want := removeWhiteSpaces(readFile(t, filepath.Join(dirName, "compare"+out)))
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
