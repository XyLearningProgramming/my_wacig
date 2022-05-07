package my_vm

import "monkey/my_object"

var (
	NULL  = &my_object.Null{}
	TRUE  = &my_object.Boolean{Value: true}
	FALSE = &my_object.Boolean{Value: false}
)

func booleanToInt(in bool) (out int64) {
	if in {
		return 1
	}
	return 0
}

func booleanToConstObj(in bool) (out my_object.Object) {
	if in {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj my_object.Object) bool {
	switch obj := obj.(type) {
	case *my_object.Boolean:
		return obj.Value
	case *my_object.Null:
		return false
	default:
		return true
	}
}
