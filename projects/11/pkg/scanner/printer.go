package scanner

import (
	"fmt"
	"io"
)

func WriteXml(tokens []Token, w io.Writer) error {
	_, err := w.Write([]byte("<tokens>\n"))
	if err != nil {
		return err
	}

	for _, t := range tokens {
		l := fmt.Sprintf("<%s> %s </%s>\n", t.TypeName, t.Lexeme, t.TypeName)
		_, err = w.Write([]byte(l))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("</tokens>\n"))
	if err != nil {
		return err
	}

	return nil
}
