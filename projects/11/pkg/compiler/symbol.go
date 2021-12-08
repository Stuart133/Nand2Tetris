package compiler

const (
	FIELD = iota
	STATIC
	ARGUMENT
	LOCAL
)

type symbol struct {
	typ  string
	kind int
}

type symbolTable struct {
	table map[string]symbol
	count map[int]int
}

func (t *symbolTable) AddSymbol(name, typ string, kind int) {
	t.table[name] = symbol{
		typ:  typ,
		kind: kind,
	}
	t.count[kind]++
}
