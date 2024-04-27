//{{.comp}}
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@{{.comp_verbose}}_COUNTER_{{.counter}}
D;{{.comp_operator}}

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_{{.comp_verbose}}_COUNTER_{{.counter}}
0;JMP

( {{.comp_verbose}}_COUNTER_{{.counter}} )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_{{.comp_verbose}}_COUNTER_{{.counter}} )
