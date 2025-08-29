package stats

type ArticleLengthStat struct {
	ArticleId         uint
	length            int
	wordCount         int
	wordsPerCharRatio float64
}

type ArticleLengthStatResult struct {
	Stats   []ArticleLengthStat
	Average float64
}

type ArticleLengthAnalytics struct {
	Min    float64
	Max    float64
	StdDev float64
	Median float64
}

// We never check if Stats is nil in any of these.
// Feels like my Java days are back.

func (alsr *ArticleLengthStatResult) PushStat(s *ArticleLengthStat) {
	alsr.Stats = append(alsr.Stats, *s)
}

func (alsr *ArticleLengthStatResult) ComputeAverage() {
	sum := 0
	for _, wc := range alsr.Stats {
		sum += wc.WordCount()
	}
	// A non zero sum means we got at least 1 result
	// So no divide by 0 is possible
	if sum > 0 {
		alsr.Average = float64(sum) / float64(len(alsr.Stats))
	}
}

func NewArticleLengthStat(
	articleId uint,
	length int,
	wordCount int,
) *ArticleLengthStat {
	return &ArticleLengthStat{
		ArticleId:         articleId,
		length:            length,
		wordCount:         wordCount,
		wordsPerCharRatio: float64(wordCount) / float64(length),
	}
}

func (als *ArticleLengthStat) WordsPerCharRatio() float64 {
	return float64(als.WordCount()) / float64(als.Length())
}

func (als *ArticleLengthStat) Length() int {
	return als.length
}

func (als *ArticleLengthStat) WordCount() int {
	return als.wordCount
}
