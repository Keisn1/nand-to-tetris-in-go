//push pointer {{.arg}}
//push pointer {{.segment}}
// *SP={{.segment_register_name}}
@{{.segment_register}}
D=M
@SP
A=M
M=D

// *SP++
@SP
M=M+1
