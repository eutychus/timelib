package timelib

import (
	"fmt"
	"strings"
)

// TimezoneIDIsValid checks if timezone ID is valid
func TimezoneIDIsValid(timezone string, tzdb *TzDB) bool {
	if tzdb == nil {
		return false
	}

	// Check if timezone exists in the database
	for _, entry := range tzdb.Index {
		if entry.ID == timezone {
			return true
		}
	}
	return false
}

// ParseTzfile parses timezone file from the database
func ParseTzfile(timezone string, tzdb *TzDB, errorCode *int) (*TzInfo, error) {
	return UpdateParseTzfile(timezone, tzdb, errorCode)
}

// TzinfoDtor frees timezone info resources
func TzinfoDtor(tz *TzInfo) {
	if tz == nil {
		return
	}

	// Free allocated resources
	tz.Name = ""
	tz.TimezoneAbbr = ""
	tz.PosixString = ""

	if tz.PosixInfo != nil {
		if tz.PosixInfo.Std != "" {
			tz.PosixInfo.Std = ""
		}
		if tz.PosixInfo.Dst != "" {
			tz.PosixInfo.Dst = ""
		}
		if tz.PosixInfo.DstBegin != nil {
			tz.PosixInfo.DstBegin = nil
		}
		if tz.PosixInfo.DstEnd != nil {
			tz.PosixInfo.DstEnd = nil
		}
		tz.PosixInfo = nil
	}

	// Clear slices
	tz.Trans = tz.Trans[:0]
	tz.TransIdx = tz.TransIdx[:0]
	tz.Type = tz.Type[:0]
	tz.LeapTimes = tz.LeapTimes[:0]
}

// TzinfoClone deep-clones a TzInfo structure
func TzinfoClone(tz *TzInfo) *TzInfo {
	if tz == nil {
		return nil
	}

	clone := &TzInfo{
		Name:         tz.Name,
		TimezoneAbbr: tz.TimezoneAbbr,
		PosixString:  tz.PosixString,
		Bc:           tz.Bc,
	}

	// Clone bit32 and bit64 structures
	clone.Bit32 = tz.Bit32
	clone.Bit64 = tz.Bit64

	// Clone slices
	if len(tz.Trans) > 0 {
		clone.Trans = make([]int64, len(tz.Trans))
		copy(clone.Trans, tz.Trans)
	}

	if len(tz.TransIdx) > 0 {
		clone.TransIdx = make([]uint8, len(tz.TransIdx))
		copy(clone.TransIdx, tz.TransIdx)
	}

	if len(tz.Type) > 0 {
		clone.Type = make([]TTInfo, len(tz.Type))
		copy(clone.Type, tz.Type)
	}

	if len(tz.LeapTimes) > 0 {
		clone.LeapTimes = make([]TLInfo, len(tz.LeapTimes))
		copy(clone.LeapTimes, tz.LeapTimes)
	}

	// Clone location info
	clone.Location = tz.Location

	// Clone POSIX info if present
	if tz.PosixInfo != nil {
		clone.PosixInfo = &PosixStr{
			Std:              tz.PosixInfo.Std,
			StdOffset:        tz.PosixInfo.StdOffset,
			Dst:              tz.PosixInfo.Dst,
			DstOffset:        tz.PosixInfo.DstOffset,
			TypeIndexStdType: tz.PosixInfo.TypeIndexStdType,
			TypeIndexDstType: tz.PosixInfo.TypeIndexDstType,
		}

		if tz.PosixInfo.DstBegin != nil {
			clone.PosixInfo.DstBegin = &PosixTransInfo{
				Type: tz.PosixInfo.DstBegin.Type,
				Days: tz.PosixInfo.DstBegin.Days,
				Hour: tz.PosixInfo.DstBegin.Hour,
			}
			clone.PosixInfo.DstBegin.Mwd = tz.PosixInfo.DstBegin.Mwd
		}

		if tz.PosixInfo.DstEnd != nil {
			clone.PosixInfo.DstEnd = &PosixTransInfo{
				Type: tz.PosixInfo.DstEnd.Type,
				Days: tz.PosixInfo.DstEnd.Days,
				Hour: tz.PosixInfo.DstEnd.Hour,
			}
			clone.PosixInfo.DstEnd.Mwd = tz.PosixInfo.DstEnd.Mwd
		}
	}

	return clone
}

// TimestampIsInDst checks if timestamp is in DST for the given timezone
func TimestampIsInDst(ts int64, tz *TzInfo) int {
	isDst, err := IsTimestampInDST(ts, tz)
	if err != nil {
		return -1
	}
	return isDst
}

