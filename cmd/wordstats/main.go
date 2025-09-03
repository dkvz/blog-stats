package main

import (
	"flag"
	"fmt"
	"net/http"
	"slices"
	"sort"

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
	cliArgs, err := cli.ParseCliArgs()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
	}

	dbs, err := db.NewDBSqlite(config.DbPath)
	if err != nil {
		fmt.Println("encountered DB error")
		panic(err)
	}

	if cliArgs.Mode == 0 {
		runModePlot(dbs, cliArgs.IgnoredIds)
	}
}

func runModePlot(dbs *db.DbSqlite, ignoredIds []uint) {
	ids, err := dbs.AllPublishedArticleIds()
	if err != nil {
		fmt.Println("encountered DB error getting article ids")
		panic(err)
	}

	// Not sure if deleting from a slice is faster than just
	// creating a new one. My intuition is that it's not.
	// If it was important I'd check with a benchmark.
	if len(ignoredIds) > 0 {
		var newIds []uint
		for _, i := range ids {
			if !slices.Contains(ignoredIds, i) {
				newIds = append(newIds, i)
			}
		}
		ids = newIds
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

	// Compute word count stats:
	wcStats := make([]float64, len(results.Stats))
	for i, r := range results.Stats {
		wcStats[i] = float64(r.WordCount())
	}
	wcStatsC := stats.ComputeStats(wcStats)
	fmt.Printf("\nWord count stats:\n%s\n\n", wcStatsC)

	// Compute length stats:
	lengthStats := make([]float64, len(results.Stats))
	for i, r := range results.Stats {
		lengthStats[i] = float64(r.Length())
	}
	lengthStatsC := stats.ComputeStats(lengthStats)
	fmt.Printf("\nArticle length stats:\n%s\n\n", lengthStatsC)

	// Create a slice with the ratios to compute stats:
	ratios := make([]float64, len(results.Stats))
	for i, r := range results.Stats {
		ratios[i] = r.WordsPerCharRatio()
	}
	ratioStats := stats.ComputeStats(ratios)
	fmt.Printf("\nRatio stats:\n%s\n\n", ratioStats)

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

	// Create the map to create the plot
	lengthToRatio := make(map[float64]float64, len(ratios))
	lengthToWC := make(map[float64]float64, len(wcStats))
	for i, l := range lengthStats {
		lengthToRatio[l] = ratios[i]
		lengthToWC[l] = wcStats[i]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lToR := stats.GenerateScatterPlot(lengthToRatio, "Word Count per Length Ratio as function of Length")
		lToWC := stats.GenerateScatterPlot(lengthToWC, "Word Count per length")
		stats.GeneratePlotPage(w, lToR, lToWC)
	})
	fmt.Println("\nStarted HTTP server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
