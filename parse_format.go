package timelib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatParser handles format-based parsing
type FormatParser struct {
	input     string
	format    string
	position  int
	formatPos int
	time      *Time
	errors    *ErrorContainer
	config    *FormatConfig
}

// NewFormatParser creates a new format parser
func NewFormatParser(input, format string, config *FormatConfig) *FormatParser {
	if config == nil {
		config = &FormatConfig{
			FormatMap:  []FormatSpecifier{},
			PrefixChar: 0,
		}
	}
	return &FormatParser{
		input:  input,
		format: format,
		time:   TimeCtor(),
		errors: &ErrorContainer{},
		config: config,
	}
}

// ParseFromFormatWithConfig parses input according to format with custom config
func ParseFromFormatWithConfig(format, input string, config *FormatConfig) (*Time, *ErrorContainer) {
	parser := NewFormatParser(input, format, config)
	result := parser.Parse()
	return result, parser.errors
}

// ParseFromFormat parses input according to format using default format specifiers
func ParseFromFormat(format, input string) (*Time, *ErrorContainer) {
	// Use default format configuration with standard format specifiers
	defaultConfig := &FormatConfig{
		FormatMap: []FormatSpecifier{
			{'Y', TIMELIB_FORMAT_YEAR_FOUR_DIGIT},
			{'y', TIMELIB_FORMAT_YEAR_TWO_DIGIT},
			{'m', TIMELIB_FORMAT_MONTH_TWO_DIGIT_PADDED},
			{'n', TIMELIB_FORMAT_MONTH_TWO_DIGIT},
			{'d', TIMELIB_FORMAT_DAY_TWO_DIGIT_PADDED},
			{'j', TIMELIB_FORMAT_DAY_TWO_DIGIT},
			{'H', TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX},
			{'h', TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX_PADDED},
			{'g', TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX},
			{'G', TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX_PADDED},
			{'i', TIMELIB_FORMAT_MINUTE_TWO_DIGIT},
			{'s', TIMELIB_FORMAT_SECOND_TWO_DIGIT},
			{'u', TIMELIB_FORMAT_MICROSECOND_SIX_DIGIT},
			{'e', TIMELIB_FORMAT_TIMEZONE_OFFSET},
			{'P', TIMELIB_FORMAT_TIMEZONE_OFFSET},
			{'p', TIMELIB_FORMAT_TIMEZONE_OFFSET},
			{'T', TIMELIB_FORMAT_TIMEZONE_OFFSET},
			{'O', TIMELIB_FORMAT_TIMEZONE_OFFSET},
			{'Z', TIMELIB_FORMAT_TIMEZONE_OFFSET_MINUTES},
			{'F', TIMELIB_FORMAT_TEXTUAL_MONTH_FULL},
			{'M', TIMELIB_FORMAT_TEXTUAL_MONTH_3_LETTER},
			{'D', TIMELIB_FORMAT_TEXTUAL_DAY_3_LETTER},
			{'l', TIMELIB_FORMAT_TEXTUAL_DAY_FULL},
			{'a', TIMELIB_FORMAT_MERIDIAN},
			{'A', TIMELIB_FORMAT_MERIDIAN},
			{'z', TIMELIB_FORMAT_DAY_OF_YEAR},
			{'U', TIMELIB_FORMAT_EPOCH_SECONDS},
			{' ', TIMELIB_FORMAT_WHITESPACE},
			{'\\', TIMELIB_FORMAT_ESCAPE},
			{'+', TIMELIB_FORMAT_ALLOW_EXTRA_CHARACTERS},
			{'#', TIMELIB_FORMAT_ANY_SEPARATOR},
			{'?', TIMELIB_FORMAT_RANDOM_CHAR},
			{'!', TIMELIB_FORMAT_RESET_ALL},
			{'|', TIMELIB_FORMAT_RESET_ALL_WHEN_NOT_SET},
			{'*', TIMELIB_FORMAT_SKIP_TO_SEPARATOR},
			{'B', TIMELIB_FORMAT_YEAR_ISO},
			{'b', TIMELIB_FORMAT_DAY_OF_WEEK_ISO},
			{'V', TIMELIB_FORMAT_WEEK_OF_YEAR_ISO},
			{'v', TIMELIB_FORMAT_MILLISECOND_THREE_DIGIT},
			{'S', TIMELIB_FORMAT_DAY_SUFFIX},
		},
		PrefixChar: 0,
	}
	return ParseFromFormatWithConfig(format, input, defaultConfig)
}

