package vmtrans_test

import (
	vmtrans "hack/vm_translator"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CodeWriter(t *testing.T) {
	t.Run("Write Arithmetic commands to output file", func(t *testing.T) {
		type testCase struct {
			name                string
			cmdType, arg1, arg2 string
			want                string
		}

		testCases := []testCase{
			{
				name:    "test push static x",
				cmdType: vmtrans.C_PUSH, arg1: "static", arg2: "2",
				want: `//push static 2
// D=TEST.2
@TEST.2
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test pop static x",
				cmdType: vmtrans.C_POP, arg1: "static", arg2: "2",
				want: `//pop static 2
// D=*SP--
@SP
AM=M-1
D=M

@TEST.2
M=D
`,
			},
			{
				name:    "test sub",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "sub", arg2: "",
				want: `//sub
//SP--; *R13=*SP
@SP
AM=M-1
D=M
@SP
A=M-1
M=M-D
`,
			},
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
				cmdType: vmtrans.C_PUSH, arg1: "that", arg2: "2",
				want: `//push that 2
//D=*(THAT + 2)
@2
D=A
@R4								// *R4=THAT
A=M+D
D=M

// *SP=*(THAT + 2)
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test pop that x",
				cmdType: vmtrans.C_POP, arg1: "that", arg2: "10",
				want: `//pop that 10
//*R13=*(THAT + 10)
@10
D=A
@R4								// *R4=THAT
D=M+D
@R13
M=D

// *SP--
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
				cmdType: vmtrans.C_PUSH, arg1: "this", arg2: "2",
				want: `//push this 2
//D=*(THIS + 2)
@2
D=A
@R3								// *R3=THIS
A=M+D
D=M

// *SP=*(THIS + 2)
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test pop this x",
				cmdType: vmtrans.C_POP, arg1: "this", arg2: "10",
				want: `//pop this 10
//*R13=*(THIS + 10)
@10
D=A
@R3								// *R3=THIS
D=M+D
@R13
M=D

// *SP--
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
				cmdType: vmtrans.C_PUSH, arg1: "argument", arg2: "2",
				want: `//push argument 2
//D=*(ARG + 2)
@2
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 2)
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test pop argument x",
				cmdType: vmtrans.C_POP, arg1: "argument", arg2: "10",
				want: `//pop argument 10
//*R13=*(ARG + 10)
@10
D=A
@R2								// *R2=ARG
D=M+D
@R13
M=D

// *SP--
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
//D=*(LCL + 2)
@2
D=A
@R1								// *R1=LCL
A=M+D
D=M

// *SP=*(LCL + 2)
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test pop local x",
				cmdType: vmtrans.C_POP, arg1: "local", arg2: "2",
				want: `//pop local 2
//*R13=*(LCL + 2)
@2
D=A
@R1								// *R1=LCL
D=M+D
@R13
M=D

// *SP--
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
				name:    "test push temp x",
				cmdType: vmtrans.C_PUSH, arg1: "temp", arg2: "2",
				want: `//push temp 2
//D=*(5 + 2)
@2
D=A
@5 								// temp base
A=D+A
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test pop temp x",
				cmdType: vmtrans.C_POP, arg1: "temp", arg2: "2",
				want: `//pop temp 2
//*R13=5 + 2
@2 								// offset
D=A
@5 								// temp base
D=D+A
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
