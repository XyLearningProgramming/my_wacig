package my_ast

import (
	token "monkey/my_token"
	"strconv"
	"strings"
)

type Node interface {
	// DebugString: for debugging purposes, usually literal value of a token;
	// if an expression, ususally literal value of the first token;
	// if an statement, usually call DebugString() on its first expression;
	// if an ident, literal value of its `Value` field;
	DebugString() string
	// String: output the formatted codes if legal;
	// it has the following rules:
	// one statement per line;
	// each line will end with semicolon;
	// each token has a space in between;
	String() string
}

const (
	NodeStringNewLine    = "\n"
	NodeStringSemiColon  = ";"
	NodeStringTokenSpace = " "
)

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

// root node

type Program struct {
	Statements []Statement
}

func (p *Program) DebugString() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].DebugString()
	}
	return ""
}

func (p *Program) String() string {
	sb := strings.Builder{}
	for idx, s := range p.Statements {
		sb.WriteString(s.String())
		if idx != len(p.Statements)-1 {
			sb.WriteString(NodeStringNewLine)
		}
	}
	return sb.String()
}

// statements

type LetStatement struct {
	Ident *Identifier
	Value Expression
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) DebugString() string {
	return l.Ident.DebugString()
}

func (l *LetStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(token.LookupKeywords(token.LET))
	sb.WriteString(NodeStringTokenSpace)
	sb.WriteString(l.Ident.Value)
	sb.WriteString(NodeStringTokenSpace)
	sb.WriteString(token.REASSIGN)
	sb.WriteString(NodeStringTokenSpace)
	sb.WriteString(l.Value.String())
	sb.WriteString(NodeStringSemiColon)
	return sb.String()
}

type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) DebugString() string {
	return r.Value.DebugString()
}

func (r *ReturnStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(token.LookupKeywords(token.RETURN))
	sb.WriteString(NodeStringTokenSpace)
	sb.WriteString(r.Value.String())
	sb.WriteString(NodeStringSemiColon)
	return sb.String()
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) DebugString() string {
	if e.Expression == nil {
		return ""
	}
	return e.Expression.DebugString()
}

func (e *ExpressionStatement) String() string {
	if e.Expression == nil {
		return ""
	}
	return e.Expression.String() + ";"
}

// expressions

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) DebugString() string {
	return i.Value
}

func (i *Identifier) String() string {
	return i.Value
}

type Integer struct {
	Value uint64
}

func (i *Integer) expressionNode() {}

func (i *Integer) DebugString() string {
	return strconv.FormatUint(i.Value, 10)
}

func (i *Integer) String() string {
	return strconv.FormatUint(i.Value, 10)
}

type Float struct {
	Value float64
}

func (f *Float) expressionNode() {}

func (f *Float) DebugString() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

func (f *Float) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

type PrefixOperator string

const (
	PREOP_MINUS PrefixOperator = token.MINUS
	PREOP_BANG  PrefixOperator = token.BANG
)

type PrefixExpression struct {
	Operator PrefixOperator
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) DebugString() string {
	return string(p.Operator)
}
func (p *PrefixExpression) String() string {
	sb := strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(string(p.Operator))
	sb.WriteString(p.Right.String())
	sb.WriteRune(')')
	return sb.String()
}

type InfixOperator string

const (
	INOP_MINUS      InfixOperator = token.MINUS
	INOP_PLUS       InfixOperator = token.PLUS
	INOP_ASTERISK   InfixOperator = token.ASTERISK
	INOP_SLASH      InfixOperator = token.SLASH
	INOP_LT         InfixOperator = token.LT
	INOP_GT         InfixOperator = token.GT
	INOP_EQ         InfixOperator = token.EQ
	INOP_NOT_EQ     InfixOperator = token.NOT_EQ
	INOP_CALL       InfixOperator = token.LPAREN
	INOP_INDEX      InfixOperator = token.LBRACKET
	INOP_INDEXCOLON InfixOperator = token.COLON
	INOP_GTE        InfixOperator = token.GTE
	INOP_LTE        InfixOperator = token.LTE
	INOP_REASSIGN   InfixOperator = token.REASSIGN
)

type InfixExpression struct {
	Operator InfixOperator
	Left     Expression
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) DebugString() string {
	return string(i.Operator)
}

func (i *InfixExpression) String() string {
	sb := strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(i.Left.String())
	sb.WriteString(string(i.Operator))
	sb.WriteString(i.Right.String())
	sb.WriteRune(')')
	return sb.String()
}

type Boolean struct {
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) DebugString() string {
	return strconv.FormatBool(b.Value)
}

func (b *Boolean) String() string {
	return strconv.FormatBool(b.Value)
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) DebugString() string {
	return token.LBRACE
}

func (b *BlockStatement) String() string {
	sb := strings.Builder{}
	sb.WriteRune('{')
	for _, stmt := range b.Statements {
		sb.WriteString(stmt.String())
	}
	sb.WriteRune('}')
	return sb.String()
}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}

func (i *IfExpression) DebugString() string {
	return token.IF
}

func (i *IfExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString("if")
	sb.WriteString("(")
	sb.WriteString(i.Condition.String())
	sb.WriteString(")")
	sb.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		sb.WriteString("else")
		sb.WriteString(i.Alternative.String())
	}
	return sb.String()
}

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *Function) expressionNode() {}

func (f *Function) DebugString() string {
	if len(f.Parameters) > 0 {
		return f.Parameters[0].DebugString()
	}
	return ""
}

