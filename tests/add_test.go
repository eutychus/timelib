package tests

import (
	"strings"
	"testing"

	timelib "github.com/eutychus/timelib"
)

const SECS_PER_HOUR = 3600

// testAddWithOffset is a helper that tests adding an interval to a time with a fixed offset
func testAddWithOffset(t *testing.T, offset int32, from string, period string, expectedY, expectedM, expectedD, expectedH, expectedI, expectedS int64) {
	t.Helper()

	// Parse the from time
	tFrom, err := timelib.StrToTime(from, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", from, err)
	}

	// No need to FillHoles - the date strings in tests are complete

	// Set the offset
	tFrom.ZoneType = timelib.TIMELIB_ZONETYPE_OFFSET
	tFrom.Z = offset
	tFrom.Dst = 0  // Initialize Dst when setting timezone offset

	// Update timestamp
	tFrom.UpdateTS(nil)

	// Parse the period/interval
	var iPeriod *timelib.RelTime
	periodStr := period
	invert := false

	// Handle negative intervals
	if strings.HasPrefix(period, "-") {
		invert = true
		periodStr = period[1:]
	}

	// Parse interval
	errorsContainer := &timelib.ErrorContainer{}
	_, _, iPeriod, _, intervalErr := timelib.Strtointerval(periodStr, errorsContainer)
	if intervalErr != nil {
		t.Fatalf("Failed to parse interval '%s': %v", periodStr, intervalErr)
	}

	if invert {
		iPeriod.Invert = true
	}

	// Perform the addition using AddWall
	tAdded := tFrom.AddWall(iPeriod)

	// Check results
	if tAdded.Y != expectedY {
		t.Errorf("Year mismatch: expected %d, got %d", expectedY, tAdded.Y)
	}
	if tAdded.M != expectedM {
		t.Errorf("Month mismatch: expected %d, got %d", expectedM, tAdded.M)
	}
	if tAdded.D != expectedD {
		t.Errorf("Day mismatch: expected %d, got %d", expectedD, tAdded.D)
	}
	if tAdded.H != expectedH {
		t.Errorf("Hour mismatch: expected %d, got %d", expectedH, tAdded.H)
	}
	if tAdded.I != expectedI {
		t.Errorf("Minute mismatch: expected %d, got %d", expectedI, tAdded.I)
	}
	if tAdded.S != expectedS {
		t.Errorf("Second mismatch: expected %d, got %d", expectedS, tAdded.S)
	}
}

// testAddWithTimezone is a helper that tests adding an interval to a time with a timezone
func testAddWithTimezone(t *testing.T, tzid string, from string, period string, expectedY, expectedM, expectedD, expectedH, expectedI, expectedS int64) {
	t.Helper()

	// Parse timezone
	var errCode int
	tzi, err := timelib.ParseTzfile(tzid, timelib.BuiltinDB(), &errCode)
	if err != nil {
		t.Fatalf("Failed to parse timezone '%s': %v", tzid, err)
	}

	// Parse the from time
	tFrom, err := timelib.StrToTime(from, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", from, err)
	}

	// No need to FillHoles - the date strings in tests are complete

	// Update timestamp with timezone
	tFrom.UpdateTS(tzi)

	// Parse the period/interval
	var iPeriod *timelib.RelTime
	periodStr := period
	invert := false

	// Handle negative intervals
	if strings.HasPrefix(period, "-") {
		invert = true
		periodStr = period[1:]
	}

	// Parse interval
	errorsContainer := &timelib.ErrorContainer{}
	_, _, iPeriod, _, err = timelib.Strtointerval(periodStr, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", periodStr, err)
	}

	if invert {
		iPeriod.Invert = true
	}

	// Perform the addition using AddWall
	tAdded := tFrom.AddWall(iPeriod)

	// Check results
	if tAdded.Y != expectedY {
		t.Errorf("Year mismatch: expected %d, got %d", expectedY, tAdded.Y)
	}
	if tAdded.M != expectedM {
		t.Errorf("Month mismatch: expected %d, got %d", expectedM, tAdded.M)
	}
	if tAdded.D != expectedD {
		t.Errorf("Day mismatch: expected %d, got %d", expectedD, tAdded.D)
	}
	if tAdded.H != expectedH {
		t.Errorf("Hour mismatch: expected %d, got %d", expectedH, tAdded.H)
	}
	if tAdded.I != expectedI {
		t.Errorf("Minute mismatch: expected %d, got %d", expectedI, tAdded.I)
	}
	if tAdded.S != expectedS {
		t.Errorf("Second mismatch: expected %d, got %d", expectedS, tAdded.S)
	}
}

