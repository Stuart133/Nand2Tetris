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
	source      []scanner.Token
	current     int
	writer      VmWriter
	global      symbolTable
	subroutine  symbolTable
	branchCount int
	funcReturn  int
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
			c.class()
			c.writer.WriteLine()
		} else {
			panic("Unexpected token")
		}
	}

	return nil
}

func (c *Compiler) class() {
	id := c.consume(scanner.IDENTIFIER)
	c.match(scanner.LEFT_BRACE)

	c.global = newSymbolTable()

	for !c.isAtEnd() && c.peek(scanner.STATIC, scanner.FIELD) {
		c.classVarDec()
	}

	for !c.isAtEnd() && c.peek(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD) {
		c.subroutineDec(id.Lexeme)
	}

	c.match(scanner.RIGHT_BRACE)
}

func (c *Compiler) classVarDec() {
	kind := c.consume(scanner.STATIC, scanner.FIELD)
	if kind.Lexeme == "static" {
		c.varInner(STATIC)
	} else {
		c.varInner(FIELD)
	}
}

func (c *Compiler) subroutineDec(className string) {
	c.subroutine = newSymbolTable()
	c.branchCount = 0
	c.writer.WriteLine()

	var t int
	if c.peek(scanner.METHOD) {
		t = c.consume(scanner.METHOD).Type
		c.subroutine.addSymbol("this", className, ARGUMENT)
	} else {
		// TODO: Constructors
		t = c.consume(scanner.CONSTRUCTOR, scanner.FUNCTION).Type
	}

	// TODO: Handle return type
	rt := c.consume(scanner.VOID, scanner.IDENTIFIER, scanner.INT, scanner.CHAR, scanner.BOOLEAN)
	c.funcReturn = rt.Type
	name := c.consume(scanner.IDENTIFIER)

	c.match(scanner.LEFT_PAREN)
	c.parameterList()
	c.match(scanner.RIGHT_PAREN)

	c.subroutineBody(className, name.Lexeme, t)
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

func (p *Compiler) subroutineBody(className, subroutineName string, subType int) {
	p.match(scanner.LEFT_BRACE)

	nVar := 0
	for !p.isAtEnd() && p.peek(scanner.VAR) {
		c := p.varDec()
		nVar += c
	}

	p.writer.WriteFunction(fmt.Sprintf("%s.%s", className, subroutineName), nVar)
	if subType == scanner.METHOD {
		p.writer.WritePush(ARGUMENT, 0)
		p.writer.WritePop(POINTER, 0)
	}

	p.statements()
	p.match(scanner.RIGHT_BRACE)
}

func (p *Compiler) varDec() int {
	p.match(scanner.VAR)
	return p.varInner(LOCAL)
}

func (c *Compiler) statements() {
	for !c.isAtEnd() && !c.peek(scanner.RIGHT_BRACE) {
		switch {
		case c.check(scanner.LET):
			c.letStatement()
		case c.check(scanner.IF):
			c.ifStatement()
		case c.check(scanner.WHILE):
			c.whileStatement()
		case c.check(scanner.DO):
			c.doStatement()
		case c.check(scanner.RETURN):
			c.returnStatement()
		default:
			panic(fmt.Sprintf("Unexpected symbol: %v", c.source[c.current]))
		}
	}
}

func (c *Compiler) letStatement() {
	name := c.consume(scanner.IDENTIFIER)
	symbol, _ := c.getSymbol(name.Lexeme)

	if c.check(scanner.LEFT_BRACKET) {
		c.writer.WritePush(symbol.kind, symbol.count)
		c.expression()
		c.match(scanner.RIGHT_BRACKET)
		c.writer.WriteArithmetic("+")

		c.match(scanner.EQUALS)
		c.expression()
		c.match(scanner.SEMICOLON)
		c.writer.WritePop(TEMP, 0)
		c.writer.WritePop(POINTER, 1)
		c.writer.WritePush(TEMP, 0)
		c.writer.WritePop(THAT, 0)
	} else {
		c.match(scanner.EQUALS)
		c.expression()
		c.match(scanner.SEMICOLON)
		c.writer.WritePop(symbol.kind, symbol.count)
	}
}

func (c *Compiler) ifStatement() {
	bc := c.branchCount
	c.branchCount += 2

	c.match(scanner.LEFT_PAREN)
	c.expression()
	c.match(scanner.RIGHT_PAREN)
	c.writer.WriteArithmetic("~")
	c.writer.WriteIf(fmt.Sprintf("L%d", bc))

	c.match(scanner.LEFT_BRACE)
	c.statements()
	c.match(scanner.RIGHT_BRACE)
	c.writer.WriteGoto(fmt.Sprintf("L%d", bc+1))

	c.writer.WriteLabel(fmt.Sprintf("L%d", bc))
	if c.check(scanner.ELSE) {
		c.match(scanner.LEFT_BRACE)
		c.statements()
		c.match(scanner.RIGHT_BRACE)
	}
	c.writer.WriteLabel(fmt.Sprintf("L%d", bc+1))

	c.branchCount += 2
}

func (c *Compiler) whileStatement() {
	bc := c.branchCount
	c.branchCount += 2

	c.match(scanner.LEFT_PAREN)
	c.writer.WriteLabel(fmt.Sprintf("L%d", bc))
	c.expression()
	c.match(scanner.RIGHT_PAREN)
	c.writer.WriteArithmetic("~")
	c.writer.WriteIf(fmt.Sprintf("L%d", bc+1))

	c.match(scanner.LEFT_BRACE)
	c.statements()
	c.writer.WriteGoto(fmt.Sprintf("L%d", bc))
	c.match(scanner.RIGHT_BRACE)

	c.writer.WriteLabel(fmt.Sprintf("L%d", bc+1))
}

func (c *Compiler) doStatement() {
	c.subroutineCallInner()
	c.match(scanner.SEMICOLON)

	c.writer.WritePop(TEMP, 0)
}

func (c *Compiler) returnStatement() {
	if !c.check(scanner.SEMICOLON) {
		c.expression()
		c.match(scanner.SEMICOLON)
	}

	if c.funcReturn == scanner.VOID {
		c.writer.WriteConstPush("0")
	}

	c.writer.WriteReturn()
}

func (c *Compiler) expressionList() int {
	count := 0
	for !c.isAtEnd() && !c.peek(scanner.RIGHT_PAREN) {
		for !c.isAtEnd() {
			count++
			c.expression()

			if !c.check(scanner.COMMA) {
				break
			}
		}
	}

	return count
}

func (c *Compiler) expression() {
	c.term()

	for !c.isAtEnd() && c.peek(scanner.PLUS, scanner.MINUS, scanner.STAR, scanner.SLASH, scanner.AND, scanner.OR, scanner.LESS_THAN, scanner.GREATER_THAN, scanner.EQUALS) {
		op := c.consume(scanner.PLUS, scanner.MINUS, scanner.STAR, scanner.SLASH, scanner.AND, scanner.OR, scanner.LESS_THAN, scanner.GREATER_THAN, scanner.EQUALS)
		c.term()
		c.writer.WriteArithmetic(op.Lexeme)
	}
}

func (c *Compiler) term() {
	if c.check(scanner.LEFT_PAREN) {
		c.expression()
		c.match(scanner.RIGHT_PAREN)
	} else if c.peek(scanner.MINUS, scanner.NOT) {
		op := c.consume(scanner.NOT, scanner.MINUS)
		c.term()
		if op.Type == scanner.MINUS {
			c.writer.WriteArithmetic("-1")
		} else {
			c.writer.WriteArithmetic(op.Lexeme)
		}
	} else if c.peekAhead(scanner.LEFT_BRACKET) {
		v := c.consume(scanner.IDENTIFIER)
		symbol, _ := c.getSymbol(v.Lexeme)
		c.writer.WritePush(symbol.kind, symbol.count)
		c.match(scanner.LEFT_BRACKET)
		c.expression()
		c.match(scanner.RIGHT_BRACKET)
		c.writer.WriteArithmetic("+")
		c.writer.WritePop(POINTER, 1)
		c.writer.WritePush(THAT, 0)
	} else if c.peekAhead(scanner.DOT, scanner.LEFT_PAREN) {
		c.subroutineCallInner()
	} else if c.peek(scanner.THIS, scanner.IDENTIFIER) {
		v := c.consume(scanner.THIS, scanner.IDENTIFIER)
		symbol, _ := c.getSymbol(v.Lexeme)
		c.writer.WritePush(symbol.kind, symbol.count)
	} else if c.peek(scanner.INT_CONST) {
		n := c.consume(scanner.INT_CONST)
		c.writer.WriteConstPush(n.Lexeme)
	} else if c.peek(scanner.STRING_CONST) {
		s := c.consume(scanner.STRING_CONST)
		c.writer.WriteConstPush(fmt.Sprintf("%d", len(s.Lexeme)))
		c.writer.WriteCall("String.new", 1)
		for i := range s.Lexeme {
			c.writer.WriteConstPush(fmt.Sprintf("%d", s.Lexeme[i]))
			c.writer.WriteCall("String.appendChar", 2)
		}
	} else {
		n := c.consume(scanner.TRUE, scanner.FALSE, scanner.NULL)
		c.writer.WriteKeywordPush(n.Lexeme)
	}
}

func (c *Compiler) varInner(kind int) int {
	typ := c.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENTIFIER)
	name := c.consume(scanner.IDENTIFIER)

	var table symbolTable
	if kind == LOCAL {
		table = c.subroutine
	} else {
		table = c.global
	}

	table.addSymbol(name.Lexeme, typ.Lexeme, kind)

	count := 1
	for !c.isAtEnd() && c.check(scanner.COMMA) {
		count++
		name = c.consume(scanner.IDENTIFIER)
		table.addSymbol(name.Lexeme, typ.Lexeme, kind)
	}

	c.match(scanner.SEMICOLON)

	return count
}

