# Data Type Parser

This is a tiny parser to parse data types from warehouses. It's intended to be a
naive generic way to convert a datatype to an AST and that way making it easy to
fetch recursive data types.

```go
func parse() {
    lexer := dtp.NewLexer([]byte("ARARY<STRUCT<bar INT NOT NULL, baz STRING>>"))
    parser := dtp.Parser{Lexer: lexer}
    ast := parser.Parse()

    j, _ := json.MarshalIndent(ast, "", "  ")
    fmt.Println(string(j))
}
```

```sh
[
  {
    "Name": "",
    "DataType": "STRUCT",
    "Children": [
      {
        "Name": "bar",
        "DataType": "INT",
        "Children": null,
        "ExtraTokens": [
          "NOT NULL"
        ]
      },
      {
        "Name": "baz",
        "DataType": "STRING",
        "Children": null,
        "ExtraTokens": []
      }
    ],
    "ExtraTokens": []
  }
]
```

See the [example](./example) folder for more usage.
