# gentab

Dynamic table generation for golang using https://github.com/olekukonko/tablewriter

The package has one method `PrintTable(w io.Writer, in interface{})` which will accept a struct or slice of structs and create a cli table with it.

It uses tags `\`tab:""\`` for structs to have personalized headers for the tables.

If you want to override any of the Header, Body or Footer implementation you need to implement the functions `Header() ([]string, error)`, `Body() ([][]string, error)` or `Footer() ([]string, error)`.
