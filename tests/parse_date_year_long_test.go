package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateYearLong tests parsing dates with very large year values
// Reference: tests/c/parse_date.cpp year_long_00 to year_long_09
func TestParseDateYearLong(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
		expectH int64
		expectI int64
		expectS int64
	}{
		{"year_long_00", "+10000-01-01T00:00:00", 10000, 1, 1, 0, 0, 0},
		{"year_long_01", "+99999-01-01T00:00:00", 99999, 1, 1, 0, 0, 0},
		{"year_long_02", "+100000-01-01T00:00:00", 100000, 1, 1, 0, 0, 0},
		{"year_long_03", "+4294967296-01-01T00:00:00", 4294967296, 1, 1, 0, 0, 0},
		{"year_long_04", "+9223372036854775807-01-01T00:00:00", 9223372036854775807, 1, 1, 0, 0, 0},
		{"year_long_05", "-10000-01-01T00:00:00", -10000, 1, 1, 0, 0, 0},
		{"year_long_06", "-99999-01-01T00:00:00", -99999, 1, 1, 0, 0, 0},
		{"year_long_07", "-100000-01-01T00:00:00", -100000, 1, 1, 0, 0, 0},
		{"year_long_08", "-4294967296-01-01T00:00:00", -4294967296, 1, 1, 0, 0, 0},
		{"year_long_09", "-9223372036854775807-01-01T00:00:00", -9223372036854775807, 1, 1, 0, 0, 0},
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
		})
	}
}
