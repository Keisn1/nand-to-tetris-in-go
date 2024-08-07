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

//pop static StaticTest.8
// D=*SP--
@SP
AM=M-1
D=M

@StaticTest.8
M=D

//pop static StaticTest.3
// D=*SP--
@SP
AM=M-1
D=M

@StaticTest.3
M=D

//pop static StaticTest.1
// D=*SP--
@SP
AM=M-1
D=M

@StaticTest.1
M=D

//push static StaticTest.3
// D=StaticTest.3
@StaticTest.3
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push static StaticTest.1
// D=StaticTest.1
@StaticTest.1
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//sub
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M-D

//push static StaticTest.8
// D=StaticTest.8
@StaticTest.8
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
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

