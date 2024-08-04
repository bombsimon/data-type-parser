# Data Type Parser

This is a tiny parser to parse data types from warehouses. It's intended to be a
naive generic way to convert a datatype to an AST and that way making it easy to
fetch recursive data types.

> [!NOTE]
> This is in early development and mostly serves as a proof of concept. A lot of
> types and syntax is missing so it's not useful for production use!

```go
func parse() {
    parser := dtp.NewParser("ARRAY<STRUCT<foo INT NOT NULL, bar STRING, baz RECORD<a INT, b FLOAT64>>>")
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
            "Name": "foo",
            "DataType": "INT",
            "ExtraTokens": [
              "NOT NULL"
            ]
          },
          {
            "Name": "bar",
            "DataType": "STRING"
          },
          {
            "Name": "baz",
            "DataType": "RECORD",
            "Children": [
              {
                "Name": "a",
                "DataType": "INT"
              },
              {
                "Name": "b",
                "DataType": "FLOAT64"
              }
            ]
          }
        ]
      }
    ]
  }
]
```

See the [example](./example) folder for more usage.
