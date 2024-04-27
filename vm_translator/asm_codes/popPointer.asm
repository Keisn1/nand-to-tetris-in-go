//pop pointer {{.arg}}
//pop pointer {{.segment}}
// *SP={{.segment_register_name}}
@SP
AM=M-1
D=M

@{{.segment_register}}
M=D
