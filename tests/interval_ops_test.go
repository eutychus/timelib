package tests

import (
	"fmt"
	"testing"

	timelib "github.com/eutychus/timelib"
)

// testAddWall tests interval addition with wall clock semantics
func testAddWall(t *testing.T, base, tzid, interval string, us int64, invert bool) *timelib.Time {
	t.Helper()

	testDirectory, err := timelib.Zoneinfo("files")
	if err != nil {
		t.Fatalf("Failed to load timezone directory: %v", err)
	}

	var errorCode int
	tzi, err := timelib.ParseTzfile(tzid, testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile %s: %v", tzid, err)
	}

	tBase, errorsBase := timelib.Strtotime(base)
	if errorsBase != nil && errorsBase.ErrorCount > 0 {
		t.Fatalf("Failed to parse time '%s': %v", base, errorsBase.ErrorMessages)
	}

	tBase.UpdateTS(tzi)
	timelib.SetTimezone(tBase, tzi)
	tBase.UpdateFromSSE()

	// Parse microseconds from base string if present (e.g., "2000-01-01 00:00:01.500000")
	// Strtotime may not parse fractional seconds properly
	// This must be done AFTER UpdateFromSSE to avoid being cleared
	parseBaseMicroseconds(base, tBase)

	errorsContainer := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval(interval, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", interval, err)
	}

	p.US = us
	p.Invert = invert

	changed := tBase.AddWall(p)
	return changed
}

// parseBaseMicroseconds extracts microseconds from time string if present
func parseBaseMicroseconds(base string, t *timelib.Time) {
	// Parse microseconds from base string if present (e.g., "2000-01-01 00:00:01.500000")
	// Strtotime may not parse fractional seconds properly
	if len(base) > 20 && base[19] == '.' {
		// Extract microseconds from string
		var usParsed int64
		_, err := fmt.Sscanf(base[20:], "%d", &usParsed)
		if err == nil {
			// Pad to 6 digits if needed
			digitCount := len(base) - 20
			for digitCount < 6 {
				usParsed *= 10
				digitCount++
			}
			t.US = usParsed
		}
	}
}

// testSubWall tests interval subtraction with wall clock semantics
func testSubWall(t *testing.T, base, tzid, interval string, us int64, invert bool) *timelib.Time {
	t.Helper()

	testDirectory, err := timelib.Zoneinfo("files")
	if err != nil {
		t.Fatalf("Failed to load timezone directory: %v", err)
	}

	var errorCode int
	tzi, err := timelib.ParseTzfile(tzid, testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse tzfile %s: %v", tzid, err)
	}

	tBase, errorsBase := timelib.Strtotime(base)
	if errorsBase != nil && errorsBase.ErrorCount > 0 {
		t.Fatalf("Failed to parse time '%s': %v", base, errorsBase.ErrorMessages)
	}

	tBase.UpdateTS(tzi)
	timelib.SetTimezone(tBase, tzi)
	tBase.UpdateFromSSE()

	// Parse microseconds from base string if present
	// This must be done AFTER UpdateFromSSE to avoid being cleared
	parseBaseMicroseconds(base, tBase)

	errorsContainer := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval(interval, errorsContainer)
	if err != nil {
		t.Fatalf("Failed to parse interval '%s': %v", interval, err)
	}

	p.US = us
	p.Invert = invert

	changed := tBase.SubWall(p)
	return changed
}

// checkInterval verifies the result of an interval operation
func checkInterval(t *testing.T, result *timelib.Time, y, m, d, h, i, s int64, us int64, sse int64) {
	t.Helper()
	if result.Y != y {
		t.Errorf("Expected year=%d, got %d", y, result.Y)
	}
	if result.M != m {
		t.Errorf("Expected month=%d, got %d", m, result.M)
	}
	if result.D != d {
		t.Errorf("Expected day=%d, got %d", d, result.D)
	}
	if result.H != h {
		t.Errorf("Expected hour=%d, got %d", h, result.H)
	}
	if result.I != i {
		t.Errorf("Expected minute=%d, got %d", i, result.I)
	}
	if result.S != s {
		t.Errorf("Expected second=%d, got %d", s, result.S)
	}
	if result.US != us {
		t.Errorf("Expected microsecond=%d, got %d", us, result.US)
	}
	if result.Sse != sse {
		t.Errorf("Expected sse=%d, got %d", sse, result.Sse)
	}
}

// PHP bug 80998: Microsecond handling in intervals
// https://bugs.php.net/bug.php?id=80998

func TestIntervalPhp80998_1a(t *testing.T) {
	result := testAddWall(t, "2021-04-05 14:00:00", "UTC", "PT10799S", 999999, true)
	checkInterval(t, result, 2021, 4, 5, 11, 0, 0, 1, 1617620400)
}

func TestIntervalPhp80998_1b(t *testing.T) {
	result := testAddWall(t, "2021-04-05 14:00:00", "UTC", "PT10800S", 0, true)
	checkInterval(t, result, 2021, 4, 5, 11, 0, 0, 0, 1617620400)
}

