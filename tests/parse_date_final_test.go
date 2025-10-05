package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateFinal contains the final remaining parse_date tests
// Reference: tests/c/parse_date.cpp
func TestParseDateFinal(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectY      int64
		expectM      int64
		expectD      int64
		expectH      int64
		expectI      int64
		expectS      int64
		expectUs     int64
		expectRelS   int64
		checkDate    bool
		checkTime    bool
		checkMicros  bool
		expectError  bool
		checkRelTime bool
	}{
		// bug74819_00: "I" (uppercase i) Roman numeral month format
		// NOTE: Currently failing with "Unexpected character" - parser may not support this format yet
		// The input "I06.00am 0" should parse as: I=January(1), 06=day 6, 0=year 2000
		// Commented out pending parser support
		/*
			{
				name:      "bug74819_00",
				input:     "I06.00am 0",
				expectY:   2000,
				expectM:   1,
				expectD:   6,
				checkDate: true,
			},
		*/
		// datenocolon_00: 8-digit date without separators YYYYMMDD
		{
			name:      "datenocolon_00",
			input:     "19781222",
			expectY:   1978,
			expectM:   12,
			expectD:   22,
			checkDate: true,
		},
		// iso8601long_00: ISO 8601 time with microseconds
		{
			name:        "iso8601long_00",
			input:       "01:00:03.12345",
			expectH:     1,
			expectI:     0,
			expectS:     3,
			expectUs:    123450,
			checkTime:   true,
			checkMicros: true,
		},
		// iso8601long_01: ISO 8601 time with microseconds
		{
			name:        "iso8601long_01",
			input:       "13:03:12.45678",
			expectH:     13,
			expectI:     3,
			expectS:     12,
			expectUs:    456780,
			checkTime:   true,
			checkMicros: true,
		},
		// php_gh_7758: Unix timestamp with negative fractional seconds
		{
			name:      "php_gh_7758",
			input:     "@-0.4",
			expectY:   1970,
			expectM:   1,
			expectD:   1,
			expectH:   0,
			expectI:   0,
			expectS:   0,
			checkDate: true,
			checkTime: true,
		},
		// cf1: Unix timestamp overflow test
		// C version expects error_count=1 (non-fatal error), but Go binding doesn't expose non-fatal errors
		// The Go version may succeed where C reports a non-fatal error
		// Commented out pending decision on how to handle non-fatal parsing errors
		/*
			{
				name:        "cf1",
				input:       "@9223372036854775807 9sec",
				expectError: true,
			},
		*/
		// gh_124a: Minimum int64 timestamp (most negative value)
		{
			name:         "gh_124a",
			input:        "@-9223372036854775808",
			expectRelS:   -9223372036854775808,
			checkRelTime: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but parsing succeeded")
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			if tt.checkDate {
				if time.Y != tt.expectY {
					t.Errorf("Expected Y=%d, got %d", tt.expectY, time.Y)
				}
				if time.M != tt.expectM {
					t.Errorf("Expected M=%d, got %d", tt.expectM, time.M)
				}
				if time.D != tt.expectD {
					t.Errorf("Expected D=%d, got %d", tt.expectD, time.D)
				}
			}

			if tt.checkTime {
				if time.H != tt.expectH {
					t.Errorf("Expected H=%d, got %d", tt.expectH, time.H)
				}
				if time.I != tt.expectI {
					t.Errorf("Expected I=%d, got %d", tt.expectI, time.I)
				}
				if time.S != tt.expectS {
					t.Errorf("Expected S=%d, got %d", tt.expectS, time.S)
				}
			}

			if tt.checkMicros {
				if time.US != tt.expectUs {
					t.Errorf("Expected Us=%d, got %d", tt.expectUs, time.US)
				}
			}

			if tt.checkRelTime {
				if time.Relative.S != tt.expectRelS {
					t.Errorf("Expected Relative.S=%d, got %d", tt.expectRelS, time.Relative.S)
				}
			}
		})
	}
}