func (c *Compiler) subroutineCallInner() {
	name := c.consume(scanner.IDENTIFIER).Lexeme
	if c.check(scanner.DOT) {
		subroutineName := c.consume(scanner.IDENTIFIER)
		symbol, v := c.getSymbol(name)
		if !v {
			name = fmt.Sprintf("%s.%s", name, subroutineName.Lexeme)
		} else {
			c.writer.WritePush(symbol.kind, symbol.count)
			name = fmt.Sprintf("%s.%s", symbol.typ, subroutineName.Lexeme)
		}
	} else {
		symbol, _ := c.getSymbol("this")
		c.writer.WritePush(symbol.kind, symbol.count)
	}

	c.match(scanner.LEFT_PAREN)
	count := c.expressionList()
	c.match(scanner.RIGHT_PAREN)

	c.writer.WriteCall(name, count)
}

func (c *Compiler) getSymbol(name string) (symbol, bool) {
	s, v := c.subroutine.getSymbol(name)

	if !v {
		s, v = c.global.getSymbol(name)
	}

	return s, v
}

func (p *Compiler) consume(types ...int) scanner.Token {
	c := p.source[p.current]

	for _, t := range types {
		if c.Type == t {
			p.advance()
			return c
		}
	}

	panic(fmt.Sprintf("Unexpected token: %v\n", c))
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
