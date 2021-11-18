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
	case s.CommandType == parser.POP:
		o = append(o, buildPop(s.Arg1, s.Arg2))
	case s.CommandType == parser.ADD:
		o = append(o, buildBinaryOperator("M=M+D"))
	case s.CommandType == parser.SUB:
		o = append(o, buildBinaryOperator("M=M-D"))
	case s.CommandType == parser.NEGATE:
		o = append(o, buildUnaryOperator("M=-D"))
	case s.CommandType == parser.AND:
		o = append(o, buildBinaryOperator("M=M&D"))
	case s.CommandType == parser.OR:
		o = append(o, buildBinaryOperator("M=M|D"))
	case s.CommandType == parser.NOT:
		o = append(o, buildUnaryOperator("M=!D"))
	case s.CommandType == parser.EQUAL:
		o = append(o, buildComp("JEQ"))
	case s.CommandType == parser.GREATER:
		o = append(o, buildComp("JGT"))
	case s.CommandType == parser.LESS:
		o = append(o, buildComp("JLT"))
	}

	return strings.Join(o, "\n")
}

func buildPop(segment string, i int) string {
	return strings.Join([]string{
		popValue(),
		saveTmp(0),
		buildSegment(segment, i),
		saveTmp(1),
		loadTmp(0, "D"),
		loadTmp(1, "A"),
		"M=D",
	}, "\n")
}

func saveTmp(i int) string {
	return strings.Join([]string{
		fmt.Sprintf("@R%d", 13+i),
		"M=D",
	}, "\n")
}

func loadTmp(i int, r string) string {
	return strings.Join([]string{
		fmt.Sprintf("@R%d", 13+i),
		fmt.Sprintf("%s=M", r),
	}, "\n")
}

func buildPush(segment string, i int) string {
	return strings.Join([]string{
		buildLoadSegment(segment, i),
		"@SP",
		"A=M",
		"M=D",
		spInc(),
	}, "\n")
}

func buildLoadSegment(segment string, i int) string {
	var seg string
	seg = buildSegment(segment, i)

	if segment != "constant" {
		seg += "\nD=M"
	}

	return seg
}

func buildSegment(segment string, i int) string {
	var seg []string
	switch {
	case segment == "constant":
		seg = []string{
			fmt.Sprintf("@%d", i),
			"D=A",
		}
	case segment == "temp":
		seg = []string{
			fmt.Sprintf("@%d", 5+i),
			"D=A",
		}
	case segment == "local":
		seg = []string{
			buildAccess("@LCL", i),
		}
	case segment == "argument":
		seg = []string{
			buildAccess("@ARG", i),
		}
	case segment == "this":
		seg = []string{
			buildAccess("@THIS", i),
		}
	case segment == "that":
		seg = []string{
			buildAccess("@THAT", i),
		}
	}
	return strings.Join(seg, "\n")
}

func buildAccess(l string, i int) string {
	return strings.Join([]string{
		fmt.Sprintf("@%d", i),
		"D=A",
		l,
		"AD=M+D",
	}, "\n")
}

func buildBinaryOperator(op string) string {
	return strings.Join([]string{
		popValue(),
		"A=A-1",
		op,
	}, "\n")
}

func buildUnaryOperator(op string) string {
	return strings.Join([]string{
		popValue(),
		op,
		spInc(),
	}, "\n")
}

func buildComp(comp string) string {
	compCount++

	return strings.Join([]string{
		popValue(),
		"A=A-1",
		"D=M-D",
		"M=-1", // Set output to true - A false comp will overwrite
		fmt.Sprintf("@COMP%d", compCount),
		fmt.Sprintf("D;%s", comp),
		"@SP",
		"A=M-1",
		"M=0", // No jump means false comp
		fmt.Sprintf("(COMP%d)", compCount),
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

func endProgram() string {
	return strings.Join([]string{
		"(END)",
		"@END",
		"0;JMP",
	}, "\n")
}
