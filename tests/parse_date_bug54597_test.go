package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateBug54597 tests parsing dates with leading zero years
// Reference: tests/c/parse_date.cpp bug54597_00 to bug54597_07
// Bug 54597: 4-digit years with leading zeros (0099) should be parsed as-is, not as 2-digit years
func TestParseDateBug54597(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"bug54597_00", "January 0099", 99, 1, 1},
		{"bug54597_01", "January 1, 0099", 99, 1, 1},
		{"bug54597_02", "0099-1", 99, 1, 1},
		{"bug54597_03", "0099-January", 99, 1, 1},
		{"bug54597_04", "0099-Jan", 99, 1, 1},
		{"bug54597_05", "January 1099", 1099, 1, 1},
		{"bug54597_06", "January 1, 1299", 1299, 1, 1},
		{"bug54597_07", "1599-1", 1599, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			if time.Y != tt.expectY {
				t.Errorf("Expected Y=%d, got %d", tt.expectY, time.Y)
			}
			if time.M != tt.expectM {
				t.Errorf("Expected M=%d, got %d", tt.expectM, time.M)
			}
			if time.D != tt.expectD {
				t.Errorf("Expected D=%d, got %d", tt.expectD, time.D)
			}
		})
	}
}
