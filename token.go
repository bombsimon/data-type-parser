package dtp

import (
	"bytes"
	"strconv"
	"unicode"
)

type Lexer struct {
	data []byte
	pos  int
}

func NewLexer(data []byte) *Lexer {
	return &Lexer{
		data: data,
		pos:  0,
	}
}

type TokenType string

const (
	TokenIdent   TokenType = "IDENT"
	TokenNumber  TokenType = "NUMBER"
	TokenLess    TokenType = "<"
	TokenGreater TokenType = ">"
	TokenLParen  TokenType = "("
	TokenRParen  TokenType = ")"
	TokenComma   TokenType = ","
	TokenStruct  TokenType = "STRUCT"
	TokenArray   TokenType = "ARRAY"
	TokenNotNull TokenType = "NOT NULL"
	TokenNewline TokenType = "\n"
	TokenEOF     TokenType = ""
)

//nolint:gochecknoglobals // Needed for now
var TokenMap = map[string]TokenType{
	"STRUCT": TokenStruct,
	"ARRAY":  TokenArray,
}

type Token struct {
	Type  TokenType
	Value string
}

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

			tokenType, ok := TokenMap[buf.String()]
			if !ok {
				tokenType = TokenIdent
			}

			return &Token{
				Type:  tokenType,
				Value: buf.String(),
			}
		}
	}

	return nil
}

func (l *Lexer) Peek() *Token {
	pos := l.pos
	defer func() {
		l.pos = pos
	}()

	token := l.Next()

	return token
}
