package stats

type ArticleLengthStat struct {
	Length    int
	WordCount int
}

type ArticleLengthStatResult struct {
	Stats   []ArticleLengthStat
	Average float64
}

// We never check if Stats is nil in any of these.
// Feels like my Java days are back.

func (alsr *ArticleLengthStatResult) PushStat(s *ArticleLengthStat) {
	alsr.Stats = append(alsr.Stats, *s)
}

func (alsr *ArticleLengthStatResult) ComputeAverage() {
	sum := 0
	for _, wc := range alsr.Stats {
		sum += wc.WordCount
	}
	// A non zero sum means we got at least 1 result
	// So no divide by 0 is possible
	if sum > 0 {
		alsr.Average = float64(sum) / float64(len(alsr.Stats))
	}
}

func NewArticleLengthStatResult() *ArticleLengthStatResult {
	return &ArticleLengthStatResult{
		Stats: make([]ArticleLengthStat, 0),
	}
}
