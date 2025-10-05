package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateBug51096 tests "first day" / "last day" relative parsing
// Reference: tests/c/parse_date.cpp bug51096_00 to bug51096_06
// Bug 51096: Proper handling of "first day", "last day" with month modifiers
func TestParseDateBug51096(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectRelY           int64
		expectRelM           int64
		expectRelD           int64
		expectRelH           int64
		expectRelI           int64
		expectRelS           int64
		expectFirstLastDayOf int
		expectH              int64
		expectI              int64
		expectS              int64
		checkFirstLastDayOf  bool
	}{
		{"bug51096_00", "first day", 0, 0, 1, 0, 0, 0, 0, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, false},
		{"bug51096_01", "last day", 0, 0, -1, 0, 0, 0, 0, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, false},
		{"bug51096_02", "next month", 0, 1, 0, 0, 0, 0, 0, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, false},
		{"bug51096_03", "first day next month", 0, 1, 1, 0, 0, 0, 0, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, false},
		{"bug51096_04", "first day of next month", 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, true},
		{"bug51096_05", "last day next month", 0, 1, -1, 0, 0, 0, 0, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, false},
		{"bug51096_06", "last day of next month", 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			// Check relative fields
			if time.Relative.Y != tt.expectRelY {
				t.Errorf("Expected Relative.Y=%d, got %d", tt.expectRelY, time.Relative.Y)
			}
			if time.Relative.M != tt.expectRelM {
				t.Errorf("Expected Relative.M=%d, got %d", tt.expectRelM, time.Relative.M)
			}
			if time.Relative.D != tt.expectRelD {
				t.Errorf("Expected Relative.D=%d, got %d", tt.expectRelD, time.Relative.D)
			}
			if time.Relative.H != tt.expectRelH {
				t.Errorf("Expected Relative.H=%d, got %d", tt.expectRelH, time.Relative.H)
			}
			if time.Relative.I != tt.expectRelI {
				t.Errorf("Expected Relative.I=%d, got %d", tt.expectRelI, time.Relative.I)
			}
			if time.Relative.S != tt.expectRelS {
				t.Errorf("Expected Relative.S=%d, got %d", tt.expectRelS, time.Relative.S)
			}

			// Check FirstLastDayOf if specified
			if tt.checkFirstLastDayOf {
				if time.Relative.FirstLastDayOf != tt.expectFirstLastDayOf {
					t.Errorf("Expected Relative.FirstLastDayOf=%d, got %d", tt.expectFirstLastDayOf, time.Relative.FirstLastDayOf)
				}
			}

			// Check absolute time fields
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
