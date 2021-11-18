package assembly

import (
	"fmt"
	"strings"
	"vm-translator/pkg/parser"
)

var compCount = 0

// Could probably use a string build type pattern here if string allocs are an issue
func Assemble(s []parser.Statement) string {
	var asm string
	for i := range s {
		asm += fmt.Sprintf("%s\n\n", getAssembly(s[i]))
	}
	asm += endProgram()

	return asm
}

func getAssembly(s parser.Statement) string {
	o := make([]string, 0)
	o = append(o, fmt.Sprintf("//%s", s.RawStatement))
	switch {
	case s.CommandType == parser.PUSH:
		o = append(o, buildPush(s.Arg1, s.Arg2))
	case s.CommandType == parser.ADD:
		o = append(o, buildBinaryOperator("M=M+D"))
	case s.CommandType == parser.SUB:
		o = append(o, buildBinaryOperator("M=D-M"))
	case s.CommandType == parser.NEGATE:
		o = append(o, buildNegate())
	case s.CommandType == parser.AND:
		o = append(o, buildBinaryOperator("M=M&D"))
	case s.CommandType == parser.OR:
		o = append(o, buildBinaryOperator("M=M|D"))
	case s.CommandType == parser.NOT:
		return ""
	case s.CommandType == parser.EQUAL:
		o = append(o, buildComp("JEQ"))
	case s.CommandType == parser.GREATER:
		o = append(o, buildComp("JGT"))
	case s.CommandType == parser.LESS:
		o = append(o, buildComp("JLT"))
	}

	return strings.Join(o, "\n")
}

func buildPush(segment string, i int) string {
	return strings.Join([]string{
		buildSegment(segment, i),
		"@SP",
		"A=M",
		"M=D",
		spInc(),
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

func buildBinaryOperator(op string) string {
	return strings.Join([]string{
		popValue(),
		"A=A-1",
		op,
	}, "\n")
}

func buildNegate() string {
	return strings.Join([]string{
		popValue(),
	}, "\n")
}

func buildComp(comp string) string {
	compCount++

	return strings.Join([]string{
		popValue(),
		"A=A-1",
		"D=D-M",
		"M=-1", // Set output to true - A false comp will overwrite
		fmt.Sprintf("@COMP%d", compCount),
		fmt.Sprintf("D;%s", comp),
		"M=0", // No jump means false comp
		fmt.Sprintf("(@COMP%d)", compCount),
	}, "\n")
}

// Puts the top stack value into D & decrements the SP
func popValue() string {
	return strings.Join([]string{
		"@SP",
		"AM=M-1",
		"D=M",
	}, "\n")
}

func spInc() string {
	return strings.Join([]string{
		"@SP",
		"M=M+1",
	}, "\n")
}

func spDec() string {
	return strings.Join([]string{
		"@SP",
		"M=M-1",
	}, "\n")
}

func endProgram() string {
	return strings.Join([]string{
		"(END)",
		"@END",
		"0;JMP",
	}, "\n")
}
