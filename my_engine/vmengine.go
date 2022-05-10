package my_engine

import (
	"monkey/my_ast"
	"monkey/my_compiler"
	"monkey/my_lexer"
	"monkey/my_object"
	"monkey/my_parser"
	"monkey/my_vm"
)

type vmEngine struct{}

func NewVMEngine() Engine {
	return &vmEngine{}
}

func (vme *vmEngine) Evaluate(code string) (my_object.Object, error) {
	program, err := parse(code)
	if err != nil {
		return nil, err
	}
	comp := my_compiler.New()
	err = comp.Compile(program)
	if err != nil {
		return nil, err
	}
	virtualMachine := my_vm.New(comp.ByteCode())
	err = virtualMachine.Run()
	if err != nil {
		return nil, err
	}
	stackTop := virtualMachine.LastPoppedStackItem()
	return stackTop, nil
}

func parse(line string) (*my_ast.Program, error) {
	l := my_lexer.New(line)
	p := my_parser.New(l)

	program := p.Parse()
	if p.Error() != nil {
		return nil, p.Error()
	}
	return program, nil
}
