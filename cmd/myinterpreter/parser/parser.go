package parser

/*
Precedence is from lowest to highest

Name        Operators    Associates
Equality    == !=        Left
Comparison  > >= < <=    Left
Term        - +          Left
Factor      / *          Left
Unary       ! -          Right

Grammar:
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ; // * means zero or more
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")"
*/

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == t
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current += 1
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) matchAny(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) expression() Expr {
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.matchAny(scanner.FALSE) {
		return NewLiteral(false)
	}

	if p.matchAny(scanner.TRUE) {
		return NewLiteral(true)
	}

	if p.matchAny(scanner.NIL) {
		return NewLiteral(nil)
	}

	if p.matchAny(scanner.NUMBER, scanner.STRING) {
		return NewLiteral(p.previous().Literal)
	}

	return NewLiteral(nil)
}

func (p *Parser) Parse() Expr {
	return p.expression()
}
