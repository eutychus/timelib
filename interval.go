package timelib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Strtointerval parses a subset of an ISO 8601 intervals specification string into its constituent parts
func Strtointerval(s string, errors *ErrorContainer) (*Time, *Time, *RelTime, int, error) {
	if s == "" {
		if errors != nil {
			errors.addError(TIMELIB_ERROR_EMPTY_STRING, "Empty interval string")
		}
		return nil, nil, nil, 0, fmt.Errorf("empty interval string")
	}

	// Trim whitespace
	s = strings.TrimSpace(s)

	// Try to match different ISO 8601 interval patterns

	// Pattern 1: Recurring interval (R5/2007-03-01T13:00:00Z/P1Y2M3DT4H5M6S)
	// This must be checked first, before general slash handling
	if strings.HasPrefix(s, "R") {
		// Extract recurrence count
		recurrenceRegex := regexp.MustCompile(`^R(\d+)/(.+)$`)
		matches := recurrenceRegex.FindStringSubmatch(s)

		if len(matches) == 3 {
			recurrences, err := strconv.Atoi(matches[1])
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid recurrence count")
				}
				return nil, nil, nil, 0, fmt.Errorf("invalid recurrence count")
			}

			// Parse the remaining interval specification (which should be a start datetime + duration)
			intervalPart := matches[2]

			// Check if this is a start datetime + duration format
			if strings.Count(intervalPart, "/") == 1 {
				// This should be start datetime + duration
				parts := strings.Split(intervalPart, "/")
				if !strings.HasPrefix(parts[0], "P") && strings.HasPrefix(parts[1], "P") {
					// Parse start datetime
					begin, err := StrToTime(parts[0], nil)
					var parseErrors *ErrorContainer
					if err != nil {
						parseErrors = &ErrorContainer{
							ErrorCount: 1,
							ErrorMessages: []ErrorMessage{{Message: err.Error()}},
						}
					}
					if parseErrors != nil && parseErrors.ErrorCount > 0 {
						if errors != nil {
							errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
							errors.ErrorCount += parseErrors.ErrorCount
						}
						return nil, nil, nil, 0, fmt.Errorf("failed to parse start time")
					}

					// Parse duration
					period, err := parseISODuration(parts[1], errors)
					if err != nil {
						return nil, nil, nil, 0, err
					}

					// Calculate end time from begin + period
					endCopy := *begin
					end := &endCopy
					end = end.Add(period)

					return begin, end, period, recurrences, nil
				}
			}

			// If we get here, the format is not supported
			if errors != nil {
				errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid recurring interval format")
			}
			return nil, nil, nil, 0, fmt.Errorf("invalid recurring interval format")
		}

		// Handle infinite recurrence (R/...)
		infiniteRegex := regexp.MustCompile(`^R/(.+)$`)
		matches = infiniteRegex.FindStringSubmatch(s)
		if len(matches) == 2 {
			// Parse the remaining interval specification as infinite recurrence
			intervalPart := matches[1]

			// Check if this is a start datetime + duration format
			if strings.Count(intervalPart, "/") == 1 {
				// This should be start datetime + duration
				parts := strings.Split(intervalPart, "/")
				if !strings.HasPrefix(parts[0], "P") && strings.HasPrefix(parts[1], "P") {
					// Parse start datetime
					begin, err := StrToTime(parts[0], nil)
					var parseErrors *ErrorContainer
					if err != nil {
						parseErrors = &ErrorContainer{
							ErrorCount: 1,
							ErrorMessages: []ErrorMessage{{Message: err.Error()}},
						}
					}
					if parseErrors != nil && parseErrors.ErrorCount > 0 {
						if errors != nil {
							errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
							errors.ErrorCount += parseErrors.ErrorCount
						}
						return nil, nil, nil, 0, fmt.Errorf("failed to parse start time")
					}

					// Parse duration
					period, err := parseISODuration(parts[1], errors)
					if err != nil {
						return nil, nil, nil, 0, err
					}

					// Calculate end time from begin + period
					endCopy := *begin
					end := &endCopy
					end = end.Add(period)

					return begin, end, period, 0, nil // 0 recurrences means infinite
				}
			}

			// If we get here, the format is not supported
			if errors != nil {
				errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid infinite recurring interval format")
			}
			return nil, nil, nil, 0, fmt.Errorf("invalid infinite recurring interval format")
		}
	}

	// Pattern 2: Duration only (P1Y2M3DT4H5M6S)
	if strings.HasPrefix(s, "P") && !strings.Contains(s, "/") {
		period, err := parseISODuration(s, errors)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		return nil, nil, period, 0, nil
	}

	// Pattern 3: Check for slash-separated intervals
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) != 2 {
			if errors != nil {
				errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid interval format: expected exactly one '/' separator")
			}
			return nil, nil, nil, 0, fmt.Errorf("invalid interval format")
		}

		// Pattern 3a: Start datetime and duration (2007-03-01T13:00:00Z/P1Y2M3DT4H5M6S)
		if !strings.HasPrefix(parts[0], "P") && strings.HasPrefix(parts[1], "P") {
			// Start datetime + duration
			begin, err := StrToTime(parts[0], nil)
			var parseErrors *ErrorContainer
			if err != nil {
				parseErrors = &ErrorContainer{
					ErrorCount: 1,
					ErrorMessages: []ErrorMessage{{Message: err.Error()}},
				}
			}
			if parseErrors != nil && parseErrors.ErrorCount > 0 {
				if errors != nil {
					errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
					errors.ErrorCount += parseErrors.ErrorCount
				}
				return nil, nil, nil, 0, fmt.Errorf("failed to parse start time")
			}

			period, err := parseISODuration(parts[1], errors)
			if err != nil {
				return nil, nil, nil, 0, err
			}

			// Calculate end time from begin + period
			endCopy := *begin
			end := &endCopy
			end = end.Add(period)

			return begin, end, period, 0, nil
		}

		// Pattern 3b: Duration and end datetime (P1Y2M3DT4H5M6S/2008-05-11T15:30:00Z)
		if strings.HasPrefix(parts[0], "P") && !strings.HasPrefix(parts[1], "P") {
			// Duration + end datetime
			period, err := parseISODuration(parts[0], errors)
			if err != nil {
				return nil, nil, nil, 0, err
			}

			end, err := StrToTime(parts[1], nil)
			var parseErrors *ErrorContainer
			if err != nil {
				parseErrors = &ErrorContainer{
					ErrorCount: 1,
					ErrorMessages: []ErrorMessage{{Message: err.Error()}},
				}
			}
			if parseErrors != nil && parseErrors.ErrorCount > 0 {
				if errors != nil {
					errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
					errors.ErrorCount += parseErrors.ErrorCount
				}
				return nil, nil, nil, 0, fmt.Errorf("failed to parse end time")
			}

			// Calculate begin time from end - period
			beginCopy := *end
			begin := &beginCopy
			begin = begin.Sub(period)

			return begin, end, period, 0, nil
		}

		// Pattern 3c: Start and end datetime (2007-03-01T13:00:00Z/2008-05-11T15:30:00Z)
		if !strings.HasPrefix(parts[0], "P") && !strings.HasPrefix(parts[1], "P") {
			// Parse start time
			begin, err := StrToTime(parts[0], nil)
			if err != nil {
				if errors != nil {
					errors.ErrorCount++
					errors.ErrorMessages = append(errors.ErrorMessages, ErrorMessage{Message: err.Error()})
				}
				return nil, nil, nil, 0, fmt.Errorf("failed to parse start time: %w", err)
			}

			// Parse end time
			end, err := StrToTime(parts[1], nil)
			if err != nil {
				if errors != nil {
					errors.ErrorCount++
					errors.ErrorMessages = append(errors.ErrorMessages, ErrorMessage{Message: err.Error()})
				}
				return nil, nil, nil, 0, fmt.Errorf("failed to parse end time: %w", err)
			}

			return begin, end, nil, 0, nil
		}
	}

	// Pattern 4: Unrecognized format
	if strings.HasPrefix(s, "R") {
		// Extract recurrence count
		recurrenceRegex := regexp.MustCompile(`^R(\d+)/(.+)$`)
		matches := recurrenceRegex.FindStringSubmatch(s)

		if len(matches) == 3 {
			recurrences, err := strconv.Atoi(matches[1])
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid recurrence count")
				}
				return nil, nil, nil, 0, fmt.Errorf("invalid recurrence count")
			}

			// Parse the remaining interval specification (which should be a start datetime + duration)
			// We need to handle this carefully to avoid the slash-counting issue
			intervalPart := matches[2]

			// Check if this is a start datetime + duration format
			if strings.Count(intervalPart, "/") == 1 {
				// This should be start datetime + duration
				parts := strings.Split(intervalPart, "/")
				if !strings.HasPrefix(parts[0], "P") && strings.HasPrefix(parts[1], "P") {
					// Parse start datetime
					begin, err := StrToTime(parts[0], nil)
					var parseErrors *ErrorContainer
					if err != nil {
						parseErrors = &ErrorContainer{
							ErrorCount: 1,
							ErrorMessages: []ErrorMessage{{Message: err.Error()}},
						}
					}
					if parseErrors != nil && parseErrors.ErrorCount > 0 {
						if errors != nil {
							errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
							errors.ErrorCount += parseErrors.ErrorCount
						}
						return nil, nil, nil, 0, fmt.Errorf("failed to parse start time")
					}

					// Parse duration
					period, err := parseISODuration(parts[1], errors)
					if err != nil {
						return nil, nil, nil, 0, err
					}

					// Calculate end time from begin + period
					endCopy := *begin
					end := &endCopy
					end = end.Add(period)

					return begin, end, period, recurrences, nil
				}
			}

			// If we get here, the format is not supported
			if errors != nil {
				errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid recurring interval format")
			}
			return nil, nil, nil, 0, fmt.Errorf("invalid recurring interval format")
		}

		// Handle infinite recurrence (R/...)
		infiniteRegex := regexp.MustCompile(`^R/(.+)$`)
		matches = infiniteRegex.FindStringSubmatch(s)
		if len(matches) == 2 {
			// Parse the remaining interval specification as infinite recurrence
			intervalPart := matches[1]

			// Check if this is a start datetime + duration format
			if strings.Count(intervalPart, "/") == 1 {
				// This should be start datetime + duration
				parts := strings.Split(intervalPart, "/")
				if !strings.HasPrefix(parts[0], "P") && strings.HasPrefix(parts[1], "P") {
					// Parse start datetime
					begin, err := StrToTime(parts[0], nil)
					var parseErrors *ErrorContainer
					if err != nil {
						parseErrors = &ErrorContainer{
							ErrorCount: 1,
							ErrorMessages: []ErrorMessage{{Message: err.Error()}},
						}
					}
					if parseErrors != nil && parseErrors.ErrorCount > 0 {
						if errors != nil {
							errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
							errors.ErrorCount += parseErrors.ErrorCount
						}
						return nil, nil, nil, 0, fmt.Errorf("failed to parse start time")
					}

					// Parse duration
					period, err := parseISODuration(parts[1], errors)
					if err != nil {
						return nil, nil, nil, 0, err
					}

					// Calculate end time from begin + period
					endCopy := *begin
					end := &endCopy
					end = end.Add(period)

					return begin, end, period, 0, nil // 0 recurrences means infinite
				}
			}

			// If we get here, the format is not supported
			if errors != nil {
				errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid infinite recurring interval format")
			}
			return nil, nil, nil, 0, fmt.Errorf("invalid infinite recurring interval format")
		}
	}

	if errors != nil {
		errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Unrecognized interval format")
	}
	return nil, nil, nil, 0, fmt.Errorf("unrecognized interval format")
}

