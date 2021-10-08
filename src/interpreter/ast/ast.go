package ast

import "github.com/alex-davis-808/go-interpreter/src/interpreter/token"

// every node in the AST will have to implement this Node interface
type Node interface {
	// nodes must provide a TokenLiteral() method which returns the literal val of the token it's associated with
	// will only used for debugging and testing
	TokenLiteral() string
}

// Statement and Expression interfaces define types of Nodes
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program node will be the root node for every AST the parser creates
// A Program is just a series of statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	// it will have a 'LET' token type and it's literal "let"
	Token token.Token
	// holds name of binding
	Name *Identifier
	// stores expression
	Value Expression
}

// methods to satisfy the Statement and Node interfaces
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	// will have an 'IDENT' token
	Token token.Token
	Value string
}

// methods to satisfy the Expression and Node interfaces
// Indentifiers will be treated as expressions because they will produce value in some situations
// ex: let x = valueProducingIdentifier
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	// a 'RETURN' token
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
