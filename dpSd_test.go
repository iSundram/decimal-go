package decimal

import "testing"

func TestDpSd(t *testing.T) {
	// decimal.js aliases: decimalPlaces() == dp(), precision() == sd().
	// Go exposes only Dp() and Sd(), so both JS alias assertions map to
	// the same Go method here.
	tf := func(n any, dp, sd float64, zs ...bool) {
		t.Helper()
		x := New(n)
		assertEq(t, dp, x.Dp())      // dp()
		assertEq(t, dp, x.Dp())      // decimalPlaces()
		assertEq(t, sd, x.Sd(zs...)) // sd(zs)
		assertEq(t, sd, x.Sd(zs...)) // precision(zs)
	}

	setCfg(20, 4, -7, 21, 9e15, -9e15, 1, false)

	tf(0.0, 0.0, 1.0)
	tf(math_NegZero(), 0.0, 1.0)
	tf(nan(), nan(), nan())
	tf(inf(1), nan(), nan())
	tf(inf(-1), nan(), nan())
	tf(1.0, 0.0, 1.0)
	tf(-1.0, 0.0, 1.0)

	tf(100.0, 0.0, 1.0)
	tf(100.0, 0.0, 1.0, false)
	tf(100.0, 0.0, 1.0, false) // zs: false
	tf(100.0, 0.0, 3.0, true)  // zs: 1
	tf(100.0, 0.0, 3.0, true)  // zs: true

	tf("0.0012345689", 10.0, 8.0)
	tf("0.0012345689", 10.0, 8.0, false)
	tf("0.0012345689", 10.0, 8.0, false)
	tf("0.0012345689", 10.0, 8.0, true)
	tf("0.0012345689", 10.0, 8.0, true)

	tf("987654321000000.0012345689000001", 16.0, 31.0, false)
	tf("987654321000000.0012345689000001", 16.0, 31.0, true)

	tf("1e+123", 0.0, 1.0)
	tf("1e+123", 0.0, 124.0, true)
	tf("1e-123", 123.0, 1.0)
	tf("1e-123", 123.0, 1.0, true)

	tf("9.9999e+9000000000000000", 0.0, 5.0, false)
	tf("9.9999e+9000000000000000", 0.0, 9000000000000001.0, true)
	tf("-9.9999e+9000000000000000", 0.0, 5.0, false)
	tf("-9.9999e+9000000000000000", 0.0, 9000000000000001.0, true)

	tf("1e-9000000000000000", 9e15, 1.0, false)
	tf("1e-9000000000000000", 9e15, 1.0, true)
	tf("-1e-9000000000000000", 9e15, 1.0, false)
	tf("-1e-9000000000000000", 9e15, 1.0, true)

	tf("55325252050000000000000000000000.000000004534500000001", 21.0, 53.0)

	// The following JS exception assertions are unportable: sd(null),
	// sd(2), sd('3') and sd({}) are compile-time type errors in Go, whose
	// Sd takes an optional bool.
	//   tx(function () {new Decimal(1).precision(null)}, "new Decimal(1).precision(null)");
	//   tx(function () {new Decimal(1).sd(null)}, "new Decimal(1).sd(null)");
	//   tx(function () {new Decimal(1).sd(2)}, "new Decimal(1).sd(2)");
	//   tx(function () {new Decimal(1).sd('3')}, "new Decimal(1).sd('3')");
	//   tx(function () {new Decimal(1).sd({})}, "new Decimal(1).sd({})");
}
