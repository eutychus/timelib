// Package timelib implements POSIX timezone string parsing and transition calculations.
//
// POSIX Timezone String Format:
//
//	std offset [dst [offset] [,start[/time],end[/time]]]
//
// Where:
//   - std: Standard time abbreviation (e.g., "EST", "<+03>")
//   - offset: Hours west of UTC (POSIX convention: positive = west)
//   - dst: Daylight saving time abbreviation (optional)
//   - start/end: DST transition rules in one of three formats:
//     1. Jn: Julian day n (1-365), February 29 never counted
//     2. n: Julian day n (0-365), February 29 counted in leap years
//     3. Mm.w.d: Month m, week w (1-5, where 5=last), day of week d (0=Sunday)
//
// Examples:
//   - EST5EDT,M3.2.0,M11.1.0: US Eastern time
//   - CST6CDT,M3.2.0/2,M11.1.0/2: US Central time with 2am transitions
//   - <+03>-3: UTC+3 with no DST
package timelib

import (
	"errors"
	"fmt"
	"strconv"
)

// findTTInfoIndex finds a type index matching offset, isDst, and abbreviation.
// This is critical for correctly mapping POSIX abbreviations (like "EDT") to the
// corresponding TTInfo entries in the timezone file, ensuring proper timezone
// display when using POSIX rules for dates beyond the transition table.
//
// Matches C function: find_ttinfo_index
func findTTInfoIndex(tz *TzInfo, offset int32, isDst int, abbr string) int {
	if tz == nil {
		return -1
	}

	for i, t := range tz.Type {
		if t.Offset == offset && t.IsDst == isDst {
			// Extract abbreviation from timezone_abbr
			if int(t.AbbrIdx) < len(tz.TimezoneAbbr) {
				typeAbbr := ""
				for j := int(t.AbbrIdx); j < len(tz.TimezoneAbbr); j++ {
					if tz.TimezoneAbbr[j] == 0 {
						typeAbbr = tz.TimezoneAbbr[int(t.AbbrIdx):j]
						break
					}
				}
				if typeAbbr == abbr {
					return i
				}
			}
		}
	}

	return -1
}

// addAbbr adds an abbreviation to the timezone abbreviation string
// Returns the index where the abbreviation was added
// Matches C function: add_abbr
func addAbbr(tz *TzInfo, abbr string) int {
	if tz == nil {
		return 0
	}

	oldLength := len(tz.TimezoneAbbr)
	// Append abbreviation + null terminator
	tz.TimezoneAbbr += abbr + "\x00"
	tz.Bit64.Charcnt = uint64(len(tz.TimezoneAbbr))

	return oldLength
}

// addNewTTInfoIndex adds a new TTInfo type to the timezone
// Returns the index of the newly added type
// Matches C function: add_new_ttinfo_index
func addNewTTInfoIndex(tz *TzInfo, offset int32, isDst int, abbr string) int {
	if tz == nil {
		return -1
	}

	// Add abbreviation and get its index
	abbrIdx := addAbbr(tz, abbr)

	// Create new type
	newType := TTInfo{
		Offset:  offset,
		IsDst:   isDst,
		AbbrIdx: abbrIdx,
		IsStd:   0,
		IsUtc:   0,
	}

	// Append to type array
	tz.Type = append(tz.Type, newType)
	tz.Bit64.Typecnt++

	return int(tz.Bit64.Typecnt - 1)
}

