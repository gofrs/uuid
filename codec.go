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

import "fmt"

// FromBytes returns a UUID generated from the raw byte slice input.
// It will return an error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (UUID, error) {
	u := UUID{}
	err := u.UnmarshalBinary(input)
	return u, err
}

// FromBytesOrNil returns a UUID generated from the raw byte slice input.
// Same behavior as FromBytes(), but returns uuid.Nil instead of an error.
func FromBytesOrNil(input []byte) UUID {
	// The logic here is duplicated from UnmarshalBinary as there is unnecessary
	// overhead generating errors which would be checked and discarded.
	if len(input) != Size {
		return Nil
	}

	uuid := UUID{}
	copy(uuid[:], input)

	return uuid
}

func fromHexChar(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 255
}

// parseBytes parses UUID text representation from a byte slice and populates
// the provided UUID reference. It centralizes the parsing logic so that both
// Parse (string input) and UnmarshalText ([]byte input) can delegate to this
// single implementation, eliminating code duplication.
//
// Supported formats and validation rules are identical to those documented
// in UnmarshalText.
func parseBytes(b []byte, u *UUID) error {
	// Fast-path: ensure we don't accidentally mutate the caller's slice.
	// We will only reslice, never modify the underlying bytes.
	switch len(b) {
	case 32: // hash
	case 36: // canonical
	case 34, 38:
		if b[0] != '{' || b[len(b)-1] != '}' {
			return fmt.Errorf("%w %q", ErrIncorrectFormatInString, b)
		}
		b = b[1 : len(b)-1]
	case 41, 45:
		if string(b[:9]) != "urn:uuid:" {
			return fmt.Errorf("%w %q", ErrIncorrectFormatInString, b[:9])
		}
		b = b[9:]
	default:
		return fmt.Errorf("%w %d in string %q", ErrIncorrectLength, len(b), b)
	}

	// canonical (36 chars with dashes at fixed positions)
	if len(b) == 36 {
		if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
			return fmt.Errorf("%w %q", ErrIncorrectFormatInString, b)
		}
		for i, x := range [16]byte{
			0, 2, 4, 6,
			9, 11,
			14, 16,
			19, 21,
			24, 26, 28, 30, 32, 34,
		} {
			v1 := fromHexChar(b[x])
			v2 := fromHexChar(b[x+1])
			if v1|v2 == 255 {
				return ErrInvalidFormat
			}
			u[i] = (v1 << 4) | v2
		}
		return nil
	}

	// hash-like (32 hex chars, no dashes)
	for i := 0; i < 32; i += 2 {
		v1 := fromHexChar(b[i])
		v2 := fromHexChar(b[i+1])
		if v1|v2 == 255 {
			return ErrInvalidFormat
		}
		u[i/2] = (v1 << 4) | v2
	}
	return nil
}

// Parse parses the UUID stored in the string text. Parsing and supported
// formats are the same as UnmarshalText.
func (u *UUID) Parse(s string) error {
	return parseBytes([]byte(s), u)
}

// FromString returns a UUID parsed from the input string.
// Input is expected in a form accepted by UnmarshalText.
func FromString(text string) (UUID, error) {
	var u UUID
	err := u.Parse(text)
	return u, err
}

// FromStringOrNil returns a UUID parsed from the input string.
// Same behavior as FromString(), but returns uuid.Nil instead of an error.
func FromStringOrNil(input string) UUID {
	uuid, err := FromString(input)
	if err != nil {
		return Nil
	}
	return uuid
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by the String() method.
func (u UUID) MarshalText() ([]byte, error) {
	var buf [36]byte
	encodeCanonical(buf[:], u)
	return buf[:], nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It delegates the actual parsing to parseBytes to avoid code duplication.
func (u *UUID) UnmarshalText(b []byte) error {
	return parseBytes(b, u)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() ([]byte, error) {
	return u.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return an error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != Size {
		return fmt.Errorf("%w, got %d bytes", ErrIncorrectByteLength, len(data))
	}
	copy(u[:], data)

	return nil
}
