package parser

import (
	"fmt"
	"io"
)

func WriteXml(stmts []SyntaxNode, w io.Writer, i int) error {
	for _, node := range stmts {
		if len(node.Nodes) == 0 {
			l := format(node, i)
			_, err := w.Write([]byte(l))
			if err != nil {
				return err
			}

			continue
		}

		l := indent(i, fmt.Sprintf("<%s>\n", node.TypeName))
		_, err := w.Write([]byte(l))
		if err != nil {
			return err
		}

		err = WriteXml(node.Nodes, w, i+1)
		if err != nil {
			return err
		}

		l = indent(i, fmt.Sprintf("</%s>\n", node.TypeName))
		_, err = w.Write([]byte(l))
		if err != nil {
			return err
		}

	}

	return nil
}

func format(n SyntaxNode, i int) string {
	if n.Lexeme == "" {
		lower := indent(i, fmt.Sprintf("</%s>\n", n.TypeName))
		return indent(i, fmt.Sprintf("<%s>\n%s", n.TypeName, lower))
	} else {
		return indent(i, fmt.Sprintf("<%s> %s </%s>\n", n.TypeName, n.Lexeme, n.TypeName))
	}
}

func indent(indent int, c string) string {
	t := ""
	for i := 0; i < indent; i++ {
		t += "  "
	}

	return fmt.Sprintf("%s%s", t, c)
}
