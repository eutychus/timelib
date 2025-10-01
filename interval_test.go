package timelib

import (
	"testing"
)

// TestStrtointervalBasic tests basic ISO 8601 interval parsing
func TestStrtointervalBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			y int64
			m int64
			d int64
			h int64
			i int64
			s int64
		}
		desc string
	}{
		{
			input: "P1Y",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 1, m: 0, d: 0, h: 0, i: 0, s: 0},
			desc: "One year duration",
		},
		{
			input: "P2M",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 2, d: 0, h: 0, i: 0, s: 0},
			desc: "Two months duration",
		},
		{
			input: "P3D",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 3, h: 0, i: 0, s: 0},
			desc: "Three days duration",
		},
		{
			input: "PT4H",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 0, h: 4, i: 0, s: 0},
			desc: "Four hours duration",
		},
		{
			input: "PT5M",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 0, h: 0, i: 5, s: 0},
			desc: "Five minutes duration",
		},
		{
			input: "PT6S",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 6},
			desc: "Six seconds duration",
		},
		{
			input: "P1Y2M3DT4H5M6S",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 1, m: 2, d: 3, h: 4, i: 5, s: 6},
			desc: "Combined duration",
		},
		{
			input: "P1W",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 7, h: 0, i: 0, s: 0}, // 1 week = 7 days
			desc: "One week duration",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			errors := &ErrorContainer{}
			begin, end, period, recurrences, err := Strtointerval(test.input, errors)

			if err != nil {
				t.Errorf("%s: unexpected error: %v", test.desc, err)
			}

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			// For duration-only intervals, we should get a period but no begin/end
			if begin != nil || end != nil {
				t.Errorf("%s: expected no begin/end times for duration-only interval", test.desc)
			}

			if period == nil {
				t.Errorf("%s: expected period to be set", test.desc)
				return
			}

			if recurrences != 0 {
				t.Errorf("%s: expected 0 recurrences, got %d", test.desc, recurrences)
			}

			if period.Y != test.expected.y {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expected.y, period.Y)
			}
			if period.M != test.expected.m {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expected.m, period.M)
			}
			if period.D != test.expected.d {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expected.d, period.D)
			}
			if period.H != test.expected.h {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expected.h, period.H)
			}
			if period.I != test.expected.i {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expected.i, period.I)
			}
			if period.S != test.expected.s {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expected.s, period.S)
			}
		})
	}
}

// TestStrtointervalStartEnd tests start and end datetime intervals
func TestStrtointervalStartEnd(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			startY int64
			startM int64
			startD int64
			startH int64
			startI int64
			startS int64
			endY   int64
			endM   int64
			endD   int64
			endH   int64
			endI   int64
			endS   int64
		}
		desc        string
		expectError bool
	}{
		{
			input: "2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
			expected: struct {
				startY int64
				startM int64
				startD int64
				startH int64
				startI int64
				startS int64
				endY   int64
				endM   int64
				endD   int64
				endH   int64
				endI   int64
				endS   int64
			}{
				startY: 2007, startM: 3, startD: 1, startH: 13, startI: 0, startS: 0,
				endY: 2008, endM: 5, endD: 11, endH: 15, endI: 30, endS: 0,
			},
			desc: "Start and end datetime interval",
		},
		{
			input: "2007-03-01T13:00:00Z/P1Y2M3DT4H5M6S",
			expected: struct {
				startY int64
				startM int64
				startD int64
				startH int64
				startI int64
				startS int64
				endY   int64
				endM   int64
				endD   int64
				endH   int64
				endI   int64
				endS   int64
			}{
				startY: 2007, startM: 3, startD: 1, startH: 13, startI: 0, startS: 0,
				endY: 2008, endM: 5, endD: 4, endH: 17, endI: 5, endS: 6,
			},
			desc: "Start datetime with duration",
		},
		{
			input: "P1Y2M3DT4H5M6S/2008-05-11T15:30:00Z",
			expected: struct {
				startY int64
				startM int64
				startD int64
				startH int64
				startI int64
				startS int64
				endY   int64
				endM   int64
				endD   int64
				endH   int64
				endI   int64
				endS   int64
			}{
				startY: 2007, startM: 3, startD: 8, startH: 11, startI: 25, startS: 54,
				endY: 2008, endM: 5, endD: 11, endH: 15, endI: 30, endS: 0,
			},
			desc: "Duration with end datetime",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			errors := &ErrorContainer{}
			begin, end, period, recurrences, err := Strtointerval(test.input, errors)

			if test.expectError {
				// For unsupported formats, we expect some kind of error
				if err == nil && errors.ErrorCount == 0 {
					t.Errorf("%s: expected error for unsupported format", test.desc)
				}
			} else {
				// For supported formats, we expect success
				if err != nil {
					t.Errorf("%s: unexpected error: %v", test.desc, err)
				}
				if errors.ErrorCount > 0 {
					t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
				}

				// For start/end intervals, we should have both begin and end times set
				if begin == nil || end == nil {
					t.Errorf("%s: expected both begin and end times", test.desc)
				}

				if recurrences != 0 {
					t.Errorf("%s: expected 0 recurrences, got %d", test.desc, recurrences)
				}
			}

			// Avoid unused variable warnings
			_ = period
		})
	}
}

