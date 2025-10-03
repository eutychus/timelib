package tests

import (
	"strings"
	"testing"

	timelib "github.com/eutychus/timelib"
)

// testSubWithOffset tests time subtraction with a fixed offset timezone
// Corresponds to C++ test_sub(offset, from, period) helper function
func testSubWithOffset(t *testing.T, offset int32, from string, period string) *timelib.Time {
	t.Helper()

	tFrom, err := timelib.StrToTime(from, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", from, err)
	}

	// No need to FillHoles - the date strings in tests are complete
	tFrom.ZoneType = timelib.TIMELIB_ZONETYPE_OFFSET
	tFrom.Z = offset
	tFrom.Dst = 0  // Initialize Dst when setting timezone offset

	tFrom.UpdateTS(nil)

	// Parse interval
	var iPeriod *timelib.RelTime
	periodStr := period
	invert := false

	if strings.HasPrefix(period, "-") {
		invert = true
		periodStr = period[1:]
	}

	errorsContainer := &timelib.ErrorContainer{}
	_, _, iPeriod, _, err = timelib.Strtointerval(periodStr, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", periodStr, err)
	}

	if invert {
		iPeriod.Invert = true
	}

	tSubbed := tFrom.SubWall(iPeriod)
	return tSubbed
}

// testSubWithTimezone tests time subtraction with a timezone
// Corresponds to C++ test_sub(tzid, from, period) helper function
func testSubWithTimezone(t *testing.T, tzid string, from string, period string) *timelib.Time {
	t.Helper()

	var errCode int
	tzi, err := timelib.ParseTzfile(tzid, timelib.BuiltinDB(), &errCode)
	if err != nil {
		t.Fatalf("Failed to parse timezone '%s': %v", tzid, err)
	}

	tFrom, err := timelib.StrToTime(from, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", from, err)
	}

	// No need to FillHoles - the date strings in tests are complete
	tFrom.UpdateTS(tzi)

	// Parse interval
	var iPeriod *timelib.RelTime
	periodStr := period
	invert := false

	if strings.HasPrefix(period, "-") {
		invert = true
		periodStr = period[1:]
	}

	errorsContainer := &timelib.ErrorContainer{}
	_, _, iPeriod, _, err2 := timelib.Strtointerval(periodStr, errorsContainer)
	if err2 != nil {
		t.Fatalf("Failed to parse interval '%s': %v", periodStr, err2)
	}

	if invert {
		iPeriod.Invert = true
	}

	tSubbed := tFrom.SubWall(iPeriod)
	return tSubbed
}

// checkResult verifies the subtraction result matches expected values
func checkResult(t *testing.T, result *timelib.Time, y, m, d, h, i, s int) {
	t.Helper()
	if result.Y != int64(y) {
		t.Errorf("Expected year=%d, got %d", y, result.Y)
	}
	if result.M != int64(m) {
		t.Errorf("Expected month=%d, got %d", m, result.M)
	}
	if result.D != int64(d) {
		t.Errorf("Expected day=%d, got %d", d, result.D)
	}
	if result.H != int64(h) {
		t.Errorf("Expected hour=%d, got %d", h, result.H)
	}
	if result.I != int64(i) {
		t.Errorf("Expected minute=%d, got %d", i, result.I)
	}
	if result.S != int64(s) {
		t.Errorf("Expected second=%d, got %d", s, result.S)
	}
}

// SECS_PER_HOUR constant for offset calculations

// ============================================================================
// DST Fall Tests with Fixed Offset (Type2)
// Tests around 2010-11-07 DST transition in America/New_York
// ============================================================================

func TestTimeSubFallType2PrevType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "P0Y1M2DT16H19M40S")
	checkResult(t, result, 2010, 10, 4, 2, 18, 48)
}

func TestTimeSubFallType2PrevType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType2PrevType3Redodt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "P0Y0M0DT6H34M5S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType2PrevType3Redost(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "P0Y0M0DT7H36M16S")
	checkResult(t, result, 2010, 11, 6, 17, 38, 28)
}

func TestTimeSubFallType2PrevType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 03:16:55", "P0Y0M0DT9H38M27S")
	checkResult(t, result, 2010, 11, 6, 17, 38, 28)
}

func TestTimeSubFallType2PrevType3Post(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType2DtType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "-P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType2DtType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:15:35", "P0Y0M0DT0H5M15S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType2DtType3Redodt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "P0Y0M0DT1H2M13S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType2DtType3Redost(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "P0Y0M0DT2H4M24S")
	checkResult(t, result, 2010, 11, 6, 23, 10, 20)
}

func TestTimeSubFallType2DtType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 03:16:55", "P0Y0M0DT4H6M35S")
	checkResult(t, result, 2010, 11, 6, 23, 10, 20)
}

