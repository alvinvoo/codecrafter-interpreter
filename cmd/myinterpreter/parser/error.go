package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

// handling static errors
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

// handling runtime errors
type RuntimeError struct {
	token   scanner.Token
	message string
}

func NewRuntimeError(token scanner.Token, message string) RuntimeError {
	return RuntimeError{token: token, message: message}
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", r.message, r.token.Line)
}