// GetTimeZoneInfo returns timezone offset information for a timestamp
func GetTimeZoneInfo(ts int64, tz *TzInfo) *TimeOffset {
	if tz == nil {
		return nil
	}

	var offset int32 = 0
	var transitionTime int64 = 0
	var abbr string = "GMT"
	var isDst uint = 0

	// Fetch timezone offset using the helper function
	to, tt := fetchTimezoneOffset(tz, ts)
	if to != nil {
		offset = to.Offset
		transitionTime = tt
		isDst = uint(to.IsDst)
		if int(to.AbbrIdx) < len(tz.TimezoneAbbr) {
			// Extract null-terminated string from timezone_abbr
			abbr = extractNullTerminatedString(tz.TimezoneAbbr, int(to.AbbrIdx))
		}
	} else {
		offset = 0
		// Use TimezoneAbbr if available, otherwise fall back to Name
		if tz.TimezoneAbbr != "" {
			abbr = tz.TimezoneAbbr
		} else {
			abbr = tz.Name
		}
		isDst = 0
		transitionTime = 0
	}

	return &TimeOffset{
		Offset:         offset,
		LeapSecs:       0, // TODO: implement leap seconds
		IsDst:          int(isDst),
		Abbr:           abbr,
		TransitionTime: transitionTime,
	}
}

// Helper function to extract a null-terminated string from a byte array
func extractNullTerminatedString(data string, offset int) string {
	if offset >= len(data) {
		return ""
	}
	for i := offset; i < len(data); i++ {
		if data[i] == 0 {
			return data[offset:i]
		}
	}
	return data[offset:]
}

// GetTimeZoneOffsetInfo returns detailed timezone offset information
func GetTimeZoneOffsetInfo(ts int64, tz *TzInfo, offset *int32, transitionTime *int64, isDst *uint) int {
	if tz == nil {
		return 0 // Failure
	}

	// Use the internal helper function
	success := getTimeZoneOffsetInfo(ts, tz, offset, transitionTime, isDst)
	if success {
		return 1
	}
	return 0
}

// GetCurrentOffset returns the current UTC offset for the given time
func GetCurrentOffset(t *Time) int64 {
	offset, err := GetCurrentOffsetForTime(t)
	if err != nil {
		return 0
	}
	return int64(offset)
}

// SameTimezone checks if two times have the same timezone
func SameTimezone(one, two *Time) bool {
	if one == nil || two == nil {
		return false
	}

	// Check if zone types are the same
	if one.ZoneType != two.ZoneType {
		return false
	}

	switch one.ZoneType {
	case TIMELIB_ZONETYPE_OFFSET, TIMELIB_ZONETYPE_ABBR:
		// For offset and abbreviation types, check if the effective offset is the same
		offset1 := one.Z + int32(one.Dst*3600)
		offset2 := two.Z + int32(two.Dst*3600)
		return offset1 == offset2

	case TIMELIB_ZONETYPE_ID:
		// For ID types, check if timezone names are the same
		if one.TzInfo == nil || two.TzInfo == nil {
			return one.TzInfo == two.TzInfo // Both nil means same
		}
		return one.TzInfo.Name == two.TzInfo.Name

	case TIMELIB_ZONETYPE_NONE:
		return true // Both have no timezone info
	}

	return false
}

// BuiltinDB is now implemented in tzdata_builtin.go

// TimezoneIdentifiersList returns a list of timezone identifiers
func TimezoneIdentifiersList(tzdb *TzDB, count *int) []TzDBIndexEntry {
	if tzdb == nil {
		if count != nil {
			*count = 0
		}
		return nil
	}

	if count != nil {
		*count = tzdb.IndexSize
	}

	// Return a copy of the index
	result := make([]TzDBIndexEntry, len(tzdb.Index))
	copy(result, tzdb.Index)
	return result
}

// Zoneinfo scans directory for timezone files and builds a database
func Zoneinfo(directory string) (*TzDB, error) {
	return ZoneinfoDir(directory)
}

// TimezoneIDFromAbbr gets timezone ID from abbreviation
func TimezoneIDFromAbbr(abbr string, gmtoffset int64, isDst int) string {
	if abbr == "" {
		return ""
	}

	// Handle UTC/GMT case-insensitively
	lowerAbbr := strings.ToLower(abbr)
	if lowerAbbr == "utc" || lowerAbbr == "gmt" {
		return "UTC"
	}

	// This would use the timezone mapping data from timezonemap.h and fallbackmap.h
	// For now, return empty string to indicate no match found
	return ""
}

