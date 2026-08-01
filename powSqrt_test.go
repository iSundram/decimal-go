package decimal

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestPowSqrt(t *testing.T) {
	Default.Config(&Config{
		ToExpNeg: i64(-7),
		ToExpPos: i64(21),
		MinE:     i64(-9e15),
		MaxE:     i64(9e15),
	})

	// In JS the loop condition is `total < 10000`, where total counts the
	// assertions performed; each iteration performs one assertion.
	for total := 0; total < 10000; total++ {
		// Get a random value in the range [0,1) with a random number of
		// significant digits in the range [1, 40], as a string in
		// exponential format.
		e := Default.Random(int64(rand.Float64()*40 + 1)).ToExponential()

		// Change exponent to a non-zero value of random length in the
		// range (-9e15, 9e15).
		n := strconv.FormatInt(int64(math.Floor(rand.Float64()*9e15)), 10)
		sign := ""
		if rand.Float64() < 0.5 {
			sign = "-"
		}
		r := New(e[:strings.IndexByte(e, 'e')+1] + sign + n[int(rand.Float64()*float64(len(n))):])

		// Random rounding mode.
		Default.Rounding = int64(rand.Float64() * 9)

		// Random precision in the range [1, 40].
		Default.Precision = int64(rand.Float64()*40 + 1)

		p := r.Pow(0.5)

		// sqrt is much faster than pow(0.5)
		s := r.Sqrt()

		assertEq(t, p.ValueOf(), s.ValueOf())
	}
}
