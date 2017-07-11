package gentab_test

import (
	"github.com/simonschneider/gentab"
	"os"
	"reflect"
	"strconv"
)

type (
	info struct {
		name        string
		secret      string `tab:"-"`
		description string `tab:"desc"`
	}
	container struct {
		info
		number int64 `tab:"number of smth"`
	}
	containers []container
)

func (c *containers) Footer() ([]string, error) {
	sum := int64(0)
	for _, n := range *c {
		sum += n.number
	}
	return []string{"", "Total", strconv.Itoa(int(sum))}, nil
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
		},
		container{
			info: info{
				name:        "good",
				secret:      "sweet",
				description: "bye",
			},
			number: int64(4),
		},
	}

	gentab.PrintTable(os.Stdout, cont, []reflect.Type{reflect.TypeOf(info{}), reflect.TypeOf(container{})})
	// Output:
	// +-------+-------+----------------+
	// | NAME  | DESC  | NUMBER OF SMTH |
	// +-------+-------+----------------+
	// | hello | world |              2 |
	// | good  | bye   |              4 |
	// +-------+-------+----------------+
	// |         TOTAL |       6        |
	// +-------+-------+----------------+
}
