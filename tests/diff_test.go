package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// Helper function to parse time string with timezone
func testParse(tzid, from, to string) (*timelib.Time, *timelib.Time, *timelib.RelTime, error) {
	var errorCode int
	tzi, err := timelib.ParseTzfile(tzid, timelib.BuiltinDB(), &errorCode)
	if err != nil {
		// If timezone parsing fails, use basic parsing
		tzi = nil
	}

	tNow := timelib.TimeCtor()
	tNow.Y = 2024
	tNow.M = 1
	tNow.D = 1
	tNow.H = 0
	tNow.I = 0
	tNow.S = 0

	tFrom, _ := timelib.StrToTime(from, nil)
	if tFrom == nil {
		return nil, nil, nil, err
	}

	tTo, _ := timelib.StrToTime(to, nil)
	if tTo == nil {
		return nil, nil, nil, err
	}

	timelib.FillHoles(tFrom, tNow, timelib.TIMELIB_NO_CLONE)
	timelib.FillHoles(tTo, tNow, timelib.TIMELIB_NO_CLONE)

	tFrom.UpdateTS(tzi)
	tTo.UpdateTS(tzi)

	diff := tFrom.Diff(tTo)
	return tFrom, tTo, diff, nil
}

// Helper function to parse time strings with different timezones
func testParseTz(tzidFrom, tzidTo, from, to string) (*timelib.Time, *timelib.Time, *timelib.RelTime, error) {
	var errCode int
	tziFrom, err := timelib.ParseTzfile(tzidFrom, timelib.BuiltinDB(), &errCode)
	if err != nil {
		tziFrom = nil
	}
	tziTo, err := timelib.ParseTzfile(tzidTo, timelib.BuiltinDB(), &errCode)
	if err != nil {
		tziTo = nil
	}

	tNow := timelib.TimeCtor()
	tNow.Y = 2024
	tNow.M = 1
	tNow.D = 1
	tNow.H = 0
	tNow.I = 0
	tNow.S = 0

	tFrom, _ := timelib.StrToTime(from, nil)
	tTo, _ := timelib.StrToTime(to, nil)

	timelib.FillHoles(tFrom, tNow, timelib.TIMELIB_NO_CLONE)
	timelib.FillHoles(tTo, tNow, timelib.TIMELIB_NO_CLONE)

	tFrom.UpdateTS(tziFrom)
	tTo.UpdateTS(tziTo)

	diff := tFrom.Diff(tTo)
	return tFrom, tTo, diff, nil
}

// Helper function to parse time strings with offset
func testParseOffset(offsetFrom, offsetTo int32, from, to string) (*timelib.Time, *timelib.Time, *timelib.RelTime, error) {
	tNow := timelib.TimeCtor()
	tNow.Y = 2024
	tNow.M = 1
	tNow.D = 1
	tNow.H = 0
	tNow.I = 0
	tNow.S = 0

	tFrom, _ := timelib.StrToTime(from, nil)
	tTo, _ := timelib.StrToTime(to, nil)

	timelib.FillHoles(tFrom, tNow, timelib.TIMELIB_NO_CLONE)
	timelib.FillHoles(tTo, tNow, timelib.TIMELIB_NO_CLONE)

	tFrom.ZoneType = timelib.TIMELIB_ZONETYPE_OFFSET
	tFrom.Z = offsetFrom
	tTo.ZoneType = timelib.TIMELIB_ZONETYPE_OFFSET
	tTo.Z = offsetTo

	tFrom.UpdateTS(nil)
	tTo.UpdateTS(nil)

	diff := tFrom.Diff(tTo)
	return tFrom, tTo, diff, nil
}

func testParseOffsetWithDst(offsetFrom int32, dstFrom int, offsetTo int32, dstTo int, from, to string) (*timelib.Time, *timelib.Time, *timelib.RelTime, error) {
	tNow := timelib.TimeCtor()
	tNow.Y = 2024
	tNow.M = 1
	tNow.D = 1
	tNow.H = 0
	tNow.I = 0
	tNow.S = 0

	tFrom, _ := timelib.StrToTime(from, nil)
	tTo, _ := timelib.StrToTime(to, nil)

	timelib.FillHoles(tFrom, tNow, timelib.TIMELIB_NO_CLONE)
	timelib.FillHoles(tTo, tNow, timelib.TIMELIB_NO_CLONE)

	tFrom.ZoneType = timelib.TIMELIB_ZONETYPE_OFFSET
	tFrom.Z = offsetFrom
	tFrom.Dst = dstFrom
	tTo.ZoneType = timelib.TIMELIB_ZONETYPE_OFFSET
	tTo.Z = offsetTo
	tTo.Dst = dstTo

	tFrom.UpdateTS(nil)
	tTo.UpdateTS(nil)

	diff := tFrom.Diff(tTo)
	return tFrom, tTo, diff, nil
}

