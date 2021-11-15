package parser

import (
	"fmt"
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
		fmt.Println(lines[i])
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
	default:
		return Statement{}
	}
}

// Yeah yeah I know we should handle errors properly
func getIntArg(a string) int {
	i, _ := strconv.Atoi(a)

	return i
}
