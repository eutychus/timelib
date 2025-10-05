package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateRelativeExtended adds the missing relative tests from C
// C tests relative_00 through relative_60
// Existing TestParseDateRelative has 18 tests
// This file adds 47 more specific tests from C covering:
// - Sign combinations (relative_00 to relative_14)
// - "ago" keyword (relative_15, relative_16, relative_24, relative_30, etc.)
// - Ordinal words (relative_18 to relative_21, relative_41)
// - Spacing variations (relative_28 to relative_36)
// - Weekday tests (relative_42 to relative_58)
// - Complex combinations (relative_40)

// TestParseDateRelativeSigns tests sign combinations (+, -, ++, --, +-, -+, etc.)
// C tests relative_00 through relative_14
func TestParseDateRelativeSigns(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectRelS int64
	}{
		{"relative_00", "2 secs", 2},
		{"relative_01", "+2 sec", 2},
		{"relative_02", "-2 secs", -2},
		{"relative_03", "++2 sec", 2},
		{"relative_04", "+-2 secs", -2},
		{"relative_05", "-+2 sec", -2},
		{"relative_06", "--2 secs", 2},
		{"relative_07", "+++2 sec", 2},
		{"relative_08", "++-2 secs", -2},
		{"relative_09", "+-+2 sec", -2},
		{"relative_10", "+--2 secs", 2},
		{"relative_11", "-++2 sec", -2},
		{"relative_12", "-+-2 secs", 2},
		{"relative_13", "--+2 sec", 2},
		{"relative_14", "---2 secs", -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			if time.Relative.Y != 0 {
				t.Errorf("Expected Relative.Y=0, got %d", time.Relative.Y)
			}
			if time.Relative.M != 0 {
				t.Errorf("Expected Relative.M=0, got %d", time.Relative.M)
			}
			if time.Relative.D != 0 {
				t.Errorf("Expected Relative.D=0, got %d", time.Relative.D)
			}
			if time.Relative.H != 0 {
				t.Errorf("Expected Relative.H=0, got %d", time.Relative.H)
			}
			if time.Relative.I != 0 {
				t.Errorf("Expected Relative.I=0, got %d", time.Relative.I)
			}
			if time.Relative.S != tt.expectRelS {
				t.Errorf("Expected Relative.S=%d, got %d", tt.expectRelS, time.Relative.S)
			}
		})
	}
}

// TestParseDateRelativeAgo tests "ago" keyword that negates the sign
// C tests relative_15, relative_16, relative_24, relative_30
func TestParseDateRelativeAgo(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectRelY int64
		expectRelM int64
		expectRelD int64
		expectRelH int64
		expectRelI int64
		expectRelS int64
	}{
		{"relative_15", "+2 sec ago", 0, 0, 0, 0, 0, -2},
		{"relative_16", "2 secs ago", 0, 0, 0, 0, 0, -2},
		{"relative_24", "+2 days ago", 0, 0, -2, 0, 0, 0},
		{"relative_30", "+ 2 days ago", 0, 0, -2, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

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
		})
	}
}

// TestParseDateRelativeOrdinal tests ordinal words (first, second, third, next)
// C tests relative_18 through relative_22, relative_41
func TestParseDateRelativeOrdinal(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectRelY int64
		expectRelM int64
		expectRelD int64
		expectRelH int64
		expectRelI int64
		expectRelS int64
	}{
		{"relative_18", "first second", 0, 0, 0, 0, 0, 1},
		{"relative_19", "next second", 0, 0, 0, 0, 0, 1},
		{"relative_20", "second second", 0, 0, 0, 0, 0, 2},
		{"relative_21", "third second", 0, 0, 0, 0, 0, 3},
		{"relative_22", "-3 seconds", 0, 0, 0, 0, 0, -3},
		{"relative_41", "first month", 0, 1, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

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
		})
	}
}

// TestParseDateRelativeUnits tests various units (days, fortnight, weeks)
// C tests relative_23, relative_25, relative_26, relative_27
func TestParseDateRelativeUnits(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectRelY int64
		expectRelM int64
		expectRelD int64
		expectRelH int64
		expectRelI int64
		expectRelS int64
	}{
		{"relative_23", "+2 days", 0, 0, 2, 0, 0, 0},
		{"relative_25", "-2 days", 0, 0, -2, 0, 0, 0},
		{"relative_26", "-3 fortnight", 0, 0, -42, 0, 0, 0},
		{"relative_27", "+12 weeks", 0, 0, 84, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

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
		})
	}
}