func TestTimeSubFallType2DtType3Post(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M1DT20H49M39S")
	checkResult(t, result, 2010, 11, 6, 23, 10, 20)
}

func TestTimeSubFallType2RedodtType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "-P0Y0M0DT6H34M5S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType2RedodtType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "-P0Y0M0DT1H2M13S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType2RedodtType3Redodt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:15:35", "P0Y0M0DT0H3M2S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType2RedodtType3Redost(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "P0Y0M0DT1H2M11S")
	checkResult(t, result, 2010, 11, 7, 0, 12, 33)
}

func TestTimeSubFallType2RedodtType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 03:16:55", "P0Y0M0DT3H4M22S")
	checkResult(t, result, 2010, 11, 7, 0, 12, 33)
}

func TestTimeSubFallType2RedodtType3Post(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M1DT19H47M26S")
	checkResult(t, result, 2010, 11, 7, 0, 12, 33)
}

func TestTimeSubFallType2RedostType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "-P0Y0M0DT7H36M16S")
	checkResult(t, result, 2010, 11, 7, 2, 14, 44)
}

func TestTimeSubFallType2RedostType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "-P0Y0M0DT2H4M24S")
	checkResult(t, result, 2010, 11, 7, 2, 14, 44)
}

func TestTimeSubFallType2RedostType3Redodt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "-P0Y0M0DT1H2M11S")
	checkResult(t, result, 2010, 11, 7, 2, 14, 44)
}

func TestTimeSubFallType2RedostType3Redost(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:16:54", "P0Y0M0DT0H2M10S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType2RedostType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 03:16:55", "P0Y0M0DT2H2M11S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType2RedostType3Post(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M1DT18H45M15S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType2StType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "-P0Y0M0DT9H38M27S")
	checkResult(t, result, 2010, 11, 7, 4, 16, 55)
}

func TestTimeSubFallType2StType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "-P0Y0M0DT4H6M35S")
	checkResult(t, result, 2010, 11, 7, 4, 16, 55)
}

func TestTimeSubFallType2StType3Redodt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "-P0Y0M0DT3H4M22S")
	checkResult(t, result, 2010, 11, 7, 4, 16, 55)
}

func TestTimeSubFallType2StType3Redost(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "-P0Y0M0DT2H2M11S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType2StType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 05:19:56", "P0Y0M0DT2H3M1S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType2StType3Post(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType2PostType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-06 18:38:28", "-P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 11, 8, 19, 59, 59)
}

func TestTimeSubFallType2PostType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 00:10:20", "-P0Y0M1DT20H49M39S")
	checkResult(t, result, 2010, 11, 8, 20, 59, 59)
}

func TestTimeSubFallType2PostType3Redodt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:12:33", "-P0Y0M1DT19H47M26S")
	checkResult(t, result, 2010, 11, 8, 20, 59, 59)
}

func TestTimeSubFallType2PostType3Redost(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:14:44", "-P0Y0M1DT18H45M15S")
	checkResult(t, result, 2010, 11, 8, 19, 59, 59)
}

func TestTimeSubFallType2PostType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 03:16:55", "-P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 11, 8, 19, 59, 59)
}

func TestTimeSubFallType2PostType3Post(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-08 19:59:59", "P0Y0M0DT1H2M4S")
	checkResult(t, result, 2010, 11, 8, 18, 57, 55)
}

func TestTimeSubFallType2DtsecType3Stsec(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-11-07 01:00:00", "P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 11, 7, 0, 59, 59)
}

func TestTimeSubFallType2StsecType3Dtsec(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-11-07 01:59:59", "-P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 11, 7, 2, 0, 0)
}

// ============================================================================
// DST Fall Tests with Timezone (Type3)
// Tests around 2010-11-07 DST transition in America/New_York
// ============================================================================

func TestTimeSubFallType3PrevType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "P0Y1M2DT16H19M40S")
	checkResult(t, result, 2010, 10, 4, 2, 18, 48)
}

func TestTimeSubFallType3PrevType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType3PrevType3Redodt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "P0Y0M0DT6H34M5S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType3PrevType3Redost(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT7H36M16S")
	checkResult(t, result, 2010, 11, 6, 17, 38, 28)
}

func TestTimeSubFallType3PrevType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "P0Y0M0DT9H38M27S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType3PrevType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 11, 6, 18, 38, 28)
}

func TestTimeSubFallType3DtType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "-P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType3DtType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 00:15:35", "P0Y0M0DT0H5M15S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType3DtType3Redodt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "P0Y0M0DT1H2M13S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType3DtType3Redost(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT2H4M24S")
	checkResult(t, result, 2010, 11, 6, 23, 10, 20)
}

