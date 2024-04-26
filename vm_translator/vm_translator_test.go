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
