package runtime

import (
	"math"
	"runtime"
	"slices"

	"github.com/dkvz/blog-stats/pkg/cli"
	"github.com/dkvz/blog-stats/pkg/db"
	"github.com/dkvz/blog-stats/pkg/stats"
)

func LengthStatsForIds(ids []uint, dbs *db.DbSqlite, cliArgs *cli.CliArgs) (*stats.ArticleLengthStatResult, error) {
	// TODO: Having one item per routine is certainly ineffective and
	// we should lower the thread count further than that.
	routinesCount := min(len(ids), runtime.NumCPU())

	resChan := make(chan *stats.ArticleLengthStatResult)
	errChan := make(chan error)

	// Not sure if deleting from a slice is faster than just
	// creating a new one. My intuition is that it's not.
	// If it was important I'd check with a benchmark.
	if cliArgs != nil && len(cliArgs.IgnoredIds) > 0 {
		var newIds []uint
		for _, i := range ids {
			if !slices.Contains(cliArgs.IgnoredIds, i) {
				newIds = append(newIds, i)
			}
		}
		ids = newIds
	}

	itemsCount := len(ids)
	n := int(math.Ceil(float64(itemsCount) / float64(routinesCount)))
	remainingRoutines := routinesCount
	intSLength := int(cliArgs.StartLength)
	intELength := int(cliArgs.EndLength)

	for sliceAt := 0; remainingRoutines > 0; sliceAt += n {
		// Start the goroutines
		go lengthStatsForSlice(
			ids[sliceAt:sliceAt+n],
			dbs,
			int(intSLength),
			int(intELength),
			resChan,
			errChan,
		)
		remainingItems := itemsCount - sliceAt
		if remainingRoutines <= remainingItems {
			// Change how many items we give to routines:
			n = int(math.Ceil(float64(remainingItems) / float64(remainingRoutines)))
		}

		remainingRoutines--
	}

	// Process results
	// Might get stuck if one routine gets stuck as well
	final := &stats.ArticleLengthStatResult{}
	// for len(final.Stats) < len(ids) {
	for remainingRoutines < routinesCount {
		select {
		case result := <-resChan:
			remainingRoutines++
			final.Stats = append(final.Stats, result.Stats...)
		case err := <-errChan:
			return nil, err
		}
	}

	return final, nil
}

// Meant to be used as a goroutine
// We do not check the validity of startLength and
// endLength at all in this function
func lengthStatsForSlice(
	ids []uint,
	dbs *db.DbSqlite,
	startLength int,
	endLength int,
	resChan chan<- *stats.ArticleLengthStatResult,
	errChan chan<- error,
) {
	res := &stats.ArticleLengthStatResult{}

	for _, id := range ids {
		content, err := dbs.ArticleContentById(id)
		if err != nil {
			errChan <- err
			continue
		}

		// The length to use has to be calculated after UTF-16
		// conversion because JS uses UTF-16 and the factors are
		// to be ultimately used in JS.
		contentLength := stats.LengthUTF16(content)
		if contentLength >= startLength {
			if endLength <= 0 || (endLength > 0 && contentLength < endLength) {
				stat := stats.NewArticleLengthStat(id, contentLength, stats.WordCount(content))
				res.PushStat(stat)
			}
		}
	}

	// Thought of calculating local averages here to merge
	// them later but I need other averages and I also need
	// those for the std dev caculation.

	resChan <- res
}
