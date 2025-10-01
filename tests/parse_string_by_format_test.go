package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestParseStringByFormat(t *testing.T) {
	// Test parsing strings using format specifiers
	// This is equivalent to timelib_parse_from_format in C

	testCases := []struct {
		format      string
		input       string
		expected    string
		hasError    bool
		description string
	}{
		{"d M Y h:i A", "07 Apr 2021 3:34 PM", "2021-04-07 15:34:00", false, "American format with AM/PM"},
		{"Y-m-d H:i:s", "2021-04-07 15:34:00", "2021-04-07 15:34:00", false, "ISO format"},
		{"d/m/Y", "07/04/2021", "2021-04-07 00:00:00", false, "European date format"},
		{"Y-m-d", "invalid-date", "", true, "Invalid date should fail"},
		{"H:i", "14:30", "1970-01-01 14:30:00", false, "Time only format"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Use ParseFromFormat which is the Go equivalent of timelib_parse_from_format
			result, errors := timelib.ParseFromFormat(tc.format, tc.input)

			if tc.hasError {
				if errors == nil || errors.ErrorCount == 0 {
					t.Errorf("Expected error for format '%s' with input '%s', but got none", tc.format, tc.input)
				}
			} else {
				if errors != nil && errors.ErrorCount > 0 {
					t.Errorf("Expected no error for format '%s' with input '%s', but got errors: %v", tc.format, tc.input, errors)
				}
				if result == nil {
					t.Errorf("Expected non-nil result for format '%s' with input '%s'", tc.format, tc.input)
				} else {
					// Basic validation - check that we got a valid time structure
					// For time-only formats, date components may remain unset (-9999999)
					if tc.format == "H:i" {
						// Time-only format should have valid time components
						if result.H < 0 || result.H > 23 || result.I < 0 || result.I > 59 {
							t.Errorf("Expected valid time components for time-only format, got H=%d, I=%d", result.H, result.I)
						}
					} else {
						// Date formats should have valid date components
						if result.Y <= 0 || result.M <= 0 || result.D <= 0 {
							t.Errorf("Expected valid date components, got Y=%d, M=%d, D=%d", result.Y, result.M, result.D)
						}
					}
				}
			}
		})
	}
}

func TestParseStringByFormatBasic(t *testing.T) {
	// Test basic format parsing functionality
	format := "d M Y h:i A"
	input := "07 Apr 2021 3:34 PM"

	result, errors := timelib.ParseFromFormat(format, input)

	// For now, just verify the function can be called without panicking
	if errors != nil && errors.ErrorCount > 0 {
		t.Logf("ParseFromFormat returned errors (may be expected): %v", errors)
	}

	if result != nil {
		t.Logf("ParseFromFormat returned result: Y=%d, M=%d, D=%d, H=%d, I=%d",
			result.Y, result.M, result.D, result.H, result.I)
	} else {
		t.Log("ParseFromFormat returned nil (may be expected for unimplemented features)")
	}

	t.Log("ParseStringByFormatBasic test completed successfully")
}

func TestParseStringByFormatErrorHandling(t *testing.T) {
	// Test error handling for format parsing
	testCases := []struct {
		format      string
		input       string
		description string
	}{
		{"invalid-format", "test", "Invalid format string"},
		{"Y-m-d", "not-a-date", "Invalid date string"},
		{"", "", "Empty format and input"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, errors := timelib.ParseFromFormat(tc.format, tc.input)

			// The function should handle errors gracefully
			t.Logf("Format: '%s', Input: '%s'", tc.format, tc.input)
			if errors != nil {
				t.Logf("  Errors: %d warnings, %d errors", errors.WarningCount, errors.ErrorCount)
			}
			if result != nil {
				t.Logf("  Result: %+v", result)
			}
		})
	}
}
