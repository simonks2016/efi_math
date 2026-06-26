package agg

import (
	"github.com/simonks2016/efi_math/core"
	"github.com/simonks2016/efi_math/finance"
	"github.com/simonks2016/efi_math/series"
)

// FieldSelector describes one numeric field extracted from an input item.
type FieldSelector[T any] struct {
	Name string
	Fn   func(T) float64
}

// Field creates a named numeric field selector for NewGroup.
func Field[T any](name string, fn func(T) float64) FieldSelector[T] {
	return FieldSelector[T]{
		Name: name,
		Fn:   fn,
	}
}

// GroupAgg stores grouped, multi-field numeric series.
type GroupAgg[T any, K comparable] struct {
	data   map[K]map[string][]float64
	counts map[K]int
}

// GroupView exposes all fields under one group.
type GroupView[K comparable] struct {
	key    K
	data   map[string][]float64
	count  int
	exists bool
}

// FieldView exposes one field across all groups.
type FieldView[K comparable] struct {
	name string
	data map[K][]float64
}

// GroupFieldView exposes one field under one group.
type GroupFieldView struct {
	values []float64
	exists bool
}

// NewGroup builds a multi-field grouped aggregator in one pass over data.
func NewGroup[T any, K comparable](
	data []T,
	keyFunc func(T) K,
	fields ...FieldSelector[T],
) *GroupAgg[T, K] {
	out := &GroupAgg[T, K]{
		data:   make(map[K]map[string][]float64),
		counts: make(map[K]int),
	}

	for _, item := range data {
		key := keyFunc(item)
		out.counts[key]++
		if _, ok := out.data[key]; !ok {
			out.data[key] = make(map[string][]float64, len(fields))
		}
		for _, field := range fields {
			value := 0.0
			if field.Fn != nil {
				value = field.Fn(item)
			}
			out.data[key][field.Name] = append(out.data[key][field.Name], value)
		}
	}

	return out
}

// Group selects one group. Missing groups return an empty view.
func (a *GroupAgg[T, K]) Group(key K) GroupView[K] {
	if a == nil {
		return GroupView[K]{key: key}
	}
	data, ok := a.data[key]
	return GroupView[K]{
		key:    key,
		data:   data,
		count:  a.counts[key],
		exists: ok,
	}
}

// Field selects one field across all groups. Missing fields return an empty view.
func (a *GroupAgg[T, K]) Field(name string) FieldView[K] {
	out := make(map[K][]float64)
	if a == nil {
		return FieldView[K]{name: name, data: out}
	}
	for key, fields := range a.data {
		values, ok := fields[name]
		if !ok {
			continue
		}
		out[key] = values
	}
	return FieldView[K]{name: name, data: out}
}

// Count returns the item count for each group.
func (a *GroupAgg[T, K]) Count() map[K]float64 {
	out := make(map[K]float64)
	if a == nil {
		return out
	}
	for key, count := range a.counts {
		out[key] = float64(count)
	}
	return out
}

// Sum returns field sums for every group.
func (a *GroupAgg[T, K]) Sum() map[K]map[string]float64 {
	return a.applyAll(core.Sum)
}

// Mean returns field means for every group.
func (a *GroupAgg[T, K]) Mean() map[K]map[string]float64 {
	return a.applyAll(core.Mean)
}

// Max returns field maximum values for every group.
func (a *GroupAgg[T, K]) Max() map[K]map[string]float64 {
	return a.applyAll(core.Max)
}

// Min returns field minimum values for every group.
func (a *GroupAgg[T, K]) Min() map[K]map[string]float64 {
	return a.applyAll(core.Min)
}

// First returns field first values for every group.
func (a *GroupAgg[T, K]) First() map[K]map[string]float64 {
	return a.applyAll(series.First)
}

// Last returns field last values for every group.
func (a *GroupAgg[T, K]) Last() map[K]map[string]float64 {
	return a.applyAll(series.Last)
}

// Range returns field ranges for every group.
func (a *GroupAgg[T, K]) Range() map[K]map[string]float64 {
	return a.applyAll(core.Range)
}

// Change returns field last-first changes for every group.
func (a *GroupAgg[T, K]) Change() map[K]map[string]float64 {
	return a.applyAll(series.Change)
}

// PercentChange returns field percent changes for every group.
func (a *GroupAgg[T, K]) PercentChange() map[K]map[string]float64 {
	return a.applyAll(series.PercentChange)
}

// Return returns field simple returns for every group.
func (a *GroupAgg[T, K]) Return() map[K]map[string]float64 {
	return a.applyAll(finance.Return)
}

// LogReturn returns field log returns for every group.
func (a *GroupAgg[T, K]) LogReturn() map[K]map[string]float64 {
	return a.applyAll(finance.LogReturn)
}

func (a *GroupAgg[T, K]) applyAll(fn func([]float64) float64) map[K]map[string]float64 {
	out := make(map[K]map[string]float64)
	if a == nil {
		return out
	}
	for key, fields := range a.data {
		out[key] = make(map[string]float64, len(fields))
		for name, values := range fields {
			out[key][name] = fn(values)
		}
	}
	return out
}

// Field selects one field under this group. Missing fields return an empty view.
func (g GroupView[K]) Field(name string) GroupFieldView {
	if !g.exists || g.data == nil {
		return GroupFieldView{}
	}
	values, ok := g.data[name]
	if !ok {
		return GroupFieldView{}
	}
	return GroupFieldView{values: values, exists: true}
}

// Count returns the item count for this group.
func (g GroupView[K]) Count() float64 {
	if !g.exists {
		return 0
	}
	return float64(g.count)
}

