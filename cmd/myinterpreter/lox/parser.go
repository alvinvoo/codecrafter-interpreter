package lox

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

// new rules for statements
program        → statement* EOF ;

statement      → exprStmt
               | printStmt ;

exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
*/

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/util"
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
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return NewLiteral(nil), err
	}

	for p.matchAny(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return NewLiteral(nil), err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return NewLiteral(nil), err
	}

	for p.matchAny(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return NewLiteral(nil), err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return NewLiteral(nil), err
	}

	for p.matchAny(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return NewLiteral(nil), err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return NewLiteral(nil), err
	}

	for p.matchAny(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return NewLiteral(nil), err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
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

	return NewLiteral(nil), fmt.Errorf(util.Error(p.peek(), "Expect expression"))
}

func (p *Parser) statement() (Stmt, error) {
	if p.matchAny(scanner.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.SEMICOLON, "';'")
	if err != nil {
		return nil, err
	}

	return NewPrint(value), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.SEMICOLON, "';'")
	if err != nil {
		return nil, err
	}

	return NewExpression(value), nil
}

func (p *Parser) ParseExpr() (Expr, error) {
	return p.expression()
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}
