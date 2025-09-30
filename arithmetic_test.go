package timelib

import (
	"testing"
)

func TestTimeAdd(t *testing.T) {
	tests := []struct {
		name     string
		base     *Time
		interval *RelTime
		expected *Time
	}{
		{
			name:     "Add days",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{D: 5},
			expected: &Time{Y: 2023, M: 6, D: 20, H: 10, I: 30, S: 45},
		},
		{
			name:     "Add months",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{M: 3},
			expected: &Time{Y: 2023, M: 9, D: 15, H: 10, I: 30, S: 45},
		},
		{
			name:     "Add years",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{Y: 2},
			expected: &Time{Y: 2025, M: 6, D: 15, H: 10, I: 30, S: 45},
		},
		{
			name:     "Add hours",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{H: 5},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 15, I: 30, S: 45},
		},
		{
			name:     "Add minutes",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{I: 45},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 11, I: 15, S: 45},
		},
		{
			name:     "Add seconds",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{S: 30},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 10, I: 31, S: 15},
		},
		{
			name:     "Add microseconds",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45, US: 500000},
			interval: &RelTime{US: 600000},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 46, US: 100000},
		},
		{
			name:     "Add across month boundary",
			base:     &Time{Y: 2023, M: 6, D: 28, H: 10, I: 30, S: 45},
			interval: &RelTime{D: 5},
			expected: &Time{Y: 2023, M: 7, D: 3, H: 10, I: 30, S: 45},
		},
		{
			name:     "Add across year boundary",
			base:     &Time{Y: 2023, M: 12, D: 28, H: 10, I: 30, S: 45},
			interval: &RelTime{D: 5},
			expected: &Time{Y: 2024, M: 1, D: 2, H: 10, I: 30, S: 45},
		},
		{
			name:     "Add leap year day",
			base:     &Time{Y: 2024, M: 2, D: 28, H: 10, I: 30, S: 45},
			interval: &RelTime{D: 1},
			expected: &Time{Y: 2024, M: 2, D: 29, H: 10, I: 30, S: 45},
		},
		{
			name:     "Add complex interval",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{Y: 1, M: 2, D: 3, H: 4, I: 5, S: 6},
			expected: &Time{Y: 2024, M: 8, D: 18, H: 14, I: 35, S: 51},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.base.Add(test.interval)

			if result.Y != test.expected.Y {
				t.Errorf("Year: expected %d, got %d", test.expected.Y, result.Y)
			}
			if result.M != test.expected.M {
				t.Errorf("Month: expected %d, got %d", test.expected.M, result.M)
			}
			if result.D != test.expected.D {
				t.Errorf("Day: expected %d, got %d", test.expected.D, result.D)
			}
			if result.H != test.expected.H {
				t.Errorf("Hour: expected %d, got %d", test.expected.H, result.H)
			}
			if result.I != test.expected.I {
				t.Errorf("Minute: expected %d, got %d", test.expected.I, result.I)
			}
			if result.S != test.expected.S {
				t.Errorf("Second: expected %d, got %d", test.expected.S, result.S)
			}
			if result.US != test.expected.US {
				t.Errorf("Microsecond: expected %d, got %d", test.expected.US, result.US)
			}
		})
	}
}

func TestTimeSub(t *testing.T) {
	tests := []struct {
		name     string
		base     *Time
		interval *RelTime
		expected *Time
	}{
		{
			name:     "Subtract days",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{D: 5},
			expected: &Time{Y: 2023, M: 6, D: 10, H: 10, I: 30, S: 45},
		},
		{
			name:     "Subtract months",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{M: 3},
			expected: &Time{Y: 2023, M: 3, D: 15, H: 10, I: 30, S: 45},
		},
		{
			name:     "Subtract hours",
			base:     &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			interval: &RelTime{H: 5},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 5, I: 30, S: 45},
		},
		{
			name:     "Subtract across month boundary",
			base:     &Time{Y: 2023, M: 6, D: 5, H: 10, I: 30, S: 45},
			interval: &RelTime{D: 10},
			expected: &Time{Y: 2023, M: 5, D: 26, H: 10, I: 30, S: 45},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.base.Sub(test.interval)

			if result.Y != test.expected.Y {
				t.Errorf("Year: expected %d, got %d", test.expected.Y, result.Y)
			}
			if result.M != test.expected.M {
				t.Errorf("Month: expected %d, got %d", test.expected.M, result.M)
			}
			if result.D != test.expected.D {
				t.Errorf("Day: expected %d, got %d", test.expected.D, result.D)
			}
			if result.H != test.expected.H {
				t.Errorf("Hour: expected %d, got %d", test.expected.H, result.H)
			}
			if result.I != test.expected.I {
				t.Errorf("Minute: expected %d, got %d", test.expected.I, result.I)
			}
			if result.S != test.expected.S {
				t.Errorf("Second: expected %d, got %d", test.expected.S, result.S)
			}
		})
	}
}

