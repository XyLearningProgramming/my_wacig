package my_code

import "fmt"

type Definition struct {
	// Name: make opcode readable
	Name string
	// OperandWidths: length: number of operands op code takes to execute;
	// element: number of bytes each operand takes up
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	// OpConstant: 2 bytes operands meaning holding 2**(16) max constants
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
	OpPop:      {"OpPop", []int{}},
	OpSub:      {"OpSub", []int{}},
	OpDiv:      {"OpDiv", []int{}},
	OpMul:      {"OpMul", []int{}},
	OpTrue:     {"OpTrue", []int{}},
	OpFalse:    {"OpFalse", []int{}},
	OpEqual:    {"OpEqual", []int{}},
	OpNotEqual: {"OpNotEqual", []int{}},
	OpGT:       {"OpGreaterThan", []int{}},
	OpGTE:      {"OpGreaterThanEqual", []int{}},
	OpMinus:    {"OpMinus", []int{}},
	OpBang:     {"OpBang", []int{}},
	// OpJumpNotTruthy: 1 operand with 2 bytes
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	// OpJump: 1 operand with 2 bytes
	OpJump: {"OpJump", []int{2}},
	OpNull: {"OpNull", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
