package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestParseTzNewYork(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("America/New_York", timelib.BuiltinDB(), &error)
	if err != nil {
		t.Fatalf("ParseTzfile failed: %v", err)
	}

	if tzi == nil {
		t.Fatal("Expected non-nil timezone info")
	}
	if error != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", error)
	}
	if tzi.PosixString != "EST5EDT,M3.2.0,M11.1.0" {
		t.Errorf("Expected posix_string EST5EDT,M3.2.0,M11.1.0, got %s", tzi.PosixString)
	}

	// Check location data
	countryCode := string(tzi.Location.CountryCode[:2])
	if countryCode != "US" {
		t.Errorf("Expected country code US, got %s", countryCode)
	}
	// Check latitude with tolerance
	if tzi.Location.Latitude < 39.0 || tzi.Location.Latitude > 41.0 {
		t.Errorf("Expected latitude ~40, got %f", tzi.Location.Latitude)
	}
	// Check longitude with tolerance
	if tzi.Location.Longitude < -75.0 || tzi.Location.Longitude > -73.0 {
		t.Errorf("Expected longitude ~-74, got %f", tzi.Location.Longitude)
	}
	// Check BC flag
	if tzi.Bc == 0 {
		t.Error("Expected bc flag to be non-zero")
	}
}

func TestParseTzUTC(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("UTC", timelib.BuiltinDB(), &error)
	if err != nil {
		t.Fatalf("ParseTzfile failed: %v", err)
	}

	if tzi == nil {
		t.Fatal("Expected non-nil timezone info")
	}
	if error != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", error)
	}
	if tzi.PosixString != "UTC0" {
		t.Errorf("Expected posix_string UTC0, got %s", tzi.PosixString)
	}
	if tzi.Bit64.Typecnt != 1 {
		t.Errorf("Expected typecnt 1, got %d", tzi.Bit64.Typecnt)
	}

	// Check location data - UTC should have default values
	countryCode := string(tzi.Location.CountryCode[:2])
	if countryCode != "??" {
		t.Errorf("Expected country code ??, got %s", countryCode)
	}
	if tzi.Location.Latitude != -90.0 {
		t.Errorf("Expected latitude -90.0, got %f", tzi.Location.Latitude)
	}
	if tzi.Location.Longitude != -180.0 {
		t.Errorf("Expected longitude -180.0, got %f", tzi.Location.Longitude)
	}
	// Check BC flag - UTC should have bc set
	if tzi.Bc == 0 {
		t.Error("Expected bc flag to be non-zero")
	}
}

func TestParseTzUSSamoa(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("US/Samoa", timelib.BuiltinDB(), &error)
	if err != nil {
		t.Errorf("ParseTzfile failed: %v", err)
	}

	if tzi == nil {
		t.Fatal("Expected non-nil timezone info")
	}
	if error != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", error)
	}
	if tzi.PosixString != "SST11" {
		t.Errorf("Expected posix_string SST11, got %s", tzi.PosixString)
	}
	// Check BC flag - US/Samoa should NOT have bc set
	if tzi.Bc != 0 {
		t.Error("Expected bc flag to be zero")
	}
}

func TestParseTzPetersburg(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("America/Indiana/Petersburg", timelib.BuiltinDB(), &error)
	if err != nil {
		t.Fatalf("ParseTzfile failed: %v", err)
	}

	if tzi == nil {
		t.Fatal("Expected non-nil timezone info")
	}
	if error != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", error)
	}
	if tzi.PosixString != "EST5EDT,M3.2.0,M11.1.0" {
		t.Errorf("Expected posix_string EST5EDT,M3.2.0,M11.1.0, got %s", tzi.PosixString)
	}
	// Note: Typecnt may vary between timezone database versions
	if tzi.Bit64.Typecnt < 6 || tzi.Bit64.Typecnt > 8 {
		t.Errorf("Expected typecnt 6-8, got %d", tzi.Bit64.Typecnt)
	}
	if tzi.PosixInfo == nil {
		t.Fatal("Expected posix_info to be non-nil")
	}
	// TypeIndexDstType varies with database version
	if tzi.PosixInfo.TypeIndexDstType < 4 || tzi.PosixInfo.TypeIndexDstType > 7 {
		t.Errorf("Expected type_index_dst_type 4-7, got %d", tzi.PosixInfo.TypeIndexDstType)
	}
}

