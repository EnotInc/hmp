package parser

import (
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("incorrect len of statements, got = %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ast.Exporession, got = %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("not an ast.Identifier, got = %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("not foobar, got = %s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("not foobar, got = %s", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("len is %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.statements[0] is %T", program.Statements[0])
	}

	lit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exproession is %T", stmt.Expression)
	}

	if lit.Value != 5 {
		t.Fatalf("value is %d", lit.Value)
	}

	if lit.TokenLiteral() != "5" {
		t.Fatalf("TokenLiteral is %s", lit.TokenLiteral())
	}
}

func TestPersingPrefixExpressions(t *testing.T) {
	prefixtests := []struct {
		input    string
		operator string
		integer  int64
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
	}

	for _, tt := range prefixtests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len is %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program statement is %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression is %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("want %s but got %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integer) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intg, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("got %t", il)
		return false
	}

	if intg.Value != value {
		t.Errorf("want %d, got %d", value, intg.Value)
		return false
	}

	if intg.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("want %d, got %s", value, intg.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 <= 3;", 5, "<=", 3},
		{"5 >= 3;", 5, ">=", 3},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len is %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program statement is %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression is %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("want %s, but got %s", exp.Operator, tt.operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}

}

func TestOperatorPrecenseParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"c + a * b", "(c + (a * b))"},
		{"1 + 2 + 3;", "((1 + 2) + 3)"},
		{"1 + 2 * 3;", "(1 + (2 * 3))"},
		{"(1 + 2) * 3;", "((1 + 2) * 3)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected = %s, got = %s", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.statement[0] is %T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("stmt.Exporession is %T", stmt.Expression)
		}

		if boolean.Value != tt.exp {
			t.Fatalf("want %t, but got %t", boolean.Value, tt.exp)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is %T", stmt.Expression)
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("len(exp.Consequence.Statements) = %d", len(exp.Consequence.Statements))
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal(program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatal(stmt.Expression)
	}

	body, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal(fn.Body.Statements[0])
	}

	t.Log(body.String())
}

func TestCallExpressionParser(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal(program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatal(stmt.Expression)
	}

	for _, e := range exp.Arguments {
		t.Log(e.String())
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello there"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	lit, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("got %T", stmt.Expression)
	}

	if lit.Value != "hello there" {
		t.Errorf("want %s, but got %s", input, lit.Value)
	}
}
