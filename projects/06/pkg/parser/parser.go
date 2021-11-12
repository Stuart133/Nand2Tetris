package parser

import (
	"fmt"
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
			// Remove any trailing comment & whitespace
			l := strings.Split(rawLines[i], "//")[0]
			l = strings.TrimSpace(l)

			lines = append(lines, l)
		}
	}

	return Parser{
		lines:    lines,
		position: -1,
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

	if strings.HasPrefix(l, "@") {
		return l[1:]
	} else {
		fmt.Println(l[1 : len(l)-1])
		return l[1 : len(l)-1]
	}
}

func (p *Parser) Comp() string {
	l := p.getLine()

	c := getPart(l, "=", 1)
	if c != "" {
		cj := getPart(c, ";", 0)
		if cj == "" {
			return c
		} else {
			return cj
		}
	} else {
		// If there was no dest field, we know there must be a jump, so return here
		c = getPart(l, ";", 0)
		return c
	}
}

func (p *Parser) Dest() string {
	l := p.getLine()

	return getPart(l, "=", 0)
}

func (p *Parser) Jump() string {
	l := p.getLine()

	return getPart(l, ";", 1)
}

func getPart(l, sep string, loc int) string {
	if strings.Contains(l, sep) {
		return strings.Split(l, sep)[loc]
	}

	return ""
}

func (p *Parser) Advance() {
	if p.HasMoreLines() {
		p.position++
	}
}

func (p *Parser) Reset() {
	p.position = -1
}

func (p *Parser) HasMoreLines() bool {
	return len(p.lines)-1 != p.position
}
