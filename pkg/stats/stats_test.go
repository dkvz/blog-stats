package stats

import "testing"

func TestSliceAnalytics(t *testing.T) {
	testSlice := []float64{2.3, 55.8, 10203.734, 22.0, 550.0, 280.3}

	expected := &SliceAnalytics{
		Average: 1852.3556666666666,
		Min:     2.3,
		Max:     10203.734,
		StdDev:  3739.731039075835,
		Median:  168.05,
	}

	stats := ComputeStats(testSlice)
	if stats.Average != expected.Average {
		t.Errorf("TestSliceAnalytics: wrong average %v instead of %v", stats.Average, expected.Average)
	}
	if stats.Min != expected.Min {
		t.Errorf("TestSliceAnalytics: wrong min %v instead of %v", stats.Min, expected.Min)
	}
	if stats.Max != expected.Max {
		t.Errorf("TestSliceAnalytics: wrong max %v instead of %v", stats.Max, expected.Max)
	}
	if stats.Median != expected.Median {
		t.Errorf("TestSliceAnalytics: wrong median %v instead of %v", stats.Median, expected.Median)
	}
	if stats.StdDev != expected.StdDev {
		t.Errorf("TestSliceAnalytics: wrong std dev %v instead of %v", stats.StdDev, expected.StdDev)
	}
}
