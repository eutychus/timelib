package timelib

import "math"

const INT64_MIN = math.MinInt64

// Add adds the relative time information 'interval' to the base time 't'.
// This can be a relative time as created by 'timelib_diff', but also by more
// complex statements such as "next workday".
func (t *Time) Add(interval *RelTime) *Time {
	result := &Time{
		Y:             t.Y,
		M:             t.M,
		D:             t.D,
		H:             t.H,
		I:             t.I,
		S:             t.S,
		US:            t.US,
		Z:             t.Z,
		TzAbbr:        t.TzAbbr,
		TzInfo:        t.TzInfo,
		Dst:           t.Dst,
		Relative:      t.Relative, // RelTime is a value type, so this copies it
		Sse:           t.Sse,
		HaveTime:      t.HaveTime,
		HaveDate:      t.HaveDate,
		HaveZone:      t.HaveZone,
		HaveRelative:  t.HaveRelative,
		HaveWeeknrDay: t.HaveWeeknrDay,
		SseUptodate:   t.SseUptodate,
		TimUptodate:   t.TimUptodate,
		IsLocaltime:   t.IsLocaltime,
		ZoneType:      t.ZoneType,
	}

	// Add years
	if interval.Y != 0 {
		result.Y += interval.Y
	}

	// Add months
	if interval.M != 0 {
		result.M += interval.M
		// Handle month overflow/underflow
		for result.M > 12 {
			result.M -= 12
			result.Y++
		}
		for result.M < 1 {
			result.M += 12
			result.Y--
		}
	}

	// Add days
	if interval.D != 0 {
		result.D += interval.D
		// Normalize days (handle month boundaries)
		timelib_do_normalize(result)
	}

	// Add hours, minutes, seconds
	if interval.H != 0 {
		result.H += interval.H
	}
	if interval.I != 0 {
		result.I += interval.I
	}
	if interval.S != 0 {
		result.S += interval.S
	}
	if interval.US != 0 {
		result.US += interval.US
	}

	// Handle microseconds overflow
	if result.US >= 1000000 || result.US < 0 {
		extraSeconds := result.US / 1000000
		result.S += extraSeconds
		result.US = result.US % 1000000
		if result.US < 0 {
			result.US += 1000000
			result.S--
		}
	}

	// Handle seconds overflow
	if result.S >= 60 || result.S < 0 {
		extraMinutes := result.S / 60
		result.I += extraMinutes
		result.S = result.S % 60
		if result.S < 0 {
			result.S += 60
			result.I--
		}
	}

	// Handle minutes overflow
	if result.I >= 60 || result.I < 0 {
		extraHours := result.I / 60
		result.H += extraHours
		result.I = result.I % 60
		if result.I < 0 {
			result.I += 60
			result.H--
		}
	}

	// Handle hours overflow
	if result.H >= 24 || result.H < 0 {
		extraDays := result.H / 24
		result.D += extraDays
		result.H = result.H % 24
		if result.H < 0 {
			result.H += 24
			result.D--
		}
	}

	// Final normalization for any day overflow from hour adjustments
	if result.D < 1 || result.D > DaysInMonth(result.Y, result.M) {
		timelib_do_normalize(result)
	}

	return result
}

// Sub subtracts the relative time information 'interval' from the base time 't'.
// Unlike with 'Add', this does not support more complex statements such as "next workday".
func (t *Time) Sub(interval *RelTime) *Time {
	// Create inverted interval
	inverted := &RelTime{
		Y:                   interval.Y,
		M:                   interval.M,
		D:                   interval.D,
		H:                   interval.H,
		I:                   interval.I,
		S:                   interval.S,
		US:                  interval.US,
		Weekday:             interval.Weekday,
		WeekdayBehavior:     interval.WeekdayBehavior,
		FirstLastDayOf:      interval.FirstLastDayOf,
		Invert:              interval.Invert,
		Days:                interval.Days,
		Special:             interval.Special,
		HaveWeekdayRelative: interval.HaveWeekdayRelative,
		HaveSpecialRelative: interval.HaveSpecialRelative,
	}
	inverted.Invert = !inverted.Invert

	// Negate all values
	inverted.Y = -inverted.Y
	inverted.M = -inverted.M
	inverted.D = -inverted.D
	inverted.H = -inverted.H
	inverted.I = -inverted.I
	inverted.S = -inverted.S
	inverted.US = -inverted.US

	return t.Add(inverted)
}

