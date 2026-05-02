package timelib

import "math"

const INT64_MIN = math.MinInt64

const (
	daysPerEra        = 146097
	yearsPerEra       = 400
	hinnantEpochShift = 719468
)

var (
	cDaysInMonth     = [...]int64{31, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	cDaysInMonthLeap = [...]int64{31, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

// Clone creates a deep copy of the Time struct
func (t *Time) Clone() *Time {
	clone := &Time{
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
		Relative:      t.Relative,
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
	return clone
}

// Add adds the relative time information 'interval' to the base time 't'.
// This can be a relative time as created by 'timelib_diff', but also by more
// complex statements such as "next workday".
func (t *Time) Add(interval *RelTime) *Time {
	var bias int64 = 1

	// Clone the time
	result := t.Clone()

	// Handle weekday or special relatives
	if interval.HaveWeekdayRelative || interval.HaveSpecialRelative {
		// Copy entire interval to relative field
		result.Relative = *interval
	} else {
		// Apply bias if interval is inverted
		if interval.Invert {
			bias = -1
		}
		// Clear relative and set individual fields with bias
		result.Relative = RelTime{}
		result.Relative.Y = interval.Y * bias
		result.Relative.M = interval.M * bias
		result.Relative.D = interval.D * bias
		result.Relative.H = interval.H * bias
		result.Relative.I = interval.I * bias
		result.Relative.S = interval.S * bias
		result.Relative.US = interval.US * bias
	}

	result.HaveRelative = true
	result.SseUptodate = false

	// Convert date+time+relative to SSE
	result.UpdateTS(nil)

	// Convert SSE back to date+time fields
	result.UpdateFromSSE()
	result.HaveRelative = false

	return result
}

// Sub subtracts the relative time information 'interval' from the base time 't'.
// Unlike with 'Add', this does not support more complex statements such as "next workday".
func (t *Time) Sub(interval *RelTime) *Time {
	var bias int64 = 1

	// Clone the time
	result := t.Clone()

	// Apply bias if interval is inverted
	if interval.Invert {
		bias = -1
	}

	// Clear relative and set individual fields negated with bias
	result.Relative = RelTime{}
	result.Relative.Y = -(interval.Y * bias)
	result.Relative.M = -(interval.M * bias)
	result.Relative.D = -(interval.D * bias)
	result.Relative.H = -(interval.H * bias)
	result.Relative.I = -(interval.I * bias)
	result.Relative.S = -(interval.S * bias)
	result.Relative.US = -(interval.US * bias)

	result.HaveRelative = true
	result.SseUptodate = false

	// Convert date+time+relative to SSE
	result.UpdateTS(nil)

	// Convert SSE back to date+time fields
	result.UpdateFromSSE()
	result.HaveRelative = false

	return result
}

func sortOldToNew(one, two **Time, rt *RelTime) {
	if (*one).ZoneType == TIMELIB_ZONETYPE_ID && (*two).ZoneType == TIMELIB_ZONETYPE_ID &&
		(*one).TzInfo != nil && (*two).TzInfo != nil &&
		(*one).TzInfo.Name == (*two).TzInfo.Name {
		if (*one).Y > (*two).Y ||
			((*one).Y == (*two).Y && (*one).M > (*two).M) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D > (*two).D) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H > (*two).H) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H == (*two).H && (*one).I > (*two).I) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H == (*two).H && (*one).I == (*two).I && (*one).S > (*two).S) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H == (*two).H && (*one).I == (*two).I && (*one).S == (*two).S && (*one).US > (*two).US) {
			*one, *two = *two, *one
			rt.Invert = true
		}
		return
	}

	if (*one).Sse > (*two).Sse ||
		((*one).Sse == (*two).Sse && (*one).US > (*two).US) {
		*one, *two = *two, *one
		rt.Invert = true
		return
	}

	if (*one).Sse == (*two).Sse && (*one).US == (*two).US {
		if (*one).Y > (*two).Y ||
			((*one).Y == (*two).Y && (*one).M > (*two).M) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D > (*two).D) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H > (*two).H) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H == (*two).H && (*one).I > (*two).I) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H == (*two).H && (*one).I == (*two).I && (*one).S > (*two).S) ||
			((*one).Y == (*two).Y && (*one).M == (*two).M && (*one).D == (*two).D && (*one).H == (*two).H && (*one).I == (*two).I && (*one).S == (*two).S && (*one).US > (*two).US) {
			*one, *two = *two, *one
			rt.Invert = true
		}
	}
}

