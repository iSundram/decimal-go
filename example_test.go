package decimal_test

import (
	"encoding/json"
	"fmt"

	dec "github.com/iSundram/decimal-go"
)

// mk returns a cloned constructor with the decimal.js default settings. Each
// example builds its own so that examples are deterministic regardless of
// which tests ran before them on whatever shared Default.
func mk() *dec.Constructor {
	return dec.Default.Clone(&dec.Config{
		Precision: dec.I64(20),
		Rounding:  dec.I64(dec.RoundHalfUp),
		ToExpNeg:  dec.I64(-7),
		ToExpPos:  dec.I64(21),
		MaxE:      dec.I64(9e15),
		MinE:      dec.I64(-9e15),
	})
}

func Example() {
	// Money calculation without float rounding surprises.
	c := dec.Default.Clone(&dec.Config{
		Precision: dec.I64(12),
		Rounding:  dec.I64(dec.RoundHalfUp),
		ToExpNeg:  dec.I64(-100),
		ToExpPos:  dec.I64(100),
	})
	price := c.New("19.99")
	total := price.Times(c.New("3"))
	fmt.Println(total)
	// Output:
	// 59.97
}

func ExampleNew() {
	fmt.Println(mk().New("1.5").Plus(mk().New("2.25")))
	fmt.Println(mk().New(42).Div(mk().New(8)))
	fmt.Println(mk().New("0x1p4"))
	// Output:
	// 3.75
	// 5.25
	// 16
}

func ExampleDecimal_ToFixed() {
	c := mk()
	n := c.New("3.14159265358979323846")
	fmt.Println(n.ToFixed(4))
	fmt.Println(n.ToExponential(3))
	fmt.Println(n.ToPrecision(6))
	// Output:
	// 3.1416
	// 3.142e+0
	// 3.14159
}

func ExampleDecimal_Sqrt() {
	c := mk()
	fmt.Println(c.New("2").Sqrt())
	fmt.Println(c.New("16").Sqrt())
	fmt.Println(c.New("-1").Sqrt())
	// Output:
	// 1.4142135623730950488
	// 4
	// NaN
}

func ExampleDecimal_Pow() {
	c := mk()
	fmt.Println(c.New("2").Pow(c.New("10")))
	fmt.Println(c.New("0").Pow(c.New("0")))
	fmt.Println(c.New("-1").Pow(c.New("0.5")))
	// Output:
	// 1024
	// 1
	// NaN
}

func ExampleDecimal_ValueOf() {
	// String omits the sign of -0; ValueOf keeps it (decimal.js parity).
	c := mk()
	x := c.New("-0")
	fmt.Println(x.String())
	fmt.Println(x.ValueOf())
	// Output:
	// 0
	// -0
}

func ExampleDecimal_MarshalJSON() {
	type invoice struct {
		Total dec.Decimal `json:"total"`
	}
	inv := invoice{Total: *mk().New("1234.5000")}
	b, _ := json.Marshal(inv)
	fmt.Println(string(b))
	// Output:
	// {"total":"1234.5"}
}

func ExampleConstructor_Clone() {
	// Clones are how you get a constructor whose settings differ from Default.
	base := mk()
	round := base.Clone(&dec.Config{Precision: dec.I64(8)})
	x := round.New("1.00000001")
	fmt.Println(x.Times(x)) // rounds to 8 significant digits
	// Output:
	// 1
}

func ExampleDecimal_Log() {
	c := mk()
	p := c.New("1000")
	fmt.Println(c.Log10(p))
	fmt.Println(p.Log(c.New("10")))
	// Output:
	// 3
	// 3
}
