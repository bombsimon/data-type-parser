package dtp

import (
	"bytes"
	"strconv"
	"unicode"
)

// Lexer represents a lexer/tokenizer that parses the input stream into valuable
// tokens that can be used by the parser.
type Lexer struct {
	data []byte
	pos  int
}

// NewLexer creates a new lexer over the passed data.
func NewLexer(data []byte) *Lexer {
	return &Lexer{
		data: data,
		pos:  0,
	}
}

// TokenType represents a type of token.
type TokenType string

// List of available tokens.
const (
	TokenIdent     TokenType = "IDENT"
	TokenNumber    TokenType = "NUMBER"
	TokenLess      TokenType = "<"
	TokenGreater   TokenType = ">"
	TokenLParen    TokenType = "("
	TokenRParen    TokenType = ")"
	TokenComma     TokenType = ","
	TokenContainer TokenType = "CONTAINER"
	TokenNotNull   TokenType = "NOT NULL"
	TokenNewline   TokenType = "\n"
	TokenEOF       TokenType = ""
)

//nolint:gochecknoglobals // Needed for now
var TokenMap = map[string]struct{}{
	"STRUCT": {},
	"ARRAY":  {},
	"RANGE":  {},
	"RECORD": {},
}

// A token is a wrapper over a token for tokens that can contain values. This is
// most of the time only idents or numbers where the value would be the actual
// value from the stream and the type would be a known token type.
type Token struct {
	Type  TokenType
	Value string
}

// Next returns the next token in the stream or nil if the stream is ended.
func (l *Lexer) Next() *Token {
	if l.pos >= len(l.data) {
		return nil
	}

	v := l.data[l.pos]

	for unicode.IsSpace(rune(v)) {
		l.pos++
		if l.pos >= len(l.data) {
			return nil
		}

		v = l.data[l.pos]
	}

	switch v {
	case ',':
		l.pos++
		return &Token{Type: TokenComma}
	case '<':
		l.pos++
		return &Token{Type: TokenLess}
	case '>':
		l.pos++
		return &Token{Type: TokenGreater}
	case '(':
		l.pos++
		return &Token{Type: TokenLParen}
	case ')':
		l.pos++
		return &Token{Type: TokenRParen}
	case '\n':
		l.pos++
		return &Token{Type: TokenNewline}
	default:
		if unicode.IsLetter(rune(v)) || unicode.IsDigit((rune(v))) {
			buf := bytes.Buffer{}
			start := l.pos

		I:
			for {
				l.pos++
				if l.pos >= len(l.data) {
					break
				}

				posVal := rune(l.data[l.pos])
				switch posVal {
				case ' ', '<', '>', '(', ')', ',', '\n':
					break I
				default:
					continue
				}
			}

			buf.Write(l.data[start:l.pos])

			if _, err := strconv.Atoi(buf.String()); err == nil {
				return &Token{
					Type:  TokenNumber,
					Value: buf.String(),
				}
			}

			tokenType := TokenIdent
			if _, tokenIsContainer := TokenMap[buf.String()]; tokenIsContainer {
				tokenType = TokenContainer
			}

			return &Token{
				Type:  tokenType,
				Value: buf.String(),
			}
		}
	}

	return nil
}

// Peek peeks at the next token in the stream without consuming it.
func (l *Lexer) Peek() *Token {
	pos := l.pos
	defer func() {
		l.pos = pos
	}()

	token := l.Next()

	return token
}
