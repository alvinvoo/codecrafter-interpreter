package parser

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

func TestExpression(t *testing.T) {
	b := NewBinary(
		NewLiteral("1"),
		scanner.NewToken(scanner.PLUS, "+", ""),
		NewLiteral("2"),
	)
	c := NewBinary(
		NewUnary(
			scanner.NewToken(scanner.MINUS, "-", "null"),
			NewLiteral(123),
		),
		scanner.NewToken(scanner.STAR, "*", "null"),
		NewGrouping(NewLiteral(45.67)),
	)

	tests := []struct {
		expr Expr
		want string
	}{
		{
			expr: c,
			want: "(* (- 123) (group 45.67))",
		},
		{
			expr: b,
			want: "(+ 1 2)",
		},
	}

	for _, test := range tests {
		if test.expr.Accept(NewAstPrinter()).(string) == "" {
			t.Errorf("Expression() = %q, want non-empty string", test.expr.Accept(NewAstPrinter()).(string))
		}

		if test.expr.Accept(NewAstPrinter()).(string) != test.want {
			t.Errorf("Expression() = %q, want %q", test.expr.Accept(NewAstPrinter()).(string), test.want)
		}
	}
}
