/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2015-2023 Derick Rethans
 * Copyright (c) 2018 MongoDB, Inc.
 * Copyright (c) 2025 Go Port
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package timelib

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// Token constants and error codes are defined in parse_date_constants.go and timelib.go

const (
	BSIZE = 8192
	MAX_ABBR_LEN = 6
)

const (
	TIMELIB_TIME_PART_KEEP = 0
	TIMELIB_TIME_PART_DONT_KEEP = 1
)

// TimeUnit constants for relative units - also in timelib.go

// Note: TIMELIB_SPECIAL_* constants are defined in timelib.go

// Scanner represents the state of the parser
type Scanner struct {
	fd      int
	lim     *byte
	str     []byte
	ptr     *byte
	cur     *byte
	tok     *byte
	pos     *byte
	line    uint
	len     uint
	errors  *ErrorContainer
	time    *Time
	tzdb    *TzDB
}

// LookupTable represents a generic lookup table
type LookupTable struct {
	Name  string
	Type  int
	Value int
}

// RelUnit represents a relative unit definition
type RelUnit struct {
	Name       string
	Unit       int
	Multiplier int
}

// Note: TzLookupTable is defined in timelib.go

// TzGetWrapper is a function type for timezone lookup
type TzGetWrapper func(string, *TzDB, *int) (*TzInfo, error)

// Timezone lookup tables
var timelibTimezoneUTC = []TzLookupTable{
	{"utc", 0, 0, "UTC"},
}

// Relative unit lookup table
var timelibRelunitLookup = []RelUnit{
	{"ms", TIMELIB_MICROSEC, 1000},
	{"msec", TIMELIB_MICROSEC, 1000},
	{"msecs", TIMELIB_MICROSEC, 1000},
	{"millisecond", TIMELIB_MICROSEC, 1000},
	{"milliseconds", TIMELIB_MICROSEC, 1000},
	{"µs", TIMELIB_MICROSEC, 1},
	{"usec", TIMELIB_MICROSEC, 1},
	{"usecs", TIMELIB_MICROSEC, 1},
	{"µsec", TIMELIB_MICROSEC, 1},
	{"µsecs", TIMELIB_MICROSEC, 1},
	{"microsecond", TIMELIB_MICROSEC, 1},
	{"microseconds", TIMELIB_MICROSEC, 1},
	{"sec", TIMELIB_SECOND, 1},
	{"secs", TIMELIB_SECOND, 1},
	{"second", TIMELIB_SECOND, 1},
	{"seconds", TIMELIB_SECOND, 1},
	{"min", TIMELIB_MINUTE, 1},
	{"mins", TIMELIB_MINUTE, 1},
	{"minute", TIMELIB_MINUTE, 1},
	{"minutes", TIMELIB_MINUTE, 1},
	{"hour", TIMELIB_HOUR, 1},
	{"hours", TIMELIB_HOUR, 1},
	{"day", TIMELIB_DAY, 1},
	{"days", TIMELIB_DAY, 1},
	{"week", TIMELIB_DAY, 7},
	{"weeks", TIMELIB_DAY, 7},
	{"fortnight", TIMELIB_DAY, 14},
	{"fortnights", TIMELIB_DAY, 14},
	{"forthnight", TIMELIB_DAY, 14},
	{"forthnights", TIMELIB_DAY, 14},
	{"month", TIMELIB_MONTH, 1},
	{"months", TIMELIB_MONTH, 1},
	{"year", TIMELIB_YEAR, 1},
	{"years", TIMELIB_YEAR, 1},

	{"mondays", TIMELIB_WEEKDAY, 1},
	{"monday", TIMELIB_WEEKDAY, 1},
	{"mon", TIMELIB_WEEKDAY, 1},
	{"tuesdays", TIMELIB_WEEKDAY, 2},
	{"tuesday", TIMELIB_WEEKDAY, 2},
	{"tue", TIMELIB_WEEKDAY, 2},
	{"wednesdays", TIMELIB_WEEKDAY, 3},
	{"wednesday", TIMELIB_WEEKDAY, 3},
	{"wed", TIMELIB_WEEKDAY, 3},
	{"thursdays", TIMELIB_WEEKDAY, 4},
	{"thursday", TIMELIB_WEEKDAY, 4},
	{"thu", TIMELIB_WEEKDAY, 4},
	{"fridays", TIMELIB_WEEKDAY, 5},
	{"friday", TIMELIB_WEEKDAY, 5},
	{"fri", TIMELIB_WEEKDAY, 5},
	{"saturdays", TIMELIB_WEEKDAY, 6},
	{"saturday", TIMELIB_WEEKDAY, 6},
	{"sat", TIMELIB_WEEKDAY, 6},
	{"sundays", TIMELIB_WEEKDAY, 0},
	{"sunday", TIMELIB_WEEKDAY, 0},
	{"sun", TIMELIB_WEEKDAY, 0},

	{"weekday", TIMELIB_SPECIAL, TIMELIB_SPECIAL_WEEKDAY},
	{"weekdays", TIMELIB_SPECIAL, TIMELIB_SPECIAL_WEEKDAY},
}

// Relative text lookup table
var timelibReltextLookup = []LookupTable{
	{"first", 0, 1},
	{"next", 0, 1},
	{"second", 0, 2},
	{"third", 0, 3},
	{"fourth", 0, 4},
	{"fifth", 0, 5},
	{"sixth", 0, 6},
	{"seventh", 0, 7},
	{"eight", 0, 8},
	{"eighth", 0, 8},
	{"ninth", 0, 9},
	{"tenth", 0, 10},
	{"eleventh", 0, 11},
	{"twelfth", 0, 12},
	{"last", 0, -1},
	{"previous", 0, -1},
	{"this", 1, 0},
}

// Month lookup table
var timelibMonthLookup = []LookupTable{
	{"jan", 0, 1},
	{"feb", 0, 2},
	{"mar", 0, 3},
	{"apr", 0, 4},
	{"may", 0, 5},
	{"jun", 0, 6},
	{"jul", 0, 7},
	{"aug", 0, 8},
	{"sep", 0, 9},
	{"sept", 0, 9},
	{"oct", 0, 10},
	{"nov", 0, 11},
	{"dec", 0, 12},
	{"i", 0, 1},
	{"ii", 0, 2},
	{"iii", 0, 3},
	{"iv", 0, 4},
	{"v", 0, 5},
	{"vi", 0, 6},
	{"vii", 0, 7},
	{"viii", 0, 8},
	{"ix", 0, 9},
	{"x", 0, 10},
	{"xi", 0, 11},
	{"xii", 0, 12},

	{"january", 0, 1},
	{"february", 0, 2},
	{"march", 0, 3},
	{"april", 0, 4},
	{"may", 0, 5},
	{"june", 0, 6},
	{"july", 0, 7},
	{"august", 0, 8},
	{"september", 0, 9},
	{"october", 0, 10},
	{"november", 0, 11},
	{"december", 0, 12},
}

// Helper functions

func allocErrorMessage(messages *[]ErrorMessage, count *int) *ErrorMessage {
	isPow2 := (*count & (*count - 1)) == 0

	if isPow2 {
		allocSize := 1
		if *count > 0 {
			allocSize = *count * 2
		}
		newMessages := make([]ErrorMessage, allocSize)
		if *messages != nil {
			copy(newMessages, *messages)
		}
		*messages = newMessages
	}
	msg := &(*messages)[*count]
	*count++
	return msg
}