// TimezoneAbbreviationsList returns known timezone abbreviations
func TimezoneAbbreviationsList() []TzLookupTable {
	// This would return the actual timezone abbreviations list
	// For now, return a basic list with common abbreviations
	return []TzLookupTable{
		{Name: "UTC", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: 0, FullTzName: "UTC"},
		{Name: "GMT", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: 0, FullTzName: "GMT"},
		{Name: "EST", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: -5 * 3600, FullTzName: "America/New_York"},
		{Name: "EDT", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: -4 * 3600, FullTzName: "America/New_York"},
		{Name: "CST", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: -6 * 3600, FullTzName: "America/Chicago"},
		{Name: "CDT", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: -5 * 3600, FullTzName: "America/Chicago"},
		{Name: "PST", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: -8 * 3600, FullTzName: "America/Los_Angeles"},
		{Name: "PDT", Type: TIMELIB_ZONETYPE_OFFSET, GmtOffset: -7 * 3600, FullTzName: "America/Los_Angeles"},
	}
}

// DumpTzinfo displays debugging information about timezone info
func DumpTzinfo(tz *TzInfo) {
	if tz == nil {
		fmt.Println("Timezone info is nil")
		return
	}

	fmt.Printf("Timezone: %s\n", tz.Name)
	fmt.Printf("Bit32: ttisgmtcnt=%d, ttisstdcnt=%d, leapcnt=%d, timecnt=%d, typecnt=%d, charcnt=%d\n",
		tz.Bit32.Ttisgmtcnt, tz.Bit32.Ttisstdcnt, tz.Bit32.Leapcnt, tz.Bit32.Timecnt, tz.Bit32.Typecnt, tz.Bit32.Charcnt)
	fmt.Printf("Bit64: ttisgmtcnt=%d, ttisstdcnt=%d, leapcnt=%d, timecnt=%d, typecnt=%d, charcnt=%d\n",
		tz.Bit64.Ttisgmtcnt, tz.Bit64.Ttisstdcnt, tz.Bit64.Leapcnt, tz.Bit64.Timecnt, tz.Bit64.Typecnt, tz.Bit64.Charcnt)
	fmt.Printf("Transitions: %d\n", len(tz.Trans))
	fmt.Printf("Types: %d\n", len(tz.Type))
	fmt.Printf("Leap times: %d\n", len(tz.LeapTimes))
	fmt.Printf("Abbreviation: %s\n", tz.TimezoneAbbr)
	fmt.Printf("Location: %s (%.6f, %.6f)\n", tz.Location.CountryCode, tz.Location.Latitude, tz.Location.Longitude)

	if tz.PosixString != "" {
		fmt.Printf("POSIX string: %s\n", tz.PosixString)
	}
}

// ParsePosixStr parses a POSIX timezone string
// This is a wrapper that calls ParsePosixString with nil TzInfo
func ParsePosixStr(posix string) (*PosixStr, error) {
	return ParsePosixString(posix, nil)
}

// PosixStrDtor frees POSIX string resources
func PosixStrDtor(ps *PosixStr) {
	if ps == nil {
		return
	}

	ps.Std = ""
	ps.Dst = ""

	if ps.DstBegin != nil {
		ps.DstBegin = nil
	}
	if ps.DstEnd != nil {
		ps.DstEnd = nil
	}
}

// Timezone Transition Calculation Functions
//
// These functions calculate DST transitions for timezones using POSIX rules.
// POSIX rules are used for timestamps beyond the last transition in the tzfile,
// allowing timezone calculations to extend indefinitely into the future.

// countLeapYears counts leap years from year 1 to given year.
// Because this is for Jan 1 (before Feb 29), we subtract 1 from the year
// before calculating to avoid counting the current year's potential leap day.
//
// Matches C function: count_leap_years
func countLeapYears(y int64) int64 {
	// Because we want this for Jan 1, the leap day hasn't happened yet, so
	// subtract one of year before we calculate
	y--
	return (y / 4) - (y / 100) + (y / 400)
}

// TsAtStartOfYear returns Unix timestamp at start of given year (Jan 1, 00:00:00)
// Matches C function: timelib_ts_at_start_of_year
func TsAtStartOfYear(year int64) int64 {
	epochLeapYears := countLeapYears(1970)
	currentLeapYears := countLeapYears(year)

	return SECS_PER_DAY * (
		((year - 1970) * DAYS_PER_YEAR) +
		currentLeapYears -
		epochLeapYears)
}