// Parse performs the format parsing
func (p *FormatParser) Parse() *Time {
	p.position = 0
	p.formatPos = 0

	for p.formatPos < len(p.format) {
		formatChar := rune(p.format[p.formatPos])

		// Handle escape character
		if formatChar == '\\' {
			if p.formatPos+1 < len(p.format) {
				p.formatPos++ // Move to the character after backslash
				expectedChar := rune(p.format[p.formatPos])
				// Check if we have enough input left
				if p.position >= len(p.input) {
					p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, fmt.Sprintf("Expected character '%c' after escape, but input ended", expectedChar))
					return p.time
				}
				if !p.matchCharacter(expectedChar) {
					p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, fmt.Sprintf("Expected character '%c' after escape, got '%c'", expectedChar, p.input[p.position]))
					return p.time
				}
				p.position++  // Advance input position past the matched character
				p.formatPos++ // Move past the escaped character
				continue
			} else {
				// Backslash at end of format string - treat as literal
				if p.position >= len(p.input) || p.input[p.position] != '\\' {
					p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected backslash character")
					return p.time
				}
				p.position++
				p.formatPos++
				continue
			}
		}

		// Find format specifier
		spec := p.findFormatSpecifier(formatChar)
		if spec == nil {
			// Literal character - must match exactly
			if !p.matchCharacter(formatChar) {
				p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Literal character mismatch")
				return p.time
			}
			p.formatPos++
			p.position++
			continue
		}

		// Handle format specifier
		if !p.parseFormatSpecifier(spec) {
			return p.time
		}
		p.formatPos++
	}

	// Convert ISO week dates to calendar dates if we have ISO week information
	p.convertISOWeekToDate()

	// If we parsed time components but no date components, set date components to 0
	if p.time.HaveTime && !p.time.HaveDate {
		if p.time.Y == TIMELIB_UNSET {
			p.time.Y = 0
		}
		if p.time.M == TIMELIB_UNSET {
			p.time.M = 0
		}
		if p.time.D == TIMELIB_UNSET {
			p.time.D = 0
		}
	}

	// If we parsed any time components, set unset time components to 0
	if p.time.HaveTime {
		if p.time.H == TIMELIB_UNSET {
			p.time.H = 0
		}
		if p.time.I == TIMELIB_UNSET {
			p.time.I = 0
		}
		if p.time.S == TIMELIB_UNSET {
			p.time.S = 0
		}
		if p.time.US == 0 { // US is initialized to 0, not TIMELIB_UNSET
			p.time.US = 0
		}
	}

	return p.time
}

// findFormatSpecifier finds a format specifier by character
func (p *FormatParser) findFormatSpecifier(char rune) *FormatSpecifier {
	for _, spec := range p.config.FormatMap {
		if rune(spec.Specifier) == char {
			return &spec
		}
	}
	return nil
}

// isSeparatorChar checks if a character is a separator
func isSeparatorChar(char rune) bool {
	separators := []rune{':', '.', '-', '/', ',', ' ', '\t'}
	for _, sep := range separators {
		if char == sep {
			return true
		}
	}
	return false
}

// matchCharacter matches a literal character with separator flexibility
func (p *FormatParser) matchCharacter(char rune) bool {
	if p.position >= len(p.input) {
		return false
	}

	inputChar := rune(p.input[p.position])

	// If both characters are separators, consider it a match
	if isSeparatorChar(char) && isSeparatorChar(inputChar) {
		return true
	}

	// Otherwise, require exact match
	return inputChar == char
}

