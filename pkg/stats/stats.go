package stats

import "slices"

// Computes various stats values from the give slice
// Also returns a new slice with values sorted in ascending order
func ComputeStatsAndSort(data []float64) (*ArticleLengthAnalytics, []float64) {
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)
	slices.Sort(dataCopy)

	// Compute std dev

	return nil, nil
}

func ComputeAverage(data []float64) float64 {
	var sum float64
	for _, d := range data {
		sum += d
	}
	// A non zero sum means we got at least 1 result
	// So no divide by 0 is possible
	if sum > 0 {
		return float64(sum) / float64(len(data))
	}
	return 0
}
