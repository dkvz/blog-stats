package stats

type ArticleLengthStat struct {
	Length    int
	WordCount int
}

type ArticleLengthStatResult struct {
	stats   []ArticleLengthStat
	Average float64
}

// TODO: Implement a push method for the results
