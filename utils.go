package timelib

import (
	"fmt"
)

// TimeDtor frees time structure resources
func TimeDtor(t *Time) {
	if t == nil {
		return
	}

	// Clear string fields
	t.TzAbbr = ""

	// Clear timezone info
	if t.TzInfo != nil {
		TzinfoDtor(t.TzInfo)
		t.TzInfo = nil
	}

	// Reset all fields to default values
	t.Y = TIMELIB_UNSET
	t.M = TIMELIB_UNSET
	t.D = TIMELIB_UNSET
	t.H = TIMELIB_UNSET
	t.I = TIMELIB_UNSET
	t.S = TIMELIB_UNSET
	t.US = 0
	t.Z = 0
	t.Dst = -1
	t.Sse = 0
	t.HaveTime = false
	t.HaveDate = false
	t.HaveZone = false
	t.HaveRelative = false
	t.HaveWeeknrDay = false
	t.SseUptodate = false
	t.TimUptodate = false
	t.IsLocaltime = false
	t.ZoneType = TIMELIB_ZONETYPE_NONE
}

// TimeClone creates a deep copy of a Time structure
func TimeClone(orig *Time) *Time {
	if orig == nil {
		return nil
	}

	clone := &Time{
		Y:             orig.Y,
		M:             orig.M,
		D:             orig.D,
		H:             orig.H,
		I:             orig.I,
		S:             orig.S,
		US:            orig.US,
		Z:             orig.Z,
		TzAbbr:        orig.TzAbbr,
		Dst:           orig.Dst,
		Relative:      orig.Relative, // RelTime is a value type, so this copies it
		Sse:           orig.Sse,
		HaveTime:      orig.HaveTime,
		HaveDate:      orig.HaveDate,
		HaveZone:      orig.HaveZone,
		HaveRelative:  orig.HaveRelative,
		HaveWeeknrDay: orig.HaveWeeknrDay,
		SseUptodate:   orig.SseUptodate,
		TimUptodate:   orig.TimUptodate,
		IsLocaltime:   orig.IsLocaltime,
		ZoneType:      orig.ZoneType,
	}

	// Clone timezone info if present
	if orig.TzInfo != nil {
		clone.TzInfo = TzinfoClone(orig.TzInfo)
	}

	return clone
}

// RelTimeDtor frees relative time structure resources
func RelTimeDtor(rt *RelTime) {
	if rt == nil {
		return
	}

	// Reset all fields to default values
	rt.Y = 0
	rt.M = 0
	rt.D = 0
	rt.H = 0
	rt.I = 0
	rt.S = 0
	rt.US = 0
	rt.Weekday = -1
	rt.WeekdayBehavior = 0
	rt.FirstLastDayOf = 0
	rt.Invert = false
	rt.Days = 0
	rt.Special.Type = 0
	rt.Special.Amount = 0
	rt.HaveWeekdayRelative = false
	rt.HaveSpecialRelative = false
}

// RelTimeClone creates a deep copy of a RelTime structure
func RelTimeClone(orig *RelTime) *RelTime {
	if orig == nil {
		return nil
	}

	return &RelTime{
		Y:                   orig.Y,
		M:                   orig.M,
		D:                   orig.D,
		H:                   orig.H,
		I:                   orig.I,
		S:                   orig.S,
		US:                  orig.US,
		Weekday:             orig.Weekday,
		WeekdayBehavior:     orig.WeekdayBehavior,
		FirstLastDayOf:      orig.FirstLastDayOf,
		Invert:              orig.Invert,
		Days:                orig.Days,
		Special:             orig.Special,
		HaveWeekdayRelative: orig.HaveWeekdayRelative,
		HaveSpecialRelative: orig.HaveSpecialRelative,
	}
}

// TimeOffsetDtor frees time offset structure resources
func TimeOffsetDtor(to *TimeOffset) {
	if to == nil {
		return
	}

	to.Abbr = ""
	to.Offset = 0
	to.LeapSecs = 0
	to.IsDst = 0
	to.TransitionTime = 0
}

// ErrorContainerDtor frees error container resources
func ErrorContainerDtor(errors *ErrorContainer) {
	if errors == nil {
		return
	}

	// Clear error messages
	errors.ErrorMessages = errors.ErrorMessages[:0]
	errors.WarningMessages = errors.WarningMessages[:0]
	errors.ErrorCount = 0
	errors.WarningCount = 0
}

