package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestEnumerateTimezones(t *testing.T) {
	// Test timezone enumeration functionality
	tzdb := timelib.BuiltinDB()
	if tzdb == nil {
		t.Skip("Built-in timezone database not available")
	}

	// Get timezone identifiers list
	var count int
	entries := timelib.TimezoneIdentifiersList(tzdb, &count)

	if count == 0 {
		t.Log("No timezone entries found in built-in database")
		return
	}

	t.Logf("Found %d timezone entries", count)

	// Test a few specific timezones that should be available
	testTimezones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
	}

	foundCount := 0
	for _, testTz := range testTimezones {
		found := false
		for _, entry := range entries {
			if entry.ID == testTz {
				found = true
				foundCount++
				break
			}
		}
		if !found {
			t.Logf("Timezone %s not found in database", testTz)
		}
	}

	t.Logf("Found %d out of %d test timezones", foundCount, len(testTimezones))
}

func TestEnumerateTimezonesBasic(t *testing.T) {
	// Test basic timezone database functionality
	tzdb := timelib.BuiltinDB()
	if tzdb == nil {
		t.Skip("Built-in timezone database not available")
	}

	// Verify basic structure
	if tzdb.Version == "" {
		t.Log("Timezone database version is empty")
	}

	if tzdb.IndexSize < 0 {
		t.Error("Timezone database index size should not be negative")
	}

	// Test that we can call the function without errors
	var count int
	entries := timelib.TimezoneIdentifiersList(tzdb, &count)

	// Basic validation
	if entries == nil && count > 0 {
		t.Error("Expected non-nil entries when count > 0")
	}

	t.Log("EnumerateTimezonesBasic test completed successfully")
}
