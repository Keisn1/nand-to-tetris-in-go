// SP=256
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
