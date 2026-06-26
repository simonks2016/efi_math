package agg

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/simonks2016/efi_math/core"
	"github.com/simonks2016/efi_math/finance"
	"github.com/simonks2016/efi_math/series"
)

type Impl[in any, K comparable] struct {
	dataMap   map[string][]float64
	mu        sync.RWMutex
	keyFunc   func(in) K
	valueFunc func(in) float64
	timeFn    func(in) time.Time
	sorted    bool
}

type groupItem struct {
	value float64
	time  time.Time
}

// Option configures an Aggregator implementation.
type Option[in any] func(*options[in])

type options[in any] struct {
	timeFn func(in) time.Time
}

// WithTime sorts each group by timestamp before sequence-aware calculations.
func WithTime[in any](timeFn func(in) time.Time) Option[in] {
	return func(o *options[in]) {
		o.timeFn = timeFn
	}
}

// NewAggregator builds a grouped Aggregator from arbitrary input records.
func NewAggregator[in any, K comparable](
	data []in,
	keyFunc func(in) K,
	valueFunc func(in) float64,
	opts ...Option[in],
) Aggregator {
	cfg := options[in]{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	groups := make(map[string][]groupItem)
	for _, item := range data {
		key := keyToString(keyFunc(item))
		groups[key] = append(groups[key], groupItem{
			value: valueFunc(item),
			time:  timeFrom(item, cfg.timeFn),
		})
	}

	dataMap := make(map[string][]float64, len(groups))
	for key, items := range groups {
		if cfg.timeFn != nil {
			sort.SliceStable(items, func(i, j int) bool {
				return items[i].time.Before(items[j].time)
			})
		}
		values := make([]float64, len(items))
		for i, item := range items {
			values[i] = item.value
		}
		dataMap[key] = values
	}

	return &Impl[in, K]{
		dataMap:   dataMap,
		keyFunc:   keyFunc,
		valueFunc: valueFunc,
		timeFn:    cfg.timeFn,
		sorted:    cfg.timeFn != nil,
	}
}

// New is an alias for NewAggregator.
func New[in any, K comparable](
	data []in,
	keyFunc func(in) K,
	valueFunc func(in) float64,
	opts ...Option[in],
) Aggregator {
	return NewAggregator(data, keyFunc, valueFunc, opts...)
}

// Count returns grouped element counts.
func (a *Impl[in, K]) Count() map[string]float64 {
	return a.applyScalar(func(x []float64) float64 {
		return float64(core.Count(x))
	})
}

// Sum returns grouped sums.
func (a *Impl[in, K]) Sum() map[string]float64 {
	return a.applyScalar(core.Sum)
}

// Mean returns grouped means.
func (a *Impl[in, K]) Mean() map[string]float64 {
	return a.applyScalar(core.Mean)
}

// Median returns grouped medians.
func (a *Impl[in, K]) Median() map[string]float64 {
	return a.applyScalar(core.Median)
}

// Max returns grouped maximum values.
func (a *Impl[in, K]) Max() map[string]float64 {
	return a.applyScalar(core.Max)
}

// Min returns grouped minimum values.
func (a *Impl[in, K]) Min() map[string]float64 {
	return a.applyScalar(core.Min)
}

// First returns grouped first values.
func (a *Impl[in, K]) First() map[string]float64 {
	return a.applyScalar(series.First)
}

// Last returns grouped last values.
func (a *Impl[in, K]) Last() map[string]float64 {
	return a.applyScalar(series.Last)
}

// Var returns grouped population variances.
func (a *Impl[in, K]) Var() map[string]float64 {
	return a.applyScalar(core.Var)
}

// Std returns grouped population standard deviations.
func (a *Impl[in, K]) Std() map[string]float64 {
	return a.applyScalar(core.Std)
}

// Range returns grouped max-min ranges.
func (a *Impl[in, K]) Range() map[string]float64 {
	return a.applyScalar(core.Range)
}

// Quantile returns grouped q quantiles.
func (a *Impl[in, K]) Quantile(q float64) map[string]float64 {
	return a.applyScalar(func(x []float64) float64 {
		return core.Quantile(x, q)
	})
}

// Diff returns grouped adjacent differences.
func (a *Impl[in, K]) Diff() map[string][]float64 {
	return a.applySlice(series.Diff)
}

// SumAbsDiff returns grouped sums of absolute adjacent differences.
func (a *Impl[in, K]) SumAbsDiff() map[string]float64 {
	return a.applyScalar(series.SumAbsDiff)
}

// AbsDiffSum returns grouped sums of absolute adjacent differences.
func (a *Impl[in, K]) AbsDiffSum() map[string]float64 {
	return a.SumAbsDiff()
}

// SumSquaredDiff returns grouped sums of squared adjacent differences.
func (a *Impl[in, K]) SumSquaredDiff() map[string]float64 {
	return a.applyScalar(series.SumSquaredDiff)
}

// Gradient returns grouped sequence gradients.
func (a *Impl[in, K]) Gradient() map[string]float64 {
	return a.applyScalar(series.Gradient)
}

// FlipRate returns grouped direction flip rates.
func (a *Impl[in, K]) FlipRate() map[string]float64 {
	return a.applyScalar(series.FlipRate)
}

// Change returns grouped last-first changes.
func (a *Impl[in, K]) Change() map[string]float64 {
	return a.applyScalar(series.Change)
}

// PercentChange returns grouped (last-first)/first changes.
func (a *Impl[in, K]) PercentChange() map[string]float64 {
	return a.applyScalar(series.PercentChange)
}

// NormalizedChange returns grouped changes normalized by mean.
func (a *Impl[in, K]) NormalizedChange() map[string]float64 {
	return a.applyScalar(series.NormalizedChange)
}

// Return returns grouped simple returns.
func (a *Impl[in, K]) Return() map[string]float64 {
	return a.applyScalar(finance.Return)
}

// LogReturn returns grouped log returns.
func (a *Impl[in, K]) LogReturn() map[string]float64 {
	return a.applyScalar(finance.LogReturn)
}

// MaxDrawdown returns grouped maximum drawdowns.
func (a *Impl[in, K]) MaxDrawdown() map[string]float64 {
	return a.applyScalar(finance.MaxDrawdown)
}

func (a *Impl[in, K]) applyScalar(fn func([]float64) float64) map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make(map[string]float64, len(a.dataMap))
	for key, values := range a.dataMap {
		out[key] = fn(values)
	}
	return out
}

func (a *Impl[in, K]) applySlice(fn func([]float64) []float64) map[string][]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make(map[string][]float64, len(a.dataMap))
	for key, values := range a.dataMap {
		out[key] = fn(values)
	}
	return out
}

func keyToString[K comparable](key K) string {
	if s, ok := any(key).(string); ok {
		return s
	}
	return fmt.Sprint(key)
}

func timeFrom[in any](item in, timeFn func(in) time.Time) time.Time {
	if timeFn == nil {
		return time.Time{}
	}
	return timeFn(item)
}