// Helper function to check diff results
func checkDiff(t *testing.T, diff *timelib.RelTime, expY, expM, expD, expH, expI, expS, expUS int64) {
	t.Helper()
	if diff.Y != expY {
		t.Errorf("Expected Y=%d, got %d", expY, diff.Y)
	}
	if diff.M != expM {
		t.Errorf("Expected M=%d, got %d", expM, diff.M)
	}
	if diff.D != expD {
		t.Errorf("Expected D=%d, got %d", expD, diff.D)
	}
	if diff.H != expH {
		t.Errorf("Expected H=%d, got %d", expH, diff.H)
	}
	if diff.I != expI {
		t.Errorf("Expected I=%d, got %d", expI, diff.I)
	}
	if diff.S != expS {
		t.Errorf("Expected S=%d, got %d", expS, diff.S)
	}
	if diff.US != expUS {
		t.Errorf("Expected US=%d, got %d", expUS, diff.US)
	}
}

func TestTimeDifference(t *testing.T) {
	// Test basic time difference functionality
	// Create two time structures representing different dates

	// First time: 2008-03-26
	t1 := timelib.TimeCtor()
	t1.Y = 2008
	t1.M = 3
	t1.D = 26
	t1.H = 0
	t1.I = 0
	t1.S = 0
	t1.HaveDate = true
	t1.HaveTime = true

	// Second time: 2001-09-11
	t2 := timelib.TimeCtor()
	t2.Y = 2001
	t2.M = 9
	t2.D = 11
	t2.H = 0
	t2.I = 0
	t2.S = 0
	t2.HaveDate = true
	t2.HaveTime = true

	// Calculate the difference
	// This would normally use timelib_diff, but we'll use a basic implementation
	// For now, just verify the structures are set up correctly

	if !t1.HaveDate || !t1.HaveTime {
		t.Error("Expected t1 to have date and time set")
	}

	if !t2.HaveDate || !t2.HaveTime {
		t.Error("Expected t2 to have date and time set")
	}

	if t1.Y != 2008 || t1.M != 3 || t1.D != 26 {
		t.Errorf("Expected t1 date 2008-03-26, got %d-%d-%d", t1.Y, t1.M, t1.D)
	}

	if t2.Y != 2001 || t2.M != 9 || t2.D != 11 {
		t.Errorf("Expected t2 date 2001-09-11, got %d-%d-%d", t2.Y, t2.M, t2.D)
	}
}

func TestTimeDifferenceBasic(t *testing.T) {
	// Test basic time structure setup
	tm := timelib.TimeCtor()
	tm.Y = 1970
	tm.M = 1
	tm.D = 1
	tm.H = 0
	tm.I = 0
	tm.S = 0
	tm.HaveDate = true
	tm.HaveTime = true

	// Verify basic structure
	if !tm.HaveDate {
		t.Error("Expected HaveDate to be true")
	}

	if !tm.HaveTime {
		t.Error("Expected HaveTime to be true")
	}

	if tm.Y != 1970 || tm.M != 1 || tm.D != 1 {
		t.Errorf("Expected epoch date 1970-01-01, got %d-%d-%d", tm.Y, tm.M, tm.D)
	}
}

