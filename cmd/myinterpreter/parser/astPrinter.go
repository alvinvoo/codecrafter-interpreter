package parser

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

func (ap AstPrinter) visitLiteralExpr(l Expr) interface{} {
	var ret string
	if ll, ok := l.(Literal); ok {
		if ll.Value == nil {
			ret = "nil"
		}

		ret = ll.Value.(string)
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
		ret = ap.parenthesize(uu.Token.Lexeme, uu.Expression)
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
