//pop static {{.x}}
// D=*SP--
@SP
AM=M-1
D=M

@{{.filename}}.{{.x}}
M=D
