package timelib

// Token type constants for ISO interval parsing
// Extracted from parse_iso_intervals.re
// Note: These are local to the ISO interval parser and may differ
// from similar constants in parse_date_constants.go
const (
	// EOI is already defined in parse_date_constants.go
	// EOI = 257

	// ISO interval specific tokens
	TIMELIB_PERIOD          = 260
	TIMELIB_ISO_DATE_INTRVL = 261 // Renamed to avoid conflict with parse_date's TIMELIB_ISO_DATE
	// TIMELIB_ERROR is already defined in parse_date_constants.go
	// TIMELIB_ERROR = 999
)
