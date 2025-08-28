package runtime

import (
	"fmt"
	"math"
	"runtime"

	"github.com/dkvz/blog-stats/pkg/db"
	"github.com/dkvz/blog-stats/pkg/stats"
)

func LengthStatsForIds(ids []uint, dbs *db.DbSqlite) (*stats.ArticleLengthStatResult, error) {
	// TODO: Having one item per routine is certainly ineffective and
	// we should lower the thread count further than that.
	routinesCount := min(len(ids), runtime.NumCPU())

	resChan := make(chan *stats.ArticleLengthStatResult)
	errChan := make(chan error)

	itemsCount := len(ids)
	n := int(math.Ceil(float64(itemsCount) / float64(routinesCount)))
	remainingRoutines := routinesCount

	for sliceAt := 0; remainingRoutines > 0; sliceAt += n {
		// Start the goroutines
		fmt.Printf("Routine starting at %v up to %v included\n", sliceAt, sliceAt+n)
		go lengthStatsForSlice(ids[sliceAt:sliceAt+n], dbs, resChan, errChan)
		remainingItems := itemsCount - sliceAt
		if remainingRoutines <= remainingItems {
			// Change how many items we give to routines:
			n = int(math.Ceil(float64(remainingItems) / float64(remainingRoutines)))
		}

		remainingRoutines--
	}

	// Process results
	// Might get stuck if one routine gets stuck as well
	final := stats.NewArticleLengthStatResult()
	for len(final.Stats) < len(ids) {
		select {
		case result := <-resChan:
			fmt.Printf("Result in: %v items; avg %v\n", len(result.Stats), result.Average)
			final.Stats = append(final.Stats, result.Stats...)
		case err := <-errChan:
			return nil, err
		}
	}

	// Compute stats:
	final.ComputeAverage()

	return final, nil
}

// Meant to be used as a goroutine
func lengthStatsForSlice(
	ids []uint,
	dbs *db.DbSqlite,
	resChan chan<- *stats.ArticleLengthStatResult,
	errChan chan<- error,
) {
	res := stats.NewArticleLengthStatResult()

	for _, id := range ids {
		content, err := dbs.ArticleContentById(id)
		if err != nil {
			errChan <- err
			continue
		}

		// Compute word count:
		stat := &stats.ArticleLengthStat{
			WordCount: stats.WordCount(content),
			Length:    len(*content),
		}
		res.PushStat(stat)
	}

	res.ComputeAverage()

	resChan <- res
}
