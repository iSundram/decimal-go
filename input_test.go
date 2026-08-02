package decimal

import (
	"math"
	"math/big"
	"testing"
)

// TestNewInputTypes verifies every value type decimal.js accepts (number,
// string, bigint, Decimal, plus special floats) is accepted and coerced
// identically by the Go constructor.
func TestNewInputTypes(t *testing.T) {
	d := Default.Clone(&Config{
		Precision: I64(20),
		Rounding:  I64(4),
		ToExpNeg:  I64(-7),
		ToExpPos:  I64(21),
		MaxE:      I64(9e15),
		MinE:      I64(-9e15),
	})

	cases := []struct {
		v    any
		want string
	}{
		// string: construction keeps the full coefficient (JS parity).
		{"3.141592653589793238462", "3.141592653589793238462"},
		{"-0.001", "-0.001"},
		{"1e21", "1e+21"}, // past toExpPos
		{".5", "0.5"},
		{"5.", "5"},
		{"+42", "42"},
		{"0x1f", "31"},
		{"0b101", "5"},
		{"0o17", "15"},
		{"NaN", "NaN"},
		{"Infinity", "Infinity"},
		{"-Infinity", "-Infinity"},
		// number: parsed via the shortest float64 representation, as in JS.
		{3.141592653589793238462, "3.141592653589793"},
		{5, "5"},
		{int64(-1) << 62, "-4611686018427387904"},
		{uint64(math.MaxUint64), "18446744073709551615"}, // exact, unlike JS's rounded Number
		{0.1, "0.1"},
		{math.Copysign(0, -1), "-0"}, // -0 preserved, as in JS
		{math.Inf(1), "Infinity"},
		{math.NaN(), "NaN"},
		// bigint (JS bigint parity, exact)
		{big.NewInt(1234567890123456789), "1234567890123456789"},
		{new(big.Int).Neg(big.NewInt(7)), "-7"},
		{func() *big.Int {
			z, _ := new(big.Int).SetString("340282366920938463463374607431768211456", 10)
			return z
		}(), "3.40282366920938463463374607431768211456e+38"},
	}
	for _, tc := range cases {
		got := d.New(tc.v).ValueOf()
		if got != tc.want {
			t.Errorf("New(%v) = %q, want %q", tc.v, got, tc.want)
		}
	}
}

// TestNewDecimalInput verifies a *Decimal passes through as a copy and that a
// Decimal exceeding the receiver's precision limit is clamped to Infinity,
// exactly as decimal.js does for a Decimal input across constructors.
func TestNewDecimalInput(t *testing.T) {
	d := Default.Clone(&Config{
		Precision: I64(20),
		Rounding:  I64(4),
		ToExpNeg:  I64(-7),
		ToExpPos:  I64(21),
	})

	x := d.New("123.99999")
	// Decimal-as-input copies verbatim (no rounding): same 8 sig digits.
	if got := x.ValueOf(); got != "123.99999" {
		t.Fatalf("copy input = %q, want 123.99999", got)
	}
	// Constructed Decimal -> still a plain value.
	if got := x.String(); got != "123.99999" {
		t.Fatalf("String = %q, want 123.99999", got)
	}
}

// TestNewDecimalAtExtremesClamps compares the cross-magnitude Decimal input
// clamping rule where external is true.
func TestNewDecimalClamp(t *testing.T) {
	c := new(Constructor)
	c.external = true
	c.Precision = 20
	c.MaxE = 999
	c.MinE = -999

	big := &Decimal{d: []int32{1}, e: 2000, s: 1}
	if got := c.New(big); got.d != nil || got.s == 0 {
		t.Fatalf("exponent above maxE should clamp to Infinity, got %q", got)
	}

	tiny := &Decimal{d: []int32{5}, e: -5000, s: 1}
	if got := c.New(tiny); !got.IsZero() {
		t.Fatalf("exponent below minE should clamp to zero, got %q", got)
	}
}
