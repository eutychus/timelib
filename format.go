package timelib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseFromFormatWithMap parses a date/time string according to a format with format specifier mapping
func ParseFromFormatWithMap(format, input string, tzdb *TzDB, tzWrapper func(string, *TzDB, *int) (*TzInfo, error), formatConfig *FormatConfig) (*Time, *ErrorContainer) {
	errors := &ErrorContainer{
		ErrorMessages:   make([]ErrorMessage, 0),
		WarningMessages: make([]ErrorMessage, 0),
	}

	// Validate inputs
	if format == "" {
		errors.addError(TIMELIB_ERROR_EMPTY_STRING, "Format string is empty")
		return nil, errors
	}

	if input == "" {
		errors.addError(TIMELIB_ERROR_EMPTY_STRING, "Input string is empty")
		return nil, errors
	}

	// Use our custom format parser instead of the basic Go parser
	result, parseErrors := ParseFromFormatWithConfig(format, input, formatConfig)
	if parseErrors != nil && parseErrors.ErrorCount > 0 {
		// Copy errors from parseErrors to our errors container
		for _, errMsg := range parseErrors.ErrorMessages {
			errors.ErrorMessages = append(errors.ErrorMessages, errMsg)
		}
		errors.ErrorCount += parseErrors.ErrorCount
		return nil, errors
	}

	return result, errors
}

// parseWithFormatMap parses using format specifier mapping
func parseWithFormatMap(format, input string, result *Time, formatConfig *FormatConfig, errors *ErrorContainer) error {
	// This is a simplified implementation
	// Full implementation would handle all format specifiers

	formatRunes := []rune(format)
	inputRunes := []rune(input)
	formatLen := len(formatRunes)
	inputLen := len(inputRunes)

	formatPos := 0
	inputPos := 0

	for formatPos < formatLen && inputPos < inputLen {
		formatChar := formatRunes[formatPos]

		// Check if this is a format specifier
		if formatConfig.PrefixChar != 0 && formatChar == rune(formatConfig.PrefixChar) {
			if formatPos+1 < formatLen {
				formatPos++
				specifier := formatRunes[formatPos]

				// Find the format specifier code
				var specifierCode FormatSpecifierCode
				found := false
				for _, fs := range formatConfig.FormatMap {
					if fs.Specifier == byte(specifier) {
						specifierCode = fs.Code
						found = true
						break
					}
				}

				if !found {
					errors.addError(TIMELIB_ERR_INVALID_SPECIFIER, fmt.Sprintf("Invalid format specifier: %c", specifier))
					return fmt.Errorf("invalid format specifier")
				}

				// Parse according to specifier code
				consumed, err := parseFormatSpecifier(specifierCode, inputRunes[inputPos:], result, errors)
				if err != nil {
					return err
				}

				inputPos += consumed
				formatPos++
				continue
			}
		}

		// Handle literal characters
		if formatChar != inputRunes[inputPos] {
			if formatConfig.PrefixChar == 0 || formatChar != rune(formatConfig.PrefixChar) {
				errors.addError(TIMELIB_ERR_FORMAT_LITERAL_MISMATCH,
					fmt.Sprintf("Literal mismatch at position %d: expected '%c', got '%c'",
						inputPos, formatChar, inputRunes[inputPos]))
				return fmt.Errorf("format literal mismatch")
			}
		}

		formatPos++
		inputPos++
	}

	// Check for trailing data if not allowed
	if inputPos < inputLen && !formatConfig.AllowExtraCharacters {
		errors.addError(TIMELIB_ERR_TRAILING_DATA, "Trailing data found")
		return fmt.Errorf("trailing data found")
	}

	return nil
}

