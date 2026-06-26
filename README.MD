# eFiMath / efi_math

通用数学与金融计算库，面向 Go 项目中的基础统计、序列统计、金融指标和泛型分组聚合。

A general-purpose math and finance calculation library for Go, covering basic statistics, ordered series metrics, financial indicators, and generic grouped aggregation.

> GitHub: `github.com/simonks2016/efi_math`

## ✨ 特性 / Features

| 中文 | English |
| --- | --- |
| 🧮 只处理通用数值，不绑定业务模型 | Numeric-first design with no business model coupling |
| 📊 `core` 提供无顺序基础统计 | `core` provides order-independent statistics |
| 📈 `series` 提供有顺序序列统计 | `series` provides order-aware sequence metrics |
| 💹 `finance` 提供常见金融指标 | `finance` provides common financial indicators |
| 🧩 `agg` 支持任意结构体的泛型分组聚合 | `agg` supports generic grouped aggregation over arbitrary records |
| ⏱️ `WithTime` 可按组内时间排序后计算序列指标 | `WithTime` sorts each group by time before sequence calculations |

## 📦 安装 / Installation

```bash
go get github.com/simonks2016/efi_math
```

导入示例：

Import examples:

```go
import (
	"github.com/simonks2016/efi_math/agg"
	"github.com/simonks2016/efi_math/core"
	"github.com/simonks2016/efi_math/finance"
	"github.com/simonks2016/efi_math/series"
)
```

## 🚀 快速开始 / Quick Start

### 基础统计 / Basic Statistics

```go
package main

import (
	"fmt"

	"github.com/simonks2016/efi_math/core"
)

func main() {
	x := []float64{1, 2, 3, 4}

	fmt.Println(core.Sum(x))      // 10
	fmt.Println(core.Mean(x))     // 2.5
	fmt.Println(core.Quantile(x, 0.5)) // 2.5
}
```

### 序列统计 / Ordered Series Metrics

```go
package main

import (
	"fmt"

	"github.com/simonks2016/efi_math/series"
)

func main() {
	x := []float64{0.1, 0.2, 0.1}

	fmt.Println(series.Diff(x))           // [0.1 -0.1]
	fmt.Println(series.SumAbsDiff(x))     // 0.2
	fmt.Println(series.SumSquaredDiff(x)) // 0.02
	fmt.Println(series.Change(x))         // 0
	fmt.Println(series.FlipRate(x))       // 1
}
```

### 金融指标 / Financial Indicators

```go
package main

import (
	"fmt"

	"github.com/simonks2016/efi_math/finance"
)

func main() {
	equity := []float64{100, 150, 120, 200}

	fmt.Println(finance.Return(equity))      // 1
	fmt.Println(finance.MaxDrawdown(equity)) // 0.2
	fmt.Println(finance.Calmar(equity))      // 5
}
```

### 泛型分组聚合 / Generic Grouped Aggregation

```go
package main

import (
	"fmt"

	"github.com/simonks2016/efi_math/agg"
)

type Trade struct {
	IsBuy  bool
	Price  float64
	Volume float64
}

func main() {
	trades := []Trade{
		{IsBuy: true, Price: 10, Volume: 2},
		{IsBuy: false, Price: 20, Volume: 3},
		{IsBuy: true, Price: 12, Volume: 4},
		{IsBuy: false, Price: 25, Volume: 1},
	}

	tradeAgg := agg.NewGroup(
		trades,
		func(trade Trade) string {
			if trade.IsBuy {
				return "buy"
			}
			return "sell"
		},
		agg.Field("volume", func(trade Trade) float64 {
			return trade.Volume
		}),
		agg.Field("notional", func(trade Trade) float64 {
			return trade.Price * trade.Volume
		}),
	)

	fmt.Println(tradeAgg.Field("volume").Sum())               // map[buy:6 sell:4]
	fmt.Println(tradeAgg.Field("notional").Sum())             // map[buy:68 sell:85]
	fmt.Println(tradeAgg.Group("buy").Field("volume").Sum())  // 6
	fmt.Println(tradeAgg.Group("sell").Field("notional").Sum()) // 85
	fmt.Println(tradeAgg.Group("buy").Count())                // 2
}
```

## 🧱 模块结构 / Package Layout

```text
efi_math/
  core/      // 无顺序基础统计 / order-independent statistics
  series/    // 有顺序序列统计 / order-aware series metrics
  finance/   // 金融语义指标 / financial indicators
  agg/       // 泛型分组聚合器 / generic grouped aggregator
```

