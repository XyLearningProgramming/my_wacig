package my_vm

import (
	"fmt"
	"monkey/my_ast"
	"monkey/my_code"
	"monkey/my_compiler"
	"monkey/my_object"
)

const StackSize = 2048

type VM struct {
	constants    []my_object.Object
	instructions my_code.Instructions
	stack        []my_object.Object
	sp           int // sp: points to the next value, top of stack is stack[sp-1]
}

func New(byteCode *my_compiler.ByteCode) *VM {
	return &VM{
		constants:    byteCode.Constants,
		instructions: byteCode.Instructions,
		stack:        make([]my_object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	// fetch-decode-execute
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := my_code.Opcode(vm.instructions[ip])
		switch op {
		case my_code.OpConstant:
			constIndex := my_code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case my_code.OpAdd:
			rightObj := vm.pop()
			leftObj := vm.pop()
			// below are almost the same as my_evaluator/eval_infix.go
			switch leftObj := leftObj.(type) {
			case *my_object.Integer:
				switch rightObj := rightObj.(type) {
				case *my_object.Integer:
					vm.push(&my_object.Integer{Value: leftObj.Value + rightObj.Value})
				case *my_object.Boolean:
					vm.push(&my_object.Integer{Value: leftObj.Value + booleanToIntObject(rightObj).Value})
				case *my_object.Float:
					vm.push(&my_object.Float{Value: integerToFloatObject(leftObj).Value + rightObj.Value})
				case *my_object.Null:
					return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
				default:
					return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
				}
			case *my_object.Boolean:
				switch rightObj := rightObj.(type) {
				case *my_object.Integer:
					vm.push(&my_object.Integer{Value: booleanToIntObject(leftObj).Value + rightObj.Value})
				case *my_object.Boolean:
					vm.push(&my_object.Integer{Value: booleanToIntObject(leftObj).Value + booleanToIntObject(rightObj).Value})
				case *my_object.Float:
					vm.push(&my_object.Float{Value: booleanToFloatObject(leftObj).Value + rightObj.Value})
				case *my_object.Null:
					return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
				default:
					return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
				}
			case *my_object.Float:
				switch rightObj := rightObj.(type) {
				case *my_object.Integer:
					vm.push(&my_object.Float{Value: leftObj.Value + integerToFloatObject(rightObj).Value})
				case *my_object.Boolean:
					vm.push(&my_object.Integer{Value: int64(leftObj.Value) + int64(booleanToFloatObject(rightObj).Value)})
				case *my_object.Float:
					vm.push(&my_object.Float{Value: leftObj.Value + rightObj.Value})
				case *my_object.Null:
					return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
				default:
					return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
				}
			case *my_object.String:
				if rightObj, ok := rightObj.(*my_object.String); ok {
					vm.push(&my_object.String{Value: leftObj.Value + rightObj.Value})
				}
				return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
			case *my_object.Null:
				return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
			default:
				return fmt.Errorf("unknown operator: %s%s%s", leftObj.Type(), "+", rightObj.Type())
			}
		}
	}
	return nil
}

func (vm *VM) StackTop() my_object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) push(obj my_object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

func (vm *VM) pop() my_object.Object {
	if vm.sp == 0 {
		return NULL
	}
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj
}

var (
	TRUE             = &my_object.Boolean{Value: true}
	FALSE            = &my_object.Boolean{Value: false}
	TRUE_AS_ONE      = &my_object.Integer{Value: 1}
	FALSE_AS_ZERO    = &my_object.Integer{Value: 0}
	TRUE_AS_ONE_FL   = &my_object.Float{Value: 1}
	FALSE_AS_ZERO_FL = &my_object.Float{Value: 0}
	NULL             = &my_object.Null{}
)

func booleanNodeToObject(b *my_ast.Boolean) *my_object.Boolean {
	if b.Value {
		return TRUE
	}
	return FALSE
}

func booleanToIntObject(b *my_object.Boolean) *my_object.Integer {
	if b == TRUE {
		return TRUE_AS_ONE
	}
	return FALSE_AS_ZERO
}

func booleanToFloatObject(b *my_object.Boolean) *my_object.Float {
	if b == TRUE {
		return TRUE_AS_ONE_FL
	}
	return FALSE_AS_ZERO_FL
}

func integerToFloatObject(i *my_object.Integer) *my_object.Float {
	return &my_object.Float{Value: float64(i.Value)}
}

func nativeBoolToBooleanObject(nb bool) *my_object.Boolean {
	if nb {
		return TRUE
	}
	return FALSE
}
