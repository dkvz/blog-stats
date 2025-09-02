package runtime

import (
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
	final := &stats.ArticleLengthStatResult{}
	for len(final.Stats) < len(ids) {
		select {
		case result := <-resChan:
			final.Stats = append(final.Stats, result.Stats...)
		case err := <-errChan:
			return nil, err
		}
	}

	return final, nil
}

// Meant to be used as a goroutine
func lengthStatsForSlice(
	ids []uint,
	dbs *db.DbSqlite,
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
		stat := stats.NewArticleLengthStat(id, contentLength, stats.WordCount(content))
		res.PushStat(stat)
	}

	// Thought of calculating local averages here to merge
	// them later but I need other averages and I also need
	// those for the std dev caculation.

	resChan <- res
}
