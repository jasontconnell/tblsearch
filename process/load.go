package process

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jasontconnell/sqlhelp"
	"github.com/jasontconnell/tblsearch/data"
	_ "github.com/microsoft/go-mssqldb"
)

var NoSearchError = errors.New("Can't search")

func load(connstr string, table data.Table) ([]data.RowData, error) {
	db, err := sql.Open("mssql", connstr)
	defer db.Close()

	if err != nil {
		return nil, err
	}

	if !canSearch(table) {
		return nil, NoSearchError
	}

	query := buildQuery(table)

	records, err := sqlhelp.GetResultSet(db, query)

	if err != nil {
		return nil, fmt.Errorf("executing sql %s: %w", query, err)
	}

	rows := []data.RowData{}
	for _, r := range records {
		rd := data.RowData{Table: table}

		for _, c := range table.Columns {
			cd := data.ColData{Column: c}
			val := r[c.Name]
			if c.Nullable && val == nil {
				cd.IsNull = true
			} else if c.Searchable {
				cd.Value = val.(string)
			}

			rd.ColData = append(rd.ColData, cd)
		}

		rows = append(rows, rd)
	}

	return rows, nil
}

func canSearch(table data.Table) bool {
	searchable := false
	for _, c := range table.Columns {
		if c.Searchable {
			searchable = true
			break
		}
	}
	return searchable
}

func buildQuery(table data.Table) string {
	sqlfmt := `select * from [%s]`

	return fmt.Sprintf(sqlfmt, table.Name)
}
