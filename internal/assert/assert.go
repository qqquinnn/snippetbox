package assert

import (
	"reflect"
	"testing"
)

// Determines if two values are equivalent.
func Equal[T any](t *testing.T, got, want T) {
	t.Helper()

	if !isEqual(got, want) {
		t.Errorf("got: %v; want: %v", got, want)
	}
}

func isEqual[T any](got, want T) bool {
	if isNil(got) && isNil(want) {
		return true
	}

	return reflect.DeepEqual(got, want)
}

// Determines if two values are not equivalent.
func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got == want {
		t.Errorf("got %v; expected values to be different", got)
	}
}

// Checks that a value is true.
func True(t *testing.T, got bool) {
	t.Helper()

	if !got {
		t.Errorf("got: false; want: true")
	}
}

// Checks that a value is false.
func False(t *testing.T, got bool) {
	t.Helper()

	if got {
		t.Errorf("got: true; want: false")
	}
}

// Checks that a value is nil.
func Nil(t *testing.T, got any) {
	t.Helper()

	if got != nil {
		t.Errorf("got: %v; want: nil", got)
	}
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	}

	return false
}

// Checks that a value is not nil.
func NotNil(t *testing.T, got any) {
	t.Helper()

	if got == nil {
		t.Errorf("got: nil; want: non-nil")
	}
}
