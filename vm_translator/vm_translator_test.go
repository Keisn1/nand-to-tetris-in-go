package vmtrans_test

import (
	vmtrans "hack/vm_translator"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_CodeWriter_Branching(t *testing.T) {
// 	t.Run("Write Arithmetic commands to output file", func(t *testing.T) {
// 		type testCase struct {
// 			name                string
// 			cmdType, arg1, arg2 string
// 			want                string
// 		}

// 		filename := "./test.asm"
// 		testCases := []testCase{
// 			// {
// 			// 	name:    "test label command",
// 			// 	cmdType: vmtrans.C_LABEL, arg1: "LOOP_2", arg2: "",
// 			// 	want: `( test.LOOP_2 )`,
// 			// },
// 			// {
// 			// 	name:    "test label command",
// 			// 	cmdType: vmtrans.C_LABEL, arg1: "LOOP_START", arg2: "",
// 			// 	want: `( test.LOOP_START )`,
// 			// },
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				cw := vmtrans.NewCodeWriter(filename)
// 				cw.WriteArithmetic(tc.cmdType, tc.arg1, tc.arg2)

// 				got := MustReadFile(t, filename)
// 				assert.Equal(t, tc.want, got)
// 				os.Remove(filename)
// 			})
// 		}
// 	})
// }

func Test_CodeWriter_Arithtmetic(t *testing.T) {
	t.Run("Write Arithmetic commands to output file", func(t *testing.T) {
		type testCase struct {
			name                string
			cmdType, arg1, arg2 string
			want                string
		}

		testCases := []testCase{
			{
				name:    "test not",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "not", arg2: "",
				want: `//not
@SP
A=M-1
M=!M
`,
			},
			{
				name:    "test neg",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "neg", arg2: "",
				want: `//neg
@SP
A=M-1
M=-M
`,
			},
			{
				name:    "test gt",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "gt", arg2: "",
				want: `//gt
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_GREATER_THAN_COUNTER_0
D;JGT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_GREATER_THAN_COUNTER_0
0;JMP

( IS_GREATER_THAN_COUNTER_0 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_GREATER_THAN_COUNTER_0 )
`,
			},
			{
				name:    "test lt",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "lt", arg2: "",
				want: `//lt
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_LOWER_THAN_COUNTER_0
D;JLT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_LOWER_THAN_COUNTER_0
0;JMP

( IS_LOWER_THAN_COUNTER_0 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_LOWER_THAN_COUNTER_0 )
`,
			},
			{
				name:    "test eq",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "eq", arg2: "",
				want: `//eq
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_EQUAL_COUNTER_0
D;JEQ

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_EQUAL_COUNTER_0
0;JMP

( IS_EQUAL_COUNTER_0 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_EQUAL_COUNTER_0 )
`,
			},
			{
				name:    "test or",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "or", arg2: "",
				want: `//or
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M|D
`,
			},
			{
				name:    "test and",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "and", arg2: "",
				want: `//and
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M&D
`,
			},
			{
				name:    "test sub",
				cmdType: vmtrans.C_ARITHMETIC, arg1: "sub", arg2: "",
				want: `//sub
//SP--; D=M
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
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D
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

func Test_CodeWriter_Segments(t *testing.T) {
	t.Run("Write Arithmetic commands to output file", func(t *testing.T) {
		type testCase struct {
			name                string
			cmdType, arg1, arg2 string
			want                string
		}

		testCases := []testCase{
			{
				name:    "test pop pointer 1",
				cmdType: vmtrans.C_POP, arg1: "pointer", arg2: "1",
				want: `//pop pointer 1
//pop pointer that
// *SP=THAT
@SP
AM=M-1
D=M

@R4
M=D
`,
			},
			{
				name:    "test pop pointer 0",
				cmdType: vmtrans.C_POP, arg1: "pointer", arg2: "0",
				want: `//pop pointer 0
//pop pointer this
// *SP=THIS
@SP
AM=M-1
D=M

@R3
M=D
`,
			},
			{
				name:    "test push pointer 1",
				cmdType: vmtrans.C_PUSH, arg1: "pointer", arg2: "1",
				want: `//push pointer 1
//push pointer that
// *SP=THAT
@R4
D=M
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
			{
				name:    "test push pointer 0",
				cmdType: vmtrans.C_PUSH, arg1: "pointer", arg2: "0",
				want: `//push pointer 0
//push pointer this
// *SP=THIS
@R3
D=M
@SP
A=M
M=D

// *SP++
@SP
M=M+1
`,
			},
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
