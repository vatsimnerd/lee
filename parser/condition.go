package parser

import (
	"fmt"

	"github.com/vatsimnerd/lee/lexer"
)

type (
	Operator struct {
		Type  OperatorType
		Token *lexer.Token
	}

	Identifier struct {
		Name  string
		Token *lexer.Token
	}

	Value struct {
		String *string
		Number *float64
		Token  *lexer.Token
	}

	Condition[T any] struct {
		Identifier  *Identifier
		Operator    *Operator
		Value       *Value
		MatcherFunc Matcher[T]
	}
)

func (c Condition[T]) String() string {
	return fmt.Sprintf("C{%s %s %s}", c.Identifier.Name, c.Operator.Token.Literal, c.Value.Token.Literal)
}
