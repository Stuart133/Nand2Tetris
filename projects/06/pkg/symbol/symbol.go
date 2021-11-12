package symbol

type Table struct {
	symbols map[string]int
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