// Diff calculates the difference between two times.
// The result is a timelib_rel_time structure that describes how you can
// convert from 'one' to 'two' with 'timelib_add'. This does *not* necessarily
// mean that you can go from 'two' to 'one' by using 'timelib_sub' due to the
// way months and days are calculated.
func (t *Time) Diff(other *Time) *RelTime {
	// Determine if we need to invert the result
	if TimeCompare(t, other) > 0 {
		// t is after other, so we need to calculate the difference in reverse
		// and mark it as inverted
		temp := &RelTime{}

		// Calculate differences in reverse (t is after other)
		temp.Y = t.Y - other.Y
		temp.M = t.M - other.M
		temp.D = t.D - other.D
		temp.H = t.H - other.H
		temp.I = t.I - other.I
		temp.S = t.S - other.S
		temp.US = t.US - other.US

		// Handle negative differences
		if temp.US < 0 {
			temp.S--
			temp.US += 1000000
		}
		if temp.S < 0 {
			temp.I--
			temp.S += 60
		}
		if temp.I < 0 {
			temp.H--
			temp.I += 60
		}
		if temp.H < 0 {
			temp.D--
			temp.H += 24
		}

		// Handle month/year overflow - only if we have negative months
		if temp.M < 0 {
			temp.Y--
			temp.M += 12
		}

		// Handle day overflow - only if we have negative days
		if temp.D < 0 {
			// Get days in previous month of the "from" time (t)
			prevMonth := t.M - 1
			prevYear := t.Y
			if prevMonth < 1 {
				prevMonth += 12
				prevYear--
			}
			temp.D += DaysInMonth(prevYear, prevMonth)
			temp.M-- // Decrease month since we borrowed days

			// Handle month underflow after borrowing
			if temp.M < 0 {
				temp.Y--
				temp.M += 12
			}
		}

		// Calculate total days difference
		temp.Days = int64(timelib_diff_days(other, t))
		temp.Invert = true

		return temp
	}

	diff := &RelTime{}

	// Calculate differences (other is after t)
	diff.Y = other.Y - t.Y
	diff.M = other.M - t.M
	diff.D = other.D - t.D
	diff.H = other.H - t.H
	diff.I = other.I - t.I
	diff.S = other.S - t.S
	diff.US = other.US - t.US

	// Handle negative differences
	if diff.US < 0 {
		diff.S--
		diff.US += 1000000
	}
	if diff.S < 0 {
		diff.I--
		diff.S += 60
	}
	if diff.I < 0 {
		diff.H--
		diff.I += 60
	}
	if diff.H < 0 {
		diff.D--
		diff.H += 24
	}

	// Handle month/year overflow - only if we have negative months
	if diff.M < 0 {
		diff.Y--
		diff.M += 12
	}

	// Handle day overflow - only if we have negative days
	if diff.D < 0 {
		// Get days in previous month of the "to" time (other)
		prevMonth := other.M - 1
		prevYear := other.Y
		if prevMonth < 1 {
			prevMonth += 12
			prevYear--
		}
		diff.D += DaysInMonth(prevYear, prevMonth)
		diff.M-- // Decrease month since we borrowed days

		// Handle month underflow after borrowing
		if diff.M < 0 {
			diff.Y--
			diff.M += 12
		}
	}

	// Calculate total days difference
	diff.Days = int64(timelib_diff_days(t, other))

	return diff
}

