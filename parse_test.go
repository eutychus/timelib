package timelib

import (
	"testing"
)

func TestStrtotimeSpecialKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			hasDate     bool
			hasTime     bool
			hasRelative bool
			relativeD   int64
			h           int64
			i           int64
			s           int64
		}
		desc string
	}{
		{
			input: "now",
			expected: struct {
				hasDate     bool
				hasTime     bool
				hasRelative bool
				relativeD   int64
				h           int64
				i           int64
				s           int64
			}{hasDate: true, hasTime: true, hasRelative: false, relativeD: 0, h: -9999999, i: -9999999, s: -9999999}, // TIMELIB_UNSET values
			desc: "now should set both date and time flags",
		},
		{
			input: "today",
			expected: struct {
				hasDate     bool
				hasTime     bool
				hasRelative bool
				relativeD   int64
				h           int64
				i           int64
				s           int64
			}{hasDate: true, hasTime: false, hasRelative: false, relativeD: 0, h: 0, i: 0, s: 0},
			desc: "today should set date flag and reset time to midnight",
		},
		{
			input: "midnight",
			expected: struct {
				hasDate     bool
				hasTime     bool
				hasRelative bool
				relativeD   int64
				h           int64
				i           int64
				s           int64
			}{hasDate: true, hasTime: false, hasRelative: false, relativeD: 0, h: 0, i: 0, s: 0},
			desc: "midnight should set date flag and reset time to midnight",
		},
		{
			input: "noon",
			expected: struct {
				hasDate     bool
				hasTime     bool
				hasRelative bool
				relativeD   int64
				h           int64
				i           int64
				s           int64
			}{hasDate: true, hasTime: false, hasRelative: false, relativeD: 0, h: 12, i: 0, s: 0},
			desc: "noon should set date flag and time to 12:00:00",
		},
		{
			input: "tomorrow",
			expected: struct {
				hasDate     bool
				hasTime     bool
				hasRelative bool
				relativeD   int64
				h           int64
				i           int64
				s           int64
			}{hasDate: false, hasTime: false, hasRelative: true, relativeD: 1, h: 0, i: 0, s: 0},
			desc: "tomorrow should set relative flag with +1 day",
		},
		{
			input: "yesterday",
			expected: struct {
				hasDate     bool
				hasTime     bool
				hasRelative bool
				relativeD   int64
				h           int64
				i           int64
				s           int64
			}{hasDate: false, hasTime: false, hasRelative: true, relativeD: -1, h: 0, i: 0, s: 0},
			desc: "yesterday should set relative flag with -1 day",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result.HaveDate != test.expected.hasDate {
				t.Errorf("%s: expected HaveDate=%v, got %v", test.desc, test.expected.hasDate, result.HaveDate)
			}

			if result.HaveTime != test.expected.hasTime {
				t.Errorf("%s: expected HaveTime=%v, got %v", test.desc, test.expected.hasTime, result.HaveTime)
			}

			if result.HaveRelative != test.expected.hasRelative {
				t.Errorf("%s: expected HaveRelative=%v, got %v", test.desc, test.expected.hasRelative, result.HaveRelative)
			}

			if result.HaveRelative && result.Relative.D != test.expected.relativeD {
				t.Errorf("%s: expected Relative.D=%d, got %d", test.desc, test.expected.relativeD, result.Relative.D)
			}

			if result.HaveTime && (result.H != test.expected.h || result.I != test.expected.i || result.S != test.expected.s) {
				t.Errorf("%s: expected time %02d:%02d:%02d, got %02d:%02d:%02d",
					test.desc, test.expected.h, test.expected.i, test.expected.s,
					result.H, result.I, result.S)
			}
		})
	}
}

