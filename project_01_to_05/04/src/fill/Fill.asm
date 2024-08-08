// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

// Init
    @color
    M=0

    (READ_INPUT)
    @KBD
    D=M

    @WHITE
    D;JEQ

    (BLACK)
        @color
        D=M

        // Screen already black? Go back to READ_INPUT
        @READ_INPUT
        D+1;JEQ

        @color
        M=-1

        @FILL_SCREEN
        0;JMP

    (WHITE)
        @color
        D=M

        // Screen already white? Go back to READ_INPUT
        @READ_INPUT
        D;JEQ

        @color
        M=0

        @FILL_SCREEN
        0;JMP

    (FILL_SCREEN)
        // Initialize
        // setting n to 8192
        // 32 * 256 = 8192
        @8192
        D=A

        @n
        M=D

        // init screen variable which is going to be incremented
        @SCREEN
        D=A

        @n
        M=M+D

        @screen
        M=D

        (FILL_LOOP)
            // set screen to color
            @color
            D=M

            @screen
            A=M
            M=D

            // increment screen
            @screen
            M=M+1

            // compare i to n
            @n
            D=M

            @screen
            D=D-M

            @FILL_LOOP
            D;JGT

            @READ_INPUT
            0;JMP
