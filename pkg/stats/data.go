package stats

import "fmt"

type ArticleLengthStat struct {
	ArticleId         uint
	length            int
	wordCount         int
	wordsPerCharRatio float64
}

type ArticleLengthPrediction struct {
	ArticleLengthStat
	predictedWordCount  int
	distanceToWordCount int
}

func NewArticleLengthPrediction(
	s *ArticleLengthStat,
	predictedWC int,
) *ArticleLengthPrediction {
	var wc int
	if s != nil {
		wc = s.WordCount()
	}

	return &ArticleLengthPrediction{
		ArticleLengthStat:   *s,
		predictedWordCount:  predictedWC,
		distanceToWordCount: predictedWC - wc,
	}
}

type ArticleLengthStatResult struct {
	Stats []ArticleLengthStat
}

type SliceAnalytics struct {
	Min     float64
	Max     float64
	StdDev  float64
	Median  float64
	Average float64
}

func (sl *SliceAnalytics) String() string {
	return fmt.Sprintf(
		"Avg: %f\tMin: %f\tMax: %f\nStdDev: %f\tMed: %f",
		sl.Average,
		sl.Min,
		sl.Max,
		sl.StdDev,
		sl.Median,
	)
}

// We never check if Stats is nil in any of these.
// Feels like my Java days are back.

func (alsr *ArticleLengthStatResult) PushStat(s *ArticleLengthStat) {
	alsr.Stats = append(alsr.Stats, *s)
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
