package agg

// Aggregator exposes grouped numeric statistics.
type Aggregator interface {
	Count() map[string]float64
	Sum() map[string]float64
	Mean() map[string]float64
	Median() map[string]float64
	Max() map[string]float64
	Min() map[string]float64
	First() map[string]float64
	Last() map[string]float64
	Var() map[string]float64
	Std() map[string]float64
	Range() map[string]float64
	Quantile(q float64) map[string]float64
	Diff() map[string][]float64
	SumAbsDiff() map[string]float64
	AbsDiffSum() map[string]float64
	SumSquaredDiff() map[string]float64
	Gradient() map[string]float64
	FlipRate() map[string]float64
	Change() map[string]float64           // last - first
	PercentChange() map[string]float64    // (last - first) / first
	NormalizedChange() map[string]float64 // (last - first) / mean
	Return() map[string]float64           // same as PercentChange, 或者你明确成 simple return
	LogReturn() map[string]float64
	MaxDrawdown() map[string]float64
}
