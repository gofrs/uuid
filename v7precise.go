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
	"errors"
	"io"
	"time"
)

// PreciseGenerator is implemented by generators that can also produce
// method 3 UUIDs, the ones carrying sub-millisecond precision in rand_a.
//
// It is kept separate from Generator rather than folded into it so that
// this stays a pure addition: anything already implementing Generator
// still does, and a generator that has no notion of sub-millisecond time
// is not obliged to invent one.
type PreciseGenerator interface {
	NewV7Precise() (UUID, error)
	NewV7AtTimePrecise(atTime time.Time) (UUID, error)
}

// ErrNoPreciseGenerator is returned by the package-level NewV7Precise and
// NewV7AtTimePrecise when DefaultGenerator has been replaced by a
// generator that does not implement PreciseGenerator. Reported rather than
// quietly falling back to a method 1 UUID, which would look right and
// carry no sub-millisecond time at all.
var ErrNoPreciseGenerator = errors.New("uuid: DefaultGenerator does not implement PreciseGenerator")

// interface check -- build will fail if *Gen doesn't satisfy PreciseGenerator
var _ PreciseGenerator = (*Gen)(nil)

// subMsUnits is how many parts of a millisecond rand_a divides into when it
// carries sub-millisecond precision: 12 bits, so 4096 steps of about 244ns.
const subMsUnits = maxV7Counter + 1

// nsPerMs is the number of nanoseconds in the millisecond that unix_ts_ms
// counts in, the unit subMsFraction scales down from.
const nsPerMs = int64(time.Millisecond)

// subMsFraction returns the position of t within its millisecond, expressed
// in 4096ths of that millisecond, as RFC 9562 section 6.2 method 3
// specifies for rand_a.
//
// The result is always in [0, 4095]: a time before the epoch has a negative
// nanosecond remainder, which is folded forwards into the millisecond it
// falls in rather than producing a negative fraction.
func subMsFraction(t time.Time) uint16 {
	ns := t.UnixNano() % nsPerMs
	if ns < 0 {
		ns += nsPerMs
	}
	return uint16(ns * subMsUnits / nsPerMs)
}

// NewV7Precise returns a k-sortable UUID that carries sub-millisecond
// precision, by filling rand_a with the fraction of the millisecond the
// timestamp falls in rather than with a counter - RFC 9562 section 6.2,
// method 3.
//
// NewV7 spends those same 12 bits on a monotonic counter (method 1), which
// orders UUIDs within a millisecond but says nothing about when inside it
// they were created. Method 3 instead resolves the timestamp itself to
// about 244ns, which matters when what reads the UUID back wants the time
// rather than the order - PostgreSQL's uuid_extract_timestamp and
// TimescaleDB's uuid_timestamp_micros both read rand_a this way.
//
// UUIDs from a single generator remain strictly increasing. Where the
// fraction alone does not advance - two calls inside the same 244ns step,
// or a clock that moved backwards - it is incremented instead, carrying
// into the timestamp when it overflows, the same combination of method 3
// with method 1 that PostgreSQL's own uuidv7() uses. The timestamp is then
// slightly ahead of the real one, as it is for NewV7 past 4096 UUIDs in a
// millisecond.
func NewV7Precise() (UUID, error) {
	g, ok := DefaultGenerator.(PreciseGenerator)
	if !ok {
		return Nil, ErrNoPreciseGenerator
	}
	return g.NewV7Precise()
}

// NewV7AtTimePrecise returns a UUID for the provided time with
// sub-millisecond precision in rand_a, as described for NewV7Precise.
//
// Unlike NewV7Precise, the timestamp is encoded exactly as given: the same
// time always produces the same 60 bits of timestamp and fraction, and
// nothing is incremented to keep it ahead of a previous call. That is what
// makes it usable for stamping a UUID with a time that came from somewhere
// else - a measurement's own clock, say - where a UUID that disagrees with
// the timestamp it was built from would defeat the purpose.
//
// The cost of that exactness is that two calls for the same time, or for
// times inside the same 244ns step, are ordered only by their random bits.
// Use NewV7Precise where the generator owns the clock and ordering has to
// hold.
//
// Times outside the range the 48-bit millisecond field can represent are
// pinned to the nearest end of it, as they are for NewV7AtTime.
func NewV7AtTimePrecise(atTime time.Time) (UUID, error) {
	g, ok := DefaultGenerator.(PreciseGenerator)
	if !ok {
		return Nil, ErrNoPreciseGenerator
	}
	return g.NewV7AtTimePrecise(atTime)
}

// NewV7Precise returns a k-sortable UUID with sub-millisecond precision in
// rand_a. See the package-level NewV7Precise for the details.
func (g *Gen) NewV7Precise() (UUID, error) {
	return g.newV7Precise(g.epochFunc(), true)
}

// NewV7AtTimePrecise returns a UUID for atTime with sub-millisecond
// precision in rand_a, encoded exactly as given. See the package-level
// NewV7AtTimePrecise for the details.
func (g *Gen) NewV7AtTimePrecise(atTime time.Time) (UUID, error) {
	return g.newV7Precise(atTime, false)
}

// newV7Precise builds a method 3 V7 UUID for atTime. When monotonic is set
// the result is held strictly above the previous one, which may push the
// encoded time ahead of atTime; otherwise atTime is encoded as given.
func (g *Gen) newV7Precise(atTime time.Time, monotonic bool) (UUID, error) {
	var u UUID

	// Same pinning as newV7: a time the timestamp field cannot represent is
	// replaced by the nearest end of the range rather than the unrelated
	// value it would wrap to.
	ms := min(uint64(max(atTime.UnixMilli(), 0)), maxV7Timestamp)
	frac := subMsFraction(atTime)

	if monotonic {
		ms, frac = g.nextV7PreciseSequence(ms, frac)
	}

	u[0] = byte(ms >> 40)
	u[1] = byte(ms >> 32)
	u[2] = byte(ms >> 24)
	u[3] = byte(ms >> 16)
	u[4] = byte(ms >> 8)
	u[5] = byte(ms)

	// rand_a holds the sub-millisecond fraction. Its top 4 bits are
	// overwritten by SetVersion below, which is why the fraction is 12 bits
	// wide and not 16.
	u[6] = byte(frac >> 8)
	u[7] = byte(frac)

	u.SetVersion(V7)

	if _, err := io.ReadFull(g.rand, u[8:16]); err != nil {
		return Nil, err
	}
	u.SetVariant(VariantRFC9562)

	return u, nil
}

// nextV7PreciseSequence returns a (timestamp, fraction) pair strictly above
// the one it returned last, so that UUIDs from this generator keep
// increasing even when the clock does not.
func (g *Gen) nextV7PreciseSequence(ms uint64, frac uint16) (uint64, uint16) {
	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	// Strictly newer than what was encoded last: take it as it is, which is
	// the ordinary case and keeps the timestamp exact.
	if !g.v7pSeeded || ms > g.v7pLastMs || (ms == g.v7pLastMs && frac > g.v7pLastFrac) {
		g.v7pSeeded = true
		g.v7pLastMs, g.v7pLastFrac = ms, frac
		return ms, frac
	}

	// The clock has not advanced past the last value - it repeated, or went
	// backwards. Step the fraction instead, carrying into the millisecond
	// once the 12 bits are exhausted.
	if g.v7pLastFrac >= maxV7Counter {
		g.v7pLastMs, g.v7pLastFrac = g.v7pLastMs+1, 0
	} else {
		g.v7pLastFrac++
	}

	return g.v7pLastMs, g.v7pLastFrac
}
