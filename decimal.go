package decimal

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// Rounding modes.
const (
	RoundUp        = 0 // Away from zero.
	RoundDown      = 1 // Towards zero.
	RoundCeil      = 2 // Towards +Infinity.
	RoundFloor     = 3 // Towards -Infinity.
	RoundHalfUp    = 4 // Nearest; ties up.
	RoundHalfDown  = 5 // Nearest; ties down.
	RoundHalfEven  = 6 // Nearest; ties to even.
	RoundHalfCeil  = 7 // Nearest; ties ceil.
	RoundHalfFloor = 8 // Nearest; ties floor.

	// Euclid is a modulo mode (not a rounding mode): Euclidean division.
	// q = sign(y) * floor(x / abs(y)); the remainder is always non-negative.
	Euclid = 9
)

const (
	expLimit  = 9e15 // 0 to 9e15
	maxDigits = 1e9  // 0 to 1e9

	base       int64 = 1e7
	logBase    int64 = 7
	maxInteger       = 9007199254740991 // 2^53 - 1
)

const (
	decimalError            = "[DecimalError] "
	invalidArgument         = decimalError + "Invalid argument: "
	precisionLimitExceeded  = decimalError + "Precision limit exceeded"
	errCryptoUnavailableMsg = decimalError + "crypto unavailable"
)

// Internal constants (1025 digits each).
const ln10Str = "2.3025850929940456840179914546843642076011014886287729760333279009675726096773524802359972050895982983419677840422862486334095254650828067566662873690987816894829072083255546808437998948262331985283935053089653777326288461633662222876982198867465436674744042432743651550489343149393914796194044002221051017141748003688084012647080685567743216228355220114804663715659121373450747856947683463616792101806445070648000277502684916746550586856935673420670581136429224554405758925724208241314695689016758940256776311356919292033376587141660230105703089634572075440370847469940168269282808481184289314848524948644871927809676271275775397027668605952496716674183485704422507197965004714951050492214776567636938662976979522110718264549734772662425709429322582798502585509785265383207606726317164309505995087807523710333101197857547331541421808427543863591778117054309827482385045648019095610299291824318237525357709750539565187697510374970888692180205189339507238539205144634197265287286965110862571492198849978748873771345686209167058"
const piStr = "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989380952572010654858632789"

var (
	ln10Precision = int64(len(ln10Str) - 1)
	piPrecision   = int64(len(piStr) - 1)

	ln10Num *Decimal
	piNum   *Decimal
)

// Internal per-constructor state lives on the *Constructor so that each clone
// is safe to use concurrently (mirroring the decimal.js guidance to create a
// cloned constructor per concurrent context). In decimal.js these are module
// globals; Go ports must not share mutable flags across goroutines.
var (
	isBinaryRe  = regexp.MustCompile(`(?i)^0b([01]+(\.[01]*)?|\.[01]+)(p[+-]?\d+)?$`)
	isHexRe     = regexp.MustCompile(`(?i)^0x([0-9a-f]+(\.[0-9a-f]*)?|\.[0-9a-f]+)(p[+-]?\d+)?$`)
	isOctalRe   = regexp.MustCompile(`(?i)^0o([0-7]+(\.[0-7]*)?|\.[0-7]+)(p[+-]?\d+)?$`)
	isDecimalRe = regexp.MustCompile(`(?i)^(\d+(\.\d*)?|\.\d+)(e[+-]?\d+)?$`)
)

const numerals = "0123456789abcdef"

// Constructor holds the configuration of a Decimal constructor, mirroring
// the per-constructor properties (precision, rounding, ...) of decimal.js.
// Use Default for the package-level default constructor, or create new
// constructors with Constructor.Clone.
type Constructor struct {
	// The maximum number of significant digits of the result of a
	// calculation or base conversion. 1 to maxDigits.
	Precision int64
	// The rounding mode used when rounding to Precision. 0 to 8.
	Rounding int64
	// The modulo mode used by Mod. 0 to 9.
	Modulo int64
	// The exponent value at and beneath which String returns exponential
	// notation. 0 to -expLimit.
	ToExpNeg int64
	// The exponent value at and above which String returns exponential
	// notation. 0 to expLimit.
	ToExpPos int64
	// The minimum exponent value, beneath which underflow to zero occurs.
	MinE int64
	// The maximum exponent value, above which overflow to Infinity occurs.
	MaxE int64
	// Whether to use cryptographically-secure random number generation.
	Crypto bool

	// inexact is set by divide when invoked for base conversion. It is
	// per-constructor so clones can be used from separate goroutines.
	inexact bool
	// quadrant is set by toLessThanHalfPi for the trig functions.
	quadrant int
	// external mirrors decimal.js's module-level flag: it is false only while
	// a constructor is mid-operation, suppressing the overflow-underflow
	// conversions applied by New to internal intermediate results.
	external bool
}

