package timelib

// ParseFromFormatWithMap parses a date/time string according to a format with format specifier mapping
func ParseFromFormatWithMap(format, input string, tzdb *TzDB, tzWrapper func(string, *TzDB, *int) (*TzInfo, error), formatConfig *FormatConfig) (*Time, *ErrorContainer) {
	errors := &ErrorContainer{
		ErrorMessages:   make([]ErrorMessage, 0),
		WarningMessages: make([]ErrorMessage, 0),
	}

	// Validate inputs
	if format == "" {
		errors.addError(TIMELIB_ERR_EMPTY_STRING, "Format string is empty")
		return nil, errors
	}

	if input == "" {
		errors.addError(TIMELIB_ERR_EMPTY_STRING, "Input string is empty")
		return nil, errors
	}

	// Use our custom format parser instead of the basic Go parser
	result, parseErrors := ParseFromFormatWithConfig(format, input, formatConfig)
	if parseErrors != nil && parseErrors.ErrorCount > 0 {
		// Copy errors from parseErrors to our errors container
		errors.ErrorMessages = append(errors.ErrorMessages, parseErrors.ErrorMessages...)
		errors.ErrorCount += parseErrors.ErrorCount
		return nil, errors
	}

	return result, errors
}

// ParseOptions holds parsing configuration options
type ParseOptions struct {
	AllowExtraChars bool
	StrictMode      bool
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
