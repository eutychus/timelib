package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestGetOffsetInfoUTCPlusOffset(t *testing.T) {
	var offset int32
	time := &timelib.Time{}
	timelib.SetTimezoneFromOffset(time, 5*3600)

	time.Unixtime2local(1483280063)
	result := timelib.GetTimeZoneOffsetInfo(time.Sse, time.TzInfo, &offset, nil, nil)

	if result != 0 {
		t.Errorf("Expected 0 for offset-based timezone, got %d", result)
	}
}

func TestGetOffsetInfoAbbreviatedTimeZone(t *testing.T) {
	var offset int32
	time := &timelib.Time{}
	timelib.SetTimezoneFromAbbr(time, "EST", -5*3600, 0)

	time.Unixtime2local(1483280063)
	result := timelib.GetTimeZoneOffsetInfo(time.Sse, time.TzInfo, &offset, nil, nil)

	if result != 0 {
		t.Errorf("Expected 0 for abbreviated timezone, got %d", result)
	}
}

func TestGetOffsetInfoLondon(t *testing.T) {
	var code int
	var offset int32
	var transition int64
	var isDst uint

	tzi, err := timelib.ParseTzfile("Europe/London", timelib.BuiltinDB(), &code)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	time := &timelib.Time{}
	timelib.SetTimezone(time, tzi)

	// 1483280063 = 2017-01-01
	time.Unixtime2local(1483280063)
	result := timelib.GetTimeZoneOffsetInfo(time.Sse, time.TzInfo, &offset, &transition, &isDst)

	if result == 0 {
		t.Errorf("Expected non-zero for London timezone, got 0")
	}
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}
	if isDst != 0 {
		t.Errorf("Expected is_dst false (0), got %d", isDst)
	}
	if transition != 1477789200 { // 1477789200 = 2016-10-30
		t.Errorf("Expected transition 1477789200, got %d", transition)
	}

	// 1501074654 = 2017-07-26
	time.Unixtime2local(1501074654)
	result = timelib.GetTimeZoneOffsetInfo(time.Sse, time.TzInfo, &offset, &transition, &isDst)

	if result == 0 {
		t.Errorf("Expected non-zero for London timezone, got 0")
	}
	if offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", offset)
	}
	if isDst == 0 {
		t.Errorf("Expected is_dst true (non-zero), got %d", isDst)
	}
	if transition != 1490490000 { // 1490490000 = 2017-03-26
		t.Errorf("Expected transition 1490490000, got %d", transition)
	}
}

func TestGetOffsetInfoAmsterdam(t *testing.T) {
	var code int
	var offset int32

	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &code)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	time := &timelib.Time{}
	timelib.SetTimezone(time, tzi)

	time.Unixtime2local(1483280063)
	result := timelib.GetTimeZoneOffsetInfo(time.Sse, time.TzInfo, &offset, nil, nil)

	if result == 0 {
		t.Errorf("Expected non-zero for Amsterdam timezone, got 0")
	}
	if offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", offset)
	}

	time.Unixtime2local(1501074654)
	result = timelib.GetTimeZoneOffsetInfo(time.Sse, time.TzInfo, &offset, nil, nil)

	if result == 0 {
		t.Errorf("Expected non-zero for Amsterdam timezone, got 0")
	}
	if offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", offset)
	}
}