// Decimal is an arbitrary-precision decimal floating-point number.
//
// The zero value is not ready for use; obtain values through
// Constructor.New (or the package-level New which uses Default).
type Decimal struct {
	c *Constructor
	s int8    // 1, -1, or 0 for NaN
	e int64   // base 10 exponent of the first word of d (undefined for Inf/NaN)
	d []int32 // coefficient words, base 1e7; nil for Infinity/NaN
}

// Config is the options struct accepted by Constructor.Config and
// Constructor.Clone. Nil pointer fields are left unchanged (or inherited).
// If Defaults is true all fields are first reset to the library defaults.
type Config struct {
	Precision *int64
	Rounding  *int64
	Modulo    *int64
	ToExpNeg  *int64
	ToExpPos  *int64
	MinE      *int64
	MaxE      *int64
	Crypto    *bool
	Defaults  bool
}

// I64 is a helper for building Config pointer fields.
func I64(v int64) *int64 { return &v }

// Bool is a helper for building Config pointer fields.
func Bool(v bool) *bool { return &v }

func defaultConstructor() *Constructor {
	return &Constructor{
		Precision: 20,
		Rounding:  RoundHalfUp,
		Modulo:    RoundDown,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -expLimit,
		MaxE:      expLimit,
		Crypto:    false,
		external:  true,
	}
}

// Default is the package-level default Decimal constructor, mirroring the
// default Decimal export of decimal.js.
var Default = defaultConstructor()

func init() {
	ln10Num = Default.New(ln10Str)
	piNum = Default.New(piStr)
}

// New returns a new Decimal parsed from v using the Default constructor.
// v may be a *Decimal, string, any integer type or float.
func New(v any) *Decimal { return Default.New(v) }

// Clone creates and returns a new constructor with the same configuration
// as c, optionally overridden by cfg.
func (c *Constructor) Clone(cfg *Config) *Constructor {
	nc := defaultConstructor()
	if cfg == nil || !cfg.Defaults {
		nc.Precision = c.Precision
		nc.Rounding = c.Rounding
		nc.Modulo = c.Modulo
		nc.ToExpNeg = c.ToExpNeg
		nc.ToExpPos = c.ToExpPos
		nc.MinE = c.MinE
		nc.MaxE = c.MaxE
		nc.Crypto = c.Crypto
	}
	if cfg != nil {
		nc.Config(cfg)
	}
	return nc
}

// Config applies the given configuration settings to c.
// It panics with a "[DecimalError]" error on invalid values.
func (c *Constructor) Config(cfg *Config) *Constructor {
	if cfg == nil {
		panic(errors.New(decimalError + "Object expected"))
	}
	useDefaults := cfg.Defaults
	def := defaultConstructor()
	set := func(p *int64, v *int64, dv int64, min, max int64, name string) {
		if useDefaults {
			*p = dv
		}
		if v == nil {
			return
		}
		if *v >= min && *v <= max {
			*p = *v
		} else {
			panic(errors.New(invalidArgument + name + ": " + strconv.FormatInt(*v, 10)))
		}
	}
	set(&c.Precision, cfg.Precision, def.Precision, 1, maxDigits, "precision")
	set(&c.Rounding, cfg.Rounding, def.Rounding, 0, 8, "rounding")
	set(&c.ToExpNeg, cfg.ToExpNeg, def.ToExpNeg, -expLimit, 0, "toExpNeg")
	set(&c.ToExpPos, cfg.ToExpPos, def.ToExpPos, 0, expLimit, "toExpPos")
	set(&c.MaxE, cfg.MaxE, def.MaxE, 0, expLimit, "maxE")
	set(&c.MinE, cfg.MinE, def.MinE, -expLimit, 0, "minE")
	set(&c.Modulo, cfg.Modulo, def.Modulo, 0, 9, "modulo")

	if useDefaults {
		c.Crypto = def.Crypto
	}
	if cfg.Crypto != nil {
		c.Crypto = *cfg.Crypto
	}
	return c
}

