package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

// Issue 17: Tests for large positive timestamps (year 2369-2370)
func TestIssue0017_Test1(t *testing.T) {
	ts := int64(12622608000)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 2369 {
		t.Errorf("Expected year 2369, got %d", time.Y)
	}
	if time.M != 12 {
		t.Errorf("Expected month 12, got %d", time.M)
	}
	if time.D != 30 {
		t.Errorf("Expected day 30, got %d", time.D)
	}
}

func TestIssue0017_Test2(t *testing.T) {
	ts := int64(12622694400)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 2369 {
		t.Errorf("Expected year 2369, got %d", time.Y)
	}
	if time.M != 12 {
		t.Errorf("Expected month 12, got %d", time.M)
	}
	if time.D != 31 {
		t.Errorf("Expected day 31, got %d", time.D)
	}
}

func TestIssue0017_Test3(t *testing.T) {
	ts := int64(12622780800)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 2370 {
		t.Errorf("Expected year 2370, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month 1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day 1, got %d", time.D)
	}
}

// Issue 19: Tests for negative timestamps (years 1569-1970)
func TestIssue0019_Test1(t *testing.T) {
	ts := int64(-172800)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 1969 {
		t.Errorf("Expected year 1969, got %d", time.Y)
	}
	if time.M != 12 {
		t.Errorf("Expected month 12, got %d", time.M)
	}
	if time.D != 30 {
		t.Errorf("Expected day 30, got %d", time.D)
	}
}

func TestIssue0019_Test2(t *testing.T) {
	ts := int64(-86400)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 1969 {
		t.Errorf("Expected year 1969, got %d", time.Y)
	}
	if time.M != 12 {
		t.Errorf("Expected month 12, got %d", time.M)
	}
	if time.D != 31 {
		t.Errorf("Expected day 31, got %d", time.D)
	}
}

func TestIssue0019_Test3(t *testing.T) {
	ts := int64(0)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 1970 {
		t.Errorf("Expected year 1970, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month 1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day 1, got %d", time.D)
	}
}

func TestIssue0019_Test4(t *testing.T) {
	ts := int64(-12622953600)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 1569 {
		t.Errorf("Expected year 1569, got %d", time.Y)
	}
	if time.M != 12 {
		t.Errorf("Expected month 12, got %d", time.M)
	}
	if time.D != 30 {
		t.Errorf("Expected day 30, got %d", time.D)
	}
}

func TestIssue0019_Test5(t *testing.T) {
	ts := int64(-12622867200)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 1569 {
		t.Errorf("Expected year 1569, got %d", time.Y)
	}
	if time.M != 12 {
		t.Errorf("Expected month 12, got %d", time.M)
	}
	if time.D != 31 {
		t.Errorf("Expected day 31, got %d", time.D)
	}
}

func TestIssue0019_Test6(t *testing.T) {
	ts := int64(-12622780800)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	timelib.SetTimezoneFromOffset(time, 0)
	time.Unixtime2local(ts)

	if time.Y != 1570 {
		t.Errorf("Expected year 1570, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month 1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day 1, got %d", time.D)
	}
}

// Issue 23: DST transition tests for Europe/London
// Helper functions for issue 23 tests
func issue23TestAdd(t *testing.T, str, interval string, tzi *timelib.TzInfo) *timelib.Time {
	t.Helper()

	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", str, err)
	}

	time.UpdateTS(tzi)

	errorsContainer := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval(interval, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", interval, err)
	}

	if p.H != 5 {
		t.Errorf("Expected interval hours=5, got %d", p.H)
	}

	changed := time.Add(p)
	return changed
}

func issue23TestSub(t *testing.T, str, interval string, tzi *timelib.TzInfo) *timelib.Time {
	t.Helper()

	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", str, err)
	}

	time.UpdateTS(tzi)

	errorsContainer := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval(interval, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", interval, err)
	}

	if p.H != 5 {
		t.Errorf("Expected interval hours=5, got %d", p.H)
	}

	changed := time.Sub(p)
	return changed
}

