package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestWarnOnSlim001(t *testing.T) {
	var errorCode int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("New_York_Slim", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if errorCode != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", errorCode)
	}

	_ = tzi // Use the variable
}

func TestDontWarnOnFat001(t *testing.T) {
	var errorCode int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("New_York_Fat", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if errorCode != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", errorCode)
	}

	_ = tzi // Use the variable
}

func TestTzif4Format(t *testing.T) {
	var errorCode int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("Nicosia_TZif4", testDirectory, &errorCode)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if errorCode != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", errorCode)
	}

	if tzi.Bit64.Leapcnt != 20 {
		t.Errorf("Expected leapcnt 20, got %d", tzi.Bit64.Leapcnt)
	}
}
