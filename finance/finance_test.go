package finance

import (
	"math"
	"testing"
)

func TestFinanceReturnsAndRisk(t *testing.T) {
	prices := []float64{100, 110, 99, 120}

	assertClose(t, Return(prices), 0.2)
	assertClose(t, LogReturn(prices), math.Log(1.2))
	assertSliceClose(t, SimpleReturns(prices), []float64{0.1, -0.1, 120.0/99.0 - 1})
	assertSliceClose(t, LogReturns(prices), []float64{math.Log(1.1), math.Log(0.9), math.Log(120.0 / 99.0)})
	assertClose(t, MaxDrawdown(prices), 0.1)
	assertClose(t, Volatility([]float64{0.1, 0.2, 0.3}), math.Sqrt(0.02/3))
	assertClose(t, Sharpe([]float64{0.1, 0.2, 0.3}, 0.1), (0.2-0.1)/math.Sqrt(0.02/3))
	assertClose(t, Calmar([]float64{100, 150, 120, 200}), 1.0/0.2)
}

func TestFinanceEmptySingleAndZero(t *testing.T) {
	assertNaN(t, Return(nil))
	assertNaN(t, Return([]float64{0, 1}))
	assertNaN(t, LogReturn([]float64{1, 0}))
	assertClose(t, MaxDrawdown([]float64{1}), 0)
	assertNaN(t, Volatility(nil))
	assertNaN(t, Sharpe([]float64{1, 1}, 0))
	assertNaN(t, Calmar([]float64{1, 2}))

	simple := SimpleReturns([]float64{0, 1})
	assertNaN(t, simple[0])
	logRet := LogReturns([]float64{1, 0})
	assertNaN(t, logRet[0])
}

func TestMovingIndicators(t *testing.T) {
	x := []float64{1, 2, 3, 4}

	assertSliceClose(t, EMA(x, 3), []float64{1, 1.5, 2.25, 3.125})
	assertSliceClose(t, DMA(x, 0.5), []float64{1, 1.5, 2.25, 3.125})

	sma := SMA(x, 2)
	assertNaN(t, sma[0])
	assertClose(t, sma[1], 1.5)
	assertClose(t, sma[3], 3.5)

	upper, middle, lower := BollingerBands(x, 2, 2)
	assertNaN(t, upper[0])
	assertNaN(t, middle[0])
	assertNaN(t, lower[0])
	assertClose(t, upper[1], 2.5)
	assertClose(t, middle[1], 1.5)
	assertClose(t, lower[1], 0.5)
}

func assertSliceClose(t *testing.T, got, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got len %d, want %d", len(got), len(want))
	}
	for i := range got {
		assertClose(t, got[i], want[i])
	}
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.IsNaN(got) || math.Abs(got-want) > 1e-9 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func assertNaN(t *testing.T, got float64) {
	t.Helper()
	if !math.IsNaN(got) {
		t.Fatalf("got %v, want NaN", got)
	}
}
