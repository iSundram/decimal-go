package decimal

import (
	"database/sql"
	"encoding/json"
	"testing"
)

func TestEncodingText(t *testing.T) {
	d := Default.Clone(&Config{
		Precision: I64(20),
		Rounding:  I64(4),
		ToExpNeg:  I64(-7),
		ToExpPos:  I64(21),
	})

	cases := []struct {
		in   string
		want string
	}{
		{"123.456", "123.456"},
		{"-0", "-0"},
		{"1e21", "1e+21"},
		{"0.0000001", "1e-7"},
	}
	for _, tc := range cases {
		x := d.New(tc.in)
		b, err := x.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != tc.want {
			t.Errorf("MarshalText(%q) = %q, want %q", tc.in, b, tc.want)
		}
		var y Decimal
		if err := y.UnmarshalText(b); err != nil {
			t.Fatal(err)
		}
		if y.ValueOf() != tc.want {
			t.Errorf("round-trip UnmarshalText(%q) = %q, want %q", b, y.ValueOf(), tc.want)
		}
	}
}

func TestJSON(t *testing.T) {
	var got struct {
		Price Decimal `json:"price"`
	}
	src := `{"price":"19.99"}`
	if err := json.Unmarshal([]byte(src), &got); err != nil {
		t.Fatal(err)
	}
	// UnmarshalText feeds a zero-value Decimal with Default config.
	if got.Price.ValueOf() != "19.99" {
		t.Errorf("unmarshal = %q, want 19.99", got.Price.ValueOf())
	}
	b, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"price":"19.99"}` {
		t.Errorf("marshal = %s, want {\"price\":\"19.99\"}", b)
	}
}

func TestSQL(t *testing.T) {
	d := Default.Clone(&Config{Precision: I64(30), Rounding: I64(4)})

	// Value -> string.
	v, err := d.New("123.456789012345678901234567").Value()
	if err != nil {
		t.Fatal(err)
	}
	if s := v.(string); s != "123.456789012345678901234567" {
		t.Errorf("Value = %q", s)
	}

	// Scan from the various source types.
	for _, src := range []any{"3.14", []byte("-42"), int64(7), 2.5} {
		var y Decimal
		if err := y.Scan(src); err != nil {
			t.Fatalf("Scan(%v): %v", src, err)
		}
		if !y.IsNaN() && y.ValueOf() == "" {
			t.Fatalf("Scan(%v) produced empty value", src)
		}
	}

	// NULL -> NaN.
	var z Decimal
	if err := z.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if !z.IsNaN() {
		t.Errorf("Scan(nil) = %q, want NaN", z.ValueOf())
	}
}

var _ sql.Scanner // ensure the interface is implemented at compile time.
