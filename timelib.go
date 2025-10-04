// Package timelib provides comprehensive date and time parsing and manipulation functionality.
// This is a Go port of the C timelib library originally developed by Derick Rethans.
package timelib

import (
	"errors"
	"math"
	"time"
)

// Constants from the original C library
const (
	TIMELIB_UNSET = -9999999

	// Zone types
	TIMELIB_ZONETYPE_NONE   = 0
	TIMELIB_ZONETYPE_OFFSET = 1
	TIMELIB_ZONETYPE_ABBR   = 2
	TIMELIB_ZONETYPE_ID     = 3

	// Special relative types
	TIMELIB_SPECIAL_WEEKDAY                   = 1
	TIMELIB_SPECIAL_DAY_OF_WEEK_IN_MONTH      = 2
	TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH = 3

	// First/Last day of month
	TIMELIB_SPECIAL_FIRST_DAY_OF_MONTH = 1
	TIMELIB_SPECIAL_LAST_DAY_OF_MONTH  = 2

	// Time unit constants (from C library)
	TIMELIB_SECOND   = 1
	TIMELIB_MINUTE   = 2
	TIMELIB_HOUR     = 3
	TIMELIB_DAY      = 4
	TIMELIB_MONTH    = 5
	TIMELIB_YEAR     = 6
	TIMELIB_WEEKDAY  = 7
	TIMELIB_SPECIAL  = 8
	TIMELIB_MICROSEC = 9

	// Time constants
	MINS_PER_HOUR   = 60
	SECS_PER_DAY    = 86400
	SECS_PER_HOUR   = 3600
	USECS_PER_HOUR  = 3600000000
	DAYS_PER_WEEK   = 7
	DAYS_PER_YEAR   = 365
	DAYS_PER_LYEAR  = 366
	MONTHS_PER_YEAR = 12
)

