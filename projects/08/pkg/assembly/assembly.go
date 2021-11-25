package assembly

import (
	"fmt"
	"strings"
	"vm-translator/pkg/parser"
)

var compCount = 0
var fName = ""
var currentFn = "" // A hack - But it does work
var callCount = 0  // Also a hack

// Could probably use a string build type pattern here if string allocs are an issue
func Assemble(s []parser.Statement, fname string) string {
	fName = fname
	var asm string
	for i := range s {
		asm += fmt.Sprintf("%s\n\n", getAssembly(s[i]))
	}

	return asm
}

func AssembleInit() string {
	currentFn = "sys.Init"

	return strings.Join([]string{
		"@256",
		"D=A",
		"@SP",
		"M=D",
		// Call Sys.init
		buildCall("Sys.init", 0),

		// Guard loop (in case we ever return from Sys.init)
		"(BADTIMES)",
		"@BADTIMES",
		"0;JMP",
		"",
		"",
	}, "\n")
}

func getAssembly(s parser.Statement) string {
	o := make([]string, 0)
	o = append(o, fmt.Sprintf("//%s", s.RawStatement))
	switch {
	case s.CommandType == parser.PUSH:
		o = append(o, buildPush(s.Arg1, s.Arg2))
	case s.CommandType == parser.POP:
		o = append(o, buildPop(s.Arg1, s.Arg2))
	case s.CommandType == parser.LABEL:
		o = append(o, fmt.Sprintf("(%s$%s)", currentFn, s.Arg1))
	case s.CommandType == parser.GOTO:
		o = append(o, buildGoto(s.Arg1))
	case s.CommandType == parser.IF_GOTO:
		o = append(o, buildIfGoto(s.Arg1))
	case s.CommandType == parser.FUNCTION:
		o = append(o, buildFunction(s.Arg1, s.Arg2))
	case s.CommandType == parser.CALL:
		o = append(o, buildCall(s.Arg1, s.Arg2))
	case s.CommandType == parser.RETURN:
		o = append(o, buildReturn())
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
		pushValue(),
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
	switch {
	case segment == "constant":
		return buildDirectAccess(fmt.Sprintf("@%d", i))
	case segment == "temp":
		return buildDirectAccess(fmt.Sprintf("@%d", 5+i))
	case segment == "pointer":
		if i == 0 {
			return buildDirectAccess("@THIS")
		} else {
			return buildDirectAccess("@THAT")
		}
	case segment == "static":
		return buildDirectAccess(fmt.Sprintf("@%s.%d", fName, i))
	case segment == "local":
		return buildPointerAccess("@LCL", i)
	case segment == "argument":
		return buildPointerAccess("@ARG", i)
	case segment == "this":
		return buildPointerAccess("@THIS", i)
	case segment == "that":
		return buildPointerAccess("@THAT", i)
	default:
		return ""
	}
}

func buildDirectAccess(a string) string {
	return strings.Join([]string{
		a,
		"D=A",
	}, "\n")
}

func buildPointerAccess(l string, i int) string {
	return strings.Join([]string{
		fmt.Sprintf("@%d", i),
		"D=A",
		l,
		"AD=M+D",
	}, "\n")
}

func buildIfGoto(l string) string {
	return strings.Join([]string{
		popValue(),
		fmt.Sprintf("@%s$%s", currentFn, l),
		"D;JNE",
	}, "\n")
}

func buildGoto(l string) string {
	return strings.Join([]string{
		fmt.Sprintf("@%s$%s", currentFn, l),
		"0;JMP",
	}, "\n")
}

func buildFunction(n string, nArgs int) string {
	currentFn = n
	callCount = 0

	asm := []string{
		fmt.Sprintf("(%s)", n),
	}
	for i := 0; i < nArgs; i++ {
		asm = append(asm, buildSegment("local", i))
		asm = append(asm, "M=0")
		asm = append(asm, spInc())
	}

	return strings.Join(asm, "\n")
}

func buildCall(fn string, n int) string {
	asm := strings.Join([]string{
		// Push return addr
		fmt.Sprintf("@%s$ret.%d", currentFn, callCount),
		"D=A",
		pushValue(),
		// Push caller memory values
		"@LCL",
		"D=M",
		pushValue(),
		"@ARG",
		"D=M",
		pushValue(),
		"@THIS",
		"D=M",
		pushValue(),
		"@THAT",
		"D=M",
		pushValue(),
		// Reposition ARG
		fmt.Sprintf("@%d", n+5),
		"D=A",
		"@SP",
		"D=M-D",
		"@ARG",
		"M=D",
		// Reposition LCL
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
		// Jump to callee
		fmt.Sprintf("@%s", fn),
		"0;JMP",
		// Return addr
		fmt.Sprintf("(%s$ret.%d)", currentFn, callCount),
	}, "\n")

	callCount++
	return asm
}

func buildReturn() string {
	return strings.Join([]string{
		// Get the return address
		"@5",
		"D=A",
		"@LCL",
		"A=M",
		"A=A-D",
		"D=M",
		saveTmp(0),
		// Put the return value into arg 0
		popValue(),
		"@ARG",
		"A=M",
		"M=D",
		// Set the stack pointer to arg 1
		"D=A+1",
		"@SP",
		"M=D",
		// Restore previous call frame
		"@LCL",
		"A=M",
		"A=A-1",
		"D=M",
		"@THAT",
		"M=D",
		buildRestorePointer("THIS", 2),
		buildRestorePointer("ARG", 3),
		buildRestorePointer("LCL", 4),
		// Jump to return address
		loadTmp(0, "A"),
		"0;JMP",
	}, "\n")
}

func buildRestorePointer(seg string, i int) string {
	return strings.Join([]string{
		fmt.Sprintf("@%d", i),
		"D=A",
		"@LCL",
		"A=M",
		"A=A-D",
		"D=M",
		fmt.Sprintf("@%s", seg),
		"M=D",
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

// Pushes D to the stack & increments the SP
func pushValue() string {
	return strings.Join([]string{
		"@SP",
		"A=M",
		"M=D",
		spInc(),
	}, "\n")
}

func spInc() string {
	return strings.Join([]string{
		"@SP",
		"M=M+1",
	}, "\n")
}

func EndProgram() string {
	return strings.Join([]string{
		"(END)",
		"@END",
		"0;JMP",
	}, "\n")
}
