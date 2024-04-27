//pop {{.segment}} {{.x}}
//*R13=*({{.segment_register_name}} + {{.x}})
@{{.x}}
D=A
@{{.segment_register}}								// *{{.segment_register}}={{.segment_register_name}}
D=M+D
@R13
M=D

// *SP--
@SP
AM=M-1
D=M

// *R13=*SP
@R13
A=M
M=D
