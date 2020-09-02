package data

type Column struct {
	Name       string
	QueryName  string
	TypeName   string
	Nullable   bool
	Searchable bool
}

type Table struct {
	Name    string
	Columns []Column
}
