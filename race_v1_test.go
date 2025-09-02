package uuid

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestV1UniqueConcurrent verifies that Version-1 UUID generation remains
// collision-free under various levels of concurrent load. The test uses
// table-driven subtests to progressively increase the number of goroutines
// and UUIDs generated. We intentionally let the timestamp advance (default
// NewGen) to keep the test quick while still exercising the new
// clock-sequence logic under contention.
func TestV1UniqueConcurrent(t *testing.T) {
	cases := []struct {
		name        string
		goroutines  int
		uuidsPerGor int
	}{
		{"small", 20, 600},    // 12 000 UUIDs (baseline)
		{"medium", 100, 1000}, // 100 000 UUIDs (original failure case)
		{"large", 200, 1000},  // 200 000 UUIDs (high contention)
	}

	for _, tc := range cases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			gen := NewGen()

			var (
				wg       sync.WaitGroup
				mu       sync.Mutex
				seen     = make(map[UUID]struct{}, tc.goroutines*tc.uuidsPerGor)
				dupCount uint32
				genErr   uint32
			)

			for i := 0; i < tc.goroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < tc.uuidsPerGor; j++ {
						u, err := gen.NewV1()
						if err != nil {
							atomic.AddUint32(&genErr, 1)
							return
						}
						mu.Lock()
						if _, exists := seen[u]; exists {
							dupCount++
						} else {
							seen[u] = struct{}{}
						}
						mu.Unlock()
					}
				}()
			}

			wg.Wait()

			if genErr > 0 {
				t.Fatalf("%d errors occurred during UUID generation", genErr)
			}
			if dupCount > 0 {
				t.Fatalf("duplicate UUIDs detected: %d", dupCount)
			}
		})
	}
}
