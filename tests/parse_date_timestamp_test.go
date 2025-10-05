package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateTimestamp tests Unix timestamp parsing with microseconds
// Reference: tests/c/parse_date.cpp timestamp_00 to timestamp_09
func TestParseDateTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectY     int64
		expectM     int64
		expectD     int64
		expectH     int64
		expectI     int64
		expectS     int64
		expectRelS  int64
		expectRelUS int64
	}{
		{"timestamp_00", "@1508765076.3", 1970, 1, 1, 0, 0, 0, 1508765076, 300000},
		{"timestamp_01", "@1508765076.34", 1970, 1, 1, 0, 0, 0, 1508765076, 340000},
		{"timestamp_02", "@1508765076.347", 1970, 1, 1, 0, 0, 0, 1508765076, 347000},
		{"timestamp_03", "@1508765076.3479", 1970, 1, 1, 0, 0, 0, 1508765076, 347900},
		{"timestamp_04", "@1508765076.34795", 1970, 1, 1, 0, 0, 0, 1508765076, 347950},
		{"timestamp_05", "@1508765076.347958", 1970, 1, 1, 0, 0, 0, 1508765076, 347958},
		{"timestamp_06", "@1508765076.003", 1970, 1, 1, 0, 0, 0, 1508765076, 3000},
		{"timestamp_07", "@1508765076.0003", 1970, 1, 1, 0, 0, 0, 1508765076, 300},
		{"timestamp_08", "@1508765076.00003", 1970, 1, 1, 0, 0, 0, 1508765076, 30},
		{"timestamp_09", "@1508765076.000003", 1970, 1, 1, 0, 0, 0, 1508765076, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			// Check absolute time (should be epoch)
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

			// Check relative fields (the actual timestamp)
			if time.Relative.S != tt.expectRelS {
				t.Errorf("Expected Relative.S=%d, got %d", tt.expectRelS, time.Relative.S)
			}
			if time.Relative.US != tt.expectRelUS {
				t.Errorf("Expected Relative.US=%d, got %d", tt.expectRelUS, time.Relative.US)
			}
		})
	}
}
