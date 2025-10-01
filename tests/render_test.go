package tests

import (
	"fmt"
	"testing"

	timelib "github.com/eutychus/timelib"
)

func renderDate(d *timelib.Time) string {
	var yearSign string
	year := d.Y
	if year < 0 {
		yearSign = "-"
		year = -year
	}

	zoneStr := "??"
	if d.ZoneType == timelib.TIMELIB_ZONETYPE_OFFSET {
		zoneStr = "GMT"
	} else if d.TzAbbr != "" {
		zoneStr = d.TzAbbr
	}

	return fmt.Sprintf("%s%04d-%02d-%02d %02d:%02d:%02d.%06d %s",
		yearSign, year, d.M, d.D, d.H, d.I, d.S, d.US, zoneStr)
}

func testRender(t *testing.T, testName, expected string, ts int64, tzid string, useBuiltinDB bool) {
	var dummyError int
	var tzi *timelib.TzInfo
	var err error

	if useBuiltinDB {
		tzi, err = timelib.ParseTzfile(tzid, timelib.BuiltinDB(), &dummyError)
	} else {
		tzdb, err2 := timelib.Zoneinfo("")
		if err2 != nil {
			t.Fatalf("%s: Zoneinfo error: %v", testName, err2)
		}
		tzi, err = timelib.ParseTzfile(tzid, tzdb, &dummyError)
	}

	if err != nil {
		t.Fatalf("%s: Parse error: %v", testName, err)
	}
	if dummyError != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Fatalf("%s: Expected TIMELIB_ERROR_NO_ERROR, got %d", testName, dummyError)
	}
	if tzi == nil {
		t.Fatalf("%s: tzi is nil", testName)
	}

	time := &timelib.Time{}
	timelib.SetTimezone(time, tzi)
	time.Unixtime2local(ts)

	rendered := renderDate(time)
	if rendered != expected {
		t.Errorf("%s: Expected %s, got %s", testName, expected, rendered)
	}
}

func TestRenderPhpBug30532_01(t *testing.T) {
	testRender(t, "php_bug_30532_01", "2004-10-31 00:00:00.000000 EDT", 1099195200, "America/New_York", true)
}

func TestRenderPhpBug30532_02(t *testing.T) {
	testRender(t, "php_bug_30532_02", "2004-10-31 01:00:00.000000 EDT", 1099198800, "America/New_York", true)
}

func TestRenderPhpBug30532_03(t *testing.T) {
	testRender(t, "php_bug_30532_03", "2004-10-31 02:00:00.000000 EST", 1099206000, "America/New_York", true)
}

