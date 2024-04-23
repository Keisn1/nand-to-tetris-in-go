package vmtrans

const (
	pushLocalX = `//push local {{ .x }}
//*R13=LCL + {{ .x }}
@{{ .x }}
D=A
@R1
D=M+D
@R13
M=D

// SP--
@SP
M=M-1

// *R13=*SP
@R13`
)
