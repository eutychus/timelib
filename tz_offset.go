package timelib

import (
	"errors"
	"sort"
)

// GetOffsetInfo returns timezone offset information for a given timestamp
func GetOffsetInfo(ts int64, tzi *TzInfo) (*TTInfo, int64, error) {
	if tzi == nil {
		return nil, 0, errors.New("no timezone info")
	}

	// If no transitions, use first type
	if len(tzi.Trans) == 0 {
		if len(tzi.Type) > 0 {
			return &tzi.Type[0], ts, nil
		}
		return nil, 0, errors.New("no timezone types defined")
	}

	// Binary search for the transition
	idx := sort.Search(len(tzi.Trans), func(i int) bool {
		return tzi.Trans[i] > ts
	})

	// If before first transition
	if idx == 0 {
		if len(tzi.Type) > 0 {
			typeIdx := 0
			if len(tzi.TransIdx) > 0 {
				typeIdx = int(tzi.TransIdx[0])
			}
			if typeIdx < len(tzi.Type) {
				return &tzi.Type[typeIdx], tzi.Trans[0], nil
			}
		}
		return nil, 0, errors.New("no valid type for timestamp")
	}

	// Use the transition before this timestamp
	idx--

	// If after last transition, check POSIX rules for future DST
	if idx >= len(tzi.Trans)-1 && tzi.PosixInfo != nil {
		return getPosixOffsetInfo(ts, tzi)
	}

	// Get the type index for this transition
	if idx >= len(tzi.TransIdx) {
		return nil, 0, errors.New("transition index out of range")
	}

	typeIdx := int(tzi.TransIdx[idx])
	if typeIdx >= len(tzi.Type) {
		return nil, 0, errors.New("type index out of range")
	}

	transTime := tzi.Trans[idx]
	return &tzi.Type[typeIdx], transTime, nil
}

// getPosixOffsetInfo calculates offset using POSIX rules for future dates
func getPosixOffsetInfo(ts int64, tzi *TzInfo) (*TTInfo, int64, error) {
	if tzi.PosixInfo == nil {
		// No POSIX rules, use last defined type
		if len(tzi.Type) > 0 {
			lastIdx := len(tzi.TransIdx) - 1
			if lastIdx >= 0 {
				typeIdx := int(tzi.TransIdx[lastIdx])
				if typeIdx < len(tzi.Type) {
					return &tzi.Type[typeIdx], tzi.Trans[len(tzi.Trans)-1], nil
				}
			}
		}
		return nil, 0, errors.New("no POSIX rules and no types")
	}

	ps := tzi.PosixInfo

	// If no DST, always use standard time
	if ps.DstBegin == nil || ps.DstEnd == nil {
		// Find standard type
		for i := range tzi.Type {
			if tzi.Type[i].IsDst == 0 {
				return &tzi.Type[i], ts, nil
			}
		}
		// Fall back to first type
		if len(tzi.Type) > 0 {
			return &tzi.Type[0], ts, nil
		}
		return nil, 0, errors.New("no types available")
	}

	// Calculate DST transitions for the year containing this timestamp
	// Convert timestamp to year
	var t Time
	t.Unixtime2gmt(ts)
	year := t.Y

	// Get transition times for this year
	beginTs, endTs, err := GetPosixTransitionsForYear(tzi, year)
	if err != nil {
		// Fall back to standard time
		for i := range tzi.Type {
			if tzi.Type[i].IsDst == 0 {
				return &tzi.Type[i], ts, nil
			}
		}
		if len(tzi.Type) > 0 {
			return &tzi.Type[0], ts, nil
		}
		return nil, 0, err
	}

	// Determine if we're in DST
	var isDst bool
	if beginTs < endTs {
		// Northern hemisphere: DST is between begin and end
		isDst = ts >= beginTs && ts < endTs
	} else {
		// Southern hemisphere: DST wraps around year boundary
		isDst = ts >= beginTs || ts < endTs
	}

	// Find appropriate type
	var targetType *TTInfo
	var transTime int64

	if isDst {
		// Use DST type
		transTime = beginTs
		if ps.TypeIndexDstType >= 0 && ps.TypeIndexDstType < len(tzi.Type) {
			targetType = &tzi.Type[ps.TypeIndexDstType]
		} else {
			// Search for DST type
			for i := range tzi.Type {
				if tzi.Type[i].IsDst != 0 {
					targetType = &tzi.Type[i]
					break
				}
			}
		}
	} else {
		// Use standard time type
		transTime = endTs
		if ps.TypeIndexStdType >= 0 && ps.TypeIndexStdType < len(tzi.Type) {
			targetType = &tzi.Type[ps.TypeIndexStdType]
		} else {
			// Search for standard type
			for i := range tzi.Type {
				if tzi.Type[i].IsDst == 0 {
					targetType = &tzi.Type[i]
					break
				}
			}
		}
	}

	if targetType == nil {
		// Fall back to first type
		if len(tzi.Type) > 0 {
			targetType = &tzi.Type[0]
		} else {
			return nil, 0, errors.New("no types available")
		}
	}

	return targetType, transTime, nil
}

