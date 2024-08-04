package dtp

import (
	"fmt"
	"strconv"
)

type Ast struct {
	Name        string      `json:",omitempty"`
	DataType    string      `json:",omitempty"`
	Size        int         `json:",omitempty"`
	Children    []Ast       `json:",omitempty"`
	ExtraTokens []TokenType `json:",omitempty"`
}

type Parser struct {
	Lexer *Lexer
}

func NewParser(data string) *Parser {
	return &Parser{
		Lexer: NewLexer([]byte(data)),
	}
}

func (p *Parser) Parse() []Ast {
	ast := []Ast{}

	for {
		a, ok := p.ParseTop()
		if !ok {
			break
		}

		ast = append(ast, a)
	}

	return ast
}

func (p *Parser) ParseTop() (Ast, bool) {
	t := p.Lexer.Next()
	if t == nil {
		return Ast{}, false
	}

	switch t.Type {
	case TokenIdent:
		i, ok := p.ParseIdent(t.Value)
		if !ok {
			return Ast{
				DataType:    t.Value,
				ExtraTokens: p.ParseExtraTokens(),
				Size:        p.ParseSize(),
			}, true
		}

		// We got an ident but without any data type, this means the data type
		// wasn't an ident but a complex type so we parse it and set our ident
		// as the name for the complex type.
		if i.DataType == "" {
			top, ok := p.ParseTop()
			if !ok {
				panic("Expected data type")
			}

			top.Name = i.Name

			return top, true
		}

		return i, ok
	case TokenStruct, TokenArray:
		p.Lexer.Next() // Consuem `<` or `(`

		return Ast{
			Name:        "",
			DataType:    string(t.Type),
			Children:    p.ParseContainer(),
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	default:
		fmt.Println("TODO!", t)
	}

	return Ast{}, false
}

func (p *Parser) ParseContainer() []Ast {
	var nodes []Ast

	for {
		child, ok := p.ParseTop()
		if !ok {
			break
		}

		nodes = append(nodes, child)

		next := p.Lexer.Peek()
		if next != nil {
			if next.Type == TokenComma {
				p.Lexer.Next() // Consume `,`
			}

			if next.Type == TokenGreater || next.Type == TokenRParen {
				break // End of container
			}
		}
	}

	p.Lexer.Next() // Consume `>` or `)`

	return nodes
}

func (p *Parser) ParseIdent(ident string) (Ast, bool) {
	// We only care about idents (name, type, nullability). If it's something
	// else it's most likely a `TokenGreater` or EOF.
	typ := p.Lexer.Peek()
	if typ == nil {
		return Ast{}, false
	}

	// Size container for ident only, not an alias
	if typ.Type == TokenLParen {
		return Ast{
			DataType:    ident,
			Size:        p.ParseSize(),
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	}

	if typ.Type != TokenIdent {
		return Ast{
			Name: ident,
		}, true
	}

	p.Lexer.Next() // Consume the data type, it's an ident

	switch typ.Type {
	case TokenIdent:
		return Ast{
			Name:        ident,
			DataType:    typ.Value,
			Size:        p.ParseSize(),
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	default:
		return Ast{}, false
	}
}

func (p *Parser) ParseSize() int {
	var (
		size int
		err  error
	)

	if lparen := p.Lexer.Peek(); lparen == nil || lparen.Type != TokenLParen {
		// No number type specified.
		return size
	}

	p.Lexer.Next() // Consume `(`

	sizeStr := p.Lexer.Next()
	if sizeStr == nil || sizeStr.Type != TokenNumber {
		panic("Expected size on type")
	}

	sizeInt, err := strconv.Atoi(sizeStr.Value)
	if err != nil {
		panic("Size is not an integer")
	}

	size = sizeInt

	p.Lexer.Next() // Consume `)`

	return size
}

func (p *Parser) ParseExtraTokens() []TokenType {
	var extraTokens []TokenType

	for {
		t := p.Lexer.Peek()
		if t == nil {
			return extraTokens
		}

		switch t.Type {
		case TokenComma:
			return extraTokens
		case TokenIdent:
			if t.Value == "NOT" {
				extraTokens = append(extraTokens, TokenNotNull)

				p.Lexer.Next() // NOT
				p.Lexer.Next() // NULL
			}
		case TokenGreater:
			return extraTokens
		default:
			return extraTokens
		}
	}
}
