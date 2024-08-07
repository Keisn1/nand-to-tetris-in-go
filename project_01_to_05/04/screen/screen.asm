// Get Screen address
// set that address to -1 (would be 1..1 16 bits)
// Go to Screen + 32 (would be 512 bits forward)
// do that in a Loop for n times
// Set Screen address
// set n
// Do for loop
    @SCREEN
    D=A

    @arr
    M=D

    @50
    D=A

    @n
    M=D

    @i
    M=0

    (LOOP)
    // check i < n
    @i
    D=M

    @n
    D=D-M

    @END
    D;JEQ

    // setting arr to -1
    @arr
    A=M
    M=-1

    // i++
    @i
    M=M+1

    @32 // go forward 32 addresses
    D=A

    @arr
    M=M+D

    @LOOP
    0;JMP

    (END)
    @END
    0;JMP
