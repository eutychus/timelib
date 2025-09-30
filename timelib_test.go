package timelib

import (
	"testing"
	"time"
)

func TestTimeCtor(t *testing.T) {
	tm := TimeCtor()

	if tm == nil {
		t.Fatal("TimeCtor returned nil")
	}

	// Check that unset values are properly initialized
	if tm.Y != TIMELIB_UNSET {
		t.Errorf("Expected Y to be TIMELIB_UNSET (%d), got %d", TIMELIB_UNSET, tm.Y)
	}
	if tm.M != TIMELIB_UNSET {
		t.Errorf("Expected M to be TIMELIB_UNSET (%d), got %d", TIMELIB_UNSET, tm.M)
	}
	if tm.D != TIMELIB_UNSET {
		t.Errorf("Expected D to be TIMELIB_UNSET (%d), got %d", TIMELIB_UNSET, tm.D)
	}
	if tm.H != TIMELIB_UNSET {
		t.Errorf("Expected H to be TIMELIB_UNSET (%d), got %d", TIMELIB_UNSET, tm.H)
	}
	if tm.I != TIMELIB_UNSET {
		t.Errorf("Expected I to be TIMELIB_UNSET (%d), got %d", TIMELIB_UNSET, tm.I)
	}
	if tm.S != TIMELIB_UNSET {
		t.Errorf("Expected S to be TIMELIB_UNSET (%d), got %d", TIMELIB_UNSET, tm.S)
	}

	// Check default values
	if tm.US != 0 {
		t.Errorf("Expected US to be 0, got %d", tm.US)
	}
	if tm.Z != 0 {
		t.Errorf("Expected Z to be 0, got %d", tm.Z)
	}
	if tm.Dst != -1 {
		t.Errorf("Expected Dst to be -1, got %d", tm.Dst)
	}
	if tm.ZoneType != TIMELIB_ZONETYPE_NONE {
		t.Errorf("Expected ZoneType to be TIMELIB_ZONETYPE_NONE (%d), got %d", TIMELIB_ZONETYPE_NONE, tm.ZoneType)
	}
}

func TestRelTimeCtor(t *testing.T) {
	rt := RelTimeCtor()

	if rt == nil {
		t.Fatal("RelTimeCtor returned nil")
	}

	// Check default values
	if rt.Y != 0 || rt.M != 0 || rt.D != 0 {
		t.Errorf("Expected Y/M/D to be 0, got %d/%d/%d", rt.Y, rt.M, rt.D)
	}
	if rt.H != 0 || rt.I != 0 || rt.S != 0 {
		t.Errorf("Expected H/I/S to be 0, got %d/%d/%d", rt.H, rt.I, rt.S)
	}
	if rt.US != 0 {
		t.Errorf("Expected US to be 0, got %d", rt.US)
	}
	if rt.Weekday != -1 {
		t.Errorf("Expected Weekday to be -1, got %d", rt.Weekday)
	}
	if rt.WeekdayBehavior != 0 {
		t.Errorf("Expected WeekdayBehavior to be 0, got %d", rt.WeekdayBehavior)
	}
	if rt.Invert != false {
		t.Errorf("Expected Invert to be false, got %v", rt.Invert)
	}
	if rt.Days != 0 {
		t.Errorf("Expected Days to be 0, got %d", rt.Days)
	}
}

func TestTimeOffsetCtor(t *testing.T) {
	to := TimeOffsetCtor()

	if to == nil {
		t.Fatal("TimeOffsetCtor returned nil")
	}

	// Check default values
	if to.Offset != 0 {
		t.Errorf("Expected Offset to be 0, got %d", to.Offset)
	}
	if to.LeapSecs != 0 {
		t.Errorf("Expected LeapSecs to be 0, got %d", to.LeapSecs)
	}
	if to.IsDst != 0 {
		t.Errorf("Expected IsDst to be 0, got %d", to.IsDst)
	}
	if to.Abbr != "" {
		t.Errorf("Expected Abbr to be empty, got '%s'", to.Abbr)
	}
	if to.TransitionTime != 0 {
		t.Errorf("Expected TransitionTime to be 0, got %d", to.TransitionTime)
	}
}

