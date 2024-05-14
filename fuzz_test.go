package uuid

import (
	"regexp"
	"testing"
)

var seeds = []string{
	"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
	"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	"6ba7b8109dad11d180b400c04fd430c8",
	"{6ba7b8109dad11d180b400c04fd430c8}",
	"urn:uuid:6ba7b8109dad11d180b400c04fd430c8",
}

const uuidPattern = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

// FuzzFromStringFunc is a fuzz testing suite that exercises the FromString functionn
func FuzzFromStringFunc(f *testing.F) {
	for _, seed := range seeds {
		f.Add(seed)
	}
	uuidRegexp, err := regexp.Compile(uuidPattern)
	if err != nil {
		f.Error("uuid regexp failed to compile")
	}
	f.Fuzz(func(t *testing.T, payload string) {
		u, err := FromString(payload)
		if err != nil {
			if !uuidRegexp.MatchString(u.String()) {
				t.Errorf("%s resulted in invalid uuid %s", payload, u.String())
			}
		}
		// otherwise, allow to pass if no panic
	})
}

// FuzzFromStringOrNil is a fuzz testing suite that exercises the FromStringOrNil functionn
func FuzzFromStringOrNilFunc(f *testing.F) {
	for _, seed := range seeds {
		f.Add(seed)
	}
	uuidRegexp, err := regexp.Compile(uuidPattern)
	if err != nil {
		f.Error("uuid regexp failed to compile")
	}
	f.Fuzz(func(t *testing.T, payload string) {
		u := FromStringOrNil(payload)
		if u != Nil {
			if !uuidRegexp.MatchString(u.String()) {
				t.Errorf("%s resulted in invalid uuid %s", payload, u.String())
			}
		}
		// otherwise, allow to pass if no panic
	})
}
