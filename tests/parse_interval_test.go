package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestParseInterval(t *testing.T) {
	// Test parsing ISO 8601 interval strings
	// This is equivalent to timelib_strtointerval in C

	testCases := []struct {
		input       string
		description string
	}{
		{"P1Y2M3DT4H5M6S", "Basic duration"},
		{"P1W", "Week duration"},
		{"P1Y2M", "Year and month duration"},
		{"PT4H5M", "Time duration"},
		{"2021-01-01T00:00:00Z/2021-12-31T23:59:59Z", "Start and end datetime"},
		{"2021-01-01T00:00:00Z/P1Y", "Start datetime with duration"},
		{"P1Y/2021-12-31T23:59:59Z", "Duration with end datetime"},
		{"R5/2021-01-01T00:00:00Z/P1M", "Recurring interval"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Use Strtointerval which is the Go equivalent of timelib_strtointerval
			errors := &timelib.ErrorContainer{}
			begin, end, period, recurrences, err := timelib.Strtointerval(tc.input, errors)

			// Basic validation - function should not crash
			t.Logf("Input: %s", tc.input)
			t.Logf("  Begin: %+v", begin)
			t.Logf("  End: %+v", end)
			t.Logf("  Period: %+v", period)
			t.Logf("  Recurrences: %d", recurrences)
			t.Logf("  Error: %v", err)

			if errors != nil && errors.ErrorCount > 0 {
				t.Logf("  Parsing errors: %d", errors.ErrorCount)
			}
		})
	}
}

func TestParseIntervalBasic(t *testing.T) {
	// Test basic interval parsing functionality
	input := "P1Y2M3DT4H5M6S"
	errors := &timelib.ErrorContainer{}

	begin, end, period, recurrences, err := timelib.Strtointerval(input, errors)

	// For now, just verify the function can be called without panicking
	t.Logf("Input: %s", input)
	t.Logf("  Begin: %+v", begin)
	t.Logf("  End: %+v", end)
	t.Logf("  Period: %+v", period)
	t.Logf("  Recurrences: %d", recurrences)
	t.Logf("  Error: %v", err)

	if errors != nil && errors.ErrorCount > 0 {
		t.Logf("  Parsing errors: %d", errors.ErrorCount)
	}

	t.Log("ParseIntervalBasic test completed successfully")
}

func TestParseIntervalErrorHandling(t *testing.T) {
	// Test error handling for interval parsing
	testCases := []struct {
		input       string
		description string
	}{
		{"", "Empty string"},
		{"invalid", "Invalid format"},
		{"P", "Incomplete duration"},
		{"2021-01-01", "Date only (no interval)"},
		{"P1X", "Invalid duration unit"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			errors := &timelib.ErrorContainer{}
			begin, end, period, recurrences, err := timelib.Strtointerval(tc.input, errors)

			// The function should handle errors gracefully
			t.Logf("Input: '%s'", tc.input)
			t.Logf("  Begin: %+v", begin)
			t.Logf("  End: %+v", end)
			t.Logf("  Period: %+v", period)
			t.Logf("  Recurrences: %d", recurrences)
			t.Logf("  Error: %v", err)
			t.Logf("  Errors: %d warnings, %d errors", errors.WarningCount, errors.ErrorCount)
		})
	}
}
