package units

import (
	"errors"
	"math"
	"testing"
	"time"
)

// relTol is the relative tolerance for round-trip conversion checks.
// A constructor + accessor pair performs one multiplication and one division
// by a power of ten, so the accumulated error is a couple of ULPs — 1e-12
// relative is orders of magnitude above that while still catching any real
// scaling mistake (a wrong factor is off by 1e3 at minimum).
const relTol = 1e-12

func approxEqual(a, b float64) bool {
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	scale := math.Max(math.Abs(a), math.Abs(b))
	return diff <= relTol*scale
}

// --- Compile-time checking: what the defined types actually catch ---
//
// The point of Power/Energy over raw float64 is that dimension mixing does
// not compile. None of the following lines build (uncomment any one of them
// to see the compiler reject it):
//
//	var p Power = KWh(1)              // Energy is not Power
//	var e Energy = KW(1)              // Power is not Energy
//	_ = KW(1) + KWh(1)                // invalid operation: mismatched types
//	var raw float64 = KW(1)           // no silent decay back to float64
//	func f(pricePerKwh float64) {}; f(KWh(1)) // can't pass Energy as float64
//
// Crossing the dimension must go through the explicit operators:
//
//	var e Energy = KW(1).Over(time.Hour) // OK: P × t → E
//	var p Power = KWh(1).Per(time.Hour)  // OK: E / t → P
//
// This is exactly the class of bug the fleet has today with suffix-named
// float64 fields (PricePerKwh, totalSupplyKWh, raw MW): the compiler cannot
// see the unit, so a kW value flows into a MW field unnoticed. With defined
// types that mistake is a build failure, not a production data bug.

func TestPowerRoundTrip(t *testing.T) {
	values := []float64{0, 1, -1, 0.1, 1.5, 123.456, 1e-9, 3.14159e8}
	for _, v := range values {
		if got := W(v).W(); got != v {
			t.Errorf("W(%v).W() = %v, want exactly %v", v, got, v)
		}
		if got := KW(v).KW(); !approxEqual(got, v) {
			t.Errorf("KW(%v).KW() = %v, want ~%v", v, got, v)
		}
		if got := MW(v).MW(); !approxEqual(got, v) {
			t.Errorf("MW(%v).MW() = %v, want ~%v", v, got, v)
		}
	}
}

func TestEnergyRoundTrip(t *testing.T) {
	values := []float64{0, 1, -1, 0.1, 1.5, 123.456, 1e-9, 3.14159e8}
	for _, v := range values {
		if got := Wh(v).Wh(); got != v {
			t.Errorf("Wh(%v).Wh() = %v, want exactly %v", v, got, v)
		}
		if got := KWh(v).KWh(); !approxEqual(got, v) {
			t.Errorf("KWh(%v).KWh() = %v, want ~%v", v, got, v)
		}
		if got := MWh(v).MWh(); !approxEqual(got, v) {
			t.Errorf("MWh(%v).MWh() = %v, want ~%v", v, got, v)
		}
	}
}

func TestCrossUnitScaling(t *testing.T) {
	// Integral scale factors are exact in float64 — no tolerance here.
	if got := KW(1).W(); got != 1000 {
		t.Errorf("KW(1).W() = %v, want 1000", got)
	}
	if got := MW(1).KW(); got != 1000 {
		t.Errorf("MW(1).KW() = %v, want 1000", got)
	}
	if got := MW(1).W(); got != 1e6 {
		t.Errorf("MW(1).W() = %v, want 1e6", got)
	}
	if got := KWh(1).Wh(); got != 1000 {
		t.Errorf("KWh(1).Wh() = %v, want 1000", got)
	}
	if got := MWh(1).KWh(); got != 1000 {
		t.Errorf("MWh(1).KWh() = %v, want 1000", got)
	}
}

func TestPowerOver(t *testing.T) {
	tests := []struct {
		name string
		p    Power
		d    time.Duration
		want Energy
	}{
		{"1 W for 1 h is 1 Wh", W(1), time.Hour, Wh(1)},
		{"500 kW for 2 h is 1 MWh", KW(500), 2 * time.Hour, MWh(1)},
		{"2 kW for 30 min is 1 kWh", KW(2), 30 * time.Minute, KWh(1)},
		{"anything over 0 s is 0 Wh", MW(3), 0, Wh(0)},
		{"negative power integrates negative", KW(-1), time.Hour, KWh(-1)},
	}
	for _, tt := range tests {
		if got := tt.p.Over(tt.d); !approxEqual(got.Wh(), tt.want.Wh()) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestEnergyPer(t *testing.T) {
	tests := []struct {
		name string
		e    Energy
		d    time.Duration
		want Power
	}{
		{"1 Wh over 1 h is 1 W", Wh(1), time.Hour, W(1)},
		{"1 MWh over 2 h is 500 kW", MWh(1), 2 * time.Hour, KW(500)},
		{"1 kWh over 30 min is 2 kW", KWh(1), 30 * time.Minute, KW(2)},
	}
	for _, tt := range tests {
		if got := tt.e.Per(tt.d); !approxEqual(got.W(), tt.want.W()) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOverPerRoundTrip(t *testing.T) {
	durations := []time.Duration{time.Second, 15 * time.Minute, time.Hour, 24 * time.Hour}
	for _, d := range durations {
		p := KW(736.5)
		if got := p.Over(d).Per(d); !approxEqual(got.W(), p.W()) {
			t.Errorf("Over(%v).Per(%v) = %v, want ~%v", d, d, got, p)
		}
	}
}

func TestPerZeroDurationIsCaughtByValidate(t *testing.T) {
	// E / 0 is not a panic in float64 — it is ±Inf (or NaN for 0/0).
	// The contract is that Validate rejects it at the boundary.
	if err := KWh(1).Per(0).Validate(); !errors.Is(err, ErrInfinite) {
		t.Errorf("KWh(1).Per(0).Validate() = %v, want ErrInfinite", err)
	}
	if err := Wh(0).Per(0).Validate(); !errors.Is(err, ErrNaN) {
		t.Errorf("Wh(0).Per(0).Validate() = %v, want ErrNaN", err)
	}
}

func TestValidate(t *testing.T) {
	if err := KW(42).Validate(); err != nil {
		t.Errorf("KW(42).Validate() = %v, want nil", err)
	}
	if err := MWh(-3.5).Validate(); err != nil {
		t.Errorf("MWh(-3.5).Validate() = %v, want nil", err)
	}
	if err := W(math.NaN()).Validate(); !errors.Is(err, ErrNaN) {
		t.Errorf("W(NaN).Validate() = %v, want ErrNaN", err)
	}
	if err := W(math.Inf(1)).Validate(); !errors.Is(err, ErrInfinite) {
		t.Errorf("W(+Inf).Validate() = %v, want ErrInfinite", err)
	}
	if err := Wh(math.Inf(-1)).Validate(); !errors.Is(err, ErrInfinite) {
		t.Errorf("Wh(-Inf).Validate() = %v, want ErrInfinite", err)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		got  string
		want string
	}{
		{W(1500).String(), "1500 W"},
		{KW(1.5).String(), "1500 W"},
		{MW(2).String(), "2e+06 W"},
		{W(0).String(), "0 W"},
		{Wh(1500).String(), "1500 Wh"},
		{KWh(1.5).String(), "1500 Wh"},
		{Wh(-0.25).String(), "-0.25 Wh"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("String() = %q, want %q", tt.got, tt.want)
		}
	}
}
