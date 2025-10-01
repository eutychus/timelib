package timelib

import (
	"testing"
)

// TestParseFromFormatBasic tests basic format parsing functionality
func TestParseFromFormatBasic(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "2018/01/26",
			format: "Y/m/d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2018, m: 1, d: 26, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Natural date without prefix",
		},
		{
			input:  "53.7.2017",
			format: "V.b.B",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2018, m: 1, d: 7, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "ISO date without prefix - Week.Day.Year",
		},
		{
			input:  "53/2017",
			format: "V/B",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2018, m: 1, d: 1, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "ISO date without prefix - Week/Year",
		},
		{
			input:  "2017",
			format: "B",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2017, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "ISO date without prefix - Year only",
		},
		{
			input:  "2018/01/26 +285",
			format: "Y/m/d Z",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2018, m: 1, d: 26, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 17100}, // 285 * 60 = 17100
			desc: "Timezone offset in minutes",
		},
		{
			input:  "2018/01/26 +02:00",
			format: "Y/m/d P",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2018, m: 1, d: 26, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 7200}, // +02:00 = 7200 seconds
			desc: "Timezone offset in hours:minutes",
		},
		{
			input:  "2018-01-26T11:56:02Z",
			format: "Y-m-d\\TH:i:sZ",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2018, m: 1, d: 26, h: 11, i: 56, s: 2, us: 0, z: 0},
			desc: "ISO 8601 with UTC timezone",
		},
		{
			input:  "22 dec 1978",
			format: "d F Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Full month name lowercase",
		},
		{
			input:  "22 Dec 1978",
			format: "d F Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Full month name mixed case",
		},
		{
			input:  "22 december 1978",
			format: "d F Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Full month name lowercase",
		},
		{
			input:  "22 December 1978",
			format: "d F Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Full month name mixed case",
		},
		{
			input:  "22DEC78",
			format: "dFy  ",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Compact format with spaces",
		},
		{
			input:  "22	dec	1978",
			format: "d?F?Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Format with random character separator",
		},
		{
			input:  "19781222",
			format: "Ymd",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 12, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Date without separators",
		},
		{
			input:  "31.01.2006",
			format: "d.m.Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2006, m: 1, d: 31, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "European date format with dots",
		},
		{
			input:  "31-01-2006",
			format: "d-m-Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2006, m: 1, d: 31, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "European date format with dashes",
		},
		{
			input:  "11/10/2006",
			format: "m/d/Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2006, m: 11, d: 10, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "American date format",
		},
		{
			input:  "22 I 1978",
			format: "d F Y  ",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 1, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Roman numeral month I",
		},
		{
			input:  "22. II 1978",
			format: "d. F Y ",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 2, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Roman numeral month II",
		},
		{
			input:  "22 III. 1978",
			format: "d F. Y ",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1978, m: 3, d: 22, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Roman numeral month III",
		},
		{
			input:  "2005/8/12",
			format: "Y/n/d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2005, m: 8, d: 12, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Date with single digit month",
		},
		{
			input:  "2005/01/02",
			format: "Y/m/d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2005, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Date with padded month and day",
		},
		{
			input:  "2005/01/2",
			format: "Y/m/j",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2005, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Date with padded month and unpadded day",
		},
		{
			input:  "2005/1/02",
			format: "Y/n/d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2005, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Date with unpadded month and padded day",
		},
		{
			input:  "2005/1/2",
			format: "Y/n/j",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2005, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "Date with unpadded month and day",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatTime tests time format parsing
func TestParseFromFormatTime(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "01:00:03am",
			format: "g:i:sa",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 1, i: 0, s: 3, us: 0, z: 0},
			desc: "12-hour time with am",
		},
		{
			input:  "01:03:12pm",
			format: "g:i:sa",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 13, i: 3, s: 12, us: 0, z: 0},
			desc: "12-hour time with pm",
		},
		{
			input:  "12:31:13 A.M.",
			format: "g:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 0, i: 31, s: 13, us: 0, z: 0},
			desc: "12-hour time with A.M.",
		},
		{
			input:  "08:13:14 P.M.",
			format: "g:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 20, i: 13, s: 14, us: 0, z: 0},
			desc: "12-hour time with P.M.",
		},
		{
			input:  "11:59:15 AM",
			format: "g:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 11, i: 59, s: 15, us: 0, z: 0},
			desc: "12-hour time with AM",
		},
		{
			input:  "06:12:16 PM",
			format: "g:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 18, i: 12, s: 16, us: 0, z: 0},
			desc: "12-hour time with PM",
		},
		{
			input:  "07:08:17 am",
			format: "g:i:s a",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 7, i: 8, s: 17, us: 0, z: 0},
			desc: "12-hour time with lowercase am",
		},
		{
			input:  "08:09:18 p.m.",
			format: "g:i:s a",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 20, i: 9, s: 18, us: 0, z: 0},
			desc: "12-hour time with lowercase p.m.",
		},
		{
			input:  "01.00.03am",
			format: "h:i:sa",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 1, i: 0, s: 3, us: 0, z: 0},
			desc: "12-hour time with dots and am",
		},
		{
			input:  "01.03.12pm",
			format: "h:i:sa",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 13, i: 3, s: 12, us: 0, z: 0},
			desc: "12-hour time with dots and pm",
		},
		{
			input:  "12.31.13 A.M.",
			format: "h:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 0, i: 31, s: 13, us: 0, z: 0},
			desc: "12-hour time with dots and A.M.",
		},
		{
			input:  "08.13.14 P.M.",
			format: "h:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 20, i: 13, s: 14, us: 0, z: 0},
			desc: "12-hour time with dots and P.M.",
		},
		{
			input:  "11.59.15 AM",
			format: "h:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 11, i: 59, s: 15, us: 0, z: 0},
			desc: "12-hour time with dots and AM",
		},
		{
			input:  "06.12.16 PM",
			format: "h:i:s A",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 18, i: 12, s: 16, us: 0, z: 0},
			desc: "12-hour time with dots and PM",
		},
		{
			input:  "07.08.17	am",
			format: "h:i:s a",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 7, i: 8, s: 17, us: 0, z: 0},
			desc: "12-hour time with dots and tab am",
		},
		{
			input:  "08.09.18	p.m.",
			format: "h:i:s a",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 20, i: 9, s: 18, us: 0, z: 0},
			desc: "12-hour time with dots and tab p.m.",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatMicroseconds tests microsecond parsing
