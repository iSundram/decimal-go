package decimal

import (
	"math"
	"testing"
)

func TestToNumber(t *testing.T) {
	setCfg(20, 4, -7, 21, 9e15, -9e15, 1, false)

	// Positive zero
	tfp := func(n any) {
		t.Helper()
		assert(t, 1/New(n).Float64() == inf(1))
	}

	tfp("0")
	tfp("0.0")
	tfp("0.000000000000")
	tfp("0e+0")
	tfp("0e-0")
	tfp("1e-9000000000000000")

	// Negative zero
	tfn := func(n any) {
		t.Helper()
		assert(t, 1/New(n).Float64() == inf(-1))
	}

	tfn("-0")
	tfn("-0.0")
	tfn("-0.000000000000")
	tfn("-0e+0")
	tfn("-0e-0")
	tfn("-1e-9000000000000000")

	tf := func(n any, expected float64) {
		t.Helper()
		assertEq(t, expected, New(n).Float64())
	}

	tf(inf(1), inf(1))
	tf("Infinity", inf(1))
	tf(inf(-1), inf(-1))
	tf("-Infinity", inf(-1))
	tf(nan(), nan())
	tf("NaN", nan())

	tf(1.0, 1.0)
	tf("1", 1.0)
	tf("1.0", 1.0)
	tf("1e+0", 1.0)
	tf("1e-0", 1.0)

	tf(-1.0, -1.0)
	tf("-1", -1.0)
	tf("-1.0", -1.0)
	tf("-1e+0", -1.0)
	tf("-1e-0", -1.0)

	tf("123.456789876543", 123.456789876543)
	tf("-123.456789876543", -123.456789876543)

	tf("1.1102230246251565e-16", 1.1102230246251565e-16)
	tf("-1.1102230246251565e-16", -1.1102230246251565e-16)

	tf("9007199254740991", 9007199254740991.0)
	tf("-9007199254740991", -9007199254740991.0)

	tf("5e-324", 5e-324)
	tf("1.7976931348623157e+308", 1.7976931348623157e+308)

	tf("9.999999e+9000000000000000", inf(1))
	tf("-9.999999e+9000000000000000", inf(-1))
	tf("1e-9000000000000000", 0.0)
	// JS: t('-1e-9000000000000000', -0) — the result must be -0.
	assert(t, math.Signbit(New("-1e-9000000000000000").Float64()))
	assertEq(t, 0.0, New("-1e-9000000000000000").Float64())
}
