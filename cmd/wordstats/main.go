package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dkvz/blog-stats/pkg/cli"
	"github.com/dkvz/blog-stats/pkg/db"
	"github.com/dkvz/blog-stats/pkg/stats"
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

	iMode := 0

	switch *mode {
	case "plot":
		iMode = 0
	case "factors":
		iMode = 1
	case "verify":
		iMode = 2
	default:
		flag.Usage()
		return
	}

	dbs, err := db.NewDBSqlite(config.DbPath)
	if err != nil {
		fmt.Println("encountered DB error")
		panic(err)
	}

	if iMode == 0 {
		runModePlot(dbs)
	}
}

func runModePlot(dbs *db.DbSqlite) error {
	ids, err := dbs.AllPublishedArticleIds()
	if err != nil {
		fmt.Println("encountered DB error getting article ids")
		panic(err)
	}

	content, err := dbs.ArticleContentById(ids[0])
	if err != nil {
		fmt.Println("encountered DB error getting article content")
	}
	count := stats.WordCount(content)

	fmt.Printf("Word Count for article: %v\n", count)

	return nil
}
