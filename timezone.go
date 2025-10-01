package timelib

import (
	"errors"
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

	// This is a simplified implementation
	// Full timezone info would require parsing transitions
	offset := &TimeOffset{
		Offset:         0, // Would be calculated from timezone data
		LeapSecs:       0,
		IsDst:          0,
		Abbr:           tz.TimezoneAbbr,
		TransitionTime: ts,
	}

	return offset
}

// GetTimeZoneOffsetInfo returns detailed timezone offset information
func GetTimeZoneOffsetInfo(ts int64, tz *TzInfo, offset *int32, transitionTime *int64, isDst *uint) int {
	if tz == nil {
		return 0 // Failure
	}

	// This is a simplified implementation
	// Full implementation would parse timezone transitions

	if offset != nil {
		*offset = 0 // Would be calculated from timezone data
	}
	if transitionTime != nil {
		*transitionTime = ts
	}
	if isDst != nil {
		*isDst = 0
	}

	return 1 // Success
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
func ParsePosixStr(posix string) (*PosixStr, error) {
	if posix == "" {
		return nil, errors.New("empty POSIX string")
	}

	// This is a simplified POSIX string parser
	// Full implementation would be more complex

	ps := &PosixStr{}

	// Basic parsing - this would need to be expanded for full POSIX support
	parts := strings.Fields(posix)
	if len(parts) < 1 {
		return nil, errors.New("invalid POSIX string format")
	}

	// Parse standard time abbreviation and offset
	ps.Std = parts[0]

	// For now, return basic structure
	// Full POSIX parsing would be implemented here

	return ps, nil
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

// GetTransitionsForYear calculates DST transitions for a specific year
func GetTransitionsForYear(tz *TzInfo, year int64, transitions *PosixTransitions) {
	if tz == nil || transitions == nil {
		return
	}

	// This is a simplified implementation
	// Full implementation would calculate actual DST transitions based on timezone rules

	transitions.Count = 0

	// For now, return empty transitions
	// Full implementation would parse timezone rules and calculate transitions
}
