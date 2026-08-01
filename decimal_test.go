package decimal

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestDecimal(t *testing.T) {
	setCfg(40, 4, -9e15, 9e15, 9e15, -9e15, 1, false)

	tf := func(coefficient []int, exponent int64, sign int8, n any) {
		t.Helper()
		assertEqProps(t, coefficient, exponent, sign, New(n))
	}

	tf([]int{0}, 0, 1, 0.0)
	tf([]int{0}, 0, -1, math_NegZero())
	tf([]int{1}, 0, -1, -1.0)
	tf([]int{10}, 1, -1, -10.0)

	tf([]int{1}, 0, 1, 1.0)
	tf([]int{10}, 1, 1, 10.0)
	tf([]int{100}, 2, 1, 100.0)
	tf([]int{1000}, 3, 1, 1000.0)
	tf([]int{10000}, 4, 1, 10000.0)
	tf([]int{100000}, 5, 1, 100000.0)
	tf([]int{1000000}, 6, 1, 1000000.0)

	tf([]int{1}, 7, 1, 10000000.0)
	tf([]int{10}, 8, 1, 100000000.0)
	tf([]int{100}, 9, 1, 1000000000.0)
	tf([]int{1000}, 10, 1, 10000000000.0)
	tf([]int{10000}, 11, 1, 100000000000.0)
	tf([]int{100000}, 12, 1, 1000000000000.0)
	tf([]int{1000000}, 13, 1, 10000000000000.0)

	tf([]int{1}, 14, -1, -100000000000000.0)
	tf([]int{10}, 15, -1, -1000000000000000.0)
	tf([]int{100}, 16, -1, -10000000000000000.0)
	tf([]int{1000}, 17, -1, -100000000000000000.0)
	tf([]int{10000}, 18, -1, -1000000000000000000.0)
	tf([]int{100000}, 19, -1, -10000000000000000000.0)
	tf([]int{1000000}, 20, -1, -100000000000000000000.0)

	tf([]int{1000000}, -1, 1, 1e-1)
	tf([]int{100000}, -2, -1, -1e-2)
	tf([]int{10000}, -3, 1, 1e-3)
	tf([]int{1000}, -4, -1, -1e-4)
	tf([]int{100}, -5, 1, 1e-5)
	tf([]int{10}, -6, -1, -1e-6)
	tf([]int{1}, -7, 1, 1e-7)

	tf([]int{1000000}, -8, 1, 1e-8)
	tf([]int{100000}, -9, -1, -1e-9)
	tf([]int{10000}, -10, 1, 1e-10)
	tf([]int{1000}, -11, -1, -1e-11)
	tf([]int{100}, -12, 1, 1e-12)
	tf([]int{10}, -13, -1, -1e-13)
	tf([]int{1}, -14, 1, 1e-14)

	tf([]int{1000000}, -15, 1, 1e-15)
	tf([]int{100000}, -16, -1, -1e-16)
	tf([]int{10000}, -17, 1, 1e-17)
	tf([]int{1000}, -18, -1, -1e-18)
	tf([]int{100}, -19, 1, 1e-19)
	tf([]int{10}, -20, -1, -1e-20)
	tf([]int{1}, -21, 1, 1e-21)

	tf([]int{9}, 0, 1, "9")
	tf([]int{99}, 1, -1, "-99")
	tf([]int{999}, 2, 1, "999")
	tf([]int{9999}, 3, -1, "-9999")
	tf([]int{99999}, 4, 1, "99999")
	tf([]int{999999}, 5, -1, "-999999")
	tf([]int{9999999}, 6, 1, "9999999")

	tf([]int{9, 9999999}, 7, -1, "-99999999")
	tf([]int{99, 9999999}, 8, 1, "999999999")
	tf([]int{999, 9999999}, 9, -1, "-9999999999")
	tf([]int{9999, 9999999}, 10, 1, "99999999999")
	tf([]int{99999, 9999999}, 11, -1, "-999999999999")
	tf([]int{999999, 9999999}, 12, 1, "9999999999999")
	tf([]int{9999999, 9999999}, 13, -1, "-99999999999999")

	tf([]int{9, 9999999, 9999999}, 14, 1, "999999999999999")
	tf([]int{99, 9999999, 9999999}, 15, -1, "-9999999999999999")
	tf([]int{999, 9999999, 9999999}, 16, 1, "99999999999999999")
	tf([]int{9999, 9999999, 9999999}, 17, -1, "-999999999999999999")
	tf([]int{99999, 9999999, 9999999}, 18, 1, "9999999999999999999")
	tf([]int{999999, 9999999, 9999999}, 19, -1, "-99999999999999999999")
	tf([]int{9999999, 9999999, 9999999}, 20, 1, "999999999999999999999")

	// Test base conversion.

	t2 := func(expected string, n any) {
		t.Helper()
		assertEq(t, expected, New(n).ValueOf())
	}

	randInt := func() int64 {
		return int64(math.Floor(rand.Float64() * 0x20000000000000 / math.Pow(10, float64(int64(rand.Float64()*16)))))
	}

	// Test random integers against Number.prototype.toString(base).
	for i, k := 0, int64(0); i < 127; i++ {
		k = randInt()
		t2(strconv.FormatInt(k, 10), "0b"+strconv.FormatInt(k, 2))
		k = randInt()
		t2(strconv.FormatInt(k, 10), "0B"+strconv.FormatInt(k, 2))
		k = randInt()
		t2(strconv.FormatInt(k, 10), "0o"+strconv.FormatInt(k, 8))
		k = randInt()
		t2(strconv.FormatInt(k, 10), "0O"+strconv.FormatInt(k, 8))
		k = randInt()
		t2(strconv.FormatInt(k, 10), "0x"+strconv.FormatInt(k, 16))
		k = randInt()
		t2(strconv.FormatInt(k, 10), "0X"+strconv.FormatInt(k, 16))
	}

	// Binary.
	t2("0", "0b0")
	t2("0", "0B0")
	t2("-5", "-0b101")
	t2("5", "+0b101")
	t2("1.5", "0b1.1")
	t2("-1.5", "-0b1.1")

	t2("18181", "0b100011100000101.00")
	t2("-12.5", "-0b1100.10")
	t2("343872.5", "0b1010011111101000000.10")
	t2("-328.28125", "-0b101001000.010010")
	t2("-341919.144535064697265625", "-0b1010011011110011111.0010010100000000010")
	t2("97.10482025146484375", "0b1100001.000110101101010110000")
	t2("-120914.40625", "-0b11101100001010010.01101")
	t2("8080777260861123367657", "0b1101101100000111101001111111010001111010111011001010100101001001011101001")

	// Octal.
	t2("8", "0o10")
	t2("-8.5", "-0O010.4")
	t2("8.5", "+0O010.4")
	t2("-262144.000000059604644775390625", "-0o1000000.00000001")
	t2("572315667420.390625", "0o10250053005734.31")

	// Hex.
	t2("1", "0x00001")
	t2("255", "0xff")
	t2("-15.5", "-0Xf.8")
	t2("15.5", "+0Xf.8")
	t2("-16777216.00000000023283064365386962890625", "-0x1000000.00000001")
	t2("325927753012307620476767402981591827744994693483231017778102969592507", "0xc16de7aa5bf90c3755ef4dea45e982b351b6e00cd25a82dcfe0646abb")

	// Test parsing.

	tx := func(fn func(), msg string) {
		t.Helper()
		assertException(t, fn, msg)
	}

	t2("NaN", nan())
	t2("NaN", -nan())
	t2("NaN", "NaN")
	t2("NaN", "-NaN")
	t2("NaN", "+NaN")

	tx(func() { New(" NaN") }, "' NaN'")
	tx(func() { New("NaN ") }, "'NaN '")
	tx(func() { New(" NaN ") }, "' NaN '")
	tx(func() { New(" -NaN") }, "' -NaN'")
	tx(func() { New(" +NaN") }, "' +NaN'")
	tx(func() { New("-NaN ") }, "'-NaN '")
	tx(func() { New("+NaN ") }, "'+NaN '")
	tx(func() { New(".NaN") }, "'.NaN'")
	tx(func() { New("NaN.") }, "'NaN.'")

	t2("Infinity", inf(1))
	t2("-Infinity", inf(-1))
	t2("Infinity", "Infinity")
	t2("-Infinity", "-Infinity")
	t2("Infinity", "+Infinity")

	tx(func() { New(" Infinity") }, "' Infinity '")
	tx(func() { New("Infinity ") }, "'Infinity '")
	tx(func() { New(" Infinity ") }, "' Infinity '")
	tx(func() { New(" -Infinity") }, "' -Infinity'")
	tx(func() { New(" +Infinity") }, "' +Infinity'")
	tx(func() { New(".Infinity") }, "'.Infinity'")
	tx(func() { New("Infinity.") }, "'Infinity.'")

	t2("0", 0.0)
	t2("-0", math_NegZero())
	t2("0", "0")
	t2("-0", "-0")
	t2("0", "0.")
	t2("-0", "-0.")
	t2("0", "0.0")
	t2("-0", "-0.0")
	t2("0", "0.00000000")
	t2("-0", "-0.0000000000000000000000")

	tx(func() { New(" 0") }, "' 0'")
	tx(func() { New("0 ") }, "'0 '")
	tx(func() { New(" 0 ") }, "' 0 '")
	tx(func() { New("0-") }, "'0-'")
	tx(func() { New(" -0") }, "' -0'")
	tx(func() { New("-0 ") }, "'-0 '")
	tx(func() { New("+0 ") }, "'+0 '")
	tx(func() { New(" +0") }, "' +0'")
	tx(func() { New(" .0") }, "' .0'")
	tx(func() { New("0. ") }, "'0. '")
	tx(func() { New("+-0") }, "'+-0'")
	tx(func() { New("-+0") }, "'-+0'")
	tx(func() { New("--0") }, "'--0'")
	tx(func() { New("++0") }, "'++0'")
	tx(func() { New(".-0") }, "'.-0'")
	tx(func() { New(".+0") }, "'.+0'")
	tx(func() { New("0 .") }, "'0 .'")
	tx(func() { New(". 0") }, "'. 0'")
	tx(func() { New("..0") }, "'..0'")
	tx(func() { New("+.-0") }, "'+.-0'")
	tx(func() { New("-.+0") }, "'-.+0'")
	tx(func() { New("+. 0") }, "'+. 0'")
	tx(func() { New(".0.") }, "'.0.'")

	t2("1", 1.0)
	t2("-1", -1.0)
	t2("1", "1")
	t2("-1", "-1")
	t2("0.1", ".1")
	t2("0.1", ".1")
	t2("-0.1", "-.1")
	t2("0.1", "+.1")
	t2("1", "1.")
	t2("1", "1.0")
	t2("-1", "-1.")
	t2("1", "+1.")
	t2("-1", "-1.0000")
	t2("1", "1.0000")
	t2("1", "1.00000000")
	t2("-1", "-1.000000000000000000000000")
	t2("1", "+1.000000000000000000000000")

	tx(func() { New(" 1") }, "' 1'")
	tx(func() { New("1 ") }, "'1 '")
	tx(func() { New(" 1 ") }, "' 1 '")
	tx(func() { New("1-") }, "'1-'")
	tx(func() { New(" -1") }, "' -1'")
	tx(func() { New("-1 ") }, "'-1 '")
	tx(func() { New(" +1") }, "' +1'")
	tx(func() { New("+1 ") }, "'+1'")
	tx(func() { New(".1.") }, "'.1.'")
	tx(func() { New("+-1") }, "'+-1'")
	tx(func() { New("-+1") }, "'-+1'")
	tx(func() { New("--1") }, "'--1'")
	tx(func() { New("++1") }, "'++1'")
	tx(func() { New(".-1") }, "'.-1'")
	tx(func() { New(".+1") }, "'.+1'")
	tx(func() { New("1 .") }, "'1 .'")
	tx(func() { New(". 1") }, "'. 1'")
	tx(func() { New("..1") }, "'..1'")
	tx(func() { New("+.-1") }, "'+.-1'")
	tx(func() { New("-.+1") }, "'-.+1'")
	tx(func() { New("+. 1") }, "'+. 1'")
	tx(func() { New("-. 1") }, "'-. 1'")
	tx(func() { New("1..") }, "'1..'")
	tx(func() { New("+1..") }, "'+1..'")
	tx(func() { New("-1..") }, "'-1..'")
	tx(func() { New("-.1.") }, "'-.1.'")
	tx(func() { New("+.1.") }, "'+.1.'")
	tx(func() { New(".-10.") }, "'.-10.'")
	tx(func() { New(".+10.") }, "'.+10.'")
	tx(func() { New(". 10.") }, "'. 10.'")

	t2("123.456789", 123.456789)
	t2("-123.456789", -123.456789)
	t2("-123.456789", "-123.456789")
	t2("123.456789", "123.456789")
	t2("123.456789", "+123.456789")

	// JS undefined and null both port to nil.
	tx(func() { New(nil) }, "void 0")
	tx(func() { New("undefined") }, "'undefined'")
	tx(func() { New(nil) }, "null")
	tx(func() { New("null") }, "'null'")
	// JS {} and [].
	tx(func() { New(map[string]any{}) }, "{}")
	tx(func() { New([]any{}) }, "[]")
	// JS function () {}.
	tx(func() { New(func() {}) }, "function () {}")
	// JS new Date and new RegExp.
	tx(func() { New(time.Now()) }, "new Date")
	tx(func() { New(regexp.MustCompile("")) }, "new RegExp")
	tx(func() { New("") }, "''")
	tx(func() { New(" ") }, "' '")
	tx(func() { New("nan") }, "'nan'")
	tx(func() { New("23e") }, "'23e'")
	tx(func() { New("e4") }, "'e4'")
	tx(func() { New("ff") }, "'ff'")
	tx(func() { New("0xg") }, "'oxg'")
	tx(func() { New("0Xfi") }, "'0Xfi'")
	tx(func() { New("++45") }, "'++45'")
	tx(func() { New("--45") }, "'--45'")
	tx(func() { New("9.99--") }, "'9.99--'")
	tx(func() { New("9.99++") }, "'9.99++'")
	tx(func() { New("0 0") }, "'0 0'")
}
