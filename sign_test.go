package decimal

import "testing"

func TestSign(t *testing.T) {
	tf := func(n any, expected float64) {
		t.Helper()
		assertEq(t, expected, Default.Sign(n))
	}

	tf(nan(), nan())
	tf("NaN", nan())
	tf(inf(1), 1.0)
	tf(inf(-1), -1.0)
	tf("Infinity", 1.0)
	tf("-Infinity", -1.0)

	assert(t, 1/Default.Sign("0") == inf(1))
	assert(t, 1/Default.Sign(New("0")) == inf(1))
	assert(t, 1/Default.Sign("-0") == inf(-1))
	assert(t, 1/Default.Sign(New("-0")) == inf(-1))

	tf("0", 0.0)
	tf("-0", math_NegZero())
	tf("1", 1.0)
	tf("-1", -1.0)
	tf("9.99", 1.0)
	tf("-9.99", -1.0)
}