## 🧮 core：无顺序基础统计 / Order-Independent Statistics

`core` 只处理 `[]float64`，不关心时间顺序。

`core` only works with `[]float64` and does not depend on input order.

| 函数 / Function | 说明 / Description |
| --- | --- |
| `Count(x []float64) int` | 元素数量 / number of values |
| `Sum(x []float64) float64` | 求和 / sum |
| `Mean(x []float64) float64` | 算术平均值 / arithmetic mean |
| `Max(x []float64) float64` | 最大值 / maximum |
| `Min(x []float64) float64` | 最小值 / minimum |
| `Range(x []float64) float64` | 极差，`Max-Min` / range, `Max-Min` |
| `Var(x []float64) float64` | 总体方差，除以 `n` / population variance divided by `n` |
| `Std(x []float64) float64` | 总体标准差 / population standard deviation |
| `Median(x []float64) float64` | 中位数 / median |
| `Quantile(x []float64, q float64) float64` | 分位数，`q` 在 `[0,1]` / quantile with `q` in `[0,1]` |

## 📈 series：有顺序序列统计 / Order-Aware Series Metrics

`series` 也只处理 `[]float64`，但会保留输入顺序。

`series` also works with `[]float64`, but the input order is meaningful.

| 函数 / Function | 说明 / Description |
| --- | --- |
| `First(x []float64) float64` | 第一个值 / first value |
| `Last(x []float64) float64` | 最后一个值 / last value |
| `Diff(x []float64) []float64` | 相邻差分，`x[i]-x[i-1]` / adjacent differences |
| `AbsDiff(x []float64) []float64` | 相邻差分绝对值 / absolute adjacent differences |
| `SquaredDiff(x []float64) []float64` | 相邻差分平方 / squared adjacent differences |
| `SumAbsDiff(x []float64) float64` | 相邻绝对差分和 / sum of absolute adjacent differences |
| `SumSquaredDiff(x []float64) float64` | 相邻平方差分和 / sum of squared adjacent differences |
| `Change(x []float64) float64` | `last-first` |
| `PercentChange(x []float64) float64` | `(last-first)/first` |
| `NormalizedChange(x []float64) float64` | `(last-first)/mean(x)` |
| `Gradient(x []float64) float64` | `(last-first)/(len(x)-1)` |
| `FlipRate(x []float64) float64` | 相邻非零差分方向翻转比例 / ratio of sign flips among non-zero adjacent differences |

## 💹 finance：金融指标 / Financial Indicators

`finance` 处理价格、净值、收益率等数值序列。它只提供金融计算语义，不绑定任何仓位、订单、PnL 或业务模型。

`finance` works with price, equity, and return sequences. It provides financial calculation semantics without coupling to positions, orders, PnL, or business models.

| 函数 / Function | 说明 / Description |
| --- | --- |
| `Return(x []float64) float64` | 简单收益，`(last-first)/first` / simple return |
| `LogReturn(x []float64) float64` | 对数收益，`ln(last/first)` / log return |
| `SimpleReturns(prices []float64) []float64` | 相邻价格简单收益 / adjacent simple returns |
| `LogReturns(prices []float64) []float64` | 相邻价格对数收益 / adjacent log returns |
| `MaxDrawdown(equity []float64) float64` | 最大回撤，返回正数 / maximum drawdown as a positive ratio |
| `Volatility(returns []float64) float64` | 收益率标准差 / standard deviation of returns |
| `Sharpe(returns []float64, riskFreeRate float64) float64` | `(mean(returns)-riskFreeRate)/std(returns)` |
| `Calmar(equity []float64) float64` | `Return(equity)/MaxDrawdown(equity)` |
| `EMA(x []float64, period int) []float64` | 指数移动平均 / exponential moving average |
| `DMA(x []float64, alpha float64) []float64` | 动态移动平均 / dynamic moving average |
| `SMA(x []float64, period int) []float64` | 简单移动平均 / simple moving average |
| `BollingerBands(x []float64, period int, k float64)` | 布林带，返回 `upper/middle/lower` / Bollinger Bands returning `upper/middle/lower` |

## 🧩 agg：泛型分组聚合器 / Generic Grouped Aggregator

`agg` 接受任意 `[]T`，通过 `keyFn` 分组，通过多个 `Field` 提取数值字段，然后调用 `core`、`series`、`finance` 完成计算。

`agg` accepts arbitrary `[]T`, groups records with `keyFn`, extracts multiple numeric fields with `Field`, and delegates calculations to `core`, `series`, and `finance`.