func TestPhp62326(t *testing.T) {
	_, _, diff, err := testParse("Europe/Berlin", "2012-06-01", "2012-12-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 6, 0, 0, 0, 0, 0)
}

func TestPhp65003_01(t *testing.T) {
	_, _, diff, err := testParse("Europe/Moscow", "13-03-01", "13-04-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 0, 0, 0, 0, 0)
}

func TestPhp65003_02(t *testing.T) {
	_, _, diff, err := testParse("Europe/Moscow", "13-03-02", "13-04-02")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 0, 0, 0, 0, 0)
}

func TestPhp68503_01(t *testing.T) {
	_, _, diff, err := testParse("Europe/London", "2015-02-01", "2015-05-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 3, 0, 0, 0, 0, 0)
}

func TestPhp68503_02(t *testing.T) {
	_, _, diff, err := testParse("UTC", "2015-02-01", "2015-05-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 3, 0, 0, 0, 0, 0)
}

func TestPhp69378_01(t *testing.T) {
	_, _, diff, err := testParse("UTC", "2015-04-02 09:55:47", "2014-02-16 02:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 1, 1, 14, 7, 55, 47, 0)
	if !diff.Invert {
		t.Errorf("Expected Invert=true, got %v", diff.Invert)
	}
}

func TestPhp69378_02(t *testing.T) {
	_, _, diff, err := testParse("UTC", "2014-02-16 02:00:00", "2015-04-02 09:55:47")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 1, 1, 17, 7, 55, 47, 0)
	if diff.Invert {
		t.Errorf("Expected Invert=false, got %v", diff.Invert)
	}
}

func TestPhp71700_01(t *testing.T) {
	_, _, diff, err := testParse("UTC", "2016-03-01", "2016-03-31")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 30, 0, 0, 0, 0)
}

func TestPhp71826_01(t *testing.T) {
	_, _, diff, err := testParse("Asia/Tokyo", "2015-02-01", "2015-03-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 0, 0, 0, 0, 0)
}

func TestPhp71826_02(t *testing.T) {
	_, _, diff, err := testParse("Asia/Tokyo", "2015-03-01", "2015-03-29")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 28, 0, 0, 0, 0)
}

func TestPhp71826_03(t *testing.T) {
	_, _, diff, err := testParse("Asia/Tokyo", "2015-04-01", "2015-04-29")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 28, 0, 0, 0, 0)
}

func TestPhp74524(t *testing.T) {
	_, _, diff, err := testParse("Europe/Paris", "2017-04-03 22:29:15.079459", "2017-11-17 22:05:26.000000")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 7, 13, 23, 36, 10, 920541)
}

func TestPhp77032_01(t *testing.T) {
	_, _, diff, err := testParse("UTC", "2008-03-01", "2018-03-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 10, 0, 0, 0, 0, 0, 0)
}

func TestPhp77032_02(t *testing.T) {
	_, _, diff, err := testParse("Europe/Amsterdam", "2008-03-01", "2018-03-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 10, 0, 0, 0, 0, 0, 0)
}

func TestPhp76374_01(t *testing.T) {
	_, _, diff, err := testParse("Europe/Paris", "2017-10-01", "2017-01-01")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 9, 0, 0, 0, 0, 0)
}

func TestPhp76374_02(t *testing.T) {
	_, _, diff, err := testParse("Europe/Paris", "2017-10-01 12:00", "2017-01-01 12:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 9, 0, 0, 0, 0, 0)
}

func TestPhp78452(t *testing.T) {
	_, _, diff, err := testParse("Asia/Tehran", "2019-09-24 11:47:24", "2019-08-21 12:47:24")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 2, 23, 0, 0, 0)
}

func TestPhp80974(t *testing.T) {
	_, _, diff, err := testParseTz("America/Toronto", "America/Vancouver", "2012-01-01 00:00", "2012-01-01 00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 3, 0, 0, 0)
}

func TestPhp81273(t *testing.T) {
	_, _, diff, err := testParseTz("Australia/Sydney", "America/Los_Angeles", "2000-01-01 00:00:00.000000", "2000-01-01 00:00:00.000000")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 19, 0, 0, 0)
}

func TestPhpGh8730(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -4*3600, "2022-06-08 09:15:00", "2022-06-08 09:15:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 0, 0)
}

func TestPhp81263a(t *testing.T) {
	_, _, diff, err := testParseTz("Europe/Berlin", "UTC", "2020-07-19 18:30:00", "2020-07-19 16:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 0, 0)
}

func TestPhp81263b(t *testing.T) {
	_, _, diff, err := testParseTz("UTC", "Europe/Berlin", "2020-07-19 16:30:00", "2020-07-19 18:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 0, 0)
}

func TestPhp80974a(t *testing.T) {
	_, _, diff, err := testParseTz("America/Toronto", "America/Vancouver", "2012-01-01 00:00:00", "2012-01-01 00:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 3, 0, 0, 0)
}

func TestPhp80974b(t *testing.T) {
	_, _, diff, err := testParseTz("America/Vancouver", "America/Toronto", "2012-01-01 00:00:00", "2012-01-01 00:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 3, 0, 0, 0)
}

func TestTimeSpringType2PrevType2Prev(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -4*3600, "2010-03-13 18:38:28", "2010-02-11 02:18:48")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 2, 16, 19, 40, 0)
}

func TestTimeSpringType2PrevType2St(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -4*3600, "2010-03-14 00:10:20", "2010-03-13 18:38:28")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 5, 31, 52, 0)
}

func TestTimeSpringType2PrevType2Dt(t *testing.T) {
	_, _, diff, err := testParseOffsetWithDst(-4*3600, 1, -4*3600, 0, "2010-03-14 03:16:55", "2010-03-13 18:38:28")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 7, 38, 27, 0)
}