// parseFormatSpecifier handles parsing based on format specifier type
func (p *FormatParser) parseFormatSpecifier(spec *FormatSpecifier) bool {
	switch spec.Code {
	case TIMELIB_FORMAT_YEAR_FOUR_DIGIT:
		return p.parseYearFourDigit()
	case TIMELIB_FORMAT_YEAR_TWO_DIGIT:
		return p.parseYearTwoDigit()
	case TIMELIB_FORMAT_YEAR_ISO:
		return p.parseYearISO()
	case TIMELIB_FORMAT_MONTH_TWO_DIGIT_PADDED:
		return p.parseMonthTwoDigitPadded()
	case TIMELIB_FORMAT_MONTH_TWO_DIGIT:
		return p.parseMonthTwoDigit()
	case TIMELIB_FORMAT_DAY_TWO_DIGIT_PADDED:
		return p.parseDayTwoDigitPadded()
	case TIMELIB_FORMAT_DAY_TWO_DIGIT:
		return p.parseDayTwoDigit()
	case TIMELIB_FORMAT_DAY_OF_WEEK_ISO:
		return p.parseDayOfWeekISO()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX:
		return p.parseHour24()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX_PADDED:
		return p.parseHour24()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX:
		return p.parseHour12()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX_PADDED:
		return p.parseHour12()
	case TIMELIB_FORMAT_MINUTE_TWO_DIGIT:
		return p.parseMinute()
	case TIMELIB_FORMAT_SECOND_TWO_DIGIT:
		return p.parseSecond()
	case TIMELIB_FORMAT_MICROSECOND_SIX_DIGIT:
		return p.parseMicrosecond()
	case TIMELIB_FORMAT_MILLISECOND_THREE_DIGIT:
		return p.parseMillisecond()
	case TIMELIB_FORMAT_TIMEZONE_OFFSET:
		return p.parseTimezoneOffset()
	case TIMELIB_FORMAT_TIMEZONE_OFFSET_MINUTES:
		return p.parseTimezoneOffsetMinutes()
	case TIMELIB_FORMAT_TEXTUAL_MONTH_FULL:
		return p.parseTextualMonthFull()
	case TIMELIB_FORMAT_TEXTUAL_MONTH_3_LETTER:
		return p.parseTextualMonthShort()
	case TIMELIB_FORMAT_TEXTUAL_DAY_3_LETTER:
		return p.parseTextualDayShort()
	case TIMELIB_FORMAT_TEXTUAL_DAY_FULL:
		return p.parseTextualDayFull()
	case TIMELIB_FORMAT_MERIDIAN:
		return p.parseMeridian()
	case TIMELIB_FORMAT_EPOCH_SECONDS:
		return p.parseEpochSeconds()
	case TIMELIB_FORMAT_DAY_OF_YEAR:
		return p.parseDayOfYear()
	case TIMELIB_FORMAT_WEEK_OF_YEAR_ISO:
		return p.parseWeekOfYearISO()
	case TIMELIB_FORMAT_DAY_SUFFIX:
		return p.parseDaySuffix()
	case TIMELIB_FORMAT_WHITESPACE:
		return p.parseWhitespace()
	case TIMELIB_FORMAT_SEPARATOR:
		return p.parseSeparator()
	case TIMELIB_FORMAT_RANDOM_CHAR:
		return p.parseRandomChar()
	case TIMELIB_FORMAT_ALLOW_EXTRA_CHARACTERS:
		return p.parseAllowExtra()
	case TIMELIB_FORMAT_ANY_SEPARATOR:
		return p.parseAnySeparator()
	case TIMELIB_FORMAT_SKIP_TO_SEPARATOR:
		return p.parseSkipToSeparator()
	case TIMELIB_FORMAT_RESET_ALL:
		return p.parseResetAll()
	case TIMELIB_FORMAT_RESET_ALL_WHEN_NOT_SET:
		return p.parseResetAllWhenNotSet()
	case TIMELIB_FORMAT_ESCAPE:
		return p.parseEscape()
	default:
		p.addError(TIMELIB_ERR_INVALID_SPECIFIER, "Unsupported format specifier")
		return false
	}
}

// parseYearFourDigit parses 4-digit year
func (p *FormatParser) parseYearFourDigit() bool {
	if p.position+4 > len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected 4-digit year")
		return false
	}

	yearStr := p.input[p.position : p.position+4]
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid year")
		return false
	}

	p.time.Y = year
	p.position += 4
	return true
}

// parseYearTwoDigit parses 2-digit year
func (p *FormatParser) parseYearTwoDigit() bool {
	if p.position+2 > len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected 2-digit year")
		return false
	}

	yearStr := p.input[p.position : p.position+2]
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid year")
		return false
	}

	// Handle 2-digit year conversion (same as C implementation)
	if year < 70 {
		year += 2000
	} else if year < 100 {
		year += 1900
	}

	p.time.Y = year
	p.position += 2
	return true
}

// parseMonthTwoDigitPadded parses 2-digit padded month
func (p *FormatParser) parseMonthTwoDigitPadded() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected month")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	monthStr := p.input[start:p.position]
	month, err := strconv.ParseInt(monthStr, 10, 64)
	if err != nil || month < 1 || month > 12 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid month")
		return false
	}

	p.time.M = month
	return true
}

// parseMonthTwoDigit parses 1-2 digit month
func (p *FormatParser) parseMonthTwoDigit() bool {
	// Find the end of the number
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected month")
		return false
	}

	monthStr := p.input[start:p.position]
	month, err := strconv.ParseInt(monthStr, 10, 64)
	if err != nil || month < 1 || month > 12 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid month")
		return false
	}

	p.time.M = month
	return true
}

