package lexer

import (
	"github.com/alex-davis-808/go-interpreter/src/interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char being looked at
}

// Instantiate a new Lexer with the given input
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// initialize position, readPosition and ch
	l.readChar()
	return l
}

// gives us next character and advance position
// only supports ASCII to limit complexity
func (l *Lexer) readChar() {
	// above syntax is how you assign a method to a struct
	if l.readPosition >= len(l.input) {
		// ASCII for "NUL"
		l.ch = 0
	} else {
		// set current character
		l.ch = l.input[l.readPosition]
	}
	// increment position
	l.position = l.readPosition
	l.readPosition += 1
}

// return a TokenType and Literal for current byte then increment w/ readChar()
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		// read identifier if ch is legal symbol
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	// while loop
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	// checks that byte is within letter ASCII ranges
	// inclusion of '_' allows us to use it in identifiers
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