func addWarning(s *Scanner, errorCode int, errorMsg string) {
	message := allocErrorMessage(&s.errors.WarningMessages, &s.errors.WarningCount)
	message.ErrorCode = errorCode
	if s.tok != nil {
		// Calculate position safely using byte pointer arithmetic
		message.Position = int(uintptr(unsafe.Pointer(s.tok)) - uintptr(unsafe.Pointer(&s.str[0])))
		message.Character = *s.tok
	} else {
		message.Position = 0
		message.Character = 0
	}
	message.Message = errorMsg
}

func addError(s *Scanner, errorCode int, errorMsg string) {
	message := allocErrorMessage(&s.errors.ErrorMessages, &s.errors.ErrorCount)
	message.ErrorCode = errorCode
	if s.tok != nil {
		// Calculate position safely using byte pointer arithmetic
		message.Position = int(uintptr(unsafe.Pointer(s.tok)) - uintptr(unsafe.Pointer(&s.str[0])))
		message.Character = *s.tok
	} else {
		message.Position = 0
		message.Character = 0
	}
	message.Message = errorMsg
}

func addPbfWarning(s *Scanner, errorCode int, errorMsg string, sptr, cptr string) {
	message := allocErrorMessage(&s.errors.WarningMessages, &s.errors.WarningCount)
	message.ErrorCode = errorCode
	if len(cptr) > 0 && len(sptr) > 0 {
		message.Position = len(sptr) - len(cptr)
		message.Character = byte(cptr[0])
	}
	message.Message = errorMsg
}

func addPbfError(s *Scanner, errorCode int, errorMsg string, sptr, cptr string) {
	message := allocErrorMessage(&s.errors.ErrorMessages, &s.errors.ErrorCount)
	message.ErrorCode = errorCode
	if len(cptr) > 0 && len(sptr) > 0 {
		message.Position = len(sptr) - len(cptr)
		message.Character = byte(cptr[0])
	}
	message.Message = errorMsg
}

func timelibMeridian(ptr *string, h int64) int64 {
	var retval int64 = 0

	for len(*ptr) > 0 && !strings.ContainsRune("AaPp", rune((*ptr)[0])) {
		*ptr = (*ptr)[1:]
	}
	if len(*ptr) == 0 {
		return retval
	}

	if (*ptr)[0] == 'a' || (*ptr)[0] == 'A' {
		if h == 12 {
			retval = -12
		}
	} else if h != 12 {
		retval = 12
	}
	*ptr = (*ptr)[1:]
	if len(*ptr) > 0 && (*ptr)[0] == '.' {
		*ptr = (*ptr)[1:]
	}
	if len(*ptr) > 0 && ((*ptr)[0] == 'M' || (*ptr)[0] == 'm') {
		*ptr = (*ptr)[1:]
	}
	if len(*ptr) > 0 && (*ptr)[0] == '.' {
		*ptr = (*ptr)[1:]
	}
	return retval
}

func timelibMeridianWithCheck(ptr *string, h int64) int64 {
	var retval int64 = 0

	for len(*ptr) > 0 && !strings.ContainsRune("AaPp", rune((*ptr)[0])) {
		*ptr = (*ptr)[1:]
	}
	if len(*ptr) == 0 {
		return TIMELIB_UNSET
	}

	if (*ptr)[0] == 'a' || (*ptr)[0] == 'A' {
		if h == 12 {
			retval = -12
		}
	} else if h != 12 {
		retval = 12
	}
	*ptr = (*ptr)[1:]
	if len(*ptr) > 0 && (*ptr)[0] == '.' {
		*ptr = (*ptr)[1:]
		if len(*ptr) == 0 || ((*ptr)[0] != 'm' && (*ptr)[0] != 'M') {
			return TIMELIB_UNSET
		}
		*ptr = (*ptr)[1:]
		if len(*ptr) == 0 || (*ptr)[0] != '.' {
			return TIMELIB_UNSET
		}
		*ptr = (*ptr)[1:]
	} else if len(*ptr) > 0 && ((*ptr)[0] == 'm' || (*ptr)[0] == 'M') {
		*ptr = (*ptr)[1:]
	} else {
		return TIMELIB_UNSET
	}
	return retval
}

func timelibString(s *Scanner) string {
	if s.tok == nil || s.cur == nil {
		return ""
	}

	// Calculate positions safely using byte pointer arithmetic
	tokPos := int(uintptr(unsafe.Pointer(s.tok)) - uintptr(unsafe.Pointer(&s.str[0])))
	curPos := int(uintptr(unsafe.Pointer(s.cur)) - uintptr(unsafe.Pointer(&s.str[0])))

	// Additional bounds checking
	if tokPos < 0 || curPos > len(s.str) || tokPos > curPos {
		return ""
	}

	return string(s.str[tokPos:curPos])
}

func timelibGetNrEx(ptr *string, maxLength int, scannedLength *int) int64 {
	if len(*ptr) == 0 {
		return TIMELIB_UNSET
	}

	for len(*ptr) > 0 && ((*ptr)[0] < '0' || (*ptr)[0] > '9') {
		if (*ptr)[0] == '\x00' {
			return TIMELIB_UNSET
		}
		*ptr = (*ptr)[1:]
	}

	begin := *ptr
	length := 0

	for len(*ptr) > 0 && (*ptr)[0] >= '0' && (*ptr)[0] <= '9' && length < maxLength {
		*ptr = (*ptr)[1:]
		length++
	}

	if scannedLength != nil {
		*scannedLength = length
	}

	if length == 0 {
		return TIMELIB_UNSET
	}

	numStr := begin[:length]
	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return TIMELIB_UNSET
	}
	return val
}

func timelibGetNr(ptr *string, maxLength int) int64 {
	return timelibGetNrEx(ptr, maxLength, nil)
}

func timelibSkipDaySuffix(ptr *string) {
	if len(*ptr) == 0 {
		return
	}
	if unicode.IsSpace(rune((*ptr)[0])) {
		return
	}
	if len(*ptr) >= 2 {
		suffix := strings.ToLower((*ptr)[:2])
		if suffix == "nd" || suffix == "rd" || suffix == "st" || suffix == "th" {
			*ptr = (*ptr)[2:]
		}
	}
}

func timelibGetFracNr(ptr *string) int64 {
	if len(*ptr) == 0 {
		return TIMELIB_UNSET
	}

	for len(*ptr) > 0 && (*ptr)[0] != '.' && (*ptr)[0] != ':' && ((*ptr)[0] < '0' || (*ptr)[0] > '9') {
		if (*ptr)[0] == '\x00' {
			return TIMELIB_UNSET
		}
		*ptr = (*ptr)[1:]
	}

	if len(*ptr) == 0 {
		return TIMELIB_UNSET
	}

	begin := *ptr

	for len(*ptr) > 0 && ((*ptr)[0] == '.' || (*ptr)[0] == ':' || ((*ptr)[0] >= '0' && (*ptr)[0] <= '9')) {
		*ptr = (*ptr)[1:]
	}

	if *ptr == begin {
		return TIMELIB_UNSET
	}

	fracLen := len(begin) - len(*ptr)
	if fracLen <= 1 {
		return TIMELIB_UNSET
	}

	numStr := begin[1:fracLen]
	// Parse as integer, not as "0.xxx" fraction
	// C code: strtod(str, NULL) where str is just the digits "123456"
	// This gives 123456.0, not 0.123456
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return TIMELIB_UNSET
	}

	return int64(val * math.Pow(10, 7-float64(fracLen)))
}

