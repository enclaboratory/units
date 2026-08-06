package units

import (
	"errors"
	"fmt"
	"math"
)

// Sentinel errors returned (wrapped) by the Validate methods.
var (
	// ErrNaN reports a value that is NaN.
	ErrNaN = errors.New("value is NaN")
	// ErrInfinite reports a value that is +Inf or -Inf.
	ErrInfinite = errors.New("value is infinite")
)

// validate is the shared finite-value check behind Power.Validate and
// Energy.Validate.
func validate(v float64, kind string) error {
	switch {
	case math.IsNaN(v):
		return fmt.Errorf("units: %s: %w", kind, ErrNaN)
	case math.IsInf(v, 0):
		return fmt.Errorf("units: %s: %w", kind, ErrInfinite)
	}
	return nil
}
