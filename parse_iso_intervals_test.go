package timelib

import (
	"testing"
)

func TestParseIsoIntervalPeriod(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, *Time, *Time, *RelTime, int, *ErrorContainer)
	}{
		{
			name:    "Simple period",
			input:   "P1Y2M10D",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if period == nil {
					t.Fatal("Expected period to be non-nil")
				}
				if period.Y != 1 {
					t.Errorf("Expected year 1, got %d", period.Y)
				}
				if period.M != 2 {
					t.Errorf("Expected month 2, got %d", period.M)
				}
				if period.D != 10 {
					t.Errorf("Expected day 10, got %d", period.D)
				}
			},
		},
		{
			name:    "Period with time",
			input:   "P1Y2M10DT2H30M",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if period == nil {
					t.Fatal("Expected period to be non-nil")
				}
				if period.Y != 1 {
					t.Errorf("Expected year 1, got %d", period.Y)
				}
				if period.M != 2 {
					t.Errorf("Expected month 2, got %d", period.M)
				}
				if period.D != 10 {
					t.Errorf("Expected day 10, got %d", period.D)
				}
				if period.H != 2 {
					t.Errorf("Expected hour 2, got %d", period.H)
				}
				if period.I != 30 {
					t.Errorf("Expected minute 30, got %d", period.I)
				}
			},
		},
		{
			name:    "Period with weeks",
			input:   "P4W",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if period == nil {
					t.Fatal("Expected period to be non-nil")
				}
				if period.D != 28 {
					t.Errorf("Expected day 28 (4 weeks), got %d", period.D)
				}
			},
		},
		{
			name:    "Period with seconds",
			input:   "PT30S",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if period == nil {
					t.Fatal("Expected period to be non-nil")
				}
				if period.S != 30 {
					t.Errorf("Expected second 30, got %d", period.S)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			begin, end, period, recur, errors := ParseIsoInterval(tt.input)

			if (errors != nil && errors.ErrorCount > 0) != tt.wantErr {
				t.Errorf("ParseIsoInterval() error count = %d, wantErr %v", errors.ErrorCount, tt.wantErr)
				if errors != nil && errors.ErrorCount > 0 {
					for _, e := range errors.ErrorMessages {
						t.Logf("  Error: %s at position %d", e.Message, e.Position)
					}
				}
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, begin, end, period, recur, errors)
			}
		})
	}
}

func TestParseIsoIntervalDates(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, *Time, *Time, *RelTime, int, *ErrorContainer)
	}{
		{
			name:    "Basic datetime",
			input:   "20080301T130000Z",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if begin == nil {
					t.Fatal("Expected begin to be non-nil")
				}
				if begin.Y != 2008 {
					t.Errorf("Expected year 2008, got %d", begin.Y)
				}
				if begin.M != 3 {
					t.Errorf("Expected month 3, got %d", begin.M)
				}
				if begin.D != 1 {
					t.Errorf("Expected day 1, got %d", begin.D)
				}
				if begin.H != 13 {
					t.Errorf("Expected hour 13, got %d", begin.H)
				}
				if begin.I != 0 {
					t.Errorf("Expected minute 0, got %d", begin.I)
				}
				if begin.S != 0 {
					t.Errorf("Expected second 0, got %d", begin.S)
				}
			},
		},
		{
			name:    "Extended datetime",
			input:   "2008-03-01T13:00:00Z",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if begin == nil {
					t.Fatal("Expected begin to be non-nil")
				}
				if begin.Y != 2008 {
					t.Errorf("Expected year 2008, got %d", begin.Y)
				}
				if begin.M != 3 {
					t.Errorf("Expected month 3, got %d", begin.M)
				}
				if begin.D != 1 {
					t.Errorf("Expected day 1, got %d", begin.D)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			begin, end, period, recur, errors := ParseIsoInterval(tt.input)

			if (errors != nil && errors.ErrorCount > 0) != tt.wantErr {
				t.Errorf("ParseIsoInterval() error count = %d, wantErr %v", errors.ErrorCount, tt.wantErr)
				if errors != nil && errors.ErrorCount > 0 {
					for _, e := range errors.ErrorMessages {
						t.Logf("  Error: %s at position %d", e.Message, e.Position)
					}
				}
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, begin, end, period, recur, errors)
			}
		})
	}
}

func TestParseIsoIntervalRecurrences(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, *Time, *Time, *RelTime, int, *ErrorContainer)
	}{
		{
			name:    "Recurrence with period",
			input:   "R5",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if recur != 5 {
					t.Errorf("Expected recurrences 5, got %d", recur)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			begin, end, period, recur, errors := ParseIsoInterval(tt.input)

			if (errors != nil && errors.ErrorCount > 0) != tt.wantErr {
				t.Errorf("ParseIsoInterval() error count = %d, wantErr %v", errors.ErrorCount, tt.wantErr)
				if errors != nil && errors.ErrorCount > 0 {
					for _, e := range errors.ErrorMessages {
						t.Logf("  Error: %s at position %d", e.Message, e.Position)
					}
				}
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, begin, end, period, recur, errors)
			}
		})
	}
}

func TestParseIsoIntervalCombinedRep(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, *Time, *Time, *RelTime, int, *ErrorContainer)
	}{
		{
			name:    "Combined representation",
			input:   "P0001-02-03T04:05:06",
			wantErr: false,
			checkFn: func(t *testing.T, begin, end *Time, period *RelTime, recur int, errors *ErrorContainer) {
				if period == nil {
					t.Fatal("Expected period to be non-nil")
				}
				if period.Y != 1 {
					t.Errorf("Expected year 1, got %d", period.Y)
				}
				if period.M != 2 {
					t.Errorf("Expected month 2, got %d", period.M)
				}
				if period.D != 3 {
					t.Errorf("Expected day 3, got %d", period.D)
				}
				if period.H != 4 {
					t.Errorf("Expected hour 4, got %d", period.H)
				}
				if period.I != 5 {
					t.Errorf("Expected minute 5, got %d", period.I)
				}
				if period.S != 6 {
					t.Errorf("Expected second 6, got %d", period.S)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			begin, end, period, recur, errors := ParseIsoInterval(tt.input)

			if (errors != nil && errors.ErrorCount > 0) != tt.wantErr {
				t.Errorf("ParseIsoInterval() error count = %d, wantErr %v", errors.ErrorCount, tt.wantErr)
				if errors != nil && errors.ErrorCount > 0 {
					for _, e := range errors.ErrorMessages {
						t.Logf("  Error: %s at position %d", e.Message, e.Position)
					}
				}
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, begin, end, period, recur, errors)
			}
		})
	}
}
