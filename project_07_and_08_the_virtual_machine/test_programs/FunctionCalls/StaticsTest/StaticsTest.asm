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
// function Class1.set 0
( Class1.set )

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

//pop static Class1.0
// D=*SP--
@SP
AM=M-1
D=M

@Class1.0
M=D

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

//pop static Class1.1
// D=*SP--
@SP
AM=M-1
D=M

@Class1.1
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

// function Class1.get 0
( Class1.get )

//push static Class1.0
// D=Class1.0
@Class1.0
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push static Class1.1
// D=Class1.1
@Class1.1
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

// function Class2.set 0
( Class2.set )

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

//pop static Class2.0
// D=*SP--
@SP
AM=M-1
D=M

@Class2.0
M=D

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

//pop static Class2.1
// D=*SP--
@SP
AM=M-1
D=M

@Class2.1
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

// function Class2.get 0
( Class2.get )

//push static Class2.0
// D=Class2.0
@Class2.0
D=M

// *SP=D
@SP
A=M
M=D

// *SP++
@SP
M=M+1

//push static Class2.1
// D=Class2.1
@Class2.1
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

//push constant 6
//*SP=6
@6
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

// call Class1.set 2
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

// ARG = SP - 5 - 2 (=nArgs)
@SP
D=M

@5
D=D-A

@2 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Class1.set
@Class1.set
0;JEQ

// (returnAddress)
( Sys$ret.0 )

//pop temp 0
//*R13=5 + 0
@0 								// offset
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

//push constant 23
//*SP=23
@23
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 15
//*SP=15
@15
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

// call Class2.set 2
// push return address
@Sys$ret.1
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

// ARG = SP - 5 - 2 (=nArgs)
@SP
D=M

@5
D=D-A

@2 // (= nArgs)
D=D-A

@R2
M=D

// LCL = SP
@SP
D=M
@R1
M=D

// goto Class2.set
@Class2.set
0;JEQ

// (returnAddress)
( Sys$ret.1 )

//pop temp 0
//*R13=5 + 0
@0 								// offset
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

// call Class1.get 0
// push return address
@Sys$ret.2
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

// goto Class1.get
@Class1.get
0;JEQ

// (returnAddress)
( Sys$ret.2 )

// call Class2.get 0
// push return address
@Sys$ret.3
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

// goto Class2.get
@Class2.get
0;JEQ

// (returnAddress)
( Sys$ret.3 )

( Sys.WHILE )

// goto Sys.WHILE
@Sys.WHILE
0;JEQ

