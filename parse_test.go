package uuid

import "testing"

// TestParseAndUnmarshalTextParity ensures that the public Parse(string)
// helper and the encoding.TextUnmarshaler implementation share the same
// behaviour after refactor (see parseBytes). We verify both success and
// failure scenarios across multiple input variations.
func TestParseAndUnmarshalTextParity(t *testing.T) {
	cases := []struct {
		name  string
		input string
		valid bool
	}{
		{"canonical", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"hash", "6ba7b8109dad11d180b400c04fd430c8", true},
		{"bracedCanonical", "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}", true},
		{"bracedHash", "{6ba7b8109dad11d180b400c04fd430c8}", true},
		{"urnCanonical", "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"urnHash", "urn:uuid:6ba7b8109dad11d180b400c04fd430c8", true},

		{"invalidLength", "6ba7b810-9dad-11d1-80b4", false},
		{"invalidDashPositions", "6ba7b8109dad-11d1-80b4-00c04fd430c8", false},
		{"invalidHexChar", "6ba7b810-9dad-11d1-80b4-00c04fd43zzz", false},
		{"invalidURNPrefix", "urn:uuid6ba7b810-9dad-11d1-80b4-00c04fd430c8", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p1, p2 UUID

			err1 := p1.Parse(tc.input)
			err2 := p2.UnmarshalText([]byte(tc.input))

			if tc.valid {
				if err1 != nil {
					t.Fatalf("Parse returned unexpected error: %v", err1)
				}
				if err2 != nil {
					t.Fatalf("UnmarshalText returned unexpected error: %v", err2)
				}
				if p1 != p2 {
					t.Fatalf("Parse and UnmarshalText results differ: %v vs %v", p1, p2)
				}
			} else {
				if err1 == nil {
					t.Errorf("Parse should fail for invalid input")
				}
				if err2 == nil {
					t.Errorf("UnmarshalText should fail for invalid input")
				}
			}
		})
	}
}
