// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

  // Set total & loop counter to 0
  @total
  M=0
  @n
  M=0
  @R2
  M=0
(LOOP)
  // if n == R1 goto END
  @n
  D=M
  @R1
  D=D-M
  @STOP
  D;JEQ  
  // Add R0 to total
  @total
  D=M
  @R0
  D=D+M
  @total
  M=D
  // n++
  @n
  M=M+1
  // goto LOOP
  @LOOP
  0;JMP
(STOP)
  // R2 = total
  @total
  D=M
  @R2
  M=D
(END)
  @END
  0;JMP
