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
					DataType: "STRING",
				},
			},
		},
		{
			description: "Number on type",
			dataType:    "STRING(10)",
			expected: []Ast{
				{
					DataType: "STRING",
					Size:     10,
				},
			},
		},
		{
			description: "Number in struct",
			dataType:    "STRUCT<a STRING(10)>",
			expected: []Ast{
				{
					DataType: "STRUCT",
					Children: []Ast{
						{
							Name:     "a",
							DataType: "STRING",
							Size:     10,
						},
					},
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
							Name:     "id",
							DataType: "INT",
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
					DataType: "STRUCT",
					Children: []Ast{
						{
							Name:        "customer_id",
							DataType:    "STRING",
							ExtraTokens: []TokenType{TokenNotNull},
						},
						{
							Name:     "order_details",
							DataType: "ARRAY",
							Children: []Ast{
								{
									DataType: "STRUCT",
									Children: []Ast{
										{
											Name:     "product_id",
											DataType: "STRING",
										},
										{
											Name:     "quantity",
											DataType: "INT64",
										},
										{
											Name:     "price",
											DataType: "FLOAT64",
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
									Name:     "street",
									DataType: "STRING",
								},
								{
									Name:     "city",
									DataType: "STRING",
								},
								{
									Name:     "state",
									DataType: "STRING",
								},
								{
									Name:     "zip_code",
									DataType: "STRING",
								},
							},
						},
					},
				},
			},
		},
		{
			description: "Array of structs",
			dataType:    "ARRAY<STRUCT<bar INT NOT NULL, baz STRING>>",
			expected: []Ast{
				{
					DataType: "ARRAY",
					Children: []Ast{
						{
							DataType: "STRUCT",
							Children: []Ast{
								{
									Name:        "bar",
									DataType:    "INT",
									ExtraTokens: []TokenType{TokenNotNull},
								},
								{
									Name:     "baz",
									DataType: "STRING",
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
