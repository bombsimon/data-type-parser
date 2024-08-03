package dtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		description string
		dataType    string
		expected    []Ast
	}{
		{
			description: "Simple type",
			dataType:    "STRING",
			expected: []Ast{
				{
					DataType:    "STRING",
					ExtraTokens: []TokenType{},
				},
			},
		},
		{
			description: "Struct",
			dataType:    "STRUCT<id INT, name STRING NOT NULL> NOT NULL",
			expected: []Ast{
				{
					DataType:    "STRUCT",
					ExtraTokens: []TokenType{TokenNotNull},
					Children: []Ast{
						{
							Name:        "id",
							DataType:    "INT",
							ExtraTokens: []TokenType{},
						},
						{
							Name:        "name",
							DataType:    "STRING",
							ExtraTokens: []TokenType{TokenNotNull},
						},
					},
				},
			},
		},
		{
			description: "Complex combinations",
			dataType: `STRUCT<
						  customer_id STRING NOT NULL,
						  order_details ARRAY<
							STRUCT<
							  product_id STRING,
							  quantity INT64,
							  price FLOAT64
							>
						  >,
						  shipping_address STRUCT<
							street STRING,
							city STRING,
							state STRING,
							zip_code STRING
						  > NOT NULL
						>`,
			expected: []Ast{
				{
					DataType:    "STRUCT",
					ExtraTokens: []TokenType{},
					Children: []Ast{
						{
							Name:        "customer_id",
							DataType:    "STRING",
							ExtraTokens: []TokenType{TokenNotNull},
						},
						{
							Name:        "order_details",
							DataType:    "ARRAY",
							ExtraTokens: []TokenType{},
							Children: []Ast{
								{
									DataType:    "STRUCT",
									ExtraTokens: []TokenType{},
									Children: []Ast{
										{
											Name:        "product_id",
											DataType:    "STRING",
											ExtraTokens: []TokenType{},
										},
										{
											Name:        "quantity",
											DataType:    "INT64",
											ExtraTokens: []TokenType{},
										},
										{
											Name:        "price",
											DataType:    "FLOAT64",
											ExtraTokens: []TokenType{},
										},
									},
								},
							},
						},
						{
							Name:        "shipping_address",
							DataType:    "STRUCT",
							ExtraTokens: []TokenType{TokenNotNull},
							Children: []Ast{
								{
									Name:        "street",
									DataType:    "STRING",
									ExtraTokens: []TokenType{},
								},
								{
									Name:        "city",
									DataType:    "STRING",
									ExtraTokens: []TokenType{},
								},
								{
									Name:        "state",
									DataType:    "STRING",
									ExtraTokens: []TokenType{},
								},
								{
									Name:        "zip_code",
									DataType:    "STRING",
									ExtraTokens: []TokenType{},
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			lexer := NewLexer([]byte(tc.dataType))
			parser := Parser{Lexer: lexer}

			assert.Equal(t, tc.expected, parser.Parse())
		})
	}
}
