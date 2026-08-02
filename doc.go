/*
Package decimal provides an arbitrary-precision decimal floating-point type for
Go.

It is a port of the decimal.js library (v10.6.0) and preserves its behavior —
rounding modes, error conditions and string formatting rules — while exposing
an idiomatic Go API. Behavioral parity with the original library is verified
byte-for-byte by the cross-validation harness in xvalidate/ and by the ported
test suite.

# Representation

The value of a Decimal is:

	sign * coefficient * 10^exponent

where the coefficient is stored as a slice of base 1e7 words. A nil digits
slice represents Infinity (sign != 0) or NaN (sign == 0). All Decimal values
are immutable: every operation returns a new value and never mutates its
operands.

# Constructing values

The package-level New function and the Default constructor parse strings,
integers, floats, *big.Int values and other Decimals:

	x := decimal.New("123.456")           // from string
	y := decimal.New(42)                   // from integer
	z := x.Plus(y)                         // 165.456

Invalid values cause a panic with a message prefixed by "[DecimalError]",
mirroring decimal.js error semantics.

# Configuration

The Default constructor applies the decimal.js defaults (precision 20,
half-up rounding, toExpNeg -7, toExpPos 21). A cloned constructor with
different settings is created with Constructor.Clone, or an existing one can
be reconfigured with Constructor.Config:

	c := decimal.Default.Clone(&decimal.Config{Precision: decimal.I64(50)})
	pi := c.New("3.14159...") // rounded to 50 significant digits

The rounding modes are exposed as rounding-mode constants (RoundUp,
RoundDown, RoundCeil, RoundFloor, RoundHalfUp, RoundHalfDown, RoundHalfEven,
RoundHalfCeil, RoundHalfFloor); Euclid selects Euclidean division in the
modulo operation.

# Concurrency

A Decimal instance is safe for concurrent read-only use. Constructors carry
mutable state (rounding results, intermediate flags), so follow the
decimal.js guidance: give each concurrent context its own cloned constructor.
A single *Constructor shared without synchronisation is not supported.

# API compatibility

Methods follow idiomatic Go name() names (Plus, Minus, Times, Div, Pow,
Sqrt, Sin, ...). The decimal.js long-form names (DividedBy, NaturalLogarithm,
SquareRoot, ...), ToJSON and MarshalJSON are provided as aliases in parity.go.
String returns the decimal.js toString() form; ValueOf returns the valueOf()
form (which, unlike String, keeps the leading minus sign of -0).
*/
package decimal