// parseFormatSpecifier parses a specific format specifier
func parseFormatSpecifier(code FormatSpecifierCode, input []rune, result *Time, errors *ErrorContainer) (int, error) {
	switch code {
	case TIMELIB_FORMAT_YEAR_FOUR_DIGIT:
		return parseYear4Digit(input, result, errors)
	case TIMELIB_FORMAT_YEAR_TWO_DIGIT:
		return parseYear2Digit(input, result, errors)
	case TIMELIB_FORMAT_MONTH_TWO_DIGIT:
		return parseMonth2Digit(input, result, errors)
	case TIMELIB_FORMAT_MONTH_TWO_DIGIT_PADDED:
		return parseMonth2DigitPadded(input, result, errors)
	case TIMELIB_FORMAT_DAY_TWO_DIGIT:
		return parseDay2Digit(input, result, errors)
	case TIMELIB_FORMAT_DAY_TWO_DIGIT_PADDED:
		return parseDay2DigitPadded(input, result, errors)
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX:
		return parseHour24(input, result, errors)
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX_PADDED:
		return parseHour24Padded(input, result, errors)
	case TIMELIB_FORMAT_MINUTE_TWO_DIGIT:
		return parseMinute2Digit(input, result, errors)
	case TIMELIB_FORMAT_SECOND_TWO_DIGIT:
		return parseSecond2Digit(input, result, errors)
	case TIMELIB_FORMAT_MICROSECOND_SIX_DIGIT:
		return parseMicrosecond6Digit(input, result, errors)
	case TIMELIB_FORMAT_TIMEZONE_OFFSET:
		return parseTimezoneOffset(input, result, errors)
	default:
		errors.addError(TIMELIB_ERR_INVALID_SPECIFIER, fmt.Sprintf("Unsupported format specifier code: %d", code))
		return 0, fmt.Errorf("unsupported format specifier")
	}
}

// Individual format specifier parsers
func parseYear4Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 4 {
		errors.addError(TIMELIB_ERR_NO_FOUR_DIGIT_YEAR, "Expected 4-digit year")
		return 0, fmt.Errorf("expected 4-digit year")
	}

	year, err := strconv.ParseInt(string(input[:4]), 10, 64)
	if err != nil {
		errors.addError(TIMELIB_ERR_NO_FOUR_DIGIT_YEAR, "Invalid 4-digit year")
		return 0, fmt.Errorf("invalid 4-digit year")
	}

	result.Y = year
	return 4, nil
}

func parseYear2Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 2 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_YEAR, "Expected 2-digit year")
		return 0, fmt.Errorf("expected 2-digit year")
	}

	year, err := strconv.ParseInt(string(input[:2]), 10, 64)
	if err != nil {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_YEAR, "Invalid 2-digit year")
		return 0, fmt.Errorf("invalid 2-digit year")
	}

	// Convert 2-digit year to 4-digit
	if year < 70 {
		year += 2000
	} else {
		year += 1900
	}

	result.Y = year
	return 2, nil
}

func parseMonth2Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 2 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_MONTH, "Expected 2-digit month")
		return 0, fmt.Errorf("expected 2-digit month")
	}

	month, err := strconv.ParseInt(string(input[:2]), 10, 64)
	if err != nil || month < 1 || month > 12 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_MONTH, "Invalid 2-digit month")
		return 0, fmt.Errorf("invalid 2-digit month")
	}

	result.M = month
	return 2, nil
}

func parseMonth2DigitPadded(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	return parseMonth2Digit(input, result, errors)
}

func parseDay2Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 2 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_DAY, "Expected 2-digit day")
		return 0, fmt.Errorf("expected 2-digit day")
	}

	day, err := strconv.ParseInt(string(input[:2]), 10, 64)
	if err != nil || day < 1 || day > 31 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_DAY, "Invalid 2-digit day")
		return 0, fmt.Errorf("invalid 2-digit day")
	}

	result.D = day
	return 2, nil
}

func parseDay2DigitPadded(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	return parseDay2Digit(input, result, errors)
}

func parseHour24(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 2 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_HOUR, "Expected 2-digit hour")
		return 0, fmt.Errorf("expected 2-digit hour")
	}

	hour, err := strconv.ParseInt(string(input[:2]), 10, 64)
	if err != nil || hour < 0 || hour > 23 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_HOUR, "Invalid 2-digit hour")
		return 0, fmt.Errorf("invalid 2-digit hour")
	}

	result.H = hour
	return 2, nil
}

func parseHour24Padded(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	return parseHour24(input, result, errors)
}

func parseMinute2Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 2 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_MINUTE, "Expected 2-digit minute")
		return 0, fmt.Errorf("expected 2-digit minute")
	}

	minute, err := strconv.ParseInt(string(input[:2]), 10, 64)
	if err != nil || minute < 0 || minute > 59 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_MINUTE, "Invalid 2-digit minute")
		return 0, fmt.Errorf("invalid 2-digit minute")
	}

	result.I = minute
	return 2, nil
}

func parseSecond2Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 2 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_SECOND, "Expected 2-digit second")
		return 0, fmt.Errorf("expected 2-digit second")
	}

	second, err := strconv.ParseInt(string(input[:2]), 10, 64)
	if err != nil || second < 0 || second > 59 {
		errors.addError(TIMELIB_ERR_NO_TWO_DIGIT_SECOND, "Invalid 2-digit second")
		return 0, fmt.Errorf("invalid 2-digit second")
	}

	result.S = second
	return 2, nil
}

