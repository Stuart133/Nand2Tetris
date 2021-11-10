package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"hack-assembler/pkg/code"
	"hack-assembler/pkg/parser"
)

func main() {
	path := os.Args[1]
	if !strings.HasSuffix(path, "asm") {
		fmt.Println("File specified must be a *.asm file")
		os.Exit(1)
	}

	p, err := parser.Load(path)
	if err != nil {
		fmt.Printf("There was an error loading the assembly file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(assembleInstruction(&p))
	for i := 0; p.HasMoreLines(); i++ {
		p.Advance()
		fmt.Println(assembleInstruction(&p))
	}
}

func assembleInstruction(p *parser.Parser) string {
	if p.InstructionType() == parser.A_INSTRUCTION || p.InstructionType() == parser.L_INSTRUCTION {
		a, _ := strconv.Atoi(strings.TrimSpace(p.Symbol()))
		return fmt.Sprintf("0%015b", a)
	}

	d := code.Dest(p.Dest())
	c := code.Comp(p.Comp())
	j := code.Jump(p.Jump())

	return fmt.Sprintf("111%s%s%s", c, d, j)
}

func printSymbol(p *parser.Parser) {
	if p.InstructionType() == parser.A_INSTRUCTION || p.InstructionType() == parser.L_INSTRUCTION {
		fmt.Printf("Addr: %s\n", p.Symbol())
	}

	if p.InstructionType() == parser.C_INSTRUCTION {
		d := code.Dest(p.Dest())
		c := code.Comp(p.Comp())
		j := code.Jump(p.Jump())

		fmt.Printf("Dest: %s\n", d)
		fmt.Printf("Comp: %s\n", c)
		fmt.Printf("Jump: %s\n", j)
	}
}
