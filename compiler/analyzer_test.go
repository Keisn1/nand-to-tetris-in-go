package compiler_test

import (
	"hack/compiler"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_analyzer(t *testing.T) {
	filesToRemove := []string{}
	t.Cleanup(func() {
		for _, fp := range filesToRemove {
			os.RemoveAll(fp)
		}
	})
	t.Run("Produce xml for empty main without expression", func(t *testing.T) {
		fp := "test_files/Test1/Main.jack"
		analyzer := compiler.NewAnalyzer(fp)
		analyzer.Analyze()

		fpOut := "test_files/Test1/MainT.xml"
		filesToRemove = append(filesToRemove, fpOut)
		assert.FileExists(t, fpOut)

		regex := regexp.MustCompile(`\s+`)
		gotByte, err := os.ReadFile(fpOut)
		got := regex.ReplaceAllString(string(gotByte), "")
		assert.NoError(t, err)

		fpTest := "test_files/Test1/test_MainT.xml"
		wantByte, err := os.ReadFile(fpTest)
		want := regex.ReplaceAllString(string(wantByte), "")
		assert.NoError(t, err)

		assert.Equal(t, want, got)
	})
}
