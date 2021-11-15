package assembly

import (
	"fmt"
	"strings"
	"vm-translator/pkg/parser"
)

// Could probably use a string build type pattern here if string allocs are an issue
func Assemble(s []parser.Statement) string {
	var asm string
	for i := range s {
		asm += fmt.Sprintf("%s\n\n", getAssembly(s[i]))
	}

	return asm
}

func getAssembly(s parser.Statement) string {
	o := make([]string, 0)
	o = append(o, fmt.Sprintf("//%s", s.RawStatement))
	switch {
	case s.CommandType == parser.PUSH:
		o = append(o, buildPush(s.Arg1, s.Arg2))
	}

	return strings.Join(o, "\n")
}

func buildPush(segment string, i int) string {
	return strings.Join([]string{
		buildSegment(segment, i),
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}, "\n")
}

func buildSegment(segment string, i int) string {
	var seg []string
	switch {
	case segment == "constant":
		seg = []string{
			fmt.Sprintf("@%d", i),
			"D=A",
		}
	}

	return strings.Join(seg, "\n")
}