func TestIntervalPhp80998_2a(t *testing.T) {
	result := testAddWall(t, "2000-01-01 00:00:00", "UTC", "PT1S", 500000, true)
	checkInterval(t, result, 1999, 12, 31, 23, 59, 58, 500000, 946684798)
}

func TestIntervalPhp80998_2b(t *testing.T) {
	result := testAddWall(t, "2000-01-01 00:00:00", "UTC", "PT1S", 0, true)
	checkInterval(t, result, 1999, 12, 31, 23, 59, 59, 0, 946684799)
}

func TestIntervalPhp80998_1aSub(t *testing.T) {
	result := testSubWall(t, "2021-04-05 11:00:00", "UTC", "PT2H59M59S", 999999, true)
	checkInterval(t, result, 2021, 4, 5, 13, 59, 59, 999999, 1617631199)
}

func TestIntervalPhp80998_1bSub(t *testing.T) {
	result := testSubWall(t, "2021-04-05 11:00:00", "UTC", "PT3H", 0, true)
	checkInterval(t, result, 2021, 4, 5, 14, 0, 0, 0, 1617631200)
}

func TestIntervalPhp80998_2aSub(t *testing.T) {
	result := testSubWall(t, "2021-04-05 11:00:00", "UTC", "PT2H59M59S", 999999, false)
	checkInterval(t, result, 2021, 4, 5, 8, 0, 0, 1, 1617609600)
}

func TestIntervalPhp80998_2bSub(t *testing.T) {
	result := testSubWall(t, "2021-04-05 11:00:00", "UTC", "PT3H", 0, false)
	checkInterval(t, result, 2021, 4, 5, 8, 0, 0, 0, 1617609600)
}

// GitHub issue 8964: Negative microseconds
// https://github.com/php/php-src/issues/8964

func TestIntervalGh8964a(t *testing.T) {
	result := testAddWall(t, "2022-07-21 14:50:13", "UTC", "PT0S", -5, false)
	checkInterval(t, result, 2022, 7, 21, 14, 50, 12, 999995, 1658415012)
}

func TestIntervalGh8964b(t *testing.T) {
	result := testAddWall(t, "2022-07-21 14:50:13", "UTC", "PT0S", -5, true)
	checkInterval(t, result, 2022, 7, 21, 14, 50, 13, 5, 1658415013)
}

func TestIntervalGh8964c(t *testing.T) {
	result := testSubWall(t, "2022-07-21 14:50:13", "UTC", "PT0S", -5, false)
	checkInterval(t, result, 2022, 7, 21, 14, 50, 13, 5, 1658415013)
}

func TestIntervalGh8964d(t *testing.T) {
	result := testSubWall(t, "2022-07-21 14:50:13", "UTC", "PT0S", -5, true)
	checkInterval(t, result, 2022, 7, 21, 14, 50, 12, 999995, 1658415012)
}

func TestIntervalGh8964e(t *testing.T) {
	result := testSubWall(t, "2022-07-21 15:00:10", "UTC", "PT0S", 900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 9, 100000, 1658415609)
}

func TestIntervalGh8964f(t *testing.T) {
	result := testAddWall(t, "2022-07-21 15:00:10", "UTC", "PT0S", 900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 10, 900000, 1658415610)
}

