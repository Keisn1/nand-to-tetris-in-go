//push constant 7
//*SP=7
@7
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 8
//*SP=8
@8
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//add
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D