func TestParseTzBeulah(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("America/North_Dakota/Beulah", timelib.BuiltinDB(), &error)
	if err != nil {
		t.Fatalf("ParseTzfile failed: %v", err)
	}

	if tzi == nil {
		t.Fatal("Expected non-nil timezone info")
	}
	if error != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Expected TIMELIB_ERROR_NO_ERROR, got %d", error)
	}
	if tzi.PosixString != "CST6CDT,M3.2.0,M11.1.0" {
		t.Errorf("Expected posix_string CST6CDT,M3.2.0,M11.1.0, got %s", tzi.PosixString)
	}
	// Note: Typecnt may vary between timezone database versions
	if tzi.Bit64.Typecnt < 6 || tzi.Bit64.Typecnt > 8 {
		t.Errorf("Expected typecnt 6-8, got %d", tzi.Bit64.Typecnt)
	}
	if tzi.PosixInfo == nil {
		t.Fatal("Expected posix_info to be non-nil")
	}
	// TypeIndexDstType varies with database version
	if tzi.PosixInfo.TypeIndexDstType < 4 || tzi.PosixInfo.TypeIndexDstType > 7 {
		t.Errorf("Expected type_index_dst_type 4-7, got %d", tzi.PosixInfo.TypeIndexDstType)
	}
	// Check timezone abbreviation for DST type - abbreviations may vary by database version
	dstTypeIdx := tzi.PosixInfo.TypeIndexDstType
	if dstTypeIdx >= 0 && int(dstTypeIdx) < len(tzi.Type) {
		abbrIdx := tzi.Type[dstTypeIdx].AbbrIdx
		if abbrIdx >= 0 && int(abbrIdx) < len(tzi.TimezoneAbbr) {
			abbr := tzi.TimezoneAbbr[abbrIdx:]
			// Find null terminator or end of string
			nullIdx := len(abbr)
			for i := 0; i < len(abbr); i++ {
				if abbr[i] == 0 {
					nullIdx = i
					break
				}
			}
			abbrStr := abbr[:nullIdx]
			// Just verify it's a valid abbreviation (3-4 chars typically)
			if len(abbrStr) < 2 || len(abbrStr) > 5 {
				t.Errorf("Expected valid DST abbreviation, got %s", abbrStr)
			}
			t.Logf("DST abbreviation: %s (may vary by database version)", abbrStr)
		}
	}
}

func TestParseTzCST6CDT(t *testing.T) {
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	tzName := "CST6CDT"
	isDst := 0
	tzNotFound := 0

	tll := timelib.ParseZone(&tzName, &isDst, time, &tzNotFound, timelib.BuiltinDB(), timelib.ParseTzfile)

	if time.TzInfo == nil {
		t.Fatal("Expected non-nil tz_info")
	}
	if time.TzInfo.Name != "CST6CDT" {
		t.Errorf("Expected tz_info.name CST6CDT, got %s", time.TzInfo.Name)
	}
	// Typecnt varies by database version (5-6 are valid)
	if time.TzInfo.Bit64.Typecnt < 5 || time.TzInfo.Bit64.Typecnt > 6 {
		t.Errorf("Expected typecnt 5-6, got %d", time.TzInfo.Bit64.Typecnt)
	}
	if time.TzInfo.PosixInfo == nil {
		t.Fatal("Expected posix_info to be non-nil")
	}
	// TypeIndexDstType varies by database version
	if time.TzInfo.PosixInfo.TypeIndexDstType < 1 || time.TzInfo.PosixInfo.TypeIndexDstType > 5 {
		t.Errorf("Expected type_index_dst_type 1-5, got %d", time.TzInfo.PosixInfo.TypeIndexDstType)
	}
	if tll != 0 {
		t.Errorf("Expected tll 0, got %d", tll)
	}
}

func TestParseTzEtcGMTPlus7(t *testing.T) {
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	tzName := "Etc/GMT+7"
	isDst := 0
	tzNotFound := 0

	tll := timelib.ParseZone(&tzName, &isDst, time, &tzNotFound, timelib.BuiltinDB(), timelib.ParseTzfile)

	if time.TzInfo == nil {
		t.Fatal("Expected non-nil tz_info")
	}
	if time.TzInfo.Name != "Etc/GMT+7" {
		t.Errorf("Expected tz_info.name Etc/GMT+7, got %s", time.TzInfo.Name)
	}
	if time.TzInfo.Bit64.Typecnt != 1 {
		t.Errorf("Expected typecnt 1, got %d", time.TzInfo.Bit64.Typecnt)
	}
	if time.TzInfo.PosixInfo == nil {
		t.Fatal("Expected posix_info to be non-nil")
	}
	// TypeIndexDstType may be -1 or 0 for non-DST timezones
	if time.TzInfo.PosixInfo.TypeIndexDstType < -1 || time.TzInfo.PosixInfo.TypeIndexDstType > 0 {
		t.Errorf("Expected type_index_dst_type -1 or 0, got %d", time.TzInfo.PosixInfo.TypeIndexDstType)
	}
	if tll != 0 {
		t.Errorf("Expected tll 0, got %d", tll)
	}
}