// Offset-based tests (fall time zone)
func TestTimeFallType2PrevType3Prev(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-10-04 02:18:48", "P0Y1M2DT16H19M40S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType2PrevType3Dt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "P0Y0M0DT5H31M52S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType2PrevType3Redodt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "P0Y0M0DT6H34M5S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType2PrevType3Redost(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "P0Y0M0DT7H36M16S", 2010, 11, 7, 2, 14, 44)
}

func TestTimeFallType2PrevType3St(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "P0Y0M0DT9H38M27S", 2010, 11, 7, 4, 16, 55)
}

func TestTimeFallType2PrevType3Post(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "P0Y0M2DT1H21M31S", 2010, 11, 8, 19, 59, 59)
}

func TestTimeFallType2DtType3Prev(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "-P0Y0M0DT5H31M52S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType2DtType3Dt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "P0Y0M0DT0H5M15S", 2010, 11, 7, 0, 15, 35)
}

func TestTimeFallType2DtType3Redodt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "P0Y0M0DT1H2M13S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType2DtType3Redost(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "P0Y0M0DT2H4M24S", 2010, 11, 7, 2, 14, 44)
}

func TestTimeFallType2DtType3St(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "P0Y0M0DT4H6M35S", 2010, 11, 7, 4, 16, 55)
}

func TestTimeFallType2DtType3Post(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "P0Y0M1DT20H49M39S", 2010, 11, 8, 20, 59, 59)
}

func TestTimeFallType2RedodtType3Prev(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "-P0Y0M0DT6H34M5S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType2RedodtType3Dt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "-P0Y0M0DT1H2M13S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType2RedodtType3Redodt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "P0Y0M0DT0H3M2S", 2010, 11, 7, 1, 15, 35)
}

func TestTimeFallType2RedodtType3Redost(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "P0Y0M0DT1H2M11S", 2010, 11, 7, 2, 14, 44)
}

func TestTimeFallType2RedodtType3St(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "P0Y0M0DT3H4M22S", 2010, 11, 7, 4, 16, 55)
}

func TestTimeFallType2RedodtType3Post(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "P0Y0M1DT19H47M26S", 2010, 11, 8, 20, 59, 59)
}

func TestTimeFallType2RedostType3Prev(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "-P0Y0M0DT7H36M16S", 2010, 11, 6, 17, 38, 28)
}

func TestTimeFallType2RedostType3Dt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "-P0Y0M0DT2H4M24S", 2010, 11, 6, 23, 10, 20)
}

func TestTimeFallType2RedostType3Redodt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "-P0Y0M0DT1H2M11S", 2010, 11, 7, 0, 12, 33)
}

func TestTimeFallType2RedostType3Redost(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "P0Y0M0DT1H0M0S", 2010, 11, 7, 2, 14, 44)
}

func TestTimeFallType2RedostType3St(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "P0Y0M0DT3H2M11S", 2010, 11, 7, 4, 16, 55)
}

func TestTimeFallType2RedostType3Post(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "P0Y0M1DT18H45M15S", 2010, 11, 8, 19, 59, 59)
}

func TestTimeFallType2StType3Prev(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 04:16:55", "-P0Y0M0DT9H38M27S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType2StType3Dt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 04:16:55", "-P0Y0M0DT4H6M35S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType2StType3Redodt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 04:16:55", "-P0Y0M0DT3H4M22S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType2StType3Redost(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 04:16:55", "-P0Y0M0DT2H2M11S", 2010, 11, 7, 2, 14, 44)
}

func TestTimeFallType2StType3St(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 04:16:55", "P0Y0M0DT1H3M2S", 2010, 11, 7, 5, 19, 57)
}

func TestTimeFallType2StType3Post(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 04:16:55", "P0Y0M1DT15H43M4S", 2010, 11, 8, 19, 59, 59)
}

func TestTimeFallType2PostType3Prev(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "-P0Y0M2DT1H21M31S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType2PostType3Dt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "-P0Y0M1DT19H49M39S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType2PostType3Redodt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "-P0Y0M1DT18H47M26S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType2PostType3Redost(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "-P0Y0M1DT17H45M15S", 2010, 11, 7, 2, 14, 44)
}

func TestTimeFallType2PostType3St(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "-P0Y0M1DT15H43M4S", 2010, 11, 7, 4, 16, 55)
}

