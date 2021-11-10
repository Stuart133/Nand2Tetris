package parser

import (
	"io/ioutil"
	"strings"
)

const (
	A_INSTRUCTION = iota
	C_INSTRUCTION
	L_INSTRUCTION
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

func (p *Parser) getLine() string {
	return p.lines[p.position]
}

func (p *Parser) InstructionType() int {
	l := p.getLine()

	if strings.HasPrefix(l, "@") {
		return A_INSTRUCTION
	}

	if strings.HasPrefix(l, "(") {
		return L_INSTRUCTION
	}

	return C_INSTRUCTION
}

func (p *Parser) Symbol() string {
	l := p.getLine()

	return l[1:]
}

func (p *Parser) Advance() {
	if p.HasMoreLines() {
		p.position++
	}
}

func (p *Parser) HasMoreLines() bool {
	return len(p.lines)-1 != p.position
}