// Error codes
const (
	TIMELIB_ERROR_NO_ERROR                          = 0x00
	TIMELIB_ERROR_CANNOT_ALLOCATE                   = 0x01
	TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE = 0x02
	TIMELIB_ERROR_CORRUPT_NO_64BIT_PREAMBLE         = 0x03
	TIMELIB_ERROR_CORRUPT_NO_ABBREVIATION           = 0x04
	TIMELIB_ERROR_UNSUPPORTED_VERSION               = 0x05
	TIMELIB_ERROR_NO_SUCH_TIMEZONE                  = 0x06
	TIMELIB_ERROR_SLIM_FILE                         = 0x07
	TIMELIB_ERROR_CORRUPT_POSIX_STRING              = 0x08
	TIMELIB_ERROR_EMPTY_POSIX_STRING                = 0x09
	TIMELIB_ERROR_CANNOT_OPEN_FILE                  = 0x0A

	// Warning codes
	TIMELIB_WARN_DOUBLE_TZ     = 0x101
	TIMELIB_WARN_INVALID_TIME  = 0x102
	TIMELIB_WARN_INVALID_DATE  = 0x103
	TIMELIB_WARN_TRAILING_DATA = 0x11a

	// Parse error codes
	TIMELIB_ERR_DOUBLE_TZ                  = 0x201
	TIMELIB_ERR_TZID_NOT_FOUND             = 0x202
	TIMELIB_ERR_DOUBLE_TIME                = 0x203
	TIMELIB_ERR_DOUBLE_DATE                = 0x204
	TIMELIB_ERR_UNEXPECTED_CHARACTER       = 0x205
	TIMELIB_ERR_EMPTY_STRING               = 0x206
	TIMELIB_ERR_UNEXPECTED_DATA            = 0x207
	TIMELIB_ERR_NO_TEXTUAL_DAY             = 0x208
	TIMELIB_ERR_NO_TWO_DIGIT_DAY           = 0x209
	TIMELIB_ERR_NO_THREE_DIGIT_DAY_OF_YEAR = 0x20a
	TIMELIB_ERR_NO_TWO_DIGIT_MONTH         = 0x20b
	TIMELIB_ERR_NO_TEXTUAL_MONTH           = 0x20c
	TIMELIB_ERR_NO_TWO_DIGIT_YEAR          = 0x20d
	TIMELIB_ERR_NO_FOUR_DIGIT_YEAR         = 0x20e
	TIMELIB_ERR_NO_TWO_DIGIT_HOUR          = 0x20f
	TIMELIB_ERR_HOUR_LARGER_THAN_12        = 0x210
	TIMELIB_ERR_MERIDIAN_BEFORE_HOUR       = 0x211
	TIMELIB_ERR_NO_MERIDIAN                = 0x212
	TIMELIB_ERR_NO_TWO_DIGIT_MINUTE        = 0x213
	TIMELIB_ERR_NO_TWO_DIGIT_SECOND        = 0x214
	TIMELIB_ERR_NO_SIX_DIGIT_MICROSECOND   = 0x215
	TIMELIB_ERR_NO_SEP_SYMBOL              = 0x216
	TIMELIB_ERR_EXPECTED_ESCAPE_CHAR       = 0x217
	TIMELIB_ERR_NO_ESCAPED_CHAR            = 0x218
	TIMELIB_ERR_WRONG_FORMAT_SEP           = 0x219
	TIMELIB_ERR_TRAILING_DATA              = 0x21a
	TIMELIB_ERR_DATA_MISSING               = 0x21b
	TIMELIB_ERR_NO_THREE_DIGIT_MILLISECOND = 0x21c
	TIMELIB_ERR_NO_FOUR_DIGIT_YEAR_ISO     = 0x21d
	TIMELIB_ERR_NO_TWO_DIGIT_WEEK          = 0x21e
	TIMELIB_ERR_INVALID_WEEK               = 0x21f
	TIMELIB_ERR_NO_DAY_OF_WEEK             = 0x220
	TIMELIB_ERR_INVALID_DAY_OF_WEEK        = 0x221
	TIMELIB_ERR_INVALID_SPECIFIER          = 0x222
	TIMELIB_ERR_INVALID_TZ_OFFSET          = 0x223
	TIMELIB_ERR_FORMAT_LITERAL_MISMATCH    = 0x224
	TIMELIB_ERR_MIX_ISO_WITH_NATURAL       = 0x225
	TIMELIB_ERR_NUMBER_OUT_OF_RANGE        = 0x226
)

// Time represents a date/time structure
type Time struct {
	Y, M, D       int64   // Year, Month, Day
	H, I, S       int64   // Hour, Minute, Second
	US            int64   // Microseconds
	Z             int32   // UTC offset in seconds
	TzAbbr        string  // Timezone abbreviation (display only)
	TzInfo        *TzInfo // Timezone structure
	Dst           int     // Flag if we were parsing a DST zone
	Relative      RelTime
	Sse           int64 // Seconds since epoch
	HaveTime      bool
	HaveDate      bool
	HaveZone      bool
	HaveRelative  bool
	HaveWeeknrDay bool
	SseUptodate   bool
	TimUptodate   bool
	IsLocaltime   bool
	ZoneType      int
}

// RelTime represents relative time information
type RelTime struct {
	Y, M, D         int64 // Years, Months and Days
	H, I, S         int64 // Hours, Minutes and Seconds
	US              int64 // Microseconds
	Weekday         int   // Stores the day in 'next monday'
	Week            int   // ISO week number (1-53)
	WeekdayBehavior int   // 0: the current day should *not* be counted when advancing forwards; 1: the current day *should* be counted
	FirstLastDayOf  int
	Invert          bool  // Whether the difference should be inverted
	Days            int64 // Contains the number of *days*, instead of Y-M-D differences
	Special         struct {
		Type   int
		Amount int64
	}
	HaveWeekdayRelative bool
	HaveSpecialRelative bool
}

// TimeOffset represents timezone offset information
type TimeOffset struct {
	Offset         int32
	LeapSecs       int
	IsDst          int
	Abbr           string
	TransitionTime int64
}

