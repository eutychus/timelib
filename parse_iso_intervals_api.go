package timelib

import (
	"strings"
)

// ParseIsoInterval parses an ISO 8601 interval/duration string.
// Returns begin time, end time, period (duration), recurrences count, and any errors.
//
// ISO 8601 interval formats supported:
//   - R5/2008-03-01T13:00:00Z/P1Y2M10DT2H30M (recurrence/start/duration)
//   - R5/P1Y2M10DT2H30M/2008-03-01T13:00:00Z (recurrence/duration/end)
//   - P1Y2M10DT2H30M (duration only)
//   - 2008-03-01T13:00:00Z/2009-05-11T15:30:00Z (start/end)
func ParseIsoInterval(s string) (*Time, *Time, *RelTime, int, *ErrorContainer) {
	// Trim whitespace
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		errors := &ErrorContainer{
			ErrorCount: 1,
			ErrorMessages: []ErrorMessage{
				{Position: 0, Character: 0, Message: "Empty string"},
			},
		}
		return nil, nil, nil, 0, errors
	}

	// Initialize scanner with null-terminated string
	strBytes := make([]byte, len(s)+1)
	copy(strBytes, s)
	strBytes[len(s)] = 0 // null terminator

	scanner := &IsoIntervalScanner{
		str:    strBytes,
		errors: &ErrorContainer{},
		Begin: &Time{
			Y:            TIMELIB_UNSET,
			M:            TIMELIB_UNSET,
			D:            TIMELIB_UNSET,
			H:            TIMELIB_UNSET,
			I:            TIMELIB_UNSET,
			S:            TIMELIB_UNSET,
			US:           0,
			Z:            0,
			Dst:          0,
			IsLocaltime:  false,
			ZoneType:     TIMELIB_ZONETYPE_OFFSET,
		},
		End: &Time{
			Y:            TIMELIB_UNSET,
			M:            TIMELIB_UNSET,
			D:            TIMELIB_UNSET,
			H:            TIMELIB_UNSET,
			I:            TIMELIB_UNSET,
			S:            TIMELIB_UNSET,
			US:           0,
			Z:            0,
			Dst:          0,
			IsLocaltime:  false,
			ZoneType:     TIMELIB_ZONETYPE_OFFSET,
		},
		Period: &RelTime{
			Y:                   0,
			M:                   0,
			D:                   0,
			H:                   0,
			I:                   0,
			S:                   0,
			US:                  0,
			Weekday:             0,
			WeekdayBehavior:     0,
			FirstLastDayOf:      0,
			Invert:              false,
			Days:                TIMELIB_UNSET,
			Special:             struct{ Type int; Amount int64 }{},
			HaveWeekdayRelative: false,
			HaveSpecialRelative: false,
		},
		Recurrences: 1,
	}

	// Set up pointers
	scanner.cur = &scanner.str[0]
	scanner.lim = &scanner.str[len(s)]

	// Run the scanner in a loop
	maxIterations := len(s) * 2 // Safety limit
	iterations := 0
	for {
		t := scanIsoInterval(scanner)
		iterations++
		if t == EOI || iterations >= maxIterations {
			break
		}
	}

	// Prepare return values
	var begin *Time
	var end *Time
	var period *RelTime
	recurrences := 0

	if scanner.HaveBeginDate {
		begin = scanner.Begin
	}
	if scanner.HaveEndDate {
		end = scanner.End
	}
	if scanner.HavePeriod {
		period = scanner.Period
	}
	if scanner.HaveRecurrences {
		recurrences = scanner.Recurrences
	}

	return begin, end, period, recurrences, scanner.errors
}

// StrToInterval is a convenience function that parses an ISO 8601 interval string
// and returns the parsed components.
func StrToInterval(s string) (*Time, *Time, *RelTime, int, error) {
	begin, end, period, recurrences, errors := ParseIsoInterval(s)

	// Check for parsing errors
	if errors != nil && errors.ErrorCount > 0 {
		if len(errors.ErrorMessages) > 0 {
			return nil, nil, nil, 0, &ParseError{
				Message:  errors.ErrorMessages[0].Message,
				Position: errors.ErrorMessages[0].Position,
			}
		}
	}

	return begin, end, period, recurrences, nil
}
