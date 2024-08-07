// Read contents of RAM[KBD]
// 0 is idle

    (LOOP)
    @KBD
    D=M

    @BLACK
    D;JNE

    @SCREEN
    D=M

    @arr

    @LOOP
    0;JMP

    (BLACK)
    @LOOP
    0;JMP

