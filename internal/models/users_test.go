package models

import (
	"testing"

	"github.com/qqquinnn/snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	// Skip test if "-short" flag is provided when running.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Set up suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get connection pool to test database. Data reset for each sub-test.
			db := newTestDB(t)

			// New instance of UserModel.
			m := UserModel{db}

			// Check that return value matches expected values w/ no error.
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.Nil(t, err)
		})
	}
}