// DiffDays calculates the difference in full days between two times.
// The result is the number of full days between 'one' and 'two'. It does take
// into account 23 and 25 hour (and variants) days when the zone_type
// is TIMELIB_ZONETYPE_ID and have the same TZID for 'one' and 'two'.
func timelib_diff_days(one, two *Time) int {
	// Convert both times to epoch days
	days1 := timelib_epoch_days_from_time(one)
	days2 := timelib_epoch_days_from_time(two)

	diff := int(days2 - days1)
	if diff < 0 {
		diff = -diff
	}

	return diff
}

// timelib_do_normalize normalizes the time values (handles overflow/underflow)
func timelib_do_normalize(t *Time) {
	// Normalize microseconds
	if t.US >= 1000000 {
		extraSeconds := t.US / 1000000
		t.S += extraSeconds
		t.US = t.US % 1000000
	} else if t.US < 0 {
		t.S--
		t.US += 1000000
	}

	// Normalize seconds
	if t.S >= 60 {
		extraMinutes := t.S / 60
		t.I += extraMinutes
		t.S = t.S % 60
	} else if t.S < 0 {
		t.I--
		t.S += 60
	}

	// Normalize minutes
	if t.I >= 60 {
		extraHours := t.I / 60
		t.H += extraHours
		t.I = t.I % 60
	} else if t.I < 0 {
		t.H--
		t.I += 60
	}

	// Normalize hours
	if t.H >= 24 {
		extraDays := t.H / 24
		t.D += extraDays
		t.H = t.H % 24
	} else if t.H < 0 {
		t.D--
		t.H += 24
	}

	// Normalize months first (handle month overflow/underflow)
	for t.M > 12 {
		t.M -= 12
		t.Y++
	}
	for t.M < 1 {
		t.M += 12
		t.Y--
	}

	// Normalize days and months
	for t.D > DaysInMonth(t.Y, t.M) {
		t.D -= DaysInMonth(t.Y, t.M)
		t.M++
		if t.M > 12 {
			t.M = 1
			t.Y++
		}
	}

	for t.D < 1 {
		t.M--
		if t.M < 1 {
			t.M = 12
			t.Y--
		}
		t.D += DaysInMonth(t.Y, t.M)
	}
}

// timelib_epoch_days_from_time calculates epoch days from time
func timelib_epoch_days_from_time(t *Time) int64 {
	// Use the algorithm from howardhinnant.github.io/date_algorithms.html
	// This is a simplified version - the full implementation would be more complex
	y := t.Y
	m := t.M
	d := t.D

	// Adjust for January/February
	if m <= 2 {
		y--
		m += 12
	}

	// Calculate epoch days
	var era int64
	if y >= 0 {
		era = y / 400
	} else {
		era = (y - 399) / 400
	}
	yoe := y - era*400
	doy := (153*(m-3)+2)/5 + d - 1
	doe := yoe*365 + yoe/4 - yoe/100 + doy
	epochDays := era*146097 + doe - 719468

	return epochDays
}

// UpdateTS updates the timestamp from date/time fields
// This is the Go equivalent of timelib_update_ts
func (t *Time) UpdateTS(tzi *TzInfo) {
	// Adjust for special relative times (early adjustments)
	doAdjustSpecialEarly(t)

	// Adjust for relative time
	doAdjustRelative(t)

	// Adjust for special relative times (late adjustments)
	doAdjustSpecial(t)

	// Calculate epoch days
	epochDays := timelib_epoch_days_from_time(t)

	// Calculate seconds since epoch
	// Split into two parts to avoid overflow
	seconds := t.H*3600 + t.I*60 + t.S
	seconds += epochDays * (86400 / 2)
	seconds += epochDays * (86400 / 2)

	t.Sse = seconds

	// Adjust for timezone - this modifies t.Sse
	doAdjustTimezone(t, tzi)
	t.SseUptodate = true
	t.HaveRelative = false
	t.Relative.HaveWeekdayRelative = false
	t.Relative.HaveSpecialRelative = false
	t.Relative.FirstLastDayOf = 0
}

