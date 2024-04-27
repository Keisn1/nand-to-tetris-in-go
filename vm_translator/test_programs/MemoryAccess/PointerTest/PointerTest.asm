//push constant 3030
//*SP=3030
@3030
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop pointer 0
//pop pointer this
// *SP=THIS
@SP
AM=M-1
D=M

@R3
M=D

//push constant 3040
//*SP=3040
@3040
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop pointer 1
//pop pointer that
// *SP=THAT
@SP
AM=M-1
D=M

@R4
M=D

//push constant 32
//*SP=32
@32
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop this 2
//*R13=*(THIS + 2)
@2
D=A
@R3								// *R3=THIS
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

//push constant 46
//*SP=46
@46
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop that 6
//*R13=*(THAT + 6)
@6
D=A
@R4								// *R4=THAT
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

//push pointer 0
//push pointer this
// *SP=THIS
@R3
D=M
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push pointer 1
//push pointer that
// *SP=THAT
@R4
D=M
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

//push this 2
//D=*(THIS + 2)
@2
D=A
@R3								// *R3=THIS
A=M+D
D=M

// *SP=*(THIS + 2)
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

//push that 6
//D=*(THAT + 6)
@6
D=A
@R4								// *R4=THAT
A=M+D
D=M

// *SP=*(THAT + 6)
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