func parseMicrosecond6Digit(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	if len(input) < 6 {
		errors.addError(TIMELIB_ERR_NO_SIX_DIGIT_MICROSECOND, "Expected 6-digit microsecond")
		return 0, fmt.Errorf("expected 6-digit microsecond")
	}

	microsecond, err := strconv.ParseInt(string(input[:6]), 10, 64)
	if err != nil || microsecond < 0 || microsecond > 999999 {
		errors.addError(TIMELIB_ERR_NO_SIX_DIGIT_MICROSECOND, "Invalid 6-digit microsecond")
		return 0, fmt.Errorf("invalid 6-digit microsecond")
	}

	result.US = microsecond
	return 6, nil
}

func parseTimezoneOffset(input []rune, result *Time, errors *ErrorContainer) (int, error) {
	// Pattern: +HHMM or -HHMM
	if len(input) < 5 {
		errors.addError(TIMELIB_ERR_INVALID_TZ_OFFSET, "Invalid timezone offset format")
		return 0, fmt.Errorf("invalid timezone offset format")
	}

	sign := input[0]
	if sign != '+' && sign != '-' {
		errors.addError(TIMELIB_ERR_INVALID_TZ_OFFSET, "Timezone offset must start with + or -")
		return 0, fmt.Errorf("timezone offset must start with + or -")
	}

	hours, err := strconv.ParseInt(string(input[1:3]), 10, 64)
	if err != nil {
		errors.addError(TIMELIB_ERR_INVALID_TZ_OFFSET, "Invalid timezone offset hours")
		return 0, fmt.Errorf("invalid timezone offset hours")
	}

	minutes, err := strconv.ParseInt(string(input[3:5]), 10, 64)
	if err != nil {
		errors.addError(TIMELIB_ERR_INVALID_TZ_OFFSET, "Invalid timezone offset minutes")
		return 0, fmt.Errorf("invalid timezone offset minutes")
	}

	// Calculate total offset in seconds
	totalOffset := hours*3600 + minutes*60
	if sign == '-' {
		totalOffset = -totalOffset
	}

	result.Z = int32(totalOffset)
	result.ZoneType = TIMELIB_ZONETYPE_OFFSET
	result.HaveZone = true

	return 5, nil
}

// parseWithBasicFormat provides basic format parsing as fallback
func parseWithBasicFormat(format, input string, result *Time, errors *ErrorContainer) error {
	// This is a simplified basic format parser
	// It handles common format strings like "Y-m-d H:i:s"

	format = strings.ReplaceAll(format, "Y", "2006")
	format = strings.ReplaceAll(format, "m", "01")
	format = strings.ReplaceAll(format, "d", "02")
	format = strings.ReplaceAll(format, "H", "15")
	format = strings.ReplaceAll(format, "i", "04")
	format = strings.ReplaceAll(format, "s", "05")

	// Try to parse with Go's time.Parse
	t, err := time.Parse(format, input)
	if err != nil {
		errors.addError(TIMELIB_ERROR_UNEXPECTED_DATA, fmt.Sprintf("Failed to parse with format: %v", err))
		return fmt.Errorf("failed to parse with format: %v", err)
	}

	// Convert to our Time structure
	result.Y = int64(t.Year())
	result.M = int64(t.Month())
	result.D = int64(t.Day())
	result.H = int64(t.Hour())
	result.I = int64(t.Minute())
	result.S = int64(t.Second())
	result.US = int64(t.Nanosecond() / 1000)

	result.HaveDate = true
	result.HaveTime = true

	return nil
}

// ParseFromFormatWithOptions parses with specific options
func ParseFromFormatWithOptions(format, input string, options ParseOptions) (*Time, *ErrorContainer) {
	formatConfig := &FormatConfig{
		FormatMap:            []FormatSpecifier{},
		AllowExtraCharacters: options.AllowExtraChars,
	}

	return ParseFromFormatWithMap(format, input, nil, nil, formatConfig)
}

// Helper function to add error to ErrorContainer
func (ec *ErrorContainer) addError(code int, message string) {
	if ec == nil {
		return
	}

	ec.ErrorCount++
	ec.ErrorMessages = append(ec.ErrorMessages, ErrorMessage{
		ErrorCode: code,
		Position:  0,
		Character: 0,
		Message:   message,
	})
}
