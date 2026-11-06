// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package uuid

import (
	"bytes"
	"testing"
	"time"
)

// randA returns the 12 bit rand_a field, which method 3 fills with the
// sub-millisecond fraction.
func randA(u UUID) uint16 {
	return uint16(u[6]&0x0f)<<8 | uint16(u[7])
}

// embeddedMs returns the 48 bit unix_ts_ms field.
func embeddedMs(u UUID) int64 {
	return int64(u[0])<<40 | int64(u[1])<<32 | int64(u[2])<<24 |
		int64(u[3])<<16 | int64(u[4])<<8 | int64(u[5])
}

func TestNewV7AtTimePreciseEncodesTheFraction(t *testing.T) {
	// The reference values are what PostgreSQL and TimescaleDB read back
	// out of rand_a: the fraction of the millisecond in 4096ths.
	for _, tt := range []struct {
		name string
		at   time.Time
		want uint16
	}{
		{"start of a millisecond", time.Date(2026, 9, 24, 16, 30, 20, 100_000_000, time.UTC), 0},
		{"125us in", time.Date(2026, 9, 24, 16, 30, 20, 455_125_000, time.UTC), 512},
		{"126us in", time.Date(2026, 9, 24, 16, 30, 20, 455_126_000, time.UTC), 516},
		{"half a millisecond in", time.Date(2026, 9, 24, 16, 30, 20, 100_500_000, time.UTC), 2048},
		{"999us in", time.Date(2026, 9, 24, 16, 30, 20, 999_999_000, time.UTC), 4091},
		{"the epoch itself", time.Unix(0, 0).UTC(), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewV7AtTimePrecise(tt.at)
			if err != nil {
				t.Fatalf("NewV7AtTimePrecise: %v", err)
			}
			if got := randA(u); got != tt.want {
				t.Errorf("rand_a = %d, want %d", got, tt.want)
			}
			if got, want := embeddedMs(u), tt.at.UnixMilli(); got != want {
				t.Errorf("unix_ts_ms = %d, want %d", got, want)
			}
			if u.Version() != V7 {
				t.Errorf("version = %v, want 7", u.Version())
			}
			if u.Variant() != VariantRFC9562 {
				t.Errorf("variant = %v, want RFC9562", u.Variant())
			}
		})
	}
}

func TestNewV7AtTimePreciseIsExact(t *testing.T) {
	// The timestamp must not be nudged to keep it ahead of a previous call:
	// the point of the AtTime form is that the UUID agrees with a timestamp
	// that came from somewhere else.
	at := time.Date(2026, 9, 24, 16, 30, 20, 455_125_000, time.UTC)

	first, err := NewV7AtTimePrecise(at)
	if err != nil {
		t.Fatalf("NewV7AtTimePrecise: %v", err)
	}
	second, err := NewV7AtTimePrecise(at)
	if err != nil {
		t.Fatalf("NewV7AtTimePrecise: %v", err)
	}

	if !bytes.Equal(first[:8], second[:8]) {
		t.Errorf("timestamp bits differ between calls: %x then %x", first[:8], second[:8])
	}
	if first == second {
		t.Error("the two UUIDs are identical, rand_b should still differ")
	}

	// A time older than one already asked for is encoded as given.
	older := at.Add(-time.Hour)
	u, err := NewV7AtTimePrecise(older)
	if err != nil {
		t.Fatalf("NewV7AtTimePrecise: %v", err)
	}
	if got, want := embeddedMs(u), older.UnixMilli(); got != want {
		t.Errorf("unix_ts_ms = %d, want %d - an older time was not encoded as given", got, want)
	}
}

func TestNewV7PreciseIsMonotonic(t *testing.T) {
	// A frozen clock forces every call into the same 244ns step, which is
	// where the method 1 fallback has to take over.
	at := time.Date(2026, 9, 24, 16, 30, 20, 455_125_000, time.UTC)
	g := NewGenWithOptions(WithEpochFunc(func() time.Time { return at }))

	const count = 100
	previous := Nil
	for i := range count {
		u, err := g.NewV7Precise()
		if err != nil {
			t.Fatalf("NewV7Precise: %v", err)
		}
		if i > 0 && bytes.Compare(u[:], previous[:]) <= 0 {
			t.Fatalf("UUID %d (%s) does not sort above the one before it (%s)", i, u, previous)
		}
		previous = u
	}

	// The first one still carries the true fraction; only the repeats step.
	first, err := NewGenWithOptions(WithEpochFunc(func() time.Time { return at })).NewV7Precise()
	if err != nil {
		t.Fatalf("NewV7Precise: %v", err)
	}
	if got, want := randA(first), subMsFraction(at); got != want {
		t.Errorf("first rand_a = %d, want the real fraction %d", got, want)
	}
}

