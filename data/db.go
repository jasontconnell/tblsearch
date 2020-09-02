package data

type RowData struct {
	Table   Table
	ColData []ColData
}

type ColData struct {
	Column Column
	Value  string
	IsNull bool
}
