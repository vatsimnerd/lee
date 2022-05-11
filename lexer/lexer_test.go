package lexer

import "testing"

type testcase struct {
	input  string
	output []Token
}

var (
	validCases = []testcase{
		{
			"a < 5",
			[]Token{
				{Identifier, "a", 1, 1},
				{Less, "<", 1, 3},
				{Number, "5", 1, 5},
				{EOF, "", 1, 6},
			},
		},
		{
			"a >= 7 and b != 12.5",
			[]Token{
				{Identifier, "a", 1, 1},
				{GreaterOrEqual, ">=", 1, 3},
				{Number, "7", 1, 6},
				{And, "and", 1, 8},
				{Identifier, "b", 1, 12},
				{NotEquals, "!=", 1, 14},
				{Number, "12.5", 1, 17},
				{EOF, "", 1, 21},
			},
		},
	}
)

func TestLexerValidCases(t *testing.T) {
	for i, tc := range validCases {
		tf, err := Tokenize(tc.input, true)
		if err != nil {
			t.Errorf("error testing case %d: %v", i+1, err)
			return
		}

		for i := 0; i < len(tc.output); i++ {
			expected := tc.output[i]
			actual := tf.get(i)

			if actual == nil {
				t.Errorf(
					"error in test %d: unexpected end of stream, expected token %s, got nil",
					i+1,
					expected.String(),
				)
				return
			}

			if expected.ne(*actual) {
				t.Errorf(
					"error in test %d: unexpected token %s, expected %s",
					i+1, actual.String(),
					expected.String(),
				)
			}
		}
	}
}