func timelibGetSignedNr(s *Scanner, ptr *string, maxLength int) int64 {
	if len(*ptr) == 0 {
		addError(s, TIMELIB_ERR_UNEXPECTED_DATA, "Found unexpected data")
		return 0
	}

	for len(*ptr) > 0 && ((*ptr)[0] < '0' || (*ptr)[0] > '9') && (*ptr)[0] != '+' && (*ptr)[0] != '-' {
		if (*ptr)[0] == '\x00' {
			addError(s, TIMELIB_ERR_UNEXPECTED_DATA, "Found unexpected data")
			return 0
		}
		*ptr = (*ptr)[1:]
	}

	if len(*ptr) == 0 {
		addError(s, TIMELIB_ERR_UNEXPECTED_DATA, "Found unexpected data")
		return 0
	}

	sign := int64(1)

	for len(*ptr) > 0 && ((*ptr)[0] == '+' || (*ptr)[0] == '-') {
		if (*ptr)[0] == '-' {
			sign = -sign
		}
		*ptr = (*ptr)[1:]
	}

	for len(*ptr) > 0 && ((*ptr)[0] < '0' || (*ptr)[0] > '9') {
		if (*ptr)[0] == '\x00' {
			addError(s, TIMELIB_ERR_UNEXPECTED_DATA, "Found unexpected data")
			return 0
		}
		*ptr = (*ptr)[1:]
	}

	if len(*ptr) == 0 {
		addError(s, TIMELIB_ERR_UNEXPECTED_DATA, "Found unexpected data")
		return 0
	}

	begin := *ptr
	length := 0

	for len(*ptr) > 0 && (*ptr)[0] >= '0' && (*ptr)[0] <= '9' && length < maxLength {
		*ptr = (*ptr)[1:]
		length++
	}

	if length == 0 {
		return 0
	}

	numStr := begin[:length]
	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		addError(s, TIMELIB_ERR_NUMBER_OUT_OF_RANGE, "Number out of range")
		return 0
	}

	return sign * val
}

func timelibLookupRelativeText(ptr *string, behavior *int) int64 {
	if len(*ptr) == 0 {
		return 0
	}

	begin := *ptr

	for len(*ptr) > 0 && (((*ptr)[0] >= 'A' && (*ptr)[0] <= 'Z') || ((*ptr)[0] >= 'a' && (*ptr)[0] <= 'z')) {
		*ptr = (*ptr)[1:]
	}

	if *ptr == begin {
		return 0
	}

	word := strings.ToLower(begin[:len(begin)-len(*ptr)])

	for _, tp := range timelibReltextLookup {
		if strings.EqualFold(word, tp.Name) {
			*behavior = tp.Type
			return int64(tp.Value)
		}
	}

	return 0
}

func timelibGetRelativeText(ptr *string, behavior *int) int64 {
	for len(*ptr) > 0 && ((*ptr)[0] == ' ' || (*ptr)[0] == '\t' || (*ptr)[0] == '-' || (*ptr)[0] == '/') {
		*ptr = (*ptr)[1:]
	}
	return timelibLookupRelativeText(ptr, behavior)
}

func timelibLookupMonth(ptr *string) int64 {
	if len(*ptr) == 0 {
		return 0
	}

	begin := *ptr

	for len(*ptr) > 0 && (((*ptr)[0] >= 'A' && (*ptr)[0] <= 'Z') || ((*ptr)[0] >= 'a' && (*ptr)[0] <= 'z')) {
		*ptr = (*ptr)[1:]
	}

	if *ptr == begin {
		return 0
	}

	word := strings.ToLower(begin[:len(begin)-len(*ptr)])

	for _, tp := range timelibMonthLookup {
		if strings.EqualFold(word, tp.Name) {
			return int64(tp.Value)
		}
	}

	return 0
}

func timelibGetMonth(ptr *string) int64 {
	for len(*ptr) > 0 && ((*ptr)[0] == ' ' || (*ptr)[0] == '\t' || (*ptr)[0] == '-' || (*ptr)[0] == '.' || (*ptr)[0] == '/') {
		*ptr = (*ptr)[1:]
	}
	return timelibLookupMonth(ptr)
}

func timelibEatSpaces(ptr *string) {
	for {
		if len(*ptr) == 0 {
			break
		}
		if (*ptr)[0] == ' ' || (*ptr)[0] == '\t' {
			*ptr = (*ptr)[1:]
			continue
		}
		if len(*ptr) >= 3 && (*ptr)[0] == 0xe2 && (*ptr)[1] == 0x80 && (*ptr)[2] == 0xaf {
			*ptr = (*ptr)[3:]
			continue
		}
		if len(*ptr) >= 2 && (*ptr)[0] == 0xc2 && (*ptr)[1] == 0xa0 {
			*ptr = (*ptr)[2:]
			continue
		}
		break
	}
}

func timelibEatUntilSeparator(ptr *string) {
	if len(*ptr) > 0 {
		*ptr = (*ptr)[1:]
	}
	for len(*ptr) > 0 && !strings.ContainsRune(" \t.,:;/-0123456789", rune((*ptr)[0])) {
		*ptr = (*ptr)[1:]
	}
}

func timelibLookupRelunit(ptr *string) *RelUnit {
	if len(*ptr) == 0 {
		return nil
	}

	begin := *ptr

	for len(*ptr) > 0 && (*ptr)[0] != '\x00' && (*ptr)[0] != ' ' && (*ptr)[0] != ',' &&
		(*ptr)[0] != '\t' && (*ptr)[0] != ';' && (*ptr)[0] != ':' &&
		(*ptr)[0] != '/' && (*ptr)[0] != '.' && (*ptr)[0] != '-' &&
		(*ptr)[0] != '(' && (*ptr)[0] != ')' {
		*ptr = (*ptr)[1:]
	}

	if *ptr == begin {
		return nil
	}

	word := strings.ToLower(begin[:len(begin)-len(*ptr)])

	for i := range timelibRelunitLookup {
		if strings.EqualFold(word, timelibRelunitLookup[i].Name) {
			return &timelibRelunitLookup[i]
		}
	}

	return nil
}

func addWithOverflow(s *Scanner, e *int64, amount int64, multiplier int) {
	*e += (amount * int64(multiplier))
}

func timelibSetRelative(ptr *string, amount int64, behavior int, s *Scanner, timePart int) {
	relunit := timelibLookupRelunit(ptr)
	if relunit == nil {
		return
	}

	switch relunit.Unit {
	case TIMELIB_MICROSEC:
		addWithOverflow(s, &s.time.Relative.US, amount, relunit.Multiplier)
	case TIMELIB_SECOND:
		addWithOverflow(s, &s.time.Relative.S, amount, relunit.Multiplier)
	case TIMELIB_MINUTE:
		addWithOverflow(s, &s.time.Relative.I, amount, relunit.Multiplier)
	case TIMELIB_HOUR:
		addWithOverflow(s, &s.time.Relative.H, amount, relunit.Multiplier)
	case TIMELIB_DAY:
		addWithOverflow(s, &s.time.Relative.D, amount, relunit.Multiplier)
	case TIMELIB_MONTH:
		addWithOverflow(s, &s.time.Relative.M, amount, relunit.Multiplier)
	case TIMELIB_YEAR:
		addWithOverflow(s, &s.time.Relative.Y, amount, relunit.Multiplier)

	case TIMELIB_WEEKDAY:
		s.time.HaveRelative = true
		s.time.Relative.HaveWeekdayRelative = true
		if timePart != TIMELIB_TIME_PART_KEEP {
			s.time.HaveTime = false
			s.time.H = 0
			s.time.I = 0
			s.time.S = 0
			s.time.US = 0
		}
		if amount > 0 {
			s.time.Relative.D += (amount - 1) * 7
		} else {
			s.time.Relative.D += amount * 7
		}
		s.time.Relative.Weekday = relunit.Multiplier
		s.time.Relative.WeekdayBehavior = behavior

	case TIMELIB_SPECIAL:
		s.time.HaveRelative = true
		s.time.Relative.HaveSpecialRelative = true
		if timePart != TIMELIB_TIME_PART_KEEP {
			s.time.HaveTime = false
			s.time.H = 0
			s.time.I = 0
			s.time.S = 0
			s.time.US = 0
		}
		s.time.Relative.Special.Type = relunit.Multiplier
		s.time.Relative.Special.Amount = amount
	}
}

