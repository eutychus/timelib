package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateDateTextualExtended adds the missing datetextual tests from C
// C tests datetextual_00 through datetextual_12
// Existing TestParseDateTextualBasic only has 7 generic tests
// These are the specific format variations from the C tests
func TestParseDateDateTextualExtended(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"datetextual_00", "December 22, 1978", 1978, 12, 22},
		{"datetextual_01", "DECEMBER 22nd 1978", 1978, 12, 22},
		{"datetextual_02", "December 22. 1978", 1978, 12, 22},
		{"datetextual_03", "December 22 1978", 1978, 12, 22},
		{"datetextual_04", "Dec 22, 1978", 1978, 12, 22},
		{"datetextual_05", "DEC 22nd 1978", 1978, 12, 22},
		{"datetextual_06", "Dec 22. 1978", 1978, 12, 22},
		{"datetextual_07", "Dec 22 1978", 1978, 12, 22},
		{"datetextual_08", "December 22", timelib.TIMELIB_UNSET, 12, 22},
		{"datetextual_09", "Dec 22", timelib.TIMELIB_UNSET, 12, 22},
		{"datetextual_10", "DEC 22nd", timelib.TIMELIB_UNSET, 12, 22},
		{"datetextual_11", "December\t22\t1978", 1978, 12, 22},
		{"datetextual_12", "DEC\t22nd", timelib.TIMELIB_UNSET, 12, 22},
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

// TestParseDateRomanExtended adds the missing dateroman tests from C
// C tests dateroman_00 through dateroman_11 (formats like "22 I 1978")
// Existing TestParseDateRoman has 3 tests with different format (1978-XII-22)
func TestParseDateRomanExtended(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"dateroman_00", "22 I 1978", 1978, 1, 22},
		{"dateroman_01", "22. II 1978", 1978, 2, 22},
		{"dateroman_02", "22 III. 1978", 1978, 3, 22},
		{"dateroman_03", "22- IV- 1978", 1978, 4, 22},
		{"dateroman_04", "22 -V -1978", 1978, 5, 22},
		{"dateroman_05", "22-VI-1978", 1978, 6, 22},
		{"dateroman_06", "22.VII.1978", 1978, 7, 22},
		{"dateroman_07", "22 VIII 1978", 1978, 8, 22},
		{"dateroman_08", "22 IX 1978", 1978, 9, 22},
		{"dateroman_09", "22 X 1978", 1978, 10, 22},
		{"dateroman_10", "22 XI 1978", 1978, 11, 22},
		{"dateroman_11", "22\tXII\t1978", 1978, 12, 22},
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

// TestParseDateISO8601DateExtended adds the missing iso8601date tests from C
// C tests iso8601date_00 through iso8601date_10
// Existing TestParseDateISO8601 has some overlap but missing specific cases
func TestParseDateISO8601DateExtended(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"iso8601date_00", "1978-12-22", 1978, 12, 22},
		{"iso8601date_01", "0078-12-22", 78, 12, 22},   // 4-digit year with leading zero
		{"iso8601date_02", "078-12-22", 1978, 12, 22},  // 3-digit year
		{"iso8601date_03", "78-12-22", 1978, 12, 22},   // 2-digit year
		{"iso8601date_04", "4-4-25", 2004, 4, 25},      // Single digit year/month/day
		{"iso8601date_05", "69-4-25", 2069, 4, 25},     // Year 69 -> 2069
		{"iso8601date_06", "70-4-25", 1970, 4, 25},     // Year 70 -> 1970
		{"iso8601date_07", "1978/12/22", 1978, 12, 22}, // Slash separator
		{"iso8601date_08", "1978/02/02", 1978, 2, 2},
		{"iso8601date_09", "1978/12/02", 1978, 12, 2},
		{"iso8601date_10", "1978/02/22", 1978, 2, 22},
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

