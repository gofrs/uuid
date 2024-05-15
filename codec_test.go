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
	"regexp"
	"strings"
	"testing"
)

// codecTestData holds []byte data for a UUID we commonly use for testing.
var codecTestData = []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

// codecTestUUID is the UUID value corresponding to codecTestData.
var codecTestUUID = UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

func TestFromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := FromBytes(codecTestData)
		if err != nil {
			t.Fatal(err)
		}
		if got != codecTestUUID {
			t.Fatalf("FromBytes(%x) = %v, want %v", codecTestData, got, codecTestUUID)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		var short [][]byte
		for i := 0; i < len(codecTestData); i++ {
			short = append(short, codecTestData[:i])
		}
		var long [][]byte
		for i := 1; i < 17; i++ {
			tmp := append(codecTestData, make([]byte, i)...)
			long = append(long, tmp)
		}
		invalid := append(short, long...)
		for _, b := range invalid {
			got, err := FromBytes(b)
			if err == nil {
				t.Fatalf("FromBytes(%x): want err != nil, got %v", b, got)
			}
		}
	})
}

func TestFromBytesOrNil(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		b := []byte{4, 8, 15, 16, 23, 42}
		got := FromBytesOrNil(b)
		if got != Nil {
			t.Errorf("FromBytesOrNil(%x): got %v, want %v", b, got, Nil)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		got := FromBytesOrNil(codecTestData)
		if got != codecTestUUID {
			t.Errorf("FromBytesOrNil(%x): got %v, want %v", codecTestData, got, codecTestUUID)
		}
	})

}

type fromStringTest struct {
	input   string
	variant string
}

// Run runs the FromString test in a subtest of t, named by fst.variant.
func (fst fromStringTest) TestFromString(t *testing.T) {
	t.Run(fst.variant, func(t *testing.T) {
		got, err := FromString(fst.input)
		if err != nil {
			t.Fatalf("FromString(%q): %v", fst.input, err)
		}
		if want := codecTestUUID; got != want {
			t.Fatalf("FromString(%q) = %v, want %v", fst.input, got, want)
		}
	})
}

func (fst fromStringTest) TestUnmarshalText(t *testing.T) {
	t.Run(fst.variant, func(t *testing.T) {
		var u UUID
		err := u.UnmarshalText([]byte(fst.input))
		if err != nil {
			t.Fatalf("UnmarshalText(%q) (%s): %v", fst.input, fst.variant, err)
		}
		if want := codecTestData; !bytes.Equal(u[:], want[:]) {
			t.Fatalf("UnmarshalText(%q) (%s) = %v, want %v", fst.input, fst.variant, u, want)
		}
	})
}

// fromStringTests contains UUID variants that are expected to be parsed
// successfully by UnmarshalText / FromString.
//
// variants must be unique across elements of this slice. Please see the
// comment in fuzz.go if you change this slice or add new tests to it.
var fromStringTests = []fromStringTest{
	{
		input:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		variant: "Canonical",
	},
	{
		input:   "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		variant: "BracedCanonical",
	},
	{
		input:   "{6ba7b8109dad11d180b400c04fd430c8}",
		variant: "BracedHashlike",
	},
	{
		input:   "6ba7b8109dad11d180b400c04fd430c8",
		variant: "Hashlike",
	},
	{
		input:   "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		variant: "URNCanonical",
	},
	{
		input:   "urn:uuid:6ba7b8109dad11d180b400c04fd430c8",
		variant: "URNHashlike",
	},
}

var invalidFromStringInputs = []string{
	// short
	"6ba7b810-9dad-11d1-80b4-00c04fd430c",
	"6ba7b8109dad11d180b400c04fd430c",

	// invalid hex
	"6ba7b8109dad11d180b400c04fd430q8",

	// long
	"6ba7b810-9dad-11d1-80b4-00c04fd430c8=",
	"6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}f",
	"6ba7b810-9dad-11d1-80b4-00c04fd430c800c04fd430c8",

	// malformed in other ways
	"ba7b8109dad11d180b400c04fd430c8}",
	"6ba7b8109dad11d180b400c04fd430c86ba7b8109dad11d180b400c04fd430c8",
	"urn:uuid:{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
	"uuid:urn:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	"uuid:urn:6ba7b8109dad11d180b400c04fd430c8",
	"6ba7b8109-dad-11d1-80b4-00c04fd430c8",
	"6ba7b810-9dad1-1d1-80b4-00c04fd430c8",
	"6ba7b810-9dad-11d18-0b4-00c04fd430c8",
	"6ba7b810-9dad-11d1-80b40-0c04fd430c8",
	"6ba7b810+9dad+11d1+80b4+00c04fd430c8",
	"(6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8>",
	"zba7b810-9dad-11d1-80b4-00c04fd430c8",
	"6ba7b810-9dad11d180b400c04fd430c8",
	"6ba7b8109dad-11d180b400c04fd430c8",
	"6ba7b8109dad11d1-80b400c04fd430c8",
	"6ba7b8109dad11d180b4-00c04fd430c8",
}

