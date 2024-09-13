package lox

type Stmt interface {
	Accept(v Visitor) interface{}
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Expression {
	return Expression{Expression: expression}
}

func (e Expression) Accept(v Visitor) interface{} {
	return v.visitExpressionStmt(e)
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) Print {
	return Print{Expression: expression}
}

func (p Print) Accept(v Visitor) interface{} {
	return v.visitPrintStmt(p)
}
