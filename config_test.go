package decimal

import "testing"

func TestConfig(t *testing.T) {
	// MAX_DIGITS = maxDigits (1e9); EXP_LIMIT = expLimit (9e15).

	/*
	  precision  {number} [1, MAX_DIGITS]
	  rounding   {number} [0, 8]
	  toExpNeg   {number} [-EXP_LIMIT, 0]
	  toExpPos   {number} [0, EXP_LIMIT]
	  maxE       {number} [0, EXP_LIMIT]
	  minE       {number} [-EXP_LIMIT, 0]
	  modulo     {number} [0, 9]
	  crypto     {boolean|number} {true, false, 1, 0}
	*/

	assert(t, Default.Config(&Config{}) == Default)

	tx := func(f func(), msg string) {
		t.Helper()
		assertException(t, f, msg)
	}

	// JS passes non-object values here; in Go a nil *Config exercises the same
	// "Object expected" error path (the scalar/string cases are rejected by
	// Go's static typing at compile time).
	tx(func() { Default.Config(nil) }, "config()")
	tx(func() { Default.Config(nil) }, "config(null)")
	tx(func() { Default.Config(nil) }, "config(undefined)")
	tx(func() { Default.Config(nil) }, "config(0)")
	tx(func() { Default.Config(nil) }, "config('')")
	tx(func() { Default.Config(nil) }, "config('hi')")
	tx(func() { Default.Config(nil) }, "config('123')")

	Default.Config(&Config{
		Precision: i64(20),
		Rounding:  i64(4),
		ToExpNeg:  i64(-7),
		ToExpPos:  i64(21),
		MinE:      i64(-9e15),
		MaxE:      i64(9e15),
		Crypto:    bptr(false),
		Modulo:    i64(1),
	})

	assert(t, Default.Precision == 20)
	assert(t, Default.Rounding == 4)
	assert(t, Default.ToExpNeg == -7)
	assert(t, Default.ToExpPos == 21)
	assert(t, Default.MinE == -9e15)
	assert(t, Default.MaxE == 9e15)
	assert(t, Default.Crypto == false)
	assert(t, Default.Modulo == 1)

	Default.Config(&Config{
		Precision: i64(40),
		Rounding:  i64(4),
		ToExpNeg:  i64(-1000),
		ToExpPos:  i64(1000),
		MinE:      i64(-1e9),
		MaxE:      i64(1e9),
		//Crypto: bptr(true),    // requires crypto object
		Modulo: i64(4),
	})

	assert(t, Default.Precision == 40)
	assert(t, Default.Rounding == 4)
	assert(t, Default.ToExpNeg == -1000)
	assert(t, Default.ToExpPos == 1000)
	assert(t, Default.MinE == -1e9)
	assert(t, Default.MaxE == 1e9)
	//assert(t, Default.Crypto == true);    // requires crypto object
	assert(t, Default.Modulo == 4)

	Default.Config(&Config{
		ToExpNeg: i64(-7),
		ToExpPos: i64(21),
		MinE:     i64(-324),
		MaxE:     i64(308),
	})

	assert(t, Default.ToExpNeg == -7)
	assert(t, Default.ToExpPos == 21)
	assert(t, Default.MinE == -324)
	assert(t, Default.MaxE == 308)

	// precision

	var tf func(expected int64, cfg *Config)

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.Precision)
	}

	tf(1, &Config{Precision: i64(1)})
	tf(1, &Config{}) // {precision: void 0}
	tf(20, &Config{Precision: i64(20)})
	tf(300000, &Config{Precision: i64(300000)})
	tf(4e8, &Config{Precision: i64(4e8)})
	tf(1e9, &Config{Precision: i64(1e9)})
	tf(maxDigits, &Config{Precision: i64(maxDigits)})

	tx(func() { Default.Config(&Config{Precision: i64(0)}) }, "precision: 0")
	tx(func() { Default.Config(&Config{Precision: i64(maxDigits + 1)}) }, "precision: MAX_DIGITS + 1")
	// Not expressible in Go (Config fields are *int64): MAX_DIGITS + 0.1.
	tx(func() { Default.Config(&Config{Precision: i64(maxDigits + 1)}) }, "precision: MAX_DIGITS + 0.1")
	// Not expressible in Go (strings are rejected at compile time):
	// tx(... {precision: '0'}, "precision: '0'")
	// tx(... {precision: '1'}, "precision: '1'")
	// tx(... {precision: '123456789'}, "precision: '123456789'")
	tx(func() { Default.Config(&Config{Precision: i64(-1)}) }, "precision: -1")
	// Not expressible in Go (non-integer numbers):
	// tx(... {precision: 0.1}, "precision: 0.1")
	// tx(... {precision: 1.1}, "precision: 1.1")
	// tx(... {precision: -1.1}, "precision: -1.1")
	// tx(... {precision: 8.1}, "precision: 8.1")
	tx(func() { Default.Config(&Config{Precision: i64(1e9 + 1)}) }, "precision: 1e9 + 1")
	// Not expressible in Go: {precision: []}, {precision: {}}, {precision: ''},
	// {precision: 'hi'}, {precision: '1e+9'}, {precision: null}, {precision: NaN},
	// {precision: Infinity} (nil means "leave unchanged"; NaN/Infinity are not int64).

	tf(maxDigits, &Config{}) // {precision: void 0}

	// rounding

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.Rounding)
	}

	tf(4, &Config{}) // {rounding: void 0}
	tf(0, &Config{Rounding: i64(0)})
	tf(1, &Config{Rounding: i64(1)})
	tf(2, &Config{Rounding: i64(2)})
	tf(3, &Config{Rounding: i64(3)})
	tf(4, &Config{Rounding: i64(4)})
	tf(5, &Config{Rounding: i64(5)})
	tf(6, &Config{Rounding: i64(6)})
	tf(7, &Config{Rounding: i64(7)})
	tf(8, &Config{Rounding: i64(8)})

	tx(func() { Default.Config(&Config{Rounding: i64(-1)}) }, "rounding : -1")
	// Not expressible in Go (non-integer numbers):
	// tx(... {rounding: 0.1}, "rounding : 0.1")
	// tx(... {rounding: 8.1}, "rounding : 8.1")
	tx(func() { Default.Config(&Config{Rounding: i64(9)}) }, "rounding : 9")
	// Not expressible in Go (strings):
	// tx(... {rounding: '0'}, "rounding: '0'")
	// tx(... {rounding: '1'}, "rounding: '1'")
	// tx(... {rounding: '123456789'}, "rounding: '123456789'")
	// tx(... {rounding: 1.1}, "rounding : 1.1")
	// tx(... {rounding: -1.1}, "rounding : -1.1")
	tx(func() { Default.Config(&Config{Rounding: i64(11)}) }, "rounding : 11")
	// Not expressible in Go: {rounding: []}, {rounding: {}}, {rounding: ''},
	// {rounding: 'hi'}, {rounding: null}, {rounding: NaN}, {rounding: Infinity}.

	tf(8, &Config{}) // {rounding: void 0}

	// toExpNeg

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.ToExpNeg)
	}

	tf(-7, &Config{}) // {toExpNeg: void 0}
	tf(0, &Config{ToExpNeg: i64(0)})
	tf(-1, &Config{ToExpNeg: i64(-1)})
	tf(-999, &Config{ToExpNeg: i64(-999)})
	tf(-5675367, &Config{ToExpNeg: i64(-5675367)})
	tf(-98770170790791, &Config{ToExpNeg: i64(-98770170790791)})
	tf(-expLimit, &Config{ToExpNeg: i64(-expLimit)})

	tx(func() { Default.Config(&Config{ToExpNeg: i64(-expLimit - 1)}) }, "-EXP_LIMIT - 1")
	// Not expressible in Go: {toExpNeg: '-7'}, {toExpNeg: -0.1}, {toExpNeg: 0.1}.
	tx(func() { Default.Config(&Config{ToExpNeg: i64(1)}) }, "toExpNeg: 1")
	// Not expressible in Go: {toExpNeg: -Infinity}, {toExpNeg: NaN},
	// {toExpNeg: null}, {toExpNeg: {}}.

	tf(-expLimit, &Config{}) // {toExpNeg: void 0}

	// toExpPos

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.ToExpPos)
	}

	tf(21, &Config{}) // {toExpPos: void 0}
	tf(0, &Config{ToExpPos: i64(0)})
	tf(1, &Config{ToExpPos: i64(1)})
	tf(999, &Config{ToExpPos: i64(999)})
	tf(5675367, &Config{ToExpPos: i64(5675367)})
	tf(98770170790791, &Config{ToExpPos: i64(98770170790791)})
	tf(expLimit, &Config{ToExpPos: i64(expLimit)})

	tx(func() { Default.Config(&Config{ToExpPos: i64(expLimit + 1)}) }, "EXP_LIMIT + 1")
	// Not expressible in Go: {toExpPos: '21'}, {toExpPos: -0.1}, {toExpPos: 0.1}.
	tx(func() { Default.Config(&Config{ToExpPos: i64(-1)}) }, "toExpPos: -1")
	// Not expressible in Go: {toExpPos: Infinity}, {toExpPos: NaN},
	// {toExpPos: null}, {toExpPos: {}}.

	tf(expLimit, &Config{}) // {toExpPos: void 0}

	// maxE

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.MaxE)
	}

	tf(308, &Config{}) // {maxE: void 0}
	tf(0, &Config{MaxE: i64(0)})
	tf(1, &Config{MaxE: i64(1)})
	tf(999, &Config{MaxE: i64(999)})
	tf(5675367, &Config{MaxE: i64(5675367)})
	tf(98770170790791, &Config{MaxE: i64(98770170790791)})
	tf(expLimit, &Config{MaxE: i64(expLimit)})

	tx(func() { Default.Config(&Config{MaxE: i64(expLimit + 1)}) }, "EXP_LIMIT + 1")
	// Not expressible in Go: {maxE: '308'}, {maxE: -0.1}, {maxE: 0.1}.
	tx(func() { Default.Config(&Config{MaxE: i64(-1)}) }, "maxE: -1")
	// Not expressible in Go: {maxE: Infinity}, {maxE: NaN}, {maxE: null}, {maxE: {}}.

	tf(expLimit, &Config{}) // {maxE: void 0}

	// minE

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.MinE)
	}

	tf(-324, &Config{}) // {minE: void 0}
	tf(0, &Config{MinE: i64(0)})
	tf(-1, &Config{MinE: i64(-1)})
	tf(-999, &Config{MinE: i64(-999)})
	tf(-5675367, &Config{MinE: i64(-5675367)})
	tf(-98770170790791, &Config{MinE: i64(-98770170790791)})
	tf(-expLimit, &Config{MinE: i64(-expLimit)})

	tx(func() { Default.Config(&Config{MinE: i64(-expLimit - 1)}) }, "-EXP_LIMIT - 1")
	// Not expressible in Go: {minE: '-324'}, {minE: -0.1}, {minE: 0.1}.
	tx(func() { Default.Config(&Config{MinE: i64(1)}) }, "minE: 1")
	// Not expressible in Go: {minE: -Infinity}, {minE: NaN}, {minE: null}, {minE: {}}.

	tf(-expLimit, &Config{}) // {minE: void 0}

	// crypto

	tfc := func(expected bool, cfg *Config) {
		Default.Config(cfg)
		assert(t, Default.Crypto == expected)
	}

	tfc(false, &Config{}) // {crypto: void 0}
	// JS crypto: 0 means false.
	tfc(false, &Config{Crypto: bptr(false)}) // {crypto: 0}
	//tfc(true, &Config{Crypto: bptr(true)});    // {crypto: 1}; requires crypto object
	tfc(false, &Config{Crypto: bptr(false)})
	//tfc(true, &Config{Crypto: bptr(true)});    // {crypto: true}; requires crypto object

	// Not expressible in Go (Crypto is *bool):
	// tx(... {crypto: 'hiya'}, "crypto: 'hiya'")
	// tx(... {crypto: 'true'}, "crypto: 'true'")
	// tx(... {crypto: 'false'}, "crypto: 'false'")
	// tx(... {crypto: '0'}, "crypto: '0'")
	// tx(... {crypto: '1'}, "crypto: '1'")
	// tx(... {crypto: -1}, "crypto: -1")
	// tx(... {crypto: 0.1}, "crypto: 0.1")
	// tx(... {crypto: 1.1}, "crypto: 1.1")
	// tx(... {crypto: []}, "crypto: []")
	// tx(... {crypto: {}}, "crypto: {}")
	// tx(... {crypto: ''}, "crypto: ''")
	// tx(... {crypto: NaN}, "crypto: NaN")
	// tx(... {crypto: Infinity}, "crypto: Infinity")

	assert(t, Default.Crypto == false)

	// modulo

	tf = func(expected int64, cfg *Config) {
		Default.Config(cfg)
		assertEq(t, expected, Default.Modulo)
	}

	tf(4, &Config{}) // {modulo: void 0}
	tf(0, &Config{Modulo: i64(0)})
	tf(1, &Config{Modulo: i64(1)})
	tf(2, &Config{Modulo: i64(2)})
	tf(3, &Config{Modulo: i64(3)})
	tf(4, &Config{Modulo: i64(4)})
	tf(5, &Config{Modulo: i64(5)})
	tf(6, &Config{Modulo: i64(6)})
	tf(7, &Config{Modulo: i64(7)})
	tf(8, &Config{Modulo: i64(8)})
	tf(9, &Config{Modulo: i64(9)})

	tx(func() { Default.Config(&Config{Modulo: i64(-1)}) }, "modulo: -1")
	// Not expressible in Go (non-integer numbers): {modulo: 0.1}, {modulo: 9.1},
	// {modulo: 1.1}, {modulo: -1.1}.
	tx(func() { Default.Config(&Config{Modulo: i64(10)}) }, "modulo: 10")
	tx(func() { Default.Config(&Config{Modulo: i64(-11)}) }, "modulo: -11")
	// Not expressible in Go: {modulo: '0'}, {modulo: '1'}, {modulo: []},
	// {modulo: {}}, {modulo: ''}, {modulo: ' '}, {modulo: 'hi'}, {modulo: null},
	// {modulo: NaN}, {modulo: Infinity}.

	tf(9, &Config{}) // {modulo: void 0}

	// defaults

	Default.Config(&Config{
		Precision: i64(100),
		Rounding:  i64(2),
		ToExpNeg:  i64(-100),
		ToExpPos:  i64(200),
	})

	assert(t, Default.Precision == 100)

	Default.Config(&Config{Defaults: true})

	assert(t, Default.Precision == 20)
	assert(t, Default.Rounding == 4)
	assert(t, Default.ToExpNeg == -7)
	assert(t, Default.ToExpPos == 21)
	// JS: Decimal.defaults === undefined — the Go Config struct has no such field.

	Default.Rounding = 3

	Default.Config(&Config{Precision: i64(50), Defaults: true})

	assert(t, Default.Precision == 50)
	assert(t, Default.Rounding == 4)

	// Decimal.set is an alias for Decimal.config (function identity is not
	// testable in Go; verify Set behaves identically instead).
	Default.Set(&Config{Precision: i64(100)})
	assertEq(t, int64(100), Default.Precision)
	Default.Set(&Config{Defaults: true})
	assertEq(t, int64(20), Default.Precision)
}