// Set is an alias of Config.
func (c *Constructor) Set(cfg *Config) *Constructor { return c.Config(cfg) }

// IsDecimal returns true if v is a *Decimal.
func IsDecimal(v any) bool {
	_, ok := v.(*Decimal)
	return ok
}

// New returns a new Decimal whose value is parsed from v, which may be a
// *Decimal (copied), a string (decimal, or 0x/0b/0o-prefixed with optional
// fraction and binary exponent), an integer of any width, or a float.
// It panics with a "[DecimalError]" error for invalid values.
func (c *Constructor) New(v any) *Decimal {
	x := &Decimal{c: c}
	switch v := v.(type) {
	case *Decimal:
		x.s = v.s
		if c.external {
			if v.d == nil || v.e > c.MaxE {
				// Infinity (or NaN).
				x.e = 0
				x.d = nil
			} else if v.e < c.MinE {
				// Zero.
				x.e = 0
				x.d = []int32{0}
			} else {
				x.e = v.e
				x.d = append([]int32(nil), v.d...)
			}
		} else {
			x.e = v.e
			if v.d != nil {
				x.d = append([]int32(nil), v.d...)
			}
		}
		return x
	case string:
		return c.fromString(x, v)
	case int:
		return c.fromFloat64(x, float64(v), true, v)
	case int8:
		return c.fromFloat64(x, float64(v), true, v)
	case int16:
		return c.fromFloat64(x, float64(v), true, v)
	case int32:
		return c.fromFloat64(x, float64(v), true, v)
	case int64:
		return c.fromFloat64(x, float64(v), true, v)
	case uint:
		return c.fromFloat64(x, float64(v), true, v)
	case uint8:
		return c.fromFloat64(x, float64(v), true, v)
	case uint16:
		return c.fromFloat64(x, float64(v), true, v)
	case uint32:
		return c.fromFloat64(x, float64(v), true, v)
	case uint64:
		return c.fromFloat64(x, float64(v), true, v)
	case float32:
		return c.fromFloat64(x, float64(v), false, 0)
	case float64:
		return c.fromFloat64(x, v, false, 0)
	case *big.Int:
		// JS parity: decimal.js parses bigint by its decimal string.
		return c.fromString(x, v.String())
	}
	panic(errors.New(invalidArgument + toDisplayString(v)))
}

func toDisplayString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// fromFloat64 mirrors the JS number path of the constructor.
// If isInt is true, v holds an integer value iv (used for exact parsing of
// 64-bit integers; JS would parse their decimal string).
func (c *Constructor) fromFloat64(x *Decimal, v float64, isInt bool, iv any) *Decimal {
	if v == 0 {
		if math.Signbit(v) {
			x.s = -1
		} else {
			x.s = 1
		}
		x.e = 0
		x.d = []int32{0}
		return x
	}

	if math.IsNaN(v) {
		x.s = 0
		x.e = 0
		x.d = nil
		return x
	}

	if v < 0 {
		x.s = -1
		v = -v
	} else {
		x.s = 1
	}

	if math.IsInf(v, 0) {
		// x.s already set above.
		x.e = 0
		x.d = nil
		return x
	}

	// Fast path for small integers.
	if v == math.Trunc(v) && v < 1e7 {
		n := uint64(v)
		e := int64(0)
		for i := n; i >= 10; i /= 10 {
			e++
		}
		if c.external {
			if e > c.MaxE {
				x.e = 0
				x.d = nil
			} else if e < c.MinE {
				x.e = 0
				x.d = []int32{0}
			} else {
				x.e = e
				x.d = []int32{int32(n)}
			}
		} else {
			x.e = e
			x.d = []int32{int32(n)}
		}
		return x
	}

	if isInt {
		// Larger integers: parse their exact decimal string, mirroring the
		// JS number path (which uses v.toString() after dropping the sign).
		s := intToString(iv)
		if x.s < 0 {
			s = s[1:]
		}
		return parseDecimal(x, s)
	}

	return parseDecimal(x, floatToString(v))
}