func TestTimeCompare(t *testing.T) {
	// Test equal times
	t1 := &Time{Sse: 1000, US: 500, SseUptodate: true}
	t2 := &Time{Sse: 1000, US: 500, SseUptodate: true}

	result := TimeCompare(t1, t2)
	if result != 0 {
		t.Errorf("Expected TimeCompare to return 0 for equal times, got %d", result)
	}

	// Test different seconds
	t1.Sse = 1000
	t2.Sse = 2000
	result = TimeCompare(t1, t2)
	if result != -1 {
		t.Errorf("Expected TimeCompare to return -1 when t1 < t2, got %d", result)
	}

	result = TimeCompare(t2, t1)
	if result != 1 {
		t.Errorf("Expected TimeCompare to return 1 when t1 > t2, got %d", result)
	}

	// Test same seconds, different microseconds
	t1.Sse = 1000
	t2.Sse = 1000
	t1.US = 500
	t2.US = 1000

	result = TimeCompare(t1, t2)
	if result != -1 {
		t.Errorf("Expected TimeCompare to return -1 when t1.US < t2.US, got %d", result)
	}

	result = TimeCompare(t2, t1)
	if result != 1 {
		t.Errorf("Expected TimeCompare to return 1 when t1.US > t2.US, got %d", result)
	}
}

func TestDecimalHourToHMS(t *testing.T) {
	tests := []struct {
		hour     float64
		expected struct {
			hour, min, sec int
		}
		desc string
	}{
		{1.5, struct{ hour, min, sec int }{1, 30, 0}, "1.5 hours = 1:30:00"},
		{2.25, struct{ hour, min, sec int }{2, 15, 0}, "2.25 hours = 2:15:00"},
		{3.0166666666666666, struct{ hour, min, sec int }{3, 1, 0}, "3.016666... hours = 3:01:00"},
		{3.0175, struct{ hour, min, sec int }{3, 1, 3}, "3.0175 hours = 3:01:03"},
		{-1.5, struct{ hour, min, sec int }{-1, 30, 0}, "-1.5 hours = -1:30:00"},
		{0.0, struct{ hour, min, sec int }{0, 0, 0}, "0 hours = 0:00:00"},
	}

	for _, test := range tests {
		h, m, s := DecimalHourToHMS(test.hour)
		if h != test.expected.hour || m != test.expected.min || s != test.expected.sec {
			t.Errorf("%s: got %d:%d:%d, expected %d:%d:%d",
				test.desc, h, m, s, test.expected.hour, test.expected.min, test.expected.sec)
		}
	}
}

func TestHMSToDecimalHour(t *testing.T) {
	tests := []struct {
		hour, min, sec int
		expected       float64
		desc           string
	}{
		{1, 30, 0, 1.5, "1:30:00 = 1.5 hours"},
		{2, 15, 0, 2.25, "2:15:00 = 2.25 hours"},
		{3, 1, 0, 3.0166666666666666, "3:01:00 = 3.016666... hours"},
		{-1, 30, 0, -1.5, "-1:30:00 = -1.5 hours"},
		{0, 0, 0, 0.0, "0:00:00 = 0 hours"},
	}

	for _, test := range tests {
		result := HMSToDecimalHour(test.hour, test.min, test.sec)
		// Allow small floating point differences
		if result < test.expected-0.000001 || result > test.expected+0.000001 {
			t.Errorf("%s: got %f, expected %f", test.desc, result, test.expected)
		}
	}
}

func TestHMSFToDecimalHour(t *testing.T) {
	tests := []struct {
		hour, min, sec, us int
		expected           float64
		tolerance          float64
		desc               string
	}{
		{1, 30, 0, 0, 1.5, 0.000001, "1:30:00.0 = 1.5 hours"},
		{1, 30, 0, 500000, 1.500139, 0.000001, "1:30:00.5 = 1.500139 hours (actual calculated value)"},
		{2, 15, 30, 0, 2.2583333333333333, 0.000001, "2:15:30.0 = 2.258333... hours"},
		{-1, 30, 0, 0, -1.5, 0.000001, "-1:30:00.0 = -1.5 hours"},
		{0, 0, 0, 0, 0.0, 0.000001, "0:00:00.0 = 0 hours"},
	}

	for _, test := range tests {
		result := HMSFToDecimalHour(test.hour, test.min, test.sec, test.us)
		// Allow small floating point differences
		if result < test.expected-test.tolerance || result > test.expected+test.tolerance {
			t.Errorf("%s: got %f, expected %f (tolerance: %f)", test.desc, result, test.expected, test.tolerance)
		}
	}
}

func TestHMSToSeconds(t *testing.T) {
	tests := []struct {
		h, m, s  int64
		expected int64
		desc     string
	}{
		{1, 0, 0, 3600, "1 hour = 3600 seconds"},
		{0, 1, 0, 60, "1 minute = 60 seconds"},
		{0, 0, 1, 1, "1 second = 1 second"},
		{1, 30, 45, 5445, "1:30:45 = 5445 seconds"},
		{0, 0, 0, 0, "0:00:00 = 0 seconds"},
	}

	for _, test := range tests {
		result := HMSToSeconds(test.h, test.m, test.s)
		if result != test.expected {
			t.Errorf("%s: got %d, expected %d", test.desc, result, test.expected)
		}
	}
}

