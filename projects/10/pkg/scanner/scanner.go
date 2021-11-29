package scanner

const (
	CLASS = iota
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOOLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

type Token struct {
	Type     int
	TypeName string
	Lexeme   string
}

type Scanner struct {
	source   string
	start    int
	current  int
	reserved map[string]int
	Tokens   []Token
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:  source,
		start:   0,
		current: 0,
		Tokens:  make([]Token, 0),
		reserved: map[string]int{
			"class":       CLASS,
			"method":      METHOD,
			"function":    FUNCTION,
			"constructor": CONSTRUCTOR,
			"int":         INT,
			"boolean":     BOOLEAN,
			"char":        CHAR,
			"void":        VOID,
			"var":         VAR,
			"static":      STATIC,
			"field":       FIELD,
			"let":         LET,
			"do":          DO,
			"if":          IF,
			"else":        ELSE,
			"while":       WHILE,
			"return":      RETURN,
			"true":        TRUE,
			"false":       FALSE,
			"null":        NULL,
			"this":        THIS,
		},
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.scanToken()
		s.start = s.current
	}

	return s.Tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch {
	case c == '{':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '}':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '(':
		s.addNamedToken(SYMBOL, "symbol")
	case c == ')':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '[':
		s.addNamedToken(SYMBOL, "symbol")
	case c == ']':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '.':
		s.addNamedToken(SYMBOL, "symbol")
	case c == ',':
		s.addNamedToken(SYMBOL, "symbol")
	case c == ';':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '+':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '-':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '*':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '&':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '|':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '<':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '>':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '=':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '~':
		s.addNamedToken(SYMBOL, "symbol")
	case c == '/':
		if s.match('/') {
			s.comment()
		} else if s.match('*') {
			s.blockComment()
		} else {
			s.addNamedToken(SYMBOL, "symbol")
		}

	// Ignore whitespace
	case c == ' ':
	case c == '\n':
	case c == '\r':
	case c == '\t':

	case c == '"':
		s.string()
		s.addNamedToken(STRING_CONST, "stringConstant")

	default:
		if isDigit(c) {
			s.numeric()
			s.addNamedToken(INT_CONST, "integerConstant")
		} else if isAlpha(c) {
			s.identifier()
		}
	}
}

func (s *Scanner) addNamedToken(t int, name string) {
	lex := s.source[s.start:s.current]

	// Quick hack to remove quotes from string constants
	if t == STRING_CONST {
		lex = s.source[s.start+1 : s.current-1]
	}

	s.Tokens = append(s.Tokens, Token{
		Type:     t,
		TypeName: name,
		Lexeme:   lex,
	})
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlpaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}

func (s *Scanner) comment() {
	// Keep advancing until we hit the new line or EOF
	for !s.isAtEnd() && !s.match('\n') {
		s.advance()
	}
}

func (s *Scanner) blockComment() {
	// Keep advancing until we hit the block closer of EOF
	// Currnetly does not support nested block comments
	for !s.isAtEnd() && !(s.peek() == '*' && s.peekNext() == '/') {
		s.advance()
	}

	// Consume the block closer
	s.advance()
	s.advance()
}

func (s *Scanner) numeric() {
	// Keep advancing until the next character is not a numeric
	for isDigit(s.peek()) {
		s.advance()
	}
}

func (s *Scanner) string() {
	// Keep advancing to the end quote or EOF
	for !s.isAtEnd() && !s.match('"') {
		s.advance()
	}
}

func (s *Scanner) identifier() {
	for isAlpaNumeric(s.peek()) {
		s.advance()
	}

	// Check if we've got a reserved keyword
	i := s.source[s.start:s.current]
	t, v := s.reserved[i]
	if v {
		s.addNamedToken(t, "keyword")
	} else {
		s.addNamedToken(IDENTIFIER, "identifier")
	}
}

func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++

	return c
}

func (s *Scanner) peek() byte {
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	return s.source[s.current+1]
}

func (s *Scanner) match(c byte) bool {
	m := s.source[s.current] == c

	// Consume the character if we have a match
	if m {
		s.advance()
	}

	return m
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