func abbrSearch(word string, gmtoffset int32, isdst int) *TzLookupTable {
	// Convert to uppercase for lookup since the abbreviation table uses uppercase
	upperWord := strings.ToUpper(word)
	
	// Use the comprehensive timezone abbreviation lookup
	abbr := LookupTimezoneAbbr(upperWord, gmtoffset, isdst)
	if abbr != nil {
		return &TzLookupTable{
			Name:       abbr.TzID,
			Type:       abbr.IsDST,
			GmtOffset:  float32(abbr.OffsetSec),
			FullTzName: word, // Use original case from input, not lowercase from table
		}
	}
	
	// Fallback to UTC/GMT special case
	if strings.EqualFold("utc", word) || strings.EqualFold("gmt", word) {
		return &timelibTimezoneUTC[0]
	}
	return nil
}

func timelibLookupAbbr(ptr *string, dst *int, tzAbbr *string, found *int) int32 {
	if len(*ptr) == 0 {
		*found = 0
		return 0
	}

	begin := *ptr

	for len(*ptr) > 0 && (((*ptr)[0] >= 'A' && (*ptr)[0] <= 'Z') ||
		((*ptr)[0] >= 'a' && (*ptr)[0] <= 'z') ||
		((*ptr)[0] >= '0' && (*ptr)[0] <= '9') ||
		(*ptr)[0] == '/' || (*ptr)[0] == '_' || (*ptr)[0] == '-' || (*ptr)[0] == '+') {
		*ptr = (*ptr)[1:]
	}

	if *ptr == begin {
		*found = 0
		return 0
	}

	word := begin[:len(begin)-len(*ptr)]
	*tzAbbr = strings.ToUpper(word)

	if len(word) < MAX_ABBR_LEN {
		tp := abbrSearch(word, -1, 0)
		if tp != nil {
			value := int32(tp.GmtOffset)
			*dst = tp.Type
			value -= int32(tp.Type * 3600)
			*found = 1
			return value
		}
	}

	*found = 0
	return 0
}

func sHOUR(a float64) int32 {
	return int32(a * 3600)
}

func sMIN(a float64) int32 {
	return int32(a * 60)
}

func timelibParseTzCor(ptr *string, tzNotFound *int) int32 {
	*tzNotFound = 1

	if len(*ptr) == 0 {
		return 0
	}

	begin := *ptr

	for len(*ptr) > 0 && (unicode.IsDigit(rune((*ptr)[0])) || (*ptr)[0] == ':') {
		*ptr = (*ptr)[1:]
	}

	length := len(begin) - len(*ptr)
	if length == 0 {
		return 0
	}

	tzStr := begin[:length]

	switch length {
	case 1, 2:
		val, _ := strconv.ParseInt(tzStr, 10, 32)
		*tzNotFound = 0
		return sHOUR(float64(val))

	case 3, 4:
		if len(tzStr) > 1 && tzStr[1] == ':' {
			h, _ := strconv.ParseInt(tzStr[:1], 10, 32)
			m, _ := strconv.ParseInt(tzStr[2:], 10, 32)
			*tzNotFound = 0
			return sHOUR(float64(h)) + sMIN(float64(m))
		} else if len(tzStr) > 2 && tzStr[2] == ':' {
			h, _ := strconv.ParseInt(tzStr[:2], 10, 32)
			m, _ := strconv.ParseInt(tzStr[3:], 10, 32)
			*tzNotFound = 0
			return sHOUR(float64(h)) + sMIN(float64(m))
		} else {
			val, _ := strconv.ParseInt(tzStr, 10, 32)
			*tzNotFound = 0
			return sHOUR(float64(val/100)) + sMIN(float64(val%100))
		}

	case 5:
		if len(tzStr) > 2 && tzStr[2] != ':' {
			break
		}
		h, _ := strconv.ParseInt(tzStr[:2], 10, 32)
		m, _ := strconv.ParseInt(tzStr[3:5], 10, 32)
		*tzNotFound = 0
		return sHOUR(float64(h)) + sMIN(float64(m))

	case 6:
		val, _ := strconv.ParseInt(tzStr, 10, 32)
		*tzNotFound = 0
		return sHOUR(float64(val/10000)) + sMIN(float64((val/100)%100)) + int32(val%100)

	case 8:
		if len(tzStr) > 2 && tzStr[2] != ':' {
			break
		}
		if len(tzStr) > 5 && tzStr[5] != ':' {
			break
		}
		h, _ := strconv.ParseInt(tzStr[:2], 10, 32)
		m, _ := strconv.ParseInt(tzStr[3:5], 10, 32)
		s, _ := strconv.ParseInt(tzStr[6:8], 10, 32)
		*tzNotFound = 0
		return sHOUR(float64(h)) + sMIN(float64(m)) + int32(s)
	}

	return 0
}

func timelibParseTzMinutes(ptr *string, t *Time) int32 {
	retval := int32(TIMELIB_UNSET)

	if len(*ptr) == 0 {
		return retval
	}

	if (*ptr)[0] != '+' && (*ptr)[0] != '-' {
		return retval
	}

	sign := (*ptr)[0]
	begin := *ptr
	*ptr = (*ptr)[1:]

	for len(*ptr) > 0 && unicode.IsDigit(rune((*ptr)[0])) {
		*ptr = (*ptr)[1:]
	}

	if len(*ptr) == len(begin)-1 {
		return retval
	}

	minStr := begin[1 : len(begin)-len(*ptr)]
	minutes, err := strconv.ParseInt(minStr, 10, 32)
	if err != nil {
		return retval
	}

	t.IsLocaltime = true
	t.ZoneType = TIMELIB_ZONETYPE_OFFSET
	t.Dst = 0

	if sign == '+' {
		retval = sMIN(float64(minutes))
	} else {
		retval = -1 * sMIN(float64(minutes))
	}

	return retval
}

