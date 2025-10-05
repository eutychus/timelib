package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestParseDateTzIdentifier tests timezone identifier parsing
// C tests tz_identifier_00 through tz_identifier_10
// These test that timezone identifiers like "Europe/Amsterdam" are parsed correctly

func TestParseDateTzIdentifier(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectY      int64
		expectM      int64
		expectD      int64
		expectH      int64
		expectI      int64
		expectS      int64
		expectUs     int64
		expectTzName string
		checkDate    bool
		checkTime    bool
		checkUs      bool
	}{
		{
			name:         "tz_identifier_00",
			input:        "01:00:03.12345 Europe/Amsterdam",
			expectH:      1,
			expectI:      0,
			expectS:      3,
			expectUs:     123450,
			expectTzName: "Europe/Amsterdam",
			checkTime:    true,
			checkUs:      true,
		},
		{
			name:         "tz_identifier_01",
			input:        "01:00:03.12345 America/Indiana/Knox",
			expectH:      1,
			expectI:      0,
			expectS:      3,
			expectUs:     123450,
			expectTzName: "America/Indiana/Knox",
			checkTime:    true,
			checkUs:      true,
		},
		{
			name:         "tz_identifier_02",
			input:        "2005-07-14 22:30:41 America/Los_Angeles",
			expectY:      2005,
			expectM:      7,
			expectD:      14,
			expectH:      22,
			expectI:      30,
			expectS:      41,
			expectTzName: "America/Los_Angeles",
			checkDate:    true,
			checkTime:    true,
		},
		{
			name:         "tz_identifier_03",
			input:        "2005-07-14\t22:30:41\tAmerica/Los_Angeles",
			expectY:      2005,
			expectM:      7,
			expectD:      14,
			expectH:      22,
			expectI:      30,
			expectS:      41,
			expectTzName: "America/Los_Angeles",
			checkDate:    true,
			checkTime:    true,
		},
		{
			name:         "tz_identifier_04",
			input:        "Africa/Dar_es_Salaam",
			expectTzName: "Africa/Dar_es_Salaam",
		},
		{
			name:         "tz_identifier_05",
			input:        "Africa/Porto-Novo",
			expectTzName: "Africa/Porto-Novo",
		},
		{
			name:         "tz_identifier_06",
			input:        "America/Blanc-Sablon",
			expectTzName: "America/Blanc-Sablon",
		},
		{
			name:         "tz_identifier_07",
			input:        "America/Port-au-Prince",
			expectTzName: "America/Port-au-Prince",
		},
		{
			name:         "tz_identifier_08",
			input:        "America/Port_of_Spain",
			expectTzName: "America/Port_of_Spain",
		},
		{
			name:         "tz_identifier_09",
			input:        "Antarctica/DumontDUrville",
			expectTzName: "Antarctica/DumontDUrville",
		},
		{
			name:         "tz_identifier_10",
			input:        "Antarctica/McMurdo",
			expectTzName: "Antarctica/McMurdo",
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

			if tt.checkUs {
				if time.US != tt.expectUs {
					t.Errorf("Expected Us=%d, got %d", tt.expectUs, time.US)
				}
			}

			// Check timezone name
			if time.TzInfo == nil {
				t.Errorf("Expected TzInfo to be set, got nil")
			} else if time.TzInfo.Name != tt.expectTzName {
				t.Errorf("Expected TzInfo.Name=%s, got %s", tt.expectTzName, time.TzInfo.Name)
			}
		})
	}
}
