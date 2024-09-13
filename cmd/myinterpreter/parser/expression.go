package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

/** grammar rules
expression     → literal
               | unary
               | binary
               | grouping ;

literal        → NUMBER | STRING | "true" | "false" | "nil" ;
grouping       → "(" expression ")" ;
unary          → ( "-" | "!" ) expression ;
binary         → expression operator expression ;
operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
               | "+"  | "-"  | "*" | "/" ;


** This still grammar has ambiguity
**/

type Visitor interface {
	VisitLiteralExpr(literal Expr) interface{}
	VisitGroupingExpr(grouping Expr) interface{}
	VisitUnaryExpr(unary Expr) interface{}
	VisitBinaryExpr(binary Expr) interface{}
}

type Expr interface {
	Accept(v Visitor) interface{}
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) Literal {
	return Literal{Value: value}
}

func (l Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(l)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Grouping {
	return Grouping{Expression: expression}
}

func (g Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(g)
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func NewUnary(operator scanner.Token, right Expr) Unary {
	return Unary{Operator: operator, Right: right}
}

func (u Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(u)
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func NewBinary(left Expr, operator scanner.Token, right Expr) Binary {
	return Binary{Left: left, Operator: operator, Right: right}
}

func (b Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(b)
}
