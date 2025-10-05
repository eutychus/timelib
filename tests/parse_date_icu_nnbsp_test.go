package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateIcuNnbsp tests ICU non-breaking space handling
// C tests icu_nnbsp_* series
// These test various Unicode space characters: NNBSP (U+202F narrow no-break space) and NBSP (U+00A0 no-break space)
// The parser should handle these Unicode spaces the same as regular spaces

const (
	NBSP  = "\u00A0" // U+00A0 No-Break Space
	NNBSP = "\u202F" // U+202F Narrow No-Break Space
)

func TestParseDateIcuNnbsp(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectY   int64
		expectM   int64
		expectD   int64
		expectH   int64
		expectI   int64
		expectS   int64
		expectZ   int32
		checkDate bool
		checkTZ   bool
	}{
		// Time with NNBSP before AM/PM
		{
			name:    "icu_nnbsp_timetiny12",
			input:   "8" + NNBSP + "pm",
			expectH: 20,
			expectI: 0,
			expectS: 0,
		},
		{
			name:    "icu_nnbsp_timeshort12_01",
			input:   "8:43" + NNBSP + "pm",
			expectH: 20,
			expectI: 43,
			expectS: 0,
		},
		{
			name:    "icu_nnbsp_timeshort12_02",
			input:   "8:43" + NNBSP + NNBSP + "pm",
			expectH: 20,
			expectI: 43,
			expectS: 0,
		},
		{
			name:    "icu_nnbsp_timelong12",
			input:   "8:43.43" + NNBSP + "pm",
			expectH: 20,
			expectI: 43,
			expectS: 43,
		},
		// ISO 8601 with various space characters before timezone
		{
			name:    "icu_nnbsp_iso8601normtz_00",
			input:   "T17:21:49GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		{
			name:    "icu_nnbsp_iso8601normtz_01",
			input:   "T17:21:49" + NNBSP + "GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		{
			name:    "icu_nnbsp_iso8601normtz_02",
			input:   "T17:21:49" + NNBSP + NNBSP + "GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		{
			name:    "icu_nnbsp_iso8601normtz_03",
			input:   "T17:21:49" + NBSP + "GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		{
			name:    "icu_nnbsp_iso8601normtz_04",
			input:   "T17:21:49" + NNBSP + NBSP + "GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		{
			name:    "icu_nnbsp_iso8601normtz_05",
			input:   "T17:21:49" + NBSP + NNBSP + "GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		{
			name:    "icu_nnbsp_iso8601normtz_06",
			input:   "T17:21:49" + NBSP + NBSP + "GMT+0230",
			expectH: 17,
			expectI: 21,
			expectS: 49,
			expectZ: 9000,
			checkTZ: true,
		},
		// CLF (Common Log Format) with NNBSP
		{
			name:      "icu_nnbsp_clf_01",
			input:     "10/Oct/2000:13:55:36" + NNBSP + "-0230",
			expectY:   2000,
			expectM:   10,
			expectD:   10,
			expectH:   13,
			expectI:   55,
			expectS:   36,
			expectZ:   -9000,
			checkDate: true,
			checkTZ:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := timelib.StrToTime(tt.input, nil)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			defer timelib.TimeDtor(time)

			if tt.checkDate {
				if time.Y != tt.expectY {
					t.Errorf("Expected Y=%d, got %d", tt.expectY, time.Y)
				}
				if time.M != tt.expectM {
					t.Errorf("Expected M=%d, got %d", tt.expectM, time.M)
				}
				if time.D != tt.expectD {
					t.Errorf("Expected D=%d, got %d", tt.expectD, time.D)
				}
			}

			if time.H != tt.expectH {
				t.Errorf("Expected H=%d, got %d", tt.expectH, time.H)
			}
			if time.I != tt.expectI {
				t.Errorf("Expected I=%d, got %d", tt.expectI, time.I)
			}
			if time.S != tt.expectS {
				t.Errorf("Expected S=%d, got %d", tt.expectS, time.S)
			}

			if tt.checkTZ {
				if time.Z != tt.expectZ {
					t.Errorf("Expected Z=%d, got %d", tt.expectZ, time.Z)
				}
			}
		})
	}
}
