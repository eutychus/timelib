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

	// Handle microseconds
	if parsed.US == 0 && now.US != 0 {
		if options&TIMELIB_OVERRIDE_TIME == 0 {
			parsed.US = now.US
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
func (t *Time) AddWall(interval *RelTime) *Time {
	// For now, this is the same as Add
	// Wall time version would handle timezone transitions differently
	return t.Add(interval)
}

// SubWall subtracts relative time from base time (wall time version)
func (t *Time) SubWall(interval *RelTime) *Time {
	// For now, this is the same as Sub
	// Wall time version would handle timezone transitions differently
	return t.Sub(interval)
}

// ParseZone parses timezone information from string
func ParseZone(ptr *string, dst *int, t *Time, tzNotFound *int, tzdb *TzDB, tzWrapper func(string, *TzDB, *int) (*TzInfo, error)) int64 {
	if ptr == nil || t == nil {
		return 0
	}

	// This is a simplified timezone parser
	// Full implementation would parse various timezone formats

	// For now, return 0 offset
	if dst != nil {
		*dst = 0
	}
	if tzNotFound != nil {
		*tzNotFound = 1
	}

	return 0
}

// Constants for FillHoles options
const (
	TIMELIB_NONE          = 0x00
	TIMELIB_OVERRIDE_TIME = 0x01
	TIMELIB_NO_CLONE      = 0x02
)