func TestTimeFallType2PostType3Post(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M1DT2H3M4S", 2010, 11, 9, 22, 3, 3)
}

// Timezone-based tests (America/New_York)
func TestTimeFallType3PrevType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-10-04 02:18:48", "P0Y1M2DT16H19M40S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType3PrevType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "P0Y0M0DT5H31M52S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType3PrevType2Redodt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "P0Y0M0DT6H34M5S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType3PrevType2Redost(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "P0Y0M0DT7H36M16S", 2010, 11, 7, 1, 14, 44)
}

func TestTimeFallType3PrevType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "P0Y0M0DT8H38M27S", 2010, 11, 7, 2, 16, 55)
}

func TestTimeFallType3PrevType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "P0Y0M2DT1H21M31S", 2010, 11, 8, 19, 59, 59)
}

// Add a simple test to verify basic functionality
func TestAddBasic(t *testing.T) {
	// Parse a time
	tFrom, err := timelib.StrToTime("2010-01-01 12:00:00", nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}

	// Update timestamp to make sure SSE is valid
	tFrom.UpdateTS(nil)

	t.Logf("Before: %04d-%02d-%02d %02d:%02d:%02d SSE=%d", tFrom.Y, tFrom.M, tFrom.D, tFrom.H, tFrom.I, tFrom.S, tFrom.Sse)

	// Parse an interval of 1 day
	errorsContainer := &timelib.ErrorContainer{}
	_, _, interval, _, err := timelib.Strtointerval("P1D", errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval: %v", err)
	}

	t.Logf("Interval: Y=%d M=%d D=%d H=%d I=%d S=%d", interval.Y, interval.M, interval.D, interval.H, interval.I, interval.S)

	// Add the interval
	result := tFrom.AddWall(interval)

	t.Logf("After: %04d-%02d-%02d %02d:%02d:%02d SSE=%d", result.Y, result.M, result.D, result.H, result.I, result.S, result.Sse)

	// Check result
	if result.Y != 2010 || result.M != 1 || result.D != 2 {
		t.Errorf("Expected 2010-01-02, got %04d-%02d-%02d", result.Y, result.M, result.D)
	}
	}

// ===== Offset-based Spring Tests (type2) =====

func TestTimeSpringType2PrevType3Prev(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-02-11 02:18:48", "P0Y1M2DT16H19M40S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType2PrevType3St(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "P0Y0M0DT5H31M52S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType2PrevType3Dt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "P0Y0M0DT7H38M27S", 2010, 3, 14, 2, 16, 55)
}

func TestTimeSpringType2PrevType3Post(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "P0Y0M2DT1H21M31S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType2StType3Prev(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "-P0Y0M0DT5H31M52S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType2StType3St(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "P0Y0M0DT0H5M15S", 2010, 3, 14, 0, 15, 35)
}

func TestTimeSpringType2StType3Dt(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "P0Y0M0DT2H6M35S", 2010, 3, 14, 2, 16, 55)
}

func TestTimeSpringType2StType3Post(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "P0Y0M1DT18H49M39S", 2010, 3, 15, 18, 59, 59)
}

func TestTimeSpringType2DtType3Prev(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "-P0Y0M0DT7H38M27S", 2010, 3, 13, 19, 38, 28)
}

func TestTimeSpringType2DtType3St(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "-P0Y0M0DT2H6M35S", 2010, 3, 14, 1, 10, 20)
}

func TestTimeSpringType2DtType3Dt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "P0Y0M0DT2H3M1S", 2010, 3, 14, 5, 19, 56)
}

func TestTimeSpringType2DtType3Post(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "P0Y0M1DT16H43M4S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType2PostType3Prev(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "-P0Y0M2DT1H21M31S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType2PostType3St(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "-P0Y0M1DT18H49M39S", 2010, 3, 14, 1, 10, 20)
}

func TestTimeSpringType2PostType3Dt(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "-P0Y0M1DT16H43M4S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType2PostType3Post(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 18:57:55", "P0Y0M0DT1H2M4S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType2StsecType3Dtsec(t *testing.T) {
	testAddWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 01:59:59", "P0Y0M0DT0H0M1S", 2010, 3, 14, 2, 0, 0)
}

func TestTimeSpringType2DtsecType3Stsec(t *testing.T) {
	testAddWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:00:00", "-P0Y0M0DT0H0M1S", 2010, 3, 14, 2, 59, 59)
}

// ===== Timezone-based Fall Tests (type3) - Remaining tests =====

