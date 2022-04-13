package my_evaluator

import (
	"monkey/my_object"
)

var (
	TRUE             = &my_object.Boolean{Value: true}
	FALSE            = &my_object.Boolean{Value: false}
	TRUE_AS_ONE      = &my_object.Integer{Value: 1}
	FALSE_AS_ZERO    = &my_object.Integer{Value: 0}
	TRUE_AS_ONE_FL   = &my_object.Float{Value: 1}
	FALSE_AS_ZERO_FL = &my_object.Float{Value: 0}
	NULL             = &my_object.Null{}
	BREAK_ERROR = &my_object.Error{Message: "break outside loop"}
	CONTINUE_ERROR =  &my_object.Error{Message: "continue outside loop"}
	EMPTY_ARRAY = &my_object.Array{Elements:  []my_object.Object{}}
)