// parseDayTwoDigitPadded parses 2-digit padded day
func (p *FormatParser) parseDayTwoDigitPadded() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected day")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	dayStr := p.input[start:p.position]
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil || day < 1 || day > 31 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid day")
		return false
	}

	p.time.D = day
	return true
}

// parseDayTwoDigit parses 1-2 digit day
func (p *FormatParser) parseDayTwoDigit() bool {
	// Find the end of the number
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected day")
		return false
	}

	dayStr := p.input[start:p.position]
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil || day < 1 || day > 31 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid day")
		return false
	}

	p.time.D = day
	return true
}

// parseHour24 parses 24-hour format hour
func (p *FormatParser) parseHour24() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected hour")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	hourStr := p.input[start:p.position]
	hour, err := strconv.ParseInt(hourStr, 10, 64)
	if err != nil || hour < 0 || hour > 23 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid hour")
		return false
	}

	p.time.H = hour
	p.time.HaveTime = true
	return true
}

// parseHour12 parses 12-hour format hour (1-2 digits)
func (p *FormatParser) parseHour12() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected hour")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	hourStr := p.input[start:p.position]
	hour, err := strconv.ParseInt(hourStr, 10, 64)
	if err != nil || hour < 1 || hour > 12 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid 12-hour format")
		return false
	}

	p.time.H = hour
	p.time.HaveTime = true
	return true
}

// parseMinute parses minute
func (p *FormatParser) parseMinute() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected minute")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	minuteStr := p.input[start:p.position]
	minute, err := strconv.ParseInt(minuteStr, 10, 64)
	if err != nil || minute < 0 || minute > 59 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid minute")
		return false
	}

	p.time.I = minute
	p.time.HaveTime = true
	return true
}

// parseSecond parses second
func (p *FormatParser) parseSecond() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected second")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	secondStr := p.input[start:p.position]
	second, err := strconv.ParseInt(secondStr, 10, 64)
	if err != nil || second < 0 || second > 59 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid second")
		return false
	}

	p.time.S = second
	p.time.HaveTime = true
	return true
}

// parseMicrosecond parses microseconds
func (p *FormatParser) parseMicrosecond() bool {
	// Parse up to 6 digits directly (no separator required for 'u' format)
	start := p.position
	digits := 0
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' && digits < 6 {
		p.position++
		digits++
	}

	if digits == 0 {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected microseconds")
		return false
	}

	microStr := p.input[start:p.position]
	microseconds, err := strconv.ParseInt(microStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid microseconds")
		return false
	}

	// Pad to 6 digits
	for i := digits; i < 6; i++ {
		microseconds *= 10
	}

	p.time.US = microseconds
	p.time.HaveTime = true
	return true
}

// parseTimezoneOffset parses timezone offset
func (p *FormatParser) parseTimezoneOffset() bool {
	if p.position >= len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected timezone offset")
		return false
	}

	// Handle 'Z' for UTC
	if p.input[p.position] == 'Z' {
		p.time.Z = 0
		p.time.IsLocaltime = true
		p.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
		p.time.HaveZone = true
		p.position++
		return true
	}

	// Try to parse as textual timezone abbreviation first
	if offset := p.lookupTimezoneAbbr(); offset != -1 {
		p.time.Z = offset
		p.time.IsLocaltime = true
		p.time.ZoneType = TIMELIB_ZONETYPE_ABBR
		p.time.HaveZone = true
		return true
	}

	// Handle various timezone formats
	sign := int32(1)
	if p.input[p.position] == '+' {
		sign = 1
		p.position++
	} else if p.input[p.position] == '-' {
		sign = -1
		p.position++
	}

	// Parse hours
	if p.position+2 > len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected timezone hours")
		return false
	}

	hourStr := p.input[p.position : p.position+2]
	hours, err := strconv.ParseInt(hourStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid timezone hours")
		return false
	}
	p.position += 2

	// Check for optional colon and minutes
	var minutes int64 = 0
	if p.position < len(p.input) && p.input[p.position] == ':' {
		p.position++
		if p.position+2 > len(p.input) {
			p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected timezone minutes")
			return false
		}
		minuteStr := p.input[p.position : p.position+2]
		minutes, err = strconv.ParseInt(minuteStr, 10, 64)
		if err != nil {
			p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid timezone minutes")
			return false
		}
		p.position += 2
	} else if p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		// No colon, but 4 digits (HHMM format)
		if p.position+2 <= len(p.input) {
			minuteStr := p.input[p.position : p.position+2]
			minutes, err = strconv.ParseInt(minuteStr, 10, 64)
			if err != nil {
				p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid timezone minutes")
				return false
			}
			p.position += 2
		}
	}

	// Calculate total offset in seconds
	totalOffset := int32(sign) * (int32(hours)*3600 + int32(minutes)*60)
	p.time.Z = totalOffset
	p.time.IsLocaltime = true
	p.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
	p.time.HaveZone = true

	return true
}

