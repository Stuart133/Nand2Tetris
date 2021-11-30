package parser

import (
	"compiler/pkg/scanner"
	"fmt"
)

type SyntaxNode struct {
	TypeName string
	Lexeme   string
	Nodes    []SyntaxNode
}

type Parser struct {
	source  []scanner.Token
	current int
}

func NewParser(t []scanner.Token) Parser {
	return Parser{
		source:  t,
		current: 0,
	}
}

func (p *Parser) Parse() []SyntaxNode {
	stmts := make([]SyntaxNode, 0)
	stmts = append(stmts, p.class())
	return stmts

	for !p.isAtEnd() {
		if p.match(scanner.CLASS) {
			stmts = append(stmts, p.class())
		} else {
			panic("Unexpected token")
		}
	}

	return stmts
}

func (p *Parser) class() SyntaxNode {
	n := SyntaxNode{
		TypeName: "class",
		Nodes:    []SyntaxNode{},
	}
	n.Nodes = append(n.Nodes, p.consume(scanner.CLASS))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, "{"))

	for !p.isAtEnd() && (p.peek().Type == scanner.STATIC || p.peek().Type == scanner.FIELD) {
		n.Nodes = append(n.Nodes, p.classVarDec())
	}

	for !p.isAtEnd() &&
		(p.peek().Type == scanner.CONSTRUCTOR || p.peek().Type == scanner.FUNCTION || p.peek().Type == scanner.METHOD) {
		n.Nodes = append(n.Nodes, p.subroutineDec())
	}

	return n
}

func (p *Parser) classVarDec() SyntaxNode {
	n := SyntaxNode{
		TypeName: "classVarDec",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.STATIC, scanner.FIELD))
	n = p.varInner(n)

	return n
}

func (p *Parser) subroutineDec() SyntaxNode {
	n := SyntaxNode{
		TypeName: "subroutineDec",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD))
	n.Nodes = append(n.Nodes, p.consume(scanner.VOID, scanner.IDENTIFIER))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, "("))
	n.Nodes = append(n.Nodes, p.parameterList())
	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, ")"))
	n.Nodes = append(n.Nodes, p.subroutineBody())

	return n
}

func (p *Parser) parameterList() SyntaxNode {
	n := SyntaxNode{
		TypeName: "parameterList",
		Nodes:    []SyntaxNode{},
	}

	if p.peek().Lexeme != ")" {
		for !p.isAtEnd() {
			n.Nodes = append(n.Nodes, p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER))
			n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

			if !p.matchWithLex(scanner.SYMBOL, ",") {
				break
			}
		}
	}

	return n
}

func (p *Parser) subroutineBody() SyntaxNode {
	n := SyntaxNode{
		TypeName: "subroutineBody",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, "{"))

	for !p.isAtEnd() && p.peek().Type == scanner.VAR {
		n.Nodes = append(n.Nodes, p.varDec())
	}

	p.statements()

	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, "}"))

	return n
}

func (p *Parser) varDec() SyntaxNode {
	n := SyntaxNode{
		TypeName: "varDec",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.VAR))
	n = p.varInner(n)

	return n
}

func (p *Parser) statements() SyntaxNode {
	n := SyntaxNode{
		TypeName: "statements",
		Nodes:    []SyntaxNode{},
	}

	for !p.isAtEnd() && p.peek().Lexeme != "}" {
		switch {
		case p.peek().Type == scanner.LET:
			n.Nodes = append(n.Nodes, p.letStatement())
		}
	}

	return n
}

func (p *Parser) letStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "letStatement",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.LET))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

	// Array handling goes here

	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, "="))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER)) // TODO: Handle expressions

	return n
}

func (p *Parser) varInner(n SyntaxNode) SyntaxNode {
	n.Nodes = append(n.Nodes, p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

	for !p.isAtEnd() && p.matchWithLex(scanner.SYMBOL, ",") {
		n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	}

	n.Nodes = append(n.Nodes, p.consumeWithLex(scanner.SYMBOL, ";"))

	return n
}

func (p *Parser) consumeWithLex(t int, lex string) SyntaxNode {
	c := p.source[p.current]

	if c.Type == t && c.Lexeme == lex {
		p.advance()
		return SyntaxNode{
			TypeName: c.TypeName,
			Lexeme:   lex,
			Nodes:    []SyntaxNode{},
		}
	}

	fmt.Printf("Got token %v\n", c)
	panic("Unexpected token")
}

func (p *Parser) consume(types ...int) SyntaxNode {
	c := p.source[p.current]

	for _, t := range types {
		if c.Type == t {
			p.advance()
			return SyntaxNode{
				TypeName: c.TypeName,
				Lexeme:   c.Lexeme,
				Nodes:    []SyntaxNode{},
			}
		}
	}

	fmt.Printf("Got token %v\n", c)
	panic("Unexpected token")
}

func (p *Parser) matchWithLex(t int, lex string) bool {
	c := p.source[p.current]

	if c.Type == t && c.Lexeme == lex {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) match(t int) bool {
	return p.matchWithLex(t, p.source[p.current].Lexeme)
}

func (p *Parser) peek() scanner.Token {
	return p.source[p.current]
}

func (p *Parser) advance() scanner.Token {
	t := p.source[p.current]
	p.current++

	return t
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.source)
}