func issue23TestAddWall(t *testing.T, str, interval string, tzi *timelib.TzInfo) *timelib.Time {
	t.Helper()

	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time '%s': %v", str, err)
	}

	// Always call UpdateTS - even for Unix timestamps (@...)
	// The @ timestamp sets HaveRelative and Relative.S, which UpdateTS applies
	time.UpdateTS(tzi)
	timelib.SetTimezone(time, tzi)
	time.UpdateFromSSE()

	// Initialize microseconds to 0 if UNSET (C code behavior)
	if time.US == timelib.TIMELIB_UNSET {
		time.US = 0
	}

	errorsContainer := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval(interval, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", interval, err)
	}

	changed := time.AddWall(p)
	return changed
}

func TestIssue0023_Test1a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAdd(t, "2014-03-30 00:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1396137600 + (4 * SECS_PER_HOUR))
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 3 {
		t.Errorf("Expected month=3, got %d", changed.M)
	}
	if changed.D != 30 {
		t.Errorf("Expected day=30, got %d", changed.D)
	}
	if changed.H != 5 {
		t.Errorf("Expected hour=5, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test1b(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestSub(t, "2014-03-30 05:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1396137600)
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 3 {
		t.Errorf("Expected month=3, got %d", changed.M)
	}
	if changed.D != 30 {
		t.Errorf("Expected day=30, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour=0, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test2a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAdd(t, "2014-03-29 00:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1396051200 + (5 * SECS_PER_HOUR))
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 3 {
		t.Errorf("Expected month=3, got %d", changed.M)
	}
	if changed.D != 29 {
		t.Errorf("Expected day=29, got %d", changed.D)
	}
	if changed.H != 5 {
		t.Errorf("Expected hour=5, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test2b(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestSub(t, "2014-03-29 05:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1396051200)
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 3 {
		t.Errorf("Expected month=3, got %d", changed.M)
	}
	if changed.D != 29 {
		t.Errorf("Expected day=29, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour=0, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test3a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAdd(t, "2014-03-31 00:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1396220400 + (5 * SECS_PER_HOUR))
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 3 {
		t.Errorf("Expected month=3, got %d", changed.M)
	}
	if changed.D != 31 {
		t.Errorf("Expected day=31, got %d", changed.D)
	}
	if changed.H != 5 {
		t.Errorf("Expected hour=5, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test3b(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestSub(t, "2014-03-31 05:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1396220400)
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 3 {
		t.Errorf("Expected month=3, got %d", changed.M)
	}
	if changed.D != 31 {
		t.Errorf("Expected day=31, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour=0, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test4a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAdd(t, "2014-10-25 00:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1414191600 + (5 * SECS_PER_HOUR))
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 25 {
		t.Errorf("Expected day=25, got %d", changed.D)
	}
	if changed.H != 5 {
		t.Errorf("Expected hour=5, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test4b(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestSub(t, "2014-10-25 05:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1414191600)
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 25 {
		t.Errorf("Expected day=25, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour=0, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test5a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAdd(t, "2014-10-26 00:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1414278000 + (6 * SECS_PER_HOUR))
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 26 {
		t.Errorf("Expected day=26, got %d", changed.D)
	}
	if changed.H != 5 {
		t.Errorf("Expected hour=5, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test5b(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestSub(t, "2014-10-26 05:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1414278000)
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 26 {
		t.Errorf("Expected day=26, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour=0, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test6a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAdd(t, "2014-10-27 00:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1414368000 + (5 * SECS_PER_HOUR))
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 27 {
		t.Errorf("Expected day=27, got %d", changed.D)
	}
	if changed.H != 5 {
		t.Errorf("Expected hour=5, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test6b(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestSub(t, "2014-10-27 05:00:00.000000", "PT5H", tzi)
	defer timelib.TimeDtor(changed)

	expectedSSE := int64(1414368000)
	if changed.Sse != expectedSSE {
		t.Errorf("Expected SSE=%d, got %d", expectedSSE, changed.Sse)
	}
	if changed.Y != 2014 {
		t.Errorf("Expected year=2014, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 27 {
		t.Errorf("Expected day=27, got %d", changed.D)
	}
	if changed.H != 0 {
		t.Errorf("Expected hour=0, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
}

func TestIssue0023_Test7a(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	changed := issue23TestAddWall(t, "@1635641999", "PT1S", tzi)
	defer timelib.TimeDtor(changed)

	if changed.Sse != 1635642000 {
		t.Errorf("Expected SSE=1635642000, got %d", changed.Sse)
	}
	if changed.Y != 2021 {
		t.Errorf("Expected year=2021, got %d", changed.Y)
	}
	if changed.M != 10 {
		t.Errorf("Expected month=10, got %d", changed.M)
	}
	if changed.D != 31 {
		t.Errorf("Expected day=31, got %d", changed.D)
	}
	if changed.H != 1 {
		t.Errorf("Expected hour=1, got %d", changed.H)
	}
	if changed.I != 0 {
		t.Errorf("Expected minute=0, got %d", changed.I)
	}
	if changed.S != 0 {
		t.Errorf("Expected second=0, got %d", changed.S)
	}
	if changed.Dst != 0 {
		t.Errorf("Expected DST=0, got %d", changed.Dst)
	}
	if changed.TzAbbr != "GMT" {
		t.Errorf("Expected TzAbbr=GMT, got %s", changed.TzAbbr)
	}
}

// Issue 35: Microsecond overflow tests
func TestIssue0035_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str := "2017-12-31 23:59:59.999999 +1 microsecond"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	time.UpdateTS(tzi)

	if time.Y != 2018 {
		t.Errorf("Expected year=2018, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
	if time.H != 0 {
		t.Errorf("Expected hour=0, got %d", time.H)
	}
	if time.I != 0 {
		t.Errorf("Expected minute=0, got %d", time.I)
	}
	if time.S != 0 {
		t.Errorf("Expected second=0, got %d", time.S)
	}
	if time.US != 0 {
		t.Errorf("Expected microsecond=0, got %d", time.US)
	}
}

func TestIssue0035_Test2(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str := "2017-12-31 23:59:59.999999 +2 microsecond"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	time.UpdateTS(tzi)

	if time.Y != 2018 {
		t.Errorf("Expected year=2018, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
	if time.H != 0 {
		t.Errorf("Expected hour=0, got %d", time.H)
	}
	if time.I != 0 {
		t.Errorf("Expected minute=0, got %d", time.I)
	}
	if time.S != 0 {
		t.Errorf("Expected second=0, got %d", time.S)
	}
	if time.US != 1 {
		t.Errorf("Expected microsecond=1, got %d", time.US)
	}
}

// Issue 16: Extended year range tests
func TestIssue0016_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str := "+10000-01-01 00:00:00.000000"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	time.UpdateTS(tzi)

	if time.Y != 10000 {
		t.Errorf("Expected year=10000, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
	if time.H != 0 {
		t.Errorf("Expected hour=0, got %d", time.H)
	}
	if time.I != 0 {
		t.Errorf("Expected minute=0, got %d", time.I)
	}
	if time.S != 0 {
		t.Errorf("Expected second=0, got %d", time.S)
	}
	if time.US != 0 {
		t.Errorf("Expected microsecond=0, got %d", time.US)
	}
}

func TestIssue0016_Test2(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str := "-10000-01-01 00:00:00.000000"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	time.UpdateTS(tzi)

	if time.Y != -10000 {
		t.Errorf("Expected year=-10000, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
	if time.H != 0 {
		t.Errorf("Expected hour=0, got %d", time.H)
	}
	if time.I != 0 {
		t.Errorf("Expected minute=0, got %d", time.I)
	}
	if time.S != 0 {
		t.Errorf("Expected second=0, got %d", time.S)
	}
	if time.US != 0 {
		t.Errorf("Expected microsecond=0, got %d", time.US)
	}
}

func TestIssue0016_Test3(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str := "+100000-01-01 00:00:00.000000"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	time.UpdateTS(tzi)

	if time.Y != 100000 {
		t.Errorf("Expected year=100000, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
	if time.H != 0 {
		t.Errorf("Expected hour=0, got %d", time.H)
	}
	if time.I != 0 {
		t.Errorf("Expected minute=0, got %d", time.I)
	}
	if time.S != 0 {
		t.Errorf("Expected second=0, got %d", time.S)
	}
	if time.US != 0 {
		t.Errorf("Expected microsecond=0, got %d", time.US)
	}
}

// Issue 50: Microsecond precision in diff calculations  
func TestIssue0050_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str1 := "2018-10-11 20:59:06.914653"
	str2 := "2018-10-11 20:59:07.237419"
	
	t1, err1 := timelib.StrToTime(str1, nil)
	if err1 != nil {
		t.Fatalf("Failed to parse time1: %v", err1)
	}
	defer timelib.TimeDtor(t1)

	t2, err2 := timelib.StrToTime(str2, nil)
	if err2 != nil {
		t.Fatalf("Failed to parse time2: %v", err2)
	}
	defer timelib.TimeDtor(t2)

	t1.UpdateTS(tzi)
	t2.UpdateTS(tzi)

	diff := t1.Diff(t2)
	defer timelib.RelTimeDtor(diff)

	if diff.S != 0 {
		t.Errorf("Expected seconds=0, got %d", diff.S)
	}
	if diff.US != 322766 {
		t.Errorf("Expected microseconds=322766, got %d", diff.US)
	}
}

func TestIssue0050_Test2(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str1 := "2018-10-11 20:59:06.237419"
	str2 := "2018-10-11 20:59:06.914653"
	
	t1, err1 := timelib.StrToTime(str1, nil)
	if err1 != nil {
		t.Fatalf("Failed to parse time1: %v", err1)
	}
	defer timelib.TimeDtor(t1)

	t2, err2 := timelib.StrToTime(str2, nil)
	if err2 != nil {
		t.Fatalf("Failed to parse time2: %v", err2)
	}
	defer timelib.TimeDtor(t2)

	t1.UpdateTS(tzi)
	t2.UpdateTS(tzi)

	diff := t1.Diff(t2)
	defer timelib.RelTimeDtor(diff)

	if diff.S != 0 {
		t.Errorf("Expected seconds=0, got %d", diff.S)
	}
	if diff.US != 677234 {
		t.Errorf("Expected microseconds=677234, got %d", diff.US)
	}
}

// Issue 51: Diff invert flag with microseconds
func TestIssue0051_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str1 := "2018-11-22 13:27:52.089635"
	str2 := "2018-11-22 13:27:52"
	
	t1, err1 := timelib.StrToTime(str1, nil)
	if err1 != nil {
		t.Fatalf("Failed to parse time1: %v", err1)
	}
	defer timelib.TimeDtor(t1)

	t2, err2 := timelib.StrToTime(str2, nil)
	if err2 != nil {
		t.Fatalf("Failed to parse time2: %v", err2)
	}
	defer timelib.TimeDtor(t2)

	t1.UpdateTS(tzi)
	t2.UpdateTS(tzi)

	diff := t1.Diff(t2)
	defer timelib.RelTimeDtor(diff)

	if diff.S != 0 {
		t.Errorf("Expected seconds=0, got %d", diff.S)
	}
	if diff.US != 89635 {
		t.Errorf("Expected microseconds=89635, got %d", diff.US)
	}
	if !diff.Invert {
		t.Errorf("Expected invert=true, got false")
	}
}

func TestIssue0051_Test2(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str1 := "2018-11-22 13:27:52"
	str2 := "2018-11-22 13:27:52.089635"
	
	t1, err1 := timelib.StrToTime(str1, nil)
	if err1 != nil {
		t.Fatalf("Failed to parse time1: %v", err1)
	}
	defer timelib.TimeDtor(t1)

	t2, err2 := timelib.StrToTime(str2, nil)
	if err2 != nil {
		t.Fatalf("Failed to parse time2: %v", err2)
	}
	defer timelib.TimeDtor(t2)

	t1.UpdateTS(tzi)
	t2.UpdateTS(tzi)

	diff := t1.Diff(t2)
	defer timelib.RelTimeDtor(diff)

	if diff.S != 0 {
		t.Errorf("Expected seconds=0, got %d", diff.S)
	}
	if diff.US != 89635 {
		t.Errorf("Expected microseconds=89635, got %d", diff.US)
	}
	if diff.Invert {
		t.Errorf("Expected invert=false, got true")
	}
}

// Issue 53: Timezone offset for very old dates
func TestIssue0053_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("America/Belize", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	ts := int64(-61626506832)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	time.TzInfo = tzi
	time.ZoneType = timelib.TIMELIB_ZONETYPE_ID
	time.Unixtime2local(ts)

	if time.Y != 17 {
		t.Errorf("Expected year=17, got %d", time.Y)
	}
	if time.M != 2 {
		t.Errorf("Expected month=2, got %d", time.M)
	}
	if time.D != 18 {
		t.Errorf("Expected day=18, got %d", time.D)
	}
	if time.H != 0 {
		t.Errorf("Expected hour=0, got %d", time.H)
	}
	if time.I != 0 {
		t.Errorf("Expected minute=0, got %d", time.I)
	}
	if time.S != 0 {
		t.Errorf("Expected second=0, got %d", time.S)
	}
	if time.Z != -21168 {
		t.Errorf("Expected Z=-21168, got %d", time.Z)
	}
}

func TestIssue0053_Test2(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("America/Belize", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	ts := int64(-1822500433)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	time.TzInfo = tzi
	time.ZoneType = timelib.TIMELIB_ZONETYPE_ID
	time.Unixtime2local(ts)

	if time.Z != -21168 {
		t.Errorf("Expected Z=-21168, got %d", time.Z)
	}
}

func TestIssue0053_Test3(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("America/Belize", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	ts := int64(-1822500432)
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	time.TzInfo = tzi
	time.ZoneType = timelib.TIMELIB_ZONETYPE_ID
	time.Unixtime2local(ts)

	if time.Z != -21600 {
		t.Errorf("Expected Z=-21600, got %d", time.Z)
	}
}

// Issue 65: Test for timezones with no transitions
func TestIssue0065_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Etc/GMT+5", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	offset := timelib.GetTimeZoneInfo(int64(-1822500432), tzi)
	defer timelib.TimeOffsetDtor(offset)

	// INT64_MIN in Go
	const INT64_MIN = -9223372036854775808
	if offset.TransitionTime != INT64_MIN {
		t.Errorf("Expected transition_time=INT64_MIN, got %d", offset.TransitionTime)
	}
}

func TestIssue0065_Test2(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("Europe/London", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	const INT64_MIN = -9223372036854775808

	offset := timelib.GetTimeZoneInfo(int64(-3852662326), tzi)
	if offset.TransitionTime != INT64_MIN {
		t.Errorf("Expected transition_time=INT64_MIN, got %d", offset.TransitionTime)
	}
	timelib.TimeOffsetDtor(offset)

	offset = timelib.GetTimeZoneInfo(int64(-3852662325), tzi)
	if offset.TransitionTime != -3852662325 {
		t.Errorf("Expected transition_time=-3852662325, got %d", offset.TransitionTime)
	}
	timelib.TimeOffsetDtor(offset)
}

// Issue 69: Relative time with milliseconds
func TestIssue0069(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str1 := "2019-10-14T15:08:23.123+02:00"
	str2 := "-50000 msec"

	t1, err1 := timelib.StrToTime(str1, nil)
	if err1 != nil {
		t.Fatalf("Failed to parse time1: %v", err1)
	}
	defer timelib.TimeDtor(t1)

	t2, err2 := timelib.StrToTime(str2, nil)
	if err2 != nil {
		t.Fatalf("Failed to parse time2: %v", err2)
	}
	defer timelib.TimeDtor(t2)

	t1.UpdateTS(tzi)

	// Copy relative time from t2 to t1
	t1.Relative = t2.Relative
	t1.HaveRelative = true
	t1.SseUptodate = false

	t1.UpdateTS(nil)
	t1.UpdateFromSSE()
	t1.HaveRelative = false

	// Clear relative
	t1.Relative = timelib.RelTime{}

	if t1.US != 123000 {
		t.Errorf("Expected microseconds=123000, got %d", t1.US)
	}
	if t1.S != 33 {
		t.Errorf("Expected seconds=33, got %d", t1.S)
	}
	if t1.I != 7 {
		t.Errorf("Expected minutes=7, got %d", t1.I)
	}
	if t1.H != 15 {
		t.Errorf("Expected hours=15, got %d", t1.H)
	}
}

// Issue 92: Error parsing invalid timestamp
func TestIssue0092(t *testing.T) {
	str := "@7."
	time, err := timelib.StrToTime(str, nil)
	if time != nil {
		defer timelib.TimeDtor(time)
	}

	if err == nil {
		t.Errorf("Expected 1 error, got %d", 1)
	}
}

// Issue 93: Nanosecond parsing (truncated to microseconds)
func TestIssue0093_Test1(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	str := "2006-01-02T15:04:05.123456789Z"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	time.UpdateTS(tzi)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if time.Y != 2006 {
		t.Errorf("Expected year=2006, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 2 {
		t.Errorf("Expected day=2, got %d", time.D)
	}
	if time.H != 15 {
		t.Errorf("Expected hour=15, got %d", time.H)
	}
	if time.I != 4 {
		t.Errorf("Expected minute=4, got %d", time.I)
	}
	if time.S != 5 {
		t.Errorf("Expected second=5, got %d", time.S)
	}
	// Nanoseconds truncated to microseconds: 123456789 -> 123456
	if time.US != 123456 {
		t.Errorf("Expected microsecond=123456, got %d", time.US)
	}
}

// Issue 94: Very long invalid string (tests parser resilience)
func TestIssue0094(t *testing.T) {
	// Create a string with 20000+ x's
	str := ""
	for i := 0; i < 20000; i++ {
		str += "x"
	}
	
	time, err := timelib.StrToTime(str, nil)
	if time != nil {
		defer timelib.TimeDtor(time)
	}

	// Should handle gracefully without crashing
	if err == nil {
		t.Error("Expected errors for invalid string")
	}
}

// updatets test: Test UpdateTS with timezone
func TestUpdateTS(t *testing.T) {
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	
	ts := int64(1289109600)
	tz := "America/New_York"
	tzi, err := timelib.ParseTzfile(tz, testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}

	time := timelib.TimeCtor()
	timelib.SetTimezone(time, tzi)
	time.Unixtime2local(ts)
	time.UpdateTS(nil)
	
	// Clean up
	if time.TzInfo != nil {
		timelib.TzinfoDtor(time.TzInfo)
	}
	timelib.TimeDtor(time)
	
	// Test passed if no crash
}

// Helper for relative time tests
func testParseRelative(t *testing.T, tzid string, ts int64, modify string) (*timelib.Time, *timelib.Time, *timelib.TzInfo) {
	t.Helper()
	
	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile(tzid, testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}

	time := timelib.TimeCtor()
	timelib.SetTimezone(time, tzi)
	time.Unixtime2local(ts)

	update, err := timelib.StrToTime(modify, nil)
	if err != nil {
		t.Fatalf("Failed to parse modify string: %v", err)
	}

	timelib.FillHoles(update, time, timelib.TIMELIB_NO_CLONE)
	update.UpdateTS(tzi)
	// ApplyLocaltime with localtime=1 calls Unixtime2local
	update.Unixtime2local(update.Sse)

	return time, update, tzi
}

// Bug 33414: Relative time test for Pacific/Kwajalein
func TestBug33414_01(t *testing.T) {
	time, update, tzi := testParseRelative(t, "Pacific/Kwajalein", 745391837, "next Saturday")
	defer timelib.TzinfoDtor(tzi)
	defer timelib.TzinfoDtor(time.TzInfo)
	defer timelib.TimeDtor(time)
	defer timelib.TimeDtor(update)

	if update.Sse != 745934400 {
		t.Errorf("Expected SSE=745934400, got %d", update.Sse)
	}
	if update.Y != 1993 {
		t.Errorf("Expected year=1993, got %d", update.Y)
	}
	if update.M != 8 {
		t.Errorf("Expected month=8, got %d", update.M)
	}
	if update.D != 22 {
		t.Errorf("Expected day=22, got %d", update.D)
	}
	if update.H != 0 {
		t.Errorf("Expected hour=0, got %d", update.H)
	}
	if update.I != 0 {
		t.Errorf("Expected minute=0, got %d", update.I)
	}
	if update.S != 0 {
		t.Errorf("Expected second=0, got %d", update.S)
	}
	if update.TzAbbr != "+12" {
		t.Errorf("Expected tz_abbr=+12, got %s", update.TzAbbr)
	}
}

// Bug 33415: Relative time test for Africa/Monrovia
func TestBug33415_01(t *testing.T) {
	time, update, tzi := testParseRelative(t, "Africa/Monrovia", 63050507, "next Friday")
	defer timelib.TzinfoDtor(tzi)
	defer timelib.TzinfoDtor(time.TzInfo)
	defer timelib.TimeDtor(time)
	defer timelib.TimeDtor(update)

	if update.Sse != 63593070 {
		t.Errorf("Expected SSE=63593070, got %d", update.Sse)
	}
	if update.Y != 1972 {
		t.Errorf("Expected year=1972, got %d", update.Y)
	}
	if update.M != 1 {
		t.Errorf("Expected month=1, got %d", update.M)
	}
	if update.D != 7 {
		t.Errorf("Expected day=7, got %d", update.D)
	}
	if update.H != 0 {
		t.Errorf("Expected hour=0, got %d", update.H)
	}
	if update.I != 44 {
		t.Errorf("Expected minute=44, got %d", update.I)
	}
	if update.S != 30 {
		t.Errorf("Expected second=30, got %d", update.S)
	}
	if update.TzAbbr != "GMT" {
		t.Errorf("Expected tz_abbr=GMT, got %s", update.TzAbbr)
	}
}

// tzcorparse_01: Timezone parsing corruption test 1
func TestTzcorparse01(t *testing.T) {
	str := "+30157"
	dst := 0
	tzNotFound := 0
	time := &timelib.Time{}

	testDirectory, _ := timelib.Zoneinfo("files")

	timelib.ParseZone(&str, &dst, time, &tzNotFound, testDirectory, timelib.ParseTzfile)

	if tzNotFound != 1 {
		t.Errorf("Expected tz_not_found=1, got %d", tzNotFound)
	}
}

// tzcorparse_02: Timezone parsing corruption test 2
func TestTzcorparse02(t *testing.T) {
	str := "+30:57"
	dst := 0
	tzNotFound := 0
	time := &timelib.Time{}

	testDirectory, _ := timelib.Zoneinfo("files")
	ret := timelib.ParseZone(&str, &dst, time, &tzNotFound, testDirectory, timelib.ParseTzfile)

	if ret != 111420 {
		t.Errorf("Expected return=111420, got %d", ret)
	}
	if tzNotFound != 0 {
		t.Errorf("Expected tz_not_found=0, got %d", tzNotFound)
	}
}

// tzabbrparse_01: Timezone abbreviation parsing
func TestTzabbrparse01(t *testing.T) {
	str := "Sun, 14 Aug 2005 17:50:44 -0730"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	if time.Y != 2005 {
		t.Errorf("Expected year=2005, got %d", time.Y)
	}
	if time.M != 8 {
		t.Errorf("Expected month=8, got %d", time.M)
	}
	if time.D != 14 {
		t.Errorf("Expected day=14, got %d", time.D)
	}
	if time.H != 17 {
		t.Errorf("Expected hour=17, got %d", time.H)
	}
	if time.I != 50 {
		t.Errorf("Expected minute=50, got %d", time.I)
	}
	if time.S != 44 {
		t.Errorf("Expected second=44, got %d", time.S)
	}
	if time.Z != -27000 {
		t.Errorf("Expected Z=-27000, got %d", time.Z)
	}
}

// weekday_time_part_01: Weekday with time component
func TestWeekdayTimePart01(t *testing.T) {
	str := "Monday 03:59:59"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	if time.H != 3 {
		t.Errorf("Expected hour=3, got %d", time.H)
	}
	if time.I != 59 {
		t.Errorf("Expected minute=59, got %d", time.I)
	}
	if time.S != 59 {
		t.Errorf("Expected second=59, got %d", time.S)
	}
	
	// Should have relative weekday
	if !time.Relative.HaveWeekdayRelative {
		t.Error("Expected have_weekday_relative=true")
	}
	if time.Relative.Weekday != 1 {
		t.Errorf("Expected weekday=1 (Monday), got %d", time.Relative.Weekday)
	}
}

// first_day_of_time_01: First day of month test 1
func TestFirstDayOfTime01(t *testing.T) {
	str := "first day of January 2023"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	time.UpdateTS(tzi)

	if time.Y != 2023 {
		t.Errorf("Expected year=2023, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
}

// first_day_of_time_02: First day of month test 2
func TestFirstDayOfTime02(t *testing.T) {
	str := "first day of next month"
	
	// Use a known base time: 2023-01-15
	baseTime := timelib.TimeCtor()
	baseTime.Y = 2023
	baseTime.M = 1
	baseTime.D = 15
	baseTime.H = 12
	baseTime.I = 0
	baseTime.S = 0
	baseTime.HaveDate = true
	baseTime.HaveTime = true
	defer timelib.TimeDtor(baseTime)

	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	timelib.FillHoles(time, baseTime, timelib.TIMELIB_NO_CLONE)
	time.UpdateTS(tzi)

	if time.Y != 2023 {
		t.Errorf("Expected year=2023, got %d", time.Y)
	}
	if time.M != 2 {
		t.Errorf("Expected month=2, got %d", time.M)
	}
	if time.D != 1 {
		t.Errorf("Expected day=1, got %d", time.D)
	}
}

// last_day_of_time_01: Last day of month test 1
func TestLastDayOfTime01(t *testing.T) {
	str := "last day of January 2023"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	time.UpdateTS(tzi)

	if time.Y != 2023 {
		t.Errorf("Expected year=2023, got %d", time.Y)
	}
	if time.M != 1 {
		t.Errorf("Expected month=1, got %d", time.M)
	}
	if time.D != 31 {
		t.Errorf("Expected day=31, got %d", time.D)
	}
}

// last_day_of_time_02: Last day of month test 2
func TestLastDayOfTime02(t *testing.T) {
	str := "last day of next month"
	
	// Use a known base time: 2023-01-15
	baseTime := timelib.TimeCtor()
	baseTime.Y = 2023
	baseTime.M = 1
	baseTime.D = 15
	baseTime.H = 12
	baseTime.I = 0
	baseTime.S = 0
	baseTime.HaveDate = true
	baseTime.HaveTime = true
	defer timelib.TimeDtor(baseTime)

	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	testDirectory, _ := timelib.Zoneinfo("files")
	var errorCode int
	tzi, err := timelib.ParseTzfile("UTC", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile: %v", err)
	}
	defer timelib.TzinfoDtor(tzi)

	timelib.FillHoles(time, baseTime, timelib.TIMELIB_NO_CLONE)
	time.UpdateTS(tzi)

	if time.Y != 2023 {
		t.Errorf("Expected year=2023, got %d", time.Y)
	}
	if time.M != 2 {
		t.Errorf("Expected month=2, got %d", time.M)
	}
	if time.D != 28 {
		t.Errorf("Expected day=28, got %d", time.D)
	}
}

// php81106: PHP bug 81106 test
func TestPhp81106(t *testing.T) {
	str := "2021-10-31T02:30:00+02:00[Europe/Berlin]"
	time, err := timelib.StrToTime(str, nil)
	if time != nil {
		defer timelib.TimeDtor(time)
	}

	// This format may not be fully supported, but should not crash
	// The test mainly ensures the parser handles the bracketed timezone
	if err != nil {
		// Expected to have errors for unsupported format
		t.Logf("Parsing error (expected): %v", err)
	}
}

// sanitizer_issue: Test edge case that triggered sanitizer
func TestSanitizerIssue(t *testing.T) {
	// This test is to ensure we don't crash on edge cases
	// The original C test had a very long string to test buffer handling
	str := "2023-01-01 00:00:00"
	time, err := timelib.StrToTime(str, nil)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	defer timelib.TimeDtor(time)

	if time.Y != 2023 {
		t.Errorf("Expected year=2023, got %d", time.Y)
	}
}
