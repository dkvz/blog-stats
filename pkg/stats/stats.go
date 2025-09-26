package stats

import (
	"math"
	"sort"
)

// Computes various stats values from the give slice
// Also returns a new slice with values sorted in ascending order
func ComputeStats(data []float64) *SliceAnalytics {
	dataLen := len(data)
	if dataLen == 0 {
		return &SliceAnalytics{}
	}

	dataCopy := make([]float64, dataLen)
	copy(dataCopy, data)
	// slices.sort does the same thing I think
	sort.Float64s(dataCopy)

	avg := ComputeAverage(dataCopy)
	// Compute std dev
	var variance float64
	for _, v := range data {
		diff := v - avg
		variance += diff * diff
	}
	variance /= float64(dataLen)
	stdDev := math.Sqrt(variance)

	// remaining: median, min and max
	var median float64
	if dataLen%2 == 0 {
		// Even number of elements: average of two middle values
		median = (dataCopy[(dataLen/2)-1] + dataCopy[dataLen/2]) / 2.0
	} else {
		// Odd number of elements: middle value
		median = dataCopy[dataLen/2]
	}

	ret := &SliceAnalytics{
		Min:      dataCopy[0],
		Max:      dataCopy[dataLen-1],
		Median:   median,
		Average:  avg,
		StdDev:   stdDev,
		Variance: variance,
	}

	return ret
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

func ComputePredictionSpread(predictions []ArticleLengthPrediction) float64 {
	spread := 0.0
	for _, p := range predictions {
		spread += p.DistanceToWordCountSquared()
	}
	if spread != 0.0 {
		// Indicates we don't have 0 items and thus a
		// division by 0
		spread /= float64(len(predictions))
	}
	return spread
}

// Slices x and y have to be the same size or things will go wrong
func ComputeLinearRegForcedOrigin(x []float64, y []float64) float64 {
	var xySum, x2Sum float64
	for i, xv := range x {
		yv := y[i]
		xySum += xv * yv
		x2Sum += xv * xv
	}

	if x2Sum == 0 {
		return 0
	}

	return xySum / x2Sum
}

// Computes normal linear regression for the series
// Returns beta, alpha
// Where: y = alpha + beta * x
//
// Slices x and y have to be the same size or things will go wrong
func ComputeLinearReg(
	x []float64,
	y []float64,
	varianceX float64,
	averageX float64,
	averageY float64,
) (float64, float64) {
	if len(x)-1 == 0 || varianceX == 0 {
		// Would cause divide by 0 errors
		return 0, 0
	}

	// Compute the covariance:
	var ss, xcompensation, ycompensation float64

	for i, xv := range x {
		yv := y[i]
		xd := xv - averageX
		yd := yv - averageY
		ss += xd * yd
		xcompensation += xd
		ycompensation += yd
	}
	cov := (ss - xcompensation*ycompensation/float64(len(x))) / float64(len(x)-1)

	beta := cov / varianceX
	alpha := averageY - (beta * averageX)

	return beta, alpha
}