func TestTimeFallType3DtType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "-P0Y0M0DT5H31M52S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType3DtType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "P0Y0M0DT0H5M15S", 2010, 11, 7, 0, 15, 35)
}

func TestTimeFallType3DtType2Redodt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "P0Y0M0DT1H2M13S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType3DtType2Redost(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "P0Y0M0DT2H4M24S", 2010, 11, 7, 1, 14, 44)
}

func TestTimeFallType3DtType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "P0Y0M0DT4H6M35S", 2010, 11, 7, 3, 16, 55)
}

func TestTimeFallType3DtType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "P0Y0M1DT20H49M39S", 2010, 11, 8, 20, 59, 59)
}

func TestTimeFallType3RedodtType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "-P0Y0M0DT6H34M5S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType3RedodtType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "-P0Y0M0DT1H2M13S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType3RedodtType2Redodt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "P0Y0M0DT0H3M2S", 2010, 11, 7, 1, 15, 35)
}

func TestTimeFallType3RedodtType2Redost(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "P0Y0M0DT1H2M11S", 2010, 11, 7, 1, 14, 44)
}

func TestTimeFallType3RedodtType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "P0Y0M0DT3H4M22S", 2010, 11, 7, 3, 16, 55)
}

func TestTimeFallType3RedodtType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "P0Y0M1DT19H47M26S", 2010, 11, 8, 20, 59, 59)
}

func TestTimeFallType3RedostType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "-P0Y0M0DT7H36M16S", 2010, 11, 6, 17, 38, 28)
}

func TestTimeFallType3RedostType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "-P0Y0M0DT2H4M24S", 2010, 11, 6, 23, 10, 20)
}

func TestTimeFallType3RedostType2Redodt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "-P0Y0M0DT1H2M11S", 2010, 11, 7, 0, 12, 33)
}

func TestTimeFallType3RedostType2Redost(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT0H2M10S", 2010, 11, 7, 1, 16, 54)
}

func TestTimeFallType3RedostType2St0(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT0H2M11S", 2010, 11, 7, 1, 16, 55)
}

func TestTimeFallType3RedostType2St1(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT1H2M11S", 2010, 11, 7, 1, 16, 55)
}

func TestTimeFallType3RedostType2St2(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT2H2M11S", 2010, 11, 7, 2, 16, 55)
}

func TestTimeFallType3RedostType2St3(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT3H2M11S", 2010, 11, 7, 3, 16, 55)
}

func TestTimeFallType3RedostType2St4(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT4H2M11S", 2010, 11, 7, 4, 16, 55)
}

func TestTimeFallType3RedostType2St23(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT23H2M11S", 2010, 11, 7, 23, 16, 55)
}

func TestTimeFallType3RedostType2St24(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT24H2M11S", 2010, 11, 8, 0, 16, 55)
}

func TestTimeFallType3RedostType2St25(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT25H2M11S", 2010, 11, 8, 1, 16, 55)
}

func TestTimeFallType3RedostType2St26(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT26H2M11S", 2010, 11, 8, 2, 16, 55)
}

func TestTimeFallType3RedostType2St1_0(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M1DT0H2M11S", 2010, 11, 8, 1, 16, 55)
}

func TestTimeFallType3RedostType2St1_1(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M1DT1H2M11S", 2010, 11, 8, 2, 16, 55)
}

func TestTimeFallType3RedostType2St1_2(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M1DT2H2M11S", 2010, 11, 8, 3, 16, 55)
}

func TestTimeFallType3RedostType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M1DT18H45M15S", 2010, 11, 8, 19, 59, 59)
}

func TestTimeFallType3StType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "-P0Y0M0DT9H38M27S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType3StType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "-P0Y0M0DT4H6M35S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType3StType2Redodt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "-P0Y0M0DT3H4M22S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType3StType2Redost(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "-P0Y0M0DT2H2M11S", 2010, 11, 7, 1, 14, 44)
}

func TestTimeFallType3StType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "P0Y0M0DT2H3M1S", 2010, 11, 7, 5, 19, 56)
}

func TestTimeFallType3StType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "P0Y0M1DT16H43M4S", 2010, 11, 8, 19, 59, 59)
}

func TestTimeFallType3PostType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "-P0Y0M2DT1H21M31S", 2010, 11, 6, 18, 38, 28)
}

func TestTimeFallType3PostType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "-P0Y0M1DT20H49M39S", 2010, 11, 7, 0, 10, 20)
}

func TestTimeFallType3PostType2Redodt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "-P0Y0M1DT19H47M26S", 2010, 11, 7, 1, 12, 33)
}

