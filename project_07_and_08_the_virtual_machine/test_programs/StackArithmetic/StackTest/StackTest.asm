//push constant 17
//*SP=17
@17
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 17
//*SP=17
@17
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//eq
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_EQUAL_COUNTER_0
D;JEQ

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_EQUAL_COUNTER_0
0;JMP

( IS_EQUAL_COUNTER_0 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_EQUAL_COUNTER_0 )

//push constant 17
//*SP=17
@17
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 16
//*SP=16
@16
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//eq
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_EQUAL_COUNTER_1
D;JEQ

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_EQUAL_COUNTER_1
0;JMP

( IS_EQUAL_COUNTER_1 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_EQUAL_COUNTER_1 )

//push constant 16
//*SP=16
@16
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 17
//*SP=17
@17
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//eq
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_EQUAL_COUNTER_2
D;JEQ

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_EQUAL_COUNTER_2
0;JMP

( IS_EQUAL_COUNTER_2 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_EQUAL_COUNTER_2 )

//push constant 892
//*SP=892
@892
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 891
//*SP=891
@891
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

@IS_LOWER_THAN_COUNTER_3
D;JLT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_LOWER_THAN_COUNTER_3
0;JMP

( IS_LOWER_THAN_COUNTER_3 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_LOWER_THAN_COUNTER_3 )

//push constant 891
//*SP=891
@891
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 892
//*SP=892
@892
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

@IS_LOWER_THAN_COUNTER_4
D;JLT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_LOWER_THAN_COUNTER_4
0;JMP

( IS_LOWER_THAN_COUNTER_4 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_LOWER_THAN_COUNTER_4 )

//push constant 891
//*SP=891
@891
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 891
//*SP=891
@891
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

@IS_LOWER_THAN_COUNTER_5
D;JLT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_LOWER_THAN_COUNTER_5
0;JMP

( IS_LOWER_THAN_COUNTER_5 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_LOWER_THAN_COUNTER_5 )

//push constant 32767
//*SP=32767
@32767
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 32766
//*SP=32766
@32766
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//gt
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_GREATER_THAN_COUNTER_6
D;JGT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_GREATER_THAN_COUNTER_6
0;JMP

( IS_GREATER_THAN_COUNTER_6 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_GREATER_THAN_COUNTER_6 )

//push constant 32766
//*SP=32766
@32766
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 32767
//*SP=32767
@32767
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//gt
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_GREATER_THAN_COUNTER_7
D;JGT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_GREATER_THAN_COUNTER_7
0;JMP

( IS_GREATER_THAN_COUNTER_7 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_GREATER_THAN_COUNTER_7 )

//push constant 32766
//*SP=32766
@32766
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 32766
//*SP=32766
@32766
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//gt
//D=*(SP--)
@SP
AM=M-1
D=M

// D=*SP-*(SP-1)
@SP
A=M-1
D=M-D

@IS_GREATER_THAN_COUNTER_8
D;JGT

// *(SP-1)=0
@SP
A=M-1
M=0
@AFTER_IS_GREATER_THAN_COUNTER_8
0;JMP

( IS_GREATER_THAN_COUNTER_8 )
// *(SP-1)=-1
@SP
A=M-1
M=-1

( AFTER_IS_GREATER_THAN_COUNTER_8 )

//push constant 57
//*SP=57
@57
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 31
//*SP=31
@31
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//push constant 53
//*SP=53
@53
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

//push constant 112
//*SP=112
@112
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

//neg
@SP
A=M-1
M=-M

//and
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M&D

//push constant 82
//*SP=82
@82
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1

//or
//SP--; D=M
@SP
AM=M-1
D=M

@SP
A=M-1
M=M|D

//not
@SP
A=M-1
M=!M