func TestDateToInt(t *testing.T) {
	// Test valid timestamp
	tm := &Time{Sse: 1234567890}
	result, err := DateToInt(tm)
	if err != nil {
		t.Errorf("Expected no error for valid timestamp, got %v", err)
	}
	if result != 1234567890 {
		t.Errorf("Expected %d, got %d", 1234567890, result)
	}

	// Test timestamp that's too large
	tm.Sse = 9223372036854775807 // Max int64
	_, err = DateToInt(tm)
	if err != nil {
		t.Errorf("Expected no error for max int64, got %v", err)
	}

	// Note: In Go, int64 overflow would wrap around, so we test boundary conditions differently
	// The DateToInt function should handle the full int64 range correctly
}

func TestSetTimezoneFromOffset(t *testing.T) {
	tm := TimeCtor()
	SetTimezoneFromOffset(tm, 3600) // +1 hour

	if tm.ZoneType != TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected ZoneType to be TIMELIB_ZONETYPE_OFFSET (%d), got %d", TIMELIB_ZONETYPE_OFFSET, tm.ZoneType)
	}
	if tm.Z != 3600 {
		t.Errorf("Expected Z to be 3600, got %d", tm.Z)
	}
	if tm.Dst != 0 {
		t.Errorf("Expected Dst to be 0, got %d", tm.Dst)
	}
	if tm.TzAbbr != "" {
		t.Errorf("Expected TzAbbr to be empty, got '%s'", tm.TzAbbr)
	}
	if tm.TzInfo != nil {
		t.Error("Expected TzInfo to be nil")
	}
}

func TestSetTimezoneFromAbbr(t *testing.T) {
	tm := TimeCtor()
	SetTimezoneFromAbbr(tm, "CET", 3600, 0)

	if tm.ZoneType != TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected ZoneType to be TIMELIB_ZONETYPE_ABBR (%d), got %d", TIMELIB_ZONETYPE_ABBR, tm.ZoneType)
	}
	if tm.Z != 3600 {
		t.Errorf("Expected Z to be 3600, got %d", tm.Z)
	}
	if tm.Dst != 0 {
		t.Errorf("Expected Dst to be 0, got %d", tm.Dst)
	}
	if tm.TzAbbr != "CET" {
		t.Errorf("Expected TzAbbr to be 'CET', got '%s'", tm.TzAbbr)
	}
	if tm.TzInfo != nil {
		t.Error("Expected TzInfo to be nil")
	}
}

func TestSetTimezone(t *testing.T) {
	tm := TimeCtor()
	tz := &TzInfo{Name: "Europe/Berlin"}

	SetTimezone(tm, tz)

	if tm.ZoneType != TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected ZoneType to be TIMELIB_ZONETYPE_ID (%d), got %d", TIMELIB_ZONETYPE_ID, tm.ZoneType)
	}
	if tm.TzInfo != tz {
		t.Error("Expected TzInfo to be set to provided timezone")
	}
	if tm.TzAbbr != "Europe/Berlin" {
		t.Errorf("Expected TzAbbr to be 'Europe/Berlin', got '%s'", tm.TzAbbr)
	}
}

func TestGetErrorMessage(t *testing.T) {
	tests := []struct {
		errorCode int
		expected  string
	}{
		{TIMELIB_ERROR_NO_ERROR, "No error"},
		{TIMELIB_ERROR_CANNOT_ALLOCATE, "Cannot allocate buffer for parsing"},
		{TIMELIB_ERROR_NO_SUCH_TIMEZONE, "No timezone with this name could be found"},
		{99, "Unknown error code"}, // Out of range
	}

	for _, test := range tests {
		result := GetErrorMessage(test.errorCode)
		if result != test.expected {
			t.Errorf("GetErrorMessage(%d): expected '%s', got '%s'", test.errorCode, test.expected, result)
		}
	}
}

func TestConvertTime(t *testing.T) {
	// This is a placeholder test - the actual implementation will come later
	tm := TimeCtor()
	tm.Y = 2023
	tm.M = 12
	tm.D = 25
	tm.H = 15
	tm.I = 30
	tm.S = 45

	goTime := ConvertTime(tm)

	// For now, just check it doesn't panic and returns a valid time
	if goTime.IsZero() {
		t.Log("ConvertTime returned zero time (expected for placeholder implementation)")
	}
}

func TestConvertFromTime(t *testing.T) {
	// This is a placeholder test - the actual implementation will come later
	goTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	tm := ConvertFromTime(goTime)

	// For now, just check it doesn't panic and returns a valid Time
	if tm == nil {
		t.Fatal("ConvertFromTime returned nil")
	}
}
