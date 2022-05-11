package lexer

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=TokenType

type TokenType int

const (
	Illegal TokenType = iota
	EOF
	WhiteSpace

	Identifier
	Number
	String

	NotEquals
	Equals
	Matches
	NotMatches
	Less
	Greater
	LessOrEqual
	GreaterOrEqual

	LBrace
	RBrace

	Or
	And
)

type (
	Token struct {
		Type     TokenType
		Literal  string
		Line     int
		Position int
	}

	TokenFlow struct {
		tokens []Token
		idx    int
	}
)

func (t Token) String() string {
	literal := strings.ReplaceAll(t.Literal, "\n", "\\n")
	return fmt.Sprintf(
		"<Token type=%s literal=\"%s\" at=%d:%d>",
		t.Type.String(),
		literal,
		t.Line,
		t.Position,
	)
}

func (t Token) ne(o Token) bool {
	return t.Type != o.Type ||
		t.Literal != o.Literal ||
		t.Line != o.Line ||
		t.Position != o.Position
}

func newTokenFlow(tokens []Token, skipWhitespace bool) *TokenFlow {
	tf := TokenFlow{make([]Token, 0), 0}
	for _, t := range tokens {
		if t.Type == WhiteSpace && skipWhitespace {
			continue
		}
		tf.tokens = append(tf.tokens, t)
	}
	return &tf
}

// Reset resets the current index to zero
func (tf *TokenFlow) Reset() {
	tf.idx = 0
}

func (tf *TokenFlow) get(idx int) *Token {
	if idx >= len(tf.tokens) {
		return nil
	}
	return &tf.tokens[idx]
}

// Current returns the current token
func (tf *TokenFlow) Current() *Token {
	return tf.get(tf.idx)
}

// Next returns the following token, useful for lookahead
func (tf *TokenFlow) Next() *Token {
	return tf.get(tf.idx + 1)
}

// Advance sets the internal index to the next token
func (tf *TokenFlow) Advance() {
	tf.idx++
}
