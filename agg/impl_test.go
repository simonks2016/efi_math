package agg

import (
	"math"
	"reflect"
	"testing"
	"time"
)

type sample struct {
	symbol string
	value  float64
	ts     time.Time
}

func TestAggregatorGroupsAndStats(t *testing.T) {
	items := []sample{
		{symbol: "A", value: 1},
		{symbol: "B", value: 10},
		{symbol: "A", value: 3},
		{symbol: "A", value: 5},
		{symbol: "B", value: 20},
	}

	a := NewAggregator(items, func(x sample) string { return x.symbol }, func(x sample) float64 { return x.value })

	assertClose(t, a.Count()["A"], 3)
	assertClose(t, a.Sum()["A"], 9)
	assertClose(t, a.Mean()["A"], 3)
	assertClose(t, a.Max()["B"], 20)
	assertClose(t, a.Min()["B"], 10)
	assertClose(t, a.Var()["A"], 8.0/3.0)
	assertClose(t, a.Std()["A"], math.Sqrt(8.0/3.0))
	assertClose(t, a.Median()["A"], 3)
	assertClose(t, a.Quantile(0.5)["B"], 15)
	assertSliceClose(t, a.Diff()["A"], []float64{2, 2})
	assertClose(t, a.SumAbsDiff()["A"], 4)
	assertClose(t, a.AbsDiffSum()["A"], 4)
	assertClose(t, a.SumSquaredDiff()["A"], 8)
	assertClose(t, a.Change()["A"], 4)
	assertClose(t, a.PercentChange()["A"], 4)
	assertClose(t, a.NormalizedChange()["A"], 4.0/3.0)
	assertClose(t, a.Gradient()["A"], 2)
	assertClose(t, a.Return()["A"], 4)
	assertClose(t, a.LogReturn()["A"], math.Log(5))
}

func TestAggregatorWithTimeSortsSequenceMetrics(t *testing.T) {
	base := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)
	items := []sample{
		{symbol: "A", value: 0.2, ts: base.Add(time.Second)},
		{symbol: "A", value: 0.1, ts: base.Add(2 * time.Second)},
		{symbol: "A", value: 0.1, ts: base},
		{symbol: "B", value: 100, ts: base},
		{symbol: "B", value: 80, ts: base.Add(time.Second)},
		{symbol: "B", value: 120, ts: base.Add(2 * time.Second)},
	}

	a := NewAggregator(
		items,
		func(x sample) string { return x.symbol },
		func(x sample) float64 { return x.value },
		WithTime(func(x sample) time.Time { return x.ts }),
	)

	assertSliceClose(t, a.Diff()["A"], []float64{0.1, -0.1})
	assertClose(t, a.SumSquaredDiff()["A"], 0.02)
	assertClose(t, a.SumAbsDiff()["A"], 0.2)
	assertClose(t, a.Change()["A"], 0)
	assertClose(t, a.FlipRate()["A"], 1)
	assertClose(t, a.MaxDrawdown()["B"], 0.2)
}

func TestAggregatorWithoutTimeKeepsInputOrder(t *testing.T) {
	items := []sample{
		{symbol: "A", value: 2},
		{symbol: "A", value: 1},
		{symbol: "A", value: 3},
	}

	a := New(items, func(x sample) string { return x.symbol }, func(x sample) float64 { return x.value })

	if !reflect.DeepEqual(a.Diff()["A"], []float64{-1, 2}) {
		t.Fatalf("got %v", a.Diff()["A"])
	}
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
