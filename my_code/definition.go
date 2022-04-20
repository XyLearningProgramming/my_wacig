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
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