// TzInfo represents timezone information
type TzInfo struct {
	Name  string
	Bit32 struct {
		Ttisgmtcnt uint32
		Ttisstdcnt uint32
		Leapcnt    uint32
		Timecnt    uint32
		Typecnt    uint32
		Charcnt    uint32
	}
	Bit64 struct {
		Ttisgmtcnt uint64
		Ttisstdcnt uint64
		Leapcnt    uint64
		Timecnt    uint64
		Typecnt    uint64
		Charcnt    uint64
	}
	Trans        []int64
	TransIdx     []uint8
	Type         []TTInfo
	TimezoneAbbr string
	LeapTimes    []TLInfo
	Bc           uint8
	Location     TLocInfo
	PosixString  string
	PosixInfo    *PosixStr
}

// TTInfo represents timezone type information
type TTInfo struct {
	Offset  int32
	IsDst   int
	AbbrIdx int
	IsStd   int
	IsUtc   int
}

// TLInfo represents leap time information
type TLInfo struct {
	Trans int64
	Corr  int64
}

// TLocInfo represents location information
type TLocInfo struct {
	CountryCode [3]byte
	Latitude    float64
	Longitude   float64
	Comments    string
}

// PosixStr represents POSIX timezone string information
type PosixStr struct {
	Std              string
	StdOffset        int64
	Dst              string
	DstOffset        int64
	DstBegin         *PosixTransInfo
	DstEnd           *PosixTransInfo
	TypeIndexStdType int
	TypeIndexDstType int
}

// PosixTransInfo represents POSIX transition information
type PosixTransInfo struct {
	Type int // 1=Jn, 2=n, 3=Mm.w.d
	Days int
	Mwd  struct {
		Month int
		Week  int
		Dow   int
	}
	Hour int
}

// PosixTransitions represents POSIX transitions
type PosixTransitions struct {
	Count int
	Times [6]int64
	Types [6]int64
}

// ErrorMessage represents an error message
type ErrorMessage struct {
	ErrorCode int
	Position  int
	Character byte
	Message   string
}

// ErrorContainer holds error and warning messages
type ErrorContainer struct {
	ErrorMessages   []ErrorMessage
	WarningMessages []ErrorMessage
	ErrorCount      int
	WarningCount    int
}

// TzLookupTable represents timezone lookup table entry
type TzLookupTable struct {
	Name       string
	Type       int
	GmtOffset  float32
	FullTzName string
}

// TzDBIndexEntry represents timezone database index entry
type TzDBIndexEntry struct {
	ID  string
	Pos int
}

// TzDB represents timezone database
type TzDB struct {
	Version   string
	IndexSize int
	Index     []TzDBIndexEntry
	Data      []byte
	BaseDir   string // Base directory for file-based databases
}

// FormatSpecifier represents a format specifier
type FormatSpecifier struct {
	Specifier byte
	Code      FormatSpecifierCode
}

// FormatSpecifierCode represents format specifier codes
type FormatSpecifierCode int

