package uuid

import (
	"strings"
	"testing"
)

// TestParseURNPrefixCaseInsensitive covers the URN form documented on
// UnmarshalText: RFC 5234 section 2.3 makes ABNF literals case-insensitive and
// RFC 2141 section 3 says the "urn" label and the namespace identifier are
// case-insensitive, so all of these have to parse to the same UUID.
func TestParseURNPrefixCaseInsensitive(t *testing.T) {
	const canonical = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	want, err := FromString(canonical)
	if err != nil {
		t.Fatalf("FromString(%q): %v", canonical, err)
	}

	for _, s := range []string{
		"urn:uuid:" + canonical,
		"URN:UUID:" + canonical,
		"uRn:uUiD:" + canonical,
		"URN:uuid:" + strings.ToUpper(canonical),
	} {
		got, err := FromString(s)
		if err != nil {
			t.Errorf("FromString(%q) returned error %v, want nil", s, err)
			continue
		}
		if got != want {
			t.Errorf("FromString(%q) = %v, want %v", s, got, want)
		}

		var fromText UUID
		if err := fromText.UnmarshalText([]byte(s)); err != nil {
			t.Errorf("UnmarshalText(%q) returned error %v, want nil", s, err)
		}
	}
}
