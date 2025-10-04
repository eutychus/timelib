package timelib

// ParseDateString parses a date/time string using the re2go-generated parser.
// This is the main entry point for date parsing.
func ParseDateString(str string, tzdb *TzDB, tzWrapper TzGetWrapper) (*Time, *ErrorContainer, error) {
	if len(str) == 0 {
		return nil, nil, nil
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
			Y:  TIMELIB_UNSET,
			M:  TIMELIB_UNSET,
			D:  TIMELIB_UNSET,
			H:  TIMELIB_UNSET,
			I:  TIMELIB_UNSET,
			S:  TIMELIB_UNSET,
			US: TIMELIB_UNSET,
			Z:  TIMELIB_UNSET,
			Dst: TIMELIB_UNSET,
		},
		tzdb:   tzdb,
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

	// Initialize time to midnight if date is set but time is not (matching C FillHoles behavior)
	// This is from parse_date.re lines 2675-2679
	if s.time.HaveDate && !s.time.HaveTime {
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
	}

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
func StrToTime(str string, tzdb *TzDB) (*Time, error) {
	time, errors, err := ParseDateString(str, tzdb, nil)
	if err != nil {
		return nil, err
	}

	// Check for parsing errors
	if errors != nil && errors.ErrorCount > 0 {
		// Return first error as a simple error
		if len(errors.ErrorMessages) > 0 {
			return nil, &ParseError{
				Message:  errors.ErrorMessages[0].Message,
				Position: errors.ErrorMessages[0].Position,
			}
		}
	}

	return time, nil
}

// ParseError represents a parsing error
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}
