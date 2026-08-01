package decimal

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	mrand "math/rand/v2"
)

// Constructor-level functions, mirroring the static methods on the
// Decimal constructor in decimal.js (Decimal.abs, Decimal.add, ...).

// Abs returns |x|.
func (c *Constructor) Abs(x any) *Decimal { return c.New(x).Abs() }

// Acos returns the arccosine in radians of x.
func (c *Constructor) Acos(x any) *Decimal { return c.New(x).Acos() }

// Acosh returns the inverse hyperbolic cosine of x.
func (c *Constructor) Acosh(x any) *Decimal { return c.New(x).Acosh() }

// Add returns x + y.
func (c *Constructor) Add(x, y any) *Decimal { return c.New(x).Plus(y) }

// Asin returns the arcsine in radians of x.
func (c *Constructor) Asin(x any) *Decimal { return c.New(x).Asin() }

// Asinh returns the inverse hyperbolic sine of x.
func (c *Constructor) Asinh(x any) *Decimal { return c.New(x).Asinh() }

// Atan returns the arctangent in radians of x.
func (c *Constructor) Atan(x any) *Decimal { return c.New(x).Atan() }

// Atanh returns the inverse hyperbolic tangent of x.
func (c *Constructor) Atanh(x any) *Decimal { return c.New(x).Atanh() }

// Cbrt returns the cube root of x.
func (c *Constructor) Cbrt(x any) *Decimal { return c.New(x).Cbrt() }

// Ceil returns x rounded to an integer using RoundCeil.
func (c *Constructor) Ceil(x any) *Decimal {
	xv := c.New(x)
	return finalise(xv, xv.e+1, RoundCeil, false)
}

// Clamp returns x clamped to the range delineated by min and max.
func (c *Constructor) Clamp(x, min, max any) *Decimal {
	return c.New(x).Clamp(min, max)
}

// Cos returns the cosine of x (radians).
func (c *Constructor) Cos(x any) *Decimal { return c.New(x).Cos() }

// Cosh returns the hyperbolic cosine of x.
func (c *Constructor) Cosh(x any) *Decimal { return c.New(x).Cosh() }

// Div returns x / y.
func (c *Constructor) Div(x, y any) *Decimal { return c.New(x).Div(y) }

// Exp returns e^x.
func (c *Constructor) Exp(x any) *Decimal { return c.New(x).Exp() }

// Floor returns x rounded to an integer using RoundFloor.
func (c *Constructor) Floor(x any) *Decimal {
	xv := c.New(x)
	return finalise(xv, xv.e+1, RoundFloor, false)
}

// Hypot returns the square root of the sum of the squares of the arguments.
func (c *Constructor) Hypot(args ...any) *Decimal {
	t := c.New(0)
	external = false

	for _, a := range args {
		n := c.New(a)
		if n.d == nil {
			if n.s != 0 {
				external = true
				return c.newInf(1)
			}
			t = n
		} else if t.d != nil {
			t = t.Plus(n.Times(n))
		}
	}

	external = true

	return t.Sqrt()
}

// Ln returns the natural logarithm of x.
func (c *Constructor) Ln(x any) *Decimal { return c.New(x).Ln() }

// Log returns the logarithm of x to the base y (default: 10).
func (c *Constructor) Log(x any, y ...any) *Decimal { return c.New(x).Log(y...) }

// Log2 returns the base 2 logarithm of x.
func (c *Constructor) Log2(x any) *Decimal { return c.New(x).Log(2) }

// Log10 returns the base 10 logarithm of x.
func (c *Constructor) Log10(x any) *Decimal { return c.New(x).Log(10) }

func maxOrMin(c *Constructor, args []any, n float64) *Decimal {
	x := c.New(args[0])

	for i := 1; i < len(args); i++ {
		y := c.New(args[i])

		// NaN?
		if y.s == 0 {
			x = y
			break
		}

		k := x.Cmp(y)
		if k == n || (k == 0 && float64(x.s) == n) {
			x = y
		}
	}

	return x
}

// Max returns the maximum of the arguments.
func (c *Constructor) Max(args ...any) *Decimal { return maxOrMin(c, args, -1) }

// Min returns the minimum of the arguments.
func (c *Constructor) Min(args ...any) *Decimal { return maxOrMin(c, args, 1) }

// Mod returns x modulo y.
func (c *Constructor) Mod(x, y any) *Decimal { return c.New(x).Mod(y) }

// Mul returns x * y.
func (c *Constructor) Mul(x, y any) *Decimal { return c.New(x).Times(y) }