// FillHoles fills gaps in parsed time with reference time information
func FillHoles(parsed, now *Time, options int) {
	if parsed == nil || now == nil {
		return
	}

	// If have_date but not have_time, reset time to midnight (C behavior)
	// This is crucial for weekday relatives like "next Saturday"
	if (options&TIMELIB_OVERRIDE_TIME == 0) && parsed.HaveDate && !parsed.HaveTime {
		parsed.H = 0
		parsed.I = 0
		parsed.S = 0
		parsed.US = 0
	}

	// Handle microseconds according to C implementation
	if parsed.Y != TIMELIB_UNSET || parsed.M != TIMELIB_UNSET || parsed.D != TIMELIB_UNSET ||
		parsed.H != TIMELIB_UNSET || parsed.I != TIMELIB_UNSET || parsed.S != TIMELIB_UNSET {
		if parsed.US == TIMELIB_UNSET {
			parsed.US = 0
		}
	} else {
		if parsed.US == TIMELIB_UNSET {
			if now.US != TIMELIB_UNSET {
				parsed.US = now.US
			} else {
				parsed.US = 0
			}
		}
	}

	// Fill unset date/time fields with reference time values
	if parsed.Y == TIMELIB_UNSET {
		parsed.Y = now.Y
	}
	if parsed.M == TIMELIB_UNSET {
		parsed.M = now.M
	}
	if parsed.D == TIMELIB_UNSET {
		parsed.D = now.D
	}

	// Handle time fields
	if options&TIMELIB_OVERRIDE_TIME == 0 {
		if parsed.H == TIMELIB_UNSET {
			parsed.H = now.H
		}
		if parsed.I == TIMELIB_UNSET {
			parsed.I = now.I
		}
		if parsed.S == TIMELIB_UNSET {
			parsed.S = now.S
		}
	} else {
		// Override time - set to zeros if not set
		if parsed.H == TIMELIB_UNSET {
			parsed.H = 0
		}
		if parsed.I == TIMELIB_UNSET {
			parsed.I = 0
		}
		if parsed.S == TIMELIB_UNSET {
			parsed.S = 0
		}
	}

	// Handle timezone information
	if parsed.Z == 0 && now.Z != 0 {
		parsed.Z = now.Z
	}
	if parsed.Dst == -1 && now.Dst != -1 {
		parsed.Dst = now.Dst
	}
	if parsed.TzAbbr == "" && now.TzAbbr != "" {
		parsed.TzAbbr = now.TzAbbr
	}

	// Handle timezone info
	if parsed.TzInfo == nil && now.TzInfo != nil {
		if options&TIMELIB_NO_CLONE == 0 {
			parsed.TzInfo = TzinfoClone(now.TzInfo)
		} else {
			parsed.TzInfo = now.TzInfo
		}
	}
}

// DoRelNormalize normalizes relative time values
func DoRelNormalize(base *Time, rt *RelTime) {
	if base == nil || rt == nil {
		return
	}

	// Normalize microseconds
	if rt.US >= 1000000 {
		extraSeconds := rt.US / 1000000
		rt.S += extraSeconds
		rt.US = rt.US % 1000000
	} else if rt.US < 0 {
		rt.S--
		rt.US += 1000000
	}

	// Normalize seconds
	if rt.S >= 60 {
		extraMinutes := rt.S / 60
		rt.I += extraMinutes
		rt.S = rt.S % 60
	} else if rt.S < 0 {
		rt.I--
		rt.S += 60
	}

	// Normalize minutes
	if rt.I >= 60 {
		extraHours := rt.I / 60
		rt.H += extraHours
		rt.I = rt.I % 60
	} else if rt.I < 0 {
		rt.H--
		rt.I += 60
	}

	// Normalize hours
	if rt.H >= 24 {
		extraDays := rt.H / 24
		rt.D += extraDays
		rt.H = rt.H % 24
	} else if rt.H < 0 {
		rt.D--
		rt.H += 24
	}

	// Normalize months
	if rt.M >= 12 {
		extraYears := rt.M / 12
		rt.Y += extraYears
		rt.M = rt.M % 12
	} else if rt.M < 0 {
		rt.Y--
		rt.M += 12
	}

	// Normalize days taking into account the base date and leap years
	for rt.D < 0 {
		// Go back one month
		base.M--
		if base.M < 1 {
			base.M = 12
			base.Y--
		}
		rt.D += DaysInMonth(base.Y, base.M)
	}

	for rt.D >= DaysInMonth(base.Y, base.M) {
		rt.D -= DaysInMonth(base.Y, base.M)
		base.M++
		if base.M > 12 {
			base.M = 1
			base.Y++
		}
	}
}

