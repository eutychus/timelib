package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestGetCurrentOffsetUTCPlus0500(t *testing.T) {
	time := &timelib.Time{}
	timelib.SetTimezoneFromOffset(time, 5*3600)

	time.Unixtime2local(1483280063)
	offset := timelib.GetCurrentOffset(time)

	if offset != 18000 {
		t.Errorf("Expected offset 18000, got %d", offset)
	}
}

func TestGetCurrentOffsetEST(t *testing.T) {
	time := &timelib.Time{}
	timelib.SetTimezoneFromAbbr(time, "EST", -5*3600, 0)

	time.Unixtime2local(1483280063)
	offset := timelib.GetCurrentOffset(time)

	if offset != -18000 {
		t.Errorf("Expected offset -18000, got %d", offset)
	}
}

func TestGetCurrentOffsetEDT(t *testing.T) {
	time := &timelib.Time{}
	timelib.SetTimezoneFromAbbr(time, "EDT", -5*3600, 1)

	time.Unixtime2local(1483280063)
	offset := timelib.GetCurrentOffset(time)

	if offset != -14400 {
		t.Errorf("Expected offset -14400, got %d", offset)
	}
}

func TestGetCurrentOffsetLondon(t *testing.T) {
	var code int
	tzi, err := timelib.ParseTzfile("Europe/London", timelib.BuiltinDB(), &code)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	time := &timelib.Time{}
	timelib.SetTimezone(time, tzi)

	time.Unixtime2local(1483280063)
	offset := timelib.GetCurrentOffset(time)
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}

	time.Unixtime2local(1501074654)
	offset = timelib.GetCurrentOffset(time)
	if offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", offset)
	}
}

func TestGetCurrentOffsetAmsterdam(t *testing.T) {
	var code int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &code)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	time := &timelib.Time{}
	timelib.SetTimezone(time, tzi)

	time.Unixtime2local(1483280063)
	offset := timelib.GetCurrentOffset(time)
	if offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", offset)
	}

	time.Unixtime2local(1501074654)
	offset = timelib.GetCurrentOffset(time)
	if offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", offset)
	}
}