const (
	TIMELIB_FORMAT_ALLOW_EXTRA_CHARACTERS FormatSpecifierCode = iota
	TIMELIB_FORMAT_ANY_SEPARATOR
	TIMELIB_FORMAT_DAY_TWO_DIGIT
	TIMELIB_FORMAT_DAY_TWO_DIGIT_PADDED
	TIMELIB_FORMAT_DAY_OF_WEEK_ISO
	TIMELIB_FORMAT_DAY_OF_WEEK
	TIMELIB_FORMAT_DAY_OF_YEAR
	TIMELIB_FORMAT_DAY_SUFFIX
	TIMELIB_FORMAT_END
	TIMELIB_FORMAT_EPOCH_SECONDS
	TIMELIB_FORMAT_ESCAPE
	TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX
	TIMELIB_FORMAT_HOUR_TWO_DIGIT_12_MAX_PADDED
	TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX
	TIMELIB_FORMAT_HOUR_TWO_DIGIT_24_MAX_PADDED
	TIMELIB_FORMAT_LITERAL
	TIMELIB_FORMAT_MERIDIAN
	TIMELIB_FORMAT_MICROSECOND_SIX_DIGIT
	TIMELIB_FORMAT_MILLISECOND_THREE_DIGIT
	TIMELIB_FORMAT_MINUTE_TWO_DIGIT
	TIMELIB_FORMAT_MONTH_TWO_DIGIT
	TIMELIB_FORMAT_MONTH_TWO_DIGIT_PADDED
	TIMELIB_FORMAT_RANDOM_CHAR
	TIMELIB_FORMAT_RESET_ALL
	TIMELIB_FORMAT_RESET_ALL_WHEN_NOT_SET
	TIMELIB_FORMAT_SECOND_TWO_DIGIT
	TIMELIB_FORMAT_SEPARATOR
	TIMELIB_FORMAT_SKIP_TO_SEPARATOR
	TIMELIB_FORMAT_TEXTUAL_DAY_3_LETTER
	TIMELIB_FORMAT_TEXTUAL_DAY_FULL
	TIMELIB_FORMAT_TEXTUAL_MONTH_3_LETTER
	TIMELIB_FORMAT_TEXTUAL_MONTH_FULL
	TIMELIB_FORMAT_TIMEZONE_OFFSET
	TIMELIB_FORMAT_TIMEZONE_OFFSET_MINUTES
	TIMELIB_FORMAT_WEEK_OF_YEAR_ISO
	TIMELIB_FORMAT_WEEK_OF_YEAR
	TIMELIB_FORMAT_WHITESPACE
	TIMELIB_FORMAT_YEAR_TWO_DIGIT
	TIMELIB_FORMAT_YEAR_FOUR_DIGIT
	TIMELIB_FORMAT_YEAR_EXPANDED
	TIMELIB_FORMAT_YEAR_ISO
)

// FormatConfig represents format configuration
type FormatConfig struct {
	FormatMap            []FormatSpecifier
	PrefixChar           byte
	AllowExtraCharacters bool
}

// Common errors
var (
	ErrInvalidDate     = errors.New("invalid date")
	ErrInvalidTime     = errors.New("invalid time")
	ErrInvalidTimezone = errors.New("invalid timezone")
	ErrParseError      = errors.New("parse error")
)

// GetErrorMessage returns a static string containing an error message belonging to a specific error code
func GetErrorMessage(errorCode int) string {
	errorMessages := []string{
		"No error",
		"Cannot allocate buffer for parsing",
		"Corrupt tzfile: The transitions in the file don't always increase",
		"Corrupt tzfile: The expected 64-bit preamble is missing",
		"Corrupt tzfile: No abbreviation could be found for a transition",
		"The version used in this timezone identifier is unsupported",
		"No timezone with this name could be found",
		"A 'slim' timezone file has been detected",
		"The embedded POSIX string is not valid",
		"The embedded POSIX string is empty",
	}

	if errorCode >= 0 && errorCode < len(errorMessages) {
		return errorMessages[errorCode]
	}
	return "Unknown error code"
}

// TimeCtor creates a new Time structure
func TimeCtor() *Time {
	return &Time{
		Y:        TIMELIB_UNSET,
		M:        TIMELIB_UNSET,
		D:        TIMELIB_UNSET,
		H:        TIMELIB_UNSET,
		I:        TIMELIB_UNSET,
		S:        TIMELIB_UNSET,
		US:       0,
		Z:        0,
		Dst:      0, // Changed from -1 to match C calloc behavior (zero-initialized)
		ZoneType: TIMELIB_ZONETYPE_NONE,
	}
}

// RelTimeCtor creates a new RelTime structure
func RelTimeCtor() *RelTime {
	return &RelTime{
		Y: 0, M: 0, D: 0,
		H: 0, I: 0, S: 0,
		US:              0,
		Weekday:         -1,
		Week:            -1,
		WeekdayBehavior: 0,
		FirstLastDayOf:  0,
		Invert:          false,
		Days:            0,
	}
}

// TimeOffsetCtor creates a new TimeOffset structure
func TimeOffsetCtor() *TimeOffset {
	return &TimeOffset{
		Offset:         0,
		LeapSecs:       0,
		IsDst:          0,
		Abbr:           "",
		TransitionTime: 0,
	}
}

