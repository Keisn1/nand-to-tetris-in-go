    // initialization
    @R0
    D=M

    @n
    M=D

    @sum
    M=0

    @i
    M=1

    (LOOP)
    @R0
    D=M+1

    @i
    D=D-M

    @SAVE
    D;JLE

    @sum
    M=D+M

    @i
    M=M+1

    @LOOP
    0;JMP

    (SAVE)
    @sum
    D=M

    @R1
    M=D

    (END)
    @END
    0;JMP
