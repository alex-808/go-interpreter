package parser

import (
	"fmt"
	"strconv"

	"github.com/alex-davis-808/go-interpreter/src/interpreter/ast"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/lexer"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/token"
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// register fns for token types
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// semicolon is optional, if next token is one, call nextToken to make it the curToken
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// the heart of the Pratt Parser
func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	// retrieve prefix fn
	// if token type is just an int, it will return parseIntegerLiteral
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// tries to find infixParseFn for next token
		// if no infix is found it is because the token is not an infix token

		// if it is found however it will be used to parse, passing in leftExp
		// which was created by the prefixParseFn

		// process is repeated until it encounters a token with higher precedence

		// as long as we don't meet a semicolon or a token of higher precedence,
		// we remain in the same (  )

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		// advances token so curToken points to infix
		// this way when the infix() checks curPrecedence() it will get the
		// operator of the infix expression
		p.nextToken()

		// this infix() call will also advance the token so if there is another
		// expression that it will pass back into parseExpression(), the curToken and
		// peekToken will be pointing at the beginning of that expression
		leftExp = infix(leftExp)

	}

	// if the current precedence is greater than that of the next token or we hit a
	// semicolon, we end the for loop and return

	// this will result in operations with greater precedence being more deeply placed
	// inside of the AST

	// on the outmost level of a parseExpression() call, the precedence will always being
	// LOWEST and therefore will always proceed until a semicolon
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// advance past the current token.BANG or token.MINUS
	p.nextToken()

	// parse the right expression by passing in it's precedence to parseExpression
	// attach return value to expression.Right
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// takes in a left expression and returns an ast.InfixExpression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// get precedence of peekToken
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// get precedence of curToken
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
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