func TestRenderPhpBug32086_01(t *testing.T) {
	testRender(t, "php_bug_32086_01", "2004-11-01 00:00:00.000000 -03", 1099278000, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_02(t *testing.T) {
	testRender(t, "php_bug_32086_02", "2004-11-01 23:00:00.000000 -03", 1099360800, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_03(t *testing.T) {
	testRender(t, "php_bug_32086_03", "2004-11-01 23:59:00.000000 -03", 1099364340, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_04(t *testing.T) {
	testRender(t, "php_bug_32086_04", "2004-11-02 01:00:00.000000 -02", 1099364400, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_05(t *testing.T) {
	testRender(t, "php_bug_32086_05", "2004-11-02 02:00:00.000000 -02", 1099368000, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_06(t *testing.T) {
	testRender(t, "php_bug_32086_06", "2005-02-18 23:00:00.000000 -02", 1108774800, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_07(t *testing.T) {
	testRender(t, "php_bug_32086_07", "2005-02-19 00:00:00.000000 -02", 1108778400, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_08(t *testing.T) {
	testRender(t, "php_bug_32086_08", "2005-02-19 01:00:00.000000 -02", 1108782000, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32086_09(t *testing.T) {
	testRender(t, "php_bug_32086_09", "2005-02-20 00:00:00.000000 -03", 1108868400, "America/Sao_Paulo", true)
}

func TestRenderPhpBug32555_01(t *testing.T) {
	testRender(t, "php_bug_32555_01", "2005-04-02 02:30:00.000000 EST", 1112427000, "America/New_York", true)
}

func TestRenderPhpBug32555_02(t *testing.T) {
	testRender(t, "php_bug_32555_02", "2005-04-02 00:00:00.000000 EST", 1112418000, "America/New_York", true)
}

func TestRenderPhpBug32555_03(t *testing.T) {
	testRender(t, "php_bug_32555_03", "2005-04-03 00:00:00.000000 EST", 1112504400, "America/New_York", true)
}

func TestRenderPhpBug32588_01(t *testing.T) {
	testRender(t, "php_bug_32588_01", "2005-04-02 00:00:00.000000 GMT", 1112400000, "GMT", true)
}

func TestRenderPhpBug32588_02(t *testing.T) {
	testRender(t, "php_bug_32588_02", "2005-04-03 00:00:00.000000 GMT", 1112486400, "GMT", true)
}

func TestRenderPhpBug32588_03(t *testing.T) {
	testRender(t, "php_bug_32588_03", "2005-04-04 00:00:00.000000 GMT", 1112572800, "GMT", true)
}

func TestRenderPhpBug73294_01(t *testing.T) {
	testRender(t, "php_bug_73294_01", "-1900-06-22 00:00:00.000000 UTC", -122110502400, "UTC", true)
}

func TestRenderPhpBug73294_02(t *testing.T) {
	testRender(t, "php_bug_73294_02", "-1916-06-22 00:00:00.000000 UTC", -122615337600, "UTC", true)
}

func TestRenderFirstTransition01(t *testing.T) {
	testRender(t, "first_transition_01", "1884-08-06 06:39:33.000000 PST", -2695022427, "America/Los_Angeles", true)
}

func TestRenderFirstTransition02(t *testing.T) {
	testRender(t, "first_transition_02", "1900-08-06 06:39:33.000000 PST", -2190187227, "America/Los_Angeles", true)
}

func TestRenderFirstTransition03(t *testing.T) {
	testRender(t, "first_transition_03", "1901-08-06 06:39:33.000000 PST", -2158651227, "America/Los_Angeles", true)
}

func TestRenderFirstTransition04(t *testing.T) {
	testRender(t, "first_transition_04", "1902-08-06 06:39:33.000000 PST", -2127115227, "America/Los_Angeles", true)
}

func TestRenderFirstTransition05(t *testing.T) {
	testRender(t, "first_transition_05", "1918-02-06 06:39:33.000000 PST", -1637832027, "America/Los_Angeles", true)
}

func TestRenderFirstTransition06(t *testing.T) {
	testRender(t, "first_transition_06", "1918-08-06 06:39:33.000000 PDT", -1622197227, "America/Los_Angeles", true)
}

func TestRenderIssue0017Render2369_01(t *testing.T) {
	testRender(t, "issue0017_render_2369_01", "2369-12-31 00:00:00.000000 UTC", 12622694400, "UTC", true)
}

func TestRenderPastLeap01(t *testing.T) {
	testRender(t, "past_leap_01", "1965-01-01 00:00:00.000000 UTC", -157766400, "UTC", true)
}

func TestRenderPastLeap02(t *testing.T) {
	testRender(t, "past_leap_02", "1964-12-31 00:00:00.000000 UTC", -157852800, "UTC", true)
}

func TestRenderPastLeap03(t *testing.T) {
	testRender(t, "past_leap_03", "1964-01-31 00:00:00.000000 UTC", -186796800, "UTC", true)
}

func TestRenderPastLeap04(t *testing.T) {
	testRender(t, "past_leap_04", "1964-01-01 00:00:00.000000 UTC", -189388800, "UTC", true)
}

func TestRenderPastLeap05(t *testing.T) {
	testRender(t, "past_leap_05", "1963-12-31 00:00:00.000000 UTC", -189475200, "UTC", true)
}

func TestRender01(t *testing.T) {
	testRender(t, "render_01", "2005-05-26 23:11:59.000000 CEST", 1117141919, "Europe/Amsterdam", true)
}

func TestRender02(t *testing.T) {
	testRender(t, "render_02", "2005-05-26 22:11:59.000000 BST", 1117141919, "Europe/London", true)
}

func TestRender03(t *testing.T) {
	testRender(t, "render_03", "2005-05-27 07:11:59.000000 AEST", 1117141919, "Australia/Sydney", true)
}

func TestRender04(t *testing.T) {
	testRender(t, "render_04", "2005-10-30 00:00:00.000000 GMT", 1130630400, "GMT", true)
}

func TestRender05(t *testing.T) {
	testRender(t, "render_05", "2005-10-30 00:59:59.000000 GMT", 1130633999, "GMT", true)
}

func TestRender06(t *testing.T) {
	testRender(t, "render_06", "2005-10-30 01:00:00.000000 GMT", 1130634000, "GMT", true)
}

func TestRender07(t *testing.T) {
	testRender(t, "render_07", "2005-10-30 02:00:00.000000 CEST", 1130630400, "Europe/Oslo", true)
}

func TestRender08(t *testing.T) {
	testRender(t, "render_08", "2005-10-30 02:59:59.000000 CEST", 1130633999, "Europe/Oslo", true)
}

func TestRender09(t *testing.T) {
	testRender(t, "render_09", "2005-10-30 02:00:00.000000 CET", 1130634000, "Europe/Oslo", true)
}

func TestRender10(t *testing.T) {
	testRender(t, "render_10", "2005-10-30 01:00:00.000000 CEST", 1130626800, "Europe/Amsterdam", true)
}

func TestRender11(t *testing.T) {
	testRender(t, "render_11", "2005-10-30 02:00:00.000000 CEST", 1130630400, "Europe/Amsterdam", true)
}

func TestRender12(t *testing.T) {
	testRender(t, "render_12", "2005-10-30 02:00:00.000000 CET", 1130634000, "Europe/Amsterdam", true)
}

func TestRender13(t *testing.T) {
	testRender(t, "render_13", "2005-10-30 03:00:00.000000 CET", 1130637600, "Europe/Amsterdam", true)
}

func TestRender14(t *testing.T) {
	testRender(t, "render_14", "2005-10-30 04:00:00.000000 CET", 1130641200, "Europe/Amsterdam", true)
}

func TestRender15(t *testing.T) {
	testRender(t, "render_15", "2006-06-07 19:06:44.000000 CEST", 1149700004, "Europe/Amsterdam", true)
}

func TestRenderTransition01(t *testing.T) {
	testRender(t, "transition_01", "2008-03-30 01:00:00.000000 CET", 1206835200, "Europe/Amsterdam", true)
}

func TestRenderTransition02(t *testing.T) {
	testRender(t, "transition_02", "2008-03-30 01:59:59.000000 CET", 1206838799, "Europe/Amsterdam", true)
}

func TestRenderTransition03(t *testing.T) {
	testRender(t, "transition_03", "2008-03-30 03:00:00.000000 CEST", 1206838800, "Europe/Amsterdam", true)
}

func TestRenderTransition04(t *testing.T) {
	testRender(t, "transition_04", "2008-03-30 03:59:59.000000 CEST", 1206842399, "Europe/Amsterdam", true)
}

func TestRenderTransition05(t *testing.T) {
	testRender(t, "transition_05", "2008-03-30 03:00:00.000000 CEST", 1206838800, "Europe/Amsterdam", true)
}

func TestRenderWeirdTimezone01(t *testing.T) {
	testRender(t, "weird_timezone_01", "1970-01-04 21:47:17.000000 -0830", 368237, "Pacific/Pitcairn", true)
}

func TestRenderWeirdTimezone02(t *testing.T) {
	testRender(t, "weird_timezone_02", "1998-04-27 08:29:59.000000 UTC", 893665799, "UTC", true)
}

func TestRenderWeirdTimezone03(t *testing.T) {
	testRender(t, "weird_timezone_03", "1998-04-27 08:30:00.000000 UTC", 893665800, "UTC", true)
}

func TestRenderWeirdTimezone04(t *testing.T) {
	testRender(t, "weird_timezone_04", "1998-04-27 08:30:01.000000 UTC", 893665801, "UTC", true)
}

func TestRenderWeirdTimezone05(t *testing.T) {
	testRender(t, "weird_timezone_05", "1998-04-26 23:59:59.000000 -0830", 893665799, "Pacific/Pitcairn", true)
}

func TestRenderWeirdTimezone06(t *testing.T) {
	testRender(t, "weird_timezone_06", "1998-04-27 00:30:00.000000 -08", 893665800, "Pacific/Pitcairn", true)
}

func TestRenderWeirdTimezone07(t *testing.T) {
	testRender(t, "weird_timezone_07", "1998-04-27 00:30:01.000000 -08", 893665801, "Pacific/Pitcairn", true)
}
