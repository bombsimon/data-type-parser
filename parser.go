package dtp

import (
	"fmt"
	"strconv"
)

// Ast is a simple representation of a data type node in the yntax tree. A node
// can have a name if it's a part of a struct or record but more importantly it
// has a data type and children. For nested data types this is how you would
// resolve the nested types. It also has a list of additional tokens such as
// nullability and a dedicated field for size.
type Ast struct {
	Name        string      `json:",omitempty"`
	DataType    string      `json:",omitempty"`
	Size        int         `json:",omitempty"`
	Children    []Ast       `json:",omitempty"`
	ExtraTokens []TokenType `json:",omitempty"`
}

// Parser implements how to parse the tokens from the lexer to AST node(s).
type Parser struct {
	Lexer *Lexer
}

// NewParser will create a new parser with a lexer over the passed data.
func NewParser(data string) *Parser {
	return &Parser{
		Lexer: NewLexer([]byte(data)),
	}
}

// Parse iterates over all tokens and parses the input (data type) to an AST.
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

// ParseTop is the top level entry point which parses top level data types.
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
			top, tOk := p.ParseTop()
			if !tOk {
				panic("Expected data type")
			}

			top.Name = i.Name

			return top, true
		}

		return i, ok
	case TokenContainer:
		p.Lexer.Next() // Consuem `<` or `(`

		return Ast{
			Name:        "",
			DataType:    string(t.Value),
			Children:    p.ParseContainer(),
			ExtraTokens: p.ParseExtraTokens(),
		}, true
	default:
		fmt.Println("TODO!", t)
	}

	return Ast{}, false
}

// ParseContainer parses container types and returns a slice with all its
// children. This is types such as `ARRAY`, `STRUCT`, `RANGE` etc.
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

// ParseIdent parses a simple ident. Either just by appending size information
// and extra tokens, or by resolving actual types, either flat ones or nested
// ones.
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

	// If the next token is a container token we know the ident is the name so
	// return the name and continue to parse recursively.
	if typ.Type == TokenContainer {
		return Ast{
			Name: ident,
		}, true
	}

	// All other token types that's not an ident mean something like a closing
	// container token or similar.
	if typ.Type != TokenIdent {
		return Ast{
			DataType: ident,
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

// ParseSize tries to parse size information by looking at a container pattern
// only containing a type.
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

// ParseExtraTokens parses known extra tokens by looking for specific idents
// that is known data or multiple idents making up for a known parameter such as
// `NOT NULL` or `PRIMARY KEY`.
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
