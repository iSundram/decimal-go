// xvalid generates operation results from decimal-go and prints one line per
// case using exactly the same format as xvalidate/x.js. Diffing the two files
// verifies behavioral parity with the original decimal.js.
//
// Usage:
//
//	go run ./xvalidate/xvalid.go | sort > go.out
//	node xvalidate/x.js   | sort > js.out
//	diff go.out js.out
//
// Each case line is "op\tinput\tresult" where the input is the raw string
// from cases.txt (additional operands separated by a comma) and result is the
// decimal.js-compatible valueOf() string.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iSundram/decimal-go"
)

func main() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	c := decimal.Default.Clone(&decimal.Config{
		Precision: decimal.I64(20),
		Rounding:  decimal.I64(4),
		ToExpNeg:  decimal.I64(-7),
		ToExpPos:  decimal.I64(21),
		MaxE:      decimal.I64(9e15),
		MinE:      decimal.I64(-9e15),
	})

	raw, err := os.ReadFile("xvalidate/cases.txt")
	if err != nil {
		panic(err)
	}
	vals := strings.Fields(string(raw))

	n := len(vals)
	// Unary.
	for i := 0; i < n; i++ {
		v := vals[i]
		x := c.New(v)
		fmt.Fprintf(w, "abs\t%s\t%s\n", v, x.Abs().ValueOf())
		fmt.Fprintf(w, "neg\t%s\t%s\n", v, x.Neg().ValueOf())
		fmt.Fprintf(w, "trunc\t%s\t%s\n", v, x.Trunc().ValueOf())
		fmt.Fprintf(w, "ceil\t%s\t%s\n", v, x.Ceil().ValueOf())
		fmt.Fprintf(w, "floor\t%s\t%s\n", v, x.Floor().ValueOf())
		fmt.Fprintf(w, "round\t%s\t%s\n", v, x.Round().ValueOf())
		fmt.Fprintf(w, "sqrt\t%s\t%s\n", v, x.Sqrt().ValueOf())
		fmt.Fprintf(w, "cbrt\t%s\t%s\n", v, x.Cbrt().ValueOf())
		fmt.Fprintf(w, "exp\t%s\t%s\n", v, x.Exp().ValueOf())
		fmt.Fprintf(w, "ln\t%s\t%s\n", v, x.Ln().ValueOf())
	}
	// dp/sd return floats; format like JS number (shortest round-trip).
	for i := 0; i < n; i++ {
		v := vals[i]
		x := c.New(v)
		fmt.Fprintf(w, "dp\t%s\t%s\n", v, fmtf(x.Dp()))
		fmt.Fprintf(w, "sd\t%s\t%s\n", v, fmtf(x.Sd()))
		fmt.Fprintf(w, "toFixed2\t%s\t%s\n", v, x.ToFixed(2))
		fmt.Fprintf(w, "toExp6\t%s\t%s\n", v, x.ToExponential(6))
		fmt.Fprintf(w, "toPrec8\t%s\t%s\n", v, x.ToPrecision(8))
	}
	// binary: for each i, pair i with i+1 (wrapped).
	binary := []struct {
		op string
		f  func(a, b *decimal.Decimal) *decimal.Decimal
	}{
		{"add", func(a, b *decimal.Decimal) *decimal.Decimal { return a.Plus(b) }},
		{"sub", func(a, b *decimal.Decimal) *decimal.Decimal { return a.Minus(b) }},
		{"mul", func(a, b *decimal.Decimal) *decimal.Decimal { return a.Times(b) }},
		{"div", func(a, b *decimal.Decimal) *decimal.Decimal { return a.Div(b) }},
		{"divToInt", func(a, b *decimal.Decimal) *decimal.Decimal { return a.DivToInt(b) }},
		{"mod", func(a, b *decimal.Decimal) *decimal.Decimal { return a.Mod(b) }},
		{"pow", func(a, b *decimal.Decimal) *decimal.Decimal { return a.Pow(b) }},
	}
	for i := 0; i < n; i++ {
		a, b := vals[i], vals[(i+1)%n]
		xa, xb := c.New(a), c.New(b)
		for _, op := range binary {
			// keep stable: sort operands matter (a op b may differ from b op a),
			// so preserve order between i and i+1 deterministically.
			fmt.Fprintf(w, "%s\t%s|%s\t%s\n", op.op, a, b, op.f(xa, xb).ValueOf())
		}
		fmt.Fprintf(w, "cmp\t%s|%s\t%s\n", a, b, strconv.FormatFloat(xa.Cmp(xb), 'g', -1, 64))
	}
}

func fmtf(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
