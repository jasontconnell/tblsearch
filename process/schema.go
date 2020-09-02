package process

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jasontconnell/sqlhelp"
	"github.com/jasontconnell/tblsearch/data"
)

func getSchema(connstr string) ([]data.Table, error) {
	db, err := sql.Open("mssql", connstr)
	defer db.Close()

	if err != nil {
		return nil, err
	}
	records, err := sqlhelp.GetResultSet(db, "select TABLE_NAME from INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE'")

	if err != nil {
		return nil, err
	}

	tables := []data.Table{}
	for _, r := range records {
		t := data.Table{}
		t.Name = r["TABLE_NAME"].(string)

		cols, err := getColumns(db, t.Name)
		if err != nil {
			return tables, err
		}
		t.Columns = cols
		tables = append(tables, t)
	}

	return tables, nil
}

func getColumns(db *sql.DB, table string) ([]data.Column, error) {
	sqlfmt := "select COLUMN_NAME, IS_NULLABLE, DATA_TYPE from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME = '%v' ORDER BY ORDINAL_POSITION"
	q := fmt.Sprintf(sqlfmt, table)

	records, err := sqlhelp.GetResultSet(db, q)
	if err != nil {
		return nil, err
	}

	cols := []data.Column{}
	for _, r := range records {
		t := r["DATA_TYPE"].(string)

		if !isString(t) {
			switch t {
			case "bit", "uniqueidentifier", "int", "datetime", "float", "decimal", "bigint", "smallint", "datetime2":
				continue
			}
			continue
		}

		col := data.Column{}

		col.Name = r["COLUMN_NAME"].(string)
		col.QueryName = "[" + col.Name + "]"
		col.Nullable = r["IS_NULLABLE"].(string) != "NO"
		col.TypeName = t
		col.Searchable = true
		cols = append(cols, col)
	}

	return cols, nil
}

func isString(t string) bool {
	return strings.Contains(t, "char")
}
