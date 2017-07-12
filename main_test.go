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
		T      MyTime
	}
	containers []container
)

func (t MyTime) MarshalText() (text []byte, err error) {
	return []byte(t.Format("2006-01-02")), nil
}

func (c containers) Footer() ([]string, error) {
	sum := int64(0)
	for _, n := range c {
		sum += n.number
	}
	return []string{"", "Total", strconv.Itoa(int(sum)), " "}, nil
}

func Example() {
	cont := containers{
		container{
			info: info{
				name:        "hello",
				secret:      "pretty",
				description: "world",
			},
			number: int64(2),
			T:      MyTime{time.Unix(1355270400, 0)},
		},
		container{
			info: info{
				name:        "good",
				secret:      "sweet",
				description: "bye",
			},
			number: int64(4),
			T:      MyTime{time.Unix(1355270400, 0)},
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
