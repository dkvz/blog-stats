package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
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
		fmt.Println()
		flag.Usage()
		return
	}

	dbs, err := db.NewDBSqlite(config.DbPath)
	if err != nil {
		fmt.Println("encountered DB error")
		panic(err)
	}

	// For the moment both modes need the full length stats
	// computed:
	results := lengthStats(dbs, cliArgs)
	switch cliArgs.Mode {
	case 0:
		runModePlot(results)
	case 1:
		runVerifyMode(results, cliArgs.VerifyModeArgs)
	}
}

// Fetches and computes the length stats for the given params
// Might panic in case of error
func lengthStats(
	dbs *db.DbSqlite,
	cliArgs *cli.CliArgs,
) *stats.ArticleLengthStatResult {
	ids, err := dbs.AllPublishedArticleIds()
	if err != nil {
		fmt.Println("encountered DB error getting article ids")
		panic(err)
	}

	lengthArgs := &runtime.LengthStatsOpts{
		IgnoredIds:  cliArgs.IgnoredIds,
		StartLength: int(cliArgs.StartLength),
		EndLength:   int(cliArgs.EndLength),
	}

	results, err := runtime.LengthStatsForIds(ids, dbs, lengthArgs)
	if err != nil {
		fmt.Println("error in the subroutines")
		panic(err)
	}

	return results
}

func runVerifyMode(results *stats.ArticleLengthStatResult, args *cli.VerifyArgs) {
	var predicted []stats.ArticleLengthPrediction

	if args.RegMode {
		// Linear regression mode
		// I'm doing hazardous float64 to int conversions but whatever.
		fmt.Printf(
			"Predicting using linear regression y = %f x + %f\n",
			args.DefaultFactor,
			args.RegA,
		)
		for _, s := range results.Stats {
			// Run the prediction:
			pred := stats.NewArticleLengthPrediction(
				&s,
				int(math.Round(float64(s.Length())*args.DefaultFactor+args.RegA)),
			)
			predicted = append(predicted, *pred)
		}
	} else {
		// Use factors mode, which is a bit more complex
		// Use the first factor that matches.
		fmt.Printf(
			"Predicting using factors, default %f and %v extra factors\n",
			args.DefaultFactor,
			len(args.Factors),
		)
		for _, s := range results.Stats {
			predWC := -1.0
			if len(args.Factors) > 0 {
				for _, f := range args.Factors {
					if s.Length() >= int(f.Start) && s.Length() <= int(f.End) {
						predWC = float64(s.Length()) * f.Value
						break
					}
				}
			}

			// None of the factors matched:
			if predWC < 0 {
				predWC = float64(s.Length()) * args.DefaultFactor
			}

			pred := stats.NewArticleLengthPrediction(&s, int(predWC))
			predicted = append(predicted, *pred)
		}
	}

	// We can now sort the data according to distance
	// and compute the average spread.
	spread := stats.ComputePredictionSpread(predicted)
	// Create a slice to compute the average relative distance
	relDist := make([]float64, len(predicted))
	for i, p := range predicted {
		relDist[i] = p.DistanceRelativeToWordCount()
	}
	relDistAvg := stats.ComputeAverage(relDist)

	// Sort by distance squared to avoid the sign:
	sort.Slice(predicted, func(i, j int) bool {
		return predicted[i].DistanceRelativeToWordCount() > predicted[j].DistanceRelativeToWordCount()
	})

	fmt.Printf(
		"\nPredictions spread (std dev - lower is better): %v\n",
		math.Sqrt(spread),
	)

	fmt.Printf(
		"\nAverage for relative distances (lower is better): %v\n",
		relDistAvg,
	)

	fmt.Printf("\nID\tWC\tP.WC\tDist\tRel.dist\tLength\n")
	for _, p := range predicted {
		fmt.Printf(
			"%v\t%v\t%v\t%v\t%.4f\t%v\n",
			p.ArticleId,
			p.WordCount(),
			p.PredictedWordCount(),
			p.DistanceToWordCount(),
			p.DistanceRelativeToWordCount(),
			p.Length(),
		)
	}

}

func runModePlot(results *stats.ArticleLengthStatResult) {
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
