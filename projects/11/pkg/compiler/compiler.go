package compiler

import (
	"compiler/pkg/scanner"
	"fmt"
	"io"
)

type SyntaxNode struct {
	TypeName string
	Lexeme   string
	Nodes    []SyntaxNode
}

type Compiler struct {
	source     []scanner.Token
	current    int
	writer     VmWriter
	global     symbolTable
	subroutine symbolTable
}

func NewCompiler(t []scanner.Token, w io.Writer) Compiler {
	return Compiler{
		source:  t,
		current: 0,
		writer:  VmWriter{w: w},
	}
}

func (c *Compiler) Compile() error {
	for !c.isAtEnd() {
		if c.check(scanner.CLASS) {
			err := c.class()
			if err != nil {
				return err
			}
			err = c.writer.WriteLine()
			if err != nil {
				return err
			}
		} else {
			panic("Unexpected token")
		}
	}

	return nil
}

func (c *Compiler) class() error {
	id := c.consume(scanner.IDENTIFIER)
	c.match(scanner.LEFT_BRACE)

	c.global = newSymbolTable()

	for !c.isAtEnd() && c.peek(scanner.STATIC, scanner.FIELD) {
		c.classVarDec()
	}

	for !c.isAtEnd() && c.peek(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD) {
		err := c.subroutineDec(id.Lexeme)
		if err != nil {
			return err
		}
	}

	c.match(scanner.RIGHT_BRACE)

	return nil
}

func (c *Compiler) classVarDec() {
	kind := c.consume(scanner.STATIC, scanner.FIELD)
	if kind.Lexeme == "static" {
		c.varInner(STATIC)
	} else {
		c.varInner(FIELD)
	}
}

func (p *Compiler) subroutineDec(className string) error {
	p.subroutine = newSymbolTable()

	p.consume(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD)
	// TODO: Handle CTOR & METHOD
	p.consume(scanner.VOID, scanner.IDENTIFIER)
	// TODO: Handle return type
	name := p.consume(scanner.IDENTIFIER)

	p.match(scanner.LEFT_PAREN)
	p.parameterList()
	p.match(scanner.RIGHT_PAREN)

	err := p.subroutineBody(className, name.Lexeme)
	if err != nil {
		return err
	}

	return nil
}

func (p *Compiler) parameterList() {
	if !p.peek(scanner.RIGHT_PAREN) {
		for !p.isAtEnd() {
			typ := p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER).Lexeme
			name := p.consume(scanner.IDENTIFIER).Lexeme
			p.subroutine.addSymbol(name, typ, ARGUMENT)

			if !p.peek(scanner.COMMA) {
				break
			}

			p.match(scanner.COMMA)
		}
	}
}

func (p *Compiler) subroutineBody(className, subroutineName string) error {
	p.match(scanner.LEFT_BRACE)

	nVar := 0
	for !p.isAtEnd() && p.peek(scanner.VAR) {
		p.varDec()
		nVar++
	}

	err := p.writer.WriteFunction(fmt.Sprintf("%s.%s", className, subroutineName), nVar)
	if err != nil {
		return err
	}

	p.statements()
	p.match(scanner.RIGHT_BRACE)

	return nil
}

func (p *Compiler) varDec() {
	p.match(scanner.VAR)
	p.varInner(LOCAL)
}

func (c *Compiler) statements() error {
	var err error
	for !c.isAtEnd() && !c.peek(scanner.RIGHT_BRACE) {
		switch {
		case c.check(scanner.LET):
			err = c.letStatement()
		case c.check(scanner.IF):
			err = c.ifStatement()
		case c.check(scanner.WHILE):
			err = c.whileStatement()
		case c.check(scanner.DO):
			err = c.doStatement()
		case c.check(scanner.RETURN):
			c.returnStatement()
		default:
			panic(fmt.Sprintf("Unexpected symbol: %v", c.source[c.current]))
		}
	}

	return err
}

func (c *Compiler) letStatement() error {
	name := c.consume(scanner.IDENTIFIER)
	symbol, _ := c.getSymbol(name.Lexeme)

	// if p.check(scanner.LEFT_BRACKET) {
	// 	n.Nodes = append(n.Nodes, p.expression())
	// 	p.match(scanner.RIGHT_BRACKET)
	// }

	c.match(scanner.EQUALS)
	c.expression()
	c.match(scanner.SEMICOLON)
	err := c.writer.WritePop(symbol.kind, symbol.count)

	return err
}

// TOOD: FINISH IF
func (c *Compiler) ifStatement() error {
	c.match(scanner.LEFT_PAREN)
	c.expression()
	c.match(scanner.RIGHT_PAREN)

	c.match(scanner.LEFT_BRACE)
	err := c.statements()
	if err != nil {
		return err
	}
	c.match(scanner.RIGHT_BRACE)

	// if p.check(scanner.ELSE) {
	// 	p.match(scanner.LEFT_BRACE)
	// 	n.Nodes = append(n.Nodes, p.statements())
	// 	p.match(scanner.RIGHT_BRACE)
	// }

	return nil
}

