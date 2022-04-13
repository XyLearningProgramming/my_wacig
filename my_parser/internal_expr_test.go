package my_parser

import (
	"monkey/my_ast"
	lexer "monkey/my_lexer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExpressionStatement(t *testing.T) {
	input := "a;b; let a=b"
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	assert.NotNil(t, prog)
	assert.NotNil(t, prog.Statements)
	assert.Equal(t, 3, len(prog.Statements))
	assert.Nil(t, p.err)
	assert.Equal(t, "a", prog.Statements[0].(*my_ast.ExpressionStatement).Expression.(*my_ast.Identifier).Value)
	assert.Equal(t, "b", prog.Statements[1].(*my_ast.ExpressionStatement).Expression.(*my_ast.Identifier).Value)
	letStmt, lok := prog.Statements[2].(*my_ast.LetStatement)
	assert.True(t, lok)
	assert.Equal(t, "a", letStmt.Ident.Value)
	assert.Equal(t, "b", letStmt.Value.(*my_ast.Identifier).String())
}

func TestParseNumberStatement(t *testing.T) {
	input := "1;1.234"
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	assert.NotNil(t, prog)
	assert.NotNil(t, prog.Statements)
	assert.Equal(t, 2, len(prog.Statements))
	assert.Nil(t, p.err)
	assert.EqualValues(t, 1, prog.Statements[0].(*my_ast.ExpressionStatement).Expression.(*my_ast.Integer).Value)
	assert.EqualValues(t, 1.234, prog.Statements[1].(*my_ast.ExpressionStatement).Expression.(*my_ast.Float).Value)
}

func TestPrefixExpressionStatement(t *testing.T) {
	input := "!5;\n-15.5;"
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	assert.NotNil(t, prog)
	assert.NotNil(t, prog.Statements)
	assert.Equal(t, 2, len(prog.Statements))
	assert.Nil(t, p.err)
	prefixNode, pok := prog.Statements[0].(*my_ast.ExpressionStatement).
		Expression.(*my_ast.PrefixExpression)
	assert.True(t, pok)
	assert.EqualValues(t, my_ast.PREOP_BANG, prefixNode.Operator)
	assert.EqualValues(t, 5, prefixNode.Right.(*my_ast.Integer).Value)
	prefixNode, pok = prog.Statements[1].(*my_ast.ExpressionStatement).
		Expression.(*my_ast.PrefixExpression)
	assert.True(t, pok)
	assert.EqualValues(t, my_ast.PREOP_MINUS, prefixNode.Operator)
	assert.EqualValues(t, 15.5, prefixNode.Right.(*my_ast.Float).Value)
}

func TestInfixExpressionStatement(t *testing.T) {
	tests := []struct {
		input     string
		leftExpr  interface{}
		rightExpr interface{}
		operator  string
	}{
		{"5+5", 5, 5, "+"},
		{"5-5", 5, 5, "-"},
		{"5 * 5", 5, 5, "*"},
		{"5/\r\n5", 5, 5, "/"},
		{"5>   5", 5, 5, ">"},
		{"5\t<5", 5, 5, "<"},
		{"5== 5;", 5, 5, "=="},
		{"5 !=5", 5, 5, "!="},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		prog := p.Parse()
		assert.Nil(t, p.Error())
		assert.NotNil(t, prog)
		assert.NotNil(t, prog.Statements)
		assert.Equal(t, 1, len(prog.Statements))
		infixNode := prog.Statements[0].(*my_ast.ExpressionStatement).Expression.(*my_ast.InfixExpression)
		assert.NotNil(t, infixNode)
		assert.EqualValues(t, test.leftExpr, infixNode.Left.(*my_ast.Integer).Value)
		assert.EqualValues(t, test.leftExpr, infixNode.Right.(*my_ast.Integer).Value)
		assert.EqualValues(t, test.operator, infixNode.Operator)
	}
}

type TestWithExpect struct {
	input  string
	expect string
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []TestWithExpect{
		{"-b*c", "(-(b*c));"},
		{"a*b-c", "((a*b)-c);"},
		{"!-c", "(!(-c));"},
		{"-1+2", "((-1)+2);"},
	}
	testSingleStringedStatements(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"return true", "return true;"},
		{"true + false", "(true+false);"},
	}
	testSingleStringedStatements(t, tests)
}

func TestGroupedExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"1+(2+ 3) +4", "((1+(2+3))+4);"},
		{"(5 +5 )/2", "((5+5)/2);"},
		{"-(\t5+ \t5)", "(-(5+5));"},
		{"!(true == true)", "(!(true==true));"},
	}
	testSingleStringedStatements(t, tests)
}
func TestIfExpression(t *testing.T) {
	input := "if (x< y) { x} else\n\n{x;return y;}"
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	assert.Nil(t, p.Error())
	assert.NotNil(t, prog)
	assert.NotNil(t, prog.Statements)
	assert.Equal(t, 1, len(prog.Statements))
	assert.Equal(t, "if((x<y)){x;}else{x;return y;};", prog.Statements[0].String())
	es, eok := prog.Statements[0].(*my_ast.ExpressionStatement)
	assert.True(t, eok)
	is, iok := es.Expression.(*my_ast.IfExpression)
	assert.True(t, iok)
	assert.Equal(t, "(x<y)", is.Condition.String())
	assert.Equal(t, "x", is.Consequence.Statements[0].(*my_ast.ExpressionStatement).Expression.(*my_ast.Identifier).Value)
	assert.Equal(t, "y", is.Alternative.Statements[1].(*my_ast.ReturnStatement).Value.(*my_ast.Identifier).Value)
}

