package stats

type ArticleLengthStat struct {
	Length    int
	WordCount int
}

type ArticleLengthStatResult struct {
	Stats   []ArticleLengthStat
	Average float64
}

func (alsr *ArticleLengthStatResult) PushStat(s *ArticleLengthStat) {
	alsr.Stats = append(alsr.Stats, *s)
}
