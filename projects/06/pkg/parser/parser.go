package parser

import (
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

	rawLines := strings.Split(string(content), "\n")

	// In leiu of a full scanner pass remove comments/whitespace here
	lines := make([]string, 0)
	for i := range rawLines {
		if !strings.HasPrefix(rawLines[i], "//") && !(len(strings.TrimSpace(rawLines[i])) == 0) {
			lines = append(lines, rawLines[i])
		}
	}

	return Parser{
		lines:    lines,
		position: 0,
	}, nil
}

func (p *Parser) GetLine() string {
	return p.lines[p.position]
}

func (p *Parser) Advance() {
	if p.HasMoreLines() {
		p.position++
	}
}

func (p *Parser) HasMoreLines() bool {
	return len(p.lines)-1 != p.position
}
