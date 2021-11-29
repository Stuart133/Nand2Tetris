package main

import (
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
		fmt.Printf("Could not read directory: %v", err)
	}

	for _, f := range fi {
		if !strings.HasSuffix(f.Name(), ".jack") {
			continue
		}

		c, err := ioutil.ReadFile(fmt.Sprintf("%s\\%s", path, f.Name()))
		if err != nil {
			fmt.Printf("Could not open file: %v", err)
			os.Exit(1)
		}

		s := scanner.NewScanner(string(c))
		tokens := s.ScanTokens()
		fmt.Println(f.Name())

		p := parser.NewParser(tokens)
		stmts := p.Parse()
		fmt.Println(stmts)

		oPath := fmt.Sprintf("%s\\%s-gen.xml", path, getFileName(f.Name()))
		err = saveFile(oPath, tokens)
		if err != nil {
			fmt.Printf("There was an error writing the the output file: %v\n", err)
			os.Exit(1)
		}
	}
}

func saveFile(path string, tokens []scanner.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = scanner.WriteXml(tokens, f)
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