func timelibParseZone(ptr *string, dst *int, t *Time, tzNotFound *int, tzdb *TzDB, tzWrapper TzGetWrapper) int32 {
	retval := int32(0)
	parenCount := 0

	*tzNotFound = 0

	for len(*ptr) > 0 && ((*ptr)[0] == ' ' || (*ptr)[0] == '\t' || (*ptr)[0] == '(') {
		if (*ptr)[0] == '(' {
			parenCount++
		}
		*ptr = (*ptr)[1:]
	}

	if len(*ptr) >= 4 && (*ptr)[0] == 'G' && (*ptr)[1] == 'M' && (*ptr)[2] == 'T' &&
		((*ptr)[3] == '+' || (*ptr)[3] == '-') {
		*ptr = (*ptr)[3:]
	}

	if len(*ptr) == 0 {
		return retval
	}

	if (*ptr)[0] == '+' {
		*ptr = (*ptr)[1:]
		t.IsLocaltime = true
		t.ZoneType = TIMELIB_ZONETYPE_OFFSET
		t.Dst = 0
		retval = timelibParseTzCor(ptr, tzNotFound)
	} else if (*ptr)[0] == '-' {
		*ptr = (*ptr)[1:]
		t.IsLocaltime = true
		t.ZoneType = TIMELIB_ZONETYPE_OFFSET
		t.Dst = 0
		retval = -1 * timelibParseTzCor(ptr, tzNotFound)
	} else {
		var found int
		var offset int32
		var tzAbbr string

		t.IsLocaltime = true

		offset = timelibLookupAbbr(ptr, dst, &tzAbbr, &found)
		if found != 0 {
			t.ZoneType = TIMELIB_ZONETYPE_ABBR
			t.Dst = *dst
			t.TzAbbr = tzAbbr
		}

		if found == 0 || tzAbbr == "UTC" {
			if tzWrapper != nil {
				var errorCode int
				res, _ := tzWrapper(tzAbbr, tzdb, &errorCode)
				if res != nil {
					t.TzInfo = res
					t.ZoneType = TIMELIB_ZONETYPE_ID
					found++
				}
			}
		}

		if found == 0 {
			*tzNotFound = 1
		}
		retval = offset
	}

	for parenCount > 0 && len(*ptr) > 0 && (*ptr)[0] == ')' {
		*ptr = (*ptr)[1:]
		parenCount--
	}

	return retval
}

func processYear(y *int64, length int) {
	if *y == TIMELIB_UNSET || length >= 4 {
		return
	}
	if *y < 100 {
		if *y < 70 {
			*y += 2000
		} else {
			*y += 1900
		}
	}
}

func timelibDaynrFromWeeknr(iyear, iweek, idow int64) int64 {
	var dow, doy int64

	// Calculate the day of week for January 1st of the given year
	dow = (dayOfWeek(iyear, 1, 1) + 6) % 7
	
	// Calculate the day of year for the given ISO week and day
	doy = -dow + (iweek * 7) + idow - 7
	
	// Handle year boundaries - if doy is negative or too large, adjust the year
	if doy <= 0 {
		// Date is in previous year
		return doy
	} else if doy > 365 {
		// Check if it's a leap year and adjust accordingly
		if (iyear%4 == 0 && iyear%100 != 0) || (iyear%400 == 0) {
			if doy > 366 {
				// Date is in next year
				return doy - 366
			}
		} else {
			if doy > 365 {
				// Date is in next year
				return doy - 365
			}
		}
	}
	
	return doy
}

func dayOfWeek(y, m, d int64) int64 {
	var c1, c2, y1, y2, d1, d2 int64

	if m > 2 {
		c1 = (m - 2) / 12
	} else {
		c1 = -1
	}
	c2 = c1

	y1 = (y - c1) / 100
	y2 = (y - c1) % 100

	d1 = ((26*(m-2-12*c2) - 2) / 10) + d + y2 + (y2 / 4) + (y1 / 4) - 2*y1 + 77

	d2 = (d1 - 7*(d1/7))

	return d2
}

func scan(s *Scanner, tzGetWrapper TzGetWrapper) int {
	var str string
	var ptr string

	// re2go generic API functions use s.cur directly (not a local copy)
	YYPEEK := func() byte {
		// Check if we're at or past the limit
		if s.cur != nil && uintptr(unsafe.Pointer(s.cur)) >= uintptr(unsafe.Pointer(s.lim)) {
			return 0 // Return null byte when at/past limit
		}
		if s.cur != nil {
			return *s.cur
		}
		return 0
	}
	YYSKIP := func() {
		if s.cur != nil && uintptr(unsafe.Pointer(s.cur)) < uintptr(unsafe.Pointer(s.lim)) {
			s.cur = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s.cur)) + 1))
		}
	}
	YYBACKUP := func() {
		s.ptr = s.cur
	}
	YYRESTORE := func() {
		s.cur = s.ptr
	}
	YYLESSTHAN := func(n int) bool {
		return uintptr(unsafe.Pointer(s.lim)) - uintptr(unsafe.Pointer(s.cur)) < uintptr(n)
	}
	_ = YYRESTORE // may not be used
	_ = YYLESSTHAN // may not be used

std:
	s.tok = s.cur
	s.len = 0
