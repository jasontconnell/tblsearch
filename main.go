package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/tblsearch/conf"
	"github.com/jasontconnell/tblsearch/process"
)

func main() {
	start := time.Now()

	conffile := flag.String("c", "config.json", "the config file with connection string")
	str := flag.String("s", "", "the string to search")
	flag.Parse()

	if *str == "" {
		log.Fatal("need a string to search")
	}

	cfg := conf.LoadConfig(*conffile)

	if cfg.ConnectionString == "" {
		log.Fatal("need a connection string")
	}

	results, err := process.Search(cfg.ConnectionString, *str)
	if err != nil {
		log.Fatal("failed searching", err)
	}

	fmt.Println("found", len(results), "results")
	for _, res := range results {
		fmt.Printf("found %s in %s.%s at position %d\n", *str, res.Table.Name, res.Column.Name, res.Position)
	}
	fmt.Println("finished", time.Since(start))
}
