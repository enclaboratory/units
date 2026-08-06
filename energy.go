package units

import (
	"strconv"
	"time"
)

// Energy is an electrical energy quantity. The internal representation is
// always watt-hours (Wh).
//
// Being a defined type (not a plain float64), Energy cannot be mixed with
// Power or with raw float64 values without an explicit conversion, which
// turns unit confusion into a compile error.
type Energy float64

// Wh constructs an Energy from a value expressed in watt-hours.
func Wh(v float64) Energy { return Energy(v) }

// KWh constructs an Energy from a value expressed in kilowatt-hours.
func KWh(v float64) Energy { return Energy(v * kilo) }

// MWh constructs an Energy from a value expressed in megawatt-hours.
func MWh(v float64) Energy { return Energy(v * mega) }

// Wh returns the energy in watt-hours.
func (e Energy) Wh() float64 { return float64(e) }

// KWh returns the energy in kilowatt-hours.
func (e Energy) KWh() float64 { return float64(e) / kilo }

// MWh returns the energy in megawatt-hours.
func (e Energy) MWh() float64 { return float64(e) / mega }

// Per divides the energy by a duration, yielding the average power:
// P = E / t. MWh(1).Per(2 * time.Hour) is KW(500).
//
// A zero duration produces ±Inf (or NaN for zero energy) — reject it at the
// boundary with Validate.
func (e Energy) Per(d time.Duration) Power {
	return Power(float64(e) / d.Hours())
}

// Validate rejects NaN and ±Inf. Call it at external input boundaries
// (API payloads, CSV imports, DB reads); errors.Is works against ErrNaN
// and ErrInfinite.
func (e Energy) Validate() error { return validate(float64(e), "Energy") }

// String formats the energy in the fixed base unit (Wh). No automatic unit
// selection — display policy belongs to the consumer.
func (e Energy) String() string {
	return strconv.FormatFloat(float64(e), 'g', -1, 64) + " Wh"
}
