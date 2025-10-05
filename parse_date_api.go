package timelib

// ParseDateString parses a date/time string using the re2go-generated parser.
// This is the main entry point for date parsing.
func ParseDateString(str string, tzdb *TzDB, tzWrapper TzGetWrapper) (*Time, *ErrorContainer, error) {
	// For empty strings, create an empty time structure with an error
	// This matches C behavior where an empty string still returns a timelib_time*
	if len(str) == 0 {
		errContainer := &ErrorContainer{
			ErrorCount: 1,
			ErrorMessages: []ErrorMessage{
				{
					Position:  0,
					ErrorCode: TIMELIB_ERR_EMPTY_STRING,
					Character: 0,
					Message:   "Empty string",
				},
			},
		}
		emptyTime := &Time{
			Y:   TIMELIB_UNSET,
			M:   TIMELIB_UNSET,
			D:   TIMELIB_UNSET,
			H:   TIMELIB_UNSET,
			I:   TIMELIB_UNSET,
			S:   TIMELIB_UNSET,
			US:  TIMELIB_UNSET,
			Z:   TIMELIB_UNSET,
			Dst: TIMELIB_UNSET,
		}
		return emptyTime, errContainer, nil
	}

	// Initialize scanner with null-terminated string
	// Add null terminator for re2go
	strBytes := make([]byte, len(str)+1)
	copy(strBytes, str)
	strBytes[len(str)] = 0 // null terminator

	s := &Scanner{
		str:    strBytes,
		errors: &ErrorContainer{},
		time: &Time{
			Y:   TIMELIB_UNSET,
			M:   TIMELIB_UNSET,
			D:   TIMELIB_UNSET,
			H:   TIMELIB_UNSET,
			I:   TIMELIB_UNSET,
			S:   TIMELIB_UNSET,
			US:  TIMELIB_UNSET,
			Z:   TIMELIB_UNSET,
			Dst: TIMELIB_UNSET,
		},
		tzdb: tzdb,
	}

	// Set up pointers (lim points to the byte after the string, which is the null terminator)
	s.cur = &s.str[0]
	s.lim = &s.str[len(str)]

	// Run the scanner in a loop (like the C version)
	var t int
	maxIterations := len(str) * 2 // Safety limit
	iterations := 0
	for {
		t = scan(s, tzWrapper)
		iterations++
		if t == EOI || iterations >= maxIterations {
			break
		}
	}

	// Note: Do NOT initialize time to midnight here. The C version does this in
	// timelib_fill_holes(), not in the parser. Doing it here prevents TIMELIB_OVERRIDE_TIME
	// from working correctly. See timelib_fill_holes() in parse_date.c lines 25497-25502.

	// Initialize time fields to 0 if any time component was set (matching C behavior)
	// This is from parse_date.re lines 2615-2628
	if s.time.H != TIMELIB_UNSET || s.time.I != TIMELIB_UNSET || s.time.S != TIMELIB_UNSET || s.time.US != TIMELIB_UNSET {
		if s.time.H == TIMELIB_UNSET {
			s.time.H = 0
		}
		if s.time.I == TIMELIB_UNSET {
			s.time.I = 0
		}
		if s.time.S == TIMELIB_UNSET {
			s.time.S = 0
		}
		if s.time.US == TIMELIB_UNSET {
			s.time.US = 0
		}
	}

	return s.time, s.errors, nil
}

// StrToTime is a convenience function that parses a date/time string.
// It's similar to PHP's strtotime() function.
// Returns the parsed time and an error if parsing failed.
func StrToTime(str string, tzdb *TzDB) (*Time, error) {
	// Pass ParseTzfile as the timezone wrapper to enable timezone parsing from the string
	// This matches C behavior where timelib_strtotime passes tz_get_wrapper to scan
	time, errors, err := ParseDateString(str, tzdb, ParseTzfile)
	if err != nil {
		return nil, err
	}

	if errors != nil && errors.ErrorCount > 0 {
		for _, msg := range errors.ErrorMessages {
			if isFatalParseError(msg.ErrorCode) {
				return nil, &ParseError{
					Message:  msg.Message,
					Position: msg.Position,
				}
			}
		}
	}

	return time, nil
}

func isFatalParseError(errorCode int) bool {
	switch errorCode {
	case TIMELIB_ERR_DOUBLE_TZ,
		TIMELIB_ERR_TZID_NOT_FOUND,
		TIMELIB_ERR_DOUBLE_TIME,
		TIMELIB_ERR_DOUBLE_DATE,
		TIMELIB_ERR_EMPTY_STRING:
		return false
	default:
		return true
	}
}

// ParseError represents a parsing error
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}