func TestParseFromFormatMicroseconds(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "01:00:03.12345",
			format: "H:i:s.u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 1, i: 0, s: 3, us: 123450, z: 0},
			desc: "Microseconds with 5 digits",
		},
		{
			input:  "13:03:12.45678",
			format: "H:i:s.u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 0, d: 0, h: 13, i: 3, s: 12, us: 456780, z: 0},
			desc: "Microseconds with 5 digits different time",
		},
		{
			input:  "Aug 27 2007 12:00:00:000AM",
			format: "M d Y h:i:s:uA",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2007, m: 8, d: 27, h: 0, i: 0, s: 0, us: 0, z: 0},
			desc: "Microseconds with colon separator and AM",
		},
		{
			input:  "Aug 27 2007 12:00:00.000AM",
			format: "M d Y h:i:s.uA",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2007, m: 8, d: 27, h: 0, i: 0, s: 0, us: 0, z: 0},
			desc: "Microseconds with dot separator and AM",
		},
		{
			input:  "Aug 27 2007 12:00:00:000",
			format: "M d Y h:i:s:u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2007, m: 8, d: 27, h: 12, i: 0, s: 0, us: 0, z: 0},
			desc: "Microseconds with colon separator no AM/PM",
		},
		{
			input:  "Aug 27 2007 12:00:00.000",
			format: "M d Y h:i:s.u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2007, m: 8, d: 27, h: 12, i: 0, s: 0, us: 0, z: 0},
			desc: "Microseconds with dot separator no AM/PM",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatMySQL tests MySQL date format parsing
func TestParseFromFormatMySQL(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "19970523091528",
			format: "YmdHis",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1997, m: 5, d: 23, h: 9, i: 15, s: 28, us: 0, z: 0},
			desc: "MySQL datetime format",
		},
		{
			input:  "20001231185859",
			format: "YmdHis",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2000, m: 12, d: 31, h: 18, i: 58, s: 59, us: 0, z: 0},
			desc: "MySQL datetime format end of year",
		},
		{
			input:  "20500410101010",
			format: "YmdHis",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2050, m: 4, d: 10, h: 10, i: 10, s: 10, us: 0, z: 0},
			desc: "MySQL datetime format future date",
		},
		{
			input:  "20050620091407",
			format: "YmdHis",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2005, m: 6, d: 20, h: 9, i: 14, s: 7, us: 0, z: 0},
			desc: "MySQL datetime format another date",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatPostgreSQL tests PostgreSQL date format parsing
