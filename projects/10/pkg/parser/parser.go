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
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_BRACE))

	for !p.isAtEnd() && (p.peek() == scanner.STATIC || p.peek() == scanner.FIELD) {
		n.Nodes = append(n.Nodes, p.classVarDec())
	}

	for !p.isAtEnd() &&
		(p.peek() == scanner.CONSTRUCTOR || p.peek() == scanner.FUNCTION || p.peek() == scanner.METHOD) {
		n.Nodes = append(n.Nodes, p.subroutineDec())
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_BRACE))

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
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_PAREN))
	n.Nodes = append(n.Nodes, p.parameterList())
	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_PAREN))
	n.Nodes = append(n.Nodes, p.subroutineBody())

	return n
}

func (p *Parser) parameterList() SyntaxNode {
	n := SyntaxNode{
		TypeName: "parameterList",
		Nodes:    []SyntaxNode{},
	}

	if p.peek() != scanner.RIGHT_PAREN {
		for !p.isAtEnd() {
			n.Nodes = append(n.Nodes, p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER))
			n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

			if !p.match(scanner.COMMA) {
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

	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_BRACE))

	for !p.isAtEnd() && p.peek() == scanner.VAR {
		n.Nodes = append(n.Nodes, p.varDec())
	}

	p.statements()

	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_BRACE))

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

	for !p.isAtEnd() && p.peek() != scanner.RIGHT_PAREN {
		switch p.peek() {
		case scanner.LET:
			n.Nodes = append(n.Nodes, p.letStatement())
		case scanner.IF:
			n.Nodes = append(n.Nodes, p.ifStatement())
		case scanner.WHILE:
			n.Nodes = append(n.Nodes, p.whileStatement())
		case scanner.DO:
			n.Nodes = append(n.Nodes, p.doStatement())
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

	// TODO: Handle arrays

	n.Nodes = append(n.Nodes, p.consume(scanner.EQUALS))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER)) // TODO: Handle expressions
	n.Nodes = append(n.Nodes, p.consume(scanner.SEMICOLON))

	return n
}

func (p *Parser) ifStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "ifStatement",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.IF))
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_PAREN))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER)) // TODO: Handle expressions
	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_PAREN))
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_BRACE))
	n.Nodes = append(n.Nodes, p.statements())
	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_BRACE))

	if p.peek() == scanner.ELSE {
		n.Nodes = append(n.Nodes, p.consume(scanner.ELSE))
		n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_BRACE))
		n.Nodes = append(n.Nodes, p.statements())
		n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_BRACE))
	}

	return n
}

func (p *Parser) whileStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "whileStatement",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.WHILE))
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_PAREN))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER)) // TODO: Handle expressions
	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_PAREN))
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_BRACE))
	n.Nodes = append(n.Nodes, p.statements())
	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_BRACE))

	return n
}

func (p *Parser) doStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "doStatement",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.DO))

	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	if p.peek() == scanner.DOT {
		n.Nodes = append(n.Nodes, p.consume(scanner.DOT))
		n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	}
	n.Nodes = append(n.Nodes, p.consume(scanner.LEFT_PAREN))
	// Expr list
	n.Nodes = append(n.Nodes, p.consume(scanner.RIGHT_PAREN))

	n.Nodes = append(n.Nodes, p.consume(scanner.SEMICOLON))

	return n
}

func (p *Parser) varInner(n SyntaxNode) SyntaxNode {
	n.Nodes = append(n.Nodes, p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

	for !p.isAtEnd() && p.match(scanner.COMMA) {
		n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.SEMICOLON))

	return n
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

func (p *Parser) match(t int) bool {
	if p.source[p.current].Type == t {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) peek() int {
	return p.source[p.current].Type
}

func (p *Parser) advance() scanner.Token {
	t := p.source[p.current]
	p.current++

	return t
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.source)
}
