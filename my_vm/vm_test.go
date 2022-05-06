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
		{"1+2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVMTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []*vmTestCase{
		{"true", true},
		{"false", false},
		{"1<=2", true},
		{"1<2", true},
		{"1>=2", false},
		{"1>2", false},
		{"1==2", false},
		{"1!=2", true},
		{"true!=false", true},
		{"(1<=2)==false", false},
		{"(1>2)==true", false},
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
	}
	runVMTests(t, tests)
}

func TestNumberStringAddArithmetic(t *testing.T) {
	tests := []*vmTestCase{
		{"1.0+4", 5.0},
		{"true+true", 2},
		{"2.0+false", 2.0},
		{`"hello"+' '+"world"`, "hello world"},
	}
	runVMTests(t, tests)
}

type vmTestCase struct {
	input    string
	expected any
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
		stackElem := vm.LastPoppedStackItem()
		testExpectedObject(t, tt.expected, stackElem, "input=%s", tt.input)
	}
}

func testExpectedObject(t *testing.T, expected any, actual my_object.Object, msgAndArgs ...interface{}) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		intObj, ok := actual.(*my_object.Integer)
		assert.True(t, ok, "want integer but got: %v", actual)
		assert.EqualValues(t, expected, intObj.Value, msgAndArgs...)
	case float64:
		floatObj, ok := actual.(*my_object.Float)
		assert.True(t, ok, "want float obj, got: %s", actual.Type())
		assert.EqualValues(t, expected, floatObj.Value, msgAndArgs...)
	case bool:
		boolObj, ok := actual.(*my_object.Boolean)
		assert.True(t, ok)
		assert.EqualValues(t, expected, boolObj.Value, msgAndArgs...)
	case string:
		strObj, ok := actual.(*my_object.String)
		assert.True(t, ok)
		assert.EqualValues(t, expected, strObj.Value, msgAndArgs...)
	}
}

func parse(input string) *my_ast.Program {
	l := my_lexer.New(input)
	p := my_parser.New(l)
	return p.Parse()
}
