package my_evaluator

import (
	"fmt"
	"monkey/my_object"
)

func newError(format string, a ...any) *my_object.Error {
	return &my_object.Error{Message: fmt.Sprintf(format, a...)}
}

func isContinueError(obj my_object.Object) bool {
	return obj == CONTINUE_ERROR
}

func isBreakError(obj my_object.Object) bool {
	return obj == BREAK_ERROR
}

func isError(obj my_object.Object) bool {
	if obj != nil && obj.Type() == my_object.ERROR_OBJ {
		return true
	}
	return false
}

func tryUnwrapReturnValue(obj my_object.Object) my_object.Object {
	if returnVal, ok := obj.(*my_object.ReturnValue); ok {
		return returnVal.Value
	}
	return obj
}