func timelib_do_rel_normalize(base *Time, rt *RelTime) {
	do_range_limit(0, 1000000, 1000000, &rt.US, &rt.S)
	do_range_limit(0, 60, 60, &rt.S, &rt.I)
	do_range_limit(0, 60, 60, &rt.I, &rt.H)
	do_range_limit(0, 24, 24, &rt.H, &rt.D)
	do_range_limit(0, 12, 12, &rt.M, &rt.Y)

	baseY := base.Y
	baseM := base.M
	doRangeLimitDaysRelative(&baseY, &baseM, &rt.Y, &rt.M, &rt.D, rt.Invert)
	do_range_limit(0, 12, 12, &rt.M, &rt.Y)
}

func diffWithTzid(one, two *Time) *RelTime {
	rt := &RelTime{}
	rt.Invert = false

	sortOldToNew(&one, &two, rt)

	dstCorr := two.Z - one.Z
	dstHCorr := dstCorr / 3600
	dstMCorr := (dstCorr % 3600) / 60

	rt.Y = two.Y - one.Y
	rt.M = two.M - one.M
	rt.D = two.D - one.D
	rt.H = two.H - one.H
	rt.I = two.I - one.I
	rt.S = two.S - one.S
	rt.US = two.US - one.US

	rt.Days = int64(DiffDays(one, two))

	if two.Sse < one.Sse {
		flipped := int64Abs((rt.I*60)+(rt.S) - int64(dstCorr))
		rt.H = flipped / 3600
		rt.I = (flipped - rt.H*3600) / 60
		rt.S = flipped % 60
		rt.Invert = !rt.Invert
	}

	if rt.Invert {
		timelib_do_rel_normalize(one, rt)
	} else {
		timelib_do_rel_normalize(two, rt)
	}

	if one.Dst == 1 && two.Dst == 0 {
		if two.TzInfo != nil {
			if (two.Sse-one.Sse+int64(dstCorr)) < 86400 {
				rt.H -= int64(dstHCorr)
				rt.I -= int64(dstMCorr)
			}
		}
	} else if one.Dst == 0 && two.Dst == 1 {
		if two.TzInfo != nil {
			var transOffset int32
			var transTransitionTime int64
			success := getTimeZoneOffsetInfo(two.Sse, two.TzInfo, &transOffset, &transTransitionTime, nil)
			if success &&
				!((one.Sse+86400 > transTransitionTime) && (one.Sse+86400 <= (transTransitionTime+int64(dstCorr)))) &&
				two.Sse >= transTransitionTime &&
				((two.Sse-one.Sse+int64(dstCorr))%86400) > (two.Sse-transTransitionTime) {
				rt.H -= int64(dstHCorr)
				rt.I -= int64(dstMCorr)
			}
		}
	} else if two.Sse-one.Sse >= 86400 {
		var transOffset int32
		var transTransitionTime int64
		if getTimeZoneOffsetInfo(two.Sse-int64(two.Z), two.TzInfo, &transOffset, &transTransitionTime, nil) {
			dstCorr = one.Z - transOffset
			if two.Sse >= transTransitionTime-int64(dstCorr) && two.Sse < transTransitionTime {
				rt.D--
				rt.H = 24
			}
		}
	}

	return rt
}

func int64Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (t *Time) Diff(other *Time) *RelTime {
	if t.ZoneType == TIMELIB_ZONETYPE_ID && other.ZoneType == TIMELIB_ZONETYPE_ID &&
		t.TzInfo != nil && other.TzInfo != nil &&
		t.TzInfo.Name == other.TzInfo.Name {
		return diffWithTzid(t, other)
	}

	rt := &RelTime{}
	rt.Invert = false

	one := t
	two := other
	sortOldToNew(&one, &two, rt)

	rt.Y = two.Y - one.Y
	rt.M = two.M - one.M
	rt.D = two.D - one.D
	rt.H = two.H - one.H
	if one.ZoneType != TIMELIB_ZONETYPE_ID {
		rt.H += int64(one.Dst)
	}
	if two.ZoneType != TIMELIB_ZONETYPE_ID {
		rt.H -= int64(two.Dst)
	}
	rt.I = two.I - one.I
	rt.S = two.S - one.S - int64(two.Z) + int64(one.Z)
	rt.US = two.US - one.US

	rt.Days = int64(DiffDays(one, two))

	if rt.Invert {
		timelib_do_rel_normalize(one, rt)
	} else {
		timelib_do_rel_normalize(two, rt)
	}

	return rt
}

