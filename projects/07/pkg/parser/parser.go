package parser

import (
	"strconv"
	"strings"
)

const (
	PUSH = iota
	POP
	LABEL
	GOTO
	IF
	FUNCTION
	RETURN
	CALL
	ADD
	SUB
	NEGATE
	AND
	OR
	NOT
	EQUAL
	GREATER
	LESS
)

type Statement struct {
	CommandType  int
	Arg1         string
	Arg2         int
	RawStatement string
}

func Parse(rawLines []string) []Statement {
	lines := make([]string, 0)
	for i := range rawLines {
		if !strings.HasPrefix(rawLines[i], "//") && len(strings.TrimSpace(rawLines[i])) != 0 {
			// Remove any comments & whitespace
			l := strings.Split(rawLines[i], "//")[0]
			l = strings.TrimSpace(l)

			lines = append(lines, l)
		}
	}

	stmts := make([]Statement, len(lines))
	for i := range lines {
		stmts[i] = parseLine(lines[i])
	}

	return stmts
}

func parseLine(l string) Statement {
	cmd := strings.Split(l, " ")
	switch {
	case cmd[0] == "push":
		return Statement{
			CommandType:  PUSH,
			Arg1:         cmd[1],
			Arg2:         getIntArg(cmd[2]),
			RawStatement: l,
		}
	case cmd[0] == "add":
		return Statement{
			CommandType:  ADD,
			RawStatement: l,
		}
	case cmd[0] == "sub":
		return Statement{
			CommandType:  SUB,
			RawStatement: l,
		}
	case cmd[0] == "neg":
		return Statement{
			CommandType:  NEGATE,
			RawStatement: l,
		}
	case cmd[0] == "and":
		return Statement{
			CommandType:  AND,
			RawStatement: l,
		}
	case cmd[0] == "or":
		return Statement{
			CommandType:  OR,
			RawStatement: l,
		}
	case cmd[0] == "not":
		return Statement{
			CommandType:  NOT,
			RawStatement: l,
		}
	case cmd[0] == "eq":
		return Statement{
			CommandType:  EQUAL,
			RawStatement: l,
		}
	case cmd[0] == "gt":
		return Statement{
			CommandType:  GREATER,
			RawStatement: l,
		}
	case cmd[0] == "lt":
		return Statement{
			CommandType:  LESS,
			RawStatement: l,
		}

	default:
		return Statement{}
	}
}

// Yeah yeah I know we should handle errors properly
func getIntArg(a string) int {
	i, _ := strconv.Atoi(a)

	return i
}
