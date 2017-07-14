# DYNamicTABles

[![Build Status](https://travis-ci.org/SimonSchneider/dyntab.svg?branch=master)](https://travis-ci.org/SimonSchneider/dyntab) [![GoDoc](https://godoc.org/github.com/SimonSchneider/dyntab?status.svg)](https://godoc.org/github.com/SimonSchneider/dyntab) [![Go Report Card](https://goreportcard.com/badge/github.com/simonschneider/dyntab)](https://goreportcard.com/report/github.com/simonschneider/dyntab)

Dynamic table generation for golang using https://github.com/olekukonko/tablewriter

The package uses an Table struct to hold main data.

To create a new table `table := dyntab.NewTable()`
To add recursion into types `table.Recurse([]reflect.Type{{reflect.TypeOf(time.Now())}})`
To add data `table.SetData(data)`
To set specialized to String methods `table.Specialize([]dyntab.ToSpecialize{time.Location{}, func(i interface{}) (string, error) {return "hey"}})`
To print the table `table.PrintTo(os.Stdout)`

The most simple usage is `err := dyntab.NewTable().SetData(data).PrintTo(os.Stdout)`

The data can be a struct or slice of structs, as well as a slice of reflect.types that it should recurse into. The output will be a cli table.

It uses struct tags `tab:"name"` to have personalized headers for the fields of the table (`tab:"-"` to ignore a field).

If you want to override any of the Header, Body or Footer implementation you need to implement the functions `Header() ([]string, error)`, `Body() ([][]string, error)` or `Footer() ([]string, error)`.

Example in the godoc references.