func TestTimeSubFallType3DtType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "P0Y0M0DT4H6M35S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType3DtType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "P0Y0M1DT20H49M39S")
	checkResult(t, result, 2010, 11, 7, 0, 10, 20)
}

func TestTimeSubFallType3RedodtType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "-P0Y0M0DT6H34M5S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType3RedodtType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "-P0Y0M0DT1H2M13S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType3RedodtType3Redodt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:15:35", "P0Y0M0DT0H3M2S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType3RedodtType3Redost(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "P0Y0M0DT1H2M11S")
	checkResult(t, result, 2010, 11, 7, 0, 12, 33)
}

func TestTimeSubFallType3RedodtType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "P0Y0M0DT3H4M22S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType3RedodtType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "P0Y0M1DT19H47M26S")
	checkResult(t, result, 2010, 11, 7, 1, 12, 33)
}

func TestTimeSubFallType3RedostType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "-P0Y0M0DT7H36M16S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType3RedostType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "-P0Y0M0DT2H4M24S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType3RedostType3Redodt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "-P0Y0M0DT1H2M11S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType3RedostType3Redost(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:16:54", "P0Y0M0DT0H2M10S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType3RedostType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "P0Y0M0DT2H2M11S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType3RedostType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "P0Y0M1DT18H45M15S")
	checkResult(t, result, 2010, 11, 7, 1, 14, 44)
}

func TestTimeSubFallType3StType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "-P0Y0M0DT9H38M27S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType3StType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "-P0Y0M0DT4H6M35S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType3StType3Redodt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "-P0Y0M0DT3H4M22S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType3StType3Redost(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "-P0Y0M0DT2H2M11S")
	checkResult(t, result, 2010, 11, 7, 2, 16, 55)
}

func TestTimeSubFallType3StType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 05:19:56", "P0Y0M0DT2H3M1S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType3StType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 11, 7, 3, 16, 55)
}

func TestTimeSubFallType3PostType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-06 18:38:28", "-P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 11, 8, 19, 59, 59)
}

func TestTimeSubFallType3PostType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 00:10:20", "-P0Y0M1DT20H49M39S")
	checkResult(t, result, 2010, 11, 8, 20, 59, 59)
}

func TestTimeSubFallType3PostType3Redodt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:12:33", "-P0Y0M1DT19H47M26S")
	checkResult(t, result, 2010, 11, 8, 20, 59, 59)
}

func TestTimeSubFallType3PostType3Redost(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:14:44", "-P0Y0M1DT18H45M15S")
	checkResult(t, result, 2010, 11, 8, 19, 59, 59)
}

func TestTimeSubFallType3PostType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 03:16:55", "-P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 11, 8, 19, 59, 59)
}

func TestTimeSubFallType3PostType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-08 19:59:59", "P0Y0M0DT1H2M4S")
	checkResult(t, result, 2010, 11, 8, 18, 57, 55)
}

func TestTimeSubFallType3DtsecType3Stsec(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:00:00", "P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 11, 7, 0, 59, 59)
}

func TestTimeSubFallType3StsecType3Dtsec(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-11-07 01:59:59", "-P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 11, 7, 1, 0, 0)
}

// ============================================================================
// DST Spring Tests with Fixed Offset (Type2)
// Tests around 2010-03-14 DST transition in America/New_York
// ============================================================================

func TestTimeSubSpringType2PrevType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "P0Y1M2DT16H19M40S")
	checkResult(t, result, 2010, 2, 11, 2, 18, 48)
}

func TestTimeSubSpringType2PrevType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 3, 13, 18, 38, 28)
}

func TestTimeSubSpringType2PrevType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "P0Y0M0DT7H38M27S")
	checkResult(t, result, 2010, 3, 13, 19, 38, 28)
}

func TestTimeSubSpringType2PrevType3Post(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 3, 13, 18, 38, 28)
}

func TestTimeSubSpringType2StType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "-P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 3, 14, 0, 10, 20)
}

func TestTimeSubSpringType2StType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:15:35", "P0Y0M0DT0H5M15S")
	checkResult(t, result, 2010, 3, 14, 0, 10, 20)
}

func TestTimeSubSpringType2StType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "P0Y0M0DT2H6M35S")
	checkResult(t, result, 2010, 3, 14, 1, 10, 20)
}

func TestTimeSubSpringType2StType3Post(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "P0Y0M1DT18H49M39S")
	checkResult(t, result, 2010, 3, 14, 1, 10, 20)
}

func TestTimeSubSpringType2DtType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "-P0Y0M0DT7H38M27S")
	checkResult(t, result, 2010, 3, 14, 2, 16, 55)
}