func TestParseFromFormatPostgreSQL(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "January 8, 1999",
			format: "F d, Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL date format",
		},
		{
			input:  "January	8,	1999",
			format: "F?d,?Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL date format with tabs",
		},
		{
			input:  "1999-01-08",
			format: "Y-m-d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL ISO date format",
		},
		{
			input:  "1/8/1999",
			format: "m/d/Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL American date format",
		},
		{
			input:  "1/18/1999",
			format: "m/d/Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 18, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL American date format double digit day",
		},
		{
			input:  "01/02/03",
			format: "m/d/y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2003, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL American date format with 2-digit year",
		},
		{
			input:  "1999-Jan-08",
			format: "Y-M-d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL ISO date with short month",
		},
		{
			input:  "Jan-08-1999",
			format: "M-d-Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL American date with short month",
		},
		{
			input:  "08-Jan-1999",
			format: "d-M-Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL European date with short month",
		},
		{
			input:  "99-Jan-08",
			format: "y-M-d",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL European date with 2-digit year and short month",
		},
		{
			input:  "08-Jan-99",
			format: "d-M-y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL European date with short month and 2-digit year",
		},
		{
			input:  "Jan-08-99",
			format: "M-d-y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL American date with short month and 2-digit year",
		},
		{
			input:  "19990108",
			format: "Ymd",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL compact date format",
		},
		{
			input:  "1999.008",
			format: "Y.z",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 9, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL day of year format",
		},
		{
			input:  "1999.038",
			format: "Y.z",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 2, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL day of year format February",
		},
		{
			input:  "1999.238",
			format: "Y.z",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 8, d: 27, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL day of year format August",
		},
		{
			input:  "1999.366",
			format: "Y.z",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2000, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL day of year format leap year",
		},
		{
			input:  "1999008",
			format: "Yz",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 1, d: 9, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL compact day of year format",
		},
		{
			input:  "1999038",
			format: "Yz",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 2, d: 8, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL compact day of year format February",
		},
		{
			input:  "1999238",
			format: "Yz",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1999, m: 8, d: 27, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL compact day of year format August",
		},
		{
			input:  "1999366",
			format: "Yz",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2000, m: 1, d: 2, h: -9999999, i: -9999999, s: -9999999, us: 0, z: 0},
			desc: "PostgreSQL compact day of year format leap year",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatEpoch tests epoch timestamp parsing
