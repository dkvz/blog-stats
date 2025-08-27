package stats

type ArticleLengthStat struct {
	Length    int
	WordCount int
}

type ArticleLengthStatResult struct {
	Stats   []ArticleLengthStat
	Average float64
}
