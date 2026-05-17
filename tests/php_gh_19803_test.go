package tests

import (
	"testing"
	"github.com/eutychus/timelib"
)

func TestPhpGh19803_03(t *testing.T) {
	time, errors, err := timelib.ParseDateString(" \t\n \r ", timelib.BuiltinDB(), timelib.ParseTzfile)
	if err != nil {
		t.Fatalf("ParseDateString returned error: %v", err)
	}
	if errors == nil || errors.ErrorCount != 1 {
		t.Fatalf("Expected 1 error, got %v", errors)
	}
	if errors.ErrorMessages[0].ErrorCode != timelib.TIMELIB_ERR_EMPTY_STRING {
		t.Errorf("Expected TIMELIB_ERR_EMPTY_STRING, got %d", errors.ErrorMessages[0].ErrorCode)
	}
	if time == nil || time.Y != timelib.TIMELIB_UNSET {
		t.Errorf("Expected TIMELIB_UNSET for year, got %v", time)
	}
}

func TestPhpGh19803_04(t *testing.T) {
	// A single whitespace character
	time, errors, err := timelib.ParseDateString(" ", timelib.BuiltinDB(), timelib.ParseTzfile)
	if err != nil {
		t.Fatalf("ParseDateString returned error: %v", err)
	}
	if errors == nil || errors.ErrorCount != 1 {
		t.Fatalf("Expected 1 error, got %v", errors)
	}
	if errors.ErrorMessages[0].ErrorCode != timelib.TIMELIB_ERR_EMPTY_STRING {
		t.Errorf("Expected TIMELIB_ERR_EMPTY_STRING, got %d", errors.ErrorMessages[0].ErrorCode)
	}
	if time == nil || time.Y != timelib.TIMELIB_UNSET {
		t.Errorf("Expected TIMELIB_UNSET for year, got %v", time)
	}
}
