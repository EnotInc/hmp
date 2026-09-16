package lexer

import (
	"testing"

	"github.com/enotinc/hmp/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
		x + y;
	};

	let res = add(five, ten);

	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10;
	10 != 9;

	"foobar"
	"foobar"

	5 <= 3 // foo bar baz
	5 >= 3

	[1, 2];

	const a = 5;

	a && b
	a || b

	for {}
	`

	tests := []struct {
		expType token.TokenType
		expLit  string
	}{
		// first line
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// second line
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// func line
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},

		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},

		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.SEMICOLON, ";"},

		// result line
		{token.LET, "let"},
		{token.IDENT, "res"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},

		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},

		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foobar"},

		{token.INT, "5"},
		{token.LT_EQ, "<="},
		{token.INT, "3"},
		{token.INT, "5"},
		{token.GT_EQ, ">="},
		{token.INT, "3"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.CONST, "const"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.IDENT, "a"},
		{token.OR, "||"},
		{token.IDENT, "b"},

		{token.FOR, "for"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expType {
			t.Fatalf("tests[%d] - wrong token type. exp %q but got %q", i, tt.expType, tok.Type)
		}

		if tok.Literal != tt.expLit {
			t.Fatalf("tests[%d] - wrong token literal. exp %q but got %q", i, tt.expLit, tok.Literal)
		}
	}
}