// doAdjustTimezone adjusts the SSE based on timezone information
// This is a direct port of the C function from tm2unixtime.c
func doAdjustTimezone(tz *Time, tzi *TzInfo) {
	const SECS_PER_HOUR = 3600

	switch tz.ZoneType {
	case TIMELIB_ZONETYPE_OFFSET:
		tz.IsLocaltime = true
		tz.Sse += -int64(tz.Z)
		return

	case TIMELIB_ZONETYPE_ABBR:
		tz.IsLocaltime = true
		tz.Sse += (-int64(tz.Z) - int64(tz.Dst*SECS_PER_HOUR))
		return

	case TIMELIB_ZONETYPE_ID:
		tzi = tz.TzInfo
		fallthrough

	default:
		// No timezone in struct, fallback to reference if possible
		var currentOffset int32 = 0
		var currentTransitionTime int64 = 0
		var currentIsDst uint = 0
		var afterOffset int32 = 0
		var afterTransitionTime int64 = 0
		var adjustment int64
		var inTransition bool
		var actualOffset int32
		var actualTransitionTime int64

		if tzi == nil {
			return
		}

		getTimeZoneOffsetInfo(tz.Sse, tzi, &currentOffset, &currentTransitionTime, &currentIsDst)
		getTimeZoneOffsetInfo(tz.Sse-int64(currentOffset), tzi, &afterOffset, &afterTransitionTime, nil)
		actualOffset = afterOffset
		actualTransitionTime = afterTransitionTime

		if currentOffset == afterOffset && tz.HaveZone {
			// Make sure we're not missing a DST change because we don't know the actual offset yet
			if currentOffset >= 0 && tz.Dst != 0 && currentIsDst == 0 {
				// Timezone or its DST at or east of UTC
				var earlierOffset int32
				var earlierTransitionTime int64
				getTimeZoneOffsetInfo(tz.Sse-int64(currentOffset)-7200, tzi, &earlierOffset, &earlierTransitionTime, nil)
				if earlierOffset != afterOffset && tz.Sse-int64(earlierOffset) < afterTransitionTime {
					actualOffset = earlierOffset
					actualTransitionTime = earlierTransitionTime
				}
			} else if currentOffset <= 0 && currentIsDst != 0 && tz.Dst == 0 {
				// Timezone west of UTC
				var laterOffset int32
				var laterTransitionTime int64
				getTimeZoneOffsetInfo(tz.Sse-int64(currentOffset)+7200, tzi, &laterOffset, &laterTransitionTime, nil)
				if laterOffset != afterOffset && tz.Sse-int64(laterOffset) >= laterTransitionTime {
					actualOffset = laterOffset
					actualTransitionTime = laterTransitionTime
				}
			}
		}

		tz.IsLocaltime = true

		inTransition = (
			actualTransitionTime != INT64_MIN &&
				((tz.Sse - int64(actualOffset)) >= (actualTransitionTime + int64(currentOffset-actualOffset))) &&
				((tz.Sse - int64(actualOffset)) < actualTransitionTime))

		if currentOffset != actualOffset && !inTransition {
			adjustment = -int64(actualOffset)
		} else {
			adjustment = -int64(currentOffset)
		}

		tz.Sse += adjustment
		SetTimezone(tz, tzi)
		return
	}
}

// getTimeZoneOffsetInfo gets timezone offset information for a given timestamp
// This is a port of the C function from parse_tz.c
func getTimeZoneOffsetInfo(ts int64, tz *TzInfo, offset *int32, transitionTime *int64, isDst *uint) bool {
	if tz == nil {
		return false
	}

	to, tmpTransitionTime := fetchTimezoneOffset(tz, ts)
	if to != nil {
		if offset != nil {
			*offset = to.Offset
		}
		if isDst != nil {
			*isDst = uint(to.IsDst)
		}
		if transitionTime != nil {
			*transitionTime = tmpTransitionTime
		}
		return true
	}
	return false
}

