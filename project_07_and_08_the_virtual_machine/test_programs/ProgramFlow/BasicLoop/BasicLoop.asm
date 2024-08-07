//push constant 0
//*SP=0
@0
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

( BASICLOOP.LOOP_START )

//push argument 0
//D=*(ARG + 0)
@0
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

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

//add
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D

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

//push argument 0
//D=*(ARG + 0)
@0
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push constant 1
//*SP=1
@1
D=A
@SP
A=M
M=D

//SP++
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

//pop argument 0
//*R13=*(ARG + 0)
@0
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

//push argument 0
//D=*(ARG + 0)
@0
D=A
@R2								// *R2=ARG
A=M+D
D=M

// *SP=*(ARG + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

// if-goto LOOP_START
@SP
AM=M-1
D=M
@BASICLOOP.LOOP_START
D;JNE

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

