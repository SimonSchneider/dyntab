// Package gentab creates dynamic tables for structs or slices of structs
package gentab

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"reflect"
)

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
func PrintTable(w io.Writer, in interface{}) (err error) {
	var header, footer []string
	var body [][]string
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

func getHeader(inn interface{}) ([]string, error) {
	n, ok := inn.(TabHeader)
	if ok {
		return (n).Header()
	}
	in := reflect.Indirect(reflect.ValueOf(inn)).Interface()
	st, err := findStruct(reflect.TypeOf(in))
	if err != nil {
		return nil, err
	}
	s := []string{}
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		t := f.Tag.Get("tab")
		if t != "" {
			if t == "-" {
				continue
			}
			s = append(s, t)
		} else if t != "-" {
			s = append(s, f.Name)
		}
	}
	return s, nil
}

func findStruct(in reflect.Type) (reflect.Type, error) {
	switch in.Kind() {
	case reflect.Slice:
		return findStruct(in.Elem())
	case reflect.Struct:
		return in, nil
	}
	return nil, errors.New("no struct")
}

func getBody(inn interface{}) (s [][]string, err error) {
	n, ok := inn.(TabBody)
	if ok {
		return (n).Body()
	}
	in := reflect.Indirect(reflect.ValueOf(inn)).Interface()
	v := reflect.ValueOf(in)
	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			v.Index(i)
			new := getBodyLine(v.Index(i).Interface())
			s = append(s, new)
		}
	case reflect.Struct:
		s = append(s, getBodyLine(in))
	}
	return s, nil
}

func getBodyLine(in interface{}) (s []string) {
	st := reflect.TypeOf(in)
	for i := 0; i < st.NumField(); i++ {
		v := reflect.ValueOf(in).Field(i)
		t := st.Field(i).Tag.Get("tab")
		if t != "-" {
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
	return ""
}

func getFooter(in interface{}) ([]string, error) {
	n, ok := in.(TabFooter)
	if ok {
		return (n).Footer()
	}
	return []string{}, nil
}
