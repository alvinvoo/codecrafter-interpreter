package lox

type Visitor interface {
	visitLiteralExpr(literal Expr) interface{}
	visitGroupingExpr(grouping Expr) interface{}
	visitUnaryExpr(unary Expr) interface{}
	visitBinaryExpr(binary Expr) interface{}

	visitExpressionStmt(expression Stmt) interface{}
	visitPrintStmt(print Stmt) interface{}
}
