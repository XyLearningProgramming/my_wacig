package my_engine

import (
	"monkey/my_object"
)

type Engine interface {
	Evaluate(code string) (result my_object.Object, err error)
}
