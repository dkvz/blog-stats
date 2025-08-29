package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/dkvz/blog-stats/pkg/cli"
	"github.com/dkvz/blog-stats/pkg/db"
	"github.com/dkvz/blog-stats/pkg/runtime"
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

func runModePlot(dbs *db.DbSqlite) {
	ids, err := dbs.AllPublishedArticleIds()
	if err != nil {
		fmt.Println("encountered DB error getting article ids")
		panic(err)
	}

	results, err := runtime.LengthStatsForIds(ids, dbs)
	if err != nil {
		fmt.Println("error in the subroutines")
		panic(err)
	}

	// I'll be sorting things multiple times but that's fine
	sort.Slice(results.Stats, func(i, j int) bool {
		return results.Stats[i].WordsPerCharRatio() < results.Stats[j].WordsPerCharRatio()
	})

	// Create a slice with the ratios to compute stats:
	ratios := make([]float64, len(results.Stats))
	for i, r := range results.Stats {
		ratios[i] = r.WordsPerCharRatio()
	}
	stats := stats.ComputeStats(ratios)

	// TODO: Implement String() for SliceAnalytics

	fmt.Printf("\nID\tWC\tLength\tRatio\n")
	for _, wc := range results.Stats {
		fmt.Printf(
			"%v\t%v\t%v\t%v\n",
			wc.ArticleId,
			wc.WordCount(),
			wc.Length(),
			wc.WordsPerCharRatio(),
		)
	}

}