// parseTextualMonthFull parses full month name
func (p *FormatParser) parseTextualMonthFull() bool {
	months := []string{
		"january", "february", "march", "april", "may", "june",
		"july", "august", "september", "october", "november", "december",
	}

	// Check for Roman numerals first (longer ones first to avoid partial matches)
	romanMonths := []string{"XII", "XI", "X", "IX", "VIII", "VII", "VI", "V", "IV", "III", "II", "I"}
	for monthNum, roman := range romanMonths {
		if p.position+len(roman) <= len(p.input) {
			inputPart := p.input[p.position : p.position+len(roman)]
			if inputPart == roman {
				p.time.M = int64(12 - monthNum) // Reverse mapping since we check longest first
				p.position += len(roman)
				return true
			}
		}
	}

	// Find longest match first for regular month names
	for monthNum, monthName := range months {
		// Try case-insensitive match
		if p.position+len(monthName) <= len(p.input) {
			inputPart := strings.ToLower(p.input[p.position : p.position+len(monthName)])
			if inputPart == monthName {
				p.time.M = int64(monthNum + 1)
				p.position += len(monthName)
				return true
			}
		}
	}

	// Also check for 3-letter abbreviations
	abbrMonths := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	for monthNum, abbr := range abbrMonths {
		if p.position+len(abbr) <= len(p.input) {
			inputPart := strings.ToLower(p.input[p.position : p.position+len(abbr)])
			if inputPart == abbr {
				p.time.M = int64(monthNum + 1)
				p.position += len(abbr)
				return true
			}
		}
	}

	p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid month name")
	return false
}

// parseTextualMonthShort parses 3-letter month abbreviation
func (p *FormatParser) parseTextualMonthShort() bool {
	months := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

	if p.position+3 > len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected 3-letter month")
		return false
	}

	monthStr := strings.ToLower(p.input[p.position : p.position+3])
	for monthNum, monthName := range months {
		if monthStr == monthName {
			p.time.M = int64(monthNum + 1)
			p.position += 3
			return true
		}
	}

	p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid month abbreviation")
	return false
}

// parseTextualDayShort parses 3-letter day abbreviation
func (p *FormatParser) parseTextualDayShort() bool {
	days := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

	if p.position+3 > len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected 3-letter day")
		return false
	}

	dayStr := strings.ToLower(p.input[p.position : p.position+3])
	for dayNum, dayName := range days {
		if dayStr == dayName {
			p.time.HaveDate = true
			p.time.Relative.Weekday = dayNum
			p.position += 3
			return true
		}
	}

	p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid day abbreviation")
	return false
}

// parseTextualDayFull parses full day name
func (p *FormatParser) parseTextualDayFull() bool {
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	// Find longest match first
	for dayNum, dayName := range days {
		if p.position+len(dayName) <= len(p.input) {
			inputPart := strings.ToLower(p.input[p.position : p.position+len(dayName)])
			if inputPart == dayName {
				p.time.HaveDate = true
				p.time.Relative.Weekday = dayNum
				p.position += len(dayName)
				return true
			}
		}
	}

	p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid day name")
	return false
}

// parseMeridian parses AM/PM
func (p *FormatParser) parseMeridian() bool {
	// Check for various AM/PM formats
	formats := []struct {
		pattern string
		length  int
		isPM    bool
	}{
		{"am", 2, false},
		{"pm", 2, true},
		{"a.m.", 4, false},
		{"p.m.", 4, true},
		{"AM", 2, false},
		{"PM", 2, true},
		{"A.M.", 4, false},
		{"P.M.", 4, true},
	}

	for _, format := range formats {
		if p.position+format.length <= len(p.input) {
			inputPart := p.input[p.position : p.position+format.length]
			if strings.EqualFold(inputPart, format.pattern) {
				p.time.HaveTime = true

				// Convert to 24-hour format for PM
				if format.isPM {
					if p.time.H >= 1 && p.time.H <= 11 {
						p.time.H += 12
					} else if p.time.H == 12 {
						p.time.H = 0 // 12 PM should be 0 in 24-hour format
					}
				} else {
					// AM format - 12 AM should be 0 in 24-hour format
					if p.time.H == 12 {
						p.time.H = 0
					}
				}

				p.position += format.length
				return true
			}
		}
	}

	p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Invalid meridian")
	return false
}