// fetchPosixTimezoneOffset calculates timezone offset using POSIX string.
// This handles dates beyond the transition table using rules like "EST5EDT,M3.2.0,M11.1.0".
//
// POSIX rules allow timezone files to specify ongoing DST rules that extend
// indefinitely into the future, even when the tzfile only contains historical
// transitions. This function calculates transitions for year-1, year, and year+1
// to correctly handle timestamps near year boundaries.
//
// Matches C function: timelib_fetch_posix_timezone_offset
func fetchPosixTimezoneOffset(tz *TzInfo, ts int64) (*TTInfo, int64) {
	if tz.PosixInfo == nil {
		return nil, 0
	}

	// If there is no second (dst_end) information, the UTC offset is valid for the whole year
	if tz.PosixInfo.DstEnd == nil {
		if len(tz.Trans) > 0 {
			return &tz.Type[tz.PosixInfo.TypeIndexStdType], tz.Trans[tz.Bit64.Timecnt-1]
		}
		return &tz.Type[tz.PosixInfo.TypeIndexStdType], INT64_MIN
	}

	// Find 'year' (UTC) for 'ts'
	var dummy Time
	dummy.Unixtime2gmt(ts)
	year := dummy.Y

	// Calculate transition times for 'year-1', 'year', and 'year+1'
	transitions := PosixTransitions{}
	GetTransitionsForYear(tz, year-1, &transitions)
	GetTransitionsForYear(tz, year, &transitions)
	GetTransitionsForYear(tz, year+1, &transitions)

	// Check where the 'ts' falls in the transitions
	for i := 1; i < transitions.Count; i++ {
		if ts < transitions.Times[i] {
			typeIdx := transitions.Types[i-1]
			if typeIdx >= 0 && int(typeIdx) < len(tz.Type) {
				return &tz.Type[typeIdx], transitions.Times[i-1]
			}
			break
		}
	}

	return nil, 0
}

// fetchTimezoneOffset finds the timezone offset for a given timestamp
// This is a port of timelib_fetch_timezone_offset from parse_tz.c
func fetchTimezoneOffset(tz *TzInfo, ts int64) (*TTInfo, int64) {
	if tz == nil {
		return nil, 0
	}

	// If there are no transitions, use type 0 or POSIX info
	if tz.Bit64.Timecnt == 0 || len(tz.Trans) == 0 {
		if tz.Bit64.Typecnt == 1 && len(tz.Type) > 0 {
			return &tz.Type[0], INT64_MIN
		}
		return nil, 0
	}

	// Before first transition, use type 0
	if ts < tz.Trans[0] {
		if len(tz.Type) > 0 {
			return &tz.Type[0], INT64_MIN
		}
		return nil, 0
	}

	// After last transition, use POSIX info or last transition
	// RFC 8536: Timestamps ON the last transition use that transition's type
	// Only timestamps AFTER the last transition use POSIX rules
	if ts > tz.Trans[tz.Bit64.Timecnt-1] {
		// Use POSIX info if available
		if tz.PosixInfo != nil {
			return fetchPosixTimezoneOffset(tz, ts)
		}

		// Fall back to last transition type
		transIdx := tz.TransIdx[tz.Bit64.Timecnt-1]
		if int(transIdx) < len(tz.Type) {
			return &tz.Type[transIdx], tz.Trans[tz.Bit64.Timecnt-1]
		}
		return nil, 0
	}

	// If timestamp exactly equals the last transition, use it
	if ts == tz.Trans[tz.Bit64.Timecnt-1] {
		transIdx := tz.TransIdx[tz.Bit64.Timecnt-1]
		if int(transIdx) < len(tz.Type) {
			return &tz.Type[transIdx], tz.Trans[tz.Bit64.Timecnt-1]
		}
		return nil, 0
	}

	// Binary search for the right transition
	// RFC 8536: The type corresponding to a transition time specifies local time
	// for timestamps starting at the given transition time and continuing up to,
	// but not including, the next transition time.
	left := uint64(0)
	right := tz.Bit64.Timecnt - 1

	for right-left > 1 {
		mid := (left + right) >> 1
		if ts < tz.Trans[mid] {
			right = mid
		} else {
			left = mid
		}
	}

	// If timestamp exactly equals the right transition, use right instead of left
	// This handles the case where ts is exactly at a transition boundary
	if ts == tz.Trans[right] {
		transIdx := tz.TransIdx[right]
		if int(transIdx) < len(tz.Type) {
			return &tz.Type[transIdx], tz.Trans[right]
		}
	}

	transIdx := tz.TransIdx[left]
	if int(transIdx) < len(tz.Type) {
		return &tz.Type[transIdx], tz.Trans[left]
	}
	return nil, 0
}

