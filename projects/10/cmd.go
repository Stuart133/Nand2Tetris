package main

import (
	"compiler/pkg/scanner"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	path := os.Args[1]

	if strings.HasSuffix(path, ".jack") {
		c, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Could not open file: %v", err)
			os.Exit(1)
		}

		fmt.Println(string(c))

		s := scanner.NewScanner(string(c))
		tokens := s.ScanTokens()

		for _, t := range tokens {
			fmt.Printf("%v\n", t)
		}
	}
}
