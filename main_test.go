package gentab_test

import (
	"github.com/simonschneider/gentab"
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

	gentab.PrintTable(os.Stdout, cont, []reflect.Type{reflect.TypeOf(info{}), reflect.TypeOf(container{})})
	// Output:
	// +-------+-------+--------------+
	// | NAME  | DESC  | CONTAINER ID |
	// +-------+-------+--------------+
	// | hello | world |            1 |
	// | good  | bye   |            2 |
	// +-------+-------+--------------+
}