// parseEpochSeconds parses Unix timestamp
func (p *FormatParser) parseEpochSeconds() bool {
	// Find the end of the number
	start := p.position
	negative := false

	if p.position < len(p.input) && p.input[p.position] == '-' {
		negative = true
		p.position++
	}

	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start || (negative && p.position == start+1) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected timestamp")
		return false
	}

	timestampStr := p.input[start:p.position]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid timestamp")
		return false
	}

	// Convert Unix timestamp to date/time
	// Use Go's time package to convert the timestamp to UTC
	t := time.Unix(timestamp, 0).UTC()

	// Preserve existing microseconds if they were already set
	existingUS := p.time.US

	p.time.Y = int64(t.Year())
	p.time.M = int64(t.Month())
	p.time.D = int64(t.Day())
	p.time.H = int64(t.Hour())
	p.time.I = int64(t.Minute())
	p.time.S = int64(t.Second())

	// Only set microseconds from timestamp if they weren't already set
	if existingUS == 0 {
		p.time.US = int64(t.Nanosecond() / 1000)
	}

	p.time.HaveDate = true
	p.time.HaveTime = true

	return true
}

// parseDayOfYear parses day of year (1-366)
func (p *FormatParser) parseDayOfYear() bool {
	// Find the end of the number
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected day of year")
		return false
	}

	dayStr := p.input[start:p.position]
	dayOfYear, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil || dayOfYear < 1 || dayOfYear > 366 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid day of year")
		return false
	}

	// Convert day of year to month and day
	// This is a simplified version - full implementation would handle leap years
	if p.time.Y > 0 {
		// We have a year, so we can calculate the actual date
		isLeap := IsLeapYear(p.time.Y)

		// Check if day of year exceeds days in current year
		totalDays := 365
		if isLeap {
			totalDays = 366
		}

		adjustedDayOfYear := dayOfYear
		if adjustedDayOfYear > int64(totalDays) {
			// Roll over to next year
			p.time.Y++
			adjustedDayOfYear -= int64(totalDays)
			// Recalculate leap year for the new year
			isLeap = IsLeapYear(p.time.Y)
		}

		month, day := dayOfYearToDate(adjustedDayOfYear, isLeap)
		p.time.M = month
		p.time.D = day
	} else {
		// No year set yet, store as relative
		p.time.Relative.D = dayOfYear - 1 // Convert to 0-based
		// Store as relative day offset
		p.time.Relative.D = dayOfYear - 1 // Convert to 0-based
	}

	return true
}

// parseWhitespace parses whitespace
func (p *FormatParser) parseWhitespace() bool {
	for p.position < len(p.input) && (p.input[p.position] == ' ' || p.input[p.position] == '\t') {
		p.position++
	}
	return true
}

// parseSeparator parses separator characters
func (p *FormatParser) parseSeparator() bool {
	if p.position >= len(p.input) {
		return true // Optional separator
	}

	// Skip common separator characters
	for p.position < len(p.input) {
		char := p.input[p.position]
		if char == ' ' || char == '\t' || char == '-' || char == '/' || char == '.' || char == ',' || char == ':' || char == ';' {
			p.position++
		} else {
			break
		}
	}
	return true
}

// parseRandomChar parses any character (used with ? format)
func (p *FormatParser) parseRandomChar() bool {
	if p.position < len(p.input) {
		p.position++
	}
	return true
}

// parseAllowExtra allows extra characters
func (p *FormatParser) parseAllowExtra() bool {
	// This is a simplified implementation
	// In the full implementation, this would allow extra characters in the input
	return true
}

// parseAnySeparator allows any separator
func (p *FormatParser) parseAnySeparator() bool {
	// Skip any separator character
	if p.position < len(p.input) {
		char := p.input[p.position]
		if char == ' ' || char == '\t' || char == '-' || char == '/' || char == '.' || char == ',' || char == ':' || char == ';' {
			p.position++
		}
	}
	return true
}

