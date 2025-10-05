package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateWeek tests "this week", "last week", "next week" expressions
// C tests week_00 through week_31
// These test relative week expressions with weekday_behavior = 2

func TestParseDateWeek(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectH               int64
		expectI               int64
		expectS               int64
		expectRelY            int64
		expectRelM            int64
		expectRelD            int64
		expectWeekday         int
		expectWeekdayBehavior int
		checkTime             bool
		checkRel              bool
	}{
		// this week (week_00-07)
		{
			name:                  "week_00",
			input:                 "this week",
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
		},
		{
			name:                  "week_01",
			input:                 "this week monday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		{
			name:                  "week_02",
			input:                 "this week tuesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         2,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		{
			name:                  "week_03",
			input:                 "this week wednesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         3,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		{
			name:                  "week_04",
			input:                 "thursday this week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         4,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		{
			name:                  "week_05",
			input:                 "friday this week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         5,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		{
			name:                  "week_06",
			input:                 "saturday this week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         6,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		{
			name:                  "week_07",
			input:                 "sunday this week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectWeekday:         0,
			expectWeekdayBehavior: 2,
			checkTime:             true,
		},
		// last week (week_08-15)
		{
			name:                  "week_08",
			input:                 "last week",
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkRel:              true,
		},
		{
			name:                  "week_09",
			input:                 "last week monday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_10",
			input:                 "last week tuesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         2,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_11",
			input:                 "last week wednesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         3,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_12",
			input:                 "thursday last week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         4,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_13",
			input:                 "friday last week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         5,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_14",
			input:                 "saturday last week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         6,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_15",
			input:                 "sunday last week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         0,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		// previous week (week_16-23)
		{
			name:                  "week_16",
			input:                 "previous week",
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkRel:              true,
		},
		{
			name:                  "week_17",
			input:                 "previous week monday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_18",
			input:                 "previous week tuesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         2,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_19",
			input:                 "previous week wednesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         3,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_20",
			input:                 "thursday previous week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         4,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_21",
			input:                 "friday previous week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         5,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_22",
			input:                 "saturday previous week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         6,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_23",
			input:                 "sunday previous week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            -7,
			expectWeekday:         0,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		// next week (week_24-31)
		{
			name:                  "week_24",
			input:                 "next week",
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkRel:              true,
		},
		{
			name:                  "week_25",
			input:                 "next week monday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         1,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_26",
			input:                 "next week tuesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         2,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_27",
			input:                 "next week wednesday",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         3,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_28",
			input:                 "thursday next week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         4,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_29",
			input:                 "friday next week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         5,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_30",
			input:                 "saturday next week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         6,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
		{
			name:                  "week_31",
			input:                 "sunday next week",
			expectH:               0,
			expectI:               0,
			expectS:               0,
			expectRelY:            0,
			expectRelM:            0,
			expectRelD:            7,
			expectWeekday:         0,
			expectWeekdayBehavior: 2,
			checkTime:             true,
			checkRel:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			if tt.checkTime {
				if time.H != tt.expectH {
					t.Errorf("Expected H=%d, got %d", tt.expectH, time.H)
				}
				if time.I != tt.expectI {
					t.Errorf("Expected I=%d, got %d", tt.expectI, time.I)
				}
				if time.S != tt.expectS {
					t.Errorf("Expected S=%d, got %d", tt.expectS, time.S)
				}
			}

			if tt.checkRel {
				if time.Relative.Y != tt.expectRelY {
					t.Errorf("Expected Relative.Y=%d, got %d", tt.expectRelY, time.Relative.Y)
				}
				if time.Relative.M != tt.expectRelM {
					t.Errorf("Expected Relative.M=%d, got %d", tt.expectRelM, time.Relative.M)
				}
				if time.Relative.D != tt.expectRelD {
					t.Errorf("Expected Relative.D=%d, got %d", tt.expectRelD, time.Relative.D)
				}
			}

			if time.Relative.Weekday != tt.expectWeekday {
				t.Errorf("Expected Relative.Weekday=%d, got %d", tt.expectWeekday, time.Relative.Weekday)
			}
			if time.Relative.WeekdayBehavior != tt.expectWeekdayBehavior {
				t.Errorf("Expected Relative.WeekdayBehavior=%d, got %d", tt.expectWeekdayBehavior, time.Relative.WeekdayBehavior)
			}
		})
	}
}