// parseISODuration parses an ISO 8601 duration string (P1Y2M3DT4H5M6S)
func parseISODuration(duration string, errors *ErrorContainer) (*RelTime, error) {
	if !strings.HasPrefix(duration, "P") {
		if errors != nil {
			errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Duration must start with 'P'")
		}
		return nil, fmt.Errorf("duration must start with 'P'")
	}

	relTime := RelTimeCtor()

	// Remove the 'P' prefix
	duration = duration[1:]

	// Check for time component (starts with 'T')
	timeComponent := ""
	if strings.Contains(duration, "T") {
		parts := strings.Split(duration, "T")
		if len(parts) != 2 {
			if errors != nil {
				errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid duration format")
			}
			return nil, fmt.Errorf("invalid duration format")
		}
		duration = parts[0]
		timeComponent = parts[1]
	}

	// Parse date components (Y, M, W, D)
	if duration != "" {
		// Year
		yearRegex := regexp.MustCompile(`(\d+)Y`)
		if matches := yearRegex.FindStringSubmatch(duration); len(matches) == 2 {
			years, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid year value")
				}
				return nil, fmt.Errorf("invalid year value")
			}
			relTime.Y = years
		}

		// Month
		monthRegex := regexp.MustCompile(`(\d+)M`)
		if matches := monthRegex.FindStringSubmatch(duration); len(matches) == 2 {
			months, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid month value")
				}
				return nil, fmt.Errorf("invalid month value")
			}
			relTime.M = months
		}

		// Week
		weekRegex := regexp.MustCompile(`(\d+)W`)
		if matches := weekRegex.FindStringSubmatch(duration); len(matches) == 2 {
			weeks, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid week value")
				}
				return nil, fmt.Errorf("invalid week value")
			}
			relTime.D += weeks * 7 // Convert weeks to days
		}

		// Day
		dayRegex := regexp.MustCompile(`(\d+)D`)
		if matches := dayRegex.FindStringSubmatch(duration); len(matches) == 2 {
			days, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid day value")
				}
				return nil, fmt.Errorf("invalid day value")
			}
			relTime.D += days
		}
	}

	// Parse time components (H, M, S)
	if timeComponent != "" {
		// Hour
		hourRegex := regexp.MustCompile(`(\d+)H`)
		if matches := hourRegex.FindStringSubmatch(timeComponent); len(matches) == 2 {
			hours, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid hour value")
				}
				return nil, fmt.Errorf("invalid hour value")
			}
			relTime.H = hours
		}

		// Minute
		minuteRegex := regexp.MustCompile(`(\d+)M`)
		if matches := minuteRegex.FindStringSubmatch(timeComponent); len(matches) == 2 {
			minutes, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid minute value")
				}
				return nil, fmt.Errorf("invalid minute value")
			}
			relTime.I = minutes
		}

		// Second
		secondRegex := regexp.MustCompile(`(\d+)S`)
		if matches := secondRegex.FindStringSubmatch(timeComponent); len(matches) == 2 {
			seconds, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				if errors != nil {
					errors.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid second value")
				}
				return nil, fmt.Errorf("invalid second value")
			}
			relTime.S = seconds
		}
	}

	return relTime, nil
}

// StrtointervalWithOptions parses an interval with options
func StrtointervalWithOptions(s string, options ParseOptions) (*Time, *Time, *RelTime, int, *ErrorContainer) {
	errors := &ErrorContainer{
		ErrorMessages:   make([]ErrorMessage, 0),
		WarningMessages: make([]ErrorMessage, 0),
	}

	begin, end, period, recurrences, err := Strtointerval(s, errors)
	if err != nil {
		return begin, end, period, recurrences, errors
	}

	return begin, end, period, recurrences, errors
}
