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
	expectedConstants    []any
	expectedInstructions []my_code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []*compilerTestCase{
		{
			"1+2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpAdd),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1-2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpSub),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1*2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpMul),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1/2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpDiv),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"-2",
			[]any{2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpMinus),
				my_code.Make(my_code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFloatStringType(t *testing.T) {
	tests := []*compilerTestCase{
		{
			"1;2; 3.1415926;true;'hello world'",
			[]any{1, 2, 3.1415926, "hello world"},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpPop),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpPop),
				my_code.Make(my_code.OpConstant, 2),
				my_code.Make(my_code.OpPop),
				my_code.Make(my_code.OpTrue),
				my_code.Make(my_code.OpPop),
				my_code.Make(my_code.OpConstant, 3),
				my_code.Make(my_code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []*compilerTestCase{
		{
			"true",
			[]any{},
			[]my_code.Instructions{
				my_code.Make(my_code.OpTrue),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"false",
			[]any{},
			[]my_code.Instructions{
				my_code.Make(my_code.OpFalse),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"!false",
			[]any{},
			[]my_code.Instructions{
				my_code.Make(my_code.OpFalse),
				my_code.Make(my_code.OpBang),
				my_code.Make(my_code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestComparisonBooleanExpressions(t *testing.T) {
	tests := []*compilerTestCase{
		{
			"1<2",
			[]any{2, 1},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpGT),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1>2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpGT),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1>=2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpGTE),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1<=2",
			[]any{2, 1},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpGTE),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"1==2",
			[]any{1, 2},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpEqual),
				my_code.Make(my_code.OpPop),
			},
		},
		{
			"2!=1",
			[]any{2, 1},
			[]my_code.Instructions{
				my_code.Make(my_code.OpConstant, 0),
				my_code.Make(my_code.OpConstant, 1),
				my_code.Make(my_code.OpNotEqual),
				my_code.Make(my_code.OpPop),
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
	assert.EqualValues(t, len(concatted), len(actual), "wrong instruction length,\nwant:\n%s\ngot:\n%s\n", concatted, actual)
	for idx, ins := range concatted {
		assert.EqualValues(t, ins, actual[idx], "wrong instruction at %d,\nwant:\n%s\ngot:\n%s\n", idx, concatted, actual)
	}
}

func testConstants(t *testing.T, expected []any, actual []my_object.Object) {
	t.Helper()
	assert.EqualValues(t, len(expected), len(actual), "wrong number of constants")
	for idx, exp := range expected {
		switch exp := exp.(type) {
		case int:
			actualIntObj, ok := actual[idx].(*my_object.Integer)
			assert.True(t, ok, "expecting int object, got %s", actual[idx].Type())
			assert.EqualValues(t, exp, actualIntObj.Value)
		case float64:
			actualFloatObj, ok := actual[idx].(*my_object.Float)
			assert.True(t, ok)
			assert.EqualValues(t, exp, actualFloatObj.Value)
		case bool:
			actualBoolObj, ok := actual[idx].(*my_object.Boolean)
			assert.True(t, ok)
			assert.EqualValues(t, exp, actualBoolObj.Value)
		case string:
			actualStringObj, ok := actual[idx].(*my_object.String)
			assert.True(t, ok)
			assert.EqualValues(t, exp, actualStringObj.Value)
		}
	}
}

func parse(input string) *my_ast.Program {
	l := my_lexer.New(input)
	p := my_parser.New(l)
	return p.Parse()
}
