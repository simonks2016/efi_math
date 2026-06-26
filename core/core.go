// Package core provides order-independent numeric statistics.
package core

import (
	"math"
	"sort"
)

// Count returns the number of values in x.
func Count(x []float64) int {
	return len(x)
}

// Sum returns the sum of x. It returns NaN for an empty slice.
func Sum(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	var sum float64
	for _, v := range x {
		sum += v
	}
	return sum
}

// Mean returns the arithmetic mean of x. It returns NaN for an empty slice.
func Mean(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	return Sum(x) / float64(len(x))
}

// Max returns the maximum value in x. It returns NaN for an empty slice.
func Max(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	max := x[0]
	for _, v := range x[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the minimum value in x. It returns NaN for an empty slice.
func Min(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	min := x[0]
	for _, v := range x[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Range returns Max(x)-Min(x). It returns NaN for an empty slice.
func Range(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	return Max(x) - Min(x)
}

// Var returns the population variance of x. It returns NaN for an empty slice.
func Var(x []float64) float64 {
	if len(x) == 0 {
		return math.NaN()
	}
	mean := Mean(x)
	var sum float64
	for _, v := range x {
		d := v - mean
		sum += d * d
	}
	return sum / float64(len(x))
}

// Std returns the population standard deviation of x.
func Std(x []float64) float64 {
	v := Var(x)
	if math.IsNaN(v) {
		return math.NaN()
	}
	return math.Sqrt(v)
}

// Median returns the median value in x. It returns NaN for an empty slice.
func Median(x []float64) float64 {
	return Quantile(x, 0.5)
}

// Quantile returns the q quantile of x using linear interpolation.
func Quantile(x []float64, q float64) float64 {
	if len(x) == 0 || q < 0 || q > 1 || math.IsNaN(q) {
		return math.NaN()
	}
	values := append([]float64(nil), x...)
	sort.Float64s(values)

	if len(values) == 1 {
		return values[0]
	}

	pos := q * float64(len(values)-1)
	lo := int(math.Floor(pos))
	hi := int(math.Ceil(pos))
	if lo == hi {
		return values[lo]
	}
	weight := pos - float64(lo)
	return values[lo]*(1-weight) + values[hi]*weight
}
