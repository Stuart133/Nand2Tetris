package compiler

import (
	"fmt"
	"io"
)

type VmWriter struct {
	w io.Writer
}

func (w *VmWriter) WritePush(seg, i int) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("push %s %d\n", getSegment(seg), i)))
}

func (w *VmWriter) WriteConstPush(n string) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("push %s\n", n)))
}

func (w *VmWriter) WriteKeywordPush(kw string) {
	switch kw {
	case "null":
		_, _ = w.w.Write([]byte("push 0\n"))
	case "true":
		_, _ = w.w.Write([]byte("push 1\n"))
	case "false":
		_, _ = w.w.Write([]byte("push 0\n"))
	}
}

func (w *VmWriter) WritePop(seg, i int) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("pop %s %d\n", getSegment(seg), i)))
}

func (w *VmWriter) WriteArithmetic(symbol string) {
	switch symbol {
	case "+":
		_, _ = w.w.Write([]byte("add\n"))
	case "-":
		_, _ = w.w.Write([]byte("sub\n"))
	case "*":
		_, _ = w.w.Write([]byte("multiply\n"))
	case "/":
		_, _ = w.w.Write([]byte("divide\n"))
	case "-1":
		_, _ = w.w.Write([]byte("neg\n"))
	case "~":
		_, _ = w.w.Write([]byte("not\n"))
	case "=":
		_, _ = w.w.Write([]byte("eq\n"))
	case ">":
		_, _ = w.w.Write([]byte("gt\n"))
	case "<":
		_, _ = w.w.Write([]byte("lt\n"))
	case "&":
		_, _ = w.w.Write([]byte("and\n"))
	case "|":
		_, _ = w.w.Write([]byte("or\n"))
	case "!":
		_, _ = w.w.Write([]byte("not\n"))
	}
}

func (w *VmWriter) WriteLabel(label string) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("label %s\n", label)))
}

func (w *VmWriter) WriteGoto(label string) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("goto %s\n", label)))
}

func (w *VmWriter) WriteIf(label string) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("if-goto %s\n", label)))
}

func (w *VmWriter) WriteCall(name string, nArgs int) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("call %s %d\n", name, nArgs)))
}

func (w *VmWriter) WriteFunction(name string, nVars int) {
	_, _ = w.w.Write([]byte(fmt.Sprintf("function %s %d\n", name, nVars)))
}

func (w *VmWriter) WriteReturn() {
	_, _ = w.w.Write([]byte("return\n"))
}

func (w *VmWriter) WriteLine() {
	_, _ = w.w.Write([]byte("\n"))
}

func getSegment(seg int) string {
	switch seg {
	case STATIC:
		return "static"
	case FIELD:
		return "this"
	case ARGUMENT:
		return "argument"
	case LOCAL:
		return "local"
	case POINTER:
		return "pointer"
	default:
		panic("Unexpected segment type")
	}
}
