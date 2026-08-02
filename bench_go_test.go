package decimal

import (
	"fmt"
	"testing"
)

func printf(v int64) string { return fmt.Sprintf("%d", v) }

// These benchmarks measure the idiomatic Go API on the default constructor.
// They are intentionally "cold API" benchmarks: each iteration constructs
// operands from strings (the most common real-world usage) so the results are
// directly comparable to the JS runner in bench/js.

func benchConfig(b *testing.B, precision int64) {
	b.Helper()
	setCfg(precision, 4, -7, 21, 9e15, -9e15, 1, false)
}

// runSub repeats fn n kinds of work per op.
func runSub(b *testing.B, f func() any) {
	b.Helper()
	var sink any
	for b.Loop() {
		sink = f()
	}
	_ = sink
}

func BenchmarkNew(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			runSub(b, func() any { return New("3.1415926535897932384626433832795028841971") })
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x, y := New("1.23456789012345"), New("9.8765432109876")
			runSub(b, func() any { return x.Plus(y) })
		})
	}
}

func BenchmarkSub(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x, y := New("1.23456789012345"), New("9.8765432109876")
			runSub(b, func() any { return x.Minus(y) })
		})
	}
}

func BenchmarkMul(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x, y := New("1.23456789012345"), New("9.8765432109876")
			runSub(b, func() any { return x.Times(y) })
		})
	}
}

// ConstructFromString pins the full parse-and-apply path so it compares
// 1:1 with decimal.js, where constructing the operands is part of every op.
func BenchmarkAddFromString(b *testing.B) {
	for _, p := range []int64{20} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			runSub(b, func() any { return New("1.23456789012345").Plus(New("9.8765432109876")) })
		})
	}
}

func BenchmarkDiv(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x, y := New("1.23456789012345"), New("9.8765432109876")
			runSub(b, func() any { return x.Div(y) })
		})
	}
}

func BenchmarkDivToInt(b *testing.B) {
	benchConfig(b, 20)
	x, y := New("123456789.5"), New("0.25")
	runSub(b, func() any { return x.DivToInt(y) })
}

func BenchmarkMod(b *testing.B) {
	benchConfig(b, 20)
	x, y := New("123456789.123"), New("0.001")
	runSub(b, func() any { return x.Mod(y) })
}

func BenchmarkPowInt(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("9")
			runSub(b, func() any { return x.Pow(75) })
		})
	}
}

func BenchmarkPowFraction(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("2")
			runSub(b, func() any { return x.Pow(0.5) })
		})
	}
}

func BenchmarkSqrt(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("2")
			runSub(b, func() any { return x.Sqrt() })
		})
	}
}

func BenchmarkCbrt(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("2")
			runSub(b, func() any { return x.Cbrt() })
		})
	}
}

func BenchmarkExp(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("1")
			runSub(b, func() any { return x.Exp() })
		})
	}
}

func BenchmarkLn(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("2")
			runSub(b, func() any { return x.Ln() })
		})
	}
}

func BenchmarkLogBase10(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("12345.6789")
			runSub(b, func() any { return x.Log(10) })
		})
	}
}

func BenchmarkSin(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("0.5")
			runSub(b, func() any { return x.Sin() })
		})
	}
}

func BenchmarkCos(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("1.2")
			runSub(b, func() any { return x.Cos() })
		})
	}
}

func BenchmarkTan(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("0.7")
			runSub(b, func() any { return x.Tan() })
		})
	}
}

func BenchmarkAtan(b *testing.B) {
	for _, p := range []int64{20, 100, 1000} {
		b.Run("precision="+printf(p), func(b *testing.B) {
			benchConfig(b, p)
			x := New("1")
			runSub(b, func() any { return x.Atan() })
		})
	}
}

func BenchmarkToFixed(b *testing.B) {
	benchConfig(b, 20)
	x := New("12345.67890123456789")
	runSub(b, func() any { return x.ToFixed(8) })
}

func BenchmarkToExponential(b *testing.B) {
	benchConfig(b, 20)
	x := New("0.0000123456789")
	runSub(b, func() any { return x.ToExponential(6) })
}

func BenchmarkToPrecision(b *testing.B) {
	benchConfig(b, 20)
	x := New("123456789.123456789")
	runSub(b, func() any { return x.ToPrecision(12) })
}

func BenchmarkCmp(b *testing.B) {
	benchConfig(b, 20)
	x, y := New("1.234567890123456789"), New("1.234567890123456788")
	runSub(b, func() any { return x.Cmp(y) })
}

func BenchmarkRoundTrip(b *testing.B) {
	benchConfig(b, 20)
	runSub(b, func() any {
		x := New("12345678.901234567890123456789")
		_ = x.Plus("0.987654321098765432109")
		return x.ValueOf()
	})
}
