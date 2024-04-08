package uuid

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	var e Error
	tcs := []struct {
		err            error
		expected       string
		expectedTarget error
	}{
		{
			err:            fmt.Errorf("%w sample error: %v", ErrInvalidVersion, 123),
			expected:       "uuid: sample error: 123",
			expectedTarget: &e,
		},
		{
			err:            fmt.Errorf("%w", ErrInvalidFormat),
			expected:       "uuid: invalid UUID format",
			expectedTarget: ErrInvalidFormat,
		},
		{
			err:            fmt.Errorf("%w %q", ErrIncorrectFormatInString, "test"),
			expected:       "uuid: incorrect UUID format in string \"test\"",
			expectedTarget: ErrIncorrectFormatInString,
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			var e2 Error
			if !errors.Is(tc.err, &e2) {
				t.Error("expected error to be of a wrapped type of Error")
			}
			if !errors.Is(tc.err, tc.expectedTarget) {
				t.Errorf("expected error to be of type %v, but was %v", reflect.TypeOf(tc.expectedTarget), reflect.TypeOf(tc.err))
			}
			if tc.err.Error() != tc.expected {
				t.Errorf("expected err.Error() to be '%s' but was '%s'", tc.expected, tc.err.Error())
			}
			var uuidErr Error
			if !errors.As(tc.err, &uuidErr) {
				t.Error("expected errors.As() to work")
			}
		})
	}
}

func TestAllErrorMessages(t *testing.T) {
	tcs := []struct {
		function string
		uuidStr  string
		expected string
	}{
		{ // 34 chars - With brackets
			function: "parse",
			uuidStr:  "..................................",
			expected: "uuid: incorrect UUID format in string \"..................................\"",
		},
		{ // 41 chars - urn:uuid:
			function: "parse",
			uuidStr:  "123456789................................",
			expected: "uuid: incorrect UUID format in string \"123456789\"",
		},
		{ // other
			function: "parse",
			uuidStr:  "....",
			expected: "uuid: incorrect UUID length 4 in string \"....\"",
		},
		{ // 36 chars - canonical, but not correct format
			function: "parse",
			uuidStr:  "....................................",
			expected: "uuid: incorrect UUID format in string \"....................................\"",
		},
		{ // 36 chars - canonical, invalid data
			function: "parse",
			uuidStr:  "xx00ae9e-dae3-459f-ad0e-6b574be3f950",
			expected: "uuid: invalid UUID format",
		},
		{ // Hash like
			function: "parse",
			uuidStr:  "................................",
			expected: "uuid: invalid UUID format",
		},
		{ // Hash like, invalid
			function: "parse",
			uuidStr:  "xx00ae9edae3459fad0e6b574be3f950",
			expected: "uuid: invalid UUID format",
		},
		{ // Hash like, invalid
			function: "parse",
			uuidStr:  "xx00ae9edae3459fad0e6b574be3f950",
			expected: "uuid: invalid UUID format",
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			id := UUID{}
			err := id.Parse(tc.uuidStr)
			if err == nil {
				t.Error("expected an error")
				return
			}
			if err.Error() != tc.expected {
				t.Errorf("unexpected error '%s' != '%s'", err.Error(), tc.expected)
			}
			err = id.UnmarshalText([]byte(tc.uuidStr))
			if err == nil {
				t.Error("expected an error")
				return
			}
			if err.Error() != tc.expected {
				t.Errorf("unexpected error '%s' != '%s'", err.Error(), tc.expected)
			}
		})
	}

	// Unmarshal binary
	id := UUID{}
	b := make([]byte, 33)
	expectedErr := "uuid: UUID must be exactly 16 bytes long, got 33 bytes"
	err := id.UnmarshalBinary([]byte(b))
	if err == nil {
		t.Error("expected an error")
		return
	}
	if err.Error() != expectedErr {
		t.Errorf("unexpected error '%s' != '%s'", err.Error(), expectedErr)
	}

	// no hw address error
	netInterfaces = func() ([]net.Interface, error) {
		return nil, nil
	}
	defer func() {
		netInterfaces = net.Interfaces
	}()
	_, err = defaultHWAddrFunc()
	if err == nil {
		t.Error("expected an error")
		return
	}
	expectedErr = "uuid: no HW address found"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error '%s' != '%s'", err.Error(), expectedErr)
	}

	// scan error
	err = id.Scan(123)
	if err == nil {
		t.Error("expected an error")
		return
	}
	expectedErr = "uuid: cannot convert int to UUID"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error '%s' != '%s'", err.Error(), expectedErr)
	}

	// UUId V1 Version
	id = FromStringOrNil("e86160d3-beff-443c-b9b5-1f8197ccb12e")
	_, err = TimestampFromV1(id)
	if err == nil {
		t.Error("expected an error")
		return
	}
	expectedErr = "uuid: e86160d3-beff-443c-b9b5-1f8197ccb12e is version 4, not version 1"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error '%s' != '%s'", err.Error(), expectedErr)
	}

	// UUId V2 Version
	id = FromStringOrNil("e86160d3-beff-443c-b9b5-1f8197ccb12e")
	_, err = TimestampFromV6(id)
	if err == nil {
		t.Error("expected an error")
		return
	}
	expectedErr = "uuid: e86160d3-beff-443c-b9b5-1f8197ccb12e is version 4, not version 6"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error '%s' != '%s'", err.Error(), expectedErr)
	}

	// UUId V7 Version
	id = FromStringOrNil("e86160d3-beff-443c-b9b5-1f8197ccb12e")
	_, err = TimestampFromV7(id)
	if err == nil {
		t.Error("expected an error")
		return
	}
	// There is a "bug" in the error message, this should probably be fixed (6 -> 7)
	expectedErr = "uuid: e86160d3-beff-443c-b9b5-1f8197ccb12e is version 4, not version 6"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error '%s' != '%s'", err.Error(), expectedErr)
	}
}