func TestNewV7PreciseCarriesIntoTheTimestamp(t *testing.T) {
	// Exhausting the 12 bits inside one millisecond has to advance the
	// timestamp rather than wrap the fraction back under itself.
	at := time.Date(2026, 9, 24, 16, 30, 20, 999_999_000, time.UTC) // fraction 4091
	g := NewGenWithOptions(WithEpochFunc(func() time.Time { return at }))

	previous := Nil
	for i := range 10 {
		u, err := g.NewV7Precise()
		if err != nil {
			t.Fatalf("NewV7Precise: %v", err)
		}
		if i > 0 && bytes.Compare(u[:], previous[:]) <= 0 {
			t.Fatalf("UUID %d does not sort above the one before it", i)
		}
		previous = u
	}

	if got, want := embeddedMs(previous), at.UnixMilli(); got <= want {
		t.Errorf("unix_ts_ms = %d, want it pushed past %d once the fraction was exhausted", got, want)
	}
}

func TestNewV7PreciseTracksTheClock(t *testing.T) {
	// Where the clock does advance, the fraction has to follow it rather
	// than count.
	base := time.Date(2026, 9, 24, 16, 30, 20, 0, time.UTC)
	offsets := []time.Duration{0, 250 * time.Microsecond, 500 * time.Microsecond, 750 * time.Microsecond}

	i := 0
	g := NewGenWithOptions(WithEpochFunc(func() time.Time {
		at := base.Add(offsets[i])
		i++
		return at
	}))

	for range offsets {
		want := subMsFraction(base.Add(offsets[i]))
		u, err := g.NewV7Precise()
		if err != nil {
			t.Fatalf("NewV7Precise: %v", err)
		}
		if got := randA(u); got != want {
			t.Errorf("offset %s: rand_a = %d, want %d", offsets[i-1], got, want)
		}
	}
}

func TestSubMsFractionStaysInRange(t *testing.T) {
	// Every nanosecond of a millisecond, plus times before the epoch, have
	// to land inside the 12 bits.
	for ns := 0; ns < 1_000_000; ns += 37 {
		at := time.Unix(0, int64(ns))
		if got := subMsFraction(at); got > maxV7Counter {
			t.Fatalf("subMsFraction(%dns) = %d, above the 12 bit maximum", ns, got)
		}
	}

	for _, at := range []time.Time{
		time.Unix(-1, 0),
		time.Unix(-1, -999_999),
		time.Date(1969, 7, 20, 20, 17, 40, 123_456_789, time.UTC),
	} {
		if got := subMsFraction(at); got > maxV7Counter {
			t.Errorf("subMsFraction(%s) = %d, above the 12 bit maximum", at, got)
		}
	}
}

func TestNewV7PreciseRoundTripsEveryMicrosecond(t *testing.T) {
	// Every microsecond inside a millisecond has to survive the trip
	// through rand_a: encode with floor, decode by rounding the way
	// PostgreSQL and TimescaleDB do, and land back on the same microsecond.
	base := time.Date(2026, 9, 24, 16, 30, 20, 0, time.UTC)

	for us := range 1000 {
		at := base.Add(time.Duration(us) * time.Microsecond)
		u, err := NewV7AtTimePrecise(at)
		if err != nil {
			t.Fatalf("NewV7AtTimePrecise: %v", err)
		}

		// (fraction * 1_000_000 / 4096) nanoseconds, rounded to the nearest
		// microsecond.
		ns := int64(randA(u)) * nsPerMs / subMsUnits
		gotUs := (ns + 500) / 1000

		if gotUs != int64(us) {
			t.Errorf("%dus encoded as rand_a=%d, read back as %dus", us, randA(u), gotUs)
		}
	}
}