/*!re2c
re2c:define:YYCTYPE = byte;
re2c:define:YYPEEK = "YYPEEK()";
re2c:define:YYSKIP = "YYSKIP()";
re2c:define:YYBACKUP = "YYBACKUP()";
re2c:define:YYRESTORE = "YYRESTORE()";
re2c:define:YYLESSTHAN = "YYLESSTHAN";
re2c:yyfill:enable = 0;

any = [\000-\377];

nbsp = [\302][\240];
nnbsp = [\342][\200][\257];
space = [ \t]+ | nbsp+ | nnbsp+;
frac = "."[0-9]+;

ago = 'ago';

hour24 = [01]?[0-9] | "2"[0-4];
hour24lz = [01][0-9] | "2"[0-4];
hour12 = "0"?[1-9] | "1"[0-2];
minute = [0-5]?[0-9];
minutelz = [0-5][0-9];
second = minute | "60";
secondlz = minutelz | "60";
meridian = ([AaPp] "."? [Mm] "."?) [\000\t ];
tz = "("? [A-Za-z]{1,6} ")"? | [A-Z][a-z]+([_/-][A-Za-z]+)+;
tzcorrection = "GMT"? [+-] ((hour24 (":"? minute)?) | (hour24lz minutelz secondlz) | (hour24lz ":" minutelz ":" secondlz));

daysuf = "st" | "nd" | "rd" | "th";

month = "0"? [0-9] | "1"[0-2];
day   = (([0-2]?[0-9]) | ("3"[01])) daysuf?;
year  = [0-9]{1,4};
year2 = [0-9]{2};
year4 = [0-9]{4};
year4withsign = [+-]? [0-9]{4};
yearx = [+-] [0-9]{5,19};

dayofyear = "00"[1-9] | "0"[1-9][0-9] | [1-2][0-9][0-9] | "3"[0-5][0-9] | "36"[0-6];
weekofyear = "0"[1-9] | [1-4][0-9] | "5"[0-3];

monthlz = "0" [0-9] | "1" [0-2];
daylz   = "0" [0-9] | [1-2][0-9] | "3" [01];

dayfulls = 'sundays' | 'mondays' | 'tuesdays' | 'wednesdays' | 'thursdays' | 'fridays' | 'saturdays';
dayfull = 'sunday' | 'monday' | 'tuesday' | 'wednesday' | 'thursday' | 'friday' | 'saturday';
dayabbr = 'sun' | 'mon' | 'tue' | 'wed' | 'thu' | 'fri' | 'sat' | 'sun';
dayspecial = 'weekday' | 'weekdays';
daytext = dayfulls | dayfull | dayabbr | dayspecial;

monthfull = 'january' | 'february' | 'march' | 'april' | 'may' | 'june' | 'july' | 'august' | 'september' | 'october' | 'november' | 'december';
monthabbr = 'jan' | 'feb' | 'mar' | 'apr' | 'may' | 'jun' | 'jul' | 'aug' | 'sep' | 'sept' | 'oct' | 'nov' | 'dec';
monthroman = "I" | "II" | "III" | "IV" | "V" | "VI" | "VII" | "VIII" | "IX" | "X" | "XI" | "XII";
monthtext = monthfull | monthabbr | monthroman;

timetiny12 = hour12 space? meridian;
timeshort12 = hour12[:.]minutelz space? meridian;
timelong12 = hour12[:.]minute[:.]secondlz space? meridian;

timetiny24 = 't' hour24;
timeshort24 = 't'? hour24[:.]minute;
timelong24 =  't'? hour24[:.]minute[:.]second;
iso8601long =  't'? hour24 [:.] minute [:.] second frac;

iso8601normtz =  't'? hour24 [:.] minute [:.] secondlz space? (tzcorrection | tz);

gnunocolon       = 't'? hour24lz minutelz;
gnunocolontz     = hour24lz minutelz space? (tzcorrection | tz);
iso8601nocolon   = 't'? hour24lz minutelz secondlz;
iso8601nocolontz = hour24lz minutelz secondlz space? (tzcorrection | tz);

americanshort    = month "/" day;
american         = month "/" day "/" year;
iso8601dateslash = year4 "/" monthlz "/" daylz "/"?;
dateslash        = year4 "/" month "/" day;
iso8601date4     = year4withsign "-" monthlz "-" daylz;
iso8601date2     = year2 "-" monthlz "-" daylz;
iso8601datex     = yearx "-" monthlz "-" daylz;
gnudateshorter   = year4 "-" month;
gnudateshort     = year "-" month "-" day;
pointeddate4     = day [.\t-] month [.-] year4;
pointeddate2     = day [.\t] month "." year2;
datefull         = day ([ \t.-])* monthtext ([ \t.-])* year;
datenoday        = monthtext ([ .\t-])* year4;
datenodayrev     = year4 ([ .\t-])* monthtext;
datetextual      = monthtext ([ .\t-])* day [,.stndrh\t ]+ year;
datenoyear       = monthtext ([ .\t-])* day ([,.stndrh\t ]+|[\000]);
datenoyearrev    = day ([ .\t-])* monthtext;
datenocolon      = year4 monthlz daylz;

soap             = year4 "-" monthlz "-" daylz "T" hour24lz ":" minutelz ":" secondlz frac tzcorrection?;
xmlrpc           = year4 monthlz daylz "T" hour24 ":" minutelz ":" secondlz;
xmlrpcnocolon    = year4 monthlz daylz 't' hour24 minutelz secondlz;
wddx             = year4 "-" month "-" day "T" hour24 ":" minute ":" second;
pgydotd          = year4 [.-]? dayofyear;
pgtextshort      = monthabbr "-" daylz "-" year;
pgtextreverse    = year "-" monthabbr "-" daylz;
mssqltime        = hour12 ":" minutelz ":" secondlz [:.] [0-9]+ meridian;
isoweekday       = year4 "-"? "W" weekofyear "-"? [0-7];
isoweek          = year4 "-"? "W" weekofyear;
exif             = year4 ":" monthlz ":" daylz " " hour24lz ":" minutelz ":" secondlz;
firstdayof       = 'first day of';
lastdayof        = 'last day of';
backof           = 'back of ' hour24 (space? meridian)?;
frontof          = 'front of ' hour24 (space? meridian)?;

clf              = day "/" monthabbr "/" year4 ":" hour24lz ":" minutelz ":" secondlz space tzcorrection;

timestamp        = "@" "-"? [0-9]+;
timestampms      = "@" "-"? [0-9]+ "." [0-9]{0,6};

dateshortwithtimeshort12  = datenoyear timeshort12;
dateshortwithtimelong12   = datenoyear timelong12;
dateshortwithtimeshort  = datenoyear timeshort24;
dateshortwithtimelong   = datenoyear timelong24;
dateshortwithtimelongtz = datenoyear iso8601normtz;

reltextnumber = 'first'|'second'|'third'|'fourth'|'fifth'|'sixth'|'seventh'|'eight'|'eighth'|'ninth'|'tenth'|'eleventh'|'twelfth';
reltexttext = 'next'|'last'|'previous'|'this';
reltextunit = 'ms' | 'µs' | (('msec'|'millisecond'|'µsec'|'microsecond'|'usec'|'sec'|'second'|'min'|'minute'|'hour'|'day'|'fortnight'|'forthnight'|'month'|'year') 's'?) | 'weeks' | daytext;

relnumber = ([+-]*[ \t]*[0-9]{1,13});
relative = relnumber space? (reltextunit | 'week' );
relativetext = (reltextnumber|reltexttext) space reltextunit;
relativetextweek = reltexttext space 'week';

weekdayof        = (reltextnumber|reltexttext) space (dayfulls|dayfull|dayabbr) space 'of';

*/

