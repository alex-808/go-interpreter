package parser

import (
	"fmt"

	"github.com/alex-davis-808/go-interpreter/src/interpreter/ast"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/lexer"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/token"
)

const (
	// create an enum with each var below getting a number and '_' getting the zero
	// add the heirarchy to operators
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX // -X or !X
	CALL   // myFunction(X)
)

type Parser struct {
	l      *lexer.Lexer
	errors []string
	// similar to position and peekPosition but iterate over tokens instead of chars
	curToken  token.Token
	peekToken token.Token

	// pass in a token type to find it's prefix/infix function
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	// Instantiate new parser by passing in a lexer
	p := &Parser{l: l, errors: []string{}}

	// make the maps specified on the type
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// Read two tokens so both curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

// getter for Parser errors
func (p *Parser) Errors() []string {
	return p.errors
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

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		//TODO nil stmts are getting through this some how
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
			fmt.Println(program.Statements)
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	// expect the next token after 'let' to be an IDENT
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// set the statement name to be an Identifier version of the current token
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// expect the next token after the Identifier to be an 'ASSIGN'
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: We will skip the expressions until we encounter a semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// semicolon is optional, if next token is one, call nextToken to make it the curToken
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// retrieve prefix fn
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Just checks if current token is of the expected type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Just checks if next token is of expected type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// if peek is as expected, next token, else add error
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// adds error to Parser.errors
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

type (
	// accepts nothing and returns an 'Expression' AST node
	prefixParseFn func() ast.Expression
	// accepts an ast.Expression node (the left side) and returns an ast.Expression node
	infixParseFn func(ast.Expression) ast.Expression
)

// helpers to add prefix/infix functions to the map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
