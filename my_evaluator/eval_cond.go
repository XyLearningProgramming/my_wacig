package my_evaluator

import (
	"monkey/my_ast"
	"monkey/my_object"
)

func evalIfExpression(ie *my_ast.IfExpression, env *my_object.Environment) my_object.Object {
	enclosed := my_object.NewEnclosedEnvironment(env)
	cond := Eval(ie.Condition, enclosed)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(ie.Consequence, enclosed)
	}
	if ie.Alternative != nil {
		return Eval(ie.Alternative, enclosed)
	}
	return NULL
}

func isTruthy(obj my_object.Object) bool {
	switch obj {
	case NULL:
		fallthrough
	case FALSE:
		return false
	default:
		return true
	}
}