func testSingleStringedStatements(t *testing.T, tests []TestWithExpect) {
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		prog := p.Parse()
		assert.Nil(t, p.Error())
		assert.NotNil(t, prog)
		assert.NotNil(t, prog.Statements)
		assert.Equal(t, 1, len(prog.Statements))
		assert.Equal(t, test.expect, prog.Statements[0].String())
	}
}

func testMultipleStringedStatements(t *testing.T, tests []TestWithExpect) {
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		prog := p.Parse()
		assert.Nil(t, p.Error())
		assert.NotNil(t, prog)
		assert.NotNil(t, prog.Statements)
		sb := &strings.Builder{}
		for _, stmt := range prog.Statements {
			sb.WriteString(stmt.String())
		}
		assert.Equal(t, test.expect, sb.String())
	}
}

func TestParseLetReturn(t *testing.T) {
	input := `
	let a = b; return c
	`
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	assert.NotNil(t, prog)
	assert.NotNil(t, prog.Statements)
	assert.Equal(t, 2, len(prog.Statements))
	assert.Nil(t, p.err)
	assert.Equal(t, "a", prog.Statements[0].(*my_ast.LetStatement).Ident.Value)
	assert.Equal(t, "b", prog.Statements[0].(*my_ast.LetStatement).Value.(*my_ast.Identifier).Value)
	assert.Equal(t, "c", prog.Statements[1].(*my_ast.ReturnStatement).Value.(*my_ast.Identifier).Value)
}

func TestParseFunctionExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"fn(x,y)\n{x+y;}", "fn(x,y){(x+y);};"},
		{"fn(x){return x+2}", "fn(x){return (x+2);};"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseCallExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"fn(x,y)\n{x+y;}(1,\ta)", "fn(x,y){(x+y);}(1,a);"},
		{"add(1,2* 3, 4 + 5,add(1+2))", "add(1,(2*3),(4+5),add((1+2)));"},
	}
	testSingleStringedStatements(t, tests)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello\tWorld!\n";`
	l := lexer.New(input)
	p := New(l)
	prog := p.Parse()
	assert.NotNil(t, prog)
	assert.Nil(t, p.err)
	assert.NotNil(t, prog.Statements)
	assert.Equal(t, 1, len(prog.Statements))
	es, eok := prog.Statements[0].(*my_ast.ExpressionStatement)
	assert.True(t, eok)
	ss, sok := es.Expression.(*my_ast.StringExpression)
	assert.True(t, sok)
	assert.Equal(t, "Hello\tWorld!\n", ss.Value)
}

func TestParseArrayExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"[1, 2*2, !false]", "[1,(2*2),(!false)];"},
		{"[1, 2*2, !false] + [1]", "([1,(2*2),(!false)]+[1]);"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseArrayIndexingExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"[1+1][0]", "([(1+1)][0]);"},
		{"[1][0:1]", "([1][0:1]);"},
		{"[1][:1]", "([1][:1]);"},
		{"[1][::1+1]", "([1][::(1+1)]);"},
		{"[1][::]", "([1][::]);"},
		{"[1][]", "([1][]);"},
		{"[1][:]", "([1][:]);"},
		{"[1][1::]", "([1][1::]);"},
		{"[1][:1:]", "([1][:1:]);"},
		{"a*[1,2,3,4][b*c]*d", "((a*([1,2,3,4][(b*c)]))*d);"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseMapExpression(t *testing.T) {
	tests := []TestWithExpect{
		{`{"one": 1, two: 2+1, 3: [1,2,3][:]}`, `{one:1,two:(2+1),3:([1,2,3][:])};`},
		{"{}", "{};"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseForExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"for(;;){}", "for(;;){};"},
		{"for(let y=x+1;y!=2;let y=y+1){y*2}", "for(let y = (x+1);(y!=2);let y = (y+1)){(y*2);};"},
		{"for(;y!=2;let y=y+1){y*2}", "for(;(y!=2);let y = (y+1)){(y*2);};"},
		{"for(;;let y=y+1){y*2}", "for(;;let y = (y+1)){(y*2);};"},
		{"for(;y!=2;){y*2}", "for(;(y!=2);){(y*2);};"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseWhileExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"while(){}", "while(){};"},
		{"while(){y+2}", "while(){(y+2);};"},
		{"while(x+1<2){y+2}", "while(((x+1)<2)){(y+2);};"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseDoWhileExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"do{}while()", "do{}while();"},
		{"do{y+2}while()", "do{(y+2);}while();"},
		{"do{y+2}while(y<2)", "do{(y+2);}while((y<2));"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseGTELTEExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"a<=1+1", "(a<=(1+1));"},
		{"a>=2*2", "(a>=(2*2));"},
		{"a>=2>2", "((a>=2)>2);"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseReassignExpression(t *testing.T) {
	tests := []TestWithExpect{
		{"a=a+1", "(a=(a+1));"},
	}
	testSingleStringedStatements(t, tests)
}

func TestParseBreakContinueStatement(t *testing.T) {
	tests := []TestWithExpect{
		{"break;1+2", "break;(1+2);"},
		{"1+1;continue;1+2", "(1+1);continue;(1+2);"},
	}
	testMultipleStringedStatements(t, tests)
}
