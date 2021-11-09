package parser

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Parser struct {
	lines    []string
	position int
}

func Load(path string) (Parser, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Parser{}, err
	}

	lines := strings.Split(string(content), "/n")

	for i := range lines {
		fmt.Printf("Line: %s\n", lines[i])
	}

	return Parser{
		lines:    lines,
		position: 0,
	}, nil
}