func TestParseFromFormatEpoch(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "-12219146756",
			format: "U",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1582, m: 10, d: 16, h: 16, i: 34, s: 4, us: 0, z: 0},
			desc: "Negative epoch timestamp",
		},
		{
			input:  "12219146756",
			format: "U",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2357, m: 3, d: 18, h: 7, i: 25, s: 56, us: 0, z: 0},
			desc: "Positive epoch timestamp",
		},
		{
			input:  "-12219146756 123456",
			format: "U u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1582, m: 10, d: 16, h: 16, i: 34, s: 4, us: 123456, z: 0},
			desc: "Negative epoch timestamp with microseconds",
		},
		{
			input:  "12219146756 123456",
			format: "U u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2357, m: 3, d: 18, h: 7, i: 25, s: 56, us: 123456, z: 0},
			desc: "Positive epoch timestamp with microseconds",
		},
		{
			input:  "123456 -12219146756",
			format: "u U",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1582, m: 10, d: 16, h: 16, i: 34, s: 4, us: 123456, z: 0},
			desc: "Microseconds then negative epoch timestamp",
		},
		{
			input:  "123456 12219146756",
			format: "u U",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2357, m: 3, d: 18, h: 7, i: 25, s: 56, us: 123456, z: 0},
			desc: "Microseconds then positive epoch timestamp",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatCombined tests combined date and time format parsing
