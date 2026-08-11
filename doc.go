// Package units provides compile-time checked physical quantities for the
// enclaboratory fleet (HUS §6.1 SSoT / §11): Power (watts) and Energy
// (kilowatt-hours) as defined float64 types, so kW/MW/kWh values stop being raw
// float64 fields distinguished only by a name suffix (PricePerKwh,
// totalSupplyKWh, raw MW columns).
//
// Each type fixes one internal representation — Power is always W, Energy is
// always kWh (v0.2.0 — kWh is the settlement domain unit). Values are created through unit constructors (W, KW, MW / Wh,
// KWh, MWh) and read back through unit accessors (.W(), .KW(), .MW() /
// .Wh(), .KWh(), .MWh()), so a unit mistake is a visible conversion at the
// boundary instead of a silent scaling bug in the middle of a formula.
// Because Power and Energy are distinct defined types, mixing them
// (power + energy, assigning one to the other) is a compile error; crossing
// the dimension requires the explicit operators Power.Over (P × t → E) and
// Energy.Per (E / t → P).
//
// Validate rejects NaN/±Inf and is meant for external input boundaries
// (API payloads, CSV imports, DB reads). String always prints the base unit
// (W / Wh) — no automatic unit selection.
//
// # Non-goals
//
// These are deliberate exclusions, not missing features:
//
//  1. Money. Currency amounts are integer minor units, never float64 —
//     the canonical fleet pattern is billing's domain.Money
//     (services/billing/internal/domain/money.go in the billing repo).
//     This package must never grow a currency type.
//  2. Tariffs (₩/kWh and friends). A unit price is a composition of money
//     and a physical quantity, so it lives outside this package: consumers
//     multiply an integer minor amount by the quantity and round exactly
//     once, following the billing money policy.
//  3. Display formatting. String is a fixed base-unit representation for
//     logs and debugging. Choosing kW vs MW for humans, locale formatting,
//     and precision are presentation policy owned by each consumer.
package units
