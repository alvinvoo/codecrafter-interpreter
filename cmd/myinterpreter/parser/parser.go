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
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
*/

import (
	"fmt"

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

// matchAny will also advance the token if it matches
func (p *Parser) matchAny(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(t scanner.TokenType, expectedTokenMsg string) (scanner.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return scanner.Token{}, fmt.Errorf("expect '%s' after expression", expectedTokenMsg)
}

func (p *Parser) expression() (Expr, error) {
	return p.unary()
}

func (p *Parser) unary() (Expr, error) {
	if p.matchAny(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return NewLiteral(nil), err
		}

		return NewUnary(operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.matchAny(scanner.FALSE) {
		return NewLiteral(false), nil
	}

	if p.matchAny(scanner.TRUE) {
		return NewLiteral(true), nil
	}

	if p.matchAny(scanner.NIL) {
		return NewLiteral(nil), nil
	}

	if p.matchAny(scanner.NUMBER, scanner.STRING) {
		return NewLiteral(p.previous().Literal), nil
	}

	if p.matchAny(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return NewLiteral(nil), err
		}

		_, err = p.consume(scanner.RIGHT_PAREN, ")")
		if err != nil {
			return NewLiteral(nil), err
		}

		return NewGrouping(expr), nil
	}

	return NewLiteral(nil), fmt.Errorf("no primary expression found")
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}
