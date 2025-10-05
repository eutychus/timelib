package timelib

import (
	"time"
)

// Strtotime parses a date/time string and returns a Unix timestamp,
// matching PHP's strtotime() function signature and behavior.
//
// PHP signature: int|false strtotime(string $datetime, ?int $baseTimestamp = null)
//
// Parameters:
//   - datetime: The date/time string to parse (e.g., "2024-01-15", "tomorrow", "+1 day")
//   - baseTimestamp: Optional base Unix timestamp for relative times. Use 0 or time.Now().Unix() for current time.
//     If 0 is provided and the input is relative (like "tomorrow"), current time is used.
//
// Returns:
//   - Unix timestamp (seconds since epoch) on success
//   - -1 on parse error (equivalent to PHP's false)
//
// Examples:
//
//	ts := Strtotime("2024-01-15", 0)                    // Absolute date
//	ts := Strtotime("tomorrow", time.Now().Unix())      // Relative to now
//	ts := Strtotime("+1 day", 1704067200)               // Relative to base
//	if ts == -1 {
//	    // Parse error
//	}
func Strtotime(datetime string, baseTimestamp int64) int64 {
	// Parse the input string
	parsedTime, err := StrToTime(datetime, nil)
	if err != nil {
		return -1
	}

	// If we have relative components, we need a base time
	if parsedTime.HaveRelative {
		var baseTime *Time

		// If no base timestamp provided (0), use current time for relative parsing
		if baseTimestamp == 0 {
			baseTime = TimeCtor()
			now := time.Now()
			baseTime.Unixtime2gmt(now.Unix())
		} else {
			// Create a time structure from the base timestamp
			baseTime = TimeCtor()
			baseTime.Unixtime2gmt(baseTimestamp)
		}

		// Apply the relative components to the base time
		result := baseTime.Add(&parsedTime.Relative)
		parsedTime = result
	}

	// Make sure timestamp is calculated
	if !parsedTime.SseUptodate {
		if parsedTime.HaveTime || parsedTime.HaveDate {
			parsedTime.UpdateTS(nil)
		}
	} else if parsedTime.HaveRelative {
		parsedTime.UpdateTS(nil)
	}

	return parsedTime.Sse
}

// StrtotimeToGoTime parses a date/time string and returns a Go time.Time value,
// similar to PHP's strtotime() but returning Go's native time type.
//
// This combines PHP's strtotime() parsing with Go's time.Time for easy integration
// with Go's standard library time functions.
//
// Parameters:
//   - datetime: The date/time string to parse (e.g., "2024-01-15", "tomorrow", "+1 day")
//   - baseTime: Optional base time for relative times. Use nil for current time or time.Time{} for epoch.
//     If nil and the input is relative (like "tomorrow"), current time is used.
//
// Returns:
//   - time.Time: The parsed time in UTC
//   - error: Parse error if the string is invalid
//
// Examples:
//
//	t, err := StrtotimeToGoTime("2024-01-15 10:30:00", nil)
//	t, err := StrtotimeToGoTime("tomorrow", nil)
//	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	t, err := StrtotimeToGoTime("+1 week", &baseTime)
//	if err != nil {
//	    // Handle parse error
//	}
//	fmt.Printf("Parsed time: %s\n", t.Format(time.RFC3339))
func StrtotimeToGoTime(datetime string, baseTime *time.Time) (time.Time, error) {
	// Parse the input string
	parsedTime, err := StrToTime(datetime, nil)
	if err != nil {
		return time.Time{}, err
	}

	// If we have relative components, we need a base time
	if parsedTime.HaveRelative {
		var timelibBase *Time

		// Determine base time for relative parsing
		if baseTime == nil {
			// Use current time
			timelibBase = TimeCtor()
			now := time.Now()
			timelibBase.Unixtime2gmt(now.Unix())
		} else if baseTime.IsZero() {
			// Use epoch (0)
			timelibBase = TimeCtor()
			timelibBase.Unixtime2gmt(0)
		} else {
			// Use provided base time
			timelibBase = TimeCtor()
			timelibBase.Unixtime2gmt(baseTime.Unix())
		}

		// Apply the relative components to the base time
		result := timelibBase.Add(&parsedTime.Relative)
		parsedTime = result
	}

	// Make sure timestamp is calculated
	if !parsedTime.SseUptodate {
		if parsedTime.HaveTime || parsedTime.HaveDate {
			parsedTime.UpdateTS(nil)
		}
	} else if parsedTime.HaveRelative {
		parsedTime.UpdateTS(nil)
	}

	// Convert to Go time.Time
	// Note: Go's time.Unix() takes seconds and nanoseconds
	// Only include microseconds if they were explicitly set in the parsed time
	var nanoseconds int64
	if parsedTime.HaveTime {
		// Convert US (microseconds) to nanoseconds by multiplying by 1000
		nanoseconds = parsedTime.US * 1000
	}

	return time.Unix(parsedTime.Sse, nanoseconds).UTC(), nil
}

// StrtotimeWithTimezone is like StrtotimeToGoTime but allows specifying a timezone
// for the result. The parsing still uses UTC/specified timezone in the string,
// but the returned time.Time will be in the requested timezone.
//
// Parameters:
//   - datetime: The date/time string to parse
//   - baseTime: Optional base time for relative times
//   - location: The timezone for the result (e.g., time.UTC, time.Local, or time.LoadLocation result)
//
// Returns:
//   - time.Time: The parsed time in the specified timezone
//   - error: Parse error if the string is invalid
//
// Example:
//
//	loc, _ := time.LoadLocation("America/New_York")
//	t, err := StrtotimeWithTimezone("2024-01-15 10:30:00", nil, loc)
func StrtotimeWithTimezone(datetime string, baseTime *time.Time, location *time.Location) (time.Time, error) {
	t, err := StrtotimeToGoTime(datetime, baseTime)
	if err != nil {
		return time.Time{}, err
	}

	// Convert to the requested timezone
	return t.In(location), nil
}
