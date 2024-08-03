package dtp

import (
	"fmt"
)

type Ast struct {
	Name        string      `json:",omitempty"`
	DataType    string      `json:",omitempty"`
	Children    []Ast       `json:",omitempty"`
	ExtraTokens []TokenType `json:",omitempty"`
}

type Parser struct {
	Lexer *Lexer
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
		i, ok := p.ParseIdent()
		if !ok {
			return Ast{
				DataType:    t.Value,
				ExtraTokens: p.ParseExtraTokens(),
			}, true
		}

		return i, ok
	case TokenStruct:
		return p.ParseStruct()
	case TokenArray:
		a, ok := p.ParseArray()
		if !ok {
			panic("Expected array")
		}

		return Ast{
			Name:        "",
			DataType:    string(TokenArray),
			Children:    []Ast{a},
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	default:
		fmt.Println("TODO!", t)
	}

	return Ast{}, false
}

func (p *Parser) ParseStruct() (Ast, bool) {
	t := p.Lexer.Next()
	if t.Type != TokenLess {
		panic(fmt.Sprintf("Expected %v, got %v", TokenLess, t.Type))
	}

	children := []Ast{}

	for {
		child, ok := p.ParseIdent()
		if !ok {
			break
		}

		children = append(children, child)

		if t := p.Lexer.Peek(); t.Type == TokenComma {
			p.Lexer.Next() // Consume `,`
		}
	}

	p.Lexer.Next() // Consume >

	return Ast{
		DataType:    string(TokenStruct),
		Children:    children,
		ExtraTokens: p.ParseExtraTokens(),
	}, true
}

func (p *Parser) ParseArray() (Ast, bool) {
	t := p.Lexer.Next()
	if t.Type != TokenLess {
		panic(fmt.Sprintf("Expected %v, got %v", TokenLess, t.Type))
	}

	a, ok := p.ParseTop()

	p.Lexer.Next() // Consume >

	return a, ok
}

func (p *Parser) ParseIdent() (Ast, bool) {
	// We only care about idents (name, type, nullability). If it's something
	// else it's most likely a `TokenGreater` or EOF.
	if name := p.Lexer.Peek(); name != nil && name.Type != TokenIdent {
		return Ast{}, false
	}

	name := p.Lexer.Next()
	if name == nil {
		return Ast{}, false
	}

	typ := p.Lexer.Next()
	if typ == nil {
		return Ast{}, false
	}

	switch typ.Type {
	case TokenIdent:
		return Ast{
			Name:        name.Value,
			DataType:    typ.Value,
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	case TokenStruct:
		child, ok := p.ParseStruct()
		if !ok {
			panic("Expected child")
		}

		child.Name = name.Value

		return child, true
	case TokenArray:
		child, ok := p.ParseArray()
		if !ok {
			panic("Expected child")
		}

		return Ast{
			Name:        name.Value,
			DataType:    string(TokenArray),
			Children:    []Ast{child},
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	}

	return Ast{}, false
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