func TestFromString(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		for _, fst := range fromStringTests {
			fst.TestFromString(t)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		for _, s := range invalidFromStringInputs {
			got, err := FromString(s)
			if err == nil {
				t.Errorf("FromString(%q): want err != nil, got %v", s, got)
			}
		}
	})
}

func TestFromStringOrNil(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		s := "bad"
		got := FromStringOrNil(s)
		if got != Nil {
			t.Errorf("FromStringOrNil(%q): got %v, want Nil", s, got)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		s := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		got := FromStringOrNil(s)
		if got != codecTestUUID {
			t.Errorf("FromStringOrNil(%q): got %v, want %v", s, got, codecTestUUID)
		}
	})
}

func TestUnmarshalText(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		for _, fst := range fromStringTests {
			fst.TestUnmarshalText(t)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		for _, s := range invalidFromStringInputs {
			var u UUID
			err := u.UnmarshalText([]byte(s))
			if err == nil {
				t.Errorf("FromBytes(%q): want err != nil, got %v", s, u)
			}
		}
	})
}

// Test that UnmarshalText() and Parse() return identical errors
func TestUnmarshalTextParseErrors(t *testing.T) {
	for _, s := range invalidFromStringInputs {
		var u UUID
		e1 := u.UnmarshalText([]byte(s))
		e2 := u.Parse(s)
		if e1 == nil || e1.Error() != e2.Error() {
			t.Errorf("%q: errors don't match: UnmarshalText: %v Parse: %v", s, e1, e2)
		}
	}
}

func TestMarshalBinary(t *testing.T) {
	got, err := codecTestUUID.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, codecTestData) {
		t.Fatalf("%v.MarshalBinary() = %x, want %x", codecTestUUID, got, codecTestData)
	}
}

func TestMarshalText(t *testing.T) {
	want := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	got, err := codecTestUUID.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%v.MarshalText(): got %s, want %s", codecTestUUID, got, want)
	}
}

func TestDecodePlainWithWrongLength(t *testing.T) {
	arg := []byte{'4', '2'}

	u := UUID{}

	if u.UnmarshalText(arg) == nil {
		t.Errorf("%v.UnmarshalText(%q): should return error, but it did not", u, arg)
	}
}

func TestFromHexChar(t *testing.T) {
	const hextable = "0123456789abcdef"

	t.Run("Valid", func(t *testing.T) {
		t.Run("Lower", func(t *testing.T) {
			for i, c := range []byte(hextable) {
				x := fromHexChar(c)
				if int(x) != i {
					t.Errorf("fromHexChar(%c): got %d want %d", c, x, i)
				}
			}
		})
		t.Run("Upper", func(t *testing.T) {
			for i, c := range []byte(strings.ToUpper(hextable)) {
				x := fromHexChar(c)
				if int(x) != i {
					t.Errorf("fromHexChar(%c): got %d want %d", c, x, i)
				}
			}
		})
	})

	t.Run("Invalid", func(t *testing.T) {
		skip := make(map[byte]bool)
		for _, c := range []byte(hextable + strings.ToUpper(hextable)) {
			skip[c] = true
		}
		for i := 0; i < 256; i++ {
			c := byte(i)
			if !skip[c] {
				v := fromHexChar(c)
				if v != 255 {
					t.Errorf("fromHexChar(%c): got %d want: %d", c, v, 255)
				}
			}
		}
	})
}

var stringBenchmarkSink string

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringBenchmarkSink = codecTestUUID.String()
	}
}

func BenchmarkFromBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromBytes(codecTestData)
	}
}

func BenchmarkFromString(b *testing.B) {
	b.Run("canonical", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		}
	})
	b.Run("urn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FromString("urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		}
	})
	b.Run("braced", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FromString("{6ba7b810-9dad-11d1-80b4-00c04fd430c8}")
		}
	})
}

