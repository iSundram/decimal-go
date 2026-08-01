package decimal

import (
	"math"
	"math/rand"
	"testing"
)

func TestRandom(t *testing.T) {
	maxDigitsVal := int64(100)

	for i := 0; i < 996; i++ {
		// sd = Math.random() * maxDigits + 1 | 0
		sd := int64(rand.Float64()*float64(maxDigitsVal) + 1)

		var r *Decimal
		if rand.Float64() > 0.5 {
			Default.Precision = sd
			r = Default.Random()
		} else {
			r = Default.Random(sd)
		}

		assert(t, r.Sd() <= float64(sd) && r.Gte(0) && r.Lt(1) && r.Eq(r) && r.Eq(r.ValueOf()))
	}

	// JS passes Infinity/'-Infinity'/NaN/null to Decimal.random; in Go the
	// same out-of-range sd values come from float64->int64 conversions
	// (which yield a value below the minimum) and 0 for null.
	assertException(t, func() { Default.Random(int64(math.Inf(1))) }, "Infinity")
	assertException(t, func() { Default.Random(int64(math.Inf(-1))) }, "'-Infinity'")
	assertException(t, func() { Default.Random(int64(math.NaN())) }, "NaN")
	assertException(t, func() { Default.Random(0) }, "null")
}
