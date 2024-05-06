// function SimpleFunction.SimpleFunction.test 2
( SimpleFunction.SimpleFunction.test )

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

//push local 0
//D=*(LCL + 0)
@0
D=A
@R1								// *R1=LCL
A=M+D
D=M

// *SP=*(LCL + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push local 1
//D=*(LCL + 1)
@1
D=A
@R1								// *R1=LCL
A=M+D
D=M

// *SP=*(LCL + 1)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//add
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D

//not
@SP
A=M-1
M=!M

//push argument 0
//D=*(ARG + 0)
@0
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//add
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D

//push argument 1
//D=*(ARG + 1)
@1
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 1)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//sub
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M-D

// return
// endframe = LCL
@1
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
@1
M=D

// goto retAddr
@retAddr
A=M
0;JEQ

