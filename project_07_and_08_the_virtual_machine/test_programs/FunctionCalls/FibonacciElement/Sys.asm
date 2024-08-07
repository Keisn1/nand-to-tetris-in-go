// function Sys.init 0
( Sys.init )

//push constant 4
//*SP=4
@4
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

// call Main.fibonacci 1
// push return address
@Sys$ret.0
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
( Sys$ret.0 )

( Sys.WHILE )

// goto Sys.WHILE
@Sys.WHILE
0;JEQ

