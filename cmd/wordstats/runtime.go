package main

import (
	"runtime"

	"github.com/dkvz/blog-stats/pkg/db"
	"github.com/dkvz/blog-stats/pkg/stats"
)

// resChan chan<- stats.ArticleLengthStatResult,
// errChan chan<- error,

func LengthStatsForIds(ids []uint, dbs *db.DbSqlite) (*stats.ArticleLengthStatResult, error) {
	routinesCount := runtime.NumCPU()
	if len(ids) < routinesCount {
		// TODO: That level of multithreading is probably ineffective and
		// we should lower the thread count further.
		routinesCount = len(ids)
	}

	resChan := make(chan stats.ArticleLengthStatResult)
	errChan := make(chan error)

	for i := 0; i < routinesCount; i++ {
		// Start the goroutines
	}

	return nil, nil
}

// Meant to be used as a goroutine
func lengthStatsForSlice(
	ids []uint,
	dbs *db.DbSqlite,
	resChan chan<- stats.ArticleLengthStatResult,
	errChan chan<- error,
) {
	res := stats.ArticleLengthStatResult{}

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

	// Compute local average:
	sum := 0
	for _, wc := range res.Stats {
		sum += wc.WordCount
	}
	// A non zero sum means we got at least 1 result
	// So no divide by 0 is possible
	if sum > 0 {
		res.Average = float64(sum) / float64(len(res.Stats))
	}

	resChan <- res
}
