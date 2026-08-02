package decimal

// This file provides Go methods that mirror the long-form method names of the
// decimal.js API (e.g. dividedBy, naturalLogarithm, toHexadecimal). Each is a
// documented alias of the idiomatic Go method so that code derived from the
// original library reads naturally in Go too.

// DividedBy returns x / y.
func (x *Decimal) DividedBy(y any) *Decimal { return x.Div(y) }

// DividedToIntegerBy returns the quotient of the division of x by y rounded
// to a whole number.
func (x *Decimal) DividedToIntegerBy(y any) *Decimal { return x.DivToInt(y) }

// Modulo returns x modulo y.
func (x *Decimal) Modulo(y any) *Decimal { return x.Mod(y) }

// Negated returns -x.
func (x *Decimal) Negated() *Decimal { return x.Neg() }

// AbsoluteValue returns |x|.
func (x *Decimal) AbsoluteValue() *Decimal { return x.Abs() }

// Truncated returns the value of x truncated to a whole number.
func (x *Decimal) Truncated() *Decimal { return x.Trunc() }

// ClampedTo returns a new Decimal whose value is x clamped between min and
// max.
func (x *Decimal) ClampedTo(min, max any) *Decimal { return x.Clamp(min, max) }

// ComparedTo compares x and y.
func (x *Decimal) ComparedTo(y any) float64 { return x.Cmp(y) }

// Equals returns true if x == y.
func (x *Decimal) Equals(y any) bool { return x.Eq(y) }

// GreaterThan returns true if x > y.
func (x *Decimal) GreaterThan(y any) bool { return x.Gt(y) }

// GreaterThanOrEqualTo returns true if x >= y.
func (x *Decimal) GreaterThanOrEqualTo(y any) bool { return x.Gte(y) }

// LessThan returns true if x < y.
func (x *Decimal) LessThan(y any) bool { return x.Lt(y) }

// LessThanOrEqualTo returns true if x <= y.
func (x *Decimal) LessThanOrEqualTo(y any) bool { return x.Lte(y) }

// IsInteger returns true if x is an integer.
func (x *Decimal) IsInteger() bool { return x.IsInt() }

// IsNegative returns true if x is negative.
func (x *Decimal) IsNegative() bool { return x.IsNeg() }

// IsPositive returns true if x is positive.
func (x *Decimal) IsPositive() bool { return x.IsPos() }

// DecimalPlaces returns the number of decimal places of x.
func (x *Decimal) DecimalPlaces() float64 { return x.Dp() }

// Precision returns the number of significant digits of x. If z is true,
// trailing integer zeros are counted.
func (x *Decimal) Precision(z ...bool) float64 { return x.Sd(z...) }

// NaturalExponential returns e^x.
func (x *Decimal) NaturalExponential() *Decimal { return x.Exp() }

// NaturalLogarithm returns the natural logarithm (base e) of x.
func (x *Decimal) NaturalLogarithm() *Decimal { return x.Ln() }

// Logarithm returns the logarithm of x to the given base (default 10).
func (x *Decimal) Logarithm(bases ...any) *Decimal { return x.Log(bases...) }

// SquareRoot returns the square root of x.
func (x *Decimal) SquareRoot() *Decimal { return x.Sqrt() }

// CubeRoot returns the cube root of x.
func (x *Decimal) CubeRoot() *Decimal { return x.Cbrt() }

// ToPower returns x raised to the power y.
func (x *Decimal) ToPower(y any) *Decimal { return x.Pow(y) }

// Sine returns the sine of x in radians.
func (x *Decimal) Sine() *Decimal { return x.Sin() }

// Cosine returns the cosine of x in radians.
func (x *Decimal) Cosine() *Decimal { return x.Cos() }

// Tangent returns the tangent of x in radians.
func (x *Decimal) Tangent() *Decimal { return x.Tan() }

// HyperbolicSine returns the hyperbolic sine of x.
func (x *Decimal) HyperbolicSine() *Decimal { return x.Sinh() }

// HyperbolicCosine returns the hyperbolic cosine of x.
func (x *Decimal) HyperbolicCosine() *Decimal { return x.Cosh() }

// HyperbolicTangent returns the hyperbolic tangent of x.
func (x *Decimal) HyperbolicTangent() *Decimal { return x.Tanh() }

// InverseSine returns the arcsine of x in radians.
func (x *Decimal) InverseSine() *Decimal { return x.Asin() }

// InverseCosine returns the arccosine of x in radians.
func (x *Decimal) InverseCosine() *Decimal { return x.Acos() }

// InverseTangent returns the arctangent of x in radians.
func (x *Decimal) InverseTangent() *Decimal { return x.Atan() }

// InverseHyperbolicSine returns the inverse hyperbolic sine of x.
func (x *Decimal) InverseHyperbolicSine() *Decimal { return x.Asinh() }

// InverseHyperbolicCosine returns the inverse hyperbolic cosine of x.
func (x *Decimal) InverseHyperbolicCosine() *Decimal { return x.Acosh() }

// InverseHyperbolicTangent returns the inverse hyperbolic tangent of x.
func (x *Decimal) InverseHyperbolicTangent() *Decimal { return x.Atanh() }

// ToDecimalPlaces returns a new Decimal rounded to dp decimal places using
// rounding mode rm (default: Constructor.Rounding).
func (x *Decimal) ToDecimalPlaces(args ...int64) *Decimal { return x.ToDP(args...) }

// ToHexadecimal returns the hexadecimal representation of x to sd significant
// digits (default: Constructor.Precision).
func (x *Decimal) ToHexadecimal(sd ...int64) string { return x.ToHex(sd...) }

// ToSignificantDigits returns x rounded to sd significant digits using
// rounding mode rm (default: Constructor.Rounding).
func (x *Decimal) ToSignificantDigits(args ...int64) *Decimal { return x.ToSD(args...) }

// ToJSON returns the JSON-compatible string representation of x, as decimal.js
// defines toJSON as an alias of valueOf.
func (x *Decimal) ToJSON() string { return x.ValueOf() }

// ToString returns the string representation of x. Unlike ValueOf, for a
// negative zero the minus sign is omitted, mirroring decimal.js toString.
func (x *Decimal) ToString() string { return x.String() }

// MarshalJSON implements json.Marshaler: the value is marshalled as a JSON
// string matching the valueOf() representation (decimal.js toJSON behaviour).
func (x *Decimal) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.ValueOf() + `"`), nil
}
