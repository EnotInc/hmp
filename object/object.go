package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/enotinc/hmp/ast"
)

type ObjectType string

const (
	FUNCTION_OBJ = "FUCTION"
	INTERER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	STRING_OBJ   = "STRING"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	ERROR_OBJ = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// -==[ int ]==-
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTERER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// -==[ bool ]==-
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// -==[ null ]==-
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// -==[ return value ]==-
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

// -==[ ERRORS !]==-
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return fmt.Sprintf("Error: %s", e.Message) }

// -==[ function ]==-
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Enviroment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { \n")
	out.WriteString(f.Body.String())
	out.WriteString("}")
	out.WriteString("\n")

	return out.String()
}

// -==[ string ]==-
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
