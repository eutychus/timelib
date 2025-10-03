package timelib

import (
	"testing"
)

// TestStrtotimeWrapper verifies the new parser API matches expected behavior
func TestStrtotimeWrapper(t *testing.T) {
	// Basic test that StrToTime works
	time, err := StrToTime("2025-01-15", nil)
	if err != nil {
		t.Fatalf("StrToTime failed: %v", err)
	}
	if time.Y != 2025 || time.M != 1 || time.D != 15 {
		t.Errorf("Expected 2025-01-15, got %d-%d-%d", time.Y, time.M, time.D)
	}
}

// TestRelativeExpressions tests relative date/time parsing
func TestRelativeExpressions(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkFn    func(*testing.T, *Time)
	}{
		{
			name:  "yesterday",
			input: "yesterday",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative || tm.Relative.D != -1 {
					t.Errorf("Expected yesterday (D=-1), got D=%d, HaveRelative=%v", tm.Relative.D, tm.HaveRelative)
				}
			},
		},
		{
			name:  "tomorrow",
			input: "tomorrow",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative || tm.Relative.D != 1 {
					t.Errorf("Expected tomorrow (D=1), got D=%d", tm.Relative.D)
				}
			},
		},
		{
			name:  "+1 day",
			input: "+1 day",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative || tm.Relative.D != 1 {
					t.Errorf("Expected +1 day, got D=%d", tm.Relative.D)
				}
			},
		},
		{
			name:  "-1 week",
			input: "-1 week",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative || tm.Relative.D != -7 {
					t.Errorf("Expected -1 week (D=-7), got D=%d", tm.Relative.D)
				}
			},
		},
		{
			name:  "+2 months",
			input: "+2 months",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative || tm.Relative.M != 2 {
					t.Errorf("Expected +2 months, got M=%d", tm.Relative.M)
				}
			},
		},
		{
			name:  "next monday",
			input: "next monday",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative to be true")
				}
			},
		},
		{
			name:  "last friday",
			input: "last friday",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative to be true")
				}
			},
		},
		{
			name:  "first day of",
			input: "first day of",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative to be true")
				}
			},
		},
		{
			name:  "last day of",
			input: "last day of",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			tt.checkFn(t, time)
		})
	}
}

// TestTimeFormats tests various time format parsing
func TestTimeFormats(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(*testing.T, *Time)
	}{
		{
			name:  "noon",
			input: "noon",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime || tm.H != 12 || tm.I != 0 || tm.S != 0 {
					t.Errorf("Expected noon (12:00:00), got %d:%d:%d", tm.H, tm.I, tm.S)
				}
			},
		},
		{
			name:  "midnight",
			input: "midnight",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime {
					t.Error("Expected HaveTime to be true for midnight")
				}
				// Midnight is 0:0:0, which is correct
				if tm.H != 0 || tm.I != 0 || tm.S != 0 {
					t.Errorf("Expected midnight (00:00:00), got %d:%d:%d", tm.H, tm.I, tm.S)
				}
			},
		},
		{
			name:  "24-hour format",
			input: "14:30",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime || tm.H != 14 || tm.I != 30 {
					t.Errorf("Expected 14:30, got %d:%d", tm.H, tm.I)
				}
			},
		},
		{
			name:  "24-hour with seconds",
			input: "14:30:45",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime || tm.H != 14 || tm.I != 30 || tm.S != 45 {
					t.Errorf("Expected 14:30:45, got %d:%d:%d", tm.H, tm.I, tm.S)
				}
			},
		},
		{
			name:  "12-hour AM format",
			input: "9:30 am",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime || tm.H != 9 || tm.I != 30 {
					t.Errorf("Expected 9:30 AM, got %d:%d", tm.H, tm.I)
				}
			},
		},
		{
			name:  "12-hour PM format",
			input: "2:30 pm",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime || tm.H != 14 || tm.I != 30 {
					t.Errorf("Expected 14:30 (2:30 PM), got %d:%d", tm.H, tm.I)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			tt.checkFn(t, time)
		})
	}
}

