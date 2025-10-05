package timelib

import (
	"testing"
	"time"
)

func TestStrtotime(t *testing.T) {
	tests := []struct {
		name          string
		datetime      string
		baseTimestamp int64
		want          int64 // Expected Unix timestamp
		wantErr       bool  // If true, expect -1
	}{
		{
			name:          "absolute date",
			datetime:      "2024-01-15",
			baseTimestamp: 0,
			want:          1705276800, // 2024-01-15 00:00:00 UTC
			wantErr:       false,
		},
		{
			name:          "absolute datetime",
			datetime:      "2024-01-15 12:30:45",
			baseTimestamp: 0,
			want:          1705321845, // 2024-01-15 12:30:45 UTC
			wantErr:       false,
		},
		{
			name:          "tomorrow from base time",
			datetime:      "tomorrow",
			baseTimestamp: 1704067200, // 2024-01-01 00:00:00 UTC
			want:          1704153600, // 2024-01-02 00:00:00 UTC
			wantErr:       false,
		},
		{
			name:          "plus one day",
			datetime:      "+1 day",
			baseTimestamp: 1704067200, // 2024-01-01 00:00:00 UTC
			want:          1704153600, // 2024-01-02 00:00:00 UTC
			wantErr:       false,
		},
		{
			name:          "plus one week",
			datetime:      "+1 week",
			baseTimestamp: 1704067200, // 2024-01-01 00:00:00 UTC
			want:          1704672000, // 2024-01-08 00:00:00 UTC
			wantErr:       false,
		},
		{
			name:          "minus one day",
			datetime:      "-1 day",
			baseTimestamp: 1704067200, // 2024-01-01 00:00:00 UTC
			want:          1703980800, // 2023-12-31 00:00:00 UTC
			wantErr:       false,
		},
		{
			name:          "next monday",
			datetime:      "next monday",
			baseTimestamp: 1704067200, // 2024-01-01 00:00:00 UTC (Monday)
			want:          1704672000, // 2024-01-08 00:00:00 UTC (next Monday)
			wantErr:       false,
		},
		{
			name:          "invalid string",
			datetime:      "not a date",
			baseTimestamp: 0,
			want:          -1,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Strtotime(tt.datetime, tt.baseTimestamp)

			if tt.wantErr {
				if got != -1 {
					t.Errorf("Strtotime() = %v, want -1 for error", got)
				}
				return
			}

			if got != tt.want {
				t.Errorf("Strtotime() = %v, want %v", got, tt.want)
				// Print the dates for debugging
				t.Logf("Got: %s", time.Unix(got, 0).UTC().Format(time.RFC3339))
				t.Logf("Want: %s", time.Unix(tt.want, 0).UTC().Format(time.RFC3339))
			}
		})
	}
}

func TestStrtotimeToGoTime(t *testing.T) {
	tests := []struct {
		name     string
		datetime string
		baseTime *time.Time
		want     time.Time
		wantErr  bool
	}{
		{
			name:     "absolute date",
			datetime: "2024-01-15",
			baseTime: nil,
			want:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "absolute datetime",
			datetime: "2024-01-15 12:30:45",
			baseTime: nil,
			want:     time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "tomorrow from base time",
			datetime: "tomorrow",
			baseTime: func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
			want:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "plus one week",
			datetime: "+1 week",
			baseTime: func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
			want:     time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "invalid string",
			datetime: "not a date",
			baseTime: nil,
			want:     time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrtotimeToGoTime(tt.datetime, tt.baseTime)

			if (err != nil) != tt.wantErr {
				t.Errorf("StrtotimeToGoTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !got.Equal(tt.want) {
				t.Errorf("StrtotimeToGoTime() = %v, want %v", got, tt.want)
				t.Logf("Got: %s", got.Format(time.RFC3339))
				t.Logf("Want: %s", tt.want.Format(time.RFC3339))
			}
		})
	}
}

func TestStrtotimeWithTimezone(t *testing.T) {
	// Load a specific timezone for testing
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load America/New_York timezone")
	}

	tests := []struct {
		name     string
		datetime string
		baseTime *time.Time
		location *time.Location
		want     time.Time
		wantErr  bool
	}{
		{
			name:     "absolute date in New York timezone",
			datetime: "2024-01-15 12:00:00",
			baseTime: nil,
			location: ny,
			want:     time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC).In(ny),
			wantErr:  false,
		},
		{
			name:     "tomorrow in UTC",
			datetime: "tomorrow",
			baseTime: func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
			location: time.UTC,
			want:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrtotimeWithTimezone(tt.datetime, tt.baseTime, tt.location)

			if (err != nil) != tt.wantErr {
				t.Errorf("StrtotimeWithTimezone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !got.Equal(tt.want) {
				t.Errorf("StrtotimeWithTimezone() = %v, want %v", got, tt.want)
				t.Logf("Got: %s", got.Format(time.RFC3339))
				t.Logf("Want: %s", tt.want.Format(time.RFC3339))
			}
		})
	}
}

// TestStrtotimeWithCurrentTime tests relative times with current time as base
func TestStrtotimeWithCurrentTime(t *testing.T) {
	// Use current time as base
	now := time.Now()

	// Test "tomorrow" with current time
	ts := Strtotime("tomorrow", now.Unix())
	if ts == -1 {
		t.Fatal("Strtotime() returned error for 'tomorrow'")
	}

	// Verify it's approximately 24 hours in the future
	expected := now.Add(24 * time.Hour).Unix()
	diff := ts - expected

	// Allow for some time passage during test execution (within 2 seconds)
	if diff < -2 || diff > 2 {
		t.Errorf("Strtotime('tomorrow', now) = %v, expected ~%v (diff: %v seconds)",
			ts, expected, diff)
		t.Logf("Got: %s", time.Unix(ts, 0).Format(time.RFC3339))
		t.Logf("Expected: %s", time.Unix(expected, 0).Format(time.RFC3339))
	}
}

// TestStrtotimeUseCases demonstrates real-world use cases
func TestStrtotimeUseCases(t *testing.T) {
	// Use case 1: Parse a date string for database query
	ts := Strtotime("2024-01-01", 0)
	if ts == -1 {
		t.Fatal("Failed to parse date")
	}
	t.Logf("Database timestamp: %d", ts)

	// Use case 2: Calculate expiration time (30 days from now)
	now := time.Now().Unix()
	expiration := Strtotime("+30 days", now)
	if expiration == -1 {
		t.Fatal("Failed to calculate expiration")
	}
	t.Logf("Expiration timestamp: %d", expiration)

	// Use case 3: Get time.Time for formatting
	goTime, err := StrtotimeToGoTime("2024-06-15 10:30:00", nil)
	if err != nil {
		t.Fatal(err)
	}
	formatted := goTime.Format("January 2, 2006 at 3:04 PM")
	t.Logf("Formatted time: %s", formatted)

	// Use case 4: Parse with timezone
	loc, _ := time.LoadLocation("America/Los_Angeles")
	localTime, err := StrtotimeWithTimezone("2024-12-25 09:00:00", nil, loc)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Christmas morning in LA: %s", localTime.Format(time.RFC3339))
}