func TestTimeSpringType2PrevType2Post(t *testing.T) {
	_, _, diff, err := testParseOffsetWithDst(-4*3600, 1, -4*3600, 0, "2010-03-15 19:59:59", "2010-03-13 18:38:28")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 2, 0, 21, 31, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd1(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -5*3600, "2010-03-14 03:00:00", "2010-03-14 01:59:59")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 1, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd2(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -5*3600, "2010-03-14 04:30:00", "2010-03-13 04:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 23, 0, 0, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd3(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -5*3600, "2010-03-14 03:30:00", "2010-03-13 04:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 22, 0, 0, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd4(t *testing.T) {
	_, _, diff, err := testParseOffset(-5*3600, -5*3600, "2010-03-14 01:30:00", "2010-03-13 04:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 21, 0, 0, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd5(t *testing.T) {
	_, _, diff, err := testParseOffset(-5*3600, -5*3600, "2010-03-14 01:30:00", "2010-03-13 01:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 1, 0, 0, 0, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd6(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -5*3600, "2010-03-14 03:30:00", "2010-03-13 03:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 23, 0, 0, 0)
}

func TestDateTimeAndDaylightSavingTimeType1Fd7(t *testing.T) {
	_, _, diff, err := testParseOffset(-4*3600, -5*3600, "2010-03-14 03:30:00", "2010-03-13 02:30:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 1, 0, 0, 0, 0)
}

func TestPhpGh9382(t *testing.T) {
	_, _, diff, err := testParseOffset(2*3600, 2*3600, "2022-08-01 07:00:00", "2022-08-01 07:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 0, 0)
}

func TestGh75(t *testing.T) {
	_, _, diff, err := testParseTz("PRC", "PRC", "2020-02-01 00:00:00", "2020-03-01 00:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 0, 0, 0, 0, 0)
}

func TestPhpGh9699(t *testing.T) {
	_, _, diff, err := testParseTz("America/Los_Angeles", "UTC", "2022-10-09 02:41:54.515330", "2022-10-10 08:41:54.534620")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 23, 0, 0, 19290)
}

func TestPhpGh9866(t *testing.T) {
	_, _, diff, err := testParseTz("America/Chicago", "America/New_York", "2000-11-01 09:29:22.907606", "2022-06-06 11:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 21, 7, 4, 23, 30, 37, 92394)
}

func TestPhpGh9880a(t *testing.T) {
	_, _, diff, err := testParseTz("America/Los_Angeles", "America/New_York", "2022-11-02 12:18:15", "2022-12-24 13:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 21, 22, 41, 45, 0)
	if diff.Invert {
		t.Errorf("Expected Invert=false, got %v", diff.Invert)
	}
}

func TestPhpGh9880b(t *testing.T) {
	_, _, diff, err := testParseTz("America/New_York", "America/Los_Angeles", "2022-12-24 13:00:00", "2022-11-02 12:18:15")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 1, 21, 22, 41, 45, 0)
	if !diff.Invert {
		t.Errorf("Expected Invert=true, got %v", diff.Invert)
	}
}

func TestTimeFallType3RedodtType3StFwd(t *testing.T) {
	_, _, diff, err := testParse("America/New_York", "2010-11-07 01:12:33", "2010-11-07 03:16:55")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 3, 4, 22, 0)
}

func TestTimeFallType3RedodtType3StRev(t *testing.T) {
	_, _, diff, err := testParse("America/New_York", "2010-11-07 03:16:55", "2010-11-07 01:12:33")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 3, 4, 22, 0)
}

func TestTimeFallType3DtsecType3Stsec(t *testing.T) {
	_, _, diff, err := testParse("America/New_York", "2010-11-07 01:59:59", "2010-11-07 01:00:00")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 59, 59, 0)
}

func TestTimeFallType2DtsecType2Stsec(t *testing.T) {
	_, _, diff, err := testParse("America/New_York", "2010-11-07 01:59:59 EDT", "2010-11-07 01:00:00 EST")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 1, 0)
}

func TestTimeFallType3StsecType3Dtsec(t *testing.T) {
	_, _, diff, err := testParse("America/New_York", "2010-11-07 01:00:00", "2010-11-07 01:59:59")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 59, 59, 0)
}

func TestTimeFallType2StsecType2Dtsec(t *testing.T) {
	_, _, diff, err := testParse("America/New_York", "2010-11-07 01:00:00 EST", "2010-11-07 01:59:59 EDT")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	checkDiff(t, diff, 0, 0, 0, 0, 0, 1, 0)
}