// TestParseDateNoDayExtended adds the missing datenoday tests from C
// C tests datenoday_00 through datenoday_09
// Existing TestParseDateNoDay only has 1 test
// Note: checkTime=true for HMS fields, checkTZ=true for timezone fields (Z, Dst, TzAbbr)
func TestParseDateNoDayExtended(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectY   int64
		expectM   int64
		expectD   int64
		checkTime bool // Only check time fields if they're provided in input
		expectH   int64
		expectI   int64
		expectS   int64
		checkTZ   bool // Only check timezone fields if timezone is in input
		expectZ   int32
		expectDst int
		expectTZ  string
	}{
		{"datenoday_00", "Oct 2003", 2003, 10, 1, false, 0, 0, 0, false, 0, 0, ""},
		{"datenoday_01", "20 October 2003", 2003, 10, 20, false, 0, 0, 0, false, 0, 0, ""},
		{"datenoday_02", "Oct 03", timelib.TIMELIB_UNSET, 10, 3, false, 0, 0, 0, false, 0, 0, ""},
		{"datenoday_03", "Oct 2003 2045", 2003, 10, 1, true, 20, 45, 0, false, 0, 0, ""},
		{"datenoday_04", "Oct 2003 20:45", 2003, 10, 1, true, 20, 45, 0, false, 0, 0, ""},
		{"datenoday_05", "Oct 2003 20:45:37", 2003, 10, 1, true, 20, 45, 37, false, 0, 0, ""},
		{"datenoday_06", "20 October 2003 00:00 CEST", 2003, 10, 20, true, 0, 0, 0, true, 3600, 1, "CEST"},
		{"datenoday_07", "Oct 03 21:46m", timelib.TIMELIB_UNSET, 10, 3, true, 21, 46, 0, true, 43200, 0, "M"},
		{"datenoday_08", "Oct\t2003\t20:45", 2003, 10, 1, true, 20, 45, 0, false, 0, 0, ""},
		{"datenoday_09", "Oct\t03\t21:46m", timelib.TIMELIB_UNSET, 10, 3, true, 21, 46, 0, true, 43200, 0, "M"},
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

			// Only check time fields if they're in the input
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

			// Only check timezone fields if timezone is in the input
			if tt.checkTZ {
				if time.Z != tt.expectZ {
					t.Errorf("Expected Z=%d, got %d", tt.expectZ, time.Z)
				}
				if time.Dst != tt.expectDst {
					t.Errorf("Expected Dst=%d, got %d", tt.expectDst, time.Dst)
				}
				if tt.expectTZ != "" && time.TzAbbr != tt.expectTZ {
					t.Errorf("Expected TzAbbr=%s, got %s", tt.expectTZ, time.TzAbbr)
				}
			}
		})
	}
}

// TestParseDateSlashExtended adds the missing dateslash tests from C
// C tests dateslash_00 through dateslash_04
// Existing TestParseDateSlash has 2 tests, missing 3
func TestParseDateSlashExtended(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"dateslash_00", "2005/8/12", 2005, 8, 12},
		{"dateslash_01", "2005/01/02", 2005, 1, 2},
		{"dateslash_02", "2005/01/2", 2005, 1, 2},
		{"dateslash_03", "2005/1/02", 2005, 1, 2},
		{"dateslash_04", "2005/1/2", 2005, 1, 2},
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

// TestParseDateWeekNrExtended adds the missing weeknr tests from C
// C tests weeknr_00 through weeknr_05 (formats like "1995W051" without dash after year)
// Existing TestParseDateWeekNumbers has 5 tests with different format ("2008-W01-1" with dashes)
// FIXED: The bug in timelibDaynrFromWeeknr has been fixed - it now uses the correct DayOfWeek function
func TestParseDateWeekNrExtended(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectY    int64
		expectM    int64
		expectD    int64
		expectRelD int64
	}{
		{"weeknr_00", "1995W051", 1995, 1, 1, 29},
		{"weeknr_01", "2004W30", 2004, 1, 1, 200},
		{"weeknr_02", "1995-W051", 1995, 1, 1, 29},
		{"weeknr_03", "2004-W30", 2004, 1, 1, 200},
		{"weeknr_04", "1995W05-1", 1995, 1, 1, 29},
		{"weeknr_05", "1995-W05-1", 1995, 1, 1, 29},
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
			if time.Relative.D != tt.expectRelD {
				t.Errorf("Expected Relative.D=%d, got %d", tt.expectRelD, time.Relative.D)
			}
		})
	}
}

// TestParseDatePointedExtended adds the missing pointeddate test from C
// C tests pointeddate_00 through pointeddate_03
// Existing TestParseDatePointed has 3 tests, missing pointeddate_02 (22.12.78)
func TestParseDatePointedExtended(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"pointeddate_00", "22.12.1978", 1978, 12, 22},
		{"pointeddate_01", "22.7.1978", 1978, 7, 22},
		{"pointeddate_02", "22.12.78", 1978, 12, 22}, // This was missing
		{"pointeddate_03", "22.7.78", 1978, 7, 22},
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
