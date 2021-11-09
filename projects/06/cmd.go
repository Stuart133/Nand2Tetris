package main

import (
	"fmt"
	"os"
	"strings"

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
		fmt.Println("There was an error loading the assembly file: $v", err)
		os.Exit(1)
	}

	fmt.Printf("%d: %s\n", 0, p.GetLine())
	for i := 0; p.HasMoreLines(); i++ {
		p.Advance()
		fmt.Printf("%d: %s\n", i, p.GetLine())
	}
}
