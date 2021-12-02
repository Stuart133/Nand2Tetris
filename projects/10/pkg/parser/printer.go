package parser

import (
	"fmt"
	"io"
)

func WriteXml(stmts []SyntaxNode, w io.Writer, i int) error {
	for _, node := range stmts {
		if len(node.Nodes) == 0 {
			l := indent(i, fmt.Sprintf("<%s> %s </%s>\n", node.TypeName, node.Lexeme, node.TypeName))

			_, err := w.Write([]byte(l))
			if err != nil {
				return err
			}

			return nil
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

func indent(indent int, c string) string {
	t := ""
	for i := 0; i < indent; i++ {
		t += "\t"
	}

	return fmt.Sprintf("%s%s", t, c)
}
