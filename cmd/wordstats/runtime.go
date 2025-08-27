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

	for _, id := range ids {
		content, err := dbs.ArticleContentById(id)
		if err != nil {
			errChan <- err
			continue
		}

		// Compute word count:

	}
}
