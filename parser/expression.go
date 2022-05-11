package parser

import "github.com/vatsimnerd/lee/lexer"

type (
	Grouping[T any] struct {
		Expression *Expression[T]
	}

	LeftExpression[T any] struct {
		Condition *Condition[T]
		Grouping  *Grouping[T]
	}

	Expression[T any] struct {
		Operator *CombineOperator
		Left     *LeftExpression[T]
		Right    *Expression[T]
	}

	CombineOperator struct {
		Type  CombineOperatorType
		Token *lexer.Token
	}

	Matcher[T any]             func(model T) bool
	CompilationCallback[T any] func(c *Condition[T]) (Matcher[T], error)
)

func (g *Grouping[T]) String() string {
	return "(" + g.Expression.String() + ")"
}

func (e *Expression[T]) String() string {
	str := "Expr[ " + e.Left.String()
	if e.Operator != nil {
		str += " " + e.Operator.String() + " " + e.Right.String()
	}
	str += " ]"
	return str
}

func (le *LeftExpression[T]) String() string {
	if le.Condition != nil {
		return le.Condition.String()
	} else {
		return le.Grouping.String()
	}
}

func (co *CombineOperator) String() string {
	return co.Type.String()
}

func (e *Expression[T]) Compile(cb CompilationCallback[T]) error {
	var err error

	if e.Left.Condition != nil {
		e.Left.Condition.MatcherFunc, err = cb(e.Left.Condition)
		if err != nil {
			return err
		}
	} else {
		err = e.Left.Grouping.Expression.Compile(cb)
		if err != nil {
			return err
		}
	}

	if e.Right != nil {
		return e.Right.Compile(cb)
	}

	return nil
}

func (e *Expression[T]) Evaluate(model T) bool {
	var left bool

	if e.Left.Condition != nil {
		left = e.Left.Condition.MatcherFunc(model)
	} else {
		left = e.Left.Grouping.Expression.Evaluate(model)
	}

	if e.Right == nil {
		return left
	}

	switch e.Operator.Type {
	case And:
		if !left {
			return false
		}
	case Or:
		if left {
			return true
		}
	}

	return e.Right.Evaluate(model)
}
