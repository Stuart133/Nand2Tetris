package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"vm-translator/pkg/parser"
)

func main() {
	path := os.Args[1]
	if !strings.HasSuffix(path, ".vm") {
		fmt.Println("File specified must be a *.vm file")
		os.Exit(1)
	}

	c, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not load file: %v", err)
	}

	l := strings.Split(string(c), "\n")
	stmts := parser.Parse(l)
}
