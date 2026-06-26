// Package finance provides numeric indicators with common financial semantics.
package finance

import (
	"math"

	"github.com/simonks2016/efi_math/core"
)

// Return returns (last-first)/first for a value sequence.
func Return(x []float64) float64 {
	if len(x) < 2 || x[0] == 0 {
		return math.NaN()
	}
	return (x[len(x)-1] - x[0]) / x[0]
}

// LogReturn returns ln(last/first) for a value sequence.
func LogReturn(x []float64) float64 {
	if len(x) < 2 || x[0] <= 0 || x[len(x)-1] <= 0 {
		return math.NaN()
	}
	return math.Log(x[len(x)-1] / x[0])
}

// SimpleReturns returns price[i]/price[i-1]-1 for each adjacent pair.
func SimpleReturns(prices []float64) []float64 {
	if len(prices) < 2 {
		return []float64{}
	}
	out := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if prices[i-1] == 0 {
			out[i-1] = math.NaN()
			continue
		}
		out[i-1] = prices[i]/prices[i-1] - 1
	}
	return out
}

// LogReturns returns ln(price[i]/price[i-1]) for each adjacent pair.
func LogReturns(prices []float64) []float64 {
	if len(prices) < 2 {
		return []float64{}
	}
	out := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if prices[i-1] <= 0 || prices[i] <= 0 {
			out[i-1] = math.NaN()
			continue
		}
		out[i-1] = math.Log(prices[i] / prices[i-1])
	}
	return out
}

// MaxDrawdown returns the maximum drawdown as a positive ratio.
func MaxDrawdown(equity []float64) float64 {
	if len(equity) == 0 {
		return math.NaN()
	}
	peak := equity[0]
	var maxDD float64
	for _, v := range equity {
		if v > peak {
			peak = v
		}
		if peak == 0 {
			continue
		}
		dd := (peak - v) / peak
		if dd > maxDD {
			maxDD = dd
		}
	}
	return maxDD
}

// Volatility returns the population standard deviation of returns.
func Volatility(returns []float64) float64 {
	return core.Std(returns)
}

// Sharpe returns (mean(returns)-riskFreeRate)/std(returns).
func Sharpe(returns []float64, riskFreeRate float64) float64 {
	mean := core.Mean(returns)
	std := core.Std(returns)
	if math.IsNaN(mean) || math.IsNaN(std) || std == 0 {
		return math.NaN()
	}
	return (mean - riskFreeRate) / std
}

// Calmar returns Return(equity)/MaxDrawdown(equity).
func Calmar(equity []float64) float64 {
	ret := Return(equity)
	dd := MaxDrawdown(equity)
	if math.IsNaN(ret) || math.IsNaN(dd) || dd == 0 {
		return math.NaN()
	}
	return ret / dd
}

// EMA returns the exponential moving average using alpha=2/(period+1).
func EMA(x []float64, period int) []float64 {
	out := make([]float64, len(x))
	if len(x) == 0 {
		return out
	}
	if period <= 0 {
		fillNaN(out)
		return out
	}
	alpha := 2.0 / float64(period+1)
	return dmaWithAlpha(x, alpha)
}

// DMA returns a dynamic moving average using the caller-provided alpha.
func DMA(x []float64, alpha float64) []float64 {
	out := make([]float64, len(x))
	if len(x) == 0 {
		return out
	}
	if alpha < 0 || alpha > 1 || math.IsNaN(alpha) {
		fillNaN(out)
		return out
	}
	return dmaWithAlpha(x, alpha)
}

// SMA returns the simple moving average for each full rolling window.
func SMA(x []float64, period int) []float64 {
	out := make([]float64, len(x))
	fillNaN(out)
	if len(x) == 0 || period <= 0 {
		return out
	}
	var sum float64
	for i, v := range x {
		sum += v
		if i >= period {
			sum -= x[i-period]
		}
		if i >= period-1 {
			out[i] = sum / float64(period)
		}
	}
	return out
}

// BollingerBands returns upper, middle, and lower rolling bands.
func BollingerBands(x []float64, period int, k float64) (upper []float64, middle []float64, lower []float64) {
	upper = make([]float64, len(x))
	lower = make([]float64, len(x))
	fillNaN(upper)
	fillNaN(lower)
	middle = SMA(x, period)
	if len(x) == 0 || period <= 0 {
		return upper, middle, lower
	}
	for i := period - 1; i < len(x); i++ {
		window := x[i-period+1 : i+1]
		std := core.Std(window)
		upper[i] = middle[i] + k*std
		lower[i] = middle[i] - k*std
	}
	return upper, middle, lower
}

func dmaWithAlpha(x []float64, alpha float64) []float64 {
	out := make([]float64, len(x))
	if len(x) == 0 {
		return out
	}
	out[0] = x[0]
	for i := 1; i < len(x); i++ {
		out[i] = alpha*x[i] + (1-alpha)*out[i-1]
	}
	return out
}

func fillNaN(x []float64) {
	for i := range x {
		x[i] = math.NaN()
	}
}