// ErrorContainerCtor creates a new ErrorContainer
func ErrorContainerCtor() *ErrorContainer {
	return &ErrorContainer{
		ErrorMessages:   make([]ErrorMessage, 0),
		WarningMessages: make([]ErrorMessage, 0),
		ErrorCount:      0,
		WarningCount:    0,
	}
}

// TimeCompare compares two Time structures
func TimeCompare(t1, t2 *Time) int {
	// If both times have valid SSE, compare by SSE first
	if t1.SseUptodate && t2.SseUptodate {
		if t1.Sse == t2.Sse {
			if t1.US == t2.US {
				return 0
			}
			if t1.US < t2.US {
				return -1
			}
			return 1
		}
		if t1.Sse < t2.Sse {
			return -1
		}
		return 1
	}

	// Fall back to comparing individual date/time fields
	// Compare years
	if t1.Y != t2.Y {
		if t1.Y < t2.Y {
			return -1
		}
		return 1
	}

	// Compare months
	if t1.M != t2.M {
		if t1.M < t2.M {
			return -1
		}
		return 1
	}

	// Compare days
	if t1.D != t2.D {
		if t1.D < t2.D {
			return -1
		}
		return 1
	}

	// Compare hours
	if t1.H != t2.H {
		if t1.H < t2.H {
			return -1
		}
		return 1
	}

	// Compare minutes
	if t1.I != t2.I {
		if t1.I < t2.I {
			return -1
		}
		return 1
	}

	// Compare seconds
	if t1.S != t2.S {
		if t1.S < t2.S {
			return -1
		}
		return 1
	}

	// Compare microseconds
	if t1.US != t2.US {
		if t1.US < t2.US {
			return -1
		}
		return 1
	}

	return 0
}

// DecimalHourToHMS converts a decimal hour into hour/min/sec components
func DecimalHourToHMS(h float64) (hour, min, sec int) {
	// Matches C function: timelib_decimal_hour_to_hms in timelib.c
	swap := false

	if h < 0 {
		swap = true
		h = -h
	}

	hour = int(math.Floor(h))
	seconds := int(math.Floor((h - float64(hour)) * 3600.0))

	min = seconds / 60
	sec = seconds % 60

	if swap {
		hour = -hour
	}

	return
}

// HMSToDecimalHour converts hour/min/sec values into a decimal hour
func HMSToDecimalHour(hour, min, sec int) float64 {
	if hour >= 0 {
		return float64(hour) + float64(min)/60 + float64(sec)/3600
	}
	return float64(hour) - float64(min)/60 - float64(sec)/3600
}

// HMSFToDecimalHour converts hour/min/sec/micro sec values into a decimal hour
func HMSFToDecimalHour(hour, min, sec, us int) float64 {
	if hour >= 0 {
		return float64(hour) + float64(min)/60.0 + float64(sec)/3600.0 + float64(us)/3600000000.0
	}
	return float64(hour) - float64(min)/60.0 - float64(sec)/3600.0 - float64(us)/3600000000.0
}

// HMSToSeconds converts hour/min/sec values into seconds
func HMSToSeconds(h, m, s int64) int64 {
	return h*3600 + m*60 + s
}

// DateToInt converts the 'sse' value to an int64 type with error checking
func DateToInt(d *Time) (int64, error) {
	ts := d.Sse

	if ts < -9223372036854775808 || ts > 9223372036854775807 {
		return 0, errors.New("timestamp out of range")
	}

	return ts, nil
}

// SetTimezoneFromOffset attaches the UTC offset as time zone information
func SetTimezoneFromOffset(t *Time, utcOffset int64) {
	t.ZoneType = TIMELIB_ZONETYPE_OFFSET
	t.Z = int32(utcOffset)
	t.Dst = 0
	t.TzAbbr = ""
	t.TzInfo = nil
}

// SetTimezoneFromAbbr attaches timezone information from abbreviation
func SetTimezoneFromAbbr(t *Time, abbr string, utcOffset int64, isDst int) {
	t.ZoneType = TIMELIB_ZONETYPE_ABBR
	t.Z = int32(utcOffset)
	t.Dst = isDst
	t.TzAbbr = abbr
	t.TzInfo = nil
}

