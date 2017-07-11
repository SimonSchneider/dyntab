# gentab

Dynamic table generation for golang using https://github.com/olekukonko/tablewriter

The package has one method `PrintTable(w io.Writer, in interface{}, typesToPrint []reflect.Type)` which will accept a struct or slice of structs, as well as a slice of reflect types that it should recurse into. The output will be a cli table.

It uses tags `\`tab:""\`` for structs to have personalized headers for the tables.

If you want to override any of the Header, Body or Footer implementation you need to implement the functions `Header() ([]string, error)`, `Body() ([][]string, error)` or `Footer() ([]string, error)`.

Example:

```go
package gentab

import (
    "os"
    "reflect"
)

type (
    info struct {
        name        string
        secret      string `tab:"-"`
        description string `tab:"desc"`
    }
    container struct {
        info
        id int64 `tab:"container id"`
    }
)

func Example() {
    cont := []container{
        container{
            info: info{
                name:        "hello",
                secret:      "pretty",
                description: "world",
            },
            id: int64(1),
        },
        container{
            info: info{
                name:        "good",
                secret:      "sweet",
                description: "bye",
            },
            id: int64(2),
        },
    }

    PrintTable(os.Stdout, cont, []reflect.Type{reflect.TypeOf(info{}), reflect.TypeOf(container{})})
    // Output:
    // +-------+-------+--------------+
    // | NAME  | DESC  | CONTAINER ID |
    // +-------+-------+--------------+
    // | hello | world |            1 |
    // | good  | bye   |            2 |
    // +-------+-------+--------------+
}
```
