package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestRenderTSZoneinfo(t *testing.T) {
	testCases := []struct {
		name     string
		ts       int64
		timezone string
		wantErr  bool
	}{
		{
			name:     "Europe/Amsterdam timestamp",
			ts:       1114819200,
			timezone: "Europe/Amsterdam",
			wantErr:  false,
		},
		{
			name:     "America/New_York timestamp",
			ts:       1234567890,
			timezone: "America/New_York",
			wantErr:  false,
		},
		{
			name:     "UTC timestamp",
			ts:       1609459200,
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "Invalid timezone",
			ts:       1114819200,
			timezone: "Invalid/Timezone",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create time object
			time := timelib.TimeCtor()
			defer timelib.TimeDtor(time)

			// Parse timezone
			var dummyError int
			tzinfo, err := timelib.ParseTzfile(tc.timezone, timelib.BuiltinDB(), &dummyError)
			if tc.wantErr {
				if err == nil && dummyError == timelib.TIMELIB_ERROR_NO_ERROR {
					t.Fatalf("Expected error for timezone %s, got nil", tc.timezone)
				}
				return
			}
			if err != nil {
				t.Fatalf("Failed to parse timezone %s: %v", tc.timezone, err)
			}

			// Set timezone and convert timestamp
			timelib.SetTimezone(time, tzinfo)
			time.Unixtime2local(tc.ts)

			// Verify the time object is valid
			if time.Y == 0 && time.M == 0 && time.D == 0 {
				t.Error("Expected non-zero date values after conversion")
			}

			// The original C code calls timelib_dump_date for debugging
			// We just verify the conversion worked
			t.Logf("Timestamp %d in %s: %04d-%02d-%02d %02d:%02d:%02d",
				tc.ts, tc.timezone, time.Y, time.M, time.D, time.H, time.I, time.S)
		})
	}
}
