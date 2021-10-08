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
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// takes in curToken.Type and chooses the algorithm needed to parse the statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// set the statement name to be an Identifier version of the current token
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: We will skip the expressions until we encounter a semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Just checks if current token is of the expected type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Just checks if next token is of expected type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}

}