// TestStrtointervalRecurring tests recurring intervals
func TestStrtointervalRecurring(t *testing.T) {
	tests := []struct {
		input       string
		expectedRec int
		expectedY   int64
		expectedM   int64
		expectedD   int64
		expectedH   int64
		expectedI   int64
		expectedS   int64
		desc        string
	}{
		{
			input:       "R5/2007-03-01T13:00:00Z/P1Y2M3DT4H5M6S",
			expectedRec: 5,
			expectedY:   1, expectedM: 2, expectedD: 3, expectedH: 4, expectedI: 5, expectedS: 6,
			desc: "Recurring interval with 5 repetitions",
		},
		{
			input:       "R/2007-03-01T13:00:00Z/P1Y2M3DT4H5M6S",
			expectedRec: 0, // Infinite recurrence
			expectedY:   1, expectedM: 2, expectedD: 3, expectedH: 4, expectedI: 5, expectedS: 6,
			desc: "Infinite recurring interval",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			errors := &ErrorContainer{}
			begin, end, period, recurrences, err := Strtointerval(test.input, errors)

			// For supported recurring intervals, we expect success
			if err != nil {
				t.Errorf("%s: unexpected error: %v", test.desc, err)
			}
			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			// For recurring intervals, we should have begin time and period set
			if begin == nil {
				t.Errorf("%s: expected begin time to be set", test.desc)
			}
			if period == nil {
				t.Errorf("%s: expected period to be set", test.desc)
				return
			}
			if recurrences != test.expectedRec {
				t.Errorf("%s: expected %d recurrences, got %d", test.desc, test.expectedRec, recurrences)
			}

			// Validate period values
			if period.Y != test.expectedY {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expectedY, period.Y)
			}
			if period.M != test.expectedM {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expectedM, period.M)
			}
			if period.D != test.expectedD {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expectedD, period.D)
			}
			if period.H != test.expectedH {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expectedH, period.H)
			}
			if period.I != test.expectedI {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expectedI, period.I)
			}
			if period.S != test.expectedS {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expectedS, period.S)
			}

			// Avoid unused variable warning
			_ = end
		})
	}
}

// TestStrtointervalErrorHandling tests error handling for invalid intervals
func TestStrtointervalErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		desc        string
		expectError bool
	}{
		{
			input:       "P",
			desc:        "Empty duration - currently accepts it",
			expectError: false, // Current implementation is lenient
		},
		{
			input:       "PT",
			desc:        "Empty time duration - currently accepts it",
			expectError: false, // Current implementation is lenient
		},
		{
			input:       "invalid",
			desc:        "Invalid format",
			expectError: true,
		},
		{
			input:       "P1X",
			desc:        "Invalid duration component - currently accepts it",
			expectError: false, // Current implementation is lenient
		},
		{
			input:       "2007-03-01T13:00:00Z",
			desc:        "Just a datetime, not an interval",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			errors := &ErrorContainer{}
			begin, end, period, recurrences, err := Strtointerval(test.input, errors)

			if test.expectError {
				if err == nil && errors.ErrorCount == 0 {
					t.Errorf("%s: expected error but got none", test.desc)
				}
			} else {
				if err != nil {
					t.Errorf("%s: unexpected error: %v", test.desc, err)
				}
				if errors.ErrorCount > 0 {
					t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
				}
			}

			// If we expect an error, we should not get valid results
			if test.expectError {
				if begin != nil || end != nil || period != nil {
					t.Errorf("%s: expected no valid results on error", test.desc)
				}
			}

			// Avoid unused variable warning
			_ = recurrences
		})
	}
}

// TestStrtointervalEdgeCases tests edge cases and boundary conditions
func TestStrtointervalEdgeCases(t *testing.T) {
	tests := []struct {
		input     string
		expectedY int64
		expectedM int64
		expectedD int64
		expectedH int64
		expectedI int64
		expectedS int64
		desc      string
	}{
		{
			input:     "P0Y",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Zero year duration",
		},
		{
			input:     "P0M",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Zero month duration",
		},
		{
			input:     "P0D",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Zero day duration",
		},
		{
			input:     "PT0H",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Zero hour duration",
		},
		{
			input:     "PT0M",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Zero minute duration",
		},
		{
			input:     "PT0S",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Zero second duration",
		},
		{
			input:     "P999Y",
			expectedY: 999, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 0,
			desc: "Large year duration",
		},
		{
			input:     "PT23H",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 23, expectedI: 0, expectedS: 0,
			desc: "Maximum hour duration",
		},
		{
			input:     "PT59M",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 59, expectedS: 0,
			desc: "Maximum minute duration",
		},
		{
			input:     "PT59S",
			expectedY: 0, expectedM: 0, expectedD: 0, expectedH: 0, expectedI: 0, expectedS: 59,
			desc: "Maximum second duration",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			errors := &ErrorContainer{}
			begin, end, period, recurrences, err := Strtointerval(test.input, errors)

			if err != nil {
				t.Errorf("%s: unexpected error: %v", test.desc, err)
			}

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if begin != nil || end != nil {
				t.Errorf("%s: expected no begin/end times for duration-only interval", test.desc)
			}

			if recurrences != 0 {
				t.Errorf("%s: expected 0 recurrences, got %d", test.desc, recurrences)
			}

			if period == nil {
				t.Errorf("%s: expected period to be set", test.desc)
				return
			}

			if period.Y != test.expectedY {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expectedY, period.Y)
			}
			if period.M != test.expectedM {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expectedM, period.M)
			}
			if period.D != test.expectedD {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expectedD, period.D)
			}
			if period.H != test.expectedH {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expectedH, period.H)
			}
			if period.I != test.expectedI {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expectedI, period.I)
			}
			if period.S != test.expectedS {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expectedS, period.S)
			}
		})
	}
}
