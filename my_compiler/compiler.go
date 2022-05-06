package my_compiler

import (
	"fmt"
	"monkey/my_ast"
	"monkey/my_code"
	"monkey/my_object"
)

type Compiler struct {
	instructions my_code.Instructions
	constants    []my_object.Object
}

type ByteCode struct {
	Instructions my_code.Instructions
	Constants    []my_object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: my_code.Instructions{},
		constants:    []my_object.Object{},
	}
}

func (c *Compiler) Compile(node my_ast.Node) error {
	switch node := node.(type) {
	// statements
	case *my_ast.Program:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *my_ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(my_code.OpPop)
	// expressions
	case *my_ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "-":
			c.emit(my_code.OpMinus)
		case "!":
			c.emit(my_code.OpBang)
		default:
			return fmt.Errorf("unknown prefix operator: %s", node.Operator)
		}
	case *my_ast.InfixExpression:
		var first my_ast.Expression
		var second my_ast.Expression
		if isNodeReversed(node.Operator) {
			first = node.Right
			second = node.Left
		} else {
			first = node.Left
			second = node.Right
		}
		err := c.Compile(first)
		if err != nil {
			return err
		}
		err = c.Compile(second)
		if err != nil {
			return err
		}
		// TODO: compile infix op?
		switch node.Operator {
		case "+":
			c.emit(my_code.OpAdd)
		case "-":
			c.emit(my_code.OpSub)
		case "*":
			c.emit(my_code.OpMul)
		case "/":
			c.emit(my_code.OpDiv)
		case "<":
			fallthrough
		case ">":
			c.emit(my_code.OpGT)
		case "<=":
			fallthrough
		case ">=":
			c.emit(my_code.OpGTE)
		case "==":
			c.emit(my_code.OpEqual)
		case "!=":
			c.emit(my_code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *my_ast.Integer:
		c.emit(
			my_code.OpConstant,
			c.addConstant(&my_object.Integer{Value: int64(node.Value)}),
		)
	case *my_ast.Float:
		c.emit(
			my_code.OpConstant,
			c.addConstant(&my_object.Float{Value: node.Value}),
		)
	case *my_ast.Boolean:
		if node.Value {
			c.emit(my_code.OpTrue)
		} else {
			c.emit(my_code.OpFalse)
		}
	case *my_ast.StringExpression:
		c.emit(
			my_code.OpConstant,
			c.addConstant(&my_object.String{Value: node.Value}),
		)
	}
	return nil
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{Instructions: c.instructions, Constants: c.constants}
}

// addConstant: add to the constant pool and return index of the newly added item as identifier
func (c *Compiler) addConstant(obj my_object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op my_code.Opcode, operands ...int) (posNewIns int) {
	posNewIns = len(c.instructions)
	ins := my_code.Make(op, operands...)
	c.instructions = append(c.instructions, ins...)
	return posNewIns
}

func isNodeReversed(operator my_ast.InfixOperator) bool {
	switch operator {
	case "<=":
		fallthrough
	case "<":
		return true
	default:
		return false
	}
}
