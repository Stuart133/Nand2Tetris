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
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Could not read directory: %v", err)
	}

	asm := ""
	for _, f := range fi {
		if !strings.HasSuffix(f.Name(), ".vm") {
			continue
		}

		c, err := ioutil.ReadFile(fmt.Sprintf("%s\\%s", path, f.Name()))
		if err != nil {
			fmt.Printf("Could not load file: %v", err)
		}

		l := strings.Split(string(c), "\n")
		stmts := parser.Parse(l)

		asm += assembly.Assemble(stmts, strings.Split(f.Name(), ".")[0])
	}

	asm += assembly.EndProgram()

	oPath := fmt.Sprintf("%s\\%s.asm", path, getFileName(path))
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

func getFileName(path string) string {
	sp := strings.Split(path, "\\")

	return sp[len(sp)-1]
}
