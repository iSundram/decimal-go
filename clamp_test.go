package decimal

import "testing"

func TestClamp(t *testing.T) {
	tf := func(x, min, max any, expected string) {
		t.Helper()
		assertEq(t, expected, New(x).Clamp(min, max).ValueOf())
	}

	tf("-0", "0", "0", "-0")
	tf("-0", "-0", "0", "-0")
	tf("-0", "0", "-0", "-0")
	tf("-0", "-0", "-0", "-0")

	tf("0", "0", "0", "0")
	tf("0", "-0", "0", "0")
	tf("0", "0", "-0", "0")
	tf("0", "-0", "-0", "0")

	tf(0.0, 0.0, 1.0, "0")
	tf(-1.0, 0.0, 1.0, "0")
	tf(-2.0, 0.0, 1.0, "0")
	tf(1.0, 0.0, 1.0, "1")
	tf(2.0, 0.0, 1.0, "1")

	tf(1.0, 1.0, 1.0, "1")
	tf(-1.0, 1.0, 1.0, "1")
	tf(-1.0, -1.0, 1.0, "-1")
	tf(2.0, 1.0, 2.0, "2")
	tf(3.0, 1.0, 2.0, "2")
	tf(1.0, 0.0, 1.0, "1")
	tf(2.0, 0.0, 1.0, "1")

	tf(inf(1), 0.0, 1.0, "1")
	tf(0.0, inf(-1), 0.0, "0")
	tf(inf(-1), 0.0, 1.0, "0")
	tf(inf(-1), inf(-1), inf(1), "-Infinity")
	tf(inf(1), inf(-1), inf(1), "Infinity")
	tf(0.0, 1.0, inf(1), "1")

	tf(0.0, nan(), 1.0, "NaN")
	tf(0.0, 0.0, nan(), "NaN")
	tf(nan(), 0.0, 1.0, "NaN")
}
