package decimal

import (
	"database/sql/driver"
	"errors"
)

// MarshalText implements encoding.TextMarshaler. The textual form is the
// valueOf() representation (the same string ValueOf returns), so text
// encodings round-trip exactly, including an explicit "-0".
func (x Decimal) MarshalText() ([]byte, error) {
	return []byte(x.ValueOf()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. b is parsed as a decimal
// string using the constructor the receiver belongs to (or Default when the
// receiver is a zero-value Decimal). It returns an error if b is not a valid
// representation.
func (x *Decimal) UnmarshalText(b []byte) error {
	if x == nil {
		return errors.New(decimalError + "nil *Decimal")
	}
	c := x.c
	if c == nil {
		c = Default
	}
	*x = *c.New(string(b))
	return nil
}

// Scan implements database/sql.Scanner. src must be a type accepted by
// Constructor.New (string, integer or float64), a []byte (treated as its
// decimal string), a *Decimal (copied), or nil (stored as NaN to represent
// SQL NULL). The receiver's constructor settings are used for parsing.
func (x *Decimal) Scan(src any) error {
	if x == nil {
		return errors.New(decimalError + "Scan on nil *Decimal")
	}
	c := x.c
	if c == nil {
		c = Default
	}
	if src == nil {
		*x = *c.newNaN()
		return nil
	}
	if b, ok := src.([]byte); ok {
		xv := c.New(string(b))
		*x = *xv
		return nil
	}
	*x = *c.New(src)
	return nil
}

// Value implements driver.Valuer for database/sql. The value is returned as a
// string in the valueOf() representation so no precision is lost. A nil
// receiver returns nil (SQL NULL).
func (x *Decimal) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}
	return x.ValueOf(), nil
}
