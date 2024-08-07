//push {{.segment}} {{.x}}
//D=*({{.segment_register_name}} + {{.x}})
@{{.x}}
D=A
@{{.segment_register}}								// *{{.segment_register}}={{.segment_register_name}}
A=M+D
D=M

// *SP=*({{.segment_register_name}} + {{.x}})
@SP
A=M
M=D

// *SP++
@SP
M=M+1
