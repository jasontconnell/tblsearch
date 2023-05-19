package process

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/jasontconnell/tblsearch/data"
)

func Search(connstr, str string, reg *regexp.Regexp) ([]data.SearchResult, error) {
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
				results := findInCol(r, str, reg)

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

func findInCol(row data.RowData, str string, reg *regexp.Regexp) []data.ColData {
	cols := []data.ColData{}
	for _, c := range row.ColData {
		if c.IsNull {
			continue
		}

		if str != "" && strings.Contains(c.Value, str) {
			cols = append(cols, c)
		} else if reg != nil && reg.MatchString(c.Value) {
			cols = append(cols, c)
		}
	}
	return cols
}
