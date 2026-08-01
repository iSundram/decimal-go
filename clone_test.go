package decimal

import "testing"

func TestClone(t *testing.T) {
	setCfg(10, 4, -7, 21, 9e15, -9e15, 1, false)

	D1 := Default.Clone(&Config{Precision: i64(1)})
	D2 := Default.Clone(&Config{Precision: i64(2)})
	D3 := Default.Clone(&Config{Precision: i64(3)})
	D4 := Default.Clone(&Config{Precision: i64(4)})
	D5 := Default.Clone(&Config{Precision: i64(5)})
	D6 := Default.Clone(&Config{Precision: i64(6)})
	D7 := Default.Clone(&Config{Precision: i64(7)})
	D8 := Default.Clone(nil)
	D8.Config(&Config{Precision: i64(8)})
	D9 := Default.Clone(&Config{Precision: i64(9)})

	// JS: T.assert(Decimal.prototype === D9.prototype) — constructors in Go
	// share the same (single) method set, so there is no prototype analog.
	assert(t, Default != D9) // JS: Decimal !== D9

	x := New(5)
	x1 := D1.New(5)
	x2 := D2.New(5)
	x3 := D3.New(5)
	x4 := D4.New(5)
	x5 := D5.New(5)
	x6 := D6.New(5)
	x7 := D7.New(5)
	x8 := D8.New(5)
	x9 := D9.New(5)

	assert(t, x1.Div(3).Eq(2.0))
	assert(t, x2.Div(3).Eq(1.7))
	assert(t, x3.Div(3).Eq(1.67))
	assert(t, x4.Div(3).Eq(1.667))
	assert(t, x5.Div(3).Eq(1.6667))
	assert(t, x6.Div(3).Eq(1.66667))
	assert(t, x7.Div(3).Eq(1.666667))
	assert(t, x8.Div(3).Eq(1.6666667))
	assert(t, x9.Div(3).Eq(1.66666667))
	assert(t, x.Div(3).Eq(1.666666667))

	y := New(3)
	y1 := D1.New(3)
	y2 := D2.New(3)
	y3 := D3.New(3)
	y4 := D4.New(3)
	y5 := D5.New(3)
	y6 := D6.New(3)
	y7 := D7.New(3)
	y8 := D8.New(3)
	y9 := D9.New(3)

	assert(t, x1.Div(y1).Eq(2.0))
	assert(t, x2.Div(y2).Eq(1.7))
	assert(t, x3.Div(y3).Eq(1.67))
	assert(t, x4.Div(y4).Eq(1.667))
	assert(t, x5.Div(y5).Eq(1.6667))
	assert(t, x6.Div(y6).Eq(1.66667))
	assert(t, x7.Div(y7).Eq(1.666667))
	assert(t, x8.Div(y8).Eq(1.6666667))
	assert(t, x9.Div(y9).Eq(1.66666667))
	assert(t, x.Div(y).Eq(1.666666667))

	assert(t, x1.Div(y9).Eq(2.0))
	assert(t, x2.Div(y8).Eq(1.7))
	assert(t, x3.Div(y7).Eq(1.67))
	assert(t, x4.Div(y6).Eq(1.667))
	assert(t, x5.Div(y5).Eq(1.6667))
	assert(t, x6.Div(y4).Eq(1.66667))
	assert(t, x7.Div(y3).Eq(1.666667))
	assert(t, x8.Div(y2).Eq(1.6666667))
	assert(t, x9.Div(y1).Eq(1.66666667))

	assert(t, Default.Precision == 10)
	assert(t, D9.Precision == 9)
	assert(t, D8.Precision == 8)
	assert(t, D7.Precision == 7)
	assert(t, D6.Precision == 6)
	assert(t, D5.Precision == 5)
	assert(t, D4.Precision == 4)
	assert(t, D3.Precision == 3)
	assert(t, D2.Precision == 2)
	assert(t, D1.Precision == 1)

	assert(t, New(9.99).Eq(D5.New("9.99")))
	assert(t, !New(9.99).Eq(D3.New("-9.99")))
	assert(t, !New(123.456789).ToSD().Eq(D3.New("123.456789").ToSD()))
	assert(t, New(123.456789).Round().Eq(D3.New("123.456789").Round()))

	assert(t, New(1).c == New(1).c) // constructor === constructor
	assert(t, D9.New(1).c == D9.New(1).c)
	assert(t, New(1).c != D1.New(1).c)
	assert(t, D8.New(1).c != D9.New(1).c)

	// JS: T.assertException(function () { Decimal.clone(null) }, ...)
	// is unportable: passing nil to Clone is the valid no-override form
	// in Go, and no other value can be passed (static typing).

	// defaults: true

	Default.Config(&Config{
		Precision: i64(100),
		Rounding:  i64(2),
		ToExpNeg:  i64(-100),
		ToExpPos:  i64(200),
		Defaults:  true,
	})

	assert(t, Default.Precision == 100)
	assert(t, Default.Rounding == 2)
	assert(t, Default.ToExpNeg == -100)
	assert(t, Default.ToExpPos == 200)
	// JS: t(Decimal.defaults === undefined) — `defaults` is not a stored
	// config property; there is no analog in Go.

	D1 = Default.Clone(&Config{Defaults: true})

	assert(t, D1.Precision == 20)
	assert(t, D1.Rounding == 4)
	assert(t, D1.ToExpNeg == -7)
	assert(t, D1.ToExpPos == 21)
	// JS: t(D1.defaults === undefined) — no analog in Go.

	D2 = Default.Clone(&Config{Defaults: true, Rounding: i64(5)})

	assert(t, D2.Precision == 20)
	assert(t, D2.Rounding == 5)
	assert(t, D2.ToExpNeg == -7)
	assert(t, D2.ToExpPos == 21)

	D3 = Default.Clone(&Config{Defaults: false})

	assert(t, D3.Rounding == 2)
}