// TestDateFormats tests various date format parsing
func TestDateFormats(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(*testing.T, *Time)
	}{
		{
			name:  "ISO date",
			input: "2025-01-15",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate || tm.Y != 2025 || tm.M != 1 || tm.D != 15 {
					t.Errorf("Expected 2025-01-15, got %d-%d-%d", tm.Y, tm.M, tm.D)
				}
			},
		},
		{
			name:  "ISO datetime",
			input: "2025-01-15T10:30:00",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate || tm.Y != 2025 || tm.M != 1 || tm.D != 15 {
					t.Errorf("Expected date 2025-01-15, got %d-%d-%d", tm.Y, tm.M, tm.D)
				}
				if !tm.HaveTime || tm.H != 10 || tm.I != 30 || tm.S != 0 {
					t.Errorf("Expected time 10:30:00, got %d:%d:%d", tm.H, tm.I, tm.S)
				}
			},
		},
		{
			name:  "American date",
			input: "1/15/2025",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate || tm.M != 1 || tm.D != 15 || tm.Y != 2025 {
					t.Errorf("Expected 1/15/2025, got %d/%d/%d", tm.M, tm.D, tm.Y)
				}
			},
		},
		{
			name:  "European date (day.month.year)",
			input: "15.01.2025",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate || tm.D != 15 || tm.M != 1 || tm.Y != 2025 {
					t.Errorf("Expected 15.01.2025, got %d.%d.%d", tm.D, tm.M, tm.Y)
				}
			},
		},
		{
			name:  "Text month",
			input: "15 January 2025",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate || tm.D != 15 || tm.M != 1 || tm.Y != 2025 {
					t.Errorf("Expected 15 January 2025, got %d %d %d", tm.D, tm.M, tm.Y)
				}
			},
		},
		{
			name:  "Abbreviated month",
			input: "Jan 15, 2025",
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate || tm.M != 1 || tm.D != 15 || tm.Y != 2025 {
					t.Errorf("Expected Jan 15, 2025, got %d %d, %d", tm.M, tm.D, tm.Y)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			tt.checkFn(t, time)
		})
	}
}

// TestTimestamps tests Unix timestamp parsing
func TestTimestamps(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(*testing.T, *Time)
	}{
		{
			name:  "Unix timestamp",
			input: "@1234567890",
			checkFn: func(t *testing.T, tm *Time) {
				if tm.Y != 1970 {
					t.Errorf("Expected base year 1970 for timestamp, got %d", tm.Y)
				}
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative for timestamp")
				}
			},
		},
		{
			name:  "Unix timestamp with decimals",
			input: "@1234567890.123456",
			checkFn: func(t *testing.T, tm *Time) {
				if tm.Y != 1970 {
					t.Errorf("Expected base year 1970 for timestamp, got %d", tm.Y)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			tt.checkFn(t, time)
		})
	}
}

// TestWeekdays tests weekday-related parsing
func TestWeekdays(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Monday", "Monday", false},
		{"Tuesday", "Tuesday", false},
		{"Wednesday", "Wednesday", false},
		{"Thursday", "Thursday", false},
		{"Friday", "Friday", false},
		{"Saturday", "Saturday", false},
		{"Sunday", "Sunday", false},
		{"Mon", "Mon", false},
		{"Tue", "Tue", false},
		{"Wed", "Wed", false},
		{"Thu", "Thu", false},
		{"Fri", "Fri", false},
		{"Sat", "Sat", false},
		{"Sun", "Sun", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !time.HaveRelative {
				t.Error("Expected HaveRelative for weekday")
			}
		})
	}
}

// TestMonthNames tests month name parsing
func TestMonthNames(t *testing.T) {
	monthTests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"January", "1 January 2025", 1},
		{"February", "1 February 2025", 2},
		{"March", "1 March 2025", 3},
		{"April", "1 April 2025", 4},
		{"May", "1 May 2025", 5},
		{"June", "1 June 2025", 6},
		{"July", "1 July 2025", 7},
		{"August", "1 August 2025", 8},
		{"September", "1 September 2025", 9},
		{"October", "1 October 2025", 10},
		{"November", "1 November 2025", 11},
		{"December", "1 December 2025", 12},
	}

	for _, tt := range monthTests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if time.M != tt.expected {
				t.Errorf("Expected month %d, got %d", tt.expected, time.M)
			}
		})
	}
}

// TestComplexExpressions tests more complex date/time expressions
func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Next week", "next week", false},
		{"Last month", "last month", false},
		{"2 weeks ago", "2 weeks ago", false},
		{"3 days ago", "3 days ago", false},
		{"In 5 hours", "in 5 hours", false},
		{"5 hours", "5 hours", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && time == nil {
				t.Error("Expected non-nil time")
			}
		})
	}
}