// DiffDays calculates the difference in full days between two times.
// The result is the number of full days between 'one' and 'two'. It does take
// into account 23 and 25 hour (and variants) days when the zone_type
// is TIMELIB_ZONETYPE_ID and have the same TZID for 'one' and 'two'.
func DiffDays(one, two *Time) int {
	// Matches C function: timelib_diff_days in interval.c
	days := 0

	if SameTimezone(one, two) {
		var earliest, latest *Time
		var earliestTime, latestTime float64

		if TimeCompare(one, two) < 0 {
			earliest = one
			latest = two
		} else {
			earliest = two
			latest = one
		}

		// Convert time of day to decimal hours for comparison
		earliestTime = float64(earliest.H) + float64(earliest.I)/60.0 +
			float64(earliest.S)/3600.0 + float64(earliest.US)/3600000000.0
		latestTime = float64(latest.H) + float64(latest.I)/60.0 +
			float64(latest.S)/3600.0 + float64(latest.US)/3600000000.0

		// Calculate absolute difference in epoch days
		days1 := timelib_epoch_days_from_time(one)
		days2 := timelib_epoch_days_from_time(two)
		days = int(days2 - days1)
		if days < 0 {
			days = -days
		}

		// If latest time-of-day is earlier than earliest time-of-day,
		// we haven't completed a full day, so subtract 1
		if latestTime < earliestTime && days > 0 {
			days--
		}
	} else {
		// Different timezones: use timestamp difference
		// C fix (2295fee): clamp to INT_MAX to avoid overflow
		ddays := (one.Sse - two.Sse) / 86400.0
		if ddays < 0 {
			ddays = -ddays
		}
		if ddays > math.MaxInt32 {
			ddays = math.MaxInt32
		}
		days = int(ddays)
	}

	return days
}

// DoNormalize normalizes the time values (handles overflow/underflow)
// do_range_limit normalizes field 'a' to be within [start, end) and carries overflow/underflow to field 'b'
// This matches the C implementation's do_range_limit function
func do_range_limit(start, end, adj int64, a, b *int64) {
	if *a < start {
		// We calculate 'a + 1' first as 'start - *a - 1' causes an int64 overflow if *a is
		// INT64_MIN. 'start' is 0 in this context, and '0 - INT64_MIN > INT64_MAX'.
		a_plus_1 := *a + 1

		*b -= (start-a_plus_1)/adj + 1

		// This code adds the extra 'adj' separately, as otherwise this can overflow int64 in
		// situations where *b is near INT64_MIN.
		*a += adj * ((start - a_plus_1) / adj)
		*a += adj
	}
	if *a >= end {
		*b += *a / adj
		*a -= adj * (*a / adj)
	}
}

func incMonth(y, m *int64) {
	*m++
	if *m > 12 {
		*m -= 12
		*y++
	}
}

func decMonth(y, m *int64) {
	*m--
	if *m < 1 {
		*m += 12
		*y--
	}
}

func daysInMonthTable(leap bool, month int64) int64 {
	idx := int(month)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(cDaysInMonth) {
		idx = len(cDaysInMonth) - 1
	}
	if leap {
		return cDaysInMonthLeap[idx]
	}
	return cDaysInMonth[idx]
}

func doRangeLimitDaysRelative(baseY, baseM, y, m, d *int64, invert bool) {
	do_range_limit(1, 13, 12, baseM, baseY)

	year := *baseY
	month := *baseM

	if !invert {
		for *d < 0 {
			decMonth(&year, &month)
			days := daysInMonthTable(IsLeapYear(year), month)
			*d += days
			*m--
		}
	} else {
		for *d < 0 {
			days := daysInMonthTable(IsLeapYear(year), month)
			*d += days
			*m--
			incMonth(&year, &month)
		}
	}
}

