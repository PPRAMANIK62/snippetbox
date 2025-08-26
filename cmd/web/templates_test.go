package main

import (
	"testing"
	"time"

	"github.com/PPRAMANIK62/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 8, 25, 15, 30, 0, 0, time.UTC),
			want: "25 Aug 2025 at 15:30",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm: time.Date(2025, 8, 25, 15, 30, 0, 0,time.FixedZone("CET", 1*60*60)),
			want: "25 Aug 2025 at 14:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
