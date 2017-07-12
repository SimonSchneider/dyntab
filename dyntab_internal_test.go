package dyntab

import (
	//"fmt"
	"reflect"
	"testing"
	"time"
)

type (
	nested struct {
		name2 string
	}

	footers []footer

	footer struct {
		number int
	}

	MyInt struct{ int }

	toString struct {
		ID MyInt
	}
)

func (f footers) Footer() ([]string, error) {
	return []string{"", "hey"}, nil
}

func (f MyInt) MarshalText() ([]byte, error) {
	return []byte("new"), nil
}

var tests = []struct {
	in           interface{}
	expectedHead []string
	expectedBody [][]string
	expectedFoot []string
}{
	{
		in: struct {
			id   int
			Name string `tab:"-"`
		}{
			1, "name",
		},
		expectedHead: []string{"id"},
		expectedBody: [][]string{[]string{"1"}},
		expectedFoot: nil,
	},
	{
		in: struct {
			id int
			Ti time.Time `tab:"tim"`
		}{
			1, time.Date(2009,
				time.November, 10,
				23, 0, 0, 0, time.UTC),
		},
		expectedHead: []string{"id", "tim"},
		expectedBody: [][]string{[]string{"1",
			"2009-11-10T23:00:00Z"}},
		expectedFoot: nil,
	},
	{
		in: []struct {
			id   int    `tab:"id2"`
			Name string `tab:"nam"`
		}{
			{1, "name1"},
			{2, "name2"},
		},
		expectedHead: []string{"id2", "nam"},
		expectedBody: [][]string{
			[]string{"1", "name1"},
			[]string{"2", "name2"},
		},
		expectedFoot: nil,
	},
	{
		in: []struct {
			id int `tab:"id2"`
			nested
		}{
			{1, nested{"naming1"}},
			{2, nested{"naming2"}},
		},
		expectedHead: []string{"id2", "name2"},
		expectedBody: [][]string{
			[]string{"1", "naming1"},
			[]string{"2", "naming2"},
		},
		expectedFoot: nil,
	},
	{
		in: footers{
			{1}, {2},
		},
		expectedHead: []string{"number"},
		expectedBody: [][]string{
			[]string{"1"},
			[]string{"2"},
		},
		expectedFoot: []string{"", "hey"},
	},
	{
		in: toString{
			MyInt{1},
		},
		expectedHead: []string{"ID"},
		expectedBody: [][]string{
			[]string{"new"},
		},
		expectedFoot: nil,
	},
}

func TestGetHeader(t *testing.T) {
	typesToRecurse = []reflect.Type{reflect.TypeOf(nested{})}
	for _, test := range tests {
		ret, err := getHeader(test.in)
		if err != nil {
			t.Error("Got an error", err, "when input", test)
			continue
		}
		if len(ret) != len(test.expectedHead) {
			t.Error("Expected", test.expectedHead, "got", ret)
			continue
		}
		for i, e := range test.expectedHead {
			if e != ret[i] {
				t.Error("Expected", test.expectedHead, "got", ret)
				continue
			}
		}
	}
}

func TestGetBody(t *testing.T) {
	typesToRecurse = []reflect.Type{reflect.TypeOf(nested{})}
	for _, test := range tests {
		ret, err := getBody(test.in)
		if err != nil {
			t.Error("Got an error", err, "when input", test)
			continue
		}
		if len(ret) != len(test.expectedBody) {
			t.Error("Expected", test.expectedBody, "got", ret)
			continue
		}
		for i, er := range test.expectedBody {
			if len(er) != len(ret[i]) {
				t.Error("Expected", test.expectedBody, "got", ret)
				continue
			}
			for j, e := range er {
				if e != ret[i][j] {
					t.Error("Expected", test.expectedBody, "got", ret)
					continue
				}
			}
		}
	}
}

func TestGetFooter(t *testing.T) {
	typesToRecurse = []reflect.Type{reflect.TypeOf(nested{})}
	for _, test := range tests {
		ret, err := getFooter(test.in)
		if err != nil {
			t.Error("Got an error", err, "when input", test)
			continue
		}
		if len(ret) != len(test.expectedFoot) {
			t.Error("Expected", test.expectedFoot, "got", ret)
			continue
		}
		for i, er := range test.expectedFoot {
			if er != ret[i] {
				t.Error("Expected", test.expectedFoot, "got", ret)
				continue
			}
		}
	}
}