func doRangeLimitDays(y, m, d *int64) bool {
	if *d >= daysPerEra || *d <= -daysPerEra {
		*y += yearsPerEra * (*d / daysPerEra)
		*d -= daysPerEra * (*d / daysPerEra)
	}

	do_range_limit(1, 13, 12, m, y)

	leapYear := IsLeapYear(*y)
	daysPerMonthCurrentYear := cDaysInMonth
	if leapYear {
		daysPerMonthCurrentYear = cDaysInMonthLeap
	}

	retval := false

	for *d <= 0 && *m > 0 {
		previousMonth := *m - 1
		previousYear := *y
		if previousMonth < 1 {
			previousMonth += 12
			previousYear = (*y) - 1
		}
		leapPrev := IsLeapYear(previousYear)
		var daysPrev int64
		if leapPrev {
			daysPrev = cDaysInMonthLeap[int(previousMonth)]
		} else {
			daysPrev = cDaysInMonth[int(previousMonth)]
		}

		*d += daysPrev
		*m--
		retval = true
	}

	// Re-evaluate days-per-month for the (potentially) updated current year
	if IsLeapYear(*y) {
		daysPerMonthCurrentYear = cDaysInMonthLeap
	} else {
		daysPerMonthCurrentYear = cDaysInMonth
	}

	for *d > 0 && *m <= 12 && *d > daysPerMonthCurrentYear[int(*m)] {
		*d -= daysPerMonthCurrentYear[int(*m)]
		*m++
		retval = true
	}

	return retval
}

func magicDateCalc(t *Time) {
	if t.D < -719498 {
		return
	}

	g := t.D + hinnantEpochShift - 1
	y := (10000*g + 14780) / 3652425
	ddd := g - ((365 * y) + (y / 4) - (y / 100) + (y / 400))
	if ddd < 0 {
		y--
		ddd = g - ((365 * y) + (y / 4) - (y / 100) + (y / 400))
	}
	mi := (100*ddd + 52) / 3060
	mm := ((mi + 2) % 12) + 1
	y = y + (mi+2)/12
	dd := ddd - ((mi*306 + 5) / 10) + 1

	t.Y = y
	t.M = mm
	t.D = dd
}

