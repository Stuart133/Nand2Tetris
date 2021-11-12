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
	a.loadSymbols()
	a.p.Reset()

	i := make([]string, 0)

	for a.p.HasMoreLines() {
		a.p.Advance()
		i = append(i, a.assembleInstruction())
	}

	return i
}

func (a *Assembler) loadSymbols() {
	i := 0

	for a.p.HasMoreLines() {
		a.p.Advance()

		if a.p.InstructionType() == parser.L_INSTRUCTION {
			a.t.AddEntry(a.p.Symbol(), i)
		} else {
			i++
		}
	}

	fmt.Printf("%v", a.t)
}

func (a *Assembler) assembleInstruction() string {
	if a.p.InstructionType() == parser.L_INSTRUCTION {
		return ""
	}

	if a.p.InstructionType() == parser.A_INSTRUCTION {
		a, _ := strconv.Atoi(a.p.Symbol())
		fmt.Printf("Symbol: %d\n", a)
		return fmt.Sprintf("0%015b\n", a)
	}

	d := code.Dest(a.p.Dest())
	c := code.Comp(a.p.Comp())
	j := code.Jump(a.p.Jump())

	return fmt.Sprintf("111%s%s%s\n", c, d, j)
}
