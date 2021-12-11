package compiler

const (
	FIELD = iota
	STATIC
	ARGUMENT
	LOCAL
	POINTER
	TEMP
	THAT
)

type symbol struct {
	typ   string
	kind  int
	count int
}

type symbolTable struct {
	table map[string]symbol
	count map[int]int
}

func newSymbolTable() symbolTable {
	return symbolTable{
		table: map[string]symbol{},
		count: map[int]int{},
	}
}

func (t *symbolTable) getCount(kind int) int {
	count := 0
	for _, v := range t.table {
		if v.kind == kind {
			count++
		}
	}

	return count
}

func (t *symbolTable) addSymbol(name, typ string, kind int) {
	t.table[name] = symbol{
		typ:   typ,
		kind:  kind,
		count: t.count[kind],
	}

	t.count[kind]++
}

func (t *symbolTable) getSymbol(name string) (symbol, bool) {
	s, v := t.table[name]

	return s, v
}