// ParsePosixString parses a POSIX timezone string
// Format: std offset [dst [offset] [,start[/time],end[/time]]]
// Examples:
//   EST5EDT,M3.2.0,M11.1.0
//   CST6CDT,M3.2.0,M11.1.0
//   <+03>-3
func ParsePosixString(posix string, tz *TzInfo) (*PosixStr, error) {
	if posix == "" {
		return nil, errors.New("empty POSIX string")
	}

	ps := &PosixStr{}
	pos := 0

	// Parse standard timezone name
	name, offset, _, newPos, err := parsePosixName(posix, pos)
	if err != nil {
		return nil, fmt.Errorf("failed to parse standard name: %v", err)
	}
	ps.Std = name
	ps.StdOffset = offset
	pos = newPos

	// Check if there's DST information
	if pos >= len(posix) {
		// No DST, find matching standard type
		if tz != nil && len(tz.Type) > 0 {
			ps.TypeIndexStdType = findTTInfoIndex(tz, int32(ps.StdOffset), 0, ps.Std)
			if ps.TypeIndexStdType == -1 {
				// If not found, use first non-DST type as fallback
				for i, t := range tz.Type {
					if t.IsDst == 0 {
						ps.TypeIndexStdType = i
						break
					}
				}
				if ps.TypeIndexStdType == -1 {
					ps.TypeIndexStdType = 0
				}
			}
			ps.TypeIndexDstType = -1
		} else if tz != nil {
			ps.TypeIndexStdType = 0
			ps.TypeIndexDstType = -1
		}
		return ps, nil
	}

	// Parse DST timezone name
	if posix[pos] == ',' {
		// No DST name, use standard + 'D'
		ps.Dst = ps.Std + "D"
		ps.DstOffset = ps.StdOffset + 3600 // Default: 1 hour ahead
	} else {
		var hasOffset bool
		name, offset, hasOffset, newPos, err = parsePosixName(posix, pos)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DST name: %v", err)
		}
		ps.Dst = name
		if hasOffset {
			// Offset was explicitly specified (even if it's 0)
			ps.DstOffset = offset
		} else {
			// No offset specified, use default: 1 hour ahead of standard
			ps.DstOffset = ps.StdOffset + 3600
		}
		pos = newPos
	}

	// Parse DST transition rules
	if pos < len(posix) && posix[pos] == ',' {
		pos++ // Skip comma

		// Parse DST start rule
		start, newPos, err := parsePosixTransitionRule(posix, pos)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DST start rule: %v", err)
		}
		ps.DstBegin = start
		pos = newPos

		if pos < len(posix) && posix[pos] == ',' {
			pos++ // Skip comma

			// Parse DST end rule
			end, newPos, err := parsePosixTransitionRule(posix, pos)
			if err != nil {
				return nil, fmt.Errorf("failed to parse DST end rule: %v", err)
			}
			ps.DstEnd = end
			pos = newPos
		}
	}

	// Set type indices by finding matching types based on offset, isDst, and abbreviation
	// This matches the C function: find_ttinfo_index and add_new_ttinfo_index
	if tz != nil {
		// Find standard type index
		ps.TypeIndexStdType = findTTInfoIndex(tz, int32(ps.StdOffset), 0, ps.Std)
		if ps.TypeIndexStdType == -1 {
			// If not found, add new type (matches C add_new_ttinfo_index)
			ps.TypeIndexStdType = addNewTTInfoIndex(tz, int32(ps.StdOffset), 0, ps.Std)
		}

		// Find DST type index if DST is present
		if ps.Dst != "" {
			ps.TypeIndexDstType = findTTInfoIndex(tz, int32(ps.DstOffset), 1, ps.Dst)
			if ps.TypeIndexDstType == -1 {
				// If not found, add new type (matches C add_new_ttinfo_index)
				ps.TypeIndexDstType = addNewTTInfoIndex(tz, int32(ps.DstOffset), 1, ps.Dst)
			}
		} else {
			ps.TypeIndexDstType = -1
		}
	}

	return ps, nil
}

// parsePosixName parses a timezone name and optional offset
// Returns name, offset (in seconds, POSIX convention), new position
func parsePosixName(s string, pos int) (string, int64, bool, int, error) {
	if pos >= len(s) {
		return "", 0, false, pos, errors.New("unexpected end of string")
	}

	var name string
	var offset int64
	hasOffset := false

	// Check if name is quoted with < >
	if s[pos] == '<' {
		pos++
		end := pos
		for end < len(s) && s[end] != '>' {
			end++
		}
		if end >= len(s) {
			return "", 0, false, pos, errors.New("unclosed < in timezone name")
		}
		name = s[pos:end]
		pos = end + 1
	} else {
		// Read alphabetic characters
		end := pos
		for end < len(s) && ((s[end] >= 'A' && s[end] <= 'Z') || (s[end] >= 'a' && s[end] <= 'z')) {
			end++
		}
		if end == pos {
			return "", 0, false, pos, errors.New("expected timezone name")
		}
		name = s[pos:end]
		pos = end
	}

	// Parse offset if present
	if pos < len(s) && (s[pos] == '+' || s[pos] == '-' || (s[pos] >= '0' && s[pos] <= '9')) {
		var err error
		offset, pos, err = parsePosixOffset(s, pos)
		if err != nil {
			return name, 0, false, pos, err
		}
		hasOffset = true
	}

	return name, offset, hasOffset, pos, nil
}

