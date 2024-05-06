//push static {{.x}}
// D={{.x}}
@{{.x}}
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1
