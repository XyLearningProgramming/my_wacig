package my_evaluator

import (
	"monkey/my_ast"
	"monkey/my_object"
)

func evalForloop(node *my_ast.ForExpression, env *my_object.Environment) my_object.Object {
	// confine newly declared vars in new env
	enclosed := my_object.NewEnclosedEnvironment(env)
	if node.InitStmt != nil {
		initObj := Eval(node.InitStmt, enclosed)
		if isError(initObj) {
			return initObj
		}
	}
	var bodyObj my_object.Object = NULL
	for {
		if node.TestExpr != nil {
			testObj := Eval(node.TestExpr, enclosed)
			if isError(testObj) {
				return testObj
			}
			if !isTruthy(testObj) {
				break
			}
		}
		if node.Body != nil {
			bodyObjEvaluated := Eval(node.Body, enclosed)
			if isBreakError(bodyObjEvaluated) {
				break
			}
			if isContinueError(bodyObjEvaluated) {
				continue
			}
			if isError(bodyObjEvaluated) {
				return bodyObjEvaluated
			}
			bodyObj = bodyObjEvaluated
		}
		if node.UpdateStmt != nil {
			updateObj := Eval(node.UpdateStmt, enclosed)
			if isError(updateObj) {
				return updateObj
			}
		}
	}
	return bodyObj
}

func evalWhileLoop(node *my_ast.WhileExpression, env *my_object.Environment) my_object.Object {
	enclosed := my_object.NewEnclosedEnvironment(env)
	var bodyObj my_object.Object = NULL
	for {
		if node.TestExpr != nil {
			testObj := Eval(node.TestExpr, enclosed)
			if isError(testObj) {
				return testObj
			}
			if !isTruthy(testObj) {
				break
			}
		}
		if node.Body != nil {
			bodyObjEvaluated := Eval(node.Body, enclosed)
			if isBreakError(bodyObjEvaluated) {
				break
			}
			if isContinueError(bodyObjEvaluated) {
				continue
			}
			if isError(bodyObjEvaluated) {
				return bodyObjEvaluated
			}
			bodyObj = bodyObjEvaluated
		}

	}
	return bodyObj
}

func evalDoWhileLoop(node *my_ast.DoWhileExpression, env *my_object.Environment) my_object.Object {
	enclosed := my_object.NewEnclosedEnvironment(env)
	var bodyObj my_object.Object = NULL
	for {
		if node.Body != nil {
			bodyObjEvaluated := Eval(node.Body, enclosed)
			if isBreakError(bodyObjEvaluated) {
				break
			}
			if isContinueError(bodyObjEvaluated) {
				continue
			}
			if isError(bodyObjEvaluated) {
				return bodyObjEvaluated
			}
			bodyObj = bodyObjEvaluated
		}
		if node.TestExpr != nil {
			testObj := Eval(node.TestExpr, enclosed)
			if isError(testObj) {
				return testObj
			}
			if !isTruthy(testObj) {
				break
			}
		}
	}
	return bodyObj
}
