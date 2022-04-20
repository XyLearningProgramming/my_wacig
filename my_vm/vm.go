package my_vm

import (
	"fmt"
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