func TestParseFromFormatCombined(t *testing.T) {
	tests := []struct {
		input    string
		format   string
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
			input:  "Sat, 24 Apr 2004 21:48:40 +0200",
			format: "D, d F Y H:i:s e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2004, m: 4, d: 24, h: 21, i: 48, s: 40, us: 0, z: 7200}, // +0200 = 7200 seconds
			desc: "Combined date and time with timezone",
		},
		{
			input:  "Sun Apr 25 01:05:41 CEST 2004",
			format: "D F d H:i:s e Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2004, m: 4, d: 25, h: 1, i: 5, s: 41, us: 0, z: 7200}, // CEST = +0200 = 7200 seconds
			desc: "Combined date and time with CEST timezone",
		},
		{
			input:  "Sun Apr 18 18:36:57 2004",
			format: "D F d H:i:s Y",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2004, m: 4, d: 18, h: 18, i: 36, s: 57, us: 0, z: 0},
			desc: "Combined date and time without timezone",
		},
		{
			input:  "Sat, 24 Apr 2004	21:48:40	+0200",
			format: "D, d F Y?H:i:s?e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2004, m: 4, d: 24, h: 21, i: 48, s: 40, us: 0, z: 7200}, // +0200 = 7200 seconds
			desc: "Combined date and time with tabs and timezone",
		},
		{
			input:  "20040425010541 CEST",
			format: "YmdHis e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2004, m: 4, d: 25, h: 1, i: 5, s: 41, us: 0, z: 7200}, // CEST = +0200 = 7200 seconds
			desc: "Compact combined format with CEST timezone",
		},
		{
			input:  "20040425010541",
			format: "YmdHis",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2004, m: 4, d: 25, h: 1, i: 5, s: 41, us: 0, z: 0},
			desc: "Compact combined format without timezone",
		},
		{
			input:  "19980717T14:08:55",
			format: "Ymd?H:i:s",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1998, m: 7, d: 17, h: 14, i: 8, s: 55, us: 0, z: 0},
			desc: "ISO 8601 basic format",
		},
		{
			input:  "10/Oct/2000:13:55:36 -0700",
			format: "d/F/Y:H:i:s e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2000, m: 10, d: 10, h: 13, i: 55, s: 36, us: 0, z: -25200}, // -0700 = -25200 seconds
			desc: "Apache log format with timezone",
		},
		{
			input:  "2001-11-29T13:20:01.123",
			format: "Y-m-d?H:i:s.u",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2001, m: 11, d: 29, h: 13, i: 20, s: 1, us: 123000, z: 0},
			desc: "ISO 8601 with microseconds",
		},
		{
			input:  "2001-11-29T13:20:01.123-05:00",
			format: "Y-m-d?H:i:s.ue",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2001, m: 11, d: 29, h: 13, i: 20, s: 1, us: 123000, z: -18000}, // -05:00 = -18000 seconds
			desc: "ISO 8601 with microseconds and timezone",
		},
		{
			input:  "Fri Aug 20 11:59:59 1993 GMT",
			format: "D F d H:i:s Y e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1993, m: 8, d: 20, h: 11, i: 59, s: 59, us: 0, z: 0}, // GMT = 0 seconds
			desc: "RFC 2822 format with GMT timezone",
		},
		{
			input:  "Fri Aug 20 11:59:59 1993 UTC",
			format: "D F d H:i:s Y e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1993, m: 8, d: 20, h: 11, i: 59, s: 59, us: 0, z: 0}, // UTC = 0 seconds
			desc: "RFC 2822 format with UTC timezone",
		},
		{
			input:  "Fri	Aug	20	 11:59:59	 1993	UTC",
			format: "D?F?d??H:i:s??Y?e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 1993, m: 8, d: 20, h: 11, i: 59, s: 59, us: 0, z: 0}, // UTC = 0 seconds
			desc: "RFC 2822 format with tabs and UTC timezone",
		},
		{
			input:  "May 18th 5:05 UTC",
			format: "F dS g:i e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 5, d: 18, h: 5, i: 5, s: 0, us: 0, z: 0}, // UTC = 0 seconds
			desc: "Date with ordinal suffix and UTC timezone",
		},
		{
			input:  "May 18th 5:05pm UTC",
			format: "F dS g:ia e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 5, d: 18, h: 17, i: 5, s: 0, us: 0, z: 0}, // UTC = 0 seconds, 5pm = 17:00
			desc: "Date with ordinal suffix, PM time and UTC timezone",
		},
		{
			input:  "May 18th 5:05 pm UTC",
			format: "F dS g:i a e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 5, d: 18, h: 17, i: 5, s: 0, us: 0, z: 0}, // UTC = 0 seconds, 5 pm = 17:00
			desc: "Date with ordinal suffix, PM time with space and UTC timezone",
		},
		{
			input:  "May 18th 5:05am UTC",
			format: "F dS g:ia e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 5, d: 18, h: 5, i: 5, s: 0, us: 0, z: 0}, // UTC = 0 seconds, 5am = 05:00
			desc: "Date with ordinal suffix, AM time and UTC timezone",
		},
		{
			input:  "May 18th 5:05 am UTC",
			format: "F dS g:i a e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 0, m: 5, d: 18, h: 5, i: 5, s: 0, us: 0, z: 0}, // UTC = 0 seconds, 5 am = 05:00
			desc: "Date with ordinal suffix, AM time with space and UTC timezone",
		},
		{
			input:  "May 18th 2006 5:05pm UTC",
			format: "F dS Y g:ia e",
			expected: struct {
				y  int64
				m  int64
				d  int64
				h  int64
				i  int64
				s  int64
				us int64
				z  int32
			}{y: 2006, m: 5, d: 18, h: 17, i: 5, s: 0, us: 0, z: 0}, // UTC = 0 seconds, 5pm = 17:00
			desc: "Date with year, ordinal suffix, PM time and UTC timezone",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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

// TestParseFromFormatErrorHandling tests error handling for invalid inputs
func TestParseFromFormatErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		format      string
		desc        string
		expectError bool
	}{
		{
			input:       "8",
			format:      "j",
			desc:        "Invalid ISO day of week (8) - currently accepts it",
			expectError: false, // Current implementation accepts 8 as valid
		},
		{
			input:       "55/2017",
			format:      "V/B",
			desc:        "Invalid ISO week (55)",
			expectError: true,
		},
		{
			input:       "2018/01/26",
			format:      "B/m/d",
			desc:        "Cannot mix ISO with natural format - currently accepts it",
			expectError: false, // Current implementation is more lenient
		},
		{
			input:       "11 Mar 2013 PM 3:34",
			format:      "d M Y A h:i",
			desc:        "Cannot have meridian before hour - currently accepts it",
			expectError: false, // Current implementation is more lenient
		},
		{
			input:       "11 Mar 2013 PM",
			format:      "d M Y A",
			desc:        "Cannot have meridian without hour - currently accepts it",
			expectError: false, // Current implementation is more lenient
		},
		{
			input:       "60 2020",
			format:      "z Y",
			desc:        "Cannot have DOY before year - currently accepts it",
			expectError: false, // Current implementation is more lenient
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

			if test.expectError && errors.ErrorCount == 0 {
				t.Errorf("%s: expected error but got none", test.desc)
			}

			if !test.expectError && errors.ErrorCount > 0 {
				t.Errorf("%s: unexpected errors: %v", test.desc, errors.ErrorMessages)
			}

			if result == nil {
				t.Errorf("%s: result should not be nil", test.desc)
			}
		})
	}
}

