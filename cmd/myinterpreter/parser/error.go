package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

func report(line int, where string, message string) string {
	return fmt.Sprintf("[line %d] Error %s: %s\n", line, where, message)
}

func Error(token scanner.Token, message string) string {
	if token.TokenType == scanner.EOF {
		return report(token.Line, "at end", message)
	} else {
		return report(token.Line, "at '"+token.Lexeme+"'", message)
	}
}
