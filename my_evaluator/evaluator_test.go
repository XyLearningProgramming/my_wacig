package my_evaluator

import (
	"math"
	lexer "monkey/my_lexer"
	"monkey/my_object"
	"monkey/my_parser"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input  string
	expect any
}

type testCaseTyped struct {
	input   string
	expect  any
	refType reflect.Type
}

func TestEvalInteger(t *testing.T) {
	cases := []*testCaseTyped{
		{"5", 5, intType},
		{"10", 10, intType},
	}
	testCaseWithStruct(t, cases)
}

func TestEvalBoolean(t *testing.T) {
	cases := []testCase{
		{"true", true},
		{"false", false},
	}
	for _, c := range cases {
		evaluated := testEval(t, c.input)
		assert.EqualValues(t, c.expect, evaluated.(*my_object.Boolean).Value)
	}
}

func TestPrefixBangOperator(t *testing.T) {
	cases := []testCase{
		{"!true", false},
		{"!5", false},
		// {"!!a", true},
	}
	for _, c := range cases {
		evaluated := testEval(t, c.input)
		assert.EqualValues(t, c.expect, evaluated.(*my_object.Boolean).Value)
	}
}

func TestPrefixMinusOperator(t *testing.T) {
	cases := []testCase{
		{"-5", "-5"},
		{"5", "5"},
		{"-" + strconv.FormatUint(math.MaxUint64, 10), "1"},
		{"-" + strconv.FormatUint(math.MaxInt64, 10), "-" + strconv.FormatUint(math.MaxInt64, 10)},
		{"--" + strconv.FormatUint(math.MaxInt64, 10), strconv.FormatUint(math.MaxInt64, 10)},
		{strconv.FormatFloat(1.234, 'f', -1, 64), "1.234"},
		{strconv.FormatFloat(-1.234, 'f', -1, 64), "-1.234"},
	}
	for _, c := range cases {
		evaluated := testEval(t, c.input)
		assert.NotNil(t, evaluated, "case is: %+v", c)
		assert.EqualValues(t, c.expect, evaluated.String())
	}
}

var (
	mf        *my_object.Float
	mi        *my_object.Integer
	mb        *my_object.Boolean
	mn        *my_object.Null
	me        *my_object.Error
	ms        *my_object.String
	ma        *my_object.Array
	floatType = reflect.TypeOf(mf)
	intType   = reflect.TypeOf(mi)
	boolType  = reflect.TypeOf(mb)
	nullType  = reflect.TypeOf(mn)
	errType   = reflect.TypeOf(me)
	strType   = reflect.TypeOf(ms)
	arrType   = reflect.TypeOf(ma)
)

func TestInfixOperator(t *testing.T) {
	tests := []*testCaseTyped{
		{"1", 1, intType},
		{"2 * 2 * 2 * 2 * 2", 32, intType},
		{"-50 + 100 + -50", 0, intType},
		{"5 * 2 + 10", 20, intType},
		{"5 + 2 * 10", 25, intType},
		{"20 + 2 * -10", 0, intType},
		{"50 / 2 * 2 + 10", 60, intType},
		{"2 * (5 + 10)", 30, intType},
		{"3 * 3 * 3 + 10", 37, intType},
		{"3 * (3 * 3) + 10", 37, intType},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50, intType},
		{"5*5.0+1", 26, floatType},
		{"1.5+2-1", 2.5, floatType},
		{"true != \n\tfalse", true, boolType},
	}
	testCaseWithStruct(t, tests)
}

