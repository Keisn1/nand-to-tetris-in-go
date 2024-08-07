    // Program: Add2.asm
    // Computes RAM[2] = RAM[1] + RAM[0]
    // Usage: put values in RAM[0] & RAM[1]
    @R0
    D=M
    @R2
    M=D

    @R1
    D=M

    @R2
    M=M+D

    (END)
    @END
    0;JMP