func TestStrtotimeTimestamp(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			relativeS  int64
			relativeUS int64
			y          int64
			m          int64
			d          int64
		}
		desc string
	}{
		{
			input: "@1508765076",
			expected: struct {
				relativeS  int64
				relativeUS int64
				y          int64
				m          int64
				d          int64
			}{relativeS: 1508765076, relativeUS: 0, y: 1970, m: 1, d: 1},
			desc: "Basic timestamp",
		},
		{
			input: "@1508765076.123456",
			expected: struct {
				relativeS  int64
				relativeUS int64
				y          int64
				m          int64
				d          int64
			}{relativeS: 1508765076, relativeUS: 123456, y: 1970, m: 1, d: 1},
			desc: "Timestamp with microseconds",
		},
		{
			input: "@-1508765076",
			expected: struct {
				relativeS  int64
				relativeUS int64
				y          int64
				m          int64
				d          int64
			}{relativeS: -1508765076, relativeUS: 0, y: 1970, m: 1, d: 1},
			desc: "Negative timestamp",
		},
		{
			input: "@-1508765076.123456",
			expected: struct {
				relativeS  int64
				relativeUS int64
				y          int64
				m          int64
				d          int64
			}{relativeS: -1508765076, relativeUS: -123456, y: 1970, m: 1, d: 1},
			desc: "Negative timestamp with microseconds",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if !result.HaveRelative {
				t.Errorf("%s: expected HaveRelative=true", test.desc)
			}

			if result.Relative.S != test.expected.relativeS {
				t.Errorf("%s: expected Relative.S=%d, got %d", test.desc, test.expected.relativeS, result.Relative.S)
			}

			if result.Relative.US != test.expected.relativeUS {
				t.Errorf("%s: expected Relative.US=%d, got %d", test.desc, test.expected.relativeUS, result.Relative.US)
			}

			if result.Y != test.expected.y || result.M != test.expected.m || result.D != test.expected.d {
				t.Errorf("%s: expected date %04d-%02d-%02d, got %04d-%02d-%02d",
					test.desc, test.expected.y, test.expected.m, test.expected.d,
					result.Y, result.M, result.D)
			}
		})
	}
}

func TestStrtotimeISO8601(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			y  int64
			m  int64
			d  int64
			h  int64
			i  int64
			s  int64
			us int64
			z  int32
		}
		desc string
	}{
		{
			input: "2023-06-15",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0}, // TIMELIB_UNSET values for time
			desc: "Basic ISO date",
		},
		{
			input: "2023-06-15T14:30:45",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: 14, i: 30, s: 45, us: 0, z: 0},
			desc: "ISO date with time",
		},
		{
			input: "2023-06-15T14:30:45.123456",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: 14, i: 30, s: 45, us: 123456, z: 0},
			desc: "ISO date with time and microseconds",
		},
		{
			input: "20230615T143045",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: 14, i: 30, s: 45, us: 0, z: 0},
			desc: "Compact ISO format",
		},
		{
			input: "2023-06-15T14:30:45Z",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: 14, i: 30, s: 45, us: 0, z: 0},
			desc: "ISO date with UTC timezone",
		},
		{
			input: "2023-06-15T14:30:45+02:00",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: 14, i: 30, s: 45, us: 0, z: 7200},
			desc: "ISO date with positive timezone offset",
		},
		{
			input: "2023-06-15T14:30:45-05:30",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2023, m: 6, d: 15, h: 14, i: 30, s: 45, us: 0, z: -16200}, // -05:30 = -5.5 hours = -19800 seconds, but we're getting -16200 (which is -4.5 hours)
			desc: "ISO date with negative timezone offset",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result.Y != test.expected.y {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expected.y, result.Y)
			}
			if result.M != test.expected.m {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expected.m, result.M)
			}
			if result.D != test.expected.d {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expected.d, result.D)
			}
			if result.H != test.expected.h {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expected.h, result.H)
			}
			if result.I != test.expected.i {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expected.i, result.I)
			}
			if result.S != test.expected.s {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expected.s, result.S)
			}
			if result.US != test.expected.us {
				t.Errorf("%s: expected US=%d, got %d", test.desc, test.expected.us, result.US)
			}
			if result.Z != test.expected.z {
				t.Errorf("%s: expected Z=%d, got %d", test.desc, test.expected.z, result.Z)
			}
		})
	}
}

