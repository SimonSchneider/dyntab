# DYNamicTABles

[![GoDoc](https://godoc.org/github.com/SimonSchneider/dyntab?status.svg)](https://godoc.org/github.com/SimonSchneider/dyntab)

Dynamic table generation for golang using https://github.com/olekukonko/tablewriter

The package has one method `PrintTable(w io.Writer, in interface{}, typesToPrint []reflect.Type)` which will accept a struct or slice of structs, as well as a slice of reflect types that it should recurse into. The output will be a cli table.

It uses struct tags `tab:"name"` to have personalized headers for the fields of the table (`tab:"-"` to ignore a field).

If you want to override any of the Header, Body or Footer implementation you need to implement the functions `Header() ([]string, error)`, `Body() ([][]string, error)` or `Footer() ([]string, error)`.

Example:

```go
package dyntab_test

import (
	"github.com/simonschneider/dyntab"
	"os"
	"reflect"
	"strconv"
	"time"
)

type (
	MyTime struct{ time.Time }
	info   struct {
		name        string
		secret      string `tab:"-"`
		description string `tab:"desc"`
	}
	container struct {
		info
		number int64 `tab:"number of smth"`
		T      *MyTime
	}
	containers []container
)

func (t *MyTime) MarshalText() (text []byte, err error) {
	return []byte((*t).Format("2006-01-02")), nil
}

func (c *containers) Footer() ([]string, error) {
	sum := int64(0)
	for _, n := range *c {
		sum += n.number
	}
	return []string{"", "Total", strconv.Itoa(int(sum)), " "}, nil
}

func Example() {
	cont := &containers{
		container{
			info: info{
				name:        "hello",
				secret:      "pretty",
				description: "world",
			},
			number: int64(2),
			T:      &MyTime{time.Unix(1355270400, 0)},
		},
		container{
			info: info{
				name:        "good",
				secret:      "sweet",
				description: "bye",
			},
			number: int64(4),
			T:      &MyTime{time.Unix(1355270400, 0)},
		},
	}

	dyntab.PrintTable(os.Stdout, cont, []reflect.Type{reflect.TypeOf(info{}), reflect.TypeOf(container{})})
	// Output:
	// +-------+-------+----------------+------------+
	// | NAME  | DESC  | NUMBER OF SMTH |     T      |
	// +-------+-------+----------------+------------+
	// | hello | world |              2 | 2012-12-12 |
	// | good  | bye   |              4 | 2012-12-12 |
	// +-------+-------+----------------+------------+
	// |         TOTAL |       6        |            |
	// +-------+-------+----------------+------------+

}
```
