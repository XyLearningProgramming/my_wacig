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
	// trackedInstructions: tracking the last and before the last instructions emitted
	trackedInstructions [2]*EmittedInstruction
	symbolTable         *SymbolTable
}

type ByteCode struct {
	Instructions my_code.Instructions
	Constants    []my_object.Object
}

type EmittedInstruction struct {
	Position int
	OpCode   my_code.Opcode
}

func New() *Compiler {
	return &Compiler{
		instructions:        my_code.Instructions{},
		constants:           []my_object.Object{},
		trackedInstructions: [2]*EmittedInstruction{nil, nil},
		symbolTable:         NewSymbolTable(),
	}
}

func NewWithState(constants []my_object.Object, symbolTable *SymbolTable) *Compiler {
	return &Compiler{
		instructions:        my_code.Instructions{},
		constants:           constants,
		trackedInstructions: [2]*EmittedInstruction{nil, nil},
		symbolTable:         symbolTable,
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
	case *my_ast.BlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *my_ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		// decide what number to set to the identifier from the symbol table
		sym := c.symbolTable.Define(node.Ident.Value)
		c.emit(my_code.OpSetGlobal, sym.Index)
	// expressions
	case *my_ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		jumpNotTruthyPos := c.emit(my_code.OpJumpNotTruthy, 0)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		// remove last OpPop because block statement emitted one and if is an expression,
		// which means if should return one and only one expression value
		if c.isLastInstruction(my_code.OpPop) {
			c.removeLastInstruction()
		}
		jumpPos := c.emit(my_code.OpJump, 0)
		// change OpJumpNotTruthy operands after we knew where consequence ins ends
		c.replaceOperands(jumpNotTruthyPos, len(c.instructions))
		if node.Alternative == nil {
			c.emit(my_code.OpNull)
		} else {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.isLastInstruction(my_code.OpPop) {
				c.removeLastInstruction()
			}
		}
		// change OpJump operands after we knew where alternative ins ends
		c.replaceOperands(jumpPos, len(c.instructions))
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
	case *my_ast.Identifier:
		sym, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.String())
		}
		c.emit(my_code.OpGetGlobal, sym.Index)
	case *my_ast.Null:
		c.emit(my_code.OpNull)
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
	c.setLastInstruction(posNewIns, op)
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

func (c *Compiler) setLastInstruction(position int, op my_code.Opcode) {
	c.trackedInstructions[0] = c.trackedInstructions[1]
	c.trackedInstructions[1] = &EmittedInstruction{Position: position, OpCode: op}
}

func (c *Compiler) isLastInstruction(op my_code.Opcode) bool {
	if c.trackedInstructions[1] == nil {
		return false
	}
	return c.trackedInstructions[1].OpCode == op
}

func (c *Compiler) removeLastInstruction() {
	// should panic if illegal removing happens in compiler
	c.instructions = c.instructions[:c.trackedInstructions[1].Position]
	c.trackedInstructions[1] = c.trackedInstructions[0]
	c.trackedInstructions[0] = nil
}

func (c *Compiler) replaceOperands(pos int, operands ...int) {
	op := my_code.Opcode(c.instructions[pos])
	c.replaceInstruction(pos, op, operands...)
}

func (c *Compiler) replaceInstruction(pos int, op my_code.Opcode, operands ...int) {
	ins := my_code.Make(op, operands...)
	for i := 0; i < len(ins); i++ {
		c.instructions[pos+i] = ins[i]
	}
}