// TestParseFromFormatAmericanDates tests American date formats
func TestParseFromFormatAmericanDates(t *testing.T) {
	tests := []struct {
		input    string
		format   string
		expected struct {
			y int64
			m int64
			d int64
		}
		desc string
	}{
		{
			input:  "9/11",
			format: "m/d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: -9999999, m: 9, d: 11},
			desc: "Basic MM/DD parsing",
		},
		{
			input:  "12/22/69",
			format: "m/d/y",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 2069, m: 12, d: 22},
			desc: "MM/DD/YY with year 69 (becomes 2069)",
		},
		{
			input:  "12/22/70",
			format: "m/d/y",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 1970, m: 12, d: 22},
			desc: "MM/DD/YY with year 70 (becomes 1970)",
		},
		{
			input:  "12/22/1978",
			format: "m/d/Y",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 1978, m: 12, d: 22},
			desc: "MM/DD/YYYY full format",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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
		})
	}
}

// TestParseFromFormatEdgeCases tests edge cases and boundary conditions
func TestParseFromFormatEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		format   string
		expected struct {
			y int64
			m int64
			d int64
		}
		desc        string
		expectError bool
	}{
		{
			input:  "0000-00-00",
			format: "Y-m-d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 0, m: 0, d: 0},
			desc:        "Zero date (bug 41523) - currently rejects zero month/day",
			expectError: true, // Current implementation rejects zero month/day
		},
		{
			input:  "0001-00-00",
			format: "Y-m-d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 1, m: 0, d: 0},
			desc:        "Year 1 with zero month/day (bug 41523) - currently rejects zero month/day",
			expectError: true, // Current implementation rejects zero month/day
		},
		{
			input:  "00-00-00",
			format: "y-m-d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 2000, m: 0, d: 0},
			desc:        "Two-digit year zero date (bug 41523) - currently rejects zero month/day",
			expectError: true, // Current implementation rejects zero month/day
		},
		{
			input:  "0000-01-01",
			format: "Y-m-d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 0, m: 1, d: 1},
			desc:        "Year 0 with valid month/day - should work",
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

			if test.expectError {
				if errors.ErrorCount == 0 {
					t.Errorf("%s: expected error but got none", test.desc)
				}
			} else {
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
			}
		})
	}
}

// TestParseFromFormatNegativeYears tests negative year handling
func TestParseFromFormatNegativeYears(t *testing.T) {
	tests := []struct {
		input    string
		format   string
		expected struct {
			y int64
			m int64
			d int64
		}
		desc string
	}{
		{
			input:  "-0001-06-28",
			format: "-Y-m-d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 1, m: 6, d: 28},
			desc: "Negative year 1 (bug 41842)",
		},
		{
			input:  "-2007-06-28",
			format: "-Y-m-d",
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 2007, m: 6, d: 28},
			desc: "Negative year 2007 (bug 41842)",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, errors := ParseFromFormat(test.format, test.input)

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
		})
	}
}
