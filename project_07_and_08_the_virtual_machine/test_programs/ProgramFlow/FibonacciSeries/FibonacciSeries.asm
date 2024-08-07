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

//pop pointer 1
//pop pointer that
// *SP=THAT
@SP
AM=M-1
D=M

@R4
M=D

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

//pop that 0
//*R13=*(THAT + 0)
@0
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

//pop that 1
//*R13=*(THAT + 1)
@1
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

//push constant 2
//*SP=2
@2
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

( FibonacciSeries.MAIN_LOOP_START )

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

// if-goto FibonacciSeries.COMPUTE_ELEMENT
@SP
AM=M-1
D=M
@FibonacciSeries.COMPUTE_ELEMENT
D;JNE

// goto FibonacciSeries.END_PROGRAM
@FibonacciSeries.END_PROGRAM
0;JEQ

( FibonacciSeries.COMPUTE_ELEMENT )

//push that 0
//D=*(THAT + 0)
@0
D=A
@R4								// *R4=THAT
A=M+D
D=M

// *SP=*(THAT + 0)
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push that 1
//D=*(THAT + 1)
@1
D=A
@R4								// *R4=THAT
A=M+D
D=M

// *SP=*(THAT + 1)
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

//add
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D

//pop pointer 1
//pop pointer that
// *SP=THAT
@SP
AM=M-1
D=M

@R4
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

// goto FibonacciSeries.MAIN_LOOP_START
@FibonacciSeries.MAIN_LOOP_START
0;JEQ

( FibonacciSeries.END_PROGRAM )

