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

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

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
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
