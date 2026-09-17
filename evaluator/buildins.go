package evaluator

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/enotinc/hmp/object"
)

var buildins = map[string]*object.BuildIn{
	"atoi": {Fn: _atoi},
	"len":  {Fn: _len},

	"first": {Fn: _first},
	"last":  {Fn: _last},
	"tail":  {Fn: _tail},
	"push":  {Fn: _push},
	"pop":   {Fn: _pop},

	"rand": {Fn: _rand},

	"scanln": {Fn: _scanln},
	"print":  {Fn: _print},
	"exit":   {Fn: _exit},

	"exists": {Fn: _exists},
	"create": {Fn: _create},
	"read":   {Fn: _read},
	"write":  {Fn: _write},
	"delete": {Fn: _delete},
	"rename": {Fn: _rename},

	"args": {Fn: _args},
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
		return newError("argument to 'tail' must be ARRAY, got %s", args[0].Type())
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
		return newError("argument to 'push' must be ARRAY, got %s", args[0].Type())
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
		return newError("argument to 'pop' must be ARRAY, got %s", args[0].Type())
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

func _scanln(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("scanln fn does not accept any args, got %d", len(args))
	}

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return newError("unable to scan input, %s", err)
	}

	return &object.String{Value: input}
}

func _print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	fmt.Print("\n")

	return NULL
}

func _exit(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	os.Exit(0)

	return NULL
}

func _exists(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to 'exists' must be STRING, got %s", args[0].Type())
	}

	name := args[0].(*object.String)
	_, err := os.Stat(name.Value)
	return nativeBoolToBooleanObj(err == nil)
}

func _create(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to 'create' must be STRING, got %s", args[0].Type())
	}

	name := args[0].(*object.String)
	f, err := os.Create(name.Value)
	if err != nil {
		return newError("unable to create file '%s'. Error: %s", name.Value, err)
	}
	defer f.Close()

	return NULL
}

func _read(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to 'read' must be STRING, got %s", args[0].Type())
	}

	name := args[0].(*object.String)
	data, err := os.ReadFile(name.Value)
	if err != nil {
		return newError("unable to create file '%s'. Error: %s", name.Value, err)
	}

	return &object.String{Value: string(data)}
}

func _write(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of agruments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
		return newError("argument to 'write' must be STRING, got %s, %s", args[0].Type(), args[1].Type())
	}

	name := args[0].(*object.String)
	data := args[1].(*object.String)
	err := os.WriteFile(name.Value, []byte(data.Value), 0644)
	if err != nil {
		return newError("unable to write file '%s'. Error: %s", name.Value, err)
	}

	return NULL
}

func _rename(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of agruments. got %d, want 2", len(args))
	}

	if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
		return newError("argument to 'rename' must be STRING, got %s, %s", args[0].Type(), args[1].Type())
	}

	old := args[0].(*object.String)
	new := args[1].(*object.String)
	err := os.Rename(old.Value, new.Value)
	if err != nil {
		return newError("unable to rename file '%s'. Error: %s", old.Value, err)
	}

	return NULL
}

func _delete(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to 'delete' must be STRING, got %s", args[0].Type())
	}

	name := args[0].(*object.String)
	err := os.Remove(name.Value)
	if err != nil {
		return newError("unable to remove file '%s'. Error: %s", name.Value, err)
	}

	return NULL
}

func _args(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("args fn does not accept any args, got %d", len(args))
	}
	as := &object.Array{}
	osa := os.Args
	if len(osa) < 2 {
		return newError("arguments wasn't provided")
	}
	for i, a := range osa {
		if i == 0 { // skipping the 'hmp' call
			continue
		}
		as.Elements = append(as.Elements, &object.String{Value: a})
	}
	return as
}

func _atoi(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of agruments. got %d, want 1", len(args))
	}
	integ, ok := args[0].(*object.String)
	if !ok {
		return newError("argument to 'atoi' must be STRING, got %s", args[0].Type())
	}
	n, err := strconv.Atoi(integ.Value)
	if err != nil {
		return newError("unable to parse string %s", integ.Value)
	}
	return &object.Integer{Value: int64(n)}
}

func _rand(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of agruments. got %d, want 2", len(args))
	}
	if args[0].Type() != object.INTERER_OBJ || args[1].Type() != object.INTERER_OBJ {
		return newError("both argument to 'rand' must be INTEGER, got %s", args[0].Type())
	}
	min := args[0].(*object.Integer).Value
	max := args[1].(*object.Integer).Value

	res := rand.Int64N(max-min) + min
	return &object.Integer{Value: res}
}