// Helper function to convert day of year to month and day
func dayOfYearToDate(dayOfYear int64, isLeap bool) (int64, int64) {
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeap {
		daysInMonth[1] = 29
	}

	// C library behavior: day 1 = January 2nd, day 2 = January 3rd, etc.
	// So we need to add 1 to match the C library's expectations
	remaining := dayOfYear + 1

	// Check if the day exceeds the number of days in the year
	totalDays := 365
	if isLeap {
		totalDays = 366
	}

	if remaining > int64(totalDays) {
		// Roll over to next year - subtract the days in current year
		remaining -= int64(totalDays)
	}

	for month := 0; month < 12; month++ {
		if remaining <= int64(daysInMonth[month]) {
			return int64(month + 1), remaining
		}
		remaining -= int64(daysInMonth[month])
	}

	return 12, remaining
}

// parseYearISO parses ISO year format
func (p *FormatParser) parseYearISO() bool {
	return p.parseYearFourDigit()
}

// parseDayOfWeekISO parses ISO day of week (1-7, Monday=1)
func (p *FormatParser) parseDayOfWeekISO() bool {
	if p.position+1 > len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected ISO day of week")
		return false
	}

	dayStr := p.input[p.position : p.position+1]
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil || day < 1 || day > 7 {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid ISO day of week")
		return false
	}

	p.time.Relative.Weekday = int(day - 1) // Convert to 0-based (Sunday=0)
	p.time.HaveDate = true
	p.position += 1
	return true
}

// parseMillisecond parses 3-digit millisecond
func (p *FormatParser) parseMillisecond() bool {
	// Find decimal point or colon
	if p.position >= len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected milliseconds")
		return false
	}

	sep := p.input[p.position]
	if sep != '.' && sep != ':' {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected decimal separator")
		return false
	}
	p.position++

	// Parse up to 3 digits
	start := p.position
	digits := 0
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' && digits < 3 {
		p.position++
		digits++
	}

	if digits == 0 {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected milliseconds")
		return false
	}

	milliStr := p.input[start:p.position]
	milliseconds, err := strconv.ParseInt(milliStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid milliseconds")
		return false
	}

	// Convert to microseconds (pad to 6 digits)
	for i := digits; i < 3; i++ {
		milliseconds *= 10
	}
	p.time.US = milliseconds * 1000
	p.time.HaveTime = true
	return true
}

// parseTimezoneOffsetMinutes parses timezone offset in minutes
func (p *FormatParser) parseTimezoneOffsetMinutes() bool {
	if p.position >= len(p.input) {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected timezone offset in minutes")
		return false
	}

	// Handle 'Z' for UTC
	if p.input[p.position] == 'Z' {
		p.time.Z = 0
		p.time.IsLocaltime = true
		p.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
		p.time.HaveZone = true
		p.position++
		return true
	}

	// Handle sign
	sign := int32(1)
	if p.input[p.position] == '+' {
		sign = 1
		p.position++
	} else if p.input[p.position] == '-' {
		sign = -1
		p.position++
	}

	// Parse minutes
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERROR_UNEXPECTED_DATA, "Expected timezone minutes")
		return false
	}

	minutesStr := p.input[start:p.position]
	minutes, err := strconv.ParseInt(minutesStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid timezone minutes")
		return false
	}

	// Convert to seconds
	p.time.Z = int32(sign) * int32(minutes) * 60
	p.time.IsLocaltime = true
	p.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
	p.time.HaveZone = true
	return true
}

// parseWeekOfYearISO parses ISO week of year
func (p *FormatParser) parseWeekOfYearISO() bool {
	// Try to parse 2-digit week first
	if p.position+2 <= len(p.input) {
		weekStr := p.input[p.position : p.position+2]
		week, err := strconv.ParseInt(weekStr, 10, 64)
		if err == nil && week >= 1 && week <= 53 {
			// Store as relative information for ISO week calculations
			p.time.HaveDate = true
			p.time.Relative.Week = int(week)
			p.position += 2
			return true
		}
	}

	// Try 1-digit week
	if p.position+1 <= len(p.input) {
		weekStr := p.input[p.position : p.position+1]
		week, err := strconv.ParseInt(weekStr, 10, 64)
		if err == nil && week >= 1 && week <= 9 {
			// Store as relative information for ISO week calculations
			p.time.HaveDate = true
			p.time.Relative.Week = int(week)
			p.position += 1
			return true
		}
	}

	p.addError(TIMELIB_ERROR_NUMBER_OUT_OF_RANGE, "Invalid ISO week")
	return false
}

