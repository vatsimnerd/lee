package parser

import (
	"fmt"
	"strconv"

	"github.com/vatsimnerd/lee/lexer"
)

type parser[T any] struct {
	tokens *lexer.TokenFlow
}

func newParser[T any](tokens *lexer.TokenFlow) *parser[T] {
	return &parser[T]{tokens}
}

func unexpected(token *lexer.Token) error {
	return fmt.Errorf(
		"unexpected token %s at line %d pos %d",
		token.Literal,
		token.Line,
		token.Position,
	)
}

func unexpectedTokenType(token *lexer.Token, expected lexer.TokenType) error {
	return fmt.Errorf(
		"unexpected token %s at line %d pos %d, expected %s",
		token.Literal,
		token.Line,
		token.Position,
		expected.String(),
	)
}

func (p *parser[T]) eat(tokenType lexer.TokenType) error {
	t := p.tokens.Current()
	if t.Type != tokenType {
		return unexpectedTokenType(t, tokenType)
	}
	p.tokens.Advance()
	return nil
}

func (p *parser[T]) parseCombineOperator() (*CombineOperator, error) {
	t := p.tokens.Current()
	if opType, found := combOperators[t.Type]; found {
		p.tokens.Advance()
		return &CombineOperator{opType, t}, nil
	}
	return nil, unexpected(t)
}

func (p *parser[T]) parseGrouping() (*Grouping[T], error) {
	var err error

	err = p.eat(lexer.LBrace)
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	err = p.eat(lexer.RBrace)
	if err != nil {
		return nil, err
	}

	return &Grouping[T]{expr}, nil
}

func (p *parser[T]) parseExpression() (*Expression[T], error) {
	var err error

	expr := &Expression[T]{
		Left: &LeftExpression[T]{},
	}

	t := p.tokens.Current()
	if t.Type == lexer.LBrace {
		expr.Left.Grouping, err = p.parseGrouping()
		if err != nil {
			return nil, err
		}
	} else if t.Type == lexer.Identifier {
		expr.Left.Condition, err = p.parseCondition()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, unexpected(t)
	}

	t = p.tokens.Current()
	if t.Type == lexer.EOF || t.Type == lexer.RBrace {
		return expr, nil
	}

	if t.Type != lexer.And && t.Type != lexer.Or {
		return nil, unexpected(t)
	}

	expr.Operator, err = p.parseCombineOperator()
	if err != nil {
		return nil, err
	}

	expr.Right, err = p.parseExpression()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *parser[T]) parseCondition() (*Condition[T], error) {
	cond := &Condition[T]{}

	t := p.tokens.Current()
	if t.Type != lexer.Identifier {
		return nil, unexpectedTokenType(t, lexer.Identifier)
	}
	cond.Identifier = &Identifier{t.Literal, t}
	p.tokens.Advance()

	t = p.tokens.Current()
	if opType, found := operators[t.Type]; found {
		cond.Operator = &Operator{opType, t}
	} else {
		return nil, unexpected(t)
	}
	p.tokens.Advance()

	t = p.tokens.Current()
	if t.Type == lexer.String {
		cond.Value = &Value{String: &t.Literal, Token: t}
	} else if t.Type == lexer.Number {
		value, err := strconv.ParseFloat(t.Literal, 64)
		if err != nil {
			return nil, err
		}
		cond.Value = &Value{Number: &value, Token: t}
	} else {
		return nil, unexpected(t)
	}
	p.tokens.Advance()
	return cond, nil
}

func Parse[T any](tokens *lexer.TokenFlow) (*Expression[T], error) {
	p := newParser[T](tokens)
	return p.parseExpression()
}
