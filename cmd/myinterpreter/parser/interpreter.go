package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i Interpreter) Evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i Interpreter) visitLiteralExpr(l Expr) interface{} {
	if ll, ok := l.(Literal); ok {
		return ll.Value
	}

	return nil
}

func (i Interpreter) visitGroupingExpr(g Expr) interface{} {
	if gg, ok := g.(Grouping); ok {
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

func (i Interpreter) visitUnaryExpr(u Expr) interface{} {
	if uu, ok := u.(Unary); ok {
		right := i.Evaluate(uu.Right)

		switch uu.Operator.TokenType {
		case scanner.MINUS:
			return -(right.(float64))
		case scanner.BANG:
			return !isTruthy(right)
		}
	}

	return nil
}

func (i Interpreter) visitBinaryExpr(b Expr) interface{} {
	if bb, ok := b.(Binary); ok {
		left := i.Evaluate(bb.Left)
		right := i.Evaluate(bb.Right)

		switch bb.Operator.TokenType {
		case scanner.STAR:
			return left.(float64) * right.(float64)
		case scanner.SLASH:
			return left.(float64) / right.(float64)
		case scanner.GREATER:
			return left.(float64) > right.(float64)
		case scanner.GREATER_EQUAL:
			return left.(float64) >= right.(float64)
		case scanner.LESS:
			return left.(float64) < right.(float64)
		case scanner.LESS_EQUAL:
			return left.(float64) <= right.(float64)
		case scanner.EQUAL_EQUAL:
			return left == right
		case scanner.BANG_EQUAL:
			return left != right
		case scanner.MINUS:
			return left.(float64) - right.(float64)
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
		}
	}

	return nil
}
