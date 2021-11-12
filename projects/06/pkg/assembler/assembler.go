package assembler

import (
	"fmt"
	"hack-assembler/pkg/code"
	"hack-assembler/pkg/parser"
	"strconv"
)

func Assemble(p *parser.Parser) []string {
	i := make([]string, 0)

	i = append(i, assembleInstruction(p))
	for p.HasMoreLines() {
		p.Advance()
		i = append(i, assembleInstruction(p))
	}

	return i
}

func assembleInstruction(p *parser.Parser) string {
	if p.InstructionType() == parser.A_INSTRUCTION || p.InstructionType() == parser.L_INSTRUCTION {
		a, _ := strconv.Atoi(p.Symbol())
		return fmt.Sprintf("0%015b\n", a)
	}

	d := code.Dest(p.Dest())
	c := code.Comp(p.Comp())
	j := code.Jump(p.Jump())

	return fmt.Sprintf("111%s%s%s\n", c, d, j)
}
