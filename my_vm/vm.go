package my_vm

import (
	"fmt"
	"math"
	"monkey/my_code"
	"monkey/my_compiler"
	"monkey/my_object"
)

const (
	StackSize  = 2048
	GlobalSize = math.MaxUint16 + 1
)

type VM struct {
	stack        []my_object.Object
	sp           int // sp: points to the next value, top of stack is stack[sp-1]
	instructions my_code.Instructions
	constants    []my_object.Object
	globals      []my_object.Object
}

func New(byteCode *my_compiler.ByteCode) *VM {
	return &VM{
		sp:           0,
		stack:        make([]my_object.Object, StackSize),
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		globals:      NewGlobals(),
	}
}

func NewWithState(byteCode *my_compiler.ByteCode, globals []my_object.Object) *VM {
	return &VM{
		sp:           0,
		stack:        make([]my_object.Object, StackSize),
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		globals:      globals,
	}
}

func NewGlobals() []my_object.Object {
	return make([]my_object.Object, GlobalSize)
}

func (vm *VM) Run() error {
	// fetch-decode-execute
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := my_code.Opcode(vm.instructions[ip])
		switch op {
		// variable-related
		case my_code.OpConstant:
			constIndex := my_code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case my_code.OpSetGlobal:
			globalIdx := my_code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIdx] = vm.pop()
		case my_code.OpGetGlobal:
			globalIdx := my_code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIdx])
			if err != nil {
				return err
			}

		// calculations
		case my_code.OpBang:
			vm.executeOpBang()
		case my_code.OpMinus:
			err := vm.executeOpMinus()
			if err != nil {
				return err
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
		// functional
		case my_code.OpJump:
			jumpToPos := my_code.ReadUint16(vm.instructions[ip+1:])
			ip = int(jumpToPos) - 1 // ip has ++ after each loop
		case my_code.OpJumpNotTruthy:
			// if condition is true
			if !isTruthy(vm.pop()) {
				jumpToPos := my_code.ReadUint16(vm.instructions[ip+1:])
				ip = int(jumpToPos) - 1
			} else {
				ip += 2 // ip has ++ after each loop, so +2 jumps over the OpJumpNotTruthy
			}
		case my_code.OpPop:
			vm.pop()
		// constants
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
		case my_code.OpNull:
			err := vm.push(NULL)
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
