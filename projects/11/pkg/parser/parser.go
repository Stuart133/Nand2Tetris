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

	for !p.isAtEnd() {
		if p.check(scanner.CLASS) {
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

	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	p.match(scanner.LEFT_BRACE)

	for !p.isAtEnd() && p.peek(scanner.STATIC, scanner.FIELD) {
		n.Nodes = append(n.Nodes, p.classVarDec())
	}

	for !p.isAtEnd() && p.peek(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD) {
		n.Nodes = append(n.Nodes, p.subroutineDec())
	}

	p.match(scanner.RIGHT_BRACE)

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
	p.match(scanner.LEFT_PAREN)
	n.Nodes = append(n.Nodes, p.parameterList())
	p.match(scanner.RIGHT_PAREN)
	n.Nodes = append(n.Nodes, p.subroutineBody())

	return n
}

func (p *Parser) parameterList() SyntaxNode {
	n := SyntaxNode{
		TypeName: "parameterList",
		Nodes:    []SyntaxNode{},
	}

	if !p.peek(scanner.RIGHT_PAREN) {
		for !p.isAtEnd() {
			n.Nodes = append(n.Nodes, p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER))
			n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

			if !p.peek(scanner.COMMA) {
				break
			}

			p.match(scanner.COMMA)
		}
	}

	return n
}

func (p *Parser) subroutineBody() SyntaxNode {
	n := SyntaxNode{
		TypeName: "subroutineBody",
		Nodes:    []SyntaxNode{},
	}

	p.match(scanner.LEFT_BRACE)

	for !p.isAtEnd() && p.peek(scanner.VAR) {
		n.Nodes = append(n.Nodes, p.varDec())
	}

	n.Nodes = append(n.Nodes, p.statements())

	p.match(scanner.RIGHT_BRACE)

	return n
}

func (p *Parser) varDec() SyntaxNode {
	n := SyntaxNode{
		TypeName: "varDec",
		Nodes:    []SyntaxNode{},
	}

	p.match(scanner.VAR)
	n = p.varInner(n)

	return n
}

func (p *Parser) statements() SyntaxNode {
	n := SyntaxNode{
		TypeName: "statements",
		Nodes:    []SyntaxNode{},
	}

	for !p.isAtEnd() && !p.peek(scanner.RIGHT_BRACE) {
		switch {
		case p.check(scanner.LET):
			n.Nodes = append(n.Nodes, p.letStatement())
		case p.check(scanner.IF):
			n.Nodes = append(n.Nodes, p.ifStatement())
		case p.check(scanner.WHILE):
			n.Nodes = append(n.Nodes, p.whileStatement())
		case p.check(scanner.DO):
			n.Nodes = append(n.Nodes, p.doStatement())
		case p.check(scanner.RETURN):
			n.Nodes = append(n.Nodes, p.returnStatement())
		default:
			fmt.Println(p.source[p.current])
			panic(fmt.Sprintf("Unexpected symbol: %v", p.source[p.current]))
		}
	}

	return n
}

func (p *Parser) letStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "letStatement",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

	if p.check(scanner.LEFT_BRACKET) {
		n.Nodes = append(n.Nodes, p.expression())
		p.match(scanner.RIGHT_BRACKET)
	}

	p.match(scanner.EQUALS)
	n.Nodes = append(n.Nodes, p.expression())
	p.match(scanner.SEMICOLON)

	return n
}

func (p *Parser) ifStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "ifStatement",
		Nodes:    []SyntaxNode{},
	}

	p.match(scanner.LEFT_PAREN)
	n.Nodes = append(n.Nodes, p.expression())
	p.match(scanner.RIGHT_PAREN)
	p.match(scanner.LEFT_BRACE)
	n.Nodes = append(n.Nodes, p.statements())
	p.match(scanner.RIGHT_BRACE)

	if p.check(scanner.ELSE) {
		p.match(scanner.LEFT_BRACE)
		n.Nodes = append(n.Nodes, p.statements())
		p.match(scanner.RIGHT_BRACE)
	}

	return n
}

func (p *Parser) whileStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "whileStatement",
		Nodes:    []SyntaxNode{},
	}

	p.match(scanner.LEFT_PAREN)
	n.Nodes = append(n.Nodes, p.expression())
	p.match(scanner.RIGHT_PAREN)
	p.match(scanner.LEFT_BRACE)
	n.Nodes = append(n.Nodes, p.statements())
	p.match(scanner.RIGHT_BRACE)

	return n
}