// Pow returns x raised to the power y.
func (c *Constructor) Pow(x, y any) *Decimal { return c.New(x).Pow(y) }

// Random returns a new Decimal with a pseudo-random value equal to or
// greater than 0 and less than 1, and with sd, or Precision if sd is
// omitted, significant digits.
func (c *Constructor) Random(sds ...int64) *Decimal {
	var sd int64
	if len(sds) == 0 {
		sd = c.Precision
	} else {
		sd = sds[0]
	}
	checkInt32(sd, 1, maxDigits)

	k := int(math.Ceil(float64(sd) / float64(logBase)))
	rd := make([]int32, 0, k)

	if !c.Crypto {
		for i := 0; i < k; i++ {
			rd = append(rd, int32(mrand.IntN(10000000)))
		}
	} else {
		buf := make([]byte, 4)
		for len(rd) < k {
			if _, err := crand.Read(buf); err != nil {
				panic(err)
			}
			n := binary.LittleEndian.Uint32(buf)
			// 0 <= n < 4294967296; reject n >= 4.29e9 to avoid bias.
			if n >= 4290000000 {
				continue
			}
			// 0 <= (n % 1e7) <= 9999999
			rd = append(rd, int32(n%10000000))
		}
	}

	i := len(rd)
	klast := int64(rd[i-1])
	sdrem := sd % logBase

	// Convert trailing digits to zeros according to sd.
	if klast != 0 && sdrem != 0 {
		n := pow10(logBase - sdrem)
		rd[i-1] = int32(klast / n * n)
	}

	// Remove trailing words which are zero.
	for ; i > 0 && rd[i-1] == 0; i-- {
	}
	rd = rd[:i]

	var e int64
	// Zero?
	if i <= 0 {
		e = 0
		rd = []int32{0}
	} else {
		e = -1
		// Remove leading words which are zero and adjust exponent accordingly.
		for rd[0] == 0 {
			rd = rd[1:]
			e -= logBase
		}

		// Count the digits of the first word of rd to determine leading zeros.
		kk := int64(1)
		for n := rd[0]; n >= 10; n /= 10 {
			kk++
		}

		// Adjust the exponent for leading zeros of the first word of rd.
		if kk < logBase {
			e -= logBase - kk
		}
	}

	r := c.New(1)
	r.e = e
	r.d = rd

	return r
}

// Round returns x rounded to an integer using the constructor's rounding
// mode.
func (c *Constructor) Round(x any) *Decimal {
	xv := c.New(x)
	return finalise(xv, xv.e+1, c.Rounding, false)
}

// Sign returns 1 if x > 0, -1 if x < 0, 0 if x is 0, -0 if x is -0,
// NaN otherwise.
func (c *Constructor) Sign(x any) float64 {
	xv := c.New(x)
	if xv.d == nil {
		if xv.s == 0 {
			return math.NaN()
		}
		return float64(xv.s)
	}
	if xv.d[0] != 0 {
		return float64(xv.s)
	}
	return math.Copysign(0, float64(xv.s))
}

// Sin returns the sine of x (radians).
func (c *Constructor) Sin(x any) *Decimal { return c.New(x).Sin() }

// Sinh returns the hyperbolic sine of x.
func (c *Constructor) Sinh(x any) *Decimal { return c.New(x).Sinh() }

// Sqrt returns the square root of x.
func (c *Constructor) Sqrt(x any) *Decimal { return c.New(x).Sqrt() }

// Sub returns x - y.
func (c *Constructor) Sub(x, y any) *Decimal { return c.New(x).Minus(y) }

// Sum returns the sum of the arguments. Only the result is rounded, not
// the intermediate calculations.
func (c *Constructor) Sum(args ...any) *Decimal {
	if len(args) == 0 {
		panic(errInvalidArgs)
	}
	x := c.New(args[0])

	external = false
	for i := 1; x.s != 0 && i < len(args); i++ {
		x = x.Plus(args[i])
	}
	external = true

	return finalise(x, c.Precision, c.Rounding, false)
}

var errInvalidArgs = errors.New(decimalError + "No arguments")

// Tan returns the tangent of x (radians).
func (c *Constructor) Tan(x any) *Decimal { return c.New(x).Tan() }

// Tanh returns the hyperbolic tangent of x.
func (c *Constructor) Tanh(x any) *Decimal { return c.New(x).Tanh() }

// Trunc returns x truncated to an integer.
func (c *Constructor) Trunc(x any) *Decimal {
	xv := c.New(x)
	return finalise(xv, xv.e+1, RoundDown, false)
}
