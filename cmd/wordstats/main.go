package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dkvz/blog-stats/pkg/cli"
	"github.com/dkvz/blog-stats/pkg/db"
)

func main() {
	config, err := cli.ConfigFromDotEnv()
	if err != nil {
		panic(err)
	}

	if config.DbPath == "" {
		fmt.Println("missing DB_PATH in env")
	}

	// Parse the current mode from CLI args
	mode := flag.String("mode", "", "plot | factors | verify")
	flag.Parse()
	*mode = strings.TrimSpace(strings.ToLower(*mode))

	switch *mode {
	case "plot":
		dbs, err := db.NewDBSqlite(config.DbPath)
		if err != nil {
			fmt.Println("encountered DB error")
			panic(err)
		}
		if ids, err := dbs.AllPublishedArticleIds(); err == nil {
			fmt.Printf("%v", ids)
		} else {
			fmt.Println("encountered DB error for query")
			panic(err)
		}

	case "factors":
		fmt.Println("Factors mode")
	case "verify":
		fmt.Println("Verify mode")
	default:
		flag.Usage()
		return
	}

}
