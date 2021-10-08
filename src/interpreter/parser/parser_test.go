package parser

import (
	"github.com/alex-davis-808/go-interpreter/src/interpreter/ast"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	// should result in an AST
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	// three statements provided so the program should identify the three
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		// 									pass in test (t), statement and expectedIdentifiers
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			// stop the loop as soon as we fail a test case
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	// TokenLiteral() is a getter for the literal value of the token
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	// typecheck that letStmt is of type *ast.LetStatement
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	// check string value of LetStatement Identifier vs expectedIdentifier
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	// check string value of Identifier on Identifier
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}
