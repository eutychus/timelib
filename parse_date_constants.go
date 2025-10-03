// Code generated from parse_date.re; DO NOT EDIT.
// To regenerate, run: make parse_date_constants.go

package timelib

// Token type constants extracted from parse_date.re
const (
	EOI  = 257
	TIME = 258
	DATE = 259

	TIMELIB_XMLRPC_SOAP           = 260
	TIMELIB_TIME12                = 261
	TIMELIB_TIME24                = 262
	TIMELIB_GNU_NOCOLON           = 263
	TIMELIB_GNU_NOCOLON_TZ        = 264
	TIMELIB_ISO_NOCOLON           = 265
	TIMELIB_AMERICAN              = 266
	TIMELIB_ISO_DATE              = 267
	TIMELIB_DATE_FULL             = 268
	TIMELIB_DATE_TEXT             = 269
	TIMELIB_DATE_NOCOLON          = 270
	TIMELIB_PG_YEARDAY            = 271
	TIMELIB_PG_TEXT               = 272
	TIMELIB_PG_REVERSE            = 273
	TIMELIB_CLF                   = 274
	TIMELIB_DATE_NO_DAY           = 275
	TIMELIB_SHORTDATE_WITH_TIME   = 276
	TIMELIB_DATE_FULL_POINTED     = 277
	TIMELIB_TIME24_WITH_ZONE      = 278
	TIMELIB_ISO_WEEK              = 279
	TIMELIB_LF_DAY_OF_MONTH       = 280
	TIMELIB_WEEK_DAY_OF_MONTH     = 281
	TIMELIB_TIMEZONE              = 300
	TIMELIB_AGO                   = 301
	TIMELIB_RELATIVE              = 310
	TIMELIB_ERROR                 = 999
)

// TimeUnit constants for relative units
const (
	TIMELIB_MICROSEC = 1
	TIMELIB_SECOND   = 2
	TIMELIB_MINUTE   = 3
	TIMELIB_HOUR     = 4
	TIMELIB_DAY      = 5
	TIMELIB_MONTH    = 6
	TIMELIB_YEAR     = 7
	TIMELIB_WEEKDAY  = 8
	TIMELIB_SPECIAL  = 9
)

// Error codes
const (
	TIMELIB_ERR_DOUBLE_TIME          = 0x203
	TIMELIB_ERR_DOUBLE_DATE          = 0x204
	TIMELIB_ERR_DOUBLE_TZ            = 0x203
	TIMELIB_ERR_TZID_NOT_FOUND       = 0x202
	TIMELIB_ERR_UNEXPECTED_CHARACTER = 0x205
	TIMELIB_ERR_UNEXPECTED_DATA      = 0x207
	TIMELIB_ERR_NUMBER_OUT_OF_RANGE  = 0x226
	TIMELIB_WARN_DOUBLE_TZ           = 0x101
)
