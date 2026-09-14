package evaluator

import (
	"testing"

	"github.com/enotinc/hmp/lexer"
	"github.com/enotinc/hmp/object"
	"github.com/enotinc/hmp/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"5 * 2 - 2", 8},
		{"-50 + 100 - 50", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntergObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnviroment()

	return Eval(program, env)
}

func testIntergObject(t *testing.T, obj object.Object, expected int64) bool {
	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("got %T (%+v)", obj, obj)
		return false
	}

	if res.Value != expected {
		t.Errorf("got %d, want %d", res.Value, expected)
		return false
	}

	return true
}

func TestBooleanObject(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 < 2 == true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	res, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("got %T (+%v)", obj, obj)
		return false
	}

	if res.Value != expected {
		t.Errorf("want %t, got %t", res.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1<2) { 10 }", 10},
		{"if (1>2) { 10 }", nil},
		{"if (1>2) { 10 } else {20}", 20},
		{"if (1<2) { 10 } else {20}", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntergObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("got %T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2*5; 9;", 10},
		{"9; return 2*5; 9;", 10},
		{` if (1 > 10) {
			return 10;
		} else {
		 return 20
		 };
		`, 20},
		{` if (1 < 10) {
			return 10;
		} else {
		 return 20
		 };
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntergObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandleing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"if (1 == 1) { true + false }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar;", "identifier not found: foobar"},
		{`"hello" - "there"`, "unknown operator: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("got %T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("\nwant %s\ngot %s", tt.expected, errObj.Message)
		}
	}

}

func TestLetStatementEval(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let a = 5;a;", 5},
		{"let a = 5+5;a;", 10},
		{"let a = 5; let b = a;b", 5},
		{"let a = 5; let b = a; let c = a + b + 5;c", 15},
	}

	for _, tt := range tests {
		testIntergObject(t, testEval(tt.input), tt.expect)
	}
}

func TestFunctionObjet(t *testing.T) {
	input := "fn(x) {x + 2; };"

	evalutated := testEval(input)
	fn, ok := evalutated.(*object.Function)
	if !ok {
		t.Fatalf("got %t(%+v)", evalutated, evalutated)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("got %q, not x", fn.Parameters[0].String())
	}

	exp := "(x + 2)"
	if fn.Body.String() != exp {
		t.Fatalf("got %s", fn.Body.String())
	}
}

func TestFunctionEval(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let add = fn(x,y) {return x + y;}; add(2, 4)", 6},
	}

	for _, tt := range tests {
		testIntergObject(t, testEval(tt.input), tt.expect)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, there"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("got %T(%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello, there" {
		t.Errorf("got %s", str.Value)
	}
}

func TestBuilinFn(t *testing.T) {
	tests := []struct {
		input  string
		expect any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to 'len' not supported, got INTEGER"},
		{`len("1", "2")`, "wrong number of agruments. got 2, want 1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch exp := tt.expect.(type) {
		case int:
			testIntergObject(t, evaluated, int64(exp))
		case string:
			err, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("got %T(%+v)", evaluated, evaluated)
			}

			if err.Message != exp {
				t.Errorf("expected %q, got %q", exp, err.Message)
			}
		}
	}
}