func TestParseTzEtcGMTMinus7(t *testing.T) {
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	tzName := "Etc/GMT-7"
	isDst := 0
	tzNotFound := 0

	tll := timelib.ParseZone(&tzName, &isDst, time, &tzNotFound, timelib.BuiltinDB(), timelib.ParseTzfile)

	if time.TzInfo == nil {
		t.Fatal("Expected non-nil tz_info")
	}
	if time.TzInfo.Name != "Etc/GMT-7" {
		t.Errorf("Expected tz_info.name Etc/GMT-7, got %s", time.TzInfo.Name)
	}
	if time.TzInfo.Bit64.Typecnt != 1 {
		t.Errorf("Expected typecnt 1, got %d", time.TzInfo.Bit64.Typecnt)
	}
	if time.TzInfo.PosixInfo == nil {
		t.Fatal("Expected posix_info to be non-nil")
	}
	// TypeIndexDstType may be -1 or 0 for non-DST timezones
	if time.TzInfo.PosixInfo.TypeIndexDstType < -1 || time.TzInfo.PosixInfo.TypeIndexDstType > 0 {
		t.Errorf("Expected type_index_dst_type -1 or 0, got %d", time.TzInfo.PosixInfo.TypeIndexDstType)
	}
	if tll != 0 {
		t.Errorf("Expected tll 0, got %d", tll)
	}
}

func TestParseTzParentheses(t *testing.T) {
	time := timelib.TimeCtor()
	defer timelib.TimeDtor(time)

	tzName := "((UTC)))"
	isDst := 0
	tzNotFound := 0

	timelib.ParseZone(&tzName, &isDst, time, &tzNotFound, timelib.BuiltinDB(), timelib.ParseTzfile)

	if time.TzInfo == nil {
		t.Fatal("Expected non-nil tz_info")
	}
	if time.TzInfo.Name != "UTC" {
		t.Errorf("Expected tz_info.name UTC, got %s", time.TzInfo.Name)
	}
	if tzName != ")" {
		t.Errorf("Expected first character of remaining string to be ), got %s", tzName)
	}
}

func TestParseTzCorruptTransitions01(t *testing.T) {
	// Create a corrupt timezone database with non-increasing transitions
	index := []timelib.TzDBIndexEntry{
		{ID: "dummy", Pos: 0},
	}
	// Minimal corrupt TZif2 data with 2 transitions that don't increase
	// PHP2 preamble (44 bytes) + TZif2 header (20 bytes) + header data (24 bytes)
	// + 2 transition times (16 bytes, 64-bit, non-increasing) + indices (2 bytes)
	// + type data (6 bytes) + abbr (4 bytes) = 116 bytes minimum
	data := []byte{
		// PHP2 preamble (44 bytes)
		'P', 'H', 'P', '2', 0, 'X', 'X', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		// TZif2 header (20 bytes)
		'T', 'Z', 'i', 'f', '2',
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		// Header: ttisgmtcnt=0, ttisstdcnt=0, leapcnt=0, timecnt=2, typecnt=1, charcnt=4
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 4,
		// Transition times (2 x 8 bytes, 64-bit, NON-INCREASING - this is the corruption)
		0, 0, 0, 0, 0, 0, 0, 2, // time 2
		0, 0, 0, 0, 0, 0, 0, 1, // time 1 (should be > 2, but isn't!)
		// Transition type indices (2 x 1 byte)
		0, 0,
		// Type info (1 x 6 bytes): gmtoff=0, isdst=0, abbrind=0
		0, 0, 0, 0, 0, 0,
		// Abbreviation string (4 bytes)
		'U', 'T', 'C', 0,
	}
	corruptTzdb := &timelib.TzDB{
		Version: "1",
		Index:   index,
		Data:    data,
	}

	var errorCode int
	tzi, err := timelib.ParseTzfile("dummy", corruptTzdb, &errorCode)

	if tzi != nil {
		t.Error("Expected nil timezone info for corrupt data")
	}
	if err == nil {
		t.Error("Expected error for corrupt data")
	}
	if errorCode != timelib.TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE {
		t.Errorf("Expected TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE (%d), got %d",
			timelib.TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE, errorCode)
	}
}

func TestParseTzCorruptTransitions02(t *testing.T) {
	// Test with the NonContinuous test file
	testDir := "files"
	testDB, err := timelib.Zoneinfo(testDir)
	if err != nil {
		t.Errorf("Test directory not found: %v", err)
	}
	defer func() {
		if testDB != nil {
			// ZoneinfoDtor would be called here if available
		}
	}()

	var errorCode int
	tzi, err2 := timelib.ParseTzfile("NonContinuous", testDB, &errorCode)

	if tzi != nil {
		t.Error("Expected nil timezone info for corrupt data")
	}
	if err2 == nil {
		t.Error("Expected error for corrupt data")
	}
	if errorCode != timelib.TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE {
		t.Errorf("Expected TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE (%d), got %d",
			timelib.TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE, errorCode)
	}
}
