package parser

import (
	"testing"

	"github.com/enotinc/hmp/ast"
	"github.com/enotinc/hmp/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foo = 1233;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statement does not contains 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		exp string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.exp) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got=%q", s.TokenLiteral())
		return false
	}

	let, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got=%q", s)
		return false
	}

	if let.Name.Value != name {
		t.Errorf("want=%q, got=%q", name, let.Name.Value)
		return false
	}

	if let.Name.TokenLiteral() != name {
		t.Errorf("want=%q, got=%q", name, let.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}

	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 12;
	return 123123123;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("got len = %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returns, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stsm is not *ast.ReturnStatement, got = %T", stmt)
			continue
		}

		if returns.TokenLiteral() != "return" {
			t.Errorf("returns.TokenLiteral is not 'return', got = %q", returns.TokenLiteral())
		}
	}
}
