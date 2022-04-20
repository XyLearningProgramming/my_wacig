package my_compiler

import (
	"monkey/my_ast"
	"monkey/my_code"
	"monkey/my_lexer"
	"monkey/my_object"
	"monkey/my_parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []my_code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []*compilerTestCase{
		{
			"1+2",
			[]interface{}{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpAdd),
			},
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []*compilerTestCase) {
	t.Helper()
	for _, test := range tests {
		program := parse(test.input)
		compiler := New()
		err := compiler.Compile(program)
		assert.NoError(t, err)
		bytecode := compiler.ByteCode()
		testInstructions(t, test.expectedInstructions, bytecode.Instructions)
		testConstants(t, test.expectedConstants, bytecode.Constants)
	}
}

func testInstructions(t *testing.T, expected []my_code.Instructions, actual my_code.Instructions) {
	t.Helper()

	concatted := my_code.Instructions{}
	for _, exp := range expected {
		concatted = append(concatted, exp...)
	}
	assert.EqualValues(t, len(concatted), len(actual), "wrong instruction length,\nwant:\n%s\ngot:%s\n", concatted, actual)
	for idx, ins := range concatted {
		assert.EqualValues(t, ins, actual[idx], "wrong instruction at %d,\nwant:\n%s\ngot:\n%s\n", idx, concatted, actual)
	}
}

func testConstants(t *testing.T, expected []interface{}, actual []my_object.Object) {
	t.Helper()
	assert.EqualValues(t, len(expected), len(actual), "wrong number of constants")
	for idx, exp := range expected {
		switch exp := exp.(type) {
		case int:
			actualIntObj, ok := actual[idx].(*my_object.Integer)
			assert.True(t, ok, "expecting int object, got %s", actual[idx].Type())
			assert.EqualValues(t, exp, actualIntObj.Value)
		}
	}
}

func parse(input string) *my_ast.Program {
	l := my_lexer.New(input)
	p := my_parser.New(l)
	return p.Parse()
}
