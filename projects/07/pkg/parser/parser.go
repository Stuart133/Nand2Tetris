package parser

import "strings"

const (
	ARITHMETIC = iota
	PUSH
	POP
	LABEL
	GOTO
	IF
	FUNCTION
	RETURN
	CALL
)

const (
	ADD = "add"
)

type Statement struct {
	CommandType int
	Arg1        string
	Arg2        int
}

func Parse(rawLines []string) Statement {
	// This overallocates somewhat - Probably fine
	lines := make([]string, len(rawLines))
	for i := range rawLines {
		if len(strings.TrimSpace(rawLines[i])) != 0 {
			// Remove any comments & whitespace
			l := strings.Split(rawLines[i], "//")[0]
			l = strings.TrimSpace(l)

			lines = append(lines, l)
		}
	}

	for i := range lines {
		cmd := strings.Split(lines[i], " ")
		switch {
		case cmd[0] == "push":
			return Statement{
				CommandType: PUSH,
				Arg1:        cmd[1],
			}
		case cmd[0] == "add":
			return Statement{
				CommandType: ARITHMETIC,
				Arg1:        ADD,
			}
		}
	}

	return Statement{}
}
