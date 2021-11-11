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
	if !strings.HasSuffix(path, ".asm") {
		fmt.Println("File specified must be a *.asm file")
		os.Exit(1)
	}

	p, err := parser.Load(path)
	if err != nil {
		fmt.Printf("There was an error loading the assembly file: %v\n", err)
		os.Exit(1)
	}

	oPath := fmt.Sprintf("%s.hack", strings.Split(getFilename(path), ".")[0])
	err = assemble(&p, oPath)
	if err != nil {
		fmt.Printf("There was an error writing the the output file: %v\n", err)
	}
}

func assemble(p *parser.Parser, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(assembleInstruction(p)))
	for p.HasMoreLines() {
		p.Advance()
		_, err = f.Write([]byte(assembleInstruction(p)))
	}

	if err != nil {
		return err
	}

	return nil
}

func assembleInstruction(p *parser.Parser) string {
	if p.InstructionType() == parser.A_INSTRUCTION || p.InstructionType() == parser.L_INSTRUCTION {
		a, _ := strconv.Atoi(strings.TrimSpace(p.Symbol()))
		return fmt.Sprintf("0%015b\n", a)
	}

	d := code.Dest(p.Dest())
	c := code.Comp(p.Comp())
	j := code.Jump(p.Jump())

	return fmt.Sprintf("111%s%s%s\n", c, d, j)
}

func getFilename(path string) string {
	sp := strings.Split(path, "\\")

	return sp[len(sp)-1]
}
