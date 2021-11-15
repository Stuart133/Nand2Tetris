package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"vm-translator/pkg/assembly"
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

	for i := range stmts {
		fmt.Println(stmts[i])
	}

	asm := assembly.Assemble(stmts)

	oPath := fmt.Sprintf("%s.hack", strings.Split(getFilename(path), ".")[0])
	err = saveFile(oPath, asm)
	if err != nil {
		fmt.Printf("There was an error writing the the output file: %v\n", err)
	}
}

func saveFile(path string, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func getFilename(path string) string {
	sp := strings.Split(path, "\\")

	return sp[len(sp)-1]
}