func TestIntervalGh8964g(t *testing.T) {
	result := testSubWall(t, "2022-07-21 15:00:10", "UTC", "PT0S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 10, 900000, 1658415610)
}

func TestIntervalGh8964h(t *testing.T) {
	result := testAddWall(t, "2022-07-21 15:00:10", "UTC", "PT0S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 9, 100000, 1658415609)
}

func TestIntervalGh8964i(t *testing.T) {
	result := testSubWall(t, "2022-07-21 15:00:10", "UTC", "PT1S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 9, 900000, 1658415609)
}

func TestIntervalGh8964j(t *testing.T) {
	result := testAddWall(t, "2022-07-21 15:00:10", "UTC", "PT1S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 10, 100000, 1658415610)
}

func TestIntervalGh8964k(t *testing.T) {
	result := testSubWall(t, "2022-07-21 15:00:10", "UTC", "PT2S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 8, 900000, 1658415608)
}

func TestIntervalGh8964l(t *testing.T) {
	result := testAddWall(t, "2022-07-21 15:00:10", "UTC", "PT2S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 11, 100000, 1658415611)
}

func TestIntervalGh8964m(t *testing.T) {
	result := testSubWall(t, "2022-07-21 15:00:09.100000", "UTC", "PT0S", -900000, false)
	checkInterval(t, result, 2022, 7, 21, 15, 0, 10, 0, 1658415610)
}

// GitHub issue 9106: Microsecond addition
// https://github.com/php/php-src/issues/9106

func TestIntervalGh9106a(t *testing.T) {
	result := testAddWall(t, "2020-01-01 00:00:00.000000", "UTC", "PT1S", 500000, false)
	checkInterval(t, result, 2020, 1, 1, 0, 0, 1, 500000, 1577836801)
}

func TestIntervalGh9106b(t *testing.T) {
	result := testAddWall(t, "2020-01-01 00:00:01.500000", "UTC", "PT1S", 500000, false)
	checkInterval(t, result, 2020, 1, 1, 0, 0, 3, 0, 1577836803)
}

func TestIntervalGh9106c(t *testing.T) {
	result := testAddWall(t, "2020-01-01 00:00:03.600000", "UTC", "PT1S", 500000, false)
	checkInterval(t, result, 2020, 1, 1, 0, 0, 5, 100000, 1577836805)
}

// GitHub issue 8860: DST transitions in Europe/Amsterdam
// https://github.com/php/php-src/issues/8860

func TestIntervalGh8860a(t *testing.T) {
	result := testAddWall(t, "2022-10-30 01:00:00", "Europe/Amsterdam", "PT0H", 0, false)
	checkInterval(t, result, 2022, 10, 30, 1, 0, 0, 0, 1667084400)
	if result.Z != 7200 {
		t.Errorf("Expected z=7200, got %d", result.Z)
	}
}

func TestIntervalGh8860b(t *testing.T) {
	result := testAddWall(t, "2022-10-30 01:00:00", "Europe/Amsterdam", "PT1H", 0, false)
	checkInterval(t, result, 2022, 10, 30, 2, 0, 0, 0, 1667088000)
	if result.Z != 7200 {
		t.Errorf("Expected z=7200, got %d", result.Z)
	}
}

func TestIntervalGh8860c(t *testing.T) {
	result := testAddWall(t, "2022-10-30 02:00:00", "Europe/Amsterdam", "PT1H", 0, false)
	checkInterval(t, result, 2022, 10, 30, 3, 0, 0, 0, 1667095200)
	if result.Z != 3600 {
		t.Errorf("Expected z=3600, got %d", result.Z)
	}
}

func TestIntervalGh8860d(t *testing.T) {
	result := testAddWall(t, "2022-10-30 02:00:00", "Europe/Amsterdam", "PT3600S", 0, false)
	checkInterval(t, result, 2022, 10, 30, 3, 0, 0, 0, 1667095200)
	if result.Z != 3600 {
		t.Errorf("Expected z=3600, got %d", result.Z)
	}
}

// Long interval tests

func TestIntervalLongPositive(t *testing.T) {
	result := testAddWall(t, "2000-01-01 00:00:01.500000", "UTC", "P135000D", 0, false)
	checkInterval(t, result, 2369, 8, 14, 0, 0, 1, 500000, 12610684801)
}

func TestIntervalLongNegative(t *testing.T) {
	result := testAddWall(t, "2000-01-01 00:00:01.500000", "UTC", "P135000D", 0, true)
	checkInterval(t, result, 1630, 5, 20, 0, 0, 1, 500000, -10717315199)
}

// Leap year crossing tests

func TestIntervalLeapYearCrossing1(t *testing.T) {
	result := testAddWall(t, "2000-02-28 00:00:01.500000", "UTC", "P1D", 0, false)
	checkInterval(t, result, 2000, 2, 29, 0, 0, 1, 500000, 951782401)
}

func TestIntervalLeapYearCrossing2(t *testing.T) {
	result := testAddWall(t, "2000-02-29 00:00:01.500000", "UTC", "P1D", 0, true)
	checkInterval(t, result, 2000, 2, 28, 0, 0, 1, 500000, 951696001)
}

func TestIntervalLeapYearCrossing3(t *testing.T) {
	result := testAddWall(t, "2000-01-28 00:00:01.500000", "UTC", "P32D", 0, false)
	checkInterval(t, result, 2000, 2, 29, 0, 0, 1, 500000, 951782401)
}

func TestIntervalLeapYearCrossing4(t *testing.T) {
	result := testAddWall(t, "2000-02-29 00:00:01.500000", "UTC", "P32D", 0, true)
	checkInterval(t, result, 2000, 1, 28, 0, 0, 1, 500000, 949017601)
}

// Non-leap year crossing tests

func TestIntervalNonLeapYearCrossing1(t *testing.T) {
	result := testAddWall(t, "2001-02-28 00:00:01.500000", "UTC", "P1D", 0, false)
	checkInterval(t, result, 2001, 3, 1, 0, 0, 1, 500000, 983404801)
}

func TestIntervalNonLeapYearCrossing2(t *testing.T) {
	result := testAddWall(t, "2001-03-01 00:00:01.500000", "UTC", "P1D", 0, true)
	checkInterval(t, result, 2001, 2, 28, 0, 0, 1, 500000, 983318401)
}

func TestIntervalNonLeapYearCrossing3(t *testing.T) {
	result := testAddWall(t, "2001-01-28 00:00:01.500000", "UTC", "P32D", 0, false)
	checkInterval(t, result, 2001, 3, 1, 0, 0, 1, 500000, 983404801)
}

func TestIntervalNonLeapYearCrossing4(t *testing.T) {
	result := testAddWall(t, "2001-03-01 00:00:01.500000", "UTC", "P32D", 0, true)
	checkInterval(t, result, 2001, 1, 28, 0, 0, 1, 500000, 980640001)
}
