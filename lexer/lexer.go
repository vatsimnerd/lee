package lexer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type lexer struct {
	sc      io.RuneScanner
	line    int
	pos     int
	literal string
	tokens  []Token
}

var (
	exprWhiteSpace = regexp.MustCompile(`\s`)
	exprIdentStart = regexp.MustCompile(`[a-zA-Z_]`)
	exprIdent      = regexp.MustCompile(`[a-z0-9A-Z_]`)
)

// push current literal as a token
func (l *lexer) push(t TokenType, line int, pos int) {
	token := Token{
		Type:     t,
		Literal:  l.literal,
		Line:     line,
		Position: pos + 1,
	}
	l.tokens = append(l.tokens, token)
	l.literal = ""
}

// add a rune to the current literal and maintain line/pos counters
func (l *lexer) eat(r rune) {
	l.literal += string(r)
	if r == '\n' {
		l.line++
		l.pos = 0
	} else {
		l.pos++
	}
}

// go one rune back
func (l *lexer) rewind() {
	_ = l.sc.UnreadRune()
}

// skip one rune
func (l *lexer) advance() {
	_, _, _ = l.sc.ReadRune()
}

func (l *lexer) readNumber() error {
	pos := l.pos
	line := l.line
	dotFound := false
	for {
		r, _, err := l.sc.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if r >= '0' && r <= '9' {
			l.eat(r)
		} else if r == '.' {
			if !dotFound {
				dotFound = true
				l.eat(r)
			} else {
				l.rewind()
				break
			}
		} else {
			l.rewind()
			break
		}
	}
	l.push(Number, line, pos)
	return nil
}

func (l *lexer) readWhitespace() error {
	pos := l.pos
	line := l.line
	for {
		r, _, err := l.sc.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if !exprWhiteSpace.MatchString(string(r)) {
			l.rewind()
			break
		}
		l.eat(r)
	}
	l.push(WhiteSpace, line, pos)
	return nil
}

func (l *lexer) readIdentifier() error {
	pos := l.pos
	line := l.line
	for {
		r, _, err := l.sc.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if !exprIdent.MatchString(string(r)) {
			l.rewind()
			break
		}
		l.eat(r)
	}

	operator := strings.ToLower(l.literal)
	if operator == "or" {
		l.push(Or, line, pos)
	} else if operator == "and" {
		l.push(And, line, pos)
	} else {
		l.push(Identifier, line, pos)
	}
	return nil
}

func (l *lexer) readAnd() error {
	line := l.line
	pos := l.pos

	// read &
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	r, _, err = l.sc.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.push(Illegal, line, pos)
			return nil
		}
		return err
	}

	if r != '&' {
		l.rewind()
		l.push(Illegal, line, pos)
	}

	l.eat(r)
	l.push(And, line, pos)
	return nil
}

func (l *lexer) readOr() error {
	line := l.line
	pos := l.pos

	// read |
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	r, _, err = l.sc.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.push(Illegal, line, pos)
			return nil
		}
		return err
	}

	if r != '|' {
		l.rewind()
		l.push(Illegal, line, pos)
	}

	l.eat(r)
	l.push(Or, line, pos)
	return nil
}

func (l *lexer) readEqualsOrMatches() error {
	line := l.line
	pos := l.pos

	// read the equals symbol
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	r, _, err = l.sc.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.push(Equals, line, pos)
			return nil
		}
		return err
	}

	if r == '~' {
		l.eat(r)
		l.push(Matches, line, pos)
	} else {
		l.rewind()
		l.push(Equals, line, pos)
	}
	return nil
}

func (l *lexer) readNotEqualsOrNotMatches() error {
	line := l.line
	pos := l.pos

	// read the exclamation mark
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	r, _, err = l.sc.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.push(Illegal, line, pos)
			return nil
		}
		return err
	}

	if r == '=' {
		l.eat(r)
		l.push(NotEquals, line, pos)
	} else if r == '~' {
		l.eat(r)
		l.push(NotMatches, line, pos)
	} else {
		l.rewind()
		l.push(Illegal, line, pos)
	}
	return nil
}

func (l *lexer) readGreater() error {
	line := l.line
	pos := l.pos

	// read the > symbol
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	r, _, err = l.sc.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.push(Greater, line, pos)
			return nil
		}
		return err
	}

	if r == '=' {
		l.eat(r)
		l.push(GreaterOrEqual, line, pos)
	} else {
		l.rewind()
		l.push(Greater, line, pos)
	}
	return nil
}

func (l *lexer) readLess() error {
	line := l.line
	pos := l.pos

	// read the < symbol
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	r, _, err = l.sc.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.push(Less, line, pos)
			return nil
		}
		return err
	}

	if r == '=' {
		l.eat(r)
		l.push(LessOrEqual, line, pos)
	} else {
		l.rewind()
		l.push(Less, line, pos)
	}
	return nil
}

func (l *lexer) readStringLiteral() error {
	line := l.line
	pos := l.pos

	// read opening quote
	r, _, err := l.sc.ReadRune()
	if err != nil {
		return err
	}
	l.eat(r)

	quoteSym := r // ' or "
	skipQuote := false

	for {
		r, _, err = l.sc.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.push(Illegal, line, pos)
				return fmt.Errorf("unexpected end of file while reading a string")
			}
			return err
		}

		l.eat(r)
		if r == quoteSym && !skipQuote {
			// found closing quote, end of string literal
			break
		}

		// eat quote on next step as it's escaped
		skipQuote = r == '\\'
	}
	l.push(String, line, pos)
	return nil
}

func (l *lexer) readAll() error {
	var r rune
	var err error

	for {
		r, _, err = l.sc.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.push(EOF, l.line, l.pos)
				break
			}
			return err
		}
		l.rewind()

		if r >= '0' && r <= '9' {
			if err = l.readNumber(); err != nil {
				return err
			}
		} else if exprWhiteSpace.MatchString(string(r)) {
			if err = l.readWhitespace(); err != nil {
				return err
			}
		} else if exprIdentStart.MatchString(string(r)) {
			if err = l.readIdentifier(); err != nil {
				return err
			}
		} else if r == '=' {
			if err = l.readEqualsOrMatches(); err != nil {
				return err
			}
		} else if r == '!' {
			if err = l.readNotEqualsOrNotMatches(); err != nil {
				return err
			}
		} else if r == '>' {
			if err = l.readGreater(); err != nil {
				return err
			}
		} else if r == '<' {
			if err = l.readLess(); err != nil {
				return err
			}
		} else if r == '"' || r == '\'' {
			if err = l.readStringLiteral(); err != nil {
				return err
			}
		} else if r == '&' {
			if err = l.readAnd(); err != nil {
				return err
			}
		} else if r == '|' {
			if err = l.readOr(); err != nil {
				return err
			}
		} else if r == '(' {
			line := l.line
			pos := l.pos
			l.advance()
			l.eat(r)
			l.push(LBrace, line, pos)
		} else if r == ')' {
			line := l.line
			pos := l.pos
			l.advance()
			l.eat(r)
			l.push(RBrace, line, pos)
		} else {
			line := l.line
			pos := l.pos
			l.advance()
			l.eat(r)
			l.push(Illegal, line, pos)
		}
	}
	return nil
}

func Tokenize(data string, skipWhitespace bool) (*TokenFlow, error) {
	l := lexer{
		sc:      bytes.NewBuffer([]byte(data)),
		line:    1,
		pos:     0,
		literal: "",
		tokens:  make([]Token, 0),
	}

	err := l.readAll()
	if err != nil {
		return nil, err
	}

	return newTokenFlow(l.tokens, skipWhitespace), nil
}
