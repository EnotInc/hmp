package ast

import (
	"bytes"
	"strings"

	"github.com/enotinc/hmp/token"
)

type Node interface {
	TokenLiteral() string // this function is used only in debugging
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

/*
This Program node is going to be the rood node of every AST our parser produces.
Every valid program is a series of statements.
These statements are contaided in the Program.Statements, which is just a slive of AST nodes that implemen the Statement interfaceI
*/
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// -==[ Let Statement ]==-

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// -==[ Return Statement ]==-

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// -==[ Expresstion Statement ]==-

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// -==[ Integer Literal ]==-
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// -==[ Booleans ]==-
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// -==[ Prefix Expressions ]==-
type PrefixExpression struct {
	Token    token.Token // the prefix token, line ! or -
	Operator string
	Right    Expression
}

func (px *PrefixExpression) expressionNode()      {}
func (px *PrefixExpression) TokenLiteral() string { return px.Token.Literal }
func (px *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(px.Operator)
	out.WriteString(px.Right.String())
	out.WriteString(")")

	return out.String()
}

// -==[ Infix Expressions ]==-
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ix *InfixExpression) expressionNode()      {}
func (ix *InfixExpression) TokenLiteral() string { return ix.Token.Literal }
func (ix *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ix.Left.String())
	out.WriteString(" ")
	out.WriteString(ix.Operator)
	out.WriteString(" ")
	out.WriteString(ix.Right.String())
	out.WriteString(")")

	return out.String()
}

// -==[ if Expression ]==-
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	ALternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.ALternative != nil {
		out.WriteString("else")
		out.WriteString(ie.ALternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the '{'
	Statements []Statement
}

func (be *BlockStatement) expressionNode()      {}
func (be *BlockStatement) TokenLiteral() string { return be.Token.Literal }
func (be *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range be.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// -==[ Fucntions ]==-
type FunctionLiteral struct {
	Token      token.Token // 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// -==[ Call Function ]==-
type CallExpression struct {
	Token     token.Token // (
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString("(")
	out.WriteString(ce.Function.String())
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
