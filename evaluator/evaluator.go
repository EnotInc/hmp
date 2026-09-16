package evaluator

import (
	"fmt"

	"github.com/enotinc/hmp/ast"
	"github.com/enotinc/hmp/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObj(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Enviroment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		env.Const(node.Name.Value)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.AssignExpression:
		return evalAssignExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExprssion:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalPostfixExpression(node.Operator, left)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	}

	return nil
}

func evalForStatement(node *ast.ForStatement, env *object.Enviroment) object.Object {
	switch node.Condition.(type) {
	case *ast.InfixExpression, *ast.Boolean, *ast.Identifier:
		loopEnv := object.NewEnclosedEnviroment(env)

		for {
			cond := Eval(node.Condition, loopEnv)
			if isError(cond) {
				return cond
			}

			if !isTruthy(cond) {
				return NULL
			}

			body := evalBlockStatements(node.Body, loopEnv)
			if isError(body) {
				return body
			}
		}
	default:
		return NULL
	}

}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extEnv)
		return unwrapReturnvalue(evaluated)
	case *object.BuildIn:
		return fn.Fn(args...)
	default:
		return newError("not a fucntion: %s", fn.Type())
	}
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Enviroment {
	env := object.NewEnclosedEnviroment(fn.Env)

	for pID, param := range fn.Parameters {
		env.Set(param.Value, args[pID])
	}

	return env
}

func unwrapReturnvalue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTERER_OBJ:
		return evalArrayIdexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator now supported: %s", left.Type())
	}
}

func evalHashIndexExpression(left, index object.Object) object.Object {
	hash := left.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrayIdexExpression(left, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(array.Elements))

	if idx < 0 || idx >= max {
		return NULL
	}

	return array.Elements[idx]
}

func evalExpressions(exps []ast.Expression, env *object.Enviroment) []object.Object {
	var res []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}
	return res
}

func evalProgram(program *ast.Program, env *object.Enviroment) object.Object {
	var res object.Object
	for _, statement := range program.Statements {
		res = Eval(statement, env)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}

	}

	return res
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Enviroment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as has key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalAssignExpression(node *ast.AssignExpression, env *object.Enviroment) object.Object {
	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}

	switch left := node.Name.(type) {
	case *ast.Identifier:
		if env.IsConst(left.Value) {
			return newError("unable to reassign constants")
		}

		err := env.Assign(left.Value, val)
		if isError(err) {
			return err
		}
		return val

	case *ast.IndexExprssion:
		l := Eval(left.Left, env)
		if isError(l) {
			return l
		}

		if env.IsConst(left.Left.String()) {
			return newError("unable to reassign constants")
		}

		index := Eval(left.Index, env)
		if isError(index) {
			return index
		}

		res := evalIdexAssignment(l, index, val)
		if isError(res) {
			return res
		}
		err := env.Assign(left.Left.String(), res)
		if isError(err) {
			return err
		}
		return res

	default:
		return newError("assign operator is not available to %T", node.Name)
	}
}

func evalIdexAssignment(left object.Object, index object.Object, val object.Object) object.Object {
	switch left.Type() {
	case object.ARRAY_OBJ:
		arr := left.(*object.Array)
		idx, ok := index.(*object.Integer)
		if !ok {
			return newError("icorrect idnex: %T", index)
		}
		max := int64(len(arr.Elements))
		if idx.Value < 0 || idx.Value > max {
			return newError("index %d out of bounds", idx)
		}

		arr.Elements[idx.Value] = val
		return arr
	case object.HASH_OBJ:
		hash := left.(*object.Hash)
		key, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", index.Type())
		}

		v, ok := hash.Pairs[key.HashKey()]
		if !ok {
			hash.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: val}
			return hash
		}

		v.Value = val
		hash.Pairs[key.HashKey()] = v

		return hash
	default:
		return newError("incorrect assignment object %T", left)
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Enviroment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builin, ok := buildins[node.Value]; ok {
		return builin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Enviroment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	} else if ie.ALternative != nil {
		return Eval(ie.ALternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Enviroment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPostfixExpression(operator string, left object.Object) object.Object {
	l, ok := left.(*object.Integer)
	if !ok {
		return newError("Incorrect operator %s for %s. Could be ony used wtih INTEGERS", operator, left.Type())
	}
	switch operator {
	case "++":
		l.Value = l.Value + 1
		return l
	case "--":
		l.Value = l.Value - 1
		return l
	default:
		return newError("unknown operator: %s%s", operator, left.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTERER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTERER_OBJ && right.Type() == object.INTERER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case operator == "&&":
		if left.Type() == right.Type() && left.Type() == object.BOOLEAN_OBJ {
			left := left.(*object.Boolean).Value
			right := right.(*object.Boolean).Value
			return nativeBoolToBooleanObj(left && right)
		}
	case operator == "||":
		if left.Type() == right.Type() && left.Type() == object.BOOLEAN_OBJ {
			left := left.(*object.Boolean).Value
			right := right.(*object.Boolean).Value
			return nativeBoolToBooleanObj(left || right)
		}
	case operator == "==":
		return nativeBoolToBooleanObj(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObj(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: l + r}
	case "==":
		return nativeBoolToBooleanObj(l == r)
	case "!=":
		return nativeBoolToBooleanObj(l != r)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {

	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}

	case ">":
		return nativeBoolToBooleanObj(l > r)
	case "<":
		return nativeBoolToBooleanObj(l < r)
	case ">=":
		return nativeBoolToBooleanObj(l >= r)
	case "<=":
		return nativeBoolToBooleanObj(l <= r)
	case "==":
		return nativeBoolToBooleanObj(l == r)
	case "!=":
		return nativeBoolToBooleanObj(l != r)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
