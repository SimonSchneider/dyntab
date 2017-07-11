// Package dyntab creates dynamic tables for structs or slices of structs
package dyntab

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"reflect"
	"time"
)

var typesToPrint []reflect.Type

type (
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
)

// PrintTable prints a table of the interface
func PrintTable(w io.Writer, in interface{}, toPrint []reflect.Type) (err error) {
	var header, footer []string
	var body [][]string
	typesToPrint = toPrint
	header, err = getHeader(in)
	body, err = getBody(in)
	footer, err = getFooter(in)
	tab := tablewriter.NewWriter(w)
	tab.SetHeader(header)
	tab.AppendBulk(body)
	tab.SetFooter(footer)
	tab.Render()
	return nil
}

func getHeader(in interface{}) ([]string, error) {
	n, ok := in.(TabHeader)
	if ok {
		return (n).Header()
	}
	i := reflect.Indirect(reflect.ValueOf(in)).Interface()
	return getTypeHeader(reflect.TypeOf(i))
}

func getTypeHeader(in reflect.Type) (s []string, err error) {
	if in.Kind() == reflect.Slice {
		in = in.Elem()
	}
	if in.Kind() != reflect.Struct {
		return nil, errors.New("Not possible to find struct")
	}
	return getStructHeader(in), nil
}

func getStructHeader(in reflect.Type) (s []string) {
	for i := 0; i < in.NumField(); i++ {
		field := in.Field(i)
		t := field.Tag.Get("tab")
		if t == "-" {
			continue
		}
		if contains(typesToPrint, field.Type) {
			s = append(s, getStructHeader(field.Type)...)
		} else if t != "" {
			s = append(s, t)
		} else {
			s = append(s, field.Name)
		}
	}
	return s
}

func getBody(in interface{}) (s [][]string, err error) {
	n, ok := in.(TabBody)
	if ok {
		return (n).Body()
	}
	i := reflect.Indirect(reflect.ValueOf(in)).Interface()
	return getValueBody(reflect.ValueOf(i)), nil
}

func getValueBody(in reflect.Value) (s [][]string) {
	if in.Type().Kind() == reflect.Slice {
		for i := 0; i < in.Len(); i++ {
			v := in.Index(i)
			s = append(s, getStructBody(v))
		}
	} else if in.Type().Kind() == reflect.Struct {
		s = append(s, getStructBody(in))
	}
	return
}

func getStructBody(in reflect.Value) (s []string) {
	st := in.Type()
	for i := 0; i < st.NumField(); i++ {
		v := in.Field(i)
		field := st.Field(i)
		t := field.Tag.Get("tab")
		if t == "-" {
			continue
		}
		if contains(typesToPrint, field.Type) {
			s = append(s, getStructBody(v)...)
		} else {
			s = append(s, getString(v))
		}
	}
	return
}

func getString(in reflect.Value) string {
	switch in.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", in.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", in.Float())
	case reflect.String:
		return in.String()
	}
	if in.Type() == reflect.TypeOf(time.Time{}) {
		t, ok := in.Interface().(time.Time)
		if ok {
			return t.Format("2006-01-02")
		}
	}
	return ""
}

func getFooter(in interface{}) ([]string, error) {
	n, ok := in.(TabFooter)
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