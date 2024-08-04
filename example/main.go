package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bombsimon/dtp"
)

const dataType = `STRUCT<
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
>`

func main() {
	dt := dataType
	if len(os.Args) > 1 {
		dt = os.Args[1]
	}

	parser := dtp.NewParser(dt)
	ast := parser.Parse()

	printASTJSON(ast)
	printAST(ast, []string{})
}

func printASTJSON(ast []dtp.Ast) {
	j, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(j))
}

func printAST(ast []dtp.Ast, path []string) {
	for _, a := range ast {
		switch a.DataType {
		case string(dtp.TokenArray), string(dtp.TokenStruct):
			p := path
			if a.Name != "" {
				p = append(p, a.Name)
				fmt.Printf("%-30s %s\n", strings.Join(p, "."), a.DataType)
			}

			printAST(a.Children, p)
		default:
			p := path
			p = append(p, a.Name)

			fmt.Printf("%-30s %s\n", strings.Join(p, "."), a.DataType)
		}
	}
}
