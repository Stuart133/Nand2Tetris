package main

import (
	"fmt"
	"os"
	"strings"

	"hack-assembler/pkg/assembler"
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

	a := assembler.NewAssembler(p)
	i := a.Assemble()
	oPath := fmt.Sprintf("%s.hack", strings.Split(getFilename(path), ".")[0])
	err = saveFile(oPath, i)
	if err != nil {
		fmt.Printf("There was an error writing the the output file: %v\n", err)
	}
}

func saveFile(path string, data []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := range data {
		_, err := f.Write([]byte(data[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func getFilename(path string) string {
	sp := strings.Split(path, "\\")

	return sp[len(sp)-1]
}
