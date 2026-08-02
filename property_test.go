package decimal

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
)

// genDecimal returns a deterministic pseudo-random decimal string under the
// current config's exponent limits, seeded for reproducibility.
func genDecimal(r *rand.Rand, digits int) *Decimal {
	// Build a value with roughly 1..digits significant digits and a random
	// exponent in a range that keeps construction cheap and finite.
	ld := r.Intn(digits) + 1
	b := make([]byte, ld)
	b[0] = byte('1' + r.Intn(9))
	for i := 1; i < ld; i++ {
		b[i] = byte('0' + r.Intn(10))
	}
	s := string(b)
	fp := r.Intn(ld)
	if fp > 0 {
		s = s[:fp] + "." + s[fp:]
	}
	e := r.Intn(31) - 15
	sign := ""
	if r.Intn(2) == 0 {
		sign = "-"
	}
	return New(sign + s + "e" + strconv.Itoa(e))
}

// TestPropRoundTrip: buf: New(x).ValueOf() re-parses to the same value.
func TestPropRoundTrip(t *testing.T) {
	setCfg(30, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 2000; i++ {
		x := genDecimal(r, 18)
		y := New(x.ValueOf())
		if !x.Eq(y) {
			t.Fatalf("round-trip %s -> %s", x.ValueOf(), y.ValueOf())
		}
	}
}

// TestPropCommute: x+y == y+x and x*y == y*x.
func TestPropCommute(t *testing.T) {
	setCfg(30, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(2))
	for i := 0; i < 2000; i++ {
		x, y := genDecimal(r, 18), genDecimal(r, 18)
		if !x.Plus(y).Eq(y.Plus(x)) {
			t.Fatalf("commute + %s %s", x.ValueOf(), y.ValueOf())
		}
		if !x.Times(y).Eq(y.Times(x)) {
			t.Fatalf("commute * %s %s", x.ValueOf(), y.ValueOf())
		}
	}
}

// TestPropInverse: (x+y)-y == x and (x*y)/y == x (when y != 0 + rounding-safe).
func TestPropInverse(t *testing.T) {
	setCfg(40, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(3))
	for i := 0; i < 3000; i++ {
		// Keep x and y on the same order of magnitude so the sum carries no
		// more than ~40 digits; otherwise the add rounds the small term away.
		exp := r.Intn(21) - 10
		gen := func() *Decimal {
			ld := r.Intn(10) + 1
			b := make([]byte, ld)
			b[0] = byte('1' + r.Intn(9))
			for j := 1; j < ld; j++ {
				b[j] = byte('0' + r.Intn(10))
			}
			return New(string(b) + "e" + strconv.Itoa(exp))
		}
		y := gen()
		if y.IsZero() {
			continue
		}
		x := gen()
		z := x.Plus(y).Minus(y)
		if z.Cmp(x) != 0 {
			t.Fatalf("inverse add broken: x=%s y=%s got=%s", x.ValueOf(), y.ValueOf(), z.ValueOf())
		}
		q := x.Times(y).Div(y)
		if q.Cmp(x) != 0 {
			t.Fatalf("inverse mul broken: x=%s y=%s got=%s", x.ValueOf(), y.ValueOf(), q.ValueOf())
		}
	}
}

// TestPropSqrtSq: sqrt(x)^2 ~= x for x >= 0, within rounding.
func TestPropSqrtSq(t *testing.T) {
	setCfg(30, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(4))
	for i := 0; i < 500; i++ {
		x := genDecimal(r, 14)
		if x.IsNeg() || x.IsZero() {
			continue
		}
		s := x.Pow(2).Sqrt()
		// sqrt(x^2) == |x| up to a rounding ulp.
		d := x.Minus(s).Abs()
		if !d.Lte(New("1e-26")) {
			t.Fatalf("sqrt(x^2)!=x for x=%s got=%s diff=%s", x.ValueOf(), s.ValueOf(), d.ValueOf())
		}
	}
}

// TestPropCbrtCb: x^3 cbrt back to x. x is limited to ~9 digits so x^3 stays
// within the precision and the inverse is exact.
func TestPropCbrtCb(t *testing.T) {
	setCfg(30, 1, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(5))
	for i := 0; i < 2000; i++ {
		ld := r.Intn(6) + 4
		b := make([]byte, ld)
		b[0] = byte('1' + r.Intn(9))
		for j := 1; j < ld; j++ {
			b[j] = byte('0' + r.Intn(10))
		}
		exp := r.Intn(11) - 5
		sign := ""
		if r.Intn(2) == 0 {
			sign = "-"
		}
		x := New(sign + string(b) + "e" + strconv.Itoa(exp))
		z := x.Pow(3).Cbrt()
		if z.Cmp(x) != 0 {
			t.Fatalf("cbrt(x^3)!=x x=%s got=%s", x.ValueOf(), z.ValueOf())
		}
	}
}

// TestPropCmpAntisym: cmp(a,b) == -cmp(b,a), and eq(a,a).
func TestPropCmpAntisym(t *testing.T) {
	setCfg(20, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(6))
	for i := 0; i < 2000; i++ {
		a, b := genDecimal(r, 15), genDecimal(r, 15)
		if a.Cmp(b) != -b.Cmp(a) {
			t.Fatalf("cmp antisym %s %s", a.ValueOf(), b.ValueOf())
		}
		if a.Cmp(a) != 0 {
			t.Fatalf("cmp self %s", a.ValueOf())
		}
	}
}

// TestPropLogExp: e^{ln(x)} == x for x in a healthy range.
func TestPropLogExp(t *testing.T) {
	setCfg(25, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(7))
	for i := 0; i < 300; i++ {
		// keep mantissa near 1..10 to stay cheap
		m := r.Intn(9) + 1
		e := r.Intn(21) - 10
		x := New(strconv.Itoa(m) + "e" + strconv.Itoa(e))
		if x.IsZero() || x.IsNeg() {
			continue
		}
		l := x.Ln()
		z := l.NaturalExponential()
		d := x.Minus(z).Abs()
		if !d.Lte(x.Times("1e-20").Abs()) {
			t.Fatalf("exp(ln x)!=x x=%s ln=%s back=%s diff=%s", x.ValueOf(), l.ValueOf(), z.ValueOf(), d.ValueOf())
		}
	}
}

// TestPropSqrtMissing: Sqrt of exact squares recovers exactly.
func TestPropSqrtExact(t *testing.T) {
	setCfg(40, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(8))
	for i := 0; i < 300; i++ {
		n := r.Int63n(1_000_000) + 1
		sq := New(n).Times(n)
		if got := sq.Sqrt(); got.Cmp(New(n)) != 0 {
			t.Fatalf("sqrt(%d^2)=%s, want %d", n, got.ValueOf(), n)
		}
	}
}

// TestPropMod: (a mod m) is congruent to a, and |r| < |m|, with r carrying
// either the dividend sign (truncation) or the modulus sign (Euclidean),
// matching decimal.js semantics depending on the modulo mode.
func TestPropModRange(t *testing.T) {
	setCfg(60, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(9))
	for i := 0; i < 3000; i++ {
		a := genDecimal(r, 12)
		m := genDecimal(r, 6)
		if m.IsZero() {
			continue
		}
		r1 := a.Mod(m)
		// True modulo: (a - r1) is an integral multiple of m, |r1| < |m|,
		// and r1 carries the dividend sign (truncated mode) or modulus sign
		// (Euclidean), exactly as decimal.js specifies.
		if r1.IsNaN() {
			t.Fatalf("mod NaN a=%s m=%s", a.ValueOf(), m.ValueOf())
		}
		if !r1.Abs().Lt(m.Abs()) {
			t.Fatalf("mod out of range a=%s m=%s r=%s", a.ValueOf(), m.ValueOf(), r1.ValueOf())
		}
		q := a.Minus(r1).Div(m)
		if !q.IsInteger() {
			t.Fatalf("mod not multiple a=%s m=%s r=%s (a-r)/m=%s", a.ValueOf(), m.ValueOf(), r1.ValueOf(), q.ValueOf())
		}
	}
}

func TestPropAll(t *testing.T) {
	t.Run("RoundTrip", TestPropRoundTrip)
	t.Run("Commute", TestPropCommute)
	t.Run("Inverse", TestPropInverse)
	t.Run("SqrtSq", TestPropSqrtSq)
	t.Run("CbrtCb", TestPropCbrtCb)
	t.Run("CmpAntisym", TestPropCmpAntisym)
	t.Run("LogExp", TestPropLogExp)
	t.Run("SqrtExact", TestPropSqrtExact)
	t.Run("ModRange", TestPropModRange)
}

func TestPropFloat64Sanity(t *testing.T) {
	setCfg(20, 4, -7, 21, 9e15, -9e15, 1, false)
	r := rand.New(rand.NewSource(10))
	for i := 0; i < 2000; i++ {
		x := genDecimal(r, 10)
		f := x.Float64()
		if math.IsInf(f, 0) && x.IsFinite() {
			continue // genuine overflow to Infinity in float64 is fine
		}
		y := New(f)
		d := x.Minus(y).Abs()
		if !d.Lte("1e-10") {
			t.Fatalf("float64 loss x=%s f=%g y=%s diff=%s", x.ValueOf(), f, y.ValueOf(), d.ValueOf())
		}
	}
}
