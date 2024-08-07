// return
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
