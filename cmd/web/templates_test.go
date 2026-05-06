package main

import (
	"testing"
	"time"

	"github.com/qqquinnn/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Slice of anonymous structs containing test case name, input to function,
	// and expected output.
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2026, 5, 5, 15, 45, 0, 0, time.UTC),
			want: "05 May 2026 at 15:45 UTC",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2026, 5, 5, 15, 45, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "05 May 2026 at 14:45 UTC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hd := humanDate(tc.tm)

			assert.Equal(t, hd, tc.want)
		})
	}
}
