package decimal

import "testing"

func TestSum(t *testing.T) {
	var expected *Decimal

	ts := func(args ...any) {
		t.Helper()
		assertEqDecimal(t, expected, Default.Sum(args...))
	}

	expected = New(0)

	ts("0")
	ts("0", New(0))
	ts(1.0, 0.0, "-1")
	ts(0.0, New(-10), 0.0, 0.0, 0.0, 0.0, 0.0, 10.0)
	ts(11.0, -11.0)
	ts(1.0, "2", New(3), New("4"), -10.0)
	ts(New(-10), "9", New(0.01), 0.99)

	expected = New(10)

	ts("10")
	ts("0", New("10"))
	ts(10.0, 0.0)
	ts(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 10.0)
	ts(11.0, -1.0)
	ts(1.0, "2", New(3), New("4"))
	ts("9", New(0.01), 0.99)

	expected = New(600)

	ts(100.0, 200.0, 300.0)
	ts("100", "200", "300")
	ts(New(100), New(200), New(300))
	ts(100.0, "200", New(300))
	ts(99.9, 200.05, 300.05)

	expected = New(nan())

	ts(nan())
	ts("1", nan())
	ts(100.0, 200.0, nan())
	ts(nan(), 0.0, "9", New(0), 11.0, inf(1))
	ts(0.0, New("-Infinity"), "9", New(nan()), 11.0)
	ts(4.0, "-Infinity", 0.0, "9", New(0), inf(1), 2.0)

	expected = New(inf(1))

	ts(inf(1))
	ts(1.0, "1e10000000000000000000000000000000000000000", "4")
	ts(100.0, 200.0, "Infinity")
	ts(0.0, New("Infinity"), "9", New(0), 11.0)
	ts(0.0, "9", New(0), 11.0, inf(1))
	ts(4.0, New(inf(1)), 0.0, "9", New(0), inf(1), 2.0)

	expected = New(inf(-1))

	ts(inf(-1))
	ts(1.0, "-1e10000000000000000000000000000000000000000", "4")
	ts(100.0, 200.0, "-Infinity")
	ts(0.0, New("-Infinity"), "9", New(0), 11.0)
	ts(0.0, "9", New(0), 11.0, inf(-1))
	ts(4.0, New(inf(-1)), 0.0, "9", New(0), inf(-1), 2.0)
}
