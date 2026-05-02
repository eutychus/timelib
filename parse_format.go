package timelib

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatParser handles format-based parsing
type FormatParser struct {
	input             string
	format            string
	position          int
	formatPos         int
	time              *Time
	errors            *ErrorContainer
	config            *FormatConfig
	haveISOYear       bool
	haveISOWeek       bool
	haveISODay        bool
	haveMeridian      bool
	haveHour          bool
	haveDOY           bool
	haveYear          bool
	meridianBeforeHour bool
	doyBeforeYear     bool
	mixedISOWithNatural bool
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
			{'h', TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX},
			{'g', TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX_PADDED},
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
			{'x', TIMELIB_FORMAT_YEAR_EXPANDED},
			{'X', TIMELIB_FORMAT_YEAR_EXPANDED},
			{':', TIMELIB_FORMAT_SEPARATOR},
			{'/', TIMELIB_FORMAT_SEPARATOR},
			{'.', TIMELIB_FORMAT_SEPARATOR},
			{',', TIMELIB_FORMAT_SEPARATOR},
			{'-', TIMELIB_FORMAT_SEPARATOR},
			{'(', TIMELIB_FORMAT_SEPARATOR},
			{')', TIMELIB_FORMAT_SEPARATOR},
			{';', TIMELIB_FORMAT_SEPARATOR},
		},
		PrefixChar: 0,
	}
	return ParseFromFormatWithConfig(format, input, defaultConfig)
}

// ParseFromFormatWithPrefix parses input according to format using % as prefix character
// This matches the C version's test_parse_with_prefix behavior
func ParseFromFormatWithPrefix(format, input string) (*Time, *ErrorContainer) {
	// Use default format configuration with % prefix
	configWithPrefix := &FormatConfig{
		FormatMap: []FormatSpecifier{
			{'Y', TIMELIB_FORMAT_YEAR_FOUR_DIGIT},
			{'y', TIMELIB_FORMAT_YEAR_TWO_DIGIT},
			{'m', TIMELIB_FORMAT_MONTH_TWO_DIGIT_PADDED},
			{'n', TIMELIB_FORMAT_MONTH_TWO_DIGIT},
			{'d', TIMELIB_FORMAT_DAY_TWO_DIGIT_PADDED},
			{'j', TIMELIB_FORMAT_DAY_TWO_DIGIT},
			{'H', TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX},
			{'h', TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX},
			{'g', TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX_PADDED},
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
			{':', TIMELIB_FORMAT_SEPARATOR},
			{'/', TIMELIB_FORMAT_SEPARATOR},
			{'.', TIMELIB_FORMAT_SEPARATOR},
			{',', TIMELIB_FORMAT_SEPARATOR},
			{'-', TIMELIB_FORMAT_SEPARATOR},
			{'(', TIMELIB_FORMAT_SEPARATOR},
			{')', TIMELIB_FORMAT_SEPARATOR},
			{';', TIMELIB_FORMAT_SEPARATOR},
		},
		PrefixChar: '%',
	}
	return ParseFromFormatWithConfig(format, input, configWithPrefix)
}