func (f *Function) String() string {
	sb := strings.Builder{}
	sb.WriteString("fn")
	sb.WriteRune('(')
	for idx, p := range f.Parameters {
		sb.WriteString(p.String())
		if idx != len(f.Parameters)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(')')
	sb.WriteString(f.Body.String())
	return sb.String()
}

type CallExpression struct {
	Function  Expression // Identifier or Function
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) DebugString() string {
	return c.Function.DebugString()
}

func (c *CallExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString(c.Function.String())
	sb.WriteRune('(')
	for idx, a := range c.Arguments {
		sb.WriteString(a.String())
		if idx != len(c.Arguments)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(')')
	return sb.String()
}

type StringExpression struct {
	Value string
}

func (sl *StringExpression) DebugString() string { return sl.Value }

func (sl *StringExpression) String() string { return sl.Value }

func (sl *StringExpression) expressionNode() {}

type ArrayExpression struct {
	Elements []Expression
}

func (ae *ArrayExpression) DebugString() string {
	return ae.String()
}

func (ae *ArrayExpression) String() string {
	elements := []string{}
	for _, el := range ae.Elements {
		elements = append(elements, el.String())
	}
	return "[" + strings.Join(elements, ",") + "]"
}

func (ae *ArrayExpression) expressionNode() {}

type IndexExpression struct {
	Left            Expression
	StartIndex      Expression // value specified by user for start index
	IsSetStartIndex bool       // if start index is set when 1. user specified a colon : 2. user specified a value explicitly
	EndIndex        Expression
	IsSetEndIndex   bool
	Stride          Expression
	IsSetStride     bool
}

func (aie *IndexExpression) DebugString() string { return aie.String() }

func (aie *IndexExpression) String() string {
	sb := &strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(aie.Left.String())
	sb.WriteRune('[')
	if aie.StartIndex != nil {
		sb.WriteString(aie.StartIndex.String())
	}
	if aie.IsSetEndIndex {
		sb.WriteRune(':')
		if aie.EndIndex != nil {
			sb.WriteString(aie.EndIndex.String())
		}
	}
	if aie.IsSetStride {
		sb.WriteRune(':')
		if aie.Stride != nil {
			sb.WriteString(aie.Stride.String())
		}
	}
	sb.WriteString("])")
	return sb.String()
}

func (aie *IndexExpression) expressionNode() {}

type HashExpression struct {
	Pairs map[Expression]Expression
	Keys  []Expression
}

func (he *HashExpression) DebugString() string { return he.String() }

func (he *HashExpression) String() string {
	pairs := []string{}
	for _, k := range he.Keys {
		pairs = append(pairs, k.String()+":"+he.Pairs[k].String())
	}
	return "{" + strings.Join(pairs, ",") + "}"
}

func (he *HashExpression) expressionNode() {}

type ForExpression struct {
	InitStmt   Statement
	TestExpr   Expression
	UpdateStmt Statement
	Body       *BlockStatement
}

func (fe *ForExpression) DebugString() string { return fe.String() }

func (fe *ForExpression) String() string {
	sb := &strings.Builder{}
	sb.WriteString("for(")
	if fe.InitStmt != nil {
		sb.WriteString(fe.InitStmt.String())
	} else {
		sb.WriteRune(';')
	}
	if fe.TestExpr != nil {
		sb.WriteString(fe.TestExpr.String())
	}
	sb.WriteRune(';')
	if fe.UpdateStmt != nil {
		sb.WriteString(fe.UpdateStmt.String())
	} else {
		sb.WriteRune(';')
	}
	trimmed := strings.TrimSuffix(sb.String(), ";")
	sb.Reset()
	sb.WriteString(trimmed)
	sb.WriteRune(')')
	if fe.Body != nil {
		sb.WriteString(fe.Body.String())
	} else {
		sb.WriteRune('{')
		sb.WriteRune('}')
	}
	return sb.String()
}

func (fe *ForExpression) expressionNode() {}

type DoWhileExpression struct {
	TestExpr Expression
	Body     *BlockStatement
}

func (dw *DoWhileExpression) DebugString() string { return dw.String() }

func (dw *DoWhileExpression) String() string {
	sb := &strings.Builder{}
	sb.WriteString("do")
	if dw.Body != nil {
		sb.WriteString(dw.Body.String())
	} else {
		sb.WriteRune('{')
		sb.WriteRune('}')
	}
	sb.WriteString("while(")
	if dw.TestExpr != nil {
		sb.WriteString(dw.TestExpr.String())
	}
	sb.WriteRune(')')
	return sb.String()
}

func (dw *DoWhileExpression) expressionNode() {}

type WhileExpression struct {
	TestExpr Expression
	Body     *BlockStatement
}

func (w *WhileExpression) DebugString() string { return w.String() }

func (w *WhileExpression) String() string {
	sb := &strings.Builder{}
	sb.WriteString("while")
	sb.WriteRune('(')
	if w.TestExpr != nil {
		sb.WriteString(w.TestExpr.String())
	}
	sb.WriteRune(')')
	if w.Body != nil {
		sb.WriteString(w.Body.String())
	} else {
		sb.WriteRune('{')
		sb.WriteRune('}')
	}

	return sb.String()
}

func (w *WhileExpression) expressionNode() {}

type BreakStatement struct{}

func (b *BreakStatement) DebugString() string { return b.String() }

func (b *BreakStatement) String() string { return "break;" }

func (b *BreakStatement) statementNode() {}

type ContinueStatement struct{}

func (c *ContinueStatement) DebugString() string { return c.String() }

func (c *ContinueStatement) String() string { return "continue;" }

func (c *ContinueStatement) statementNode() {}

type Null struct{}

var NULL = &Null{}

func (n *Null) DebugString() string { return n.String() }

func (n *Null) String() string { return "null" }

func (n *Null) expressionNode() {}
