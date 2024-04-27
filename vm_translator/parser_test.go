package vmtrans_test

import (
	vmtrans "hack/vm_translator"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser(t *testing.T) {
	t.Run("Parsing command types", func(t *testing.T) {
		type testCase struct {
			cmd           string
			want_cmd_type string
			want_arg1     string
			want_arg2     string
			want_err_arg2 bool
		}

		testCases := []testCase{
			{
				cmd:           "add",
				want_cmd_type: "C_ARITHMETIC",
				want_arg1:     "add",
				want_err_arg2: true,
			},
			{
				cmd:           "add",
				want_cmd_type: "C_ARITHMETIC",
				want_arg1:     "add",
				want_err_arg2: true,
			},
			{
				cmd:           "sub",
				want_cmd_type: "C_ARITHMETIC",
				want_arg1:     "sub",
				want_err_arg2: true,
			},
			{
				cmd:           "push local 2",
				want_cmd_type: "C_PUSH",
				want_arg1:     "local",
				want_arg2:     "2",
			},
			{
				cmd:           "   push    local    2   ",
				want_cmd_type: "C_PUSH",
				want_arg1:     "local",
				want_arg2:     "2",
			},
			{
				cmd:           "push local 3",
				want_cmd_type: "C_PUSH",
				want_arg1:     "local",
				want_arg2:     "3",
			},
			{
				cmd:           "pop temp 0",
				want_cmd_type: "C_POP",
				want_arg1:     "temp",
				want_arg2:     "0",
			},
		}
		for _, tc := range testCases {
			file := createTestFile(t, tc.cmd)
			defer os.Remove(file.Name())

			p, err := vmtrans.NewParser(file.Name())
			if err != nil {
				t.Fatal(err)
			}

			p.Advance()
			gotCmdType, _ := p.CommandType()
			assert.Equal(t, tc.want_cmd_type, gotCmdType)
			gotArg1, _ := p.Arg1()
			assert.Equal(t, tc.want_arg1, gotArg1)

			if tc.want_err_arg2 {
				_, err = p.Arg2()
				assert.Error(t, err)
			} else {
				gotCmdType, _ = p.Arg2()
				assert.Equal(t, tc.want_arg2, gotCmdType)
			}
		}
	})

	t.Run("Error commandType", func(t *testing.T) {
		file := createTestFile(t, "")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		if err != nil {
			t.Fatal(err)
		}
		_, err = p.CommandType()
		assert.ErrorAs(t, err, &vmtrans.ErrNotAdvanced)
	})

	t.Run("Reading lines and put them into commands slice", func(t *testing.T) {
		type testCase struct {
			name     string
			content  string
			wantCmds []string
		}

		testCases := []testCase{
			{
				name:     "one line",
				content:  "first line",
				wantCmds: []string{"first line"},
			},
			{
				name: "two lines",
				content: `first line
second line`,
				wantCmds: []string{"first line", "second line"},
			},
			{
				name: "skipping empty line",
				content: `first line

second line`,
				wantCmds: []string{"first line", "second line"},
			},
			{
				name: "skipping comment line",
				content: `first line
// comment
second line`,
				wantCmds: []string{"first line", "second line"},
			},
			{
				name:     "Cuts comment from end of line",
				content:  `first line     // some comment`,
				wantCmds: []string{"first line"},
			},
			{
				name:     "Ignores whitespace in front",
				content:  `    first     line     // comment`,
				wantCmds: []string{"first line"},
			},
			{
				name: "Ignore newlines at the end",
				content: `first line

`,
				wantCmds: []string{"first line"},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				file := createTestFile(t, tc.content)
				defer os.Remove(file.Name())

				p, err := vmtrans.NewParser(file.Name())
				if err != nil {
					t.Fatal(err)
				}
				var got []string
				for p.HasMoreCommands() {
					p.Advance()
					curCmd, _ := p.CurrentCmd()
					got = append(got, curCmd)
				}

				assert.Equal(t, tc.wantCmds, got)
			})
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		file := createTestFile(t, "")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, p.HasMoreCommands())
	})

	t.Run("Reading before advancing", func(t *testing.T) {
		file := createTestFile(t, "")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		if err != nil {
			t.Fatal(err)
		}
		_, err = p.CurrentCmd()
		assert.ErrorAs(t, err, &vmtrans.ErrNotAdvanced)
	})

	t.Run("Advancing over the limit always gives last cmd", func(t *testing.T) {
		file := createTestFile(t, "first line")
		defer os.Remove(file.Name())

		p, err := vmtrans.NewParser(file.Name())
		if err != nil {
			t.Fatal(err)
		}
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
