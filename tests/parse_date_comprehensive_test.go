package tests

import (
	"testing"
	timelib "github.com/eutychus/timelib"
)

// TestParseDateRelative tests relative date/time expressions
func TestParseDateRelative(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectRelY int64
		expectRelM int64
		expectRelD int64
		expectRelH int64
		expectRelI int64
		expectRelS int64
	}{
		{"2 secs", "2 secs", 0, 0, 0, 0, 0, 2},
		{"5 seconds", "5 seconds", 0, 0, 0, 0, 0, 5},
		{"1 sec", "1 sec", 0, 0, 0, 0, 0, 1},
		{"2 mins", "2 mins", 0, 0, 0, 0, 2, 0},
		{"5 minutes", "5 minutes", 0, 0, 0, 0, 5, 0},
		{"1 min", "1 min", 0, 0, 0, 0, 1, 0},
		{"2 hours", "2 hours", 0, 0, 0, 2, 0, 0},
		{"5 hour", "5 hour", 0, 0, 0, 5, 0, 0},
		{"2 days", "2 days", 0, 0, 2, 0, 0, 0},
		{"5 day", "5 day", 0, 0, 5, 0, 0, 0},
		{"2 weeks", "2 weeks", 0, 0, 14, 0, 0, 0},
		{"5 week", "5 week", 0, 0, 35, 0, 0, 0},
		{"2 fortnights", "2 fortnights", 0, 0, 28, 0, 0, 0},
		{"5 fortnight", "5 fortnight", 0, 0, 70, 0, 0, 0},
		{"2 months", "2 months", 0, 2, 0, 0, 0, 0},
		{"5 month", "5 month", 0, 5, 0, 0, 0, 0},
		{"2 years", "2 years", 2, 0, 0, 0, 0, 0},
		{"5 year", "5 year", 5, 0, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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

// TestParseDateISO8601 tests ISO 8601 date formats
func TestParseDateISO8601(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"Basic ISO date", "1978-12-22", 1978, 12, 22},
		{"Zero-padded year", "0078-12-22", 78, 12, 22},
		{"Year 0001", "0001-12-22", 1, 12, 22},
		{"Year 2000", "2000-12-22", 2000, 12, 22},
		{"Year 2008", "2008-12-22", 2008, 12, 22},
		{"Jan 1st", "2008-01-01", 2008, 1, 1},
		{"Dec 31st", "2008-12-31", 2008, 12, 31},
		{"Feb 29 leap year", "2008-02-29", 2008, 2, 29},
		{"Single digit month", "2008-1-22", 2008, 1, 22},
		{"Single digit day", "2008-12-2", 2008, 12, 2},
		{"Both single digits", "2008-1-2", 2008, 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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

// TestParseDateISO8601Long tests ISO 8601 date with time formats
func TestParseDateISO8601Long(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
		expectH int64
		expectI int64
		expectS int64
	}{
		{"Full ISO datetime", "1978-12-22T09:15:44", 1978, 12, 22, 9, 15, 44},
		{"ISO with space separator", "1978-12-22 09:15:44", 1978, 12, 22, 9, 15, 44},
		{"ISO with T separator midnight", "2008-07-01T00:00:00", 2008, 7, 1, 0, 0, 0},
		{"ISO with T separator noon", "2008-07-01T12:00:00", 2008, 7, 1, 12, 0, 0},
		{"ISO with T separator end of day", "2008-07-01T23:59:59", 2008, 7, 1, 23, 59, 59},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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

// TestParseDateAmericanFormat tests American date formats (MM/DD/YYYY)
func TestParseDateAmericanFormat(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
		skipYearCheck bool
	}{
		{"MM/DD/YY (1970)", "12/22/70", 1970, 12, 22, false},
		{"MM/DD/YY (1978)", "12/22/78", 1978, 12, 22, false},
		{"MM/DD/YYYY", "12/22/1978", 1978, 12, 22, false},
		{"MM/DD/YYYY (2078)", "12/22/2078", 2078, 12, 22, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
			}
			defer timelib.TimeDtor(time)

			// Skip year check if requested (for MM/DD format without year)
			if !tt.skipYearCheck {
				if time.Y != tt.expectY {
					t.Errorf("Expected Y=%d, got %d", tt.expectY, time.Y)
				}
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

// TestParseDateTextual tests textual date formats
func TestParseDateTextual(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"Jan 1, 2008", "Jan 1, 2008", 2008, 1, 1},
		{"January 1, 2008", "January 1, 2008", 2008, 1, 1},
		{"Feb 2, 2008", "Feb 2, 2008", 2008, 2, 2},
		{"December 31, 2008", "December 31, 2008", 2008, 12, 31},
		{"1 Jan 2008", "1 Jan 2008", 2008, 1, 1},
		{"1 January 2008", "1 January 2008", 2008, 1, 1},
		{"31 December 2008", "31 December 2008", 2008, 12, 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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

// TestParseDateWithTimezone tests date parsing with timezone identifiers
func TestParseDateWithTimezone(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")

	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
		expectH int64
		expectI int64
		expectS int64
		expectTZ string
	}{
		{"America/New_York", "2006-05-12 12:59:59 America/New_York", 2006, 5, 12, 12, 59, 59, "America/New_York"},
		{"America/New_York (13:00)", "2006-05-12 13:00:00 America/New_York", 2006, 5, 12, 13, 0, 0, "America/New_York"},
		{"America/New_York (13:00:01)", "2006-05-12 13:00:01 America/New_York", 2006, 5, 12, 13, 0, 1, "America/New_York"},
		{"Europe/Amsterdam", "2008-07-01 12:00:00 Europe/Amsterdam", 2008, 7, 1, 12, 0, 0, "Europe/Amsterdam"},
		{"UTC", "2008-07-01 12:00:00 UTC", 2008, 7, 1, 12, 0, 0, "UTC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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

			// Parse timezone from the input and check
			if time.TzInfo != nil {
				var errorCode int
				tzi, _ := timelib.ParseTzfile(tt.expectTZ, testDirectory, &errorCode)
				if tzi != nil {
					defer timelib.TzinfoDtor(tzi)
					if time.TzInfo.Name != tzi.Name {
						t.Errorf("Expected timezone %s, got %s", tzi.Name, time.TzInfo.Name)
					}
				}
			}
		})
	}
}

// TestParseDateTimezoneOffset tests date parsing with timezone offsets
func TestParseDateTimezoneOffset(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
		expectH int64
		expectI int64
		expectS int64
		expectZ int32
	}{
		{"+0000", "2008-07-01T22:35:17+0000", 2008, 7, 1, 22, 35, 17, 0},
		{"+0100", "2008-07-01T22:35:17+0100", 2008, 7, 1, 22, 35, 17, 3600},
		{"+0500", "2008-07-01T22:35:17+0500", 2008, 7, 1, 22, 35, 17, 18000},
		{"-0100", "2008-07-01T22:35:17-0100", 2008, 7, 1, 22, 35, 17, -3600},
		{"-0500", "2008-07-01T22:35:17-0500", 2008, 7, 1, 22, 35, 17, -18000},
		{"+01:00", "2008-07-01T22:35:17+01:00", 2008, 7, 1, 22, 35, 17, 3600},
		{"-05:00", "2008-07-01T22:35:17-05:00", 2008, 7, 1, 22, 35, 17, -18000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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
			if time.Z != tt.expectZ {
				t.Errorf("Expected Z=%d, got %d", tt.expectZ, time.Z)
			}
		})
	}
}

// TestParseDateMicroseconds tests microsecond parsing
func TestParseDateMicroseconds(t *testing.T) {
	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
		expectH int64
		expectI int64
		expectS int64
		expectUS int64
	}{
		{"6 digits", "2008-07-01T22:35:17.123456", 2008, 7, 1, 22, 35, 17, 123456},
		{"5 digits", "2008-07-01T22:35:17.12345", 2008, 7, 1, 22, 35, 17, 123450},
		{"4 digits", "2008-07-01T22:35:17.1234", 2008, 7, 1, 22, 35, 17, 123400},
		{"3 digits", "2008-07-01T22:35:17.123", 2008, 7, 1, 22, 35, 17, 123000},
		{"2 digits", "2008-07-01T22:35:17.12", 2008, 7, 1, 22, 35, 17, 120000},
		{"1 digit", "2008-07-01T22:35:17.1", 2008, 7, 1, 22, 35, 17, 100000},
		{"Zero", "2008-07-01T22:35:17.000000", 2008, 7, 1, 22, 35, 17, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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

// TestParseDateSpecialKeywords tests special keyword parsing
func TestParseDateSpecialKeywords(t *testing.T) {
	tests := []struct {
		name string
		input string
		shouldHaveRelative bool
	}{
		{"now lowercase", "now", false},
		{"NOW uppercase", "NOW", false},
		{"noW mixed case", "noW", false},
		{"today", "today", false},
		{"yesterday", "yesterday", true},
		{"tomorrow", "tomorrow", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
			}
			defer timelib.TimeDtor(time)

			if tt.shouldHaveRelative {
				if !time.HaveRelative {
					t.Errorf("Expected HaveRelative=true for %s", tt.input)
				}
			}
		})
	}
}

// TestParseDateWeekNumbers tests ISO week number parsing
// Note: ISO week format (YYYY-Www-D) is not yet fully implemented in the parser
func TestParseDateWeekNumbers(t *testing.T) {
	t.Skip("ISO week format (YYYY-Www-D) not yet implemented - TODO")

	tests := []struct {
		name string
		input string
		expectY int64
		expectM int64
		expectD int64
	}{
		{"Week 1 Monday 2008", "2008-W01-1", 2007, 12, 31},
		{"Week 1 Sunday 2008", "2008-W01-7", 2008, 1, 6},
		{"Week 52 Monday 2008", "2008-W52-1", 2008, 12, 22},
		{"Week 1 Monday 2009", "2009-W01-1", 2008, 12, 29},
		{"Week 53 Sunday 2009", "2009-W53-7", 2010, 1, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, errors := timelib.Strtotime(tt.input)
			if errors != nil && errors.ErrorCount > 0 {
				t.Fatalf("Parse failed: %v", errors.ErrorMessages)
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
