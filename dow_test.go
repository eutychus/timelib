package timelib

import (
	"testing"
)

func TestDayOfWeek(t *testing.T) {
	tests := []struct {
		y, m, d  int64
		expected int64
		desc     string
	}{
		{1978, 12, 22, 5, "December 22, 1978 is Friday (5)"},
		{2005, 2, 19, 6, "February 19, 2005 is Saturday (6)"},
		{2023, 1, 1, 0, "January 1, 2023 is Sunday (0)"},
		{2023, 6, 15, 4, "June 15, 2023 is Thursday (4)"},
		{2000, 2, 29, 2, "February 29, 2000 (leap year) is Tuesday (2)"},
		{1900, 1, 1, 1, "January 1, 1900 is Monday (1)"},
		{1600, 1, 1, 6, "January 1, 1600 is Saturday (6)"},
		{-4713, 11, 24, 1, "November 24, -4713 is Monday (1)"}, // Julian day epoch
	}

	for _, test := range tests {
		result := DayOfWeek(test.y, test.m, test.d)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}

func TestIsoDayOfWeek(t *testing.T) {
	tests := []struct {
		y, m, d  int64
		expected int64
		desc     string
	}{
		{1978, 12, 22, 5, "December 22, 1978 is Friday (5)"},
		{2005, 2, 19, 6, "February 19, 2005 is Saturday (6)"},
		{2023, 1, 1, 7, "January 1, 2023 is Sunday (7)"},
		{2023, 6, 15, 4, "June 15, 2023 is Thursday (4)"},
		{2000, 2, 29, 2, "February 29, 2000 (leap year) is Tuesday (2)"},
		{1900, 1, 1, 1, "January 1, 1900 is Monday (1)"},
		{1600, 1, 1, 6, "January 1, 1600 is Saturday (6)"},
	}

	for _, test := range tests {
		result := IsoDayOfWeek(test.y, test.m, test.d)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}

func TestDayOfYear(t *testing.T) {
	tests := []struct {
		y, m, d  int64
		expected int64
		desc     string
	}{
		{2023, 1, 1, 0, "January 1, 2023 is day 0 of year"},
		{2023, 1, 31, 30, "January 31, 2023 is day 30 of year"},
		{2023, 2, 1, 31, "February 1, 2023 is day 31 of year"},
		{2023, 12, 31, 364, "December 31, 2023 is day 364 of year"},
		{2024, 1, 1, 0, "January 1, 2024 (leap year) is day 0 of year"},
		{2024, 2, 29, 59, "February 29, 2024 (leap year) is day 59 of year"},
		{2024, 12, 31, 365, "December 31, 2024 (leap year) is day 365 of year"},
	}

	for _, test := range tests {
		result := DayOfYear(test.y, test.m, test.d)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		y, m     int64
		expected int64
		desc     string
	}{
		{2023, 1, 31, "January 2023 has 31 days"},
		{2023, 2, 28, "February 2023 has 28 days (non-leap year)"},
		{2023, 4, 30, "April 2023 has 30 days"},
		{2024, 2, 29, "February 2024 has 29 days (leap year)"},
		{1900, 2, 28, "February 1900 has 28 days (not leap year, divisible by 100)"},
		{2000, 2, 29, "February 2000 has 29 days (leap year, divisible by 400)"},
		{2023, 12, 31, "December 2023 has 31 days"},
	}

	for _, test := range tests {
		result := DaysInMonth(test.y, test.m)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}

func TestValidTime(t *testing.T) {
	tests := []struct {
		h, i, s  int64
		expected bool
		desc     string
	}{
		{0, 0, 0, true, "00:00:00 is valid"},
		{23, 59, 59, true, "23:59:59 is valid"},
		{12, 30, 45, true, "12:30:45 is valid"},
		{24, 0, 0, false, "24:00:00 is invalid"},
		{23, 60, 0, false, "23:60:00 is invalid"},
		{23, 0, 60, false, "23:00:60 is invalid"},
		{-1, 0, 0, false, "-1:00:00 is invalid"},
		{0, -1, 0, false, "00:-1:00 is invalid"},
		{0, 0, -1, false, "00:00:-1 is invalid"},
	}

	for _, test := range tests {
		result := ValidTime(test.h, test.i, test.s)
		if result != test.expected {
			t.Errorf("%s: got %v, expected %v", test.desc, result, test.expected)
		}
	}
}

func TestValidDate(t *testing.T) {
	tests := []struct {
		y, m, d  int64
		expected bool
		desc     string
	}{
		{2023, 1, 1, true, "2023-01-01 is valid"},
		{2023, 12, 31, true, "2023-12-31 is valid"},
		{2024, 2, 29, true, "2024-02-29 is valid (leap year)"},
		{2023, 2, 29, false, "2023-02-29 is invalid (non-leap year)"},
		{2023, 4, 31, false, "2023-04-31 is invalid (April has 30 days)"},
		{2023, 0, 1, false, "2023-00-01 is invalid (month 0)"},
		{2023, 13, 1, false, "2023-13-01 is invalid (month 13)"},
		{2023, 1, 0, false, "2023-01-00 is invalid (day 0)"},
		{2023, 1, 32, false, "2023-01-32 is invalid (January has 31 days)"},
	}

	for _, test := range tests {
		result := ValidDate(test.y, test.m, test.d)
		if result != test.expected {
			t.Errorf("%s: got %v, expected %v", test.desc, result, test.expected)
		}
	}
}

func TestIsoWeekFromDate(t *testing.T) {
	tests := []struct {
		y, m, d    int64
		expectedIw int64
		expectedIy int64
		desc       string
	}{
		{2023, 1, 1, 52, 2022, "January 1, 2023 is week 52 of 2022"},
		{2023, 1, 2, 1, 2023, "January 2, 2023 is week 1 of 2023"},
		{2023, 12, 31, 52, 2023, "December 31, 2023 is week 52 of 2023"},
		{2024, 1, 1, 1, 2024, "January 1, 2024 is week 1 of 2024"},
		{2024, 12, 30, 1, 2025, "December 30, 2024 is week 1 of 2025"},
	}

	for _, test := range tests {
		iw, iy := IsoWeekFromDate(test.y, test.m, test.d)
		if iw != test.expectedIw || iy != test.expectedIy {
			t.Errorf("%s: got week %d of year %d, expected week %d of year %d",
				test.desc, iw, iy, test.expectedIw, test.expectedIy)
		}
	}
}

func TestIsoDateFromDate(t *testing.T) {
	tests := []struct {
		y, m, d    int64
		expectedIy int64
		expectedIw int64
		expectedId int64
		desc       string
	}{
		{2023, 1, 1, 2022, 52, 7, "January 1, 2023 is 2022-W52-7"},
		{2023, 1, 2, 2023, 1, 1, "January 2, 2023 is 2023-W01-1"},
		{2023, 6, 15, 2023, 24, 4, "June 15, 2023 is 2023-W24-4"},
	}

	for _, test := range tests {
		iy, iw, id := IsoDateFromDate(test.y, test.m, test.d)
		if iy != test.expectedIy || iw != test.expectedIw || id != test.expectedId {
			t.Errorf("%s: got %d-W%02d-%d, expected %d-W%02d-%d",
				test.desc, iy, iw, id, test.expectedIy, test.expectedIw, test.expectedId)
		}
	}
}

func TestDayNrFromWeekNr(t *testing.T) {
	tests := []struct {
		iy, iw, id int64
		expected   int64
		desc       string
	}{
		{2023, 1, 1, 1, "Week 1, day 1 of 2023 is day 1"},
		{2023, 1, 7, 7, "Week 1, day 7 of 2023 is day 7"},
		{2023, 2, 1, 8, "Week 2, day 1 of 2023 is day 8"},
	}

	for _, test := range tests {
		result := DayNrFromWeekNr(test.iy, test.iw, test.id)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}

func TestDateFromIsoDate(t *testing.T) {
	tests := []struct {
		iy, iw, id int64
		expectedY  int64
		expectedM  int64
		expectedD  int64
		desc       string
	}{
		{2023, 1, 1, 2023, 1, 2, "2023-W01-1 is January 2, 2023"},
		{2023, 1, 7, 2023, 1, 8, "2023-W01-7 is January 8, 2023"},
		{2023, 24, 4, 2023, 6, 15, "2023-W24-4 is June 15, 2023"},
	}

	for _, test := range tests {
		y, m, d := DateFromIsoDate(test.iy, test.iw, test.id)
		if y != test.expectedY || m != test.expectedM || d != test.expectedD {
			t.Errorf("%s: got %04d-%02d-%02d, expected %04d-%02d-%02d",
				test.desc, y, m, d, test.expectedY, test.expectedM, test.expectedD)
		}
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		y        int64
		expected bool
		desc     string
	}{
		{2024, true, "2024 is a leap year"},
		{2023, false, "2023 is not a leap year"},
		{2000, true, "2000 is a leap year (divisible by 400)"},
		{1900, false, "1900 is not a leap year (divisible by 100 but not 400)"},
		{1600, true, "1600 is a leap year (divisible by 400)"},
		{1996, true, "1996 is a leap year"},
		{1997, false, "1997 is not a leap year"},
	}

	for _, test := range tests {
		result := IsLeapYear(test.y)
		if result != test.expected {
			t.Errorf("%s: got %v, expected %v", test.desc, result, test.expected)
		}
	}
}

func TestPositiveMod(t *testing.T) {
	tests := []struct {
		x, y     int64
		expected int64
		desc     string
	}{
		{10, 7, 3, "10 % 7 = 3"},
		{-10, 7, 4, "-10 % 7 = 4 (positive result)"},
		{10, -7, 3, "10 % -7 = 3"},
		{-10, -7, 4, "-10 % -7 = 4 (positive result)"},
		{0, 7, 0, "0 % 7 = 0"},
		{7, 7, 0, "7 % 7 = 0"},
	}

	for _, test := range tests {
		result := positiveMod(test.x, test.y)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}