/*!re2c
	'yesterday'
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		s.time.Relative.D = -1
		return TIMELIB_RELATIVE
	}

	'now'
	{
		str = timelibString(s)
		ptr = str
		return TIMELIB_RELATIVE
	}

	'noon'
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		s.time.HaveTime = true
		s.time.H = 12
		return TIMELIB_RELATIVE
	}

	'midnight' | 'today'
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		return TIMELIB_RELATIVE
	}

	'tomorrow'
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		s.time.Relative.D = 1
		return TIMELIB_RELATIVE
	}

	timestamp
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.HaveDate = false
		s.time.D = 0
		s.time.M = 0
		s.time.Y = 0
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		if s.time.HaveZone {
			addError(s, TIMELIB_ERR_DOUBLE_TZ, "Double timezone specification")
			return TIMELIB_ERROR
		} else {
			s.time.HaveZone = true
		}

		i := timelibGetSignedNr(s, &ptr, 24)
		s.time.Y = 1970
		s.time.M = 1
		s.time.D = 1
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		s.time.Relative.S += i
		s.time.IsLocaltime = true
		s.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
		s.time.Z = 0
		s.time.Dst = 0

		return TIMELIB_RELATIVE
	}

	timestampms
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.HaveDate = false
		s.time.D = 0
		s.time.M = 0
		s.time.Y = 0
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		if s.time.HaveZone {
			addError(s, TIMELIB_ERR_DOUBLE_TZ, "Double timezone specification")
			return TIMELIB_ERROR
		} else {
			s.time.HaveZone = true
		}

		isNegative := len(ptr) > 1 && ptr[1] == '-'

		i := timelibGetSignedNr(s, &ptr, 24)

		ptrBefore := ptr
		us := timelibGetSignedNr(s, &ptr, 6)
		us = int64(float64(us) * math.Pow(10, 7-float64(len(ptrBefore)-len(ptr))))
		if isNegative {
			us *= -1
		}

		s.time.Y = 1970
		s.time.M = 1
		s.time.D = 1
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		s.time.Relative.S += i
		s.time.Relative.US = us
		s.time.IsLocaltime = true
		s.time.ZoneType = TIMELIB_ZONETYPE_OFFSET
		s.time.Z = 0
		s.time.Dst = 0

		return TIMELIB_RELATIVE
	}

	firstdayof | lastdayof
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0

		if len(ptr) > 0 && (ptr[0] == 'l' || ptr[0] == 'L') {
			s.time.Relative.FirstLastDayOf = TIMELIB_SPECIAL_LAST_DAY_OF_MONTH
		} else {
			s.time.Relative.FirstLastDayOf = TIMELIB_SPECIAL_FIRST_DAY_OF_MONTH
		}

		return TIMELIB_LF_DAY_OF_MONTH
	}

	backof | frontof
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		s.time.HaveTime = true

		if len(ptr) > 0 && ptr[0] == 'b' {
			s.time.H = timelibGetNr(&ptr, 2)
			s.time.I = 15
		} else {
			s.time.H = timelibGetNr(&ptr, 2) - 1
			s.time.I = 45
		}
		if len(ptr) > 0 && ptr[0] != '\x00' {
			timelibEatSpaces(&ptr)
			s.time.H += timelibMeridian(&ptr, s.time.H)
		}

		return TIMELIB_LF_DAY_OF_MONTH
	}

	weekdayof
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.Relative.HaveSpecialRelative = true

		behavior := 0
		i := timelibGetRelativeText(&ptr, &behavior)
		timelibEatSpaces(&ptr)
		if i > 0 {
			s.time.Relative.Special.Type = TIMELIB_SPECIAL_DAY_OF_WEEK_IN_MONTH
			timelibSetRelative(&ptr, i, 1, s, TIMELIB_TIME_PART_DONT_KEEP)
		} else {
			s.time.Relative.Special.Type = TIMELIB_SPECIAL_LAST_DAY_OF_WEEK_IN_MONTH
			timelibSetRelative(&ptr, i, behavior, s, TIMELIB_TIME_PART_DONT_KEEP)
		}
		return TIMELIB_WEEK_DAY_OF_MONTH
	}

	timetiny12 | timeshort12 | timelong12
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
			s.time.I = timelibGetNr(&ptr, 2)
			if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
				s.time.S = timelibGetNr(&ptr, 2)
			}
		}
		timelibEatSpaces(&ptr)
		s.time.H += timelibMeridian(&ptr, s.time.H)
		return TIMELIB_TIME12
	}

	mssqltime
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
			s.time.S = timelibGetNr(&ptr, 2)

			if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
				s.time.US = timelibGetFracNr(&ptr)
			}
		}
		timelibEatSpaces(&ptr)
		s.time.H += timelibMeridian(&ptr, s.time.H)
		return TIMELIB_TIME24_WITH_ZONE
	}

	timetiny24 | timeshort24 | timelong24 | iso8601long
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
			s.time.I = timelibGetNr(&ptr, 2)
			if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
				s.time.S = timelibGetNr(&ptr, 2)

				if len(ptr) > 0 && ptr[0] == '.' {
					s.time.US = timelibGetFracNr(&ptr)
				}
			}
		}

		if len(ptr) > 0 && ptr[0] != '\x00' {
			tzNotFound := 0
			s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
			if tzNotFound != 0 {
				addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
			}
		}
		return TIMELIB_TIME24_WITH_ZONE
	}

	gnunocolon
	{
		str = timelibString(s)
		ptr = str
		if !s.time.HaveTime {
			s.time.H = timelibGetNr(&ptr, 2)
			s.time.I = timelibGetNr(&ptr, 2)
			s.time.S = 0
			s.time.HaveTime = true
		} else {
			s.time.Y = timelibGetNr(&ptr, 4)
		}
		return TIMELIB_GNU_NOCOLON
	}

	gnunocolontz
	{
		str = timelibString(s)
		ptr = str
		if !s.time.HaveTime {
			s.time.H = timelibGetNr(&ptr, 2)
			s.time.I = timelibGetNr(&ptr, 2)
			s.time.S = 0
			s.time.HaveTime = true
			if len(ptr) > 0 && ptr[0] != '\x00' {
				tzNotFound := 0
				s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
				if tzNotFound != 0 {
					addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
				}
			}
		} else {
			s.time.Y = timelibGetNr(&ptr, 4)
		}
		return TIMELIB_GNU_NOCOLON_TZ
	}

	iso8601nocolon
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		s.time.S = timelibGetNr(&ptr, 2)

		if len(ptr) > 0 && ptr[0] != '\x00' {
			tzNotFound := 0
			s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
			if tzNotFound != 0 {
				addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
			}
		}
		return TIMELIB_ISO_NOCOLON
	}

	iso8601nocolontz
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		s.time.S = timelibGetNr(&ptr, 2)

		if len(ptr) > 0 && ptr[0] != '\x00' {
			tzNotFound := 0
			s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
			if tzNotFound != 0 {
				addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
			}
		}
		return TIMELIB_ISO_NOCOLON_TZ
	}

	americanshort | american
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && ptr[0] == '/' {
			length := 0
			s.time.Y = timelibGetNrEx(&ptr, 4, &length)
			processYear(&s.time.Y, length)
		}
		return TIMELIB_AMERICAN
	}

	iso8601date4 | iso8601dateslash | dateslash
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.Y = timelibGetSignedNr(s, &ptr, 4)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		return TIMELIB_ISO_DATE
	}

	iso8601date2
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		processYear(&s.time.Y, length)
		return TIMELIB_ISO_DATE
	}

	iso8601datex
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.Y = timelibGetSignedNr(s, &ptr, 19)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		return TIMELIB_ISO_DATE
	}

	gnudateshorter
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = 1
		processYear(&s.time.Y, length)
		return TIMELIB_ISO_DATE
	}

	gnudateshort
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		processYear(&s.time.Y, length)
		return TIMELIB_ISO_DATE
	}

	datefull
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.D = timelibGetNr(&ptr, 2)
		timelibSkipDaySuffix(&ptr)
		s.time.M = timelibGetMonth(&ptr)
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		processYear(&s.time.Y, length)
		return TIMELIB_DATE_FULL
	}

	pointeddate4
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.D = timelibGetNr(&ptr, 2)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.Y = timelibGetNr(&ptr, 4)
		return TIMELIB_DATE_FULL_POINTED
	}

	pointeddate2
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.D = timelibGetNr(&ptr, 2)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.Y = timelibGetNrEx(&ptr, 2, &length)
		processYear(&s.time.Y, length)
		return TIMELIB_DATE_FULL_POINTED
	}

	datenoday
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.M = timelibGetMonth(&ptr)
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.D = 1
		processYear(&s.time.Y, length)
		return TIMELIB_DATE_NO_DAY
	}

	datenodayrev
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.M = timelibGetMonth(&ptr)
		s.time.D = 1
		processYear(&s.time.Y, length)
		return TIMELIB_DATE_NO_DAY
	}


	datetextual | datenoyear
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.M = timelibGetMonth(&ptr)
		s.time.D = timelibGetNr(&ptr, 2)
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		processYear(&s.time.Y, length)
		return TIMELIB_DATE_TEXT
	}

	datenoyearrev
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.D = timelibGetNr(&ptr, 2)
		timelibSkipDaySuffix(&ptr)
		s.time.M = timelibGetMonth(&ptr)
		return TIMELIB_DATE_TEXT
	}

	datenocolon
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.Y = timelibGetNr(&ptr, 4)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		return TIMELIB_DATE_NOCOLON
	}

	xmlrpc | xmlrpcnocolon | soap | wddx | exif
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.HaveDate = true
		s.time.Y = timelibGetNr(&ptr, 4)
		s.time.M = timelibGetNr(&ptr, 2)
		s.time.D = timelibGetNr(&ptr, 2)
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		s.time.S = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && ptr[0] == '.' {
			s.time.US = timelibGetFracNr(&ptr)
			if len(ptr) > 0 {
				tzNotFound := 0
				s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
				if tzNotFound != 0 {
					addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
				}
			}
		}
		return TIMELIB_XMLRPC_SOAP
	}

	pgydotd
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.D = timelibGetNr(&ptr, 3)
		s.time.M = 1
		processYear(&s.time.Y, length)
		return TIMELIB_PG_YEARDAY
	}

	isoweekday
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.HaveRelative = true

		s.time.Y = timelibGetNr(&ptr, 4)
		w := timelibGetNr(&ptr, 2)
		d := timelibGetNr(&ptr, 1)
		s.time.M = 1
		s.time.D = 1
		s.time.Relative.D = timelibDaynrFromWeeknr(s.time.Y, w, d)

		return TIMELIB_ISO_WEEK
	}

	isoweek
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.HaveRelative = true

		s.time.Y = timelibGetNr(&ptr, 4)
		w := timelibGetNr(&ptr, 2)
		d := int64(1)
		s.time.M = 1
		s.time.D = 1
		s.time.Relative.D = timelibDaynrFromWeeknr(s.time.Y, w, d)

		return TIMELIB_ISO_WEEK
	}

	pgtextshort
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.M = timelibGetMonth(&ptr)
		s.time.D = timelibGetNr(&ptr, 2)
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		processYear(&s.time.Y, length)
		return TIMELIB_PG_TEXT
	}

	pgtextreverse
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		length := 0
		s.time.Y = timelibGetNrEx(&ptr, 4, &length)
		s.time.M = timelibGetMonth(&ptr)
		s.time.D = timelibGetNr(&ptr, 2)
		processYear(&s.time.Y, length)
		return TIMELIB_PG_TEXT
	}

	clf
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.HaveDate = true
		s.time.D = timelibGetNr(&ptr, 2)
		s.time.M = timelibGetMonth(&ptr)
		s.time.Y = timelibGetNr(&ptr, 4)
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		s.time.S = timelibGetNr(&ptr, 2)

		timelibEatSpaces(&ptr)

		tzNotFound := 0
		s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
		if tzNotFound != 0 {
			addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
		}
		return TIMELIB_CLF
	}

	year4
	{
		str = timelibString(s)
		ptr = str
		s.time.Y = timelibGetNr(&ptr, 4)
		return TIMELIB_CLF
	}

	ago
	{
		str = timelibString(s)
		ptr = str
		s.time.Relative.Y = 0 - s.time.Relative.Y
		s.time.Relative.M = 0 - s.time.Relative.M
		s.time.Relative.D = 0 - s.time.Relative.D
		s.time.Relative.H = 0 - s.time.Relative.H
		s.time.Relative.I = 0 - s.time.Relative.I
		s.time.Relative.S = 0 - s.time.Relative.S
		s.time.Relative.Weekday = 0 - s.time.Relative.Weekday
		if s.time.Relative.Weekday == 0 {
			s.time.Relative.Weekday = -7
		}
		if s.time.Relative.HaveSpecialRelative && s.time.Relative.Special.Type == TIMELIB_SPECIAL_WEEKDAY {
			s.time.Relative.Special.Amount = 0 - s.time.Relative.Special.Amount
		}
		return TIMELIB_AGO
	}

	daytext
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true
		s.time.Relative.HaveWeekdayRelative = true
		s.time.HaveTime = false
		s.time.H = 0
		s.time.I = 0
		s.time.S = 0
		s.time.US = 0
		relunit := timelibLookupRelunit(&ptr)
		if relunit != nil {
			s.time.Relative.Weekday = relunit.Multiplier
			if s.time.Relative.WeekdayBehavior != 2 {
				s.time.Relative.WeekdayBehavior = 1
			}
		}

		return TIMELIB_WEEKDAY
	}

	relativetextweek
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true

		for len(ptr) > 0 {
			behavior := 0
			i := timelibGetRelativeText(&ptr, &behavior)
			timelibEatSpaces(&ptr)
			timelibSetRelative(&ptr, i, behavior, s, TIMELIB_TIME_PART_DONT_KEEP)
			s.time.Relative.WeekdayBehavior = 2

			if s.time.Relative.HaveWeekdayRelative == false {
				s.time.Relative.HaveWeekdayRelative = true
				s.time.Relative.Weekday = 1
			}
		}
		return TIMELIB_RELATIVE
	}

	relativetext
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true

		for len(ptr) > 0 {
			behavior := 0
			i := timelibGetRelativeText(&ptr, &behavior)
			timelibEatSpaces(&ptr)
			timelibSetRelative(&ptr, i, behavior, s, TIMELIB_TIME_PART_DONT_KEEP)
		}
		return TIMELIB_RELATIVE
	}

	monthfull | monthabbr
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.M = timelibLookupMonth(&ptr)
		return TIMELIB_DATE_TEXT
	}

	tzcorrection | tz
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveZone {
			addError(s, TIMELIB_ERR_DOUBLE_TZ, "Double timezone specification")
			return TIMELIB_ERROR
		} else {
			s.time.HaveZone = true
		}
		timelibEatSpaces(&ptr)
		tzNotFound := 0
		s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
		if tzNotFound != 0 {
			addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
		}
		return TIMELIB_TIMEZONE
	}

	dateshortwithtimeshort12 | dateshortwithtimelong12
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.M = timelibGetMonth(&ptr)
		s.time.D = timelibGetNr(&ptr, 2)

		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && (ptr[0] == ':' || ptr[0] == '.') {
			s.time.S = timelibGetNr(&ptr, 2)

			if len(ptr) > 0 && ptr[0] == '.' {
				s.time.US = timelibGetFracNr(&ptr)
			}
		}

		s.time.H += timelibMeridian(&ptr, s.time.H)
		return TIMELIB_SHORTDATE_WITH_TIME
	}

	dateshortwithtimeshort | dateshortwithtimelong | dateshortwithtimelongtz
	{
		str = timelibString(s)
		ptr = str
		if s.time.HaveDate {
			addError(s, TIMELIB_ERR_DOUBLE_DATE, "Double date specification")
			return TIMELIB_ERROR
		}
		s.time.HaveDate = true
		s.time.M = timelibGetMonth(&ptr)
		s.time.D = timelibGetNr(&ptr, 2)

		if s.time.HaveTime {
			addError(s, TIMELIB_ERR_DOUBLE_TIME, "Double time specification")
			return TIMELIB_ERROR
		}
		s.time.HaveTime = true
		s.time.H = timelibGetNr(&ptr, 2)
		s.time.I = timelibGetNr(&ptr, 2)
		if len(ptr) > 0 && ptr[0] == ':' {
			s.time.S = timelibGetNr(&ptr, 2)

			if len(ptr) > 0 && ptr[0] == '.' {
				s.time.US = timelibGetFracNr(&ptr)
			}
		}

		if len(ptr) > 0 && ptr[0] != '\x00' {
			tzNotFound := 0
			s.time.Z = timelibParseZone(&ptr, &s.time.Dst, s.time, &tzNotFound, s.tzdb, tzGetWrapper)
			if tzNotFound != 0 {
				addError(s, TIMELIB_ERR_TZID_NOT_FOUND, "The timezone could not be found in the database")
			}
		}
		return TIMELIB_SHORTDATE_WITH_TIME
	}

	relative
	{
		str = timelibString(s)
		ptr = str
		s.time.HaveRelative = true

		for len(ptr) > 0 {
			i := timelibGetSignedNr(s, &ptr, 24)
			timelibEatSpaces(&ptr)
			timelibSetRelative(&ptr, i, 1, s, TIMELIB_TIME_PART_KEEP)
		}
		return TIMELIB_RELATIVE
	}

	[.,]
	{
		goto std
	}

	space
	{
		goto std
	}

	"\000"
	{
		return EOI
	}

	"\n"
	{
		s.pos = s.cur
		s.line++
		goto std
	}

	any
	{
		addError(s, TIMELIB_ERR_UNEXPECTED_CHARACTER, "Unexpected character")
		goto std
	}
*/
}

/*!max:re2c */
