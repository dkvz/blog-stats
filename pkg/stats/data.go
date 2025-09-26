package stats

import (
	"fmt"
	"math"
)

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
	// Divides the distance by the real word count
	// to show how they relate as some measure of
	// error, lower values are better
	distanceRelativeToWordCount float64
}

func NewArticleLengthPrediction(
	s *ArticleLengthStat,
	predictedWC int,
) *ArticleLengthPrediction {
	var wc int
	if s != nil {
		wc = s.WordCount()
	}

	distance := predictedWC - wc
	relDist := math.Abs(float64(distance) / float64(s.WordCount()))
	return &ArticleLengthPrediction{
		ArticleLengthStat:   *s,
		predictedWordCount:  predictedWC,
		distanceToWordCount: distance,
		// Pretty sure this isn't a percentage I always mix these up
		distanceRelativeToWordCount: relDist,
	}
}

func (a *ArticleLengthPrediction) PredictedWordCount() int {
	return a.predictedWordCount
}

func (a *ArticleLengthPrediction) DistanceToWordCount() int {
	return a.distanceToWordCount
}

func (a *ArticleLengthPrediction) DistanceToWordCountSquared() float64 {
	return float64(a.distanceToWordCount) * float64(a.distanceToWordCount)
}

func (a *ArticleLengthPrediction) DistanceRelativeToWordCount() float64 {
	return a.distanceRelativeToWordCount
}

type ArticleLengthStatResult struct {
	Stats []ArticleLengthStat
}

type SliceAnalytics struct {
	Min      float64
	Max      float64
	StdDev   float64
	Variance float64
	Median   float64
	Average  float64
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
