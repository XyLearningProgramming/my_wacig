package my_code

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Instructions []byte

type Opcode byte

const (
	// OpConstant: retrieves the constant using operand as index and push it onto the stack
	OpConstant Opcode = iota
)

func (ins Instructions) String() string {
	sb := &strings.Builder{}
	for i := 0; i < len(ins); {
		def, err := Lookup(ins[i])
		if err != nil {
			sb.WriteString(fmt.Sprintf("ERROR: %s\n", err))
		}
		operands, bytesRead := ReadOperands(def, ins[i+1:])
		sb.WriteString(fmt.Sprintf("%04d %s\n", i, fmtIns(def, operands)))
		i += 1 + bytesRead
	}
	return sb.String()
}

// ReadOperands: read operands out by the definition of op code from the start of instruction
func ReadOperands(def *Definition, ins Instructions) (operands []int, offset int) {
	operands = make([]int, len(def.OperandWidths))
	offset = 0
	for idx, byteWidth := range def.OperandWidths {
		switch byteWidth {
		case 2:
			operands[idx] = int(ReadUint16(ins[offset:]))
		}
		offset += byteWidth
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func Concat(instructions []Instructions) Instructions {
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}
	return concatted
}

func fmtIns(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if operandCount != len(operands) {
		return fmt.Sprintf("ERROR: operand read length %d does not match defined %d", len(operands), operandCount)
	}
	switch operandCount {
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s", def.Name)
}