// DumpDate displays debugging information about date/time
func DumpDate(d *Time, options int) {
	if d == nil {
		fmt.Println("Time structure is nil")
		return
	}

	fmt.Printf("Date/Time: %04d-%02d-%02d %02d:%02d:%02d.%06d\n",
		d.Y, d.M, d.D, d.H, d.I, d.S, d.US)
	fmt.Printf("Timezone: offset=%d, dst=%d, abbr='%s', type=%d\n",
		d.Z, d.Dst, d.TzAbbr, d.ZoneType)
	fmt.Printf("Flags: have_time=%v, have_date=%v, have_zone=%v, have_relative=%v\n",
		d.HaveTime, d.HaveDate, d.HaveZone, d.HaveRelative)
	fmt.Printf("Status: sse_uptodate=%v, tim_uptodate=%v, is_localtime=%v\n",
		d.SseUptodate, d.TimUptodate, d.IsLocaltime)
	fmt.Printf("SSE: %d\n", d.Sse)

	if options&1 != 0 {
		fmt.Printf("Relative: y=%d, m=%d, d=%d, h=%d, i=%d, s=%d, us=%d\n",
			d.Relative.Y, d.Relative.M, d.Relative.D,
			d.Relative.H, d.Relative.I, d.Relative.S, d.Relative.US)
	}

	if options&2 != 0 {
		fmt.Printf("Zone type: %d\n", d.ZoneType)
	}
}

// DumpRelTime displays debugging information about relative time
func DumpRelTime(d *RelTime) {
	if d == nil {
		fmt.Println("RelTime structure is nil")
		return
	}

	fmt.Printf("Relative Time: y=%d, m=%d, d=%d, h=%d, i=%d, s=%d, us=%d\n",
		d.Y, d.M, d.D, d.H, d.I, d.S, d.US)
	fmt.Printf("Weekday: %d, behavior: %d, first_last_day_of: %d\n",
		d.Weekday, d.WeekdayBehavior, d.FirstLastDayOf)
	fmt.Printf("Flags: invert=%v, days=%d\n", d.Invert, d.Days)
	fmt.Printf("Special: type=%d, amount=%d\n", d.Special.Type, d.Special.Amount)
	fmt.Printf("Relative flags: have_weekday_relative=%v, have_special_relative=%v\n",
		d.HaveWeekdayRelative, d.HaveSpecialRelative)
}

// TsToJulianday converts Unix timestamp to Julian Day
func TsToJulianday(ts int64) float64 {
	// Julian Day starts from -4714-11-24T12:00:00 UTC
	// Unix timestamp starts from 1970-01-01T00:00:00 UTC
	// Difference is 2440588 days
	return float64(ts)/86400.0 + 2440587.5
}

// TsToJ2000 converts Unix timestamp to J2000 epoch
func TsToJ2000(ts int64) float64 {
	// J2000 epoch starts from 2000-01-01T12:00:00 UTC
	// Unix timestamp starts from 1970-01-01T00:00:00 UTC
	// Difference is 10957.5 days
	return float64(ts)/86400.0 - 10957.5
}

