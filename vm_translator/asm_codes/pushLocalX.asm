
//push local {{.x}}
//*R13=LCL + {{.x}}
@{{.x}}
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
