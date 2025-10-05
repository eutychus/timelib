package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// TestPhpRfcInterval tests DST transition edge cases for PHP RFC compliance
// These tests ensure proper behavior when adding/subtracting intervals around
// the DST transition on 2010-11-07 in America/New_York timezone

func testAddWallPhpRfc(t *testing.T, timeStr, intervalStr, tzid string) *timelib.Time {
	// Parse timezone
	var dummyError int
	tzi, err := timelib.ParseTzfile(tzid, timelib.BuiltinDB(), &dummyError)
	if err != nil {
		t.Fatalf("Failed to parse timezone %s: %v", tzid, err)
	}
	if dummyError != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Fatalf("Timezone parse error for %s: %d", tzid, dummyError)
	}

	// Parse time string
	tm, err := timelib.StrToTime(timeStr, nil)
	if err != nil {
		t.Fatalf("Failed to parse time string %s: %v", timeStr, err)
	}

	// Set timezone and update
	tm.UpdateTS(tzi)
	timelib.SetTimezone(tm, tzi)
	tm.UpdateFromSSE()

	// Parse interval
	errors := &timelib.ErrorContainer{}
	_, _, period, _, parseErr := timelib.Strtointerval(intervalStr, errors)
	if parseErr != nil {
		t.Fatalf("Failed to parse interval %s: %v", intervalStr, parseErr)
	}

	// Add interval
	changed := tm.AddWall(period)

	return changed
}

func testSubWallPhpRfc(t *testing.T, timeStr, intervalStr, tzid string) *timelib.Time {
	// Parse timezone
	var dummyError int
	tzi, err := timelib.ParseTzfile(tzid, timelib.BuiltinDB(), &dummyError)
	if err != nil {
		t.Fatalf("Failed to parse timezone %s: %v", tzid, err)
	}
	if dummyError != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Fatalf("Timezone parse error for %s: %d", tzid, dummyError)
	}

	// Parse time string
	tm, err := timelib.StrToTime(timeStr, nil)
	if err != nil {
		t.Fatalf("Failed to parse time string %s: %v", timeStr, err)
	}

	// Set timezone and update
	tm.UpdateTS(tzi)
	timelib.SetTimezone(tm, tzi)
	tm.UpdateFromSSE()

	// Parse interval
	errors := &timelib.ErrorContainer{}
	_, _, period, _, parseErr := timelib.Strtointerval(intervalStr, errors)
	if parseErr != nil {
		t.Fatalf("Failed to parse interval %s: %v", intervalStr, parseErr)
	}

	// Subtract interval
	changed := tm.SubWall(period)

	return changed
}

func TestPhpRfcInterval_ba1(t *testing.T) {
	changed := testAddWallPhpRfc(t, "@1289109599", "PT1S", "America/New_York")

	if changed.Sse != 1289109600 {
		t.Errorf("Expected SSE 1289109600, got %d", changed.Sse)
	}
	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute 0, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba2(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 04:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba3(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 04:30:00", "PT24H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 3 {
		t.Errorf("Expected hour 3, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba4(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 04:30:00", "PT23H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 2 {
		t.Errorf("Expected hour 2, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba5(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 04:30:00", "PT22H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba6(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 04:30:00", "PT21H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba7(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 01:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba8(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 01:30:00", "P1DT1H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba9(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 04:30:00", "PT25H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba10(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 03:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 3 {
		t.Errorf("Expected hour 3, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_ba11(t *testing.T) {
	changed := testAddWallPhpRfc(t, "2010-11-06 02:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 2 {
		t.Errorf("Expected hour 2, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST 0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EST" {
		t.Errorf("Expected timezone abbr EST, got %s", changed.TzAbbr)
	}
}

// Subtraction tests (bs1-bs10)

func TestPhpRfcInterval_bs1(t *testing.T) {
	changed := testSubWallPhpRfc(t, "@1289109600", "PT1S", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 7 {
		t.Errorf("Expected day 7, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 59 {
		t.Errorf("Expected minute 59, got %d", changed.I)
	}
	if changed.S != 59 {
		t.Errorf("Expected second 59, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs2(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 04:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs3(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 03:30:00", "PT24H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs4(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 02:30:00", "PT23H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs5(t *testing.T) {
	changed := testSubWallPhpRfc(t, "@1289111400", "PT22H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs6(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 01:30:00", "PT21H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 4 {
		t.Errorf("Expected hour 4, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs7(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 01:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour 1, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs8(t *testing.T) {
	changed := testSubWallPhpRfc(t, "@1289111400", "P1DT1H", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour 0, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs9(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 03:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 3 {
		t.Errorf("Expected hour 3, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}

func TestPhpRfcInterval_bs10(t *testing.T) {
	changed := testSubWallPhpRfc(t, "2010-11-07 02:30:00", "P1D", "America/New_York")

	if changed.Y != 2010 {
		t.Errorf("Expected year 2010, got %d", changed.Y)
	}
	if changed.M != 11 {
		t.Errorf("Expected month 11, got %d", changed.M)
	}
	if changed.D != 6 {
		t.Errorf("Expected day 6, got %d", changed.D)
	}
	if changed.H != 2 {
		t.Errorf("Expected hour 2, got %d", changed.H)
	}
	if changed.I != 30 {
		t.Errorf("Expected minute 30, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second 0, got %d", changed.S)
	}
	if changed.Dst != 1 {
		t.Errorf("Expected DST 1, got %d", changed.Dst)
	}
	if changed.TzAbbr != "EDT" {
		t.Errorf("Expected timezone abbr EDT, got %s", changed.TzAbbr)
	}
}