// parsePosixOffset parses a timezone offset
// POSIX offsets are OPPOSITE of normal: positive means west of Greenwich
// Returns offset in seconds (converted to normal convention), new position
func parsePosixOffset(s string, pos int) (int64, int, error) {
	if pos >= len(s) {
		return 0, pos, nil
	}

	sign := int64(-1) // Default: positive number means west (negative offset)
	if s[pos] == '+' {
		sign = -1 // POSIX: + means west (negative offset)
		pos++
	} else if s[pos] == '-' {
		sign = 1 // POSIX: - means east (positive offset)
		pos++
	}

	// Parse hours
	hours, newPos, err := parseNumber(s, pos)
	if err != nil {
		return 0, pos, err
	}
	pos = newPos

	minutes := int64(0)
	seconds := int64(0)

	// Parse optional minutes
	if pos < len(s) && s[pos] == ':' {
		pos++
		minutes, pos, err = parseNumber(s, pos)
		if err != nil {
			return 0, pos, err
		}

		// Parse optional seconds
		if pos < len(s) && s[pos] == ':' {
			pos++
			seconds, pos, err = parseNumber(s, pos)
			if err != nil {
				return 0, pos, err
			}
		}
	}

	offset := sign * (hours*3600 + minutes*60 + seconds)
	return offset, pos, nil
}

// parsePosixTransitionRule parses a DST transition rule
// Formats:
//   Jn    - Julian day n (1-365, leap days not counted)
//   n     - Julian day n (0-365, leap days counted)
//   Mm.w.d - Month m, week w, day d
func parsePosixTransitionRule(s string, pos int) (*PosixTransInfo, int, error) {
	if pos >= len(s) {
		return nil, pos, errors.New("unexpected end of string")
	}

	info := &PosixTransInfo{
		Hour: 7200, // Default transition time is 2:00 AM (2 hours * 3600 seconds)
	}

	if s[pos] == 'J' {
		// Jn format (Julian day, no leap days)
		info.Type = 1
		pos++
		days, newPos, err := parseNumber(s, pos)
		if err != nil {
			return nil, pos, err
		}
		info.Days = int(days)
		pos = newPos
	} else if s[pos] == 'M' {
		// Mm.w.d format (month, week, day)
		info.Type = 3
		pos++

		// Parse month
		month, newPos, err := parseNumber(s, pos)
		if err != nil {
			return nil, pos, err
		}
		info.Mwd.Month = int(month)
		pos = newPos

		if pos >= len(s) || s[pos] != '.' {
			return nil, pos, errors.New("expected '.' after month")
		}
		pos++

		// Parse week
		week, newPos, err := parseNumber(s, pos)
		if err != nil {
			return nil, pos, err
		}
		info.Mwd.Week = int(week)
		pos = newPos

		if pos >= len(s) || s[pos] != '.' {
			return nil, pos, errors.New("expected '.' after week")
		}
		pos++

		// Parse day of week
		dow, newPos, err := parseNumber(s, pos)
		if err != nil {
			return nil, pos, err
		}
		info.Mwd.Dow = int(dow)
		pos = newPos
	} else if s[pos] >= '0' && s[pos] <= '9' {
		// n format (Julian day, with leap days)
		info.Type = 2
		days, newPos, err := parseNumber(s, pos)
		if err != nil {
			return nil, pos, err
		}
		info.Days = int(days)
		pos = newPos
	} else {
		return nil, pos, fmt.Errorf("unexpected character in transition rule: %c", s[pos])
	}

	// Parse optional time
	if pos < len(s) && s[pos] == '/' {
		pos++
		time, newPos, err := parsePosixTime(s, pos)
		if err != nil {
			return nil, pos, err
		}
		info.Hour = time
		pos = newPos
	}

	return info, pos, nil
}

// parsePosixTime parses a time in [+-]hh[:mm[:ss]] format
// Returns time in seconds since midnight
func parsePosixTime(s string, pos int) (int, int, error) {
	sign := 1
	if pos < len(s) && (s[pos] == '+' || s[pos] == '-') {
		if s[pos] == '-' {
			sign = -1
		}
		pos++
	}

	hours, newPos, err := parseNumber(s, pos)
	if err != nil {
		return 0, pos, err
	}
	pos = newPos

	minutes := int64(0)
	seconds := int64(0)

	if pos < len(s) && s[pos] == ':' {
		pos++
		minutes, pos, err = parseNumber(s, pos)
		if err != nil {
			return 0, pos, err
		}

		if pos < len(s) && s[pos] == ':' {
			pos++
			seconds, pos, err = parseNumber(s, pos)
			if err != nil {
				return 0, pos, err
			}
		}
	}

	time := sign * int(hours*3600+minutes*60+seconds)
	return time, pos, nil
}

