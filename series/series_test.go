package series

import (
	"math"
	"reflect"
	"testing"
)

func TestSeriesStats(t *testing.T) {
	x := []float64{0.1, 0.2, 0.1}

	assertClose(t, First(x), 0.1)
	assertClose(t, Last(x), 0.1)
	assertSliceClose(t, Diff(x), []float64{0.1, -0.1})
	assertSliceClose(t, AbsDiff(x), []float64{0.1, 0.1})
	assertSliceClose(t, SquaredDiff(x), []float64{0.01, 0.01})
	assertClose(t, SumAbsDiff(x), 0.2)
	assertClose(t, SumSquaredDiff(x), 0.02)
	assertClose(t, Change(x), 0)
	assertClose(t, PercentChange(x), 0)
	assertClose(t, NormalizedChange(x), 0)
	assertClose(t, Gradient(x), 0)
	assertClose(t, FlipRate(x), 1)
}

func TestSeriesEmptySingleAndZero(t *testing.T) {
	assertNaN(t, First(nil))
	assertNaN(t, Last(nil))
	assertTrue(t, reflect.DeepEqual(Diff(nil), []float64{}))
	assertClose(t, SumAbsDiff([]float64{1}), 0)
	assertClose(t, SumSquaredDiff([]float64{1}), 0)
	assertClose(t, Change([]float64{1}), 0)
	assertNaN(t, Gradient([]float64{1}))
	assertNaN(t, PercentChange([]float64{0, 1}))
	assertNaN(t, NormalizedChange([]float64{-1, 1}))
	assertClose(t, FlipRate([]float64{1, 1, 2, 2, 1}), 1)
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

func assertTrue(t *testing.T, ok bool) {
	t.Helper()
	if !ok {
		t.Fatal("condition is false")
	}
}