// Sum returns sums for every field in this group.
func (g GroupView[K]) Sum() map[string]float64 {
	return g.applyFields(core.Sum)
}

// Mean returns means for every field in this group.
func (g GroupView[K]) Mean() map[string]float64 {
	return g.applyFields(core.Mean)
}

// Max returns maximum values for every field in this group.
func (g GroupView[K]) Max() map[string]float64 {
	return g.applyFields(core.Max)
}

// Min returns minimum values for every field in this group.
func (g GroupView[K]) Min() map[string]float64 {
	return g.applyFields(core.Min)
}

// First returns first values for every field in this group.
func (g GroupView[K]) First() map[string]float64 {
	return g.applyFields(series.First)
}

// Last returns last values for every field in this group.
func (g GroupView[K]) Last() map[string]float64 {
	return g.applyFields(series.Last)
}

// Range returns ranges for every field in this group.
func (g GroupView[K]) Range() map[string]float64 {
	return g.applyFields(core.Range)
}

// Change returns last-first changes for every field in this group.
func (g GroupView[K]) Change() map[string]float64 {
	return g.applyFields(series.Change)
}

// PercentChange returns percent changes for every field in this group.
func (g GroupView[K]) PercentChange() map[string]float64 {
	return g.applyFields(series.PercentChange)
}

// Return returns simple returns for every field in this group.
func (g GroupView[K]) Return() map[string]float64 {
	return g.applyFields(finance.Return)
}

// LogReturn returns log returns for every field in this group.
func (g GroupView[K]) LogReturn() map[string]float64 {
	return g.applyFields(finance.LogReturn)
}

func (g GroupView[K]) applyFields(fn func([]float64) float64) map[string]float64 {
	out := make(map[string]float64)
	if !g.exists || g.data == nil {
		return out
	}
	for name, values := range g.data {
		out[name] = fn(values)
	}
	return out
}

// Count returns field counts for every group.
func (f FieldView[K]) Count() map[K]float64 {
	return f.apply(coreCount)
}

// Sum returns field sums for every group.
func (f FieldView[K]) Sum() map[K]float64 {
	return f.apply(core.Sum)
}

// Mean returns field means for every group.
func (f FieldView[K]) Mean() map[K]float64 {
	return f.apply(core.Mean)
}

// Max returns field maximum values for every group.
func (f FieldView[K]) Max() map[K]float64 {
	return f.apply(core.Max)
}

// Min returns field minimum values for every group.
func (f FieldView[K]) Min() map[K]float64 {
	return f.apply(core.Min)
}

// First returns field first values for every group.
func (f FieldView[K]) First() map[K]float64 {
	return f.apply(series.First)
}

// Last returns field last values for every group.
func (f FieldView[K]) Last() map[K]float64 {
	return f.apply(series.Last)
}

// Range returns field ranges for every group.
func (f FieldView[K]) Range() map[K]float64 {
	return f.apply(core.Range)
}

// Change returns field last-first changes for every group.
func (f FieldView[K]) Change() map[K]float64 {
	return f.apply(series.Change)
}

// PercentChange returns field percent changes for every group.
func (f FieldView[K]) PercentChange() map[K]float64 {
	return f.apply(series.PercentChange)
}

// Return returns field simple returns for every group.
func (f FieldView[K]) Return() map[K]float64 {
	return f.apply(finance.Return)
}

// LogReturn returns field log returns for every group.
func (f FieldView[K]) LogReturn() map[K]float64 {
	return f.apply(finance.LogReturn)
}

func (f FieldView[K]) apply(fn func([]float64) float64) map[K]float64 {
	out := make(map[K]float64)
	for key, values := range f.data {
		out[key] = fn(values)
	}
	return out
}

// Count returns the number of values in this group field.
func (g GroupFieldView) Count() float64 {
	if !g.exists {
		return 0
	}
	return float64(core.Count(g.values))
}

// Sum returns the sum of this group field.
func (g GroupFieldView) Sum() float64 {
	return g.applyZero(core.Sum)
}

// Mean returns the mean of this group field.
func (g GroupFieldView) Mean() float64 {
	return g.applyZero(core.Mean)
}

// Max returns the maximum value of this group field.
func (g GroupFieldView) Max() float64 {
	return g.applyZero(core.Max)
}

// Min returns the minimum value of this group field.
func (g GroupFieldView) Min() float64 {
	return g.applyZero(core.Min)
}

// First returns the first value of this group field.
func (g GroupFieldView) First() float64 {
	return g.applyZero(series.First)
}

// Last returns the last value of this group field.
func (g GroupFieldView) Last() float64 {
	return g.applyZero(series.Last)
}

// Range returns the range of this group field.
func (g GroupFieldView) Range() float64 {
	return g.applyZero(core.Range)
}

// Change returns last-first for this group field.
func (g GroupFieldView) Change() float64 {
	return g.applyZero(series.Change)
}

// PercentChange returns percent change for this group field.
func (g GroupFieldView) PercentChange() float64 {
	return g.applyZero(series.PercentChange)
}

// Return returns simple return for this group field.
func (g GroupFieldView) Return() float64 {
	return g.applyZero(finance.Return)
}

// LogReturn returns log return for this group field.
func (g GroupFieldView) LogReturn() float64 {
	return g.applyZero(finance.LogReturn)
}

func (g GroupFieldView) applyZero(fn func([]float64) float64) float64 {
	if !g.exists {
		return 0
	}
	return fn(g.values)
}

func coreCount(values []float64) float64 {
	return float64(core.Count(values))
}
