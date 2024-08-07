//push constant 10
//*SP=10
@10
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop local 0
//*R13=*(LCL + 0)
@0
D=A
@R1								// *R1=LCL
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

//push constant 21
//*SP=21
@21
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 22
//*SP=22
@22
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop argument 2
//*R13=*(ARG + 2)
@2
D=A
@R2								// *R2=ARG
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

//pop argument 1
//*R13=*(ARG + 1)
@1
D=A
@R2								// *R2=ARG
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

//push constant 36
//*SP=36
@36
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop this 6
//*R13=*(THIS + 6)
@6
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

//push constant 42
//*SP=42
@42
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 45
//*SP=45
@45
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop that 5
//*R13=*(THAT + 5)
@5
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

//pop that 2
//*R13=*(THAT + 2)
@2
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

//push constant 510
//*SP=510
@510
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//pop temp 6
//*R13=5 + 6
@6 								// offset
D=A
@5 								// temp base
D=D+A
@R13
M=D

// SP--
@SP
AM=M-1
D=M

// *R13=*SP
@R13
A=M
M=D

//push local 0
//D=*(LCL + 0)
@0
D=A
@R1								// *R1=LCL
A=M+D
D=M

// *SP=*(LCL + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push that 5
//D=*(THAT + 5)
@5
D=A
@R4								// *R4=THAT
A=M+D
D=M

// *SP=*(THAT + 5)
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

//push argument 1
//D=*(ARG + 1)
@1
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 1)
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

//push this 6
//D=*(THIS + 6)
@6
D=A
@R3								// *R3=THIS
A=M+D
D=M

// *SP=*(THIS + 6)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push this 6
//D=*(THIS + 6)
@6
D=A
@R3								// *R3=THIS
A=M+D
D=M

// *SP=*(THIS + 6)
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

//sub
//SP--; *R13=*SP
@SP
AM=M-1
D=M
@SP
A=M-1
M=M-D

//push temp 6
//D=*(5 + 6)
@6
D=A
@5 								// temp base
A=D+A
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