// Parse performs the format parsing
func (p *FormatParser) Parse() *Time {
	p.position = 0
	p.formatPos = 0
	prefixFound := false

	for p.formatPos < len(p.format) && p.position < len(p.input) {
		formatChar := rune(p.format[p.formatPos])

		// Handle prefix character if configured
		if p.config.PrefixChar != 0 {
			// Check if this is a prefix character or a literal match
			if !prefixFound && formatChar != rune(p.config.PrefixChar) {
				if p.position < len(p.input) && p.matchCharacter(formatChar) {
					p.position++
				} else {
					p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Literal character mismatch")
					if p.position < len(p.input) {
						p.position++
					}
				}
				p.formatPos++
				continue
			}

			if formatChar == rune(p.config.PrefixChar) {
				// Found prefix character
				if prefixFound {
					// Sequential prefix characters - second one is escaped/literal
					if !p.matchCharacter(formatChar) {
						p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected prefix character as literal")
						return p.time
					}
					p.position++
					p.formatPos++
					prefixFound = false
					continue
				}
				prefixFound = true
				p.formatPos++
				continue
			}

			// If we get here with prefixFound=true, the next character is a format specifier
			if prefixFound {
				spec := p.findFormatSpecifier(formatChar)
				if spec == nil {
					p.addError(TIMELIB_ERR_INVALID_SPECIFIER, fmt.Sprintf("Invalid format specifier after prefix: %c", formatChar))
					return p.time
				}
				if !p.parseFormatSpecifier(spec) {
					return p.time
				}
				p.formatPos++
				prefixFound = false
				continue
			}
		}

		// No prefix character configured - original behavior
		// Handle escape character
		if formatChar == '\\' {
			if p.formatPos+1 < len(p.format) {
				p.formatPos++ // Move to the character after backslash
				expectedChar := rune(p.format[p.formatPos])
				// Check if we have enough input left
				if p.position >= len(p.input) {
					p.addError(TIMELIB_ERR_UNEXPECTED_DATA, fmt.Sprintf("Expected character '%c' after escape, but input ended", expectedChar))
					return p.time
				}
				if !p.matchCharacter(expectedChar) {
					p.addError(TIMELIB_ERR_UNEXPECTED_DATA, fmt.Sprintf("Expected character '%c' after escape, got '%c'", expectedChar, p.input[p.position]))
					return p.time
				}
				p.position++  // Advance input position past the matched character
				p.formatPos++ // Move past the escaped character
				continue
			} else {
				// Backslash at end of format string - treat as literal
				if p.position >= len(p.input) || p.input[p.position] != '\\' {
					p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected backslash character")
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
				p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Literal character mismatch")
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

	// If month is 0 and day > 0, treat as day-of-year
	if p.time.M == 0 && p.time.D > 0 && p.time.Y > 0 {
		isLeap := IsLeapYear(p.time.Y)
		doy := p.time.D
		daysInMonth := []int64{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		if isLeap {
			daysInMonth[1] = 29
		}
		remaining := doy
		for m, dim := range daysInMonth {
			if remaining <= dim {
				p.time.M = int64(m + 1)
				p.time.D = remaining
				break
			}
			remaining -= dim
		}
		p.time.HaveDate = true
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
	}

	p.validate()

	return p.time
}

func (p *FormatParser) validate() {
	if (p.haveISOWeek || p.haveISODay) && !p.haveISOYear {
		p.addWarning(TIMELIB_ERR_UNEXPECTED_DATA, "ISO week/day without ISO year")
	}
	if p.meridianBeforeHour {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Meridian before hour")
	} else if p.haveMeridian && !p.haveHour {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Meridian without hour")
	}
	if p.doyBeforeYear {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Day of year before year")
	}
	if p.mixedISOWithNatural {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Cannot mix ISO with natural date")
	}
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
		p.haveYear = true
		if p.haveDOY {
			p.doyBeforeYear = true
		}
		return p.parseYearFourDigit()
	case TIMELIB_FORMAT_YEAR_TWO_DIGIT:
		p.haveYear = true
		if p.haveDOY {
			p.doyBeforeYear = true
		}
		return p.parseYearTwoDigit()
	case TIMELIB_FORMAT_YEAR_EXPANDED:
		p.haveYear = true
		if p.haveDOY {
			p.doyBeforeYear = true
		}
		return p.parseYearExpanded()
	case TIMELIB_FORMAT_YEAR_ISO:
		p.haveISOYear = true
		if p.haveYear {
			p.mixedISOWithNatural = true
		}
		return p.parseYearISO()
	case TIMELIB_FORMAT_MONTH_TWO_DIGIT_PADDED:
		if p.haveISOYear {
			p.mixedISOWithNatural = true
		}
		return p.parseMonthTwoDigitPadded()
	case TIMELIB_FORMAT_MONTH_TWO_DIGIT:
		if p.haveISOYear {
			p.mixedISOWithNatural = true
		}
		return p.parseMonthTwoDigit()
	case TIMELIB_FORMAT_DAY_TWO_DIGIT_PADDED:
		if p.haveISOYear {
			p.mixedISOWithNatural = true
		}
		return p.parseDayTwoDigitPadded()
	case TIMELIB_FORMAT_DAY_TWO_DIGIT:
		if p.haveISOYear {
			p.mixedISOWithNatural = true
		}
		return p.parseDayTwoDigit()
	case TIMELIB_FORMAT_DAY_OF_WEEK_ISO:
		p.haveISODay = true
		return p.parseDayOfWeekISO()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX:
		p.haveHour = true
		return p.parseHour24()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX_PADDED:
		p.haveHour = true
		return p.parseHour24()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX:
		p.haveHour = true
		return p.parseHour12()
	case TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX_PADDED:
		p.haveHour = true
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
		if p.haveISOYear {
			p.mixedISOWithNatural = true
		}
		return p.parseTextualMonthFull()
	case TIMELIB_FORMAT_TEXTUAL_MONTH_3_LETTER:
		if p.haveISOYear {
			p.mixedISOWithNatural = true
		}
		return p.parseTextualMonthShort()
	case TIMELIB_FORMAT_TEXTUAL_DAY_3_LETTER:
		return p.parseTextualDayShort()
	case TIMELIB_FORMAT_TEXTUAL_DAY_FULL:
		return p.parseTextualDayFull()
	case TIMELIB_FORMAT_MERIDIAN:
		if !p.haveHour {
			p.meridianBeforeHour = true
		}
		p.haveMeridian = true
		return p.parseMeridian()
	case TIMELIB_FORMAT_EPOCH_SECONDS:
		return p.parseEpochSeconds()
	case TIMELIB_FORMAT_DAY_OF_YEAR:
		p.haveDOY = true
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected 4-digit year")
		return false
	}

	yearStr := p.input[p.position : p.position+4]
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid year")
		return false
	}

	p.time.Y = year
	p.position += 4
	return true
}

// parseYearTwoDigit parses 2-digit year
func (p *FormatParser) parseYearTwoDigit() bool {
	if p.position+2 > len(p.input) {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected 2-digit year")
		return false
	}

	yearStr := p.input[p.position : p.position+2]
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid year")
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

// parseYearExpanded parses an expanded year with optional sign and up to 19 digits
func (p *FormatParser) parseYearExpanded() bool {
	if p.position >= len(p.input) {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected expanded year")
		return false
	}

	// Check for optional sign
	sign := int64(1)
	start := p.position
	if p.input[p.position] == '+' || p.input[p.position] == '-' {
		if p.input[p.position] == '-' {
			sign = -1
		}
		p.position++
	}

	// Parse up to 19 digits
	digitStart := p.position
	for p.position < len(p.input) && p.position-digitStart < 19 && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == digitStart {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected digits in expanded year")
		p.position = start // Reset position
		return false
	}

	yearStr := p.input[digitStart:p.position]
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid expanded year")
		p.position = start // Reset position
		return false
	}

	p.time.Y = year * sign
	p.time.HaveDate = true
	return true
}

// parseMonthTwoDigitPadded parses 2-digit padded month
func (p *FormatParser) parseMonthTwoDigitPadded() bool {
	start := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}

	if p.position == start {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected month")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	monthStr := p.input[start:p.position]
	month, err := strconv.ParseInt(monthStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid month")
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected month")
		return false
	}

	monthStr := p.input[start:p.position]
	month, err := strconv.ParseInt(monthStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid month")
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected day")
		return false
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	dayStr := p.input[start:p.position]
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid day")
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected day")
		return false
	}

	dayStr := p.input[start:p.position]
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid day")
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
		return true
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	hourStr := p.input[start:p.position]
	hour, err := strconv.ParseInt(hourStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid hour")
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
		return true
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	hourStr := p.input[start:p.position]
	hour, err := strconv.ParseInt(hourStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid 12-hour format")
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
		return true
	}

	if p.position-start > 2 {
		p.position = start + 2 // Limit to 2 digits max
	}

	minuteStr := p.input[start:p.position]
	minute, err := strconv.ParseInt(minuteStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid minute")
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
		return true
	}

	if p.position-start > 2 {
		p.position = start + 2
	}

	secondStr := p.input[start:p.position]
	second, err := strconv.ParseInt(secondStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid second")
		return false
	}

	p.time.S = second
	p.time.HaveTime = true
	return true
}

// parseMicrosecond parses microseconds
func (p *FormatParser) parseMicrosecond() bool {
	start := p.position
	digits := 0
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' && digits < 6 {
		p.position++
		digits++
	}

	if digits == 0 {
		return true
	}

	microStr := p.input[start:p.position]
	microseconds, err := strconv.ParseInt(microStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid microseconds")
		return false
	}

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
		return true
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
	if abbrResult := p.tryLookupTimezoneAbbr(); abbrResult != nil {
		p.time.Z = abbrResult.OffsetSec
		p.time.IsLocaltime = true
		p.time.ZoneType = TIMELIB_ZONETYPE_ABBR
		p.time.HaveZone = true
		p.time.TzAbbr = strings.ToUpper(abbrResult.Abbr)
		return true
	}

	// Try to parse as timezone identifier (e.g., Europe/Amsterdam, America/New_York)
	if p.tryParseTzIdentifier() {
		return true
	}

	// Handle numeric timezone offset formats
	return p.parseNumericTzOffset()
}

func (p *FormatParser) tryParseTzIdentifier() bool {
	start := p.position

	// Read a timezone identifier: letters, digits, /, _, -
	end := p.position
	for end < len(p.input) {
		c := p.input[end]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '/' || c == '_' || c == '-' || c == '+' {
			end++
		} else {
			break
		}
	}

	if end == start {
		return false
	}

	tzID := p.input[start:end]

	// Try to parse as a tzfile
	errorCode := 0
	tz, err := ParseTzfile(tzID, nil, &errorCode)
	if err == nil && tz != nil {
		p.time.TzInfo = tz
		p.time.ZoneType = TIMELIB_ZONETYPE_ID
		p.time.IsLocaltime = true
		p.time.HaveZone = true
		p.position = end
		return true
	}

	return false
}

func (p *FormatParser) parseNumericTzOffset() bool {
	sign := int32(1)
	if p.position < len(p.input) && p.input[p.position] == '+' {
		sign = 1
		p.position++
	} else if p.position < len(p.input) && p.input[p.position] == '-' {
		sign = -1
		p.position++
	}

	// Read all digits first
	digitStart := p.position
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
		p.position++
	}
	totalDigits := p.position - digitStart

	if totalDigits == 0 {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected timezone hours")
		return false
	}

	var hours, minutes int64
	var err error

	// Check if next char is a colon (HH:MM format)
	if p.position < len(p.input) && p.input[p.position] == ':' {
		// Colon format: digits before colon are hours
		hours, _ = strconv.ParseInt(p.input[digitStart:p.position], 10, 64)
		p.position++ // consume colon
		// Read minutes after colon
		minStart := p.position
		for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
			p.position++
		}
		if p.position > minStart {
			minutes, err = strconv.ParseInt(p.input[minStart:p.position], 10, 64)
			if err != nil {
				minutes = 0
			}
		}
	} else {
		// No colon format: split digits into hours and minutes
		switch totalDigits {
		case 1:
			hours, _ = strconv.ParseInt(p.input[digitStart:p.position], 10, 64)
		case 2:
			hours, _ = strconv.ParseInt(p.input[digitStart:p.position], 10, 64)
		case 3:
			hours, _ = strconv.ParseInt(p.input[digitStart:digitStart+1], 10, 64)
			minutes, _ = strconv.ParseInt(p.input[digitStart+1:p.position], 10, 64)
		default:
			hours, _ = strconv.ParseInt(p.input[digitStart:digitStart+2], 10, 64)
			minutes, _ = strconv.ParseInt(p.input[digitStart+2:p.position], 10, 64)
		}
	}

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

	p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Invalid month name")
	return false
}

// parseTextualMonthShort parses 3-letter month abbreviation
func (p *FormatParser) parseTextualMonthShort() bool {
	months := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

	if p.position+3 > len(p.input) {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected 3-letter month")
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

	p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Invalid month abbreviation")
	return false
}

// parseTextualDayShort parses 3-letter day abbreviation
func (p *FormatParser) parseTextualDayShort() bool {
	days := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

	if p.position+3 > len(p.input) {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected 3-letter day")
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

	p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Invalid day abbreviation")
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

	p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Invalid day name")
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

	p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Invalid meridian")
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected timestamp")
		return false
	}

	timestampStr := p.input[start:p.position]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid timestamp")
		return false
	}

	// Preserve existing microseconds if they were already set
	// This is important when parsing formats like "u U" (microseconds then epoch)
	existingUS := p.time.US

	// Match C behavior: set SSE and timezone info, then call update_from_sse
	// C code at parse_date.c for TIMELIB_FORMAT_EPOCH_SECONDS
	p.time.HaveZone = true
	p.time.Sse = timestamp
	p.time.IsLocaltime = true
	p.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
	p.time.Z = 0
	p.time.Dst = 0
	p.time.UpdateFromSSE()

	// Restore microseconds if they were set before parsing epoch
	if existingUS != 0 {
		p.time.US = existingUS
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected day of year")
		return false
	}

	dayStr := p.input[start:p.position]
	dayOfYear, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil || dayOfYear < 1 || dayOfYear > 366 {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid day of year")
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
// In the C implementation, this matches ONE separator character and increments position once
// It does NOT consume multiple consecutive separators (that would be parseAnySeparator or similar)
func (p *FormatParser) parseSeparator() bool {
	if p.position >= len(p.input) {
		return true // Optional separator
	}

	// Match ONE separator character (like C implementation does with ++ptr)
	char := p.input[p.position]
	if char == ' ' || char == '\t' || char == '-' || char == '/' || char == '.' || char == ',' || char == ':' || char == ';' || char == '(' || char == ')' {
		p.position++
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

func dayOfYearToDate(dayOfYear int64, isLeap bool) (int64, int64) {
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeap {
		daysInMonth[1] = 29
	}

	remaining := dayOfYear + 1

	totalDays := 365
	if isLeap {
		totalDays = 366
	}

	if remaining > int64(totalDays) {
		remaining -= int64(totalDays)
	}

	for monthIdx, dim := range daysInMonth {
		if remaining <= int64(dim) {
			return int64(monthIdx + 1), remaining
		}
		remaining -= int64(dim)
	}

	return 12, remaining
}

// parseMillisecond parses 3-digit millisecond
func (p *FormatParser) parseMillisecond() bool {
	start := p.position
	digits := 0
	for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' && digits < 3 {
		p.position++
		digits++
	}

	if digits == 0 {
		return true
	}

	milliStr := p.input[start:p.position]
	milliseconds, err := strconv.ParseInt(milliStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid milliseconds")
		return false
	}

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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected timezone offset in minutes")
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
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected timezone minutes")
		return false
	}

	minutesStr := p.input[start:p.position]
	minutes, err := strconv.ParseInt(minutesStr, 10, 64)
	if err != nil {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid timezone minutes")
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
	p.haveISOWeek = true
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

	p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid ISO week")
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

	if p.haveISOYear && p.time.Y != TIMELIB_UNSET && p.time.M == TIMELIB_UNSET && p.time.D == TIMELIB_UNSET {
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

// tryLookupTimezoneAbbr looks up timezone abbreviation and returns the entry if found
func (p *FormatParser) tryLookupTimezoneAbbr() *TimezoneAbbreviation {
	start := p.position

	// Read the full word (same characters as tryParseTzIdentifier)
	end := p.position
	for end < len(p.input) {
		c := p.input[end]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '/' || c == '_' || c == '-' || c == '+' {
			end++
		} else {
			break
		}
	}

	if end == start {
		return nil
	}

	word := p.input[start:end]

	// Only check abbreviation table if word length is less than MAX_ABBR_LEN (6)
	if len(word) < 6 {
		entry := LookupTimezoneAbbr(word, -1, -1)
		if entry != nil {
			p.position = end
			return entry
		}
	}

	return nil
}

// addError adds an error to the error container
func (p *FormatParser) parseYearISO() bool {
	p.haveISOYear = true
	return p.parseYearFourDigit()
}

func (p *FormatParser) parseDayOfWeekISO() bool {
	if p.position+1 > len(p.input) {
		p.addError(TIMELIB_ERR_UNEXPECTED_DATA, "Expected ISO day of week")
		return false
	}

	dayStr := p.input[p.position : p.position+1]
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil || day < 1 || day > 7 {
		p.addError(TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Invalid ISO day of week")
		return false
	}

	p.time.Relative.Weekday = int(day - 1)
	p.time.HaveDate = true
	p.position += 1
	return true
}

func (p *FormatParser) parseAllowExtra() bool {
	return true
}

func (p *FormatParser) parseAnySeparator() bool {
	for p.position < len(p.input) {
		c := p.input[p.position]
		if c == ' ' || c == '\t' || c == '-' || c == '/' || c == '.' || c == ',' || c == ':' || c == ';' {
			p.position++
		} else {
			break
		}
	}
	return true
}

func (p *FormatParser) addError(code int, message string) {
	p.errors.ErrorCount++
	p.errors.ErrorMessages = append(p.errors.ErrorMessages, ErrorMessage{
		ErrorCode: code,
		Position:  p.position,
		Character: 0,
		Message:   message,
	})
}

func (p *FormatParser) addWarning(code int, message string) {
	p.errors.WarningCount++
	p.errors.WarningMessages = append(p.errors.WarningMessages, ErrorMessage{
		ErrorCode: code,
		Position:  p.position,
		Character: 0,
		Message:   message,
	})
}