// doAdjustRelative applies relative time adjustments
func doAdjustRelative(t *Time) {
	if t.Relative.HaveWeekdayRelative {
		// TODO: implement weekday adjustments
	}

	timelib_do_normalize(t)

	if t.HaveRelative {
		t.US += t.Relative.US
		t.S += t.Relative.S
		t.I += t.Relative.I
		t.H += t.Relative.H
		t.D += t.Relative.D
		t.M += t.Relative.M
		t.Y += t.Relative.Y
	}

	// Handle first/last day of month
	switch t.Relative.FirstLastDayOf {
	case TIMELIB_SPECIAL_FIRST_DAY_OF_MONTH:
		t.D = 1
	case TIMELIB_SPECIAL_LAST_DAY_OF_MONTH:
		t.D = 0
		t.M++
	}

	timelib_do_normalize(t)
}

// doAdjustSpecial handles special relative time adjustments (late)
func doAdjustSpecial(t *Time) {
	if t.Relative.HaveSpecialRelative {
		switch t.Relative.Special.Type {
		case TIMELIB_SPECIAL_WEEKDAY:
			// TODO: implement weekday special adjustments
		}
	}
	timelib_do_normalize(t)
	// Clear special relative
	t.Relative.Special.Type = 0
	t.Relative.Special.Amount = 0
}

// doAdjustSpecialEarly handles special relative time adjustments (early)
func doAdjustSpecialEarly(t *Time) {
	if t.Relative.HaveSpecialRelative {
		switch t.Relative.Special.Type {
		case TIMELIB_SPECIAL_DAY_OF_WEEK_IN_MONTH:
			// TODO: implement day of week in month adjustments
		}
	}
}

// Unixtime2date converts Unix timestamp to date
// This matches the C function: timelib_unixtime2date
func Unixtime2date(ts int64, y, m, d *int64) {
	// Calculate days since algorithm's epoch (0000-03-01)
	// HINNANT_EPOCH_SHIFT = 719468 (days from 0000-03-01 to 1970-01-01)
	days := ts / SECS_PER_DAY
	days += 719468

	// Adjustment for negative time portion
	// If the time-of-day portion is negative, we need to go back one day
	t := ts % SECS_PER_DAY
	if t < 0 {
		days--
	}

	// Calculate year, month, and day using Howard Hinnant's algorithm
	// http://howardhinnant.github.io/date_algorithms.html#civil_from_days
	var era int64
	if days >= 0 {
		era = days / 146097
	} else {
		era = (days - 146097 + 1) / 146097
	}

	dayOfEra := days - era*146097
	yearOfEra := (dayOfEra - dayOfEra/1460 + dayOfEra/36524 - dayOfEra/146096) / 365
	year := yearOfEra + era*400
	dayOfYear := dayOfEra - (365*yearOfEra + yearOfEra/4 - yearOfEra/100)
	monthPortion := (5*dayOfYear + 2) / 153
	day := dayOfYear - (153*monthPortion+2)/5 + 1

	var month int64
	if monthPortion < 10 {
		month = monthPortion + 3
	} else {
		month = monthPortion - 9
	}

	if month <= 2 {
		year++
	}

	*y = year
	*m = month
	*d = day
}

