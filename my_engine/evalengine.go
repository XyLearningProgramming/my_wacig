package my_engine

import (
	"monkey/my_evaluator"
	"monkey/my_object"
)

type evalEngine struct {
	env *my_object.Environment
}

func NewEvalEngine() Engine {
	return &evalEngine{env: my_object.NewEnvironment()}
}

func (e *evalEngine) Evaluate(code string) (result my_object.Object, err error) {
	program, err := parse(code)
	if err != nil {
		return nil, err
	}
	evaluated := my_evaluator.Eval(program, e.env)
	return evaluated, nil
}
