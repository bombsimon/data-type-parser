# Data Type Parser

This is a tiny parser to parse data types from warehouses. It's intended to be a
naive generic way to convert a datatype to an AST and that way making it easy to
fetch recursive data types.

> [!NOTE]
> This is in early development and mostly serves as a proof of concept. A lot of
> types and syntax is missing so it's not useful for production use!

```go
func parse() {
    lexer := dtp.NewLexer([]byte("ARRAY<STRUCT<bar INT NOT NULL, baz STRING>>"))
    parser := dtp.Parser{Lexer: lexer}
    ast := parser.Parse()

    j, _ := json.MarshalIndent(ast, "", "  ")
    fmt.Println(string(j))
}
```

```sh
[
  {
    "DataType": "ARRAY",
    "Children": [
      {
        "DataType": "STRUCT",
        "Children": [
          {
            "Name": "bar",
            "DataType": "INT",
            "ExtraTokens": [
              "NOT NULL"
            ]
          },
          {
            "Name": "baz",
            "DataType": "STRING"
          }
        ]
      }
    ]
  }
]

```

See the [example](./example) folder for more usage.
