package vmtrans_test

import (
	vmtrans "hack/vm_translator"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser(t *testing.T) {
	t.Run("Happy line reading", func(t *testing.T) {
		type testCase struct {
			name    string
			content string
			cmds    []string
		}

		testCases := []testCase{
			{
				name:    "one line",
				content: "first line",
				cmds:    []string{"first line"},
			},
			{
				name: "two lines",
				content: `first line
second line`,
				cmds: []string{"first line", "second line"},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				file := createTestFile(t, tc.content)
				defer os.Remove(file.Name())

				p, err := vmtrans.NewParser(file.Name())
				assert.NoError(t, err)
				var got []string
				for p.HasMoreCommands() {
					p.Advance()
					curCmd, _ := p.CurrentCmd()
					got = append(got, curCmd)
				}

				assert.Equal(t, tc.cmds, got)
			})
		}

	})
	t.Run("Empty file", func(t *testing.T) {
		file := createTestFile(t, "")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		assert.NoError(t, err)
		assert.False(t, p.HasMoreCommands())
	})

	t.Run("Reading before advancing", func(t *testing.T) {
		file := createTestFile(t, "")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		assert.NoError(t, err)
		_, err = p.CurrentCmd()
		assert.Error(t, err)
	})

	t.Run("Advancing over the limit always gives last cmd", func(t *testing.T) {
		file := createTestFile(t, "first line")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		assert.NoError(t, err)
		p.Advance()
		p.Advance()
		p.Advance()
		cmd, err := p.CurrentCmd()
		assert.NoError(t, err)
		assert.Equal(t, "first line", cmd)
	})
}

func createTestFile(t *testing.T, content string) *os.File {
	file, err := os.CreateTemp("", "test.asm")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte(content))
	return file
}
