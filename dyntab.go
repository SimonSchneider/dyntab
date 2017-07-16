// Package dyntab creates dynamic tables for structs or slices of structs
package dyntab

import (
	"encoding"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"reflect"
)

type (
	//Table holds the main table data
	Table struct {
		data              interface{}
		typesToRecurse    []reflect.Type
		typesToSpecialize toSpecializes
	}

	// TabFooter interface can be implemented to override the
	// default footer creation
	TabFooter interface {
		Footer() ([]string, error)
	}

	// TabHeader interface can be implemented to override the
	// default header creation
	TabHeader interface {
		Header() ([]string, error)
	}

	// TabBody interface can be implemented to override the
	// default body creation
	TabBody interface {
		Body() ([][]string, error)
	}

	// ToSpecialize is struct for holding to string specializations
	ToSpecialize struct {
		Type     reflect.Type
		ToString func(interface{}) (string, error)
	}

	toSpecializes []ToSpecialize
)

//NewTable returns a new table ready for printing
func NewTable() *Table {
	return &Table{}
}

//SetData sets the data of the table
func (t *Table) SetData(i interface{}) *Table {
	t.data = i
	return t
}

//PrintTo prints the table to the io.Writer
func (t Table) PrintTo(w io.Writer) (err error) {
	var header, footer []string
	var body [][]string
	if t.data == nil {
		return errors.New("no data")
	}
	header, err = t.getHeader()
	if err != nil {
		return err
	}
	body, err = t.getBody()
	if err != nil {
		return err
	}
	footer, err = t.getFooter()
	if err != nil {
		return err
	}
	tab := tablewriter.NewWriter(w)
	tab.SetHeader(header)
	tab.AppendBulk(body)
	tab.SetFooter(footer)
	tab.Render()
	return nil
}

//Recurse sets the reflect.Types, Recurse slice is required so it's possible to determine what structs to print as a column and which to recurse into,
func (t *Table) Recurse(r []reflect.Type) *Table {
	t.typesToRecurse = r
	return t
}

//Specialize sets the toSpecialize types to the specialize
func (t *Table) Specialize(s []ToSpecialize) *Table {
	t.typesToSpecialize = s
	return t
}

func (t Table) getHeader() ([]string, error) {
	n, ok := (t.data).(TabHeader)
	if ok {
		return (n).Header()
	}
	i := reflect.Indirect(reflect.ValueOf(t.data)).Interface()
	return getTypeHeader(reflect.TypeOf(i), t.typesToRecurse)
}

func getTypeHeader(in reflect.Type, ttr []reflect.Type) (s []string, err error) {
	if in.Kind() == reflect.Slice {
		in = in.Elem()
	}
	if in.Kind() != reflect.Struct {
		return nil, errors.New("Not possible to find struct")
	}
	return getStructHeader(in, ttr), nil
}

func getStructHeader(in reflect.Type, ttr []reflect.Type) (s []string) {
	for i := 0; i < in.NumField(); i++ {
		field := in.Field(i)
		t := field.Tag.Get("tab")
		if t == "-" {
			continue
		}
		if contains(ttr, field.Type) {
			s = append(s, getStructHeader(field.Type, ttr)...)
		} else if t != "" {
			s = append(s, t)
		} else {
			s = append(s, field.Name)
		}
	}
	return s
}

func (t Table) getBody() (s [][]string, err error) {
	n, ok := t.data.(TabBody)
	if ok {
		return (n).Body()
	}
	i := reflect.Indirect(reflect.ValueOf(t.data)).Interface()
	return getValueBody(reflect.ValueOf(i), t.typesToRecurse, t.typesToSpecialize), nil
}

func getValueBody(in reflect.Value, ttr []reflect.Type, tts toSpecializes) (s [][]string) {
	if in.Type().Kind() == reflect.Slice {
		for i := 0; i < in.Len(); i++ {
			v := in.Index(i)
			s = append(s, getStructBody(v, ttr, tts))
		}
	} else if in.Type().Kind() == reflect.Struct {
		s = append(s, getStructBody(in, ttr, tts))
	}
	return
}

func getStructBody(in reflect.Value, ttr []reflect.Type, tts toSpecializes) (s []string) {
	st := in.Type()
	for i := 0; i < st.NumField(); i++ {
		v := in.Field(i)
		field := st.Field(i)
		t := field.Tag.Get("tab")
		if t == "-" {
			continue
		}
		if contains(ttr, field.Type) {
			s = append(s, getStructBody(v, ttr, tts)...)
		} else {
			s = append(s, getString(v, tts))
		}
	}
	return
}

func getString(in reflect.Value, tts toSpecializes) string {
	if in.CanInterface() {
		if f, ok := tts.contains(in.Type()); ok {
			s, err := f(in.Interface())
			if err == nil {
				return s
			}
		}
		t, ok := in.Interface().(encoding.TextMarshaler)
		if ok {
			s, err := t.MarshalText()
			if err == nil {
				return string(s)
			}
		}
	}
	switch in.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", in.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", in.Float())
	case reflect.String:
		return in.String()
	}
	return ""
}

func (ts toSpecializes) contains(in reflect.Type) (func(interface{}) (string, error), bool) {
	for _, t := range ts {
		if t.Type == in {
			return t.ToString, true
		}
	}
	return nil, false
}

func (t Table) getFooter() ([]string, error) {
	n, ok := t.data.(TabFooter)
	if ok {
		return (n).Footer()
	}
	return []string{}, nil
}

func contains(in []reflect.Type, t reflect.Type) bool {
	for _, a := range in {
		if a == t {
			return true
		}
	}
	return false
}