func intToString(iv any) string {
	switch n := iv.(type) {
	case int:
		return strconv.Itoa(n)
	case int8:
		return strconv.FormatInt(int64(n), 10)
	case int16:
		return strconv.FormatInt(int64(n), 10)
	case int32:
		return strconv.FormatInt(int64(n), 10)
	case int64:
		return strconv.FormatInt(n, 10)
	case uint:
		return strconv.FormatUint(uint64(n), 10)
	case uint8:
		return strconv.FormatUint(uint64(n), 10)
	case uint16:
		return strconv.FormatUint(uint64(n), 10)
	case uint32:
		return strconv.FormatUint(uint64(n), 10)
	case uint64:
		return strconv.FormatUint(n, 10)
	}
	return ""
}

// floatToString formats a finite, non-zero, positive float64 the way
// JavaScript Number.prototype.toString would (shortest round-trip digits in
// possibly exponential notation). Only the digits and exponent matter, the
// notation differences are normalised by parseDecimal.
func floatToString(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func (c *Constructor) fromString(x *Decimal, v string) *Decimal {
	if len(v) > 0 && v[0] == '-' {
		x.s = -1
		v = v[1:]
	} else {
		if len(v) > 0 && v[0] == '+' {
			v = v[1:]
		}
		x.s = 1
	}
	if isDecimalRe.MatchString(v) {
		return parseDecimal(x, v)
	}
	return parseOther(x, v)
}

func indexE(str string) int {
	for i := 0; i < len(str); i++ {
		if str[i] == 'e' || str[i] == 'E' {
			return i
		}
	}
	return -1
}

// parseExp parses a JS-number-style exponent suffix, clamping on overflow.
func parseExp(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return n
	}
	if len(s) > 0 && s[0] == '-' {
		return math.MinInt64 / 2
	}
	return math.MaxInt64 / 2
}

// Parse the value of a new Decimal x from a decimal string str.
func parseDecimal(x *Decimal, str string) *Decimal {
	e := int64(-1)
	if i := strings.IndexByte(str, '.'); i > -1 {
		str = str[:i] + str[i+1:]
		e = int64(i)
	}

	// Exponential form?
	if i := indexE(str); i > 0 {
		if e < 0 {
			e = int64(i)
		}
		e += parseExp(str[i+1:])
		str = str[:i]
	} else if e < 0 {
		// Integer.
		e = int64(len(str))
	}

	// Determine leading zeros.
	i := 0
	for i < len(str) && str[i] == '0' {
		i++
	}
	// Determine trailing zeros.
	ln := len(str)
	for ln > 0 && str[ln-1] == '0' {
		ln--
	}

	if i < ln {
		str = str[i:ln]
		n := len(str)
		x.e = e - int64(i) - 1

		// Transform base.
		d := make([]int32, 0, n/int(logBase)+1)
		i = int((x.e + 1) % logBase)
		if x.e < 0 {
			i += int(logBase)
		}
		if i < n {
			if i > 0 {
				d = append(d, int32(atoiSafe(str[:i])))
			}
			for n -= int(logBase); i < n; i += int(logBase) {
				d = append(d, int32(atoiSafe(str[i:i+int(logBase)])))
			}
			str = str[i:]
			i = int(logBase) - len(str)
		} else {
			i -= n
		}
		for ; i > 0; i-- {
			str += "0"
		}
		d = append(d, int32(atoiSafe(str)))
		x.d = d

		if x.c.external {
			// Overflow?
			if x.e > x.c.MaxE {
				// Infinity.
				x.d = nil
				x.e = 0
			} else if x.e < x.c.MinE {
				// Zero.
				x.e = 0
				x.d = []int32{0}
			}
		}
	} else {
		// Zero.
		x.e = 0
		x.d = []int32{0}
	}

	return x
}

