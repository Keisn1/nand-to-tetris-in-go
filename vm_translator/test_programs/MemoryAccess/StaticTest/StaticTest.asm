//push constant 111
//*SP=111
@111
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 333
//*SP=333
@333
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 888
//*SP=888
@888
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop static 8
// D=*SP--
@SP
AM=M-1
D=M

@STATICTEST.8
M=D

//pop static 3
// D=*SP--
@SP
AM=M-1
D=M

@STATICTEST.3
M=D

//pop static 1
// D=*SP--
@SP
AM=M-1
D=M

@STATICTEST.1
M=D

//push static 3
// D=STATICTEST.3
@STATICTEST.3
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push static 1
// D=STATICTEST.1
@STATICTEST.1
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//sub
//SP--; *R13=*SP
@SP
AM=M-1
D=M
@SP
A=M-1
M=M-D

//push static 8
// D=STATICTEST.8
@STATICTEST.8
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//add
//SP--; *R13=*SP
@SP
AM=M-1
D=M
@SP
A=M-1
M=M+D

