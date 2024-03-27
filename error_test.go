package uuid

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	tcs := []struct {
		err            error
		expected       string
		expectedTarget error
	}{
		{
			err:            fmt.Errorf("%w: sample error: %v", &ErrUUID{}, 123),
			expected:       "uuid: sample error: 123",
			expectedTarget: &ErrUUID{},
		},
		{
			err:            invalidFormat(),
			expected:       "uuid: invalid UUID format",
			expectedTarget: &ErrUUIDInvalidFormat{},
		},
		{
			err:            invalidFormatf("sample error: %v", 123),
			expected:       "uuid: sample error: 123",
			expectedTarget: &ErrUUIDInvalidFormat{},
		},
		{
			err:            fmt.Errorf("uuid error: %w", invalidFormatf("sample error: %v", 123)),
			expected:       "uuid error: uuid: sample error: 123",
			expectedTarget: &ErrUUIDInvalidFormat{},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			if !errors.Is(tc.err, &ErrUUID{}) {
				t.Error("expected error to be of a wrapped type of Error")
			}
			if !errors.Is(tc.err, tc.expectedTarget) {
				t.Errorf("expected error to be of type %v, but was %v", reflect.TypeOf(tc.expectedTarget), reflect.TypeOf(tc.err))
			}
			if tc.err.Error() != tc.expected {
				t.Errorf("expected err.Error() to be '%s' but was '%s'", tc.expected, tc.err.Error())
			}
			uuidErr := &ErrUUID{}
			if !errors.As(tc.err, &uuidErr) {
				t.Error("expected errors.As() to work")
			}
		})
	}
}
