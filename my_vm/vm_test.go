package my_vm

import (
	"monkey/my_ast"
	"monkey/my_compiler"
	"monkey/my_lexer"
	"monkey/my_object"
	"monkey/my_parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []*vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1+2", 2}, // TODO
	}
	runVMTests(t, tests)
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVMTests(t *testing.T, tests []*vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := my_compiler.New()
		err := comp.Compile(program)
		assert.NoError(t, err)
		vm := New(comp.ByteCode())
		err = vm.Run()
		assert.NoError(t, err)
		stackElem := vm.StackTop()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual my_object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		intObj, ok := actual.(*my_object.Integer)
		assert.True(t, ok)
		assert.EqualValues(t, expected, intObj.Value)
	}
}

func parse(input string) *my_ast.Program {
	l := my_lexer.New(input)
	p := my_parser.New(l)
	return p.Parse()
}