// calcTransition calculates the seconds from start of year for a POSIX transition.
// This implements all three POSIX transition rule formats:
//  1. Type 1 (Jn): Julian day without leap days (1-365)
//  2. Type 2 (n): Julian day with leap days (0-365)
//  3. Type 3 (Mm.w.d): Month/week/day format using Zeller's Congruence
//
// Matches C function: calc_transition
func calcTransition(psi *PosixTransInfo, year int64) int64 {
	if psi == nil {
		return 0
	}

	leapYear := IsLeapYear(year)

	switch psi.Type {
	case 1: // TIMELIB_POSIX_TRANS_TYPE_JULIAN_NO_FEB29 - Jn format
		value := int64(psi.Days - 1)
		if leapYear && psi.Days >= 60 {
			value++
		}
		return value * SECS_PER_DAY

	case 2: // TIMELIB_POSIX_TRANS_TYPE_JULIAN_FEB29 - n format
		return int64(psi.Days) * SECS_PER_DAY

	case 3: // TIMELIB_POSIX_TRANS_TYPE_MWD - Mm.w.d format
		// Use Zeller's Congruence to get day-of-week of first day of month
		m1 := int64((psi.Mwd.Month + 9) % 12 + 1)
		yy0 := year
		if psi.Mwd.Month <= 2 {
			yy0 = year - 1
		}
		yy1 := yy0 / 100
		yy2 := yy0 % 100
		dow := ((26*m1 - 2) / 10 + 1 + yy2 + yy2/4 + yy1/4 - 2*yy1) % 7
		if dow < 0 {
			dow += DAYS_PER_WEEK
		}

		// dow is the day-of-week of the first day of the month
		// Get the day-of-month (zero-origin) of the first "dow" day of the month
		d := int64(psi.Mwd.Dow) - dow
		if d < 0 {
			d += DAYS_PER_WEEK
		}

		// Add weeks to get to the nth occurrence
		for i := 1; i < psi.Mwd.Week; i++ {
			monthLen := int64(monthLengths[0][psi.Mwd.Month-1])
			if leapYear {
				monthLen = int64(monthLengths[1][psi.Mwd.Month-1])
			}
			if d+DAYS_PER_WEEK >= monthLen {
				break
			}
			d += DAYS_PER_WEEK
		}

		// d is the day-of-month (zero-origin) of the day we want
		value := int64(d) * SECS_PER_DAY
		for i := 0; i < psi.Mwd.Month-1; i++ {
			monthLen := monthLengths[0][i]
			if leapYear {
				monthLen = monthLengths[1][i]
			}
			value += int64(monthLen) * SECS_PER_DAY
		}

		return value
	}

	return 0
}

// monthLengths contains days in each month for normal and leap years
var monthLengths = [2][12]int{
	{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}, // normal year
	{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}, // leap year
}

// GetTransitionsForYear calculates DST transitions for a specific year.
// This function computes the exact Unix timestamps when DST begins and ends
// for the given year, based on the POSIX transition rules. The transitions
// are ordered correctly for both Northern and Southern hemispheres.
//
// The function is typically called with year-1, year, and year+1 to handle
// timestamps near year boundaries.
//
// Matches C function: timelib_get_transitions_for_year
func GetTransitionsForYear(tz *TzInfo, year int64, transitions *PosixTransitions) {
	if tz == nil || transitions == nil || tz.PosixInfo == nil {
		return
	}

	if tz.PosixInfo.DstBegin == nil || tz.PosixInfo.DstEnd == nil {
		return
	}

	yearBeginTs := TsAtStartOfYear(year)

	// Calculate begin transition
	transBegin := yearBeginTs
	transBegin += calcTransition(tz.PosixInfo.DstBegin, year)
	transBegin += int64(tz.PosixInfo.DstBegin.Hour)
	transBegin -= tz.PosixInfo.StdOffset

	// Calculate end transition
	transEnd := yearBeginTs
	transEnd += calcTransition(tz.PosixInfo.DstEnd, year)
	transEnd += int64(tz.PosixInfo.DstEnd.Hour)
	transEnd -= tz.PosixInfo.DstOffset

	// Store transitions in order
	if transBegin < transEnd {
		transitions.Times[transitions.Count] = transBegin
		transitions.Times[transitions.Count+1] = transEnd
		transitions.Types[transitions.Count] = int64(tz.PosixInfo.TypeIndexDstType)
		transitions.Types[transitions.Count+1] = int64(tz.PosixInfo.TypeIndexStdType)
	} else {
		transitions.Times[transitions.Count+1] = transBegin
		transitions.Times[transitions.Count] = transEnd
		transitions.Types[transitions.Count+1] = int64(tz.PosixInfo.TypeIndexDstType)
		transitions.Types[transitions.Count] = int64(tz.PosixInfo.TypeIndexStdType)
	}

	transitions.Count += 2
}