func TestIfElseStatements(t *testing.T) {
	tests := []*testCaseTyped{
		{"if (true) { 10 }", 10, intType},
		{"if (false) { 10 }", nil, nullType},
		{"if (1) { 10 }", 10, intType},
		{"if (1 < 2) { 10 }", 10, intType},
		{"if (1 > 2) { 10 }", nil, nullType},
		{"if (1 > 2) { 10 } else { 20 }", 20, intType},
		{"if (1 < 2) { 10 } else { 20 }", 10, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestReturnStatements(t *testing.T) {
	tests := []*testCaseTyped{
		{"return 10", 10, intType},
		{"return 2*5; 9", 10, intType},
		{"9; return 2*5; 9;", 10, intType},
		{"if(true){if(10>1){return 10;} return 1;}", 10, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestErrorHandling(t *testing.T) {
	tests := []*testCaseTyped{
		{"return -if(false){10}", "unknown operator: -NULL", errType},
		{"return NULL", "identifier not found: NULL", errType},
	}
	testCaseWithStruct(t, tests)
}

func TestLetStatements(t *testing.T) {
	tests := []*testCaseTyped{
		{"let a = 5; a;", 5, intType},
		{"let a = 5*5; 5*a", 5 * 5 * 5, intType},
		{"let a = 5; let b = a +5; b;", 10, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestEvalFunctionObject(t *testing.T) {
	input := "fn(x){x+2};"
	evaluated := testEval(t, input)
	fn, ok := evaluated.(*my_object.Function)
	assert.True(t, ok)
	assert.Equal(t, 1, len(fn.Parameters))
	assert.Equal(t, "x", fn.Parameters[0].String())
	assert.Equal(t, "{(x+2);}", fn.Body.String())
}

func TestFunctionEvaluation(t *testing.T) {
	tests := []*testCaseTyped{
		{"let identity = fn(x) { x; }; identity(5);", 5, intType},
		{"let identity = fn(x) { return x; }; identity(5);", 5, intType},
		{"let double = fn(x) { x * 2; }; double(5);", 10, intType},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10, intType},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20, intType},
		{"fn(x) { x; }(5)", 5, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestStringEvaluation(t *testing.T) {
	tests := []*testCaseTyped{
		{`"Hello\tWorld!\n"`, "Hello\tWorld!\n", strType},
		{"\"Hello\"+ \t\"World\"", "HelloWorld", strType},
		{"'Hello'- \n'World'", "unknown operator: STRING-STRING", errType},
	}
	testCaseWithStruct(t, tests)
}

func TestBuiltinLenFunction(t *testing.T) {
	tests := []*testCaseTyped{
		{`len("Hello\tWorld!\n")`, 13, intType},
		{"len(1)", "argument to len not supported: got INT", errType},
		{"len(\"one\", \"two\")", "wrong number of arguments: got=2, want=1", errType},
	}
	testCaseWithStruct(t, tests)
}

func TestArrayEvaluation(t *testing.T) {
	tests := []*testCaseTyped{
		{"[1, 2*2, 3+3]", []any{1, 4, 6}, arrType},
		{"[]", []any{}, arrType},
	}
	testCaseWithStruct(t, tests)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"[1, 2, 3, 4][0]", 1, intType},
		{"[1, 2, 3, 4][-1]", 4, intType},
		{"[1, 2, 3, 4][4]", "index 4 out of array with length 4", errType},
		{"[1, 2, 3, 4][1:2]", []any{2}, arrType},
		{"[1, 2, 3, 4][1:3]", []any{2, 3}, arrType},
		{"[1, 2, 3, 4][2:5]", []any{3, 4}, arrType},
		{"[1, 2, 3, 4][2:1]", []any{}, arrType},
		{"[1, 2, 3, 4][0:0]", []any{}, arrType},
		{"[1, 2, 3, 4][:]", []any{1, 2, 3, 4}, arrType},
		{"[1, 2, 3, 4][::]", []any{1, 2, 3, 4}, arrType},
		{"[1, 2, 3, 4][]", "array-like indexing with empty expression", errType},
		{"[1, 2, 3, 4][:5]", []any{1, 2, 3, 4}, arrType},
		{"[1, 2, 3, 4][::1]", []any{1, 2, 3, 4}, arrType},
		{"[1, 2, 3, 4][::2]", []any{1, 3}, arrType},
		{"[1, 2, 3, 4][::5]", []any{1}, arrType},
		{"[1, 2, 3, 4][::0]", "array-like indexing expecting non-zero stride", errType},
		{"[1, 2, 3, 4][::-1]", []any{4, 3, 2, 1}, arrType},
		{"[1, 2, 3, 4][::-3]", []any{4, 1}, arrType},
		{"[1, 2, 3, 4][1::-3]", []any{2}, arrType},
	}
	testCaseWithStruct(t, tests)
}

func TestHashExpression(t *testing.T) {
	input := `let two="three"; {"one": 10-9, two: 1+1, three: 3, 4: 4, true:5, false:6, 7.1: 7}`
	evaluated := testEval(t, input)
	hashObj, hok := evaluated.(*my_object.Hash)
	assert.True(t, hok)
	expected := map[my_object.HashKey]int64{
		(&my_object.String{Value: "one"}).HashKey():   1,
		(&my_object.String{Value: "two"}).HashKey():   2,
		(&my_object.String{Value: "three"}).HashKey(): 3,
		(&my_object.Integer{Value: 4}).HashKey():      4,
		(&my_object.Boolean{Value: true}).HashKey():   5,
		(&my_object.Boolean{Value: false}).HashKey():  6,
		(&my_object.Float{Value: 7.1}).HashKey():      7,
	}
	for ek, ev := range expected {
		pair, ok := hashObj.Pairs[ek]
		assert.True(t, ok, "pair: %+v", pair)
		assert.EqualValues(t, ev, pair.Value.(*my_object.Integer).Value)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"{foo:5}[\"foo\"]", 5, intType},
		{"{foo:5}[\"bar\"]", nil, nullType},
		{"{}[\"bar\"]", nil, nullType},
		{"{5:5}[5]", 5, intType},
		{"{5.0:5}[5.0]", 5, intType},
		{"{true:5}[true]", 5, intType},
		{"{true:5}[fn(x){x}]", "key type not hashable: FUNCTION", errType},
	}
	testCaseWithStruct(t, tests)
}

func TestForLoopExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"for(;false;){}", "", nullType},
		{"let b = 1;for(let a=1;a<3; a = a+1){let b = a; return b;}", 2, intType},
		{"let b = 1;for(let a=1;a<3;){a = a+1; b= a; b}", 3, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestWhileLoopExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"while(false){}", "", nullType},
		{"let a = 1;while(a<3){a = a+1; a;}", 3, intType},
		{"let a = 1;while(a<3){a = a+1; a;}", 3, intType},
		{"let a = 1;while(a==1){let a = 2; break;};a;", 1, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestDoWhileLoopExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"do{1}while(false)", 1, intType},
		{"let a=1;do{a=a+1; a;}while(a<3)", 3, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestGTELTEExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"1>=2", false, boolType},
		{"1<=0.9", false, boolType},
		{"1<=0.9<=0.1", true, boolType},
	}
	testCaseWithStruct(t, tests)
}

func TestReassignExpression(t *testing.T) {
	tests := []*testCaseTyped{
		{"let a= true; a=2; a", 2, intType},
		{"let a= true; a=2;", 2, intType},
		{"a=2", "cannot assign to undefined identifier", errType},
		{"let a= true; a>=false", true, boolType},
		{"let a = 1; do{let a= 2}while(false);a=3;a;", 3, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestBreakContinueStatements(t *testing.T) {
	tests := []*testCaseTyped{
		{"break", "break outside loop", errType},
		{"continue", "continue outside loop", errType},
		{"for(;true;){1;break;}", "", nullType},
		{"let a = 1;while(a<3){a=a+1;break;};a;", 2, intType},
		{"let a = 1;while(a<3){a=a+1;continue;}", "", nullType},
		{"let a = 1;while(a<3){a=a+1;continue;};a;", 3, intType},
	}
	testCaseWithStruct(t, tests)
}

func TestNullStatements(t *testing.T) {
	tests := []*testCaseTyped{
		{"null", "", nullType},
		{"null==null", true, boolType},
		{"null!=null;", false, boolType},
		{"null>null", false, boolType},
		{"null<=null", true, boolType},
		{"null+null", "unknown operator: NULL+NULL", errType},
	}
	testCaseWithStruct(t, tests)
}

func testCaseWithStruct(t *testing.T, tests []*testCaseTyped) {
	for _, c := range tests {
		evaluated := testEval(t, c.input)
		assert.NotNil(t, evaluated, "case is: %+v", c)
		ev := reflect.ValueOf(evaluated)
		assert.EqualValues(t, c.refType, ev.Type())
		testOneCaseWithStruct(t, c.expect, c.refType, ev)
	}
}

func testOneCaseWithStruct(t *testing.T, expect any, expectType reflect.Type, actualValue reflect.Value) {
	switch expectType {
	case floatType:
		assert.EqualValues(t, expect, actualValue.Elem().Field(0).Float())
	case intType:
		assert.EqualValues(t, expect, actualValue.Elem().Field(0).Int())
	case boolType:
		assert.EqualValues(t, expect, actualValue.Elem().Field(0).Bool())
	case nullType:
		assert.EqualValues(t, 0, actualValue.Elem().NumField())
	case errType:
		assert.EqualValues(t, expect, actualValue.Elem().Field(0).String())
	case strType:
		assert.EqualValues(t, expect, actualValue.Elem().Field(0).String())
	case arrType:
		carr, cok := expect.([]any)
		assert.True(t, cok)
		// fmt.Printf("carr = %+v: actualValue: %+v\n", carr, actualValue)
		if len(carr) == 0 {
			assert.Equal(t, 0, actualValue.Elem().Field(0).Len())
			return
		}
		for idx, celem := range carr {
			slicedActual := actualValue.Elem().Field(0).Index(idx)
			testOneCaseWithStruct(t, celem, slicedActual.Elem().Type(), slicedActual.Elem())
		}
	default:
		assert.Fail(t, "unknown ev: %T: %v: c: %+v", actualValue, actualValue)
	}
}

func testEval(t *testing.T, input string) my_object.Object {
	l := lexer.New(input)
	p := my_parser.New(l)
	prog := p.Parse()
	assert.NoError(t, p.Error())
	// fmt.Println(prog.String())
	return Eval(prog, my_object.NewEnvironment())
}
