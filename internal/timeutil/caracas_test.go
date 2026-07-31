package timeutil

import (
	"testing"
	"time"
)

func TestFechaEfectivaBCVAt(t *testing.T) {
	loc := CaracasLocation()

	cases := []struct {
		name string
		at   time.Time
		want string
	}{
		{
			name: "07:00 guarda hoy",
			at:   time.Date(2026, 7, 30, 7, 0, 0, 0, loc),
			want: "2026-07-30",
		},
		{
			name: "16:59 guarda hoy",
			at:   time.Date(2026, 7, 30, 16, 59, 0, 0, loc),
			want: "2026-07-30",
		},
		{
			name: "17:00 guarda manana",
			at:   time.Date(2026, 7, 30, 17, 0, 0, 0, loc),
			want: "2026-07-31",
		},
		{
			name: "23:30 guarda manana",
			at:   time.Date(2026, 7, 30, 23, 30, 0, 0, loc),
			want: "2026-07-31",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FechaEfectivaBCVAt(tc.at).Format("2006-01-02")
			if got != tc.want {
				t.Fatalf("FechaEfectivaBCVAt(%v) = %s, want %s", tc.at, got, tc.want)
			}
		})
	}
}
