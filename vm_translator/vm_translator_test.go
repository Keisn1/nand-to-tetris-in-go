package vmtrans_test

import (
	vmtrans "hack/vm_translator"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CodeWriter(t *testing.T) {
	t.Run("Write Arithmetic commands to output file 2", func(t *testing.T) {
		type testCase struct {
			name                string
			cmdType, arg1, arg2 string
			want                string
		}

		testCases := []testCase{
			{
				name:    "test add",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "add", arg2: "",
				want: `//add
//SP--; *R13=*SP
@SP
AM=M-1
D=M
@SP
A=M-1
M=M+D
`,
			},
			{
				name:    "test push that x",
				cmdType: vmtrans.C_PUSH, arg1: "that", arg2: "5",
				want: `//push that 5
//*R13=THAT + 5
@5
D=A
@R4
D=M+D
@R13
M=D

// SP--
@SP
AM=M-1
D=M

// *R13=*SP
@R13
A=M
M=D
`,
			},
			{
				name:    "test push this x",
				cmdType: vmtrans.C_PUSH, arg1: "this", arg2: "5",
				want: `//push this 5
//*R13=THIS + 5
@5
D=A
@R3
D=M+D
@R13
M=D

// SP--
@SP
AM=M-1
D=M

// *R13=*SP
@R13
A=M
M=D
`,
			},
			{
				name:    "test push argument x",
				cmdType: vmtrans.C_PUSH, arg1: "argument", arg2: "5",
				want: `//push argument 5
//*R13=ARG + 5
@5
D=A
@R2
D=M+D
@R13
M=D

// SP--
@SP
AM=M-1
D=M

// *R13=*SP
@R13
A=M
M=D
`,
			},
			{
				name:    "test push local x",
				cmdType: vmtrans.C_PUSH, arg1: "local", arg2: "2",
				want: `//push local 2
//*R13=LCL + 2
@2
D=A
@R1
D=M+D
@R13
M=D

// SP--
@SP
AM=M-1
D=M

// *R13=*SP
@R13
A=M
M=D
`,
			},
			{
				name:    "test push constant x",
				cmdType: vmtrans.C_PUSH, arg1: "constant", arg2: "4",
				want: `//push constant 4
//*SP=4
@4
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1
`,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				filename := "./test.asm"
				cw := vmtrans.NewCodeWriter(filename)
				cw.WriteArithmetic(tc.cmdType, tc.arg1, tc.arg2)

				got := MustReadFile(t, filename)
				assert.Equal(t, tc.want, got)
				os.Remove(filename)
			})
		}

	})
}

func MustReadFile(t *testing.T, fp string) string {
	t.Helper()
	c, err := os.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}
	return string(c)

}

func Test_Parser(t *testing.T) {
	t.Run("Parsing command types", func(t *testing.T) {
		type testCase struct {
			cmd           string
			want_c_type   string
			want_arg1     string
			want_arg2     string
			want_err_arg2 bool
		}

		testCases := []testCase{
			{
				cmd:           "add",
				want_c_type:   "C_ARITHMETIC",
				want_arg1:     "add",
				want_err_arg2: true,
			},
			{
				cmd:           "sub",
				want_c_type:   "C_ARITHMETIC",
				want_arg1:     "sub",
				want_err_arg2: true,
			},
			{
				cmd:         "push local 2",
				want_c_type: "C_PUSH",
				want_arg1:   "local",
				want_arg2:   "2",
			},
			{
				cmd:         "push local 3",
				want_c_type: "C_PUSH",
				want_arg1:   "local",
				want_arg2:   "3",
			},
			{
				cmd:         "pop temp 0",
				want_c_type: "C_POP",
				want_arg1:   "temp",
				want_arg2:   "0",
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
			got, _ := p.CommandType()
			assert.Equal(t, tc.want_c_type, got)
			got, _ = p.Arg1()
			assert.Equal(t, tc.want_arg1, got)

			if tc.want_err_arg2 {
				_, err = p.Arg2()
				assert.Error(t, err)
			} else {
				got, _ = p.Arg2()
				assert.Equal(t, tc.want_arg2, got)
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

	t.Run("Happy path, reading files", func(t *testing.T) {
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
				if err != nil {
					t.Fatal(err)
				}
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