// TestParseDateRelativeSpacing tests spacing variations (space after sign, multiple spaces)
// C tests relative_28 through relative_36
func TestParseDateRelativeSpacing(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectRelY int64
		expectRelM int64
		expectRelD int64
		expectRelH int64
		expectRelI int64
		expectRelS int64
	}{
		{"relative_28", "- 3 seconds", 0, 0, 0, 0, 0, -3},
		{"relative_29", "+ 2 days", 0, 0, 2, 0, 0, 0},
		{"relative_31", "- 2 days", 0, 0, -2, 0, 0, 0},
		{"relative_32", "- 3 fortnight", 0, 0, -42, 0, 0, 0},
		{"relative_33", "+ 12 weeks", 0, 0, 84, 0, 0, 0},
		{"relative_34", "- 2 days", 0, 0, -2, 0, 0, 0},
		{"relative_35", "-   3 fortnight", 0, 0, -42, 0, 0, 0},
		{"relative_36", "+   12 weeks", 0, 0, 84, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

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
		})
	}
}

// TestParseDateRelativeComplex tests complex combinations (relative with absolute dates)
// C test relative_40
func TestParseDateRelativeComplex(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectRelY int64
		expectRelM int64
		expectRelD int64
		expectRelH int64
		expectRelI int64
		expectRelS int64
	}{
		{"relative_40", "6 months ago 4 days", 0, -6, 4, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

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
		})
	}
}

// TestParseDateRelativeWeekday tests weekday relative expressions
// C tests relative_42 through relative_51
// Note: These tests check Relative.Weekday and Relative.WeekdayBehavior fields
func TestParseDateRelativeWeekday(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectRelD            int64
		expectWeekday         int
		expectWeekdayBehavior int
	}{
		{"relative_42", "saturday", 0, 6, 1},
		{"relative_43", "saturday ago", 0, -6, 1},
		{"relative_44", "this saturday", 0, 6, 1},
		{"relative_45", "this saturday ago", 0, -6, 1},
		{"relative_46", "last saturday", -7, 6, 0},
		{"relative_47", "last saturday ago", 7, -6, 0},
		{"relative_48", "first saturday", 0, 6, 0},
		{"relative_49", "first saturday ago", 0, -6, 0},
		{"relative_50", "next saturday", 0, 6, 0},
		{"relative_51", "next saturday ago", 0, -6, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			// Check time fields (should be 0)
			if time.H != 0 {
				t.Errorf("Expected H=0, got %d", time.H)
			}
			if time.I != 0 {
				t.Errorf("Expected I=0, got %d", time.I)
			}
			if time.S != 0 {
				t.Errorf("Expected S=0, got %d", time.S)
			}

			// Check relative D
			if time.Relative.D != tt.expectRelD {
				t.Errorf("Expected Relative.D=%d, got %d", tt.expectRelD, time.Relative.D)
			}

			// Check weekday fields
			if time.Relative.Weekday != tt.expectWeekday {
				t.Errorf("Expected Relative.Weekday=%d, got %d", tt.expectWeekday, time.Relative.Weekday)
			}
			if time.Relative.WeekdayBehavior != tt.expectWeekdayBehavior {
				t.Errorf("Expected Relative.WeekdayBehavior=%d, got %d", tt.expectWeekdayBehavior, time.Relative.WeekdayBehavior)
			}
		})
	}
}

