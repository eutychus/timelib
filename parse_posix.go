package timelib

import (
	"errors"
	"fmt"
	"strconv"
)

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
	name, offset, newPos, err := parsePosixName(posix, pos)
	if err != nil {
		return nil, fmt.Errorf("failed to parse standard name: %v", err)
	}
	ps.Std = name
	ps.StdOffset = offset
	pos = newPos

	// Check if there's DST information
	if pos >= len(posix) {
		// No DST, create type info
		if tz != nil {
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
		name, offset, newPos, err = parsePosixName(posix, pos)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DST name: %v", err)
		}
		ps.Dst = name
		if offset != 0 {
			ps.DstOffset = offset
		} else {
			ps.DstOffset = ps.StdOffset + 3600 // Default: 1 hour ahead
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

	// Set type indices
	if tz != nil && len(tz.Type) > 0 {
		// Find standard and DST types
		ps.TypeIndexStdType = 0
		ps.TypeIndexDstType = 1

		// Search for actual types that match
		for i, t := range tz.Type {
			if t.IsDst == 0 {
				ps.TypeIndexStdType = i
			} else if t.IsDst != 0 {
				ps.TypeIndexDstType = i
			}
		}
	}

	return ps, nil
}

// parsePosixName parses a timezone name and optional offset
// Returns name, offset (in seconds, POSIX convention), new position
func parsePosixName(s string, pos int) (string, int64, int, error) {
	if pos >= len(s) {
		return "", 0, pos, errors.New("unexpected end of string")
	}

	var name string
	var offset int64

	// Check if name is quoted with < >
	if s[pos] == '<' {
		pos++
		end := pos
		for end < len(s) && s[end] != '>' {
			end++
		}
		if end >= len(s) {
			return "", 0, pos, errors.New("unclosed < in timezone name")
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
			return "", 0, pos, errors.New("expected timezone name")
		}
		name = s[pos:end]
		pos = end
	}

	// Parse offset if present
	if pos < len(s) && (s[pos] == '+' || s[pos] == '-' || (s[pos] >= '0' && s[pos] <= '9')) {
		var err error
		offset, pos, err = parsePosixOffset(s, pos)
		if err != nil {
			return name, 0, pos, err
		}
	}

	return name, offset, pos, nil
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
		Hour: 2, // Default transition time is 2:00 AM
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
// in a specific year
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
