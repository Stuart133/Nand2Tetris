package compiler

const (
	FIELD = iota
	STATIC
	ARGUMENT
	LOCAL
)

type symbol struct {
	kind int
}
