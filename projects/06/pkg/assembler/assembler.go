package assembler

import (
	"fmt"
	"hack-assembler/pkg/code"
	"hack-assembler/pkg/parser"
	"hack-assembler/pkg/symbol"
	"strconv"
)

type Assembler struct {
	p parser.Parser
	t symbol.Table
}

func NewAssembler(p parser.Parser) Assembler {
	return Assembler{
		p: p,
		t: symbol.NewTable(),
	}
}

func (a *Assembler) Assemble() []string {
	i := make([]string, 0)

	i = append(i, a.assembleInstruction())
	for a.p.HasMoreLines() {
		a.p.Advance()
		i = append(i, a.assembleInstruction())
	}

	return i
}

func (a *Assembler) assembleInstruction() string {
	if a.p.InstructionType() == parser.A_INSTRUCTION || a.p.InstructionType() == parser.L_INSTRUCTION {
		a, _ := strconv.Atoi(a.p.Symbol())
		return fmt.Sprintf("0%015b\n", a)
	}

	d := code.Dest(a.p.Dest())
	c := code.Comp(a.p.Comp())
	j := code.Jump(a.p.Jump())

	return fmt.Sprintf("111%s%s%s\n", c, d, j)
}
