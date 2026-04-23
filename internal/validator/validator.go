package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Contains a map of validation error messages for form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Returns true if the FieldErrors map is empty.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds an error message to the FieldErrors map if given key is unused.
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Adds an error message to the FieldErrors map if a validation check fails.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Returns true if a value contains <n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Returns true if a value is in a list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