// TestParseDateRelativeMissing adds the 13 missing relative tests
// relative_17, relative_37-39, relative_52-60
func TestParseDateRelativeMissing(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectY               int64
		expectM               int64
		expectD               int64
		expectH               int64
		expectI               int64
		expectS               int64
		expectZ               int32
		expectDst             int
		checkAbsolute         bool
		checkTZ               bool
		expectRelY            int64
		expectRelM            int64
		expectRelD            int64
		expectRelH            int64
		expectRelI            int64
		expectRelS            int64
		checkRelative         bool
		expectWeekday         int
		expectWeekdayBehavior int
		checkWeekday          bool
		expectSpecialType     int
		expectSpecialAmount   int64
		checkSpecial          bool
	}{
		// relative_17: parse-only test (no expectations)
		{"relative_17", "0 second", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, 0, 0, 0, 0, false, 0, 0, false, 0, 0, false},

		// relative_37-39: Combined absolute + relative tests
		{"relative_37", "6 month 2004-05-05 12:15:23 CEST", 2004, 5, 5, 12, 15, 23, 3600, 1, true, true, 0, 6, 0, 0, 0, 0, true, 0, 0, false, 0, 0, false},
		{"relative_38", "2004-05-05 12:15:23 CEST 6 months", 2004, 5, 5, 12, 15, 23, 3600, 1, true, true, 0, 6, 0, 0, 0, 0, true, 0, 0, false, 0, 0, false},
		{"relative_39", "2004-05-05 12:15:23 CEST 6 months ago", 2004, 5, 5, 12, 15, 23, 3600, 1, true, true, 0, -6, 0, 0, 0, 0, true, 0, 0, false, 0, 0, false},

		// relative_52-54: Weekday with ordinals
		{"relative_52", "third saturday", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, 14, 0, 0, 0, true, 6, 0, true, 0, 0, false},
		{"relative_53", "third saturday ago", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, -14, 0, 0, 0, true, -6, 0, true, 0, 0, false},
		{"relative_54", "previous saturday", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, -7, 0, 0, 0, true, 6, 0, true, 0, 0, false},

		// relative_55-58: Special weekday tests
		{"relative_55", "this weekday", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, 0, 0, 0, 0, false, 0, 0, false, timelib.TIMELIB_SPECIAL_WEEKDAY, 0, true},
		{"relative_56", "last weekday", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, 0, 0, 0, 0, false, 0, 0, false, timelib.TIMELIB_SPECIAL_WEEKDAY, -1, true},
		{"relative_57", "next weekday", 0, 0, 0, 0, 0, 0, 0, 0, false, false, 0, 0, 0, 0, 0, 0, false, 0, 0, false, timelib.TIMELIB_SPECIAL_WEEKDAY, 1, true},
		{"relative_58", "8 weekdays ago", timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, timelib.TIMELIB_UNSET, 0, 0, true, false, 0, 0, 0, 0, 0, 0, false, 0, 0, false, timelib.TIMELIB_SPECIAL_WEEKDAY, -8, true},

		// relative_59-60: Combined with weekday_behavior
		{"relative_59", "Sun, 21 Dec 2003 20:38:33 +0000 GMT", 2003, 12, 21, 20, 38, 33, 0, 0, true, false, 0, 0, 0, 0, 0, 0, false, 0, 1, true, 0, 0, false},
		{"relative_60", "Mon, 08 May 2006 13:06:44 -0400 +30 days", 2006, 5, 8, 13, 6, 44, -14400, 0, true, true, 0, 0, 30, 0, 0, 0, true, 1, 1, true, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			// Check absolute time fields if specified
			if tt.checkAbsolute {
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
			}

			// Check timezone fields if specified
			if tt.checkTZ {
				if time.Z != tt.expectZ {
					t.Errorf("Expected Z=%d, got %d", tt.expectZ, time.Z)
				}
				if time.Dst != tt.expectDst {
					t.Errorf("Expected Dst=%d, got %d", tt.expectDst, time.Dst)
				}
			}

			// Check relative fields if specified
			if tt.checkRelative {
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
			}

			// Check weekday fields if specified
			if tt.checkWeekday {
				if time.Relative.Weekday != tt.expectWeekday {
					t.Errorf("Expected Relative.Weekday=%d, got %d", tt.expectWeekday, time.Relative.Weekday)
				}
				if time.Relative.WeekdayBehavior != tt.expectWeekdayBehavior {
					t.Errorf("Expected Relative.WeekdayBehavior=%d, got %d", tt.expectWeekdayBehavior, time.Relative.WeekdayBehavior)
				}
			}

			// Check special fields if specified
			if tt.checkSpecial {
				if time.Relative.Special.Type != tt.expectSpecialType {
					t.Errorf("Expected Relative.Special.Type=%d, got %d", tt.expectSpecialType, time.Relative.Special.Type)
				}
				if time.Relative.Special.Amount != tt.expectSpecialAmount {
					t.Errorf("Expected Relative.Special.Amount=%d, got %d", tt.expectSpecialAmount, time.Relative.Special.Amount)
				}
			}
		})
	}
}
