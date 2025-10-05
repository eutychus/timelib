package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateBug50392 tests microsecond precision in date parsing
// Reference: tests/c/parse_date.cpp bug50392_00 to bug50392_08
// Bug 50392: Microsecond precision handling
func TestParseDateBug50392(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expectY  int64
		expectM  int64
		expectD  int64
		expectH  int64
		expectI  int64
		expectS  int64
		expectUS int64
	}{
		{"bug50392_00", "2010-03-06 16:07:25", 2010, 3, 6, 16, 7, 25, 0},
		{"bug50392_01", "2010-03-06 16:07:25.1", 2010, 3, 6, 16, 7, 25, 100000},
		{"bug50392_02", "2010-03-06 16:07:25.12", 2010, 3, 6, 16, 7, 25, 120000},
		{"bug50392_03", "2010-03-06 16:07:25.123", 2010, 3, 6, 16, 7, 25, 123000},
		{"bug50392_04", "2010-03-06 16:07:25.1234", 2010, 3, 6, 16, 7, 25, 123400},
		{"bug50392_05", "2010-03-06 16:07:25.12345", 2010, 3, 6, 16, 7, 25, 123450},
		{"bug50392_06", "2010-03-06 16:07:25.123456", 2010, 3, 6, 16, 7, 25, 123456},
		{"bug50392_07", "2010-03-06 16:07:25.1234567", 2010, 3, 6, 16, 7, 25, 123456},  // truncated
		{"bug50392_08", "2010-03-06 16:07:25.12345678", 2010, 3, 6, 16, 7, 25, 123456}, // truncated
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
			if time.H != tt.expectH {
				t.Errorf("Expected H=%d, got %d", tt.expectH, time.H)
			}
			if time.I != tt.expectI {
				t.Errorf("Expected I=%d, got %d", tt.expectI, time.I)
			}
			if time.S != tt.expectS {
				t.Errorf("Expected S=%d, got %d", tt.expectS, time.S)
			}
			if time.US != tt.expectUS {
				t.Errorf("Expected US=%d, got %d", tt.expectUS, time.US)
			}
		})
	}
}
