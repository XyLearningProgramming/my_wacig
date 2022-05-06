package my_code

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisassembleString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, math.MaxUint16),
		Make(OpAdd),
	}
	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
`
	assert.EqualValues(t, expected, Concat(instructions).String())
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{math.MaxUint16}, 2},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		assert.NoError(t, err)
		operandsRead, n := ReadOperands(def, instruction[1:])
		assert.EqualValues(t, tt.bytesRead, n)
		assert.EqualValues(t, tt.operands, operandsRead)
	}
}
