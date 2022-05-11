package parser

import (
	"github.com/vatsimnerd/lee/lexer"
)

//go:generate stringer -type=OperatorType,CombineOperatorType

type (
	CombineOperatorType int
	OperatorType        int
)

const (
	And CombineOperatorType = iota
	Or

	Equals OperatorType = iota
	NotEquals
	Matches
	NotMatches
	Less
	Greater
	LessOrEqual
	GreaterOrEqual
)

var (
	combOperators = map[lexer.TokenType]CombineOperatorType{
		lexer.And: And,
		lexer.Or:  Or,
	}

	operators = map[lexer.TokenType]OperatorType{
		lexer.Equals:         Equals,
		lexer.NotEquals:      NotEquals,
		lexer.Matches:        Matches,
		lexer.NotMatches:     NotMatches,
		lexer.Less:           Less,
		lexer.LessOrEqual:    LessOrEqual,
		lexer.Greater:        Greater,
		lexer.GreaterOrEqual: GreaterOrEqual,
	}
)
