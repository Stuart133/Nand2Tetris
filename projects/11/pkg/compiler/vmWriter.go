package compiler

import (
	"fmt"
	"io"
)

type VmWriter struct {
	w io.Writer
}

func (w *VmWriter) WritePush(seg, i int) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("push %s %d\n", getSegment(seg), i)))

	return err
}

func (w *VmWriter) WriteConstPush(n string) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("push %s\n", n)))

	return err
}

func (w *VmWriter) WritePop(seg, i int) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("pop %s %d\n", getSegment(seg), i)))

	return err
}

func (w *VmWriter) WriteArithmetic(symbol string) error {
	var err error
	switch symbol {
	case "+":
		_, err = w.w.Write([]byte("add\n"))
	case "-":
		_, err = w.w.Write([]byte("sub\n"))
	case "-1":
		_, err = w.w.Write([]byte("neg\n"))
	case ">":
		_, err = w.w.Write([]byte("gt\n"))
	case "<":
		_, err = w.w.Write([]byte("lt\n"))
	case "&":
		_, err = w.w.Write([]byte("and\n"))
	case "|":
		_, err = w.w.Write([]byte("or\n"))
	case "!":
		_, err = w.w.Write([]byte("not\n"))
	}

	return err
}

func (w *VmWriter) WriteLabel(label string) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("label %s\n", label)))

	return err
}

func (w *VmWriter) WriteGoto(label string) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("goto %s\n", label)))

	return err
}

func (w *VmWriter) WriteIf(label string) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("if-goto %s\n", label)))

	return err
}

func (w *VmWriter) WriteCall(name string, nArgs int) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("call %s %d\n", name, nArgs)))

	return err
}

func (w *VmWriter) WriteFunction(name string, nVars int) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("function %s %d\n", name, nVars)))

	return err
}

func (w *VmWriter) WriteReturn() error {
	_, err := w.w.Write([]byte("return\n"))

	return err
}

func (w *VmWriter) WriteLine() error {
	_, err := w.w.Write([]byte("\n"))

	return err
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
	default:
		panic("Unexpected segment type")
	}
}