func (c *Compiler) whileStatement() error {
	c.match(scanner.LEFT_PAREN)
	c.expression()
	c.match(scanner.RIGHT_PAREN)
	c.match(scanner.LEFT_BRACE)
	err := c.statements()
	if err != nil {
		return err
	}

	c.match(scanner.RIGHT_BRACE)

	return nil
}

func (c *Compiler) doStatement() error {
	err := c.subroutineCallInner()
	if err != nil {
		return err
	}

	c.match(scanner.SEMICOLON)

	return nil
}

func (c *Compiler) returnStatement() error {
	if !c.check(scanner.SEMICOLON) {
		c.expression()
		c.match(scanner.SEMICOLON)
	}

	err := c.writer.WriteReturn()
	return err
}

func (c *Compiler) expressionList() error {
	for !c.isAtEnd() && !c.peek(scanner.RIGHT_PAREN) {
		for !c.isAtEnd() {
			err := c.expression()
			if err != nil {
				return err
			}

			if !c.check(scanner.COMMA) {
				break
			}
		}
	}

	return nil
}

func (c *Compiler) expression() error {
	err := c.term()
	if err != nil {
		return err
	}

	for !c.isAtEnd() && c.peek(scanner.PLUS, scanner.MINUS, scanner.STAR, scanner.SLASH, scanner.AND, scanner.OR, scanner.LESS_THAN, scanner.GREATER_THAN, scanner.EQUALS) {
		c.consume(scanner.PLUS, scanner.MINUS, scanner.STAR, scanner.SLASH, scanner.AND, scanner.OR, scanner.LESS_THAN, scanner.GREATER_THAN, scanner.EQUALS)
		err := c.term()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) term() error {
	var err error

	if c.check(scanner.LEFT_PAREN) {
		err = c.expression()
		c.match(scanner.RIGHT_PAREN)
	} else if c.peek(scanner.MINUS, scanner.NOT) {
		c.consume(scanner.NOT, scanner.MINUS)
		err = c.term()
	} else if c.peekAhead(scanner.LEFT_BRACKET) {
		c.consume(scanner.IDENTIFIER)
		c.match(scanner.LEFT_BRACKET)
		err = c.expression()
		c.match(scanner.RIGHT_BRACKET)
	} else if c.peekAhead(scanner.DOT, scanner.LEFT_PAREN) {
		err = c.subroutineCallInner()
	} else {
		c.consume(scanner.INT_CONST, scanner.STRING_CONST, scanner.TRUE, scanner.FALSE, scanner.NULL, scanner.THIS, scanner.IDENTIFIER)
	}

	return err
}

func (c *Compiler) varInner(kind int) {
	typ := c.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER)
	name := c.consume(scanner.IDENTIFIER)

	var table symbolTable
	if kind == LOCAL {
		table = c.subroutine
	} else {
		table = c.global
	}

	table.addSymbol(name.Lexeme, typ.Lexeme, kind)

	for !c.isAtEnd() && c.check(scanner.COMMA) {
		name = c.consume(scanner.IDENTIFIER)
		table.addSymbol(name.Lexeme, typ.Lexeme, kind)
	}

	c.match(scanner.SEMICOLON)
}

func (c *Compiler) subroutineCallInner() error {
	name := c.consume(scanner.IDENTIFIER).Lexeme
	if c.check(scanner.DOT) {
		subroutineName := c.consume(scanner.IDENTIFIER)
		symbol, v := c.getSymbol(name)
		if !v {
			name = fmt.Sprintf("%s.%s", name, subroutineName.Lexeme)
		} else {
			name = fmt.Sprintf("%s.%s", symbol.typ, subroutineName.Lexeme)
		}
	} else {
		symbol, _ := c.getSymbol("this")
		err := c.writer.WritePush(symbol.kind, symbol.count)
		if err != nil {
			return err
		}
	}

	// TODO: Wire up numbers
	c.match(scanner.LEFT_PAREN)
	c.expressionList()
	c.match(scanner.RIGHT_PAREN)

	err := c.writer.WriteCall(name, 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Compiler) getSymbol(name string) (symbol, bool) {
	s, v := c.subroutine.getSymbol(name)

	if !v {
		s, v = c.global.getSymbol(name)
	}

	return s, v
}

func (p *Compiler) consume(types ...int) SyntaxNode {
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

func (p *Compiler) match(types ...int) {
	c := p.source[p.current]

	for _, t := range types {
		if c.Type == t {
			p.advance()
			return
		}
	}

	panic(fmt.Sprintf("Unexpected token: %v", p.source[p.current]))
}

func (p *Compiler) check(t int) bool {
	if p.peek(t) {
		p.advance()
		return true
	}

	return false
}

func (p *Compiler) peek(types ...int) bool {
	c := p.source[p.current]

	for _, t := range types {
		if c.Type == t {
			return true
		}
	}

	return false
}

func (p *Compiler) peekAhead(types ...int) bool {
	c := p.source[p.current+1]

	for _, t := range types {
		if c.Type == t {
			return true
		}
	}

	return false
}

func (p *Compiler) advance() scanner.Token {
	t := p.source[p.current]
	p.current++

	return t
}

func (p *Compiler) isAtEnd() bool {
	return p.current >= len(p.source)
}
