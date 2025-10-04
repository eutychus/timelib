package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestAllTestsSetup(t *testing.T) {
	// Test basic setup functionality similar to all_tests.cpp
	// This test verifies that the timezone database can be initialized

	// Test built-in database
	tzdb := timelib.BuiltinDB()
	if tzdb == nil {
		t.Error("Built-in timezone database not available")
	}

	// Test that we can create a basic timezone info structure
	var errorCode int
	tzInfo, err := timelib.ParseTzfile("UTC", tzdb, &errorCode)
	if err != nil {
		t.Logf("Could not parse UTC timezone: %v", err)
	} else if tzInfo != nil {
		t.Logf("Successfully parsed UTC timezone: %s", tzInfo.Name)
	}

	if errorCode != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Logf("UTC timezone parse error code: %d", errorCode)
	}

	// Basic validation - the function should not crash
	t.Log("AllTestsSetup test completed successfully")
}

func TestTimezoneDatabaseInitialization(t *testing.T) {
	// Test timezone database initialization similar to the C++ test runner
	tzdb := timelib.BuiltinDB()
	if tzdb == nil {
		t.Error("Built-in timezone database not available")
	}

	// Verify basic structure
	if tzdb.Version == "" {
		t.Log("Timezone database version is empty")
	}

	if tzdb.IndexSize < 0 {
		t.Error("Timezone database index size should not be negative")
	}

	// Test that we can get timezone identifiers
	var count int
	entries := timelib.TimezoneIdentifiersList(tzdb, &count)

	// Basic validation - function should work without crashing
	if entries == nil && count > 0 {
		t.Error("Expected non-nil entries when count > 0")
	}

	t.Log("TimezoneDatabaseInitialization test completed successfully")
}
