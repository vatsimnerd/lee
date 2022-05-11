// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package lexer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Illegal-0]
	_ = x[EOF-1]
	_ = x[WhiteSpace-2]
	_ = x[Identifier-3]
	_ = x[Number-4]
	_ = x[String-5]
	_ = x[NotEquals-6]
	_ = x[Equals-7]
	_ = x[Matches-8]
	_ = x[NotMatches-9]
	_ = x[Less-10]
	_ = x[Greater-11]
	_ = x[LessOrEqual-12]
	_ = x[GreaterOrEqual-13]
	_ = x[LBrace-14]
	_ = x[RBrace-15]
	_ = x[Or-16]
	_ = x[And-17]
}

const _TokenType_name = "IllegalEOFWhiteSpaceIdentifierNumberStringNotEqualsEqualsMatchesNotMatchesLessGreaterLessOrEqualGreaterOrEqualLBraceRBraceOrAnd"

var _TokenType_index = [...]uint8{0, 7, 10, 20, 30, 36, 42, 51, 57, 64, 74, 78, 85, 96, 110, 116, 122, 124, 127}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}