package compiler

const (
	FIELD = iota
	STATIC
	ARGUMENT
	LOCAL
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

func (t *symbolTable) addSymbol(name, typ string, kind int) {
	t.table[name] = symbol{
		typ:  typ,
		kind: kind,
	}
	t.count[kind]++
}

func (t *symbolTable) getSymbol(name string) (symbol, bool) {
	s, v := t.table[name]
	if v {
		s.count = t.count[s.kind] - 1
	}

	return s, v
}
