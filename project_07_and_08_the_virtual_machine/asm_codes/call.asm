// call {{ .function_name }} {{ .n_args }}
// push return address
@{{ .caller }}$ret.{{ .function_counter }}
D=A

@SP
A=M
M=D

@SP
M=M+1

// push LCL
@R1
D=M

@SP
A=M
M=D

@SP
M=M+1

// push ARG
@R2
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THIS
@R3
D=M

@SP
A=M
M=D

@SP
M=M+1

// push THAT
@R4
D=M

@SP
A=M
M=D

@SP
M=M+1

// ARG = SP - 5 - {{ .n_args }} (=nArgs)
@SP
D=M

@5
D=D-A

@{{ .n_args }} // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto {{ .function_name }}
@{{ .function_name }}
0;JEQ

// (returnAddress)
( {{ .caller }}$ret.{{ .function_counter }} )
