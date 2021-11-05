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

  // current is the next screen word to be operated on
  @SCREEN
  D=A
  @current
  M=D
  @8192
  D=D+A
  @max
  M=D
(LOOP)
  // Is there a key press? FILL if yes, EMPTY if not
  @KBD
  D=M
  @EMPTY
  D;JEQ
  @FILL
  0;JMP
// Fill current word
(FILL)
  @current
  A=M
  M=-1
  @current
  M=M+1
  // If current == max goto reset
  D=M
  @max
  D=D-M
  @RESET
  D;JEQ
  // goto LOOP
  @LOOP
  0;JMP
// Empty the current word
(EMPTY)
  @current
  A=M
  M=0
  @current
  M=M+1
  // If current == max goto reset
  D=M
  @max
  D=D-M
  @RESET
  D;JEQ
  // goto LOOP
  @LOOP
  0;JMP
// Reset word counter to first word
(RESET) 
  @SCREEN
  D=A
  @current
  M=D
  // Back to main loop
  @LOOP
  0;JMP