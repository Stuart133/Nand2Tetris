package symbol

import "fmt"

type Table struct {
	symbols          map[string]int
	newSymbolAddress int
}

func NewTable() Table {
	t := Table{
		symbols: map[string]int{
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
			"SCREEN": 16384,
			"KBD":    24576,
		},
		newSymbolAddress: 16,
	}

	for i := 0; i < 16; i++ {
		t.symbols[fmt.Sprintf("R%d", i)] = i
	}

	return t
}

func (t *Table) AddEntry(s string, addr int) {
	t.symbols[s] = addr
}

func (t *Table) Contains(s string) bool {
	_, v := t.symbols[s]

	return v
}

func (t *Table) GetAddress(s string) int {
	_, v := t.symbols[s]
	if !v {
		t.symbols[s] = t.newSymbolAddress
		t.newSymbolAddress++
	}

	return t.symbols[s]
}