// Unixtime2gmt converts Unix timestamp to GMT
// This matches the C function: timelib_unixtime2gmt
func (t *Time) Unixtime2gmt(ts int64) {
	var y, m, d int64
	Unixtime2date(ts, &y, &m, &d)

	t.Y = y
	t.M = m
	t.D = d

	// Calculate remaining time (time of day portion)
	// For negative timestamps, we need to handle the remainder carefully
	remainder := ts % SECS_PER_DAY
	if remainder < 0 {
		remainder += SECS_PER_DAY
	}

	hours := remainder / 3600
	minutes := (remainder - hours*3600) / 60
	seconds := remainder % 60

	t.H = hours
	t.I = minutes
	t.S = seconds
	t.US = 0
	t.Z = 0
	t.Dst = 0
	t.Sse = ts
	t.SseUptodate = true
	t.IsLocaltime = false
	t.ZoneType = TIMELIB_ZONETYPE_NONE
}

// Unixtime2local converts Unix timestamp to local time
// This matches the C function: timelib_unixtime2local
func (t *Time) Unixtime2local(ts int64) {
	switch t.ZoneType {
	case TIMELIB_ZONETYPE_ABBR, TIMELIB_ZONETYPE_OFFSET:
		// For fixed offsets, just add the offset
		z := t.Z
		dst := t.Dst
		t.Unixtime2gmt(ts + int64(t.Z) + int64(t.Dst*3600))
		t.Sse = ts
		t.Z = z
		t.Dst = dst

	case TIMELIB_ZONETYPE_ID:
		if t.TzInfo != nil {
			// Get timezone offset for this timestamp
			gmtOffset := GetTimeZoneInfo(ts, t.TzInfo)
			if gmtOffset != nil {
				t.Unixtime2gmt(ts + int64(gmtOffset.Offset))

				// Restore original timestamp and set timezone info
				t.Sse = ts
				t.Dst = gmtOffset.IsDst
				t.Z = gmtOffset.Offset
				t.TzInfo = t.TzInfo
				t.TzAbbr = gmtOffset.Abbr
			} else {
				t.Unixtime2gmt(ts)
				t.Sse = ts
			}
		} else {
			t.Unixtime2gmt(ts)
			t.Sse = ts
		}

	default:
		t.Unixtime2gmt(ts)
		t.Sse = ts
	}
}

// UpdateFromSSE updates time from seconds since epoch
// This is the Go equivalent of timelib_update_from_sse
func (t *Time) UpdateFromSSE() {
	sse := t.Sse
	z := t.Z
	dst := t.Dst

	switch t.ZoneType {
	case TIMELIB_ZONETYPE_ABBR, TIMELIB_ZONETYPE_OFFSET:
		// For abbreviations and fixed offsets, just add the offset
		t.Unixtime2gmt(t.Sse + int64(t.Z) + int64(t.Dst*3600))

	case TIMELIB_ZONETYPE_ID:
		// For timezone IDs, get the actual offset at this timestamp
		if t.TzInfo != nil {
			offset := int32(0)
			GetTimeZoneOffsetInfo(t.Sse, t.TzInfo, &offset, nil, nil)
			t.Unixtime2gmt(t.Sse + int64(offset))
		} else {
			t.Unixtime2gmt(t.Sse)
		}

	default:
		t.Unixtime2gmt(t.Sse)
	}

	// Restore original values
	t.Sse = sse
	t.IsLocaltime = true
	t.HaveZone = true
	t.Z = z
	t.Dst = dst
}
