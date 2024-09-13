package lox

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap AstPrinter) parenthesize(name string, exprs ...Expr) string {
	str := "(" + name

	for _, expr := range exprs {
		str += " "
		str += expr.Accept(ap).(string)
	}

	str += ")"
	return str
}

func (ap AstPrinter) Print(expr Expr) string {
	return expr.Accept(ap).(string)
}

// literal        → NUMBER | STRING | "true" | "false" | "nil" ;
func (ap AstPrinter) visitLiteralExpr(l Expr) interface{} {
	var ret string
	if ll, ok := l.(Literal); ok {
		if ll.Value == nil {
			ret = "nil"
		} else {
			ret = fmt.Sprintf("%v", scanner.HandleNumberLiteral(ll.Value))
		}
	}

	return ret
}

func (ap AstPrinter) visitGroupingExpr(g Expr) interface{} {
	var ret string
	if gg, ok := g.(Grouping); ok {
		ret = ap.parenthesize("group", gg.Expression)
	}

	return ret
}

func (ap AstPrinter) visitUnaryExpr(u Expr) interface{} {
	var ret string
	if uu, ok := u.(Unary); ok {
		ret = ap.parenthesize(uu.Operator.Lexeme, uu.Right)
	}

	return ret
}

func (ap AstPrinter) visitBinaryExpr(b Expr) interface{} {
	var ret string
	if bb, ok := b.(Binary); ok {
		ret = ap.parenthesize(bb.Operator.Lexeme, bb.Left, bb.Right)
	}

	return ret
}

func (ap AstPrinter) visitPrintStmt(p Stmt) interface{} {
	if pp, ok := p.(Print); ok {
		return ap.parenthesize("print", pp.Expression)
	}
	return nil
}

func (ap AstPrinter) visitExpressionStmt(e Stmt) interface{} {
	if ee, ok := e.(Expression); ok {
		return ap.parenthesize(";", ee.Expression)
	}
	return nil
}