func TestTimeSubSpringType2DtType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "-P0Y0M0DT2H6M35S")
	checkResult(t, result, 2010, 3, 14, 2, 16, 55)
}

func TestTimeSubSpringType2DtType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 05:19:56", "P0Y0M0DT2H3M1S")
	checkResult(t, result, 2010, 3, 14, 3, 16, 55)
}

func TestTimeSubSpringType2DtType3Post(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 3, 14, 3, 16, 55)
}

func TestTimeSubSpringType2PostType3Prev(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-13 18:38:28", "-P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 3, 15, 19, 59, 59)
}

func TestTimeSubSpringType2PostType3St(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 00:10:20", "-P0Y0M1DT18H49M39S")
	checkResult(t, result, 2010, 3, 15, 18, 59, 59)
}

func TestTimeSubSpringType2PostType3Dt(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:16:55", "-P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 3, 15, 19, 59, 59)
}

func TestTimeSubSpringType2PostType3Post(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-15 19:59:59", "P0Y0M0DT1H2M4S")
	checkResult(t, result, 2010, 3, 15, 18, 57, 55)
}

func TestTimeSubSpringType2StsecType3Dtsec(t *testing.T) {
	result := testSubWithOffset(t, -4*SECS_PER_HOUR, "2010-03-14 03:00:00", "P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 3, 14, 2, 59, 59)
}

func TestTimeSubSpringType2DtsecType3Stsec(t *testing.T) {
	result := testSubWithOffset(t, -5*SECS_PER_HOUR, "2010-03-14 01:59:59", "-P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 3, 14, 2, 0, 0)
}

// ============================================================================
// DST Spring Tests with Timezone (Type3)
// Tests around 2010-03-14 DST transition in America/New_York
// ============================================================================

func TestTimeSubSpringType3PrevType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "P0Y1M2DT16H19M40S")
	checkResult(t, result, 2010, 2, 11, 2, 18, 48)
}

func TestTimeSubSpringType3PrevType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 3, 13, 18, 38, 28)
}

func TestTimeSubSpringType3PrevType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "P0Y0M0DT7H38M27S")
	checkResult(t, result, 2010, 3, 13, 18, 38, 28)
}

func TestTimeSubSpringType3PrevType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 3, 13, 18, 38, 28)
}

func TestTimeSubSpringType3StType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "-P0Y0M0DT5H31M52S")
	checkResult(t, result, 2010, 3, 14, 0, 10, 20)
}

func TestTimeSubSpringType3StType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 00:15:35", "P0Y0M0DT0H5M15S")
	checkResult(t, result, 2010, 3, 14, 0, 10, 20)
}

func TestTimeSubSpringType3StType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "P0Y0M0DT2H6M35S")
	checkResult(t, result, 2010, 3, 14, 0, 10, 20)
}

func TestTimeSubSpringType3StType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "P0Y0M1DT18H49M39S")
	checkResult(t, result, 2010, 3, 14, 0, 10, 20)
}

func TestTimeSubSpringType3DtType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "-P0Y0M0DT7H38M27S")
	checkResult(t, result, 2010, 3, 14, 3, 16, 55)
}

func TestTimeSubSpringType3DtType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "-P0Y0M0DT2H6M35S")
	checkResult(t, result, 2010, 3, 14, 3, 16, 55)
}

func TestTimeSubSpringType3DtType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 05:19:56", "P0Y0M0DT2H3M1S")
	checkResult(t, result, 2010, 3, 14, 3, 16, 55)
}

func TestTimeSubSpringType3DtType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 3, 14, 3, 16, 55)
}

func TestTimeSubSpringType3PostType3Prev(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-13 18:38:28", "-P0Y0M2DT1H21M31S")
	checkResult(t, result, 2010, 3, 15, 19, 59, 59)
}

func TestTimeSubSpringType3PostType3St(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 00:10:20", "-P0Y0M1DT18H49M39S")
	checkResult(t, result, 2010, 3, 15, 18, 59, 59)
}

func TestTimeSubSpringType3PostType3Dt(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 03:16:55", "-P0Y0M1DT16H43M4S")
	checkResult(t, result, 2010, 3, 15, 19, 59, 59)
}

func TestTimeSubSpringType3PostType3Post(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-15 19:59:59", "P0Y0M0DT1H2M4S")
	checkResult(t, result, 2010, 3, 15, 18, 57, 55)
}

func TestTimeSubSpringType3StsecType3Dtsec(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 03:00:00", "P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 3, 14, 1, 59, 59)
}

func TestTimeSubSpringType3DtsecType3Stsec(t *testing.T) {
	result := testSubWithTimezone(t, "America/New_York", "2010-03-14 01:59:59", "-P0Y0M0DT0H0M1S")
	checkResult(t, result, 2010, 3, 14, 3, 0, 0)
}
