package parser

import (
	"github.com/alex-davis-808/go-interpreter/src/interpreter/ast"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/lexer"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	// similar to position and peekPosition but iterate over tokens instead of chars
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	// Instantiate new parser by passing in a lexer
	p := &Parser{l: l}

	// Read two tokens so both curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

// on first call will set only peekToken. On second will set curToken as well
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Will return the root of an ast
func (p *Parser) ParseProgram() *ast.Program {
	// store reference to Program struct
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
