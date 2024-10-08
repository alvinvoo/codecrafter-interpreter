package lox

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
	return v.visitLiteralExpr(l)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Grouping {
	return Grouping{Expression: expression}
}

func (g Grouping) Accept(v Visitor) interface{} {
	return v.visitGroupingExpr(g)
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func NewUnary(operator scanner.Token, right Expr) Unary {
	return Unary{Operator: operator, Right: right}
}

func (u Unary) Accept(v Visitor) interface{} {
	return v.visitUnaryExpr(u)
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
	return v.visitBinaryExpr(b)
}
