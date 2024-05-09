package vmtrans_test

import (
	vmtrans "hack/vm_translator"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CodeWriter_Function(t *testing.T) {
	t.Run("Test boot code", func(t *testing.T) {
		want := `// SP=256
@256
D=A
@R0
M=D

// call Sys.init 0
// push return address
@Boot$ret.0
D=A

@SP
A=M
M=D

@SP
M=M+1

// push LCL
@R1
D=M

@SP
A=M
M=D

@SP
M=M+1

// push ARG
@R2
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THIS
@R3
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THAT
@R4
D=M

@SP
A=M
M=D

@SP
M=M+1

// ARG = SP - 5 - 0 (=nArgs)
@SP
D=M

@5
D=D-A

@0 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Sys.init
@Sys.init
0;JEQ

// (returnAddress)
( Boot$ret.0 )
`
		caller := "test.asm"
		filename := "./out.asm"
		cw := vmtrans.NewCodeWriter(filename, caller)
		cw.WriteBootStrap()
		got := MustReadFile(t, filename)
		assert.Equal(t, want, got)
		os.Remove(filename)
	})

	t.Run("Write function/return/call commands to output file", func(t *testing.T) {
		type testCase struct {
			name                string
			cmdType, arg1, arg2 string
			caller              string
			want                string
		}

		testCases := []testCase{
			{
				name:    "call",
				cmdType: vmtrans.C_CALL, arg1: "Foo.recursion", arg2: "4",
				caller: "test.vm",
				want: `// call Foo.recursion 4
// push return address
@Test$ret.0
D=A

@SP
A=M
M=D

@SP
M=M+1

// push LCL
@R1
D=M

@SP
A=M
M=D

@SP
M=M+1

// push ARG
@R2
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THIS
@R3
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THAT
@R4
D=M

@SP
A=M
M=D

@SP
M=M+1

// ARG = SP - 5 - 4 (=nArgs)
@SP
D=M

@5
D=D-A

@4 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Foo.recursion
@Foo.recursion
0;JEQ

// (returnAddress)
( Test$ret.0 )
`,
			},
			{
				name:    "call",
				cmdType: vmtrans.C_CALL, arg1: "Main.fibonacci", arg2: "1",
				caller: "test.vm",
				want: `// call Main.fibonacci 1
// push return address
@Test$ret.0
D=A

@SP
A=M
M=D

@SP
M=M+1

// push LCL
@R1
D=M

@SP
A=M
M=D

@SP
M=M+1

// push ARG
@R2
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THIS
@R3
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THAT
@R4
D=M

@SP
A=M
M=D

@SP
M=M+1

// ARG = SP - 5 - 1 (=nArgs)
@SP
D=M

@5
D=D-A

@1 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Main.fibonacci
@Main.fibonacci
0;JEQ

// (returnAddress)
( Test$ret.0 )
`,
			},
			{
				name:    "return",
				cmdType: vmtrans.C_RETURN, arg1: "", arg2: "",
				want: `// return
// endframe = LCL
@R1
D=M
@endframe
M=D

// retAddr = *(endframe-5)
// D=*(endframe-5)
@5
A=D-A
D=M

// retAddr = D
@retAddr
M=D

// *ARG=pop()
// D=*(SP--)
@SP
AM=M-1
D=M

// *ARG=D
@2
A=M
M=D

// SP=ARG+1
@2
D=M
@SP
M=D+1

// THAT=*(endframe-1)
@endframe
AM=M-1
D=M
@4
M=D

// THIS=*(endframe-2)
@endframe
AM=M-1
D=M
@3
M=D

// ARG=*(endframe-3)
@endframe
AM=M-1
D=M
@2
M=D

// LCL=*(endframe-4)
@endframe
AM=M-1
D=M
@R1
M=D

// goto retAddr
@retAddr
A=M
0;JEQ
`,
			},
			{
				name:    "function SimpleFunction.test 0",
				cmdType: vmtrans.C_FUNCTION, arg1: "SimpleFunction.test", arg2: "0",
				want: `// function SimpleFunction.test 0
( SimpleFunction.test )
`,
			},
			{
				name:    "function AnotherFunction.test 3",
				cmdType: vmtrans.C_FUNCTION, arg1: "AnotherFunction.test", arg2: "3",
				want: `// function AnotherFunction.test 3
( AnotherFunction.test )

// Iteration #0
// D=0 and *SP=D
@0
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

// Iteration #1
// D=0 and *SP=D
@0
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

// Iteration #2
// D=0 and *SP=D
@0
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1
`,
			},
			{
				name:    "function SimpleFunction.test 2",
				cmdType: vmtrans.C_FUNCTION, arg1: "SimpleFunction.test", arg2: "2",
				want: `// function SimpleFunction.test 2
( SimpleFunction.test )

// Iteration #0
// D=0 and *SP=D
@0
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

// Iteration #1
// D=0 and *SP=D
@0
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
				caller := tc.caller
				filename := "./out.asm"
				cw := vmtrans.NewCodeWriter(filename, caller)
				cw.Write(tc.cmdType, tc.arg1, tc.arg2)

				got := MustReadFile(t, filename)
				assert.Equal(t, tc.want, got)
				os.Remove(filename)
			})
		}
	})
}

func Test_CodeWriter_Branching(t *testing.T) {
	t.Run("Write Arithmetic commands to output file", func(t *testing.T) {
		type testCase struct {
			name                string
			cmdType, arg1, arg2 string
			want                string
		}

		testCases := []testCase{
			{
				name:    "goto MAIN_LOOP_START",
				cmdType: vmtrans.C_GOTO, arg1: "MAIN_LOOP_START", arg2: "",
				want: `// goto MAIN_LOOP_START
@MAIN_LOOP_START
0;JEQ
`,
			},
			{
				name:    "test if-goto command",
				cmdType: vmtrans.C_IF, arg1: "LOOP_START", arg2: "",
				want: `// if-goto LOOP_START
@SP
AM=M-1
D=M
@LOOP_START
D;JNE
`,
			},
			{
				name:    "test label command",
				cmdType: vmtrans.C_LABEL, arg1: "LOOP_START", arg2: "",
				want: `( LOOP_START )
`,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				filename := "./test.asm"
				cw := vmtrans.NewCodeWriter(filename, "")
				cw.Write(tc.cmdType, tc.arg1, tc.arg2)

				got := MustReadFile(t, filename)
				assert.Equal(t, tc.want, got)
				os.Remove(filename)
			})
		}
	})
}

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
				cw := vmtrans.NewCodeWriter(filename, "")
				cw.Write(tc.cmdType, tc.arg1, tc.arg2)

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
				name:    "test push static Test.x",
				cmdType: vmtrans.C_PUSH, arg1: "static", arg2: "Test.2",
				want: `//push static Test.2
// D=Test.2
@Test.2
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
				name:    "test pop static Test.x",
				cmdType: vmtrans.C_POP, arg1: "static", arg2: "Test.2",
				want: `//pop static Test.2
// D=*SP--
@SP
AM=M-1
D=M

@Test.2
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
				cw := vmtrans.NewCodeWriter(filename, "")
				cw.Write(tc.cmdType, tc.arg1, tc.arg2)

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
