//push temp {{.x}}
//D=*(5 + {{.x}})
@{{.x}}
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
