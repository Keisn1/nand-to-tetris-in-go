// SP=256
@256
D=A
@R0
M=D

// call Sys.init 0
// push return address
@Boot$ret.0
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

// ARG = SP - 5 - 0 (=nArgs)
@SP
D=M

@5
D=D-A

@0 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Sys.init
@Sys.init
0;JEQ

// (returnAddress)
( Boot$ret.0 )
// function Main.fibonacci 0
( Main.fibonacci )

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

//lt
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_LOWER_THAN_COUNTER_0
D;JLT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_LOWER_THAN_COUNTER_0
0;JMP

( IS_LOWER_THAN_COUNTER_0 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_LOWER_THAN_COUNTER_0 )

// if-goto Main.IF_TRUE
@SP
AM=M-1
D=M
@Main.IF_TRUE
D;JNE

// goto Main.IF_FALSE
@Main.IF_FALSE
0;JEQ

( Main.IF_TRUE )

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

// return
// endframe = LCL
@R1
D=M
@endframe
M=D

// retAddr = *(endframe-5)
// D=*(endframe-5)
@5
A=D-A
D=M

// retAddr = D
@retAddr
M=D

// *ARG=pop()
// D=*(SP--)
@SP
AM=M-1
D=M

// *ARG=D
@2
A=M
M=D

// SP=ARG+1
@2
D=M
@SP
M=D+1

// THAT=*(endframe-1)
@endframe
AM=M-1
D=M
@4
M=D

// THIS=*(endframe-2)
@endframe
AM=M-1
D=M
@3
M=D

// ARG=*(endframe-3)
@endframe
AM=M-1
D=M
@2
M=D

// LCL=*(endframe-4)
@endframe
AM=M-1
D=M
@R1
M=D

// goto retAddr
@retAddr
A=M
0;JEQ

( Main.IF_FALSE )

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

// call Main.fibonacci 1
// push return address
@Main$ret.0
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

// ARG = SP - 5 - 1 (=nArgs)
@SP
D=M

@5
D=D-A

@1 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Main.fibonacci
@Main.fibonacci
0;JEQ

// (returnAddress)
( Main$ret.0 )

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

// call Main.fibonacci 1
// push return address
@Main$ret.1
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

// ARG = SP - 5 - 1 (=nArgs)
@SP
D=M

@5
D=D-A

@1 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Main.fibonacci
@Main.fibonacci
0;JEQ

// (returnAddress)
( Main$ret.1 )

//add
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M+D

// return
// endframe = LCL
@R1
D=M
@endframe
M=D

// retAddr = *(endframe-5)
// D=*(endframe-5)
@5
A=D-A
D=M

// retAddr = D
@retAddr
M=D

// *ARG=pop()
// D=*(SP--)
@SP
AM=M-1
D=M

// *ARG=D
@2
A=M
M=D

// SP=ARG+1
@2
D=M
@SP
M=D+1

// THAT=*(endframe-1)
@endframe
AM=M-1
D=M
@4
M=D

// THIS=*(endframe-2)
@endframe
AM=M-1
D=M
@3
M=D

// ARG=*(endframe-3)
@endframe
AM=M-1
D=M
@2
M=D

// LCL=*(endframe-4)
@endframe
AM=M-1
D=M
@R1
M=D

// goto retAddr
@retAddr
A=M
0;JEQ

// function Sys.init 0
( Sys.init )

//push constant 4
//*SP=4
@4
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

// call Main.fibonacci 1
// push return address
@Sys$ret.0
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

// ARG = SP - 5 - 1 (=nArgs)
@SP
D=M

@5
D=D-A

@1 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Main.fibonacci
@Main.fibonacci
0;JEQ

// (returnAddress)
( Sys$ret.0 )

( Sys.WHILE )

// goto Sys.WHILE
@Sys.WHILE
0;JEQ