func TestTimeFallType3PostType2Redost(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "-P0Y0M1DT18H45M15S", 2010, 11, 7, 1, 14, 44)
}

func TestTimeFallType3PostType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "-P0Y0M1DT16H43M4S", 2010, 11, 7, 3, 16, 55)
}

func TestTimeFallType3PostType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-08 18:57:55", "P0Y0M0DT1H2M4S", 2010, 11, 8, 19, 59, 59)
}

func TestTimeFallType3DtsecType2Stsec(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:59:59", "P0Y0M0DT0H0M1S", 2010, 11, 7, 1, 0, 0)
}

func TestTimeFallType3StsecType2Dtsec(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-11-07 01:00:00", "-P0Y0M0DT0H0M1S", 2010, 11, 7, 0, 59, 59)
}

// ===== Timezone-based Spring Tests type3_prev_type2 =====

func TestTimeSpringType3PrevType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-02-11 02:18:48", "P0Y1M2DT16H19M40S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3PrevType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y0M0DT5H31M52S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType3PrevType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y0M0DT7H38M27S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType3PrevType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y0M2DT1H21M31S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType3StType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "-P0Y0M0DT5H31M52S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3StType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M0DT0H5M15S", 2010, 3, 14, 0, 15, 35)
}

func TestTimeSpringType3StType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M0DT2H6M35S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType3StType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M1DT18H49M39S", 2010, 3, 15, 18, 59, 59)
}

func TestTimeSpringType3DtType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "-P0Y0M0DT7H38M27S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3DtType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "-P0Y0M0DT2H6M35S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType3DtType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "P0Y0M0DT2H3M1S", 2010, 3, 14, 5, 19, 56)
}

func TestTimeSpringType3DtType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "P0Y0M1DT16H43M4S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType3PostType2Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "-P0Y0M2DT1H21M31S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3PostType2St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "-P0Y0M1DT18H49M39S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType3PostType2Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "-P0Y0M1DT16H43M4S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType3PostType2Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 18:57:55", "P0Y0M0DT1H2M4S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType3StsecType2Dtsec(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 01:59:59", "P0Y0M0DT0H0M1S", 2010, 3, 14, 3, 0, 0)
}

func TestTimeSpringType3DtsecType2Stsec(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:00:00", "-P0Y0M0DT0H0M1S", 2010, 3, 14, 1, 59, 59)
}

// ===== Timezone-based Spring Tests type3_prev_type3 =====

func TestTimeSpringType3PrevType3Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-02-11 02:18:48", "P0Y1M2DT16H19M40S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3PrevType3St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y0M0DT5H31M52S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType3PrevType3Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y0M0DT7H38M27S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType3PrevType3Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y0M2DT1H21M31S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType3StType3Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "-P0Y0M0DT5H31M52S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3StType3St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M0DT0H5M15S", 2010, 3, 14, 0, 15, 35)
}

func TestTimeSpringType3StType3Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M0DT2H6M35S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType3StType3Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M1DT18H49M39S", 2010, 3, 15, 18, 59, 59)
}

func TestTimeSpringType3DtType3Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "-P0Y0M0DT7H38M27S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3DtType3St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "-P0Y0M0DT2H6M35S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType3DtType3Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "P0Y0M0DT2H3M1S", 2010, 3, 14, 5, 19, 56)
}

func TestTimeSpringType3DtType3Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "P0Y0M1DT16H43M4S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType3PostType3Prev(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "-P0Y0M2DT1H21M31S", 2010, 3, 13, 18, 38, 28)
}

func TestTimeSpringType3PostType3St(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "-P0Y0M1DT18H49M39S", 2010, 3, 14, 0, 10, 20)
}

func TestTimeSpringType3PostType3Dt(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "-P0Y0M1DT16H43M4S", 2010, 3, 14, 3, 16, 55)
}

func TestTimeSpringType3PostType3Post(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-15 18:57:55", "P0Y0M0DT1H2M4S", 2010, 3, 15, 19, 59, 59)
}

func TestTimeSpringType3StsecType3Dtsec(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 01:59:59", "P0Y0M0DT0H0M1S", 2010, 3, 14, 3, 0, 0)
}

func TestTimeSpringType3DtsecType3Stsec(t *testing.T) {
	testAddWithTimezone(t, "America/New_York", "2010-03-14 03:00:00", "-P0Y0M0DT0H0M1S", 2010, 3, 14, 1, 59, 59)
}
