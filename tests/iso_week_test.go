package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestIsoWeekFromDate(t *testing.T) {
	tests := []struct {
		name    string
		year    int64
		month   int64
		day     int64
		expWeek int64
		expYear int64
	}{
		{"2023-01-01", 2023, 1, 1, 52, 2022},
		{"2023-01-02", 2023, 1, 2, 1, 2023},
		{"2023-12-31", 2023, 12, 31, 52, 2023},
		{"2024-01-01", 2024, 1, 1, 1, 2024},
		{"2024-12-30", 2024, 12, 30, 1, 2025},
		{"2024-12-31", 2024, 12, 31, 1, 2025},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			week, year := timelib.IsoWeekFromDate(tt.year, tt.month, tt.day)
			if week != tt.expWeek || year != tt.expYear {
				t.Errorf("IsoWeekFromDate(%d, %d, %d) = %d, %d; want %d, %d",
					tt.year, tt.month, tt.day, week, year, tt.expWeek, tt.expYear)
			}
		})
	}
}

func TestIsoDateFromDate(t *testing.T) {
	tests := []struct {
		name    string
		year    int64
		month   int64
		day     int64
		expYear int64
		expWeek int64
		expDay  int64
	}{
		{"2023-01-01", 2023, 1, 1, 2022, 52, 7},
		{"2023-01-02", 2023, 1, 2, 2023, 1, 1},
		{"2023-12-31", 2023, 12, 31, 2023, 52, 7},
		{"2024-01-01", 2024, 1, 1, 2024, 1, 1},
		{"2024-12-30", 2024, 12, 30, 2025, 1, 1},
		{"2024-12-31", 2024, 12, 31, 2025, 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, week, day := timelib.IsoDateFromDate(tt.year, tt.month, tt.day)
			if year != tt.expYear || week != tt.expWeek || day != tt.expDay {
				t.Errorf("IsoDateFromDate(%d, %d, %d) = %d, %d, %d; want %d, %d, %d",
					tt.year, tt.month, tt.day, year, week, day, tt.expYear, tt.expWeek, tt.expDay)
			}
		})
	}
}
