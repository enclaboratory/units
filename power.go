package units

import (
	"strconv"
	"time"
)

// Power is an electrical power quantity. The internal representation is
// always watts (W).
//
// Being a defined type (not a plain float64), Power cannot be mixed with
// Energy or with raw float64 values without an explicit conversion, which
// turns unit confusion into a compile error.
type Power float64

// Conversion factors between SI prefixes of the base units.
const (
	kilo = 1e3
	mega = 1e6
)

// W constructs a Power from a value expressed in watts.
func W(v float64) Power { return Power(v) }

// KW constructs a Power from a value expressed in kilowatts.
func KW(v float64) Power { return Power(v * kilo) }

// MW constructs a Power from a value expressed in megawatts.
func MW(v float64) Power { return Power(v * mega) }

// W returns the power in watts.
func (p Power) W() float64 { return float64(p) }

// KW returns the power in kilowatts.
func (p Power) KW() float64 { return float64(p) / kilo }

// MW returns the power in megawatts.
func (p Power) MW() float64 { return float64(p) / mega }

// Over integrates the power over a duration: E = P × t.
// KW(500).Over(2 * time.Hour) is MWh(1).
func (p Power) Over(d time.Duration) Energy {
	return Energy(float64(p) * d.Hours())
}

// Validate rejects NaN and ±Inf. Call it at external input boundaries
// (API payloads, CSV imports, DB reads); errors.Is works against ErrNaN
// and ErrInfinite.
func (p Power) Validate() error { return validate(float64(p), "Power") }

// String formats the power in the fixed base unit (W). No automatic unit
// selection — display policy belongs to the consumer.
func (p Power) String() string {
	return strconv.FormatFloat(float64(p), 'g', -1, 64) + " W"
}