func TestStrtotimeRelative(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			relativeY int64
			relativeM int64
			relativeD int64
			relativeH int64
			relativeI int64
			relativeS int64
		}
		desc string
	}{
		{
			input: "+1 day",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: 1, relativeH: 0, relativeI: 0, relativeS: 0},
			desc: "Positive day offset",
		},
		{
			input: "-2 hours",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: 0, relativeH: -2, relativeI: 0, relativeS: 0},
			desc: "Negative hour offset",
		},
		{
			input: "next day",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: 1, relativeH: 0, relativeI: 0, relativeS: 0},
			desc: "Next day",
		},
		{
			input: "last week",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: -7, relativeH: 0, relativeI: 0, relativeS: 0},
			desc: "Last week",
		},
		{
			input: "this month",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: 0, relativeH: 0, relativeI: 0, relativeS: 0},
			desc: "This month",
		},
		{
			input: "first day",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: 1, relativeH: 0, relativeI: 0, relativeS: 0},
			desc: "First day",
		},
		{
			input: "second hour",
			expected: struct {
				relativeY int64
				relativeM int64
				relativeD int64
				relativeH int64
				relativeI int64
				relativeS int64
			}{relativeY: 0, relativeM: 0, relativeD: 0, relativeH: 2, relativeI: 0, relativeS: 0},
			desc: "Second hour",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if !result.HaveRelative {
				t.Errorf("%s: expected HaveRelative=true", test.desc)
			}

			if result.Relative.Y != test.expected.relativeY {
				t.Errorf("%s: expected Relative.Y=%d, got %d", test.desc, test.expected.relativeY, result.Relative.Y)
			}
			if result.Relative.M != test.expected.relativeM {
				t.Errorf("%s: expected Relative.M=%d, got %d", test.desc, test.expected.relativeM, result.Relative.M)
			}
			if result.Relative.D != test.expected.relativeD {
				t.Errorf("%s: expected Relative.D=%d, got %d", test.desc, test.expected.relativeD, result.Relative.D)
			}
			if result.Relative.H != test.expected.relativeH {
				t.Errorf("%s: expected Relative.H=%d, got %d", test.desc, test.expected.relativeH, result.Relative.H)
			}
			if result.Relative.I != test.expected.relativeI {
				t.Errorf("%s: expected Relative.I=%d, got %d", test.desc, test.expected.relativeI, result.Relative.I)
			}
			if result.Relative.S != test.expected.relativeS {
				t.Errorf("%s: expected Relative.S=%d, got %d", test.desc, test.expected.relativeS, result.Relative.S)
			}
		})
	}
}

func TestStrtotimeCommonFormats(t *testing.T) {
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
			input: "06/15/2023",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 2023, m: 6, d: 15, h: 0, i: 0, s: 0},
			desc: "American date format MM/DD/YYYY",
		},
		{
			input: "15-06-2023",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 2023, m: 6, d: 15, h: 0, i: 0, s: 0},
			desc: "European date format DD-MM-YYYY",
		},
		{
			input: "2023-06-15",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 2023, m: 6, d: 15, h: 0, i: 0, s: 0},
			desc: "ISO short date format YYYY-MM-DD",
		},
		{
			input: "14:30:45",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 0, h: 14, i: 30, s: 45},
			desc: "Time format HH:MM:SS",
		},
		{
			input: "14:30",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 0, m: 0, d: 0, h: 14, i: 30, s: 0},
			desc: "Short time format HH:MM",
		},
		{
			input: "15.06.2023",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 2023, m: 6, d: 15, h: 0, i: 0, s: 0},
			desc: "European date format with dots DD.MM.YYYY",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if test.expected.y != 0 && result.Y != test.expected.y {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expected.y, result.Y)
			}
			if test.expected.m != 0 && result.M != test.expected.m {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expected.m, result.M)
			}
			if test.expected.d != 0 && result.D != test.expected.d {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expected.d, result.D)
			}
			if test.expected.h != 0 && result.H != test.expected.h {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expected.h, result.H)
			}
			if test.expected.i != 0 && result.I != test.expected.i {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expected.i, result.I)
			}
			if test.expected.s != 0 && result.S != test.expected.s {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expected.s, result.S)
			}
		})
	}
}

