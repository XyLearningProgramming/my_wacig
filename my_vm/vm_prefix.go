package my_vm

import (
	"fmt"
	"monkey/my_object"
)

func (vm *VM) executeOpBang() {
	switch vm.pop() {
	case NULL:
		fallthrough
	case FALSE:
		vm.push(TRUE)
	case TRUE:
		fallthrough
	default:
		vm.push(FALSE)
	}
}

func (vm *VM) executeOpMinus() error {
	switch obj := vm.pop().(type) {
	case *my_object.Float:
		obj.Value = -obj.Value
		vm.push(obj)
	case *my_object.Integer:
		obj.Value = -obj.Value
		vm.push(obj)
	case *my_object.Boolean:
		if obj == TRUE {
			vm.push(FALSE)
		} else {
			vm.push(TRUE)
		}
	default:
		return fmt.Errorf("unsupported type for negation: %s", obj.Type())
	}
	return nil
}
