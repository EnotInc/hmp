package evaluator

import (
	"fmt"

	"github.com/enotinc/hmp/object"
)

var buildins = map[string]*object.BuildIn{
	"len":   {Fn: _len},
	"first": {Fn: _first},
	"last":  {Fn: _last},
	"tail":  {Fn: _tail},
	"push":  {Fn: _push},
	"pop":   {Fn: _pop},
	"print": {Fn: _print},
}

func _len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to 'len' not supported, got %s", args[0].Type())
	}
}

func _first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'first' must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func _last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'last' must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	ln := len(arr.Elements)
	if ln > 0 {
		return arr.Elements[ln-1]
	}
	return NULL
}

func _tail(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'first' must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	ln := len(arr.Elements)
	if ln > 0 {
		n := make([]object.Object, ln-1)
		copy(n, arr.Elements[1:ln])
		return &object.Array{Elements: n}
	}
	return NULL
}

func _push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of agruments. got %d, want 2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'first' must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	ln := len(arr.Elements)

	n := make([]object.Object, ln+1)
	copy(n, arr.Elements)
	n[ln] = args[1]

	return &object.Array{Elements: n}
}

func _pop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'first' must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	ln := len(arr.Elements)
	if ln > 0 {
		n := make([]object.Object, ln-1)
		copy(n, arr.Elements[:ln-1])
		return &object.Array{Elements: n}
	}
	return NULL
}

func _print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	fmt.Print("\n")

	return NULL
}