// SetTimezone attaches timezone information from TzInfo
func SetTimezone(t *Time, tz *TzInfo) {
	gmtOffset := GetTimeZoneInfo(t.Sse, tz)
	if gmtOffset != nil {
		t.Z = gmtOffset.Offset
		t.Dst = int(gmtOffset.IsDst)
		t.TzAbbr = gmtOffset.Abbr
	} else {
		// No timezone data available, use defaults
		t.Z = 0
		t.Dst = 0
		t.TzAbbr = tz.Name
	}
	t.TzInfo = tz
	t.HaveZone = true
	t.ZoneType = TIMELIB_ZONETYPE_ID
}

// ConvertTime converts a Time structure to Go's time.Time
func ConvertTime(t *Time) time.Time {
	// For now, return zero time - this will be implemented later
	return time.Time{}
}

// Helper functions and lookup tables for date calculations

// Lookup tables for day of week calculations
var mTableCommon = [13]int{-1, 0, 3, 3, 6, 1, 4, 6, 2, 5, 0, 3, 5} // 1 = jan
var mTableLeap = [13]int{-1, 6, 2, 3, 6, 1, 4, 6, 2, 5, 0, 3, 5}   // 1 = jan

// Lookup tables for day of year calculations
var dTableCommon = [13]int{0, 0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
var dTableLeap = [13]int{0, 0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335}

// Lookup tables for days in month
var mlTableCommon = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
var mlTableLeap = [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// positiveMod returns positive modulo result
func positiveMod(x, y int64) int64 {
	if y == 0 {
		return 0
	}
	tmp := x % y
	if tmp < 0 {
		if y > 0 {
			tmp += y
		} else {
			tmp -= y // For negative divisors, subtract to get positive result
		}
	}
	return tmp
}

// centuryValue calculates century value for day of week calculation
func centuryValue(j int64) int64 {
	return 6 - positiveMod(j, 4)*2
}

// IsLeapYear determines if a year is a leap year
func IsLeapYear(y int64) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// DayOfWeekEx calculates day of week with ISO option
func DayOfWeekEx(y, m, d int64, iso bool) int64 {
	// Only valid for Gregorian calendar
	c1 := centuryValue(positiveMod(y, 400) / 100)
	y1 := positiveMod(y, 100)

	var m1 int64
	if IsLeapYear(y) {
		m1 = int64(mTableLeap[m])
	} else {
		m1 = int64(mTableCommon[m])
	}

	dow := positiveMod((c1 + y1 + m1 + (y1 / 4) + d), 7)

	if iso {
		if dow == 0 {
			dow = 7
		}
	}
	return dow
}

// DayOfWeek calculates day of week (0=Sunday..6=Saturday)
func DayOfWeek(y, m, d int64) int64 {
	return DayOfWeekEx(y, m, d, false)
}

// IsoDayOfWeek calculates ISO day of week (1=Monday, 7=Sunday)
func IsoDayOfWeek(y, m, d int64) int64 {
	return DayOfWeekEx(y, m, d, true)
}

// DayOfYear calculates day of year according to y-m-d (0=Jan 1st..364/365=Dec 31st)
func DayOfYear(y, m, d int64) int64 {
	if IsLeapYear(y) {
		return int64(dTableLeap[m]) + d - 1
	}
	return int64(dTableCommon[m]) + d - 1
}

// DaysInMonth calculates number of days in month m for year y
func DaysInMonth(y, m int64) int64 {
	if IsLeapYear(y) {
		return int64(mlTableLeap[m])
	}
	return int64(mlTableCommon[m])
}

// ValidTime checks if h, i, s fit in valid range
func ValidTime(h, i, s int64) bool {
	return h >= 0 && h <= 23 && i >= 0 && i <= 59 && s >= 0 && s <= 59
}

// ValidDate checks if y, m, d form a valid date
func ValidDate(y, m, d int64) bool {
	if m < 1 || m > 12 || d < 1 {
		return false
	}
	return d <= DaysInMonth(y, m)
}

// IsoWeekFromDate calculates ISO week from date
func IsoWeekFromDate(y, m, d int64) (iw, iy int64) {
	yLeap := IsLeapYear(y)
	prevYLeap := IsLeapYear(y - 1)
	doy := DayOfYear(y, m, d) + 1
	if yLeap && m > 2 {
		doy++
	}
	jan1weekday := DayOfWeek(y, 1, 1)
	weekday := DayOfWeek(y, m, d)
	if weekday == 0 {
		weekday = 7
	}
	if jan1weekday == 0 {
		jan1weekday = 7
	}

	// Find if Y M D falls in YearNumber Y-1, WeekNumber 52 or 53
	if doy <= (8-jan1weekday) && jan1weekday > 4 {
		iy = y - 1
		if jan1weekday == 5 || (jan1weekday == 6 && prevYLeap) {
			iw = 53
		} else {
			iw = 52
		}
	} else {
		iy = y
	}

	// Find if Y M D falls in YearNumber Y+1, WeekNumber 1
	if iy == y {
		daysInYear := int64(365)
		if yLeap {
			daysInYear = 366
		}
		yLeapInt := int64(0)
		if yLeap {
			yLeapInt = 1
		}
		if (daysInYear - (doy - yLeapInt)) < (4 - weekday) {
			iy = y + 1
			iw = 1
			return
		}
	}

	// Find if Y M D falls in YearNumber Y, WeekNumber 1 through 53
	if iy == y {
		j := doy + (7 - weekday) + (jan1weekday - 1)
		iw = j / 7
		if jan1weekday > 4 {
			iw -= 1
		}
	}

	return
}

// IsoDateFromDate calculates ISO date from date
func IsoDateFromDate(y, m, d int64) (iy, iw, id int64) {
	iw, iy = IsoWeekFromDate(y, m, d)
	id = IsoDayOfWeek(y, m, d)
	return
}

// DayNrFromWeekNr calculates day number from week number
func DayNrFromWeekNr(iy, iw, id int64) int64 {
	// Figure out the dayofweek for y-1-1
	dow := DayOfWeek(iy, 1, 1)
	// then use that to figure out the offset for day 1 of week 1
	var day int64
	if dow > 4 {
		day = 0 - (dow - 7)
	} else {
		day = 0 - dow
	}

	// Add weeks and days
	return day + ((iw - 1) * 7) + id
}

// DateFromIsoDate calculates date from ISO date
func DateFromIsoDate(iy, iw, id int64) (y, m, d int64) {
	daynr := DayNrFromWeekNr(iy, iw, id) + 1
	var table []int64
	var isLeapYear bool

	// Invariant: isLeapYear == IsLeapYear(*y)
	y = iy
	isLeapYear = IsLeapYear(y)

	// Establish invariant that daynr >= 0
	for daynr <= 0 {
		y -= 1
		isLeapYear = IsLeapYear(y)
		if isLeapYear {
			daynr += 366
		} else {
			daynr += 365
		}
	}

	// Establish invariant that daynr <= number of days in *yr
	for (isLeapYear && daynr > 366) || (!isLeapYear && daynr > 365) {
		if isLeapYear {
			daynr -= 366
		} else {
			daynr -= 365
		}
		y += 1
		isLeapYear = IsLeapYear(y)
	}

	if isLeapYear {
		table = make([]int64, 13)
		for i := 0; i < 13; i++ {
			table[i] = int64(mlTableLeap[i])
		}
	} else {
		table = make([]int64, 13)
		for i := 0; i < 13; i++ {
			table[i] = int64(mlTableCommon[i])
		}
	}

	// Establish invariant that daynr <= number of days in *m
	m = 1
	for daynr > table[m] {
		daynr -= table[m]
		m += 1
	}

	d = daynr
	return
}

// ConvertFromTime converts Go's time.Time to a Time structure
func ConvertFromTime(t time.Time) *Time {
	// For now, return basic structure - this will be implemented later
	return TimeCtor()
}
