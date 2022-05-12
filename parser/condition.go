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

func (v Value) IsString() bool {
	return v.String != nil
}

func (v Value) IsFloat() bool {
	return v.Number != nil
}

func (v Value) GetStringValue() (string, error) {
	if !v.IsString() {
		return "", fmt.Errorf("token %s has no string value", v.Token.String())
	}
	return *v.String, nil
}

func (v Value) GetFloatValue() (float64, error) {
	if !v.IsFloat() {
		return 0, fmt.Errorf("token %s has no string value", v.Token.String())
	}
	return *v.Number, nil
}

func (v Value) MustGetStringValue() string {
	str, err := v.GetStringValue()
	if err != nil {
		panic(err)
	}
	return str
}

func (v Value) MustGetFloatValue() float64 {
	f, err := v.GetFloatValue()
	if err != nil {
		panic(err)
	}
	return f
}
