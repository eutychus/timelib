package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateMicrosecondRelative tests relative microsecond expressions
// C tests microsecond_00 through microsecond_11
// These test millisecond (ms) and microsecond (µs, usec) relative values

func TestParseDateMicrosecondRelative(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expectUS int64
	}{
		// milliseconds (1 ms = 1000 microseconds)
		{"microsecond_00", "+1 ms", 1000},
		{"microsecond_01", "+3 msec", 3000},
		{"microsecond_02", "+4 msecs", 4000},
		{"microsecond_03", "+5 millisecond", 5000},
		{"microsecond_04", "+6 milliseconds", 6000},
		// microseconds
		{"microsecond_05", "+1 µs", 1},
		{"microsecond_06", "+3 usec", 3},
		{"microsecond_07", "+4 usecs", 4},
		{"microsecond_08", "+5 µsec", 5},
		{"microsecond_09", "+6 µsecs", 6},
		{"microsecond_10", "+7 microsecond", 7},
		{"microsecond_11", "+8 microseconds", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			if time.Relative.US != tt.expectUS {
				t.Errorf("Expected Relative.US=%d, got %d", tt.expectUS, time.Relative.US)
			}
		})
	}
}