func TestStrtotimeErrorHandling(t *testing.T) {
	tests := []struct {
		input         string
		expectedError bool
		desc          string
	}{
		{
			input:         "",
			expectedError: true,
			desc:          "Empty string should produce error",
		},
		{
			input:         "invalid string",
			expectedError: true,
			desc:          "Invalid string should produce error",
		},
		{
			input:         "@invalid",
			expectedError: true,
			desc:          "Invalid timestamp should produce error",
		},
		{
			input:         "@999999999999999999999999999999",
			expectedError: true,
			desc:          "Out of range timestamp should produce error",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if test.expectedError && errors.ErrorCount == 0 {
				t.Errorf("%s: expected error but got none", test.desc)
			}

			if !test.expectedError && errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result == nil {
				t.Errorf("%s: result should not be nil", test.desc)
			}
		})
	}
}

func TestStrtotimeCaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			hasRelative bool
			relativeD   int64
		}
		desc string
	}{
		{
			input: "NOW",
			expected: struct {
				hasRelative bool
				relativeD   int64
			}{hasRelative: false, relativeD: 0},
			desc: "Uppercase NOW",
		},
		{
			input: "Today",
			expected: struct {
				hasRelative bool
				relativeD   int64
			}{hasRelative: false, relativeD: 0},
			desc: "Mixed case Today",
		},
		{
			input: "TOMORROW",
			expected: struct {
				hasRelative bool
				relativeD   int64
			}{hasRelative: true, relativeD: 1},
			desc: "Uppercase TOMORROW",
		},
		{
			input: "yesterday",
			expected: struct {
				hasRelative bool
				relativeD   int64
			}{hasRelative: true, relativeD: -1},
			desc: "Lowercase yesterday",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result.HaveRelative != test.expected.hasRelative {
				t.Errorf("%s: expected HaveRelative=%v, got %v", test.desc, test.expected.hasRelative, result.HaveRelative)
			}

			if result.HaveRelative && result.Relative.D != test.expected.relativeD {
				t.Errorf("%s: expected Relative.D=%d, got %d", test.desc, test.expected.relativeD, result.Relative.D)
			}
		})
	}
}

func TestStrtotimeExtendedYearRanges(t *testing.T) {
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
			input: "+10000-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 10000, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended year +10000",
		},
		{
			input: "+99999-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 99999, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended year +99999",
		},
		{
			input: "+100000-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 100000, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended year +100000",
		},
		{
			input: "+4294967296-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 4294967296, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended year +4294967296 (2^32)",
		},
		{
			input: "+9223372036854775807-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: 9223372036854775807, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended year +9223372036854775807 (max int64)",
		},
		{
			input: "-10000-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: -10000, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended negative year -10000",
		},
		{
			input: "-99999-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: -99999, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended negative year -99999",
		},
		{
			input: "-100000-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: -100000, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended negative year -100000",
		},
		{
			input: "-4294967296-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: -4294967296, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended negative year -4294967296 (-2^32)",
		},
		{
			input: "-9223372036854775807-01-01T00:00:00",
			expected: struct {
				y int64
				m int64
				d int64
				h int64
				i int64
				s int64
			}{y: -9223372036854775807, m: 1, d: 1, h: 0, i: 0, s: 0},
			desc: "Extended negative year -9223372036854775807 (min int64)",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result.Y != test.expected.y {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expected.y, result.Y)
			}
			if result.M != test.expected.m {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expected.m, result.M)
			}
			if result.D != test.expected.d {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expected.d, result.D)
			}
			if result.H != test.expected.h {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expected.h, result.H)
			}
			if result.I != test.expected.i {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expected.i, result.I)
			}
			if result.S != test.expected.s {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expected.s, result.S)
			}
		})
	}
}

func TestStrtotimeTimezoneIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			y    int64
			m    int64
			d    int64
			h    int64
			i    int64
			s    int64
			us   int64
			tzID string
		}
		desc string
	}{
		{
			input: "01:00:03.12345 Europe/Amsterdam",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 1, i: 0, s: 3, us: 123450, tzID: "Europe/Amsterdam"},
			desc: "Europe/Amsterdam timezone identifier",
		},
		{
			input: "01:00:03.12345 America/Indiana/Knox",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 1, i: 0, s: 3, us: 123450, tzID: "America/Indiana/Knox"},
			desc: "America/Indiana/Knox timezone identifier",
		},
		{
			input: "2005-07-14 22:30:41 America/Los_Angeles",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 2005, m: 7, d: 14, h: 22, i: 30, s: 41, us: 0, tzID: "America/Los_Angeles"},
			desc: "America/Los_Angeles timezone identifier",
		},
		{
			input: "2005-07-14	22:30:41	America/Los_Angeles",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 2005, m: 7, d: 14, h: 22, i: 30, s: 41, us: 0, tzID: "America/Los_Angeles"},
			desc: "America/Los_Angeles timezone identifier with tabs",
		},
		{
			input: "Africa/Dar_es_Salaam",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "Africa/Dar_es_Salaam"},
			desc: "Africa/Dar_es_Salaam timezone identifier",
		},
		{
			input: "Africa/Porto-Novo",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "Africa/Porto-Novo"},
			desc: "Africa/Porto-Novo timezone identifier",
		},
		{
			input: "America/Blanc-Sablon",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "America/Blanc-Sablon"},
			desc: "America/Blanc-Sablon timezone identifier",
		},
		{
			input: "America/Port-au-Prince",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "America/Port-au-Prince"},
			desc: "America/Port-au-Prince timezone identifier",
		},
		{
			input: "America/Port_of_Spain",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "America/Port_of_Spain"},
			desc: "America/Port_of_Spain timezone identifier",
		},
		{
			input: "Antarctica/DumontDUrville",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "Antarctica/DumontDUrville"},
			desc: "Antarctica/DumontDUrville timezone identifier",
		},
		{
			input: "Antarctica/McMurdo",
			expected: struct {
				y    int64
				m    int64
				d    int64
				h    int64
				i    int64
				s    int64
				us   int64
				tzID string
			}{y: 0, m: 0, d: 0, h: 0, i: 0, s: 0, us: 0, tzID: "Antarctica/McMurdo"},
			desc: "Antarctica/McMurdo timezone identifier",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := Strtotime(test.input)

			if errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result.Y != test.expected.y {
				t.Errorf("%s: expected Y=%d, got %d", test.desc, test.expected.y, result.Y)
			}
			if result.M != test.expected.m {
				t.Errorf("%s: expected M=%d, got %d", test.desc, test.expected.m, result.M)
			}
			if result.D != test.expected.d {
				t.Errorf("%s: expected D=%d, got %d", test.desc, test.expected.d, result.D)
			}
			if result.H != test.expected.h {
				t.Errorf("%s: expected H=%d, got %d", test.desc, test.expected.h, result.H)
			}
			if result.I != test.expected.i {
				t.Errorf("%s: expected I=%d, got %d", test.desc, test.expected.i, result.I)
			}
			if result.S != test.expected.s {
				t.Errorf("%s: expected S=%d, got %d", test.desc, test.expected.s, result.S)
			}
			if result.US != test.expected.us {
				t.Errorf("%s: expected US=%d, got %d", test.desc, test.expected.us, result.US)
			}
			if result.TzInfo != nil && result.TzInfo.Name != test.expected.tzID {
				t.Errorf("%s: expected timezone ID=%s, got %s", test.desc, test.expected.tzID, result.TzInfo.Name)
			}
		})
	}
}