// AddWall adds relative time to base time (wall time version)
// This version properly handles timezone transitions and DST changes
func (t *Time) AddWall(interval *RelTime) *Time {
	// Clone the input time
	result := TimeClone(t)

	result.HaveRelative = true
	result.SseUptodate = false

	bias := int64(1)
	if interval.Invert {
		bias = -1
	}

	// Handle special relative times (weekday, etc.)
	if interval.HaveWeekdayRelative || interval.HaveSpecialRelative {
		// Copy the interval to result's relative field
		result.Relative = *interval

		result.UpdateTS(nil)
		result.UpdateFromSSE()
	} else {
		// Clear relative field and set Y/M/D
		result.Relative = RelTime{}
		result.Relative.Y = interval.Y * bias
		result.Relative.M = interval.M * bias
		result.Relative.D = interval.D * bias

		// Update timestamp if we have year/month/day changes
		if result.Relative.Y != 0 || result.Relative.M != 0 || result.Relative.D != 0 {
			result.UpdateTS(nil)
		}

		// Add time component
		if interval.US == 0 {
			// Simple case: no microseconds
			// Save base microseconds before UpdateFromSSE (which clears them)
			baseUS := result.US
			result.Sse += bias * hmsToSeconds(interval.H, interval.I, interval.S)
			result.UpdateFromSSE()
			// Restore base microseconds
			result.US = baseUS
		} else {
			// Complex case: with microseconds
			tempInterval := RelTimeClone(interval)

			// Normalize microseconds into seconds
			doRangeLimit(0, 1000000, 1000000, &tempInterval.US, &tempInterval.S)

			// Save base microseconds before UpdateFromSSE (which clears them)
			baseUS := result.US
			result.Sse += bias * hmsToSeconds(tempInterval.H, tempInterval.I, tempInterval.S)
			result.UpdateFromSSE()
			// Add interval microseconds to base microseconds
			result.US = baseUS + tempInterval.US*bias

			timelib_do_normalize(result)
			result.UpdateTS(nil)
		}
		timelib_do_normalize(result)
	}

	// If we have a timezone ID, set the timezone properly
	if result.ZoneType == TIMELIB_ZONETYPE_ID && result.TzInfo != nil {
		SetTimezone(result, result.TzInfo)
	}
	result.HaveRelative = false

	return result
}

// hmsToSeconds converts hours, minutes, seconds to total seconds
func hmsToSeconds(h, i, s int64) int64 {
	return (h * 3600) + (i * 60) + s
}

// doRangeLimit normalizes values within a range
// Matches C function: do_range_limit from interval.c
func doRangeLimit(start, end, adj int64, a, b *int64) {
	if *a < start {
		*b -= (start - *a - 1) / adj + 1
		*a += adj * ((start - *a - 1) / adj + 1)
	}
	if *a >= end {
		*b += *a / adj
		*a -= adj * (*a / adj)
	}
}

// SubWall subtracts relative time from base time (wall time version)
// SubWall subtracts a relative time interval from the base time using wall clock semantics.
// This properly handles DST transitions by working with wall clock time (local time)
// rather than just epoch seconds.
//
// The function handles three cases:
// 1. Weekday/special relative times: applies them through UpdateTS
// 2. Simple intervals (no microseconds): adjusts SSE directly, then updates time fields
// 3. Intervals with microseconds: adjusts both SSE and microsecond fields separately
//
// This matches the C function: timelib_sub_wall in interval.c
func (t *Time) SubWall(interval *RelTime) *Time {
	bias := int64(1)
	result := TimeClone(t)

	result.HaveRelative = true
	result.SseUptodate = false

	if interval.HaveWeekdayRelative || interval.HaveSpecialRelative {
		// For weekday/special relative times, copy the interval to relative field
		// and let UpdateTS handle it
		result.Relative = *interval
		result.UpdateTS(nil)
		result.UpdateFromSSE()
	} else {
		// Handle the invert flag
		if interval.Invert {
			bias = -1
		}

		// Clear relative field and set up date component subtraction
		result.Relative = RelTime{}
		result.Relative.Y = 0 - (interval.Y * bias)
		result.Relative.M = 0 - (interval.M * bias)
		result.Relative.D = 0 - (interval.D * bias)

		// If we have date components, update the timestamp
		if result.Relative.Y != 0 || result.Relative.M != 0 || result.Relative.D != 0 {
			result.UpdateTS(nil)
		}

		// Handle time components (H/I/S/US)
		if interval.US == 0 {
			// Simple case: no microseconds, just adjust SSE
			// Save base microseconds before UpdateFromSSE (which clears them)
			baseUS := result.US
			hmsSeconds := hmsToSeconds(interval.H, interval.I, interval.S)
			result.Sse -= bias * hmsSeconds
			result.UpdateFromSSE()
			// Restore base microseconds
			result.US = baseUS
		} else {
			// Complex case: have microseconds
			// Clone the interval and normalize microseconds
			tempInterval := RelTimeClone(interval)

			// Normalize microseconds into seconds using do_range_limit logic
			doRangeLimit(0, 1000000, 1000000, &tempInterval.US, &tempInterval.S)

			// Save base microseconds before UpdateFromSSE (which clears them)
			baseUS := result.US
			// Adjust SSE by H/I/S
			hmsSeconds := hmsToSeconds(tempInterval.H, tempInterval.I, tempInterval.S)
			result.Sse -= bias * hmsSeconds
			result.UpdateFromSSE()

			// Adjust microseconds: subtract interval microseconds from base microseconds
			result.US = baseUS - tempInterval.US*bias

			// Normalize the time fields
			timelib_do_normalize(result)
			result.UpdateTS(nil)
		}

		// Final normalization
		timelib_do_normalize(result)
	}

	// Re-apply timezone if it's a timezone ID
	if result.ZoneType == TIMELIB_ZONETYPE_ID {
		SetTimezone(result, result.TzInfo)
	}

	result.HaveRelative = false
	return result
}

