package parser

import (
	"testing"

	"github.com/vatsimnerd/lee/lexer"
)

func getParser[T any](data string) *parser[T] {
	l, err := lexer.Tokenize(data, true)
	if err != nil {
		return nil
	}
	return newParser[T](l)
}

func TestParseCondition1(t *testing.T) {
	var p *parser[string]
	var c *Condition[string]
	var err error

	p = getParser[string](`ident =~ "value"`)
	c, err = p.parseCondition()
	if err != nil {
		t.Errorf("unexpected error parsing condition: %v", err)
		return
	}

	if c.Identifier.Name != "ident" {
		t.Errorf("invalid identifier name, got %s, expected %s", c.Identifier.Name, "ident")
	}
	if c.Operator.Type != Matches {
		t.Errorf("invalid operator type, got %s, expected %s", c.Operator.Type, Matches)
	}
	if c.Value.String == nil {
		t.Errorf("string value is unexpectedly nil")
	}
	if *c.Value.String != `"value"` {
		t.Errorf("invalid value, got %v, expected %v", *c.Value.String, `"value"`)
	}
}

func TestParseCondition2(t *testing.T) {
	var p *parser[string]
	var c *Condition[string]
	var err error

	p = getParser[string](`num > 45.3`)
	c, err = p.parseCondition()
	if err != nil {
		t.Errorf("unexpected error parsing condition: %v", err)
		return
	}

	if c.Identifier.Name != "num" {
		t.Errorf("invalid identifier name, got %s, expected %s", c.Identifier.Name, "num")
	}
	if c.Operator.Type != Greater {
		t.Errorf("invalid operator type, got %s, expected %s", c.Operator.Type, Greater)
	}
	if c.Value.Number == nil {
		t.Errorf("number value is unexpectedly nil")
	}
	if *c.Value.Number != 45.3 {
		t.Errorf("invalid value, got %v, expected %v", *c.Value.Number, 45.3)
	}
}

func TestParseCondition3(t *testing.T) {
	var p *parser[string]
	var err error

	p = getParser[string](`should = fail`)
	_, err = p.parseCondition()
	if err == nil {
		t.Errorf("should throw unexpected")
		return
	}

	exp := "unexpected token fail at line 1 pos 10"
	if err.Error() != exp {
		t.Errorf("should throw %v, but throws %v", exp, err)
		return
	}
}

func TestParseCondition4(t *testing.T) {
	var p *parser[string]
	var err error

	p = getParser[string](`invalid op`)
	_, err = p.parseCondition()
	if err == nil {
		t.Errorf("should throw unexpected")
		return
	}

	exp := "unexpected token op at line 1 pos 9"
	if err.Error() != exp {
		t.Errorf("should throw %v, but throws %v", exp, err)
		return
	}
}