// DoNormalize normalizes a Time in-place, handling overflow/underflow of
// time components (microseconds, seconds, minutes, hours, days, months).
// This matches the C function: DoNormalize in tm2unixtime.c
func DoNormalize(t *Time) {
	if t.US != TIMELIB_UNSET {
		do_range_limit(0, 1000000, 1000000, &t.US, &t.S)
	}
	if t.S != TIMELIB_UNSET {
		do_range_limit(0, 60, 60, &t.S, &t.I)
		do_range_limit(0, 60, 60, &t.I, &t.H)
		do_range_limit(0, 24, 24, &t.H, &t.D)
	}

	do_range_limit(1, 13, 12, &t.M, &t.Y)

	if t.Y == 1970 && t.M == 1 && t.D != 1 {
		magicDateCalc(t)
	}

	for {
		if !doRangeLimitDays(&t.Y, &t.M, &t.D) {
			break
		}
	}

	do_range_limit(1, 13, 12, &t.M, &t.Y)
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

	// Initialize UNSET time fields to 0 before calculation (matching C behavior where fields are 0 by default)
	h := t.H
	if h == TIMELIB_UNSET {
		h = 0
	}
	i := t.I
	if i == TIMELIB_UNSET {
		i = 0
	}
	s := t.S
	if s == TIMELIB_UNSET {
		s = 0
	}

	// Calculate epoch days and SSE
	// Note: This is done in two steps to avoid overflow (matching C at tm2unixtime.c:473-481)
	// The C comment explains: "timelib_epoch_days_from_time(time) * SECS_PER_DAY with the lowest
	// limit of timelib_epoch_days_from_time() is less than the range of an int64_t. This then
	// overflows. In order to be able to still allow for any time in that day that only halfly
	// fits in the int64_t range, we add the time element first, which is always positive, and
	// then twice half the value of the earliest day as expressed as unix timestamp."
	epochDays := timelib_epoch_days_from_time(t)

	// First add the time of day (always positive)
	t.Sse = h*3600 + i*60 + s
	// Then add epoch days in two halves to avoid overflow
	t.Sse += epochDays * (SECS_PER_DAY / 2)
	t.Sse += epochDays * (SECS_PER_DAY / 2)

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

		inTransition = (actualTransitionTime != INT64_MIN &&
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

// doAdjustForWeekday adjusts the day based on relative weekday (e.g., "next Monday")
func doAdjustForWeekday(t *Time) {
	currentDow := DayOfWeek(t.Y, t.M, t.D)

	if t.Relative.WeekdayBehavior == 2 {
		// "this" behavior - stay in current week
		// To make "this week" work, where the current DOW is a "sunday"
		if currentDow == 0 && t.Relative.Weekday != 0 {
			t.Relative.Weekday -= 7
		}
		// To make "sunday this week" work, where the current DOW is not a "sunday"
		if t.Relative.Weekday == 0 && currentDow != 0 {
			t.Relative.Weekday = 7
		}
		t.D -= currentDow
		t.D += int64(t.Relative.Weekday)
		return
	}

	// Calculate difference between target and current day of week
	difference := int64(t.Relative.Weekday) - currentDow

	// Adjust difference based on whether we're going forward or backward
	if (t.Relative.D < 0 && difference < 0) ||
		(t.Relative.D >= 0 && difference <= -int64(t.Relative.WeekdayBehavior)) {
		difference += 7
	}

	if t.Relative.Weekday >= 0 {
		t.D += difference
	} else {
		// Negative weekday (shouldn't happen in our use case, but handle it)
		absWeekday := int64(t.Relative.Weekday)
		if absWeekday < 0 {
			absWeekday = -absWeekday
		}
		t.D -= (7 - (absWeekday - currentDow))
	}

	// Clear the weekday relative flag after adjustment (C behavior)
	t.Relative.HaveWeekdayRelative = false
}

// doAdjustRelative applies relative time adjustments
func doAdjustRelative(t *Time) {
	if t.Relative.HaveWeekdayRelative {
		doAdjustForWeekday(t)
	}

	DoNormalize(t)

	if t.HaveRelative {
		t.US += t.Relative.US
		t.S += t.Relative.S
		t.I += t.Relative.I
		t.H += t.Relative.H

		// Only add D if it's not UNSET
		if t.D != TIMELIB_UNSET {
			t.D += t.Relative.D
		}
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

	DoNormalize(t)
}

// doAdjustSpecial handles special relative time adjustments (late)
func doAdjustSpecial(t *Time) {
	if t.Relative.HaveSpecialRelative {
		switch t.Relative.Special.Type {
		case TIMELIB_SPECIAL_WEEKDAY:
			doAdjustSpecialWeekday(t)
		}
	}

	DoNormalize(t)
	// Clear special relative
	t.Relative.Special.Type = 0
	t.Relative.Special.Amount = 0
}

// doAdjustSpecialWeekday handles weekday special adjustments
// This matches the C function: do_adjust_special_weekday
func doAdjustSpecialWeekday(t *Time) {
	count := t.Relative.Special.Amount
	dow := DayOfWeek(t.Y, t.M, t.D)

	// Add increments of 5 weekdays as a week, leaving the DOW unchanged
	t.D += (count / 5) * 7

	// Deal with the remainder
	rem := count % 5

	if count > 0 {
		if rem == 0 {
			// Head back to Friday if we stop on the weekend
			if dow == 0 {
				t.D -= 2
			} else if dow == 6 {
				t.D -= 1
			}
		} else if dow == 6 {
			// We ended up on Saturday, but there's still work to do, so move
			// to Sunday and continue from there
			t.D += 1
		} else if dow+rem > 5 {
			// We're on a weekday, but we're going past Friday, so skip right
			// over the weekend
			t.D += 2
		}
	} else {
		// Completely mirror the forward direction. This also covers the 0
		// case, since if we start on the weekend, we want to move forward as
		// if we stopped there while going backwards
		if rem == 0 {
			if dow == 6 {
				t.D += 2
			} else if dow == 0 {
				t.D += 1
			}
		} else if dow == 0 {
			t.D -= 1
		} else if dow+rem < 1 {
			t.D -= 2
		}
	}

	t.D += rem
}

// doAdjustSpecialEarly handles special relative time adjustments (early)
func doAdjustSpecialEarly(t *Time) {
	if t.Relative.HaveSpecialRelative {
		switch t.Relative.Special.Type {
		case TIMELIB_SPECIAL_DAY_OF_WEEK_IN_MONTH:
			t.D = 1
			t.M += t.Relative.M
			t.Relative.M = 0
		case TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH:
			t.D = 1
			t.M += t.Relative.M + 1
			t.Relative.M = 0
		}
	}
	switch t.Relative.FirstLastDayOf {
	case TIMELIB_SPECIAL_FIRST_DAY_OF_MONTH:
		t.D = 1
	case TIMELIB_SPECIAL_LAST_DAY_OF_MONTH:
		t.D = 0
		t.M++
	}
	DoNormalize(t)
}

// Unixtime2date converts Unix timestamp to date
// This matches the C function: timelib_unixtime2date
func Unixtime2date(ts int64, y, m, d *int64) {
	// Calculate days since algorithm's epoch (0000-03-01)
	// HINNANT_EPOCH_SHIFT = 719468 (days from 0000-03-01 to 1970-01-01)
	days := ts / SECS_PER_DAY
	days += 719468
	//fmt.Printf("DEBUG Unixtime2date: ts=%d, days=%d\n", ts, days)

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
	//fmt.Printf("DEBUG Unixtime2date result: y=%d m=%d d=%d\n", year, month, day)
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
	// Note: Do NOT set US here - C version preserves existing microseconds
	t.Z = 0
	t.Dst = 0
	t.Sse = ts
	t.SseUptodate = true
	t.TimUptodate = true
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
		zoneType := t.ZoneType
		t.Unixtime2gmt(ts + int64(t.Z) + int64(t.Dst*3600))
		t.Sse = ts
		t.Z = z
		t.Dst = dst
		t.ZoneType = zoneType
		// Mark as localtime with zone info set (C code at unixtime2tm.c:159-160)
		t.IsLocaltime = true
		t.HaveZone = true
		return

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
				t.TzAbbr = gmtOffset.Abbr
				t.ZoneType = TIMELIB_ZONETYPE_ID
				// Mark as localtime with zone info set (C code at unixtime2tm.c:159-160)
				t.IsLocaltime = true
				t.HaveZone = true
				return
			} else {
				t.Unixtime2gmt(ts)
				t.Sse = ts
				t.ZoneType = TIMELIB_ZONETYPE_ID
			}
		} else {
			t.Unixtime2gmt(ts)
			t.Sse = ts
			// ZoneType remains TIMELIB_ZONETYPE_ID as set before
		}

	default:
		t.Unixtime2gmt(ts)
		t.Sse = ts
		// For default case, set have_zone to false (C code at unixtime2tm.c:154-156)
		t.IsLocaltime = false
		t.HaveZone = false
	}
}

