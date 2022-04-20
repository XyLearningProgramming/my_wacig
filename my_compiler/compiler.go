package my_compiler

import (
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
	case *my_ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		// TODO: compile infix op?
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
	case *my_ast.Integer:
		c.emit(
			my_code.OpConstant,
			c.addConstant(&my_object.Integer{Value: int64(node.Value)}),
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