func TestTimeDiff(t *testing.T) {
	tests := []struct {
		name     string
		one      *Time
		two      *Time
		expected *RelTime
	}{
		{
			name:     "Same time",
			one:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			two:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			expected: &RelTime{Y: 0, M: 0, D: 0, H: 0, I: 0, S: 0, US: 0, Invert: false},
		},
		{
			name:     "One day difference",
			one:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			two:      &Time{Y: 2023, M: 6, D: 16, H: 10, I: 30, S: 45},
			expected: &RelTime{Y: 0, M: 0, D: 1, H: 0, I: 0, S: 0, US: 0, Invert: false},
		},
		{
			name:     "One month difference",
			one:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			two:      &Time{Y: 2023, M: 7, D: 15, H: 10, I: 30, S: 45},
			expected: &RelTime{Y: 0, M: 1, D: 0, H: 0, I: 0, S: 0, US: 0, Invert: false},
		},
		{
			name:     "One year difference",
			one:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			two:      &Time{Y: 2024, M: 6, D: 15, H: 10, I: 30, S: 45},
			expected: &RelTime{Y: 1, M: 0, D: 0, H: 0, I: 0, S: 0, US: 0, Invert: false},
		},
		{
			name:     "Complex difference",
			one:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			two:      &Time{Y: 2023, M: 6, D: 16, H: 12, I: 45, S: 30},
			expected: &RelTime{Y: 0, M: 0, D: 1, H: 2, I: 14, S: 45, US: 0, Invert: false},
		},
		{
			name:     "Reverse order (should invert)",
			one:      &Time{Y: 2023, M: 6, D: 16, H: 10, I: 30, S: 45},
			two:      &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45},
			expected: &RelTime{Y: 0, M: 0, D: 1, H: 0, I: 0, S: 0, US: 0, Invert: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.one.Diff(test.two)

			if result.Y != test.expected.Y {
				t.Errorf("Years: expected %d, got %d", test.expected.Y, result.Y)
			}
			if result.M != test.expected.M {
				t.Errorf("Months: expected %d, got %d", test.expected.M, result.M)
			}
			if result.D != test.expected.D {
				t.Errorf("Days: expected %d, got %d", test.expected.D, result.D)
			}
			if result.H != test.expected.H {
				t.Errorf("Hours: expected %d, got %d", test.expected.H, result.H)
			}
			if result.I != test.expected.I {
				t.Errorf("Minutes: expected %d, got %d", test.expected.I, result.I)
			}
			if result.S != test.expected.S {
				t.Errorf("Seconds: expected %d, got %d", test.expected.S, result.S)
			}
			if result.US != test.expected.US {
				t.Errorf("Microseconds: expected %d, got %d", test.expected.US, result.US)
			}
			if result.Invert != test.expected.Invert {
				t.Errorf("Invert: expected %v, got %v", test.expected.Invert, result.Invert)
			}
		})
	}
}

func TestTimelibDoNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    *Time
		expected *Time
	}{
		{
			name:     "Normalize microseconds overflow",
			input:    &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 45, US: 1500000},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 46, US: 500000},
		},
		{
			name:     "Normalize seconds overflow",
			input:    &Time{Y: 2023, M: 6, D: 15, H: 10, I: 30, S: 75},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 10, I: 31, S: 15},
		},
		{
			name:     "Normalize minutes overflow",
			input:    &Time{Y: 2023, M: 6, D: 15, H: 10, I: 75, S: 30},
			expected: &Time{Y: 2023, M: 6, D: 15, H: 11, I: 15, S: 30},
		},
		{
			name:     "Normalize hours overflow",
			input:    &Time{Y: 2023, M: 6, D: 15, H: 25, I: 30, S: 45},
			expected: &Time{Y: 2023, M: 6, D: 16, H: 1, I: 30, S: 45},
		},
		{
			name:     "Normalize days overflow",
			input:    &Time{Y: 2023, M: 6, D: 35, H: 10, I: 30, S: 45},
			expected: &Time{Y: 2023, M: 7, D: 5, H: 10, I: 30, S: 45},
		},
		{
			name:     "Normalize months overflow",
			input:    &Time{Y: 2023, M: 14, D: 15, H: 10, I: 30, S: 45},
			expected: &Time{Y: 2024, M: 2, D: 15, H: 10, I: 30, S: 45},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Time{
				Y:  test.input.Y,
				M:  test.input.M,
				D:  test.input.D,
				H:  test.input.H,
				I:  test.input.I,
				S:  test.input.S,
				US: test.input.US,
			}
			timelib_do_normalize(result)

			if result.Y != test.expected.Y {
				t.Errorf("Year: expected %d, got %d", test.expected.Y, result.Y)
			}
			if result.M != test.expected.M {
				t.Errorf("Month: expected %d, got %d", test.expected.M, result.M)
			}
			if result.D != test.expected.D {
				t.Errorf("Day: expected %d, got %d", test.expected.D, result.D)
			}
			if result.H != test.expected.H {
				t.Errorf("Hour: expected %d, got %d", test.expected.H, result.H)
			}
			if result.I != test.expected.I {
				t.Errorf("Minute: expected %d, got %d", test.expected.I, result.I)
			}
			if result.S != test.expected.S {
				t.Errorf("Second: expected %d, got %d", test.expected.S, result.S)
			}
			if result.US != test.expected.US {
				t.Errorf("Microsecond: expected %d, got %d", test.expected.US, result.US)
			}
		})
	}
}

func TestUnixtime2date(t *testing.T) {
	tests := []struct {
		name     string
		ts       int64
		expected struct {
			y int64
			m int64
			d int64
		}
	}{
		{
			name: "Unix epoch",
			ts:   0,
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 1970, m: 1, d: 1},
		},
		{
			name: "Current time example",
			ts:   1686800000, // Around June 2023
			expected: struct {
				y int64
				m int64
				d int64
			}{y: 2023, m: 6, d: 15},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var y, m, d int64
			Unixtime2date(test.ts, &y, &m, &d)

			if y != test.expected.y {
				t.Errorf("Year: expected %d, got %d", test.expected.y, y)
			}
			if m != test.expected.m {
				t.Errorf("Month: expected %d, got %d", test.expected.m, m)
			}
			if d != test.expected.d {
				t.Errorf("Day: expected %d, got %d", test.expected.d, d)
			}
		})
	}
}

func TestUnixtime2gmt(t *testing.T) {
	tests := []struct {
		name     string
		ts       int64
		expected *Time
	}{
		{
			name:     "Unix epoch",
			ts:       0,
			expected: &Time{Y: 1970, M: 1, D: 1, H: 0, I: 0, S: 0, US: 0, IsLocaltime: false, ZoneType: TIMELIB_ZONETYPE_NONE},
		},
		{
			name:     "Specific timestamp",
			ts:       1686800000, // 2023-06-15 03:33:20 UTC
			expected: &Time{Y: 2023, M: 6, D: 15, H: 3, I: 33, S: 20, US: 0, IsLocaltime: false, ZoneType: TIMELIB_ZONETYPE_NONE},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := TimeCtor()
			result.Unixtime2gmt(test.ts)

			if result.Y != test.expected.Y {
				t.Errorf("Year: expected %d, got %d", test.expected.Y, result.Y)
			}
			if result.M != test.expected.M {
				t.Errorf("Month: expected %d, got %d", test.expected.M, result.M)
			}
			if result.D != test.expected.D {
				t.Errorf("Day: expected %d, got %d", test.expected.D, result.D)
			}
			if result.H != test.expected.H {
				t.Errorf("Hour: expected %d, got %d", test.expected.H, result.H)
			}
			if result.I != test.expected.I {
				t.Errorf("Minute: expected %d, got %d", test.expected.I, result.I)
			}
			if result.S != test.expected.S {
				t.Errorf("Second: expected %d, got %d", test.expected.S, result.S)
			}
			if result.US != test.expected.US {
				t.Errorf("Microsecond: expected %d, got %d", test.expected.US, result.US)
			}
			if result.IsLocaltime != test.expected.IsLocaltime {
				t.Errorf("IsLocaltime: expected %v, got %v", test.expected.IsLocaltime, result.IsLocaltime)
			}
			if result.ZoneType != test.expected.ZoneType {
				t.Errorf("ZoneType: expected %d, got %d", test.expected.ZoneType, result.ZoneType)
			}
		})
	}
}
