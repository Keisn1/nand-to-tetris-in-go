
//push constant {{ .x }}
//*SP={{ .x }}
@{{ .x }}
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1
