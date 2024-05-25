package compiler_test

import (
	"hack/compiler"
	"os"
	"testing"

	"regexp"

	"github.com/stretchr/testify/assert"
)

func Test_compileClass(t *testing.T) {
	type testCase struct {
		name        string
		inputFile   string
		compareFile string
	}

	testCases := []testCase{
		{name: "empty main", inputFile: "test_programs/own/emptyMain.jack", compareFile: "test_programs/own/emptyMain.xml"},
	}

	for _, tc := range testCases {
		input, err := os.ReadFile(tc.inputFile)
		assert.NoError(t, err, "error reading file")

		tknzr := compiler.NewTokenizer(string(input))
		engine := compiler.NewEngine(tknzr)

		tknzr.Advance()
		assert.Equal(t, compiler.KEYWORD, tknzr.TokenType())
		assert.Equal(t, compiler.CLASS, tknzr.Keyword())

		got := engine.CompileClass()
		wantByte, err := os.ReadFile(tc.compareFile)
		assert.NoError(t, err, "error reading file")

		regex := regexp.MustCompile(`\s+`)
		want := regex.ReplaceAllString(string(wantByte), "")
		got = regex.ReplaceAllString(got, "")
		assert.Equal(t, want, got)
	}
}

func removeWhiteSpaces(input string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(input, "")
}