// ParseZone parses timezone information from string
func ParseZone(ptr *string, dst *int, t *Time, tzNotFound *int, tzdb *TzDB, tzWrapper func(string, *TzDB, *int) (*TzInfo, error)) int64 {
	if ptr == nil || t == nil || *ptr == "" {
		if tzNotFound != nil {
			*tzNotFound = 1
		}
		return 0
	}

	s := *ptr
	var retval int64 = 0
	parenCount := 0

	if tzNotFound != nil {
		*tzNotFound = 0
	}

	// Skip leading whitespace and opening parentheses
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '(') {
		if s[i] == '(' {
			parenCount++
		}
		i++
	}

	// Check for GMT+/- prefix
	if i+3 <= len(s) && s[i:i+3] == "GMT" {
		i += 3
		// Skip optional whitespace after GMT
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
	}

	// Handle +/- offset
	if i < len(s) && s[i] == '+' {
		i++
		t.IsLocaltime = true
		t.ZoneType = TIMELIB_ZONETYPE_OFFSET
		t.Dst = 0

		offset, consumed, found := parseTzCor(s[i:])
		i += consumed
		if !found && tzNotFound != nil {
			*tzNotFound = 1
		}
		retval = offset
	} else if i < len(s) && s[i] == '-' {
		i++
		t.IsLocaltime = true
		t.ZoneType = TIMELIB_ZONETYPE_OFFSET
		t.Dst = 0

		offset, consumed, found := parseTzCor(s[i:])
		i += consumed
		if !found && tzNotFound != nil {
			*tzNotFound = 1
		}
		retval = -offset
	} else {
		// Try to parse timezone abbreviation or identifier
		t.IsLocaltime = true

		// First, lookup by abbreviation only
		abbr, consumed, offset, dstVal, found := lookupAbbr(s[i:])
		i += consumed

		if found {
			t.ZoneType = TIMELIB_ZONETYPE_ABBR
			t.Dst = dstVal
			if dst != nil {
				*dst = dstVal
			}
			t.TzAbbr = abbr
		}

		// If not found or if it's UTC, try as timezone identifier
		// This matches C code: if (!found || strcmp("UTC", tz_abbr) == 0)
		if !found || abbr == "UTC" {
			if tzdb != nil && tzWrapper != nil {
				var dummyErrorCode int
				res, err := tzWrapper(abbr, tzdb, &dummyErrorCode)
				if err == nil && res != nil {
					t.TzInfo = res
					t.ZoneType = TIMELIB_ZONETYPE_ID
					found = true
				}
			}
		}

		if !found && tzNotFound != nil {
			*tzNotFound = 1
		}
		retval = offset
	}

	// Skip trailing closing parentheses
	for parenCount > 0 && i < len(s) && s[i] == ')' {
		i++
		parenCount--
	}

	// Update the pointer to point to the remaining string
	*ptr = s[i:]

	return retval
}

