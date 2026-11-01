package uuid

import (
	"bytes"
	"testing"
	"time"
)

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// TestNewV7AtTimeCeilingExhaustion covers the documented behaviour of
// NewV7AtTime: a timestamp outside the range of the 48-bit field is "pinned to
// the nearest end of it". With the timestamp pinned at maxV7Timestamp, running
// out of counter room must not push the timestamp one past the end: the field
// would silently wrap to 0 and hand out UUIDs dated 1970, out of order.
func TestNewV7AtTimeCeilingExhaustion(t *testing.T) {
	gen := NewGenWithOptions(WithRandomReader(zeroReader{}))
	at := time.UnixMilli(1 << 48) // beyond the 48-bit millisecond field

	var prev UUID
	emitted := 0
	var firstErr error
	for i := 0; i < 6000; i++ {
		u, err := gen.NewV7AtTime(at)
		if err != nil {
			firstErr = err
			break
		}
		if i > 0 && bytes.Compare(u[:], prev[:]) <= 0 {
			t.Fatalf("UUIDs not strictly increasing: call %d returned %v, previous was %v", i+1, u, prev)
		}
		prev = u
		emitted++
	}

	if emitted < int(maxV7Counter) {
		t.Fatalf("emitted only %d UUIDs before stopping, want at least %d", emitted, int(maxV7Counter))
	}
	if firstErr == nil {
		t.Fatal("exhausting the timestamp ceiling returned no error")
	}
}
