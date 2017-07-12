# DYNamicTABles

[![GoDoc](https://godoc.org/github.com/SimonSchneider/dyntab?status.svg)](https://godoc.org/github.com/SimonSchneider/dyntab) [![Go Report Card](https://goreportcard.com/badge/github.com/simonschneider/dyntab)](https://goreportcard.com/report/github.com/simonschneider/dyntab)

Dynamic table generation for golang using https://github.com/olekukonko/tablewriter

The package has one method `PrintTable(w io.Writer, in interface{}, typesToPrint []reflect.Type)` which will accept a struct or slice of structs, as well as a slice of reflect types that it should recurse into. The output will be a cli table.

It uses struct tags `tab:"name"` to have personalized headers for the fields of the table (`tab:"-"` to ignore a field).

If you want to override any of the Header, Body or Footer implementation you need to implement the functions `Header() ([]string, error)`, `Body() ([][]string, error)` or `Footer() ([]string, error)`.

Example: in the godoc references
