package timelib

import (
	"path/filepath"
	"testing"
)

// TestParseTzfileNewYorkSlim tests parsing a slim timezone file
func TestParseTzfileNewYorkSlim(t *testing.T) {
	path := filepath.Join("tests", "files", "New_York_Slim")

	var errorCode int
	tzi, err := ParseTzfileFromFile(path, &errorCode)

	if err != nil {
		t.Fatalf("Failed to parse timezone file: %v (error code: %d)", err, errorCode)
	}

	if tzi == nil {
		t.Fatal("Parsed timezone info is nil")
	}

	// Check basic properties
	if tzi.PosixString == "" {
		t.Error("Expected POSIX string to be set")
	}

	// Check that we have transitions
	if len(tzi.Trans) == 0 {
		t.Error("Expected transitions to be parsed")
	}

	// Check that we have types
	if len(tzi.Type) == 0 {
		t.Error("Expected types to be parsed")
	}

	t.Logf("Parsed timezone: %d transitions, %d types, POSIX: %s",
		len(tzi.Trans), len(tzi.Type), tzi.PosixString)
}

// TestParseTzfileNewYorkFat tests parsing a fat timezone file
func TestParseTzfileNewYorkFat(t *testing.T) {
	path := filepath.Join("tests", "files", "New_York_Fat")

	var errorCode int
	tzi, err := ParseTzfileFromFile(path, &errorCode)

	if err != nil {
		t.Fatalf("Failed to parse timezone file: %v (error code: %d)", err, errorCode)
	}

	if tzi == nil {
		t.Fatal("Parsed timezone info is nil")
	}

	// Fat file should have more transitions than slim
	if len(tzi.Trans) < 100 {
		t.Errorf("Expected more transitions in fat file, got %d", len(tzi.Trans))
	}

	t.Logf("Parsed fat timezone: %d transitions, %d types",
		len(tzi.Trans), len(tzi.Type))
}

// TestParseTzfileNonContinuous tests error handling for corrupt files
func TestParseTzfileNonContinuous(t *testing.T) {
	path := filepath.Join("tests", "files", "NonContinuous")

	var errorCode int
	tzi, err := ParseTzfileFromFile(path, &errorCode)

	// Should fail with CORRUPT_TRANSITIONS_DONT_INCREASE error
	if err == nil {
		t.Error("Expected error for non-continuous transitions")
	}

	if tzi != nil {
		t.Error("Expected nil timezone info for corrupt file")
	}

	if errorCode != TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE {
		t.Errorf("Expected error code %d, got %d",
			TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE, errorCode)
	}

	t.Logf("Correctly detected corrupt file: %v (error code: %d)", err, errorCode)
}

// TestParsePosixString tests POSIX TZ string parsing
func TestParsePosixString(t *testing.T) {
	tests := []struct {
		name         string
		posix        string
		expectStd    string
		expectDst    string
		expectStdOff int64
		expectDstOff int64
		expectError  bool
	}{
		{
			name:         "EST5EDT",
			posix:        "EST5EDT,M3.2.0,M11.1.0",
			expectStd:    "EST",
			expectDst:    "EDT",
			expectStdOff: -5 * 3600,
			expectDstOff: -4 * 3600,
		},
		{
			name:         "CST6CDT",
			posix:        "CST6CDT,M3.2.0,M11.1.0",
			expectStd:    "CST",
			expectDst:    "CDT",
			expectStdOff: -6 * 3600,
			expectDstOff: -5 * 3600,
		},
		{
			name:         "UTC0",
			posix:        "UTC0",
			expectStd:    "UTC",
			expectDst:    "",
			expectStdOff: 0,
			expectDstOff: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps, err := ParsePosixString(test.posix, nil)

			if test.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if ps.Std != test.expectStd {
				t.Errorf("Expected std=%s, got %s", test.expectStd, ps.Std)
			}

			if ps.Dst != test.expectDst {
				t.Errorf("Expected dst=%s, got %s", test.expectDst, ps.Dst)
			}

			if ps.StdOffset != test.expectStdOff {
				t.Errorf("Expected stdOffset=%d, got %d", test.expectStdOff, ps.StdOffset)
			}

			if test.expectDst != "" && ps.DstOffset != test.expectDstOff {
				t.Errorf("Expected dstOffset=%d, got %d", test.expectDstOff, ps.DstOffset)
			}
		})
	}
}

// TestGetOffsetInfo tests offset calculation
func TestGetOffsetInfo(t *testing.T) {
	path := filepath.Join("tests", "files", "New_York_Slim")

	var errorCode int
	tzi, err := ParseTzfileFromFile(path, &errorCode)
	if err != nil {
		t.Fatalf("Failed to parse timezone file: %v", err)
	}

	// Test a known timestamp: 2022-07-01 12:00:00 UTC = 1656676800
	// Should be in EDT (DST)
	ts := int64(1656676800)

	info, transTime, err := GetOffsetInfo(ts, tzi)
	if err != nil {
		t.Fatalf("Failed to get offset info: %v", err)
	}

	if info == nil {
		t.Fatal("Offset info is nil")
	}

	// July should be DST in New York
	if info.IsDst == 0 {
		t.Error("Expected DST flag to be set for July in New York")
	}

	// EDT offset should be -4 hours = -14400 seconds
	if info.Offset != -14400 {
		t.Errorf("Expected offset -14400 (EDT), got %d", info.Offset)
	}

	t.Logf("Offset for %d: %d seconds, DST=%d, transition=%d",
		ts, info.Offset, info.IsDst, transTime)
}

// TestZoneinfoDir tests loading timezone database from directory
func TestZoneinfoDir(t *testing.T) {
	dir := filepath.Join("tests", "files")

	tzdb, err := ZoneinfoDir(dir)
	if err != nil {
		t.Fatalf("Failed to load timezone directory: %v", err)
	}

	if tzdb == nil {
		t.Fatal("Timezone database is nil")
	}

	if len(tzdb.Index) == 0 {
		t.Error("Expected non-empty timezone index")
	}

	t.Logf("Loaded %d timezone files from directory", len(tzdb.Index))

	// Try to parse one of the loaded files
	if len(tzdb.Index) > 0 {
		tzName := tzdb.Index[0].ID
		t.Logf("Testing timezone: %s", tzName)

		// This won't work directly since we need the full path
		// But at least we verified the index was built
	}
}

// TestTimezoneValidation tests timezone ID validation
func TestTimezoneValidation(t *testing.T) {
	dir := filepath.Join("tests", "files")
	tzdb, err := ZoneinfoDir(dir)
	if err != nil {
		t.Fatalf("Failed to load timezone directory: %v", err)
	}

	// Check for a timezone that exists
	if len(tzdb.Index) > 0 {
		validTz := tzdb.Index[0].ID
		if !TimezoneIDIsValid(validTz, tzdb) {
			t.Errorf("Expected timezone %s to be valid", validTz)
		}
	}

	// Check for a timezone that doesn't exist
	if TimezoneIDIsValid("Invalid/Nonexistent", tzdb) {
		t.Error("Expected invalid timezone to be rejected")
	}
}