// UpdateFromSSE updates time from seconds since epoch
// This is the Go equivalent of timelib_update_from_sse
func (t *Time) UpdateFromSSE() {
	sse := t.Sse
	z := t.Z
	dst := t.Dst
	zoneType := t.ZoneType

	switch t.ZoneType {
	case TIMELIB_ZONETYPE_ABBR, TIMELIB_ZONETYPE_OFFSET:
		// For abbreviations and fixed offsets, just add the offset
		//fmt.Printf("DEBUG UpdateFromSSE: t.Dst=%d, t.Z=%d, t.Sse=%d\n", t.Dst, t.Z, t.Sse)
		param := t.Sse + int64(t.Z) + int64(t.Dst*3600)
		//fmt.Printf("DEBUG UpdateFromSSE: Calling Unixtime2gmt(%d) = %d + %d + %d\n", param, t.Sse, t.Z, t.Dst*3600)
		t.Unixtime2gmt(param)

	case TIMELIB_ZONETYPE_ID:
		// For timezone IDs, get the actual offset at this timestamp
		if t.TzInfo != nil {
			var offset int32 = 0
			result := GetTimeZoneOffsetInfo(t.Sse, t.TzInfo, &offset, nil, nil)
			if result == 0 {
				// Fallback: if offset lookup fails, use the saved Z value
				offset = z
			}
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
	t.ZoneType = zoneType
}
