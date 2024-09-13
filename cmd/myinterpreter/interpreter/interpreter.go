package interpreter

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/util"
)

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i Interpreter) Evaluate(expr parser.Expr) interface{} {
	return expr.Accept(i)
}

func (i Interpreter) VisitLiteralExpr(l parser.Expr) interface{} {
	if ll, ok := l.(parser.Literal); ok {
		return ll.Value
	}

	return nil
}

func (i Interpreter) VisitGroupingExpr(g parser.Expr) interface{} {
	if gg, ok := g.(parser.Grouping); ok {
		return i.Evaluate(gg.Expression)
	}

	return nil
}

// Ruby’s simple rule: false and nil are falsey, and everything else is truthy
func isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if b, ok := obj.(bool); ok {
		return b
	}

	return true
}

func checkNumberOperand(operator scanner.Token, operand interface{}) {
	if _, ok := operand.(float64); !ok {
		panic(
			util.NewRuntimeError(operator, "Operand must be a number."),
		)
	}
}

func (i Interpreter) VisitUnaryExpr(u parser.Expr) interface{} {
	if uu, ok := u.(parser.Unary); ok {
		right := i.Evaluate(uu.Right)

		switch uu.Operator.TokenType {
		case scanner.MINUS:
			checkNumberOperand(uu.Operator, right)
			return -(right.(float64))
		case scanner.BANG:
			return !isTruthy(right)
		}
	}

	return nil
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

func checkNumberOperands(operator scanner.Token, left, right interface{}) {
	if _, ok := left.(float64); !ok {
		panic(
			util.NewRuntimeError(operator, "Left operand must be a number."),
		)
	}

	if _, ok := right.(float64); !ok {
		panic(
			util.NewRuntimeError(operator, "Right operand must be a number."),
		)
	}
}

func (i Interpreter) VisitBinaryExpr(b parser.Expr) interface{} {
	if bb, ok := b.(parser.Binary); ok {
		left := i.Evaluate(bb.Left)
		right := i.Evaluate(bb.Right)

		switch bb.Operator.TokenType {
		case scanner.MINUS:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) - right.(float64)
		case scanner.STAR:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) * right.(float64)
		case scanner.SLASH:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) / right.(float64)
		case scanner.GREATER:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) > right.(float64)
		case scanner.GREATER_EQUAL:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) >= right.(float64)
		case scanner.LESS:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) < right.(float64)
		case scanner.LESS_EQUAL:
			checkNumberOperands(bb.Operator, left, right)
			return left.(float64) <= right.(float64)
		case scanner.EQUAL_EQUAL:
			return isEqual(left, right)
		case scanner.BANG_EQUAL:
			return !isEqual(left, right)
		case scanner.PLUS:
			if l, ok := left.(float64); ok {
				if r, ok := right.(float64); ok {
					return l + r
				}
			}

			if l, ok := left.(string); ok {
				if r, ok := right.(string); ok {
					return l + r
				}
			}

			panic(
				util.NewRuntimeError(bb.Operator, "Operands must be two numbers or two strings."),
			)
		}
	}

	return nil
}
