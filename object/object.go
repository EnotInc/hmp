package object

import (
	"fmt"
)

type ObjectType string

const (
	INTERER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	ERROR_OBJ = "ERROR"
)

type Object interface {
	Type() ObjectType
	Insect() string
}

// -==[ int ]==-
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTERER_OBJ }
func (i *Integer) Insect() string   { return fmt.Sprintf("%d", i.Value) }

// -==[ bool ]==-
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Insect() string   { return fmt.Sprintf("%t", b.Value) }

// -==[ null ]==-
type Null struct {
	Value bool
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Insect() string   { return "null" }

// -==[ return value ]==-
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Insect() string   { return r.Value.Insect() }

// -==[ ERRORS !]==-
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Insect() string   { return fmt.Sprintf("Error: %s", e.Message) }
