package engine_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"hack/compiler/engine"
	"hack/compiler/token"
	"hack/compiler/vmWriter"

	"github.com/stretchr/testify/assert"
)

func Test_subroutineBody(t *testing.T) {
	t.Run("Testing happy subRoutineBodies", func(t *testing.T) {
		type testCase struct {
			name     string
			filename string
		}

		dir := "own/bodies/"
		testCases := []testCase{
			{name: "var int x;let x = 2+3;return;", filename: "letStatement2plus3"},
			{name: "var int x;let x = 3+4;return;", filename: "letStatement3plus4"},
			{name: "var int x;let x = -1;return;", filename: "unaryOperatorMinus"},
			{name: "var int x;let x = ~true;return;", filename: "unaryOperatorTilde"},
			{name: "var int x;let x = 3-4;return;", filename: "letStatement3minus4"},
			{name: "var int x, y;let x = x * y;return;", filename: "mathMulitply"},
			{name: "the body of the seven program", filename: "seven"},
			{name: "Screen.drawRectangle(x, y, x + width, y + height);", filename: "drawRectangle"},
			{name: "expression from the course", filename: "expressionFromCourse"},
			{name: "expression from the course with do", filename: "expressionFromCourseWithDo"},
			{name: "let statement function call with List", filename: "letStatementFunctionCall"},
			{name: "if without else", filename: "ifWithoutElse"},
			{name: "if with else", filename: "ifWithElse"},
			{name: "while statement", filename: "whileStatement"},
			{name: "multiple flow control statements", filename: "flowControl"},
			{name: "test that false constant is correct", filename: "testFalseToken"},
			{name: "construction and manipulation of objects", filename: "objectManipulation"},
			{name: "string handling", filename: "outputString"},
			{name: "array manipulation", filename: "arrayManipulation"},
			{name: "array manipulation 2", filename: "arrayManipulation2"},
			{name: "array manipulation 3", filename: "arrayManipulation3"},
			{name: "array manipulation 4", filename: "arrayManipulation4"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.filename+".jack")
				want := readFile(t, dir+tc.filename+".vm")

				tempDir := t.TempDir()
				out := filepath.Join(tempDir, tc.filename+".vm")
				vw := vmWriter.NewVmWriter(out)
				tknzr := token.NewTokenizer(input)
				e := engine.NewEngine(&tknzr, &vw)

				e.Tknzr.Advance()
				e.CompileSubroutineDec()
				vw.Close()

				assert.Empty(t, e.Errors)
				assert.FileExists(t, out)

				got, err := os.ReadFile(out)
				fmt.Println(string(got))
				assert.NoError(t, err)
				assert.Equal(t, want, string(got))
				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
			})
		}
	})
}

func Test_compileClass(t *testing.T) {
	t.Run("Testing happy classes", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
			dir  string
		}

		dir := "own/classes/"
		testCases := []testCase{
			{name: "2d point from course", fp: "Point2d"},
			{name: "3d point from course", fp: "Point3d"},
			{name: "method call in Point class", fp: "PointMethod"},
			{name: "seven program", fp: "Seven"},
			{name: "Main of Square game", fp: "SquareGameMain"},
			{name: "field var is object, calling method on that object", fp: "SquareGameFieldVarWithDo"},
			{name: "field var is object, calling method on that object with do", fp: "SquareGameFieldVar"},
			{name: "print the score left bottom", dir: "PrintScore/", fp: "Main"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.dir != "" {
					dir = dir + tc.dir
				}
				input := readFile(t, dir+tc.fp+".jack")
				want := readFile(t, dir+tc.fp+".vm")

				tempDir := t.TempDir()
				out := filepath.Join(tempDir, tc.fp+".vm")
				vw := vmWriter.NewVmWriter(out)
				tknzr := token.NewTokenizer(input)
				e := engine.NewEngine(&tknzr, &vw)

				e.Tknzr.Advance()
				e.CompileClass()
				vw.Close()

				assert.Empty(t, e.Errors)
				assert.FileExists(t, out)

				got, err := os.ReadFile(out)
				fmt.Println(string(got))
				assert.NoError(t, err)
				assert.Equal(t, want, string(got))
				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
			})
		}
	})
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

func assertErrorFound(t *testing.T, gotErrs []error, wantErr error) {
	t.Helper()
	foundErr := false
	for _, gotErr := range gotErrs {
		if errors.Is(gotErr, wantErr) {
			foundErr = true
		}
	}

	if !foundErr {
		t.Log("Could not find error:", wantErr, "\nPresent errors: ", gotErrs)
		t.Fail()
	}
}
