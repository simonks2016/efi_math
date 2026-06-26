package core

import (
	"math"
	"testing"
)

func TestCoreStats(t *testing.T) {
	x := []float64{1, 2, 3, 4}

	assertEqual(t, Count(x), 4)
	assertClose(t, Sum(x), 10)
	assertClose(t, Mean(x), 2.5)
	assertClose(t, Max(x), 4)
	assertClose(t, Min(x), 1)
	assertClose(t, Range(x), 3)
	assertClose(t, Var(x), 1.25)
	assertClose(t, Std(x), math.Sqrt(1.25))
	assertClose(t, Median(x), 2.5)
	assertClose(t, Quantile(x, 0.25), 1.75)
}

func TestCoreEmptyAndSingle(t *testing.T) {
	assertEqual(t, Count(nil), 0)
	assertNaN(t, Sum(nil))
	assertNaN(t, Mean(nil))
	assertNaN(t, Max(nil))
	assertNaN(t, Min(nil))
	assertNaN(t, Range(nil))
	assertNaN(t, Var(nil))
	assertNaN(t, Std(nil))
	assertNaN(t, Median(nil))
	assertNaN(t, Quantile(nil, 0.5))
	assertNaN(t, Quantile([]float64{1}, 1.5))

	x := []float64{5}
	assertClose(t, Sum(x), 5)
	assertClose(t, Var(x), 0)
	assertClose(t, Std(x), 0)
	assertClose(t, Median(x), 5)
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

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
