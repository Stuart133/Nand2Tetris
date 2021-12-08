package main

import (
	"compiler/pkg/compiler"
	"compiler/pkg/parser"
	"compiler/pkg/scanner"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	path := os.Args[1]
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Could not read directory: %v\n", err)
		os.Exit(1)
	}

	for _, f := range fi {
		if !strings.HasSuffix(f.Name(), ".jack") {
			continue
		}

		c, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, f.Name()))
		if err != nil {
			fmt.Printf("Could not open file: %v\n", err)
			os.Exit(1)
		}

		s := scanner.NewScanner(string(c))
		tokens := s.ScanTokens()

		fmt.Println(f.Name())

		p := parser.NewParser(tokens)
		stmts := p.Parse()

		oPath := fmt.Sprintf("%s\\%s.xml", path, getFileName(f.Name()))
		err = saveFile(oPath, stmts)
		if err != nil {
			fmt.Printf("There was an error writing the the output file: %v\n", err)
			os.Exit(1)
		}
		oPath = fmt.Sprintf("%s\\%s.vm", path, getFileName(f.Name()))
		err = compileAndSave(oPath, tokens)
		if err != nil {
			fmt.Printf("There was an error writing the the output file: %v\n", err)
			os.Exit(1)
		}
	}
}

func compileAndSave(path string, tokens []scanner.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	c := compiler.NewCompiler(tokens, f)
	err = c.Compile()
	if err != nil {
		return err
	}

	return nil
}

func saveFile(path string, stmts []parser.SyntaxNode) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = parser.WriteXml(stmts, f, 0)
	if err != nil {
		return err
	}

	return nil
}

func getFileName(path string) string {
	sp := strings.Split(path, "\\")
	sp = strings.Split(sp[len(sp)-1], ".")

	return sp[0]
}
