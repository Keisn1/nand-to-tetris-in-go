//push static {{.x}}
// D={{.filename}}.{{.x}}
@{{.filename}}.{{.x}}
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1