// convertISOWeekToDate converts ISO week date to calendar date
func (p *FormatParser) convertISOWeekToDate() {
	// If we have ISO week information, convert it to calendar date
	if p.time.Relative.Week > 0 && p.time.Y != TIMELIB_UNSET {
		// Default to Monday (ISO day 1) if no ISO day is set
		isoDay := 1
		if p.time.Relative.Weekday >= 0 {
			isoDay = p.time.Relative.Weekday + 1 // Convert from 0-based to 1-based
		}

		// Convert ISO week date to calendar date
		y, m, d := DateFromIsoDate(p.time.Y, int64(p.time.Relative.Week), int64(isoDay))
		p.time.Y = y
		p.time.M = m
		p.time.D = d

		// Clear the ISO week information since we've converted it
		p.time.Relative.Week = -1
		p.time.Relative.Weekday = -1
		return
	}

	// If we have a year but no specific date, default to January 2nd (as per C implementation)
	if p.time.Y != TIMELIB_UNSET && p.time.M == TIMELIB_UNSET && p.time.D == TIMELIB_UNSET {
		// Default to January 2nd for ISO year-only format
		p.time.M = 1
		p.time.D = 2
		return
	}
}

// parseSkipToSeparator skips to next separator
func (p *FormatParser) parseSkipToSeparator() bool {
	// Skip until we find a separator character
	for p.position < len(p.input) {
		char := p.input[p.position]
		if char == ' ' || char == '\t' || char == '-' || char == '/' || char == '.' || char == ',' || char == ':' || char == ';' {
			break
		}
		p.position++
	}
	return true
}

// parseResetAll resets all time components
func (p *FormatParser) parseResetAll() bool {
	// Reset all time components to unset
	p.time.Y = -9999999
	p.time.M = -9999999
	p.time.D = -9999999
	p.time.H = -9999999
	p.time.I = -9999999
	p.time.S = -9999999
	p.time.US = 0
	p.time.Z = 0
	return true
}

// parseResetAllWhenNotSet resets all unset time components
func (p *FormatParser) parseResetAllWhenNotSet() bool {
	// Reset only unset time components
	if p.time.Y == -9999999 {
		p.time.Y = 0
	}
	if p.time.M == -9999999 {
		p.time.M = 0
	}
	if p.time.D == -9999999 {
		p.time.D = 0
	}
	if p.time.H == -9999999 {
		p.time.H = 0
	}
	if p.time.I == -9999999 {
		p.time.I = 0
	}
	if p.time.S == -9999999 {
		p.time.S = 0
	}
	return true
}

// parseEscape handles escape character
func (p *FormatParser) parseEscape() bool {
	// Skip the next character (it's escaped)
	if p.position < len(p.input) {
		p.position++
	}
	return true
}

// parseDaySuffix parses day suffix like "st", "nd", "rd", "th"
func (p *FormatParser) parseDaySuffix() bool {
	if p.position+2 > len(p.input) {
		return true // Optional suffix
	}

	suffix := p.input[p.position : p.position+2]
	suffixLower := strings.ToLower(suffix)
	switch suffixLower {
	case "st", "nd", "rd", "th":
		p.position += 2
		return true
	}
	return true // Ignore if not matching
}

// lookupTimezoneAbbr looks up timezone abbreviation and returns offset in seconds
func (p *FormatParser) lookupTimezoneAbbr() int32 {
	// Common timezone abbreviations and their offsets in seconds
	timezoneAbbrs := map[string]int32{
		"UTC":  0,
		"GMT":  0,
		"CEST": 7200,   // +02:00
		"CET":  3600,   // +01:00
		"EST":  -18000, // -05:00
		"EDT":  -14400, // -04:00
		"PST":  -28800, // -08:00
		"PDT":  -25200, // -07:00
		"MST":  -25200, // -07:00
		"MDT":  -21600, // -06:00
		"CST":  -21600, // -06:00
		"CDT":  -18000, // -05:00
		"JST":  32400,  // +09:00
		"HKT":  28800,  // +08:00
		"IST":  19800,  // +05:30
		"BST":  3600,   // +01:00
	}

	// Try to match abbreviations of different lengths (longest first)
	for length := 4; length >= 3; length-- {
		if p.position+length <= len(p.input) {
			abbr := strings.ToUpper(p.input[p.position : p.position+length])
			if offset, exists := timezoneAbbrs[abbr]; exists {
				p.position += length
				return offset
			}
		}
	}

	return -1 // Not found
}

// addError adds an error to the error container
func (p *FormatParser) addError(code int, message string) {
	p.errors.ErrorCount++
	p.errors.ErrorMessages = append(p.errors.ErrorMessages, ErrorMessage{
		ErrorCode: code,
		Position:  p.position,
		Character: 0,
		Message:   message,
	})
}
