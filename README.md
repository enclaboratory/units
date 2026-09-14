# units

Compile-time checked physical quantities for enclaboratory services —
`Power` (internally watts) and `Energy` (internally kilowatt-hours) as defined
`float64` types, replacing raw floats whose unit lives only in a field-name
suffix (`PricePerKwh float64`, `totalSupplyKWh float64`, raw MW columns).

Zero dependencies — standard library only.

```
go get github.com/enclaboratory/units
```

## The model

```go
type Power float64  // internal representation: W
type Energy float64 // internal representation: kWh (v0.2.0)
```

Values enter through unit constructors and leave through unit accessors, so
every conversion is explicit and visible at the boundary:

```go
import "github.com/enclaboratory/units"

cap := units.MW(1.5)       // 1.5 MW plant capacity
load := units.KW(736.5)    // 736.5 kW load

cap.KW()                   // 1500 — explicit unit at the read site
load.MW()                  // 0.7365
```

Mixing dimensions does not compile:

```go
var p units.Power = units.KWh(1) // compile error: Energy is not Power
_ = units.KW(1) + units.KWh(1)   // compile error: mismatched types
var raw float64 = units.KW(1)    // compile error: no silent decay to float64
```

Crossing the dimension requires the explicit operators:

```go
e := units.KW(500).Over(2 * time.Hour) // P × t → Energy: 1 MWh
p := units.MWh(1).Per(2 * time.Hour)   // E / t → Power: 500 kW
```

At external input boundaries (API payloads, CSV imports, DB reads), reject
non-finite values:

```go
if err := units.KWh(payload.Supply).Validate(); err != nil {
    // errors.Is(err, units.ErrNaN) / units.ErrInfinite
}
```

`String()` prints the fixed base unit (`"1500 W"`, `"1.5 kWh"`) for logs and
debugging only.

## Non-goals

Deliberate exclusions — do not add these here:

1. **Money.** Currency is integer minor units, never `float64`. The canonical
   fleet pattern is billing's `domain.Money`
   (`services/billing/internal/domain/money.go` in the billing repo). This
   package will never grow a currency type.
2. **Tariffs (₩/kWh and friends).** A unit price composes money with a
   physical quantity, so it lives outside this package: consumers multiply an
   integer minor amount by the quantity and round exactly once, per the
   billing money policy.
3. **Automatic unit display.** `String()` is fixed base-unit output. Choosing
   kW vs MW for humans, locale formatting, and precision are presentation
   policy owned by each consumer.

## License

Internal — enclaboratory.
