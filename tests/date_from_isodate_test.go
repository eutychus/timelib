package tests

import (
	"fmt"
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestDateFromIsoDate(t *testing.T) {
	// Test data from the original C test
	tests := []struct {
		isoYear, isoWeek, isoDay                 int64
		expectedYear, expectedMonth, expectedDay int64
		description                              string
	}{
		{2014, 52, 1, 2014, 12, 22, "2014 week 52 day 1"},
		{2014, 52, 7, 2014, 12, 28, "2014 week 52 day 7"},
		{2015, 1, 1, 2014, 12, 29, "2015 week 1 day 1"},
		{2015, 1, 3, 2014, 12, 31, "2015 week 1 day 3"},
		{2015, 1, 4, 2015, 1, 1, "2015 week 1 day 4"},
		{2015, 52, 7, 2015, 12, 27, "2015 week 52 day 7"},
		{2015, 53, 1, 2015, 12, 28, "2015 week 53 day 1"},
		{2015, 53, 4, 2015, 12, 31, "2015 week 53 day 4"},
		{2015, 53, 5, 2016, 1, 1, "2015 week 53 day 5"},
		{2016, 1, 1, 2016, 1, 4, "2016 week 1 day 1"},
		{2016, 1, 3, 2016, 1, 6, "2016 week 1 day 3"},
		{2016, 1, 4, 2016, 1, 7, "2016 week 1 day 4"},
		{2016, 51, 7, 2016, 12, 25, "2016 week 51 day 7"},
		{2016, 52, 1, 2016, 12, 26, "2016 week 52 day 1"},
		{2016, 52, 4, 2016, 12, 29, "2016 week 52 day 4"},
		{2016, 52, 7, 2017, 1, 1, "2016 week 52 day 7"},
		{2017, 8, 6, 2017, 2, 25, "2017 week 8 day 6"},
		{2017, 8, 7, 2017, 2, 26, "2017 week 8 day 7"},
		{2017, 9, 1, 2017, 2, 27, "2017 week 9 day 1"},
		{2017, 9, 2, 2017, 2, 28, "2017 week 9 day 2"},
		{2017, 9, 3, 2017, 3, 1, "2017 week 9 day 3"},
		{2020, 9, 2, 2020, 2, 25, "2020 week 9 day 2"},
		{2020, 9, 3, 2020, 2, 26, "2020 week 9 day 3"},
		{2020, 9, 5, 2020, 2, 28, "2020 week 9 day 5"},
		{2020, 9, 6, 2020, 2, 29, "2020 week 9 day 6 (leap year)"},
		{2020, 9, 7, 2020, 3, 1, "2020 week 9 day 7"},
		{2043, 53, 1, 2043, 12, 28, "2043 week 53 day 1"},
		{2043, 53, 2, 2043, 12, 29, "2043 week 53 day 2"},
		{2043, 53, 3, 2043, 12, 30, "2043 week 53 day 3"},
		{2043, 53, 4, 2043, 12, 31, "2043 week 53 day 4"},
		{2043, 53, 5, 2044, 1, 1, "2043 week 53 day 5"},
		{2043, 53, 6, 2044, 1, 2, "2043 week 53 day 6"},
		{2043, 53, 7, 2044, 1, 3, "2043 week 53 day 7"},
		{2019, 0, 1, 2018, 12, 24, "2019 week 0 day 1"},
		{2019, -1, 1, 2018, 12, 17, "2019 week -1 day 1"},
		{2019, 62, 1, 2020, 3, 2, "2019 week 62 day 1"},
		{2019, 110, 1, 2021, 2, 1, "2019 week 110 day 1"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			year, month, day := timelib.DateFromIsoDate(test.isoYear, test.isoWeek, test.isoDay)

			expected := fmt.Sprintf("%04d-%02d-%02d", test.expectedYear, test.expectedMonth, test.expectedDay)
			actual := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

			if actual != expected {
				t.Errorf("DateFromIsoDate(%d, %d, %d) = %s, expected %s",
					test.isoYear, test.isoWeek, test.isoDay, actual, expected)
			} else {
				t.Logf("OK: DateFromIsoDate(%d, %d, %d) = %s",
					test.isoYear, test.isoWeek, test.isoDay, actual)
			}
		})
	}
}

func TestDateFromIsoDateBasic(t *testing.T) {
	// Test basic ISO date conversion
	year, month, day := timelib.DateFromIsoDate(2023, 1, 1)

	// Basic validation - should return valid date components
	if year < 1 || year > 9999 {
		t.Errorf("Expected valid year, got %d", year)
	}

	if month < 1 || month > 12 {
		t.Errorf("Expected valid month (1-12), got %d", month)
	}

	if day < 1 || day > 31 {
		t.Errorf("Expected valid day (1-31), got %d", day)
	}

	t.Logf("DateFromIsoDate(2023, 1, 1) = %04d-%02d-%02d", year, month, day)
}
