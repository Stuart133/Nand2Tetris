package symbol

import "fmt"

type Table struct {
	symbols map[string]int
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
	return t.symbols[s]
}