### 推荐 API：Group -> Field -> Agg / Recommended API

```go
func Field[T any](name string, fn func(T) float64) FieldSelector[T]

func NewGroup[T any, K comparable](
	data []T,
	keyFunc func(T) K,
	fields ...FieldSelector[T],
) *GroupAgg[T, K]
```

调用方式：

Usage:

```go
tradeAgg.Field("volume").Sum()               // map[K]float64
tradeAgg.Field("notional").Mean()            // map[K]float64
tradeAgg.Group("buy").Count()                // float64
tradeAgg.Group("buy").Field("volume").Sum()  // float64
```

底层数据模型：

Underlying data model:

```go
map[K]map[string][]float64
```

`K` 是 `comparable`，可以是 `string`、`bool`、`int` 等，不会强制转换成字符串。

`K` is `comparable`; it can be `string`, `bool`, `int`, and so on. It is not forced into a string key.

### 兼容 API / Compatibility API

```go
func NewAggregator[in any, K comparable](
	data []in,
	keyFunc func(in) K,
	valueFunc func(in) float64,
	opts ...Option[in],
) Aggregator

func New[in any, K comparable](
	data []in,
	keyFunc func(in) K,
	valueFunc func(in) float64,
	opts ...Option[in],
) Aggregator
```

这套接口保留给旧的单字段聚合代码使用。新业务建议使用 `NewGroup`。

This interface remains available for older single-field aggregation code. New business code should prefer `NewGroup`.

### Option for Compatibility API

| Option | 说明 / Description |
| --- | --- |
| `WithTime(func(T) time.Time)` | 每组内部按时间升序排序后再计算序列类指标 / sort each group by time before sequence-aware metrics |

如果不传 `WithTime`，每组按原始输入顺序计算。

If `WithTime` is not provided, each group keeps the original input order.

### 聚合方法 / Aggregation Methods

| 类别 / Category | 方法 / Methods |
| --- | --- |
| `GroupAgg` | `Count`, `Sum`, `Mean`, `Max`, `Min`, `First`, `Last`, `Range`, `Change`, `PercentChange`, `Return`, `LogReturn` |
| `FieldView` | `Count`, `Sum`, `Mean`, `Max`, `Min`, `First`, `Last`, `Range`, `Change`, `PercentChange`, `Return`, `LogReturn` |
| `GroupView` | `Count`, `Field`, `Sum`, `Mean`, `Max`, `Min`, `First`, `Last`, `Range`, `Change`, `PercentChange`, `Return`, `LogReturn` |
| `GroupFieldView` | `Count`, `Sum`, `Mean`, `Max`, `Min`, `First`, `Last`, `Range`, `Change`, `PercentChange`, `Return`, `LogReturn` |

## ⚠️ 边界行为 / Edge Cases

| 场景 / Case | 行为 / Behavior |
| --- | --- |
| 空数组 / empty slice | 无法计算的标量返回 `math.NaN()` / invalid scalar results return `math.NaN()` |
| 长度不足 / insufficient length | 序列指标按函数语义返回 `NaN` 或空切片 / series functions return `NaN` or empty slices according to semantics |
| 除以 0 / division by zero | 返回 `math.NaN()` / returns `math.NaN()` |
| `Quantile` 的 `q` 不在 `[0,1]` / invalid quantile `q` | 返回 `math.NaN()` / returns `math.NaN()` |
| `SMA`、`BollingerBands` 前期窗口不足 / incomplete rolling windows | 对应位置返回 `math.NaN()` / returns `math.NaN()` at those positions |

## ✅ 测试 / Testing

```bash
go test ./...
```

如果本地环境限制 Go build cache 写入，可以临时指定 `GOCACHE`：

If your environment restricts Go build cache writes, use a temporary `GOCACHE`:

```bash
GOCACHE=/private/tmp/efi_math-go-cache go test ./...
```

## 📌 设计原则 / Design Principles

- 不在底层模块中出现仓位、价格对象、PnL、Holding 等业务模型。
- `core`、`series`、`finance` 只处理数值序列。
- `agg` 只负责泛型分组、数值提取和可选时间排序。
- 数学逻辑集中在对应模块中，聚合器不重复实现计算公式。

- No business models such as positions, price objects, PnL, or holdings in low-level packages.
- `core`, `series`, and `finance` only operate on numeric sequences.
- `agg` only handles generic grouping, value extraction, and optional time sorting.
- Mathematical logic stays in the dedicated packages instead of being duplicated in the aggregator.
