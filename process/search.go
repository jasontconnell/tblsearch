package process

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jasontconnell/tblsearch/data"
)

func Search(connstr, str string) ([]data.SearchResult, error) {
	tables, err := getSchema(connstr)

	if err != nil {
		return nil, fmt.Errorf("getting tables %s: %w", connstr, err)
	}

	searchResults := []data.SearchResult{}

	var wg sync.WaitGroup
	wg.Add(len(tables))

	for _, tbl := range tables {
		go func(tbl data.Table) {
			defer wg.Done()
			rows, err := load(connstr, tbl)
			if errors.Is(err, NoSearchError) {
				//fmt.Println("can't search", tbl.Name, ", continuing")
				return
			}
			if err != nil {
				//return nil, fmt.Errorf("couldn't load from table %s: %w", tbl.Name, err)
				return
			}
			for _, r := range rows {
				results := findInCol(r, str)

				for _, res := range results {
					pos := strings.Index(res.Value, str)
					searchResults = append(searchResults, data.SearchResult{Table: tbl, Column: res.Column, Position: pos})
				}
			}
		}(tbl)
	}

	wg.Wait()

	return searchResults, nil
}

func findInCol(row data.RowData, str string) []data.ColData {
	cols := []data.ColData{}
	for _, c := range row.ColData {
		if c.IsNull {
			continue
		}

		if strings.Contains(c.Value, str) {
			cols = append(cols, c)
		}
	}
	return cols
}
