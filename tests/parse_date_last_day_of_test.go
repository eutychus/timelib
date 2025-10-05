package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateLastDayOf tests "last X of month/year" expressions
// C tests last_day_of_00 through last_day_of_05
// These test the special TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH type

func TestParseDateLastDayOf(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectY               int64
		expectM               int64
		expectD               int64
		expectRelY            int64
		expectRelM            int64
		expectRelD            int64
		expectWeekday         int
		expectWeekdayBehavior int
		expectSpecialType     int
	}{
		{
			name:                  "last_day_of_00",
			input:                 "last saturday of feb 2008",
			expectY:               2008,
			expectM:               2,
			expectD:               1,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         6,
			expectWeekdayBehavior: 0,
			expectSpecialType:     timelib.TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH,
		},
		{
			name:                  "last_day_of_01",
			input:                 "last tue of 2008-11",
			expectY:               2008,
			expectM:               11,
			expectD:               1,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         2,
			expectWeekdayBehavior: 0,
			expectSpecialType:     timelib.TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH,
		},
		{
			name:                  "last_day_of_02",
			input:                 "last sunday of sept",
			expectY:               timelib.TIMELIB_UNSET,
			expectM:               9,
			expectD:               timelib.TIMELIB_UNSET,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         0,
			expectWeekdayBehavior: 0,
			expectSpecialType:     timelib.TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH,
		},
		{
			name:                  "last_day_of_03",
			input:                 "last saturday of this month",
			expectY:               timelib.TIMELIB_UNSET,
			expectM:               timelib.TIMELIB_UNSET,
			expectD:               timelib.TIMELIB_UNSET,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         6,
			expectWeekdayBehavior: 0,
			expectSpecialType:     timelib.TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH,
		},
		{
			name:                  "last_day_of_04",
			input:                 "last thursday of last month",
			expectY:               timelib.TIMELIB_UNSET,
			expectM:               timelib.TIMELIB_UNSET,
			expectD:               timelib.TIMELIB_UNSET,
			expectRelY:            0,
			expectRelM:            -1,
			expectRelD:            -7,
			expectWeekday:         4,
			expectWeekdayBehavior: 0,
			expectSpecialType:     timelib.TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH,
		},
		{
			name:                  "last_day_of_05",
			input:                 "last wed of fourth month",
			expectY:               timelib.TIMELIB_UNSET,
			expectM:               timelib.TIMELIB_UNSET,
			expectD:               timelib.TIMELIB_UNSET,
			expectRelY:            0,
			expectRelM:            4,
			expectRelD:            -7,
			expectWeekday:         3,
			expectWeekdayBehavior: 0,
			expectSpecialType:     timelib.TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH,
		},
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
			if time.Relative.Y != tt.expectRelY {
				t.Errorf("Expected Relative.Y=%d, got %d", tt.expectRelY, time.Relative.Y)
			}
			if time.Relative.M != tt.expectRelM {
				t.Errorf("Expected Relative.M=%d, got %d", tt.expectRelM, time.Relative.M)
			}
			if time.Relative.D != tt.expectRelD {
				t.Errorf("Expected Relative.D=%d, got %d", tt.expectRelD, time.Relative.D)
			}
			if time.Relative.Weekday != tt.expectWeekday {
				t.Errorf("Expected Relative.Weekday=%d, got %d", tt.expectWeekday, time.Relative.Weekday)
			}
			if time.Relative.WeekdayBehavior != tt.expectWeekdayBehavior {
				t.Errorf("Expected Relative.WeekdayBehavior=%d, got %d", tt.expectWeekdayBehavior, time.Relative.WeekdayBehavior)
			}
			if time.Relative.Special.Type != tt.expectSpecialType {
				t.Errorf("Expected Relative.Special.Type=%d, got %d", tt.expectSpecialType, time.Relative.Special.Type)
			}
		})
	}
}
