package main

import (
	"fmt"
	"hack-assembler/pkg/parser"
	"os"
)

func main() {
	path := os.Args[1]
	if path[len(path)-3:] != "asm" {
		fmt.Println("File specified must be a *.asm file")
		os.Exit(1)
	}

	_, err := parser.Load(path)
	if err != nil {
		fmt.Println("There was an error loading the assembly file: $v", err)
		os.Exit(1)
	}
}
