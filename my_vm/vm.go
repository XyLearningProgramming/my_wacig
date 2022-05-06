package my_vm

import (
	"fmt"
	"monkey/my_code"
	"monkey/my_compiler"
	"monkey/my_object"

	"golang.org/x/exp/constraints"
)

const StackSize = 2048

var (
	NULL  = &my_object.Null{}
	TRUE  = &my_object.Boolean{Value: true}
	FALSE = &my_object.Boolean{Value: false}
)

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
		case my_code.OpPop:
			vm.pop()
		case my_code.OpBang:
			switch vm.pop() {
			case FALSE:
				vm.push(TRUE)
			case TRUE:
				fallthrough
			default:
				vm.push(FALSE)
			}
		case my_code.OpMinus:
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
		case my_code.OpAdd, my_code.OpSub, my_code.OpDiv, my_code.OpMul:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case my_code.OpGT, my_code.OpGTE, my_code.OpEqual, my_code.OpNotEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case my_code.OpTrue:
			err := vm.push(TRUE)
			if err != nil {
				return err
			}
		case my_code.OpFalse:
			err := vm.push(FALSE)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) StackTop() my_object.Object {
	if vm.sp == 0 {
		return NULL
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackItem() my_object.Object {
	if vm.sp == StackSize {
		return NULL
	}
	return vm.stack[vm.sp]
}

type arithmeticFuncs struct {
	intFunc   func(a, b int64) int64
	floatFunc func(a, b float64) float64
}

type numberConstraint interface {
	~int64 | ~float64
}

func add[T numberConstraint](a, b T) T { return a + b }
func sub[T numberConstraint](a, b T) T { return a - b }
func mul[T numberConstraint](a, b T) T { return a * b }
func div[T numberConstraint](a, b T) T { return a / b }

var opToArithFuncs = map[my_code.Opcode]arithmeticFuncs{
	my_code.OpAdd: {
		intFunc:   add[int64],
		floatFunc: add[float64],
	},
	my_code.OpSub: {
		intFunc:   sub[int64],
		floatFunc: sub[float64],
	},
	my_code.OpDiv: {
		intFunc:   div[int64],
		floatFunc: div[float64],
	},
	my_code.OpMul: {
		intFunc:   mul[int64],
		floatFunc: mul[float64],
	},
}

// executeBinaryOperation now accepts arithmetic ops like + - * /
func (vm *VM) executeBinaryOperation(op my_code.Opcode) error {
	rightObj := vm.pop()
	leftObj := vm.pop()
	// below are almost the same as my_evaluator/eval_infix.go
	switch leftObj := leftObj.(type) {
	case *my_object.Integer:
		switch rightObj := rightObj.(type) {
		case *my_object.Integer:
			vm.push(&my_object.Integer{Value: opToArithFuncs[op].intFunc(leftObj.Value, rightObj.Value)})
		case *my_object.Boolean:
			vm.push(&my_object.Integer{Value: opToArithFuncs[op].intFunc(leftObj.Value, booleanToInt(rightObj.Value))})
		case *my_object.Float:
			vm.push(&my_object.Float{Value: opToArithFuncs[op].floatFunc(float64(leftObj.Value), rightObj.Value)})
		case *my_object.Null:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		default:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.Boolean:
		switch rightObj := rightObj.(type) {
		case *my_object.Integer:
			vm.push(&my_object.Integer{Value: opToArithFuncs[op].intFunc(booleanToInt(leftObj.Value), rightObj.Value)})
		case *my_object.Boolean:
			vm.push(&my_object.Integer{Value: opToArithFuncs[op].intFunc(booleanToInt(leftObj.Value), booleanToInt(rightObj.Value))})
		case *my_object.Float:
			vm.push(&my_object.Float{Value: opToArithFuncs[op].floatFunc(float64(booleanToInt(leftObj.Value)), rightObj.Value)})
		case *my_object.Null:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		default:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.Float:
		switch rightObj := rightObj.(type) {
		case *my_object.Integer:
			vm.push(&my_object.Float{Value: opToArithFuncs[op].floatFunc(leftObj.Value, float64(rightObj.Value))})
		case *my_object.Boolean:
			vm.push(&my_object.Float{Value: opToArithFuncs[op].floatFunc(leftObj.Value, float64(booleanToInt(rightObj.Value)))})
		case *my_object.Float:
			vm.push(&my_object.Float{Value: opToArithFuncs[op].floatFunc(leftObj.Value, rightObj.Value)})
		case *my_object.Null:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		default:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.String:
		if rightObj, ok := rightObj.(*my_object.String); ok && op == my_code.OpAdd {
			vm.push(&my_object.String{Value: leftObj.Value + rightObj.Value})
		} else {
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.Null:
		return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
	default:
		return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
	}
	return nil
}

type compFuncs struct {
	intFunc   func(a, b int64) bool
	floatFunc func(a, b float64) bool
}

func equal[T comparable](a, b T) bool                     { return a == b }
func notEqual[T comparable](a, b T) bool                  { return a != b }
func greaterThan[T constraints.Ordered](a, b T) bool      { return a > b }
func greaterThanEqual[T constraints.Ordered](a, b T) bool { return a >= b }

var opToCompFuncs = map[my_code.Opcode]compFuncs{
	my_code.OpGT:       {intFunc: greaterThan[int64], floatFunc: greaterThan[float64]},
	my_code.OpGTE:      {intFunc: greaterThanEqual[int64], floatFunc: greaterThanEqual[float64]},
	my_code.OpEqual:    {intFunc: equal[int64], floatFunc: equal[float64]},
	my_code.OpNotEqual: {intFunc: notEqual[int64], floatFunc: notEqual[float64]},
}

func (vm *VM) executeComparison(op my_code.Opcode) error {
	rightObj := vm.pop()
	leftObj := vm.pop()
	switch leftObj := leftObj.(type) {
	case *my_object.Integer:
		switch rightObj := rightObj.(type) {
		case *my_object.Integer:
			vm.push(booleanToConstObj(opToCompFuncs[op].intFunc(leftObj.Value, rightObj.Value)))
		case *my_object.Boolean:
			vm.push(booleanToConstObj(opToCompFuncs[op].intFunc(leftObj.Value, booleanToInt(rightObj.Value))))
		case *my_object.Float:
			vm.push(booleanToConstObj(opToCompFuncs[op].floatFunc(float64(leftObj.Value), rightObj.Value)))
		case *my_object.Null:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		default:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.Boolean:
		switch rightObj := rightObj.(type) {
		case *my_object.Integer:
			vm.push(booleanToConstObj(opToCompFuncs[op].intFunc(booleanToInt(leftObj.Value), rightObj.Value)))
		case *my_object.Boolean:
			vm.push(booleanToConstObj(opToCompFuncs[op].intFunc(booleanToInt(leftObj.Value), booleanToInt(rightObj.Value))))
		case *my_object.Float:
			vm.push(booleanToConstObj(opToCompFuncs[op].floatFunc(float64(booleanToInt(leftObj.Value)), rightObj.Value)))
		case *my_object.Null:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		default:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.Float:
		switch rightObj := rightObj.(type) {
		case *my_object.Integer:
			vm.push(booleanToConstObj(opToCompFuncs[op].floatFunc(leftObj.Value, float64(rightObj.Value))))
		case *my_object.Boolean:
			vm.push(booleanToConstObj(opToCompFuncs[op].floatFunc(leftObj.Value, float64(booleanToInt(rightObj.Value)))))
		case *my_object.Float:
			vm.push(booleanToConstObj(opToCompFuncs[op].floatFunc(leftObj.Value, rightObj.Value)))
		case *my_object.Null:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		default:
			return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
		}
	case *my_object.String:
		return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
	case *my_object.Null:
		return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
	default:
		return fmt.Errorf("unknown operator: %s%d%s", leftObj.Type(), op, rightObj.Type())
	}
	return nil
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