// parseTzCor parses timezone correction (offset) like +0100, -05:00, etc.
func parseTzCor(s string) (offset int64, consumed int, found bool) {
	begin := 0
	end := 0

	// Find the end of the numeric/colon sequence
	for end < len(s) && (isDigit(s[end]) || s[end] == ':') {
		end++
	}

	if end == begin {
		return 0, 0, false
	}

	length := end - begin
	switch length {
	case 1, 2: // H or HH
		hours := parseInt(s[begin:end])
		return hours * 3600, end, true

	case 3: // H:M
		if s[begin+1] == ':' {
			hours := parseInt(s[begin : begin+1])
			mins := parseInt(s[begin+2 : end])
			return hours*3600 + mins*60, end, true
		}
		return 0, end, false

	case 4: // H:MM, HH:M, or HHMM
		if s[begin+1] == ':' {
			hours := parseInt(s[begin : begin+1])
			mins := parseInt(s[begin+2 : end])
			return hours*3600 + mins*60, end, true
		} else if s[begin+2] == ':' {
			hours := parseInt(s[begin : begin+2])
			mins := parseInt(s[begin+3 : end])
			return hours*3600 + mins*60, end, true
		} else {
			// HHMM
			val := parseInt(s[begin:end])
			hours := val / 100
			mins := val % 100
			return hours*3600 + mins*60, end, true
		}

	case 5: // HH:MM
		if s[begin+2] != ':' {
			return 0, end, false
		}
		hours := parseInt(s[begin : begin+2])
		mins := parseInt(s[begin+3 : end])
		return hours*3600 + mins*60, end, true

	case 6: // HHMMSS
		val := parseInt(s[begin:end])
		hours := val / 10000
		mins := (val / 100) % 100
		secs := val % 100
		return hours*3600 + mins*60 + secs, end, true

	case 8: // HH:MM:SS
		if s[begin+2] != ':' || s[begin+5] != ':' {
			return 0, end, false
		}
		hours := parseInt(s[begin : begin+2])
		mins := parseInt(s[begin+3 : begin+5])
		secs := parseInt(s[begin+6 : end])
		return hours*3600 + mins*60 + secs, end, true
	}

	return 0, end, false
}

// lookupAbbr looks up timezone abbreviation
func lookupAbbr(s string) (abbr string, consumed int, offset int64, dst int, found bool) {
	begin := 0
	end := 0

	// Only include A-Z, a-z, 0-9, /, _, -, and + in abbreviations/TZ IDs
	for end < len(s) {
		c := s[end]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '/' || c == '_' || c == '-' || c == '+') {
			break
		}
		end++
	}

	if end == begin {
		return "", 0, 0, 0, false
	}

	word := s[begin:end]

	// Try to look up in abbreviation table
	offset, dst, found = lookupTimezoneAbbr(word)

	return word, end, offset, dst, found
}

// lookupTimezoneAbbr looks up common timezone abbreviations
func lookupTimezoneAbbr(abbr string) (offset int64, dst int, found bool) {
	// This is a simplified version - the full implementation would use the abbr_search table
	// For now, we'll handle some common abbreviations

	switch abbr {
	case "UTC", "GMT":
		return 0, 0, true
	case "EST":
		return -5 * 3600, 0, true
	case "EDT":
		return -4 * 3600, 0, true  // EDT is UTC-4 (total offset, dst flag not added)
	case "CST":
		return -6 * 3600, 0, true
	case "CDT":
		return -5 * 3600, 0, true  // CDT is UTC-5 (total offset, dst flag not added)
	case "MST":
		return -7 * 3600, 0, true
	case "MDT":
		return -6 * 3600, 0, true  // MDT is UTC-6 (total offset, dst flag not added)
	case "PST":
		return -8 * 3600, 0, true
	case "PDT":
		return -7 * 3600, 0, true  // PDT is UTC-7 (total offset, dst flag not added)
	case "CET":
		return 1 * 3600, 0, true
	case "CEST":
		return 2 * 3600, 0, true   // CEST is UTC+2 (total offset, dst flag not added)
	case "BST":
		return 1 * 3600, 0, true   // BST is UTC+1 (total offset, dst flag not added)
	default:
		// Try to parse as timezone identifier (will be handled by caller)
		return 0, 0, false
	}
}

// isDigit checks if a byte is a digit
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// parseInt parses a string to int64
func parseInt(s string) int64 {
	var result int64
	for i := 0; i < len(s); i++ {
		if isDigit(s[i]) {
			result = result*10 + int64(s[i]-'0')
		}
	}
	return result
}

// Constants for FillHoles options
const (
	TIMELIB_NONE          = 0x00
	TIMELIB_OVERRIDE_TIME = 0x01
	TIMELIB_NO_CLONE      = 0x02
)