func atoiSafe(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

// strip underscores in JS numeric-separator style: /(\d)_(?=\d)/g
func stripUnderscores(str string) string {
	var b strings.Builder
	for i := 0; i < len(str); i++ {
		if str[i] == '_' && i > 0 && i+1 < len(str) &&
			str[i-1] >= '0' && str[i-1] <= '9' && str[i+1] >= '0' && str[i+1] <= '9' {
			continue
		}
		b.WriteByte(str[i])
	}
	return b.String()
}

// Parse the value of a new Decimal x from a string str which is not a
// decimal value (Infinity, NaN, or a 0x/0b/0o-prefixed value).
func parseOther(x *Decimal, str string) *Decimal {
	if strings.IndexByte(str, '_') > -1 {
		str = stripUnderscores(str)
		if isDecimalRe.MatchString(str) {
			return parseDecimal(x, str)
		}
	} else if str == "Infinity" || str == "NaN" {
		if str == "NaN" {
			x.s = 0
		}
		x.e = 0
		x.d = nil
		return x
	}

	var baseIn int64
	if isHexRe.MatchString(str) {
		baseIn = 16
		str = strings.ToLower(str)
	} else if isBinaryRe.MatchString(str) {
		baseIn = 2
	} else if isOctalRe.MatchString(str) {
		baseIn = 8
	} else {
		panic(errors.New(invalidArgument + str))
	}

	// Is there a binary exponent part?
	var p int64
	hasP := false
	i := -1
	for j := 0; j < len(str); j++ {
		if str[j] == 'p' || str[j] == 'P' {
			i = j
			break
		}
	}
	if i > 0 {
		p = parseExp(str[i+1:])
		hasP = true
		str = str[2:i]
	} else {
		str = str[2:]
	}

	// Convert str as an integer then divide the result by base raised to a
	// power such that the fraction part will be restored.
	dotIdx := strings.IndexByte(str, '.')
	isFloat := dotIdx >= 0
	c := x.c
	var divisor *Decimal
	ln := 0
	if isFloat {
		str = str[:dotIdx] + str[dotIdx+1:]
		ln = len(str)
		i = ln - dotIdx

		divisor = intPow(c, c.New(baseIn), int64(i), int64(i)*2)
	}

	xd := convertBase(str, baseIn, base)
	xe := int64(len(xd) - 1)

	// Remove trailing zeros.
	i = len(xd) - 1
	for i >= 0 && xd[i] == 0 {
		xd = xd[:i]
		i--
	}
	if i < 0 {
		return c.newZero(x.s)
	}
	x.e = getBase10Exponent(xd, xe)
	x.d = xd
	x.c.external = false

	// 4 * the number of digits of str will always be enough precision for
	// an exact conversion.
	if isFloat {
		x = divide(x, divisor, int64(ln)*4, 0, false, 0)
	}

	// Multiply by the binary exponent part if present.
	if hasP && p != 0 {
		if p > -54 && p < 54 {
			x = x.Times(math.Pow(2, float64(p)))
		} else {
			x = x.Times(Default.Pow(int64(2), p))
		}
	}
	x.c.external = true

	return x
}

// Internal constructors.

func (c *Constructor) newNaN() *Decimal {
	return &Decimal{c: c, s: 0, e: 0, d: nil}
}

func (c *Constructor) newInf(sign int8) *Decimal {
	return &Decimal{c: c, s: sign, e: 0, d: nil}
}

func (c *Constructor) newZero(sign int8) *Decimal {
	return &Decimal{c: c, s: sign, e: 0, d: []int32{0}}
}

// newSign returns a Decimal with no digits assigned yet (used by divide).
func (c *Constructor) newSign(sign int8) *Decimal {
	return &Decimal{c: c, s: sign, e: 0, d: nil}
}

func checkInt32(i, min, max int64) {
	if i < min || i > max {
		panic(errors.New(invalidArgument + strconv.FormatInt(i, 10)))
	}
}

// floorDiv returns floor(a / b) for positive b.
func floorDiv(a, b int64) int64 {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}

// word returns d[i], or 0 if i is out of range (mirrors JS `d[i] || 0`).
func word(d []int32, i int) int32 {
	if i >= 0 && i < len(d) {
		return d[i]
	}
	return 0
}
