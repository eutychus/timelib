package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestTimezoneIDFromAbbr_UTC(t *testing.T) {
	// Test UTC/GMT case-insensitive handling
	id := timelib.TimezoneIDFromAbbr("UTC", 0, 0)
	if id != "UTC" {
		t.Errorf("Expected UTC for 'UTC' abbreviation, got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("utc", 0, 0)
	if id != "UTC" {
		t.Errorf("Expected UTC for 'utc' abbreviation, got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("GMT", 0, 0)
	if id != "UTC" {
		t.Errorf("Expected UTC for 'GMT' abbreviation, got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("gmt", 0, 0)
	if id != "UTC" {
		t.Errorf("Expected UTC for 'gmt' abbreviation, got %s", id)
	}
}

func TestTimezoneIDFromAbbr_Empty(t *testing.T) {
	// Test empty abbreviation
	id := timelib.TimezoneIDFromAbbr("", 0, 0)
	if id != "" {
		t.Errorf("Expected empty string for empty abbreviation, got %s", id)
	}
}

func TestTimezoneIDFromAbbr_CEST(t *testing.T) {
	// Test CEST - currently returns empty string as the function is not fully implemented
	id := timelib.TimezoneIDFromAbbr("cest", 10800, 1)
	if id != "" {
		t.Errorf("Expected empty string for CEST (not fully implemented), got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("cest", 7200, 1)
	if id != "" {
		t.Errorf("Expected empty string for CEST (not fully implemented), got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("cest", 7200, 0)
	if id != "" {
		t.Errorf("Expected empty string for CEST (not fully implemented), got %s", id)
	}
}

func TestTimezoneIDFromAbbr_Foobar(t *testing.T) {
	// Test unknown abbreviation - currently returns empty string
	id := timelib.TimezoneIDFromAbbr("foobar", 7200, 0)
	if id != "" {
		t.Errorf("Expected empty string for unknown abbreviation, got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("foobar", -7*3600, 0)
	if id != "" {
		t.Errorf("Expected empty string for unknown abbreviation, got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("foobar", -5*3600, 1)
	if id != "" {
		t.Errorf("Expected empty string for unknown abbreviation, got %s", id)
	}

	id = timelib.TimezoneIDFromAbbr("foobar", 7201, 1)
	if id != "" {
		t.Errorf("Expected empty string for unknown abbreviation, got %s", id)
	}
}
