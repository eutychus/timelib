package timelib

import (
	"testing"
)

func TestBasicDateParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		checkFn  func(*testing.T, *Time)
	}{
		{
			name:    "time without colon (gnunocolon)",
			input:   "2025",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime {
					t.Error("Expected HaveTime to be true")
				}
				if tm.H != 20 {
					t.Errorf("Expected hour 20, got %d", tm.H)
				}
				if tm.I != 25 {
					t.Errorf("Expected minute 25, got %d", tm.I)
				}
			},
		},
		{
			name:    "yesterday",
			input:   "yesterday",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative to be true")
				}
				if tm.Relative.D != -1 {
					t.Errorf("Expected relative day -1, got %d", tm.Relative.D)
				}
			},
		},
		{
			name:    "tomorrow",
			input:   "tomorrow",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveRelative {
					t.Error("Expected HaveRelative to be true")
				}
				if tm.Relative.D != 1 {
					t.Errorf("Expected relative day 1, got %d", tm.Relative.D)
				}
			},
		},
		{
			name:    "now",
			input:   "now",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				// "now" doesn't set HaveRelative - it just returns TIMELIB_RELATIVE token
				// The actual interpretation happens later in timelib_update_ts
				// For a basic parse, it returns no flags set
			},
		},
		{
			name:    "noon",
			input:   "noon",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveTime {
					t.Error("Expected HaveTime to be true")
				}
				if tm.H != 12 {
					t.Errorf("Expected hour 12, got %d", tm.H)
				}
			},
		},
		{
			name:    "ISO date",
			input:   "2025-01-15",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate {
					t.Error("Expected HaveDate to be true")
				}
				if tm.Y != 2025 {
					t.Errorf("Expected year 2025, got %d", tm.Y)
				}
				if tm.M != 1 {
					t.Errorf("Expected month 1, got %d", tm.M)
				}
				if tm.D != 15 {
					t.Errorf("Expected day 15, got %d", tm.D)
				}
			},
		},
		{
			name:    "American date",
			input:   "1/15/2025",
			wantErr: false,
			checkFn: func(t *testing.T, tm *Time) {
				if !tm.HaveDate {
					t.Error("Expected HaveDate to be true")
				}
				if tm.M != 1 {
					t.Errorf("Expected month 1, got %d", tm.M)
				}
				if tm.D != 15 {
					t.Errorf("Expected day 15, got %d", tm.D)
				}
				if tm.Y != 2025 {
					t.Errorf("Expected year 2025, got %d", tm.Y)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("StrToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, time)
			}
		})
	}
}

func TestTimestampParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Unix timestamp",
			input:   "@1234567890",
			wantErr: false,
		},
		{
			name:    "Unix timestamp with milliseconds",
			input:   "@1234567890.123456",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("StrToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if time.Y != 1970 {
					t.Errorf("Expected base year 1970 for timestamp, got %d", time.Y)
				}
				if !time.HaveRelative {
					t.Error("Expected HaveRelative to be true for timestamp")
				}
			}
		})
	}
}

func TestRelativeDateParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "plus one day",
			input:   "+1 day",
			wantErr: false,
		},
		{
			name:    "minus one week",
			input:   "-1 week",
			wantErr: false,
		},
		{
			name:    "next monday",
			input:   "next monday",
			wantErr: false,
		},
		{
			name:    "first day of",
			input:   "first day of",
			wantErr: false,
		},
		{
			name:    "last day of",
			input:   "last day of",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := StrToTime(tt.input, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("StrToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !time.HaveRelative {
				t.Error("Expected HaveRelative to be true for relative date")
			}
		})
	}
}