func BenchmarkUnmarshalText(b *testing.B) {
	b.Run("canonical", func(b *testing.B) {
		text := []byte(Must(FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")).String())
		u := new(UUID)
		if err := u.UnmarshalText(text); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = u.UnmarshalText(text)
		}
	})
	b.Run("urn", func(b *testing.B) {
		text := []byte(Must(FromString("urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8")).String())
		u := new(UUID)
		if err := u.UnmarshalText(text); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = u.UnmarshalText(text)
		}
	})
	b.Run("braced", func(b *testing.B) {
		text := []byte(Must(FromString("{6ba7b810-9dad-11d1-80b4-00c04fd430c8}")).String())
		u := new(UUID)
		if err := u.UnmarshalText(text); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = u.UnmarshalText(text)
		}
	})
}

func BenchmarkMarshalBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		codecTestUUID.MarshalBinary()
	}
}

func BenchmarkMarshalText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		codecTestUUID.MarshalText()
	}
}

func BenchmarkParseV4(b *testing.B) {
	const text = "f52a747a-983f-45f7-90b5-e84d70f470dd"
	for i := 0; i < b.N; i++ {
		var u UUID
		if err := u.Parse(text); err != nil {
			b.Fatal(err)
		}
	}
}

const uuidPattern = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"

var fromBytesCorpus = [][]byte{
	{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
	{4, 8, 15, 16, 23, 42},
}

// FuzzFromBytesFunc is a fuzz testing suite that exercises the FromBytes function
func FuzzFromBytesFunc(f *testing.F) {
	for _, seed := range fromBytesCorpus {
		f.Add(seed)
	}
	uuidRegexp, err := regexp.Compile(uuidPattern)
	if err != nil {
		f.Error("uuid regexp failed to compile")
	}
	f.Fuzz(func(t *testing.T, payload []byte) {
		u, err := FromBytes(payload)
		if len(payload) != Size && err == nil {
			t.Errorf("%v did not result in an error", payload)
		}
		if len(payload) == Size && u == Nil {
			t.Errorf("%v resulted in Nil uuid", payload)
		}
		if len(payload) == Size && !uuidRegexp.MatchString(u.String()) {
			t.Errorf("%v resulted in invalid uuid %s", payload, u.String())
		}
		// otherwise, allow to pass if no panic
	})
}

// FuzzFromBytesOrNilFunc is a fuzz testing suite that exercises the FromBytesOrNil function
func FuzzFromBytesOrNilFunc(f *testing.F) {
	for _, seed := range fromBytesCorpus {
		f.Add(seed)
	}
	uuidRegexp, err := regexp.Compile(uuidPattern)
	if err != nil {
		f.Error("uuid regexp failed to compile")
	}
	f.Fuzz(func(t *testing.T, payload []byte) {
		u := FromBytesOrNil(payload)
		if len(payload) != Size && u != Nil {
			t.Errorf("%v resulted in non Nil uuid %s", payload, u.String())
		}
		if len(payload) == Size && u == Nil {
			t.Errorf("%v resulted Nil uuid", payload)
		}
		if len(payload) == Size && !uuidRegexp.MatchString(u.String()) {
			t.Errorf("%v resulted in invalid uuid %s", payload, u.String())
		}
		// otherwise, allow to pass if no panic
	})
}

var fromStringCorpus = []string{
	"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	"6BA7B810-9DAD-11D1-80B4-00C04FD430C8",
	"{6BA7B810-9DAD-11D1-80B4-00C04FD430C8}",
	"urn:uuid:6BA7B810-9DAD-11D1-80B4-00C04FD430C8",
	"6BA7B8109DAD11D180B400C04FD430C8",
	"{6BA7B8109DAD11D180B400C04FD430C8}",
	"urn:uuid:6BA7B8109DAD11D180B400C04FD430C8",
}

// FuzzFromStringFunc is a fuzz testing suite that exercises the FromString function
func FuzzFromStringFunc(f *testing.F) {
	for _, seed := range fromStringCorpus {
		f.Add(seed)
	}
	uuidRegexp, err := regexp.Compile(uuidPattern)
	if err != nil {
		f.Error("uuid regexp failed to compile")
	}
	f.Fuzz(func(t *testing.T, payload string) {
		u, err := FromString(payload)
		if err != nil {
			if u == Nil {
				t.Errorf("%s resulted in Nil uuid", payload)
			}
			if !uuidRegexp.MatchString(u.String()) {
				t.Errorf("%s resulted in invalid uuid %s", payload, u.String())
			}
		}
		// otherwise, allow to pass if no panic
	})
}

// FuzzFromStringOrNil is a fuzz testing suite that exercises the FromStringOrNil function
func FuzzFromStringOrNilFunc(f *testing.F) {
	for _, seed := range fromStringCorpus {
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