func (p *Parser) doStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "doStatement",
		Nodes:    []SyntaxNode{},
	}

	n = p.subroutineCallInner(n)
	p.match(scanner.SEMICOLON)

	return n
}

func (p *Parser) returnStatement() SyntaxNode {
	n := SyntaxNode{
		TypeName: "returnStatement",
		Nodes:    []SyntaxNode{},
	}

	if !p.check(scanner.SEMICOLON) {
		n.Nodes = append(n.Nodes, p.expression())
		p.match(scanner.SEMICOLON)
	}

	return n
}

func (p *Parser) expressionList() SyntaxNode {
	n := SyntaxNode{
		TypeName: "expressionList",
		Nodes:    []SyntaxNode{},
	}

	for !p.isAtEnd() && !p.peek(scanner.RIGHT_PAREN) {
		for !p.isAtEnd() {
			n.Nodes = append(n.Nodes, p.expression())

			if !p.check(scanner.COMMA) {
				break
			}
		}
	}

	return n
}

func (p *Parser) expression() SyntaxNode {
	n := SyntaxNode{
		TypeName: "expression",
		Nodes:    []SyntaxNode{},
	}

	n.Nodes = append(n.Nodes, p.term())

	for !p.isAtEnd() && p.peek(scanner.PLUS, scanner.MINUS, scanner.STAR, scanner.SLASH, scanner.AND, scanner.OR, scanner.LESS_THAN, scanner.GREATER_THAN, scanner.EQUALS) {
		n.Nodes = append(n.Nodes, p.consume(scanner.PLUS, scanner.MINUS, scanner.STAR, scanner.SLASH, scanner.AND, scanner.OR, scanner.LESS_THAN, scanner.GREATER_THAN, scanner.EQUALS))
		n.Nodes = append(n.Nodes, p.term())
	}

	return n
}

func (p *Parser) term() SyntaxNode {
	n := SyntaxNode{
		TypeName: "term",
		Nodes:    []SyntaxNode{},
	}

	if p.check(scanner.LEFT_PAREN) {
		n.Nodes = append(n.Nodes, p.expression())
		p.match(scanner.RIGHT_PAREN)
	} else if p.peek(scanner.MINUS, scanner.NOT) {
		n.Nodes = append(n.Nodes, p.consume(scanner.NOT, scanner.MINUS))
		n.Nodes = append(n.Nodes, p.term())
	} else if p.peekAhead(scanner.LEFT_BRACKET) {
		n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
		p.match(scanner.LEFT_BRACKET)
		n.Nodes = append(n.Nodes, p.expression())
		p.match(scanner.RIGHT_BRACKET)
	} else if p.peekAhead(scanner.DOT, scanner.LEFT_PAREN) {
		n = p.subroutineCallInner(n)
	} else {
		n.Nodes = append(n.Nodes, p.consume(scanner.INT_CONST, scanner.STRING_CONST, scanner.TRUE, scanner.FALSE, scanner.NULL, scanner.THIS, scanner.IDENTIFIER))
	}

	return n
}

func (p *Parser) varInner(n SyntaxNode) SyntaxNode {
	n.Nodes = append(n.Nodes, p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER))
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))

	for !p.isAtEnd() && p.check(scanner.COMMA) {
		n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	}

	p.match(scanner.SEMICOLON)

	return n
}

func (p *Parser) subroutineCallInner(n SyntaxNode) SyntaxNode {
	n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	if p.check(scanner.DOT) {
		n.Nodes = append(n.Nodes, p.consume(scanner.IDENTIFIER))
	}
	p.match(scanner.LEFT_PAREN)
	n.Nodes = append(n.Nodes, p.expressionList())
	p.match(scanner.RIGHT_PAREN)

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

func (p *Parser) match(types ...int) {
	c := p.source[p.current]

	for _, t := range types {
		if c.Type == t {
			p.advance()
			return
		}
	}

	panic(fmt.Sprintf("Unexpected token: %v", p.source[p.current]))
}

func (p *Parser) check(t int) bool {
	if p.peek(t) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) peek(types ...int) bool {
	c := p.source[p.current]

	for _, t := range types {
		if c.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) peekAhead(types ...int) bool {
	c := p.source[p.current+1]

	for _, t := range types {
		if c.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) advance() scanner.Token {
	t := p.source[p.current]
	p.current++

	return t
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.source)
}
