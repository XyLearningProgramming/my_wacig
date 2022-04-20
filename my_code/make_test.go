package my_code

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{math.MaxUint16 - 1}, []byte{byte(OpConstant), 255, 254}},
	}
	for _, test := range tests {
		instruction := Make(test.op, test.operands...)
		assert.EqualValues(t, len(test.expected), len(instruction), "instruction has wrong length: test: %+v", test)
		for i, b := range test.expected {
			assert.EqualValues(t, b, instruction[i], "wrong byte at pos: %d: test: %+v", i, test)
		}
	}
}