// parseNumber parses a number from string
func parseNumber(s string, pos int) (int64, int, error) {
	if pos >= len(s) {
		return 0, pos, errors.New("unexpected end of string")
	}

	start := pos
	for pos < len(s) && s[pos] >= '0' && s[pos] <= '9' {
		pos++
	}

	if pos == start {
		return 0, pos, errors.New("expected number")
	}

	num, err := strconv.ParseInt(s[start:pos], 10, 64)
	if err != nil {
		return 0, pos, err
	}

	return num, pos, nil
}

// CalculatePosixTransitionTime calculates the Unix timestamp for a POSIX transition
// in a specific year. This function converts a POSIX transition rule (Jn, n, or Mm.w.d)
// into an absolute Unix timestamp, adjusting for the timezone offset.
//
// Parameters:
//   - year: The year to calculate the transition for
//   - trans: The POSIX transition rule (start or end of DST)
//   - offset: The timezone offset to subtract (either std_offset or dst_offset)
//
// Returns: Unix timestamp when the transition occurs
func CalculatePosixTransitionTime(year int, trans *PosixTransInfo, offset int64) int64 {
	if trans == nil {
		return 0
	}

	var t Time

	switch trans.Type {
	case 1: // Jn - Julian day (1-365, no leap days)
		// Calculate date from Julian day
		day := trans.Days
		if IsLeapYear(int64(year)) && day >= 60 {
			day++ // Skip Feb 29
		}
		t.Y = int64(year)
		t.M = 1
		t.D = int64(day)
		// Normalize to get correct month/day
		timelib_do_normalize(&t)

	case 2: // n - Julian day (0-365, with leap days)
		t.Y = int64(year)
		t.M = 1
		t.D = int64(trans.Days + 1) // Days are 0-based
		timelib_do_normalize(&t)

	case 3: // Mm.w.d - Month, week, day
		t.Y = int64(year)
		t.M = int64(trans.Mwd.Month)
		t.D = 1

		// Find first day of month
		timelib_do_normalize(&t)

		// Find first occurrence of dow
		dow := int(DayOfWeek(t.Y, t.M, t.D))
		targetDow := trans.Mwd.Dow
		daysToAdd := (targetDow - dow + 7) % 7

		// Add weeks
		daysToAdd += (trans.Mwd.Week - 1) * 7

		// Handle week 5 (last occurrence)
		if trans.Mwd.Week == 5 {
			// Start from the last day of the month and work backwards
			t.M++
			if t.M > 12 {
				t.M = 1
				t.Y++
			}
			t.D = 1
			timelib_do_normalize(&t)
			t.D-- // Last day of previous month
			timelib_do_normalize(&t)

			dow = int(DayOfWeek(t.Y, t.M, t.D))
			daysToSubtract := (dow - targetDow + 7) % 7
			t.D -= int64(daysToSubtract)
			timelib_do_normalize(&t)
		} else {
			t.D += int64(daysToAdd)
			timelib_do_normalize(&t)
		}
	}

	// Set time
	timeInSeconds := trans.Hour
	t.H = int64(timeInSeconds / 3600)
	t.I = int64((timeInSeconds % 3600) / 60)
	t.S = int64(timeInSeconds % 60)

	// Convert to timestamp
	t.UpdateTS(nil)

	// Adjust for timezone offset
	t.Sse -= offset

	return t.Sse
}

// GetPosixTransitionsForYear calculates DST transitions for a specific year
func GetPosixTransitionsForYear(tz *TzInfo, year int64) (beginTs, endTs int64, err error) {
	if tz == nil || tz.PosixInfo == nil {
		return 0, 0, errors.New("no POSIX info available")
	}

	ps := tz.PosixInfo

	if ps.DstBegin == nil || ps.DstEnd == nil {
		return 0, 0, errors.New("no DST transitions defined")
	}

	// Calculate transition times
	beginTs = CalculatePosixTransitionTime(int(year), ps.DstBegin, ps.StdOffset)
	endTs = CalculatePosixTransitionTime(int(year), ps.DstEnd, ps.DstOffset)

	return beginTs, endTs, nil
}
