// Package series provides order-aware numeric sequence statistics.
package series

import (
	"math"

	"github.com/simonks2016/efi_math/core"
)

// First returns the first value in x. It returns NaN for an empty slice.
func First(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	return x[0]
}

// Last returns the last value in x. It returns NaN for an empty slice.
func Last(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	return x[len(x)-1]
}

// Diff returns adjacent differences x[i]-x[i-1].
func Diff(x []float64) []float64 {
	if len(x) < 2 {
		return []float64{}
	}
	out := make([]float64, len(x)-1)
	for i := 1; i < len(x); i++ {
		out[i-1] = x[i] - x[i-1]
	}
	return out
}

// AbsDiff returns absolute adjacent differences.
func AbsDiff(x []float64) []float64 {
	diffs := Diff(x)
	out := make([]float64, len(diffs))
	for i, v := range diffs {
		out[i] = math.Abs(v)
	}
	return out
}

// SquaredDiff returns squared adjacent differences.
func SquaredDiff(x []float64) []float64 {
	diffs := Diff(x)
	out := make([]float64, len(diffs))
	for i, v := range diffs {
		out[i] = v * v
	}
	return out
}

// SumAbsDiff returns the sum of absolute adjacent differences.
func SumAbsDiff(x []float64) float64 {
	var sum float64
	for _, v := range AbsDiff(x) {
		sum += v
	}
	return sum
}

// SumSquaredDiff returns the sum of squared adjacent differences.
func SumSquaredDiff(x []float64) float64 {
	var sum float64
	for _, v := range SquaredDiff(x) {
		sum += v
	}
	return sum
}

// Change returns last-first. It returns NaN for an empty slice.
func Change(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	return Last(x) - First(x)
}

// PercentChange returns (last-first)/first.
func PercentChange(x []float64) float64 {
	if len(x) == 0 || x[0] == 0 {
		return math.NaN()
	}
	return Change(x) / x[0]
}

// NormalizedChange returns (last-first)/mean(x).
func NormalizedChange(x []float64) float64 {
	mean := core.Mean(x)
	if len(x) == 0 || mean == 0 || math.IsNaN(mean) {
		return math.NaN()
	}
	return Change(x) / mean
}

// Gradient returns (last-first)/float64(len(x)-1).
func Gradient(x []float64) float64 {
	if len(x) < 2 {
		return math.NaN()
	}
	return Change(x) / float64(len(x)-1)
}

// FlipRate returns the ratio of sign flips among adjacent non-zero differences.
func FlipRate(x []float64) float64 {
	diffs := Diff(x)
	var prev int
	var transitions float64
	var flips float64

	for _, d := range diffs {
		sign := 0
		switch {
		case d > 0:
			sign = 1
		case d < 0:
			sign = -1
		default:
			continue
		}

		if prev != 0 {
			transitions++
			if sign != prev {
				flips++
			}
		}
		prev = sign
	}

	if transitions == 0 {
		return 0
	}
	return flips / transitions
}
