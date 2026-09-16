package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/enotinc/hmp/ast"
)

type ObjectType string

type BuildintFunction func(args ...Object) Object

const (
	BUILDID_OBJ  = "BUILDIN"
	FUNCTION_OBJ = "FUCTION"
	INTERER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	STRING_OBJ   = "STRING"
	ARRAY_OBJ    = "ARRAY"
	HASH_OBJ     = "HASH"
	BREAK_OBJ    = "BREAK"

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
func (e *Error) Inspect() string  { return fmt.Sprintf("\tError: %s\n", e.Message) }

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

// -==[ build id funciton ]==-
type BuildIn struct {
	Fn BuildintFunction
}

func (bi *BuildIn) Type() ObjectType { return BUILDID_OBJ }
func (bi *BuildIn) Inspect() string  { return "build in function" }

// -==[ Array ]==-
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// -==[ hash maps ]==-

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, p := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", p.Key.Inspect(), p.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// -==[ break ]==-
type Break struct{}

func (n *Break) Type() ObjectType { return BREAK_OBJ }
func (n *Break) Inspect() string  { return "null" }