// GetCurrentOffsetForTime returns the current UTC offset for the given time
func GetCurrentOffsetForTime(t *Time) (int32, error) {
	if t == nil {
		return 0, errors.New("time is nil")
	}

	switch t.ZoneType {
	case TIMELIB_ZONETYPE_OFFSET:
		return t.Z, nil

	case TIMELIB_ZONETYPE_ABBR:
		return t.Z, nil

	case TIMELIB_ZONETYPE_ID:
		if t.TzInfo != nil {
			// Make sure we have a timestamp
			if !t.SseUptodate {
				t.UpdateTS(nil)
			}

			info, _, err := GetOffsetInfo(t.Sse, t.TzInfo)
			if err != nil {
				return 0, err
			}
			return info.Offset, nil
		}
		return 0, errors.New("no timezone info")

	default:
		return 0, nil
	}
}

// IsTimestampInDST checks if timestamp is in DST for the given timezone
func IsTimestampInDST(ts int64, tzi *TzInfo) (int, error) {
	info, _, err := GetOffsetInfo(ts, tzi)
	if err != nil {
		return -1, err
	}

	return info.IsDst, nil
}

// GetTimeZoneInfoForTime returns detailed timezone offset information
func GetTimeZoneInfoForTime(ts int64, tzi *TzInfo) (offset int32, transitionTime int64, isDst int, abbr string, err error) {
	if tzi == nil {
		return 0, 0, 0, "", errors.New("no timezone info")
	}

	info, transTime, err := GetOffsetInfo(ts, tzi)
	if err != nil {
		return 0, 0, 0, "", err
	}

	// Get abbreviation
	abbrStr := ""
	if info.AbbrIdx >= 0 && info.AbbrIdx < len(tzi.TimezoneAbbr) {
		// Find null terminator
		end := info.AbbrIdx
		for end < len(tzi.TimezoneAbbr) && tzi.TimezoneAbbr[end] != '\x00' {
			end++
		}
		abbrStr = tzi.TimezoneAbbr[info.AbbrIdx:end]
	}

	return info.Offset, transTime, info.IsDst, abbrStr, nil
}

// ApplyOffsetToTime applies timezone offset to convert from UTC to local time
func ApplyOffsetToTime(t *Time) error {
	if t == nil {
		return errors.New("time is nil")
	}

	if t.TzInfo == nil {
		return nil // No timezone to apply
	}

	// Get offset
	offset, err := GetCurrentOffsetForTime(t)
	if err != nil {
		return err
	}

	// Apply offset to timestamp
	t.Sse += int64(offset)

	// Convert back to date/time components
	t.Unixtime2gmt(t.Sse)

	return nil
}

// ConvertTimezone converts time from one timezone to another
func ConvertTimezone(t *Time, fromTz, toTz *TzInfo) error {
	if t == nil {
		return errors.New("time is nil")
	}

	// Convert to UTC first
	if fromTz != nil {
		offset, err := GetCurrentOffsetForTime(t)
		if err != nil {
			return err
		}
		t.Sse -= int64(offset)
	}

	// Convert to target timezone
	t.TzInfo = toTz
	if toTz != nil {
		t.ZoneType = TIMELIB_ZONETYPE_ID
		t.HaveZone = true

		// Apply new timezone offset
		return ApplyOffsetToTime(t)
	}

	return nil
}
