package dyntab_test

import (
	"errors"
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
		Loc    time.Location `tab:"location"`
	}
	containers []container
)

func (t MyTime) MarshalText() (text []byte, err error) {
	return []byte(t.Format("2006-01-02")), nil
}

func Loc2String(i interface{}) (string, error) {
	l, ok := i.(time.Location)
	if ok {
		return "in: " + l.String(), nil
	}
	return "", errors.New("not time.location")
}

func (c containers) Footer() ([]string, error) {
	sum := int64(0)
	for _, n := range c {
		sum += n.number
	}
	return []string{"", "Total", strconv.Itoa(int(sum)), "", " "}, nil
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
			Loc:    *time.UTC,
		},
		container{
			info: info{
				name:        "good",
				secret:      "sweet",
				description: "bye",
			},
			number: int64(4),
			T:      MyTime{time.Unix(1355270400, 0)},
			Loc:    *time.UTC,
		},
	}

	dyntab.NewTable().
		SetData(cont).
		Recurse([]reflect.Type{
			reflect.TypeOf(info{}),
			reflect.TypeOf(container{})}).
		Specialize([]dyntab.ToSpecialize{{
			Type:     reflect.TypeOf(time.Location{}),
			ToString: Loc2String}}).
		PrintTo(os.Stdout)
	// Output:
	// +-------+-------+----------------+------------+----------+
	// | NAME  | DESC  | NUMBER OF SMTH |     T      | LOCATION |
	// +-------+-------+----------------+------------+----------+
	// | hello | world |              2 | 2012-12-12 | in: UTC  |
	// | good  | bye   |              4 | 2012-12-12 | in: UTC  |
	// +-------+-------+----------------+------------+----------+
	// |         TOTAL |       6        |                       |
	// +-------+-------+----------------+------------+----------+
}
