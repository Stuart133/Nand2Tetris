package main

import (
	"fmt"
	"os"
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

	printSymbol(&p)
	for i := 0; p.HasMoreLines(); i++ {
		p.Advance()
		printSymbol(&p)
	}
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
