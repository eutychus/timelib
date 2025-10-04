/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2015-2019 Derick Rethans
 * Copyright (c) 2025 Go port
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
	"unsafe"
)

// IsoIntervalScanner holds the state for parsing ISO 8601 intervals
type IsoIntervalScanner struct {
	str    []byte
	lim    *byte
	ptr    *byte
	cur    *byte
	tok    *byte
	pos    *byte
	line   int
	len    int
	errors *ErrorContainer

	Begin        *Time
	End          *Time
	Period       *RelTime
	Recurrences  int

	HavePeriod      bool
	HaveRecurrences bool
	HaveDate        bool
	HaveBeginDate   bool
	HaveEndDate     bool
}

// Helper function to add error to scanner
func addIsoError(s *IsoIntervalScanner, errorMsg string) {
	var position int
	var character byte

	if s.tok != nil {
		position = int(uintptr(unsafe.Pointer(s.tok)) - uintptr(unsafe.Pointer(&s.str[0])))
		character = *s.tok
	}

	s.errors.ErrorCount++
	s.errors.ErrorMessages = append(s.errors.ErrorMessages, ErrorMessage{
		Position:  position,
		Character: character,
		Message:   errorMsg,
	})
}

// Helper function to extract string from scanner
func timelibIsoString(s *IsoIntervalScanner) string {
	length := int(uintptr(unsafe.Pointer(s.cur)) - uintptr(unsafe.Pointer(s.tok)))
	if length <= 0 {
		return ""
	}

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		ptr := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s.tok)) + uintptr(i)))
		bytes[i] = *ptr
	}

	return string(bytes)
}

// Helper function to get unsigned number from string pointer
func timelibGetUnsignedNr(ptr *string, maxLength int) int64 {
	if ptr == nil || *ptr == "" {
		return TIMELIB_UNSET
	}

	str := *ptr
	length := 0

	for length < len(str) && length < maxLength && str[length] >= '0' && str[length] <= '9' {
		length++
	}

	if length == 0 {
		return TIMELIB_UNSET
	}

	result := int64(0)
	for i := 0; i < length; i++ {
		result = result*10 + int64(str[i]-'0')
	}

	*ptr = str[length:]
	return result
}

// scan is the main scanning function for ISO intervals
func scanIsoInterval(s *IsoIntervalScanner) int {
	var str string
	var ptr string

	// re2go generic API functions use s.cur directly
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
		return uintptr(unsafe.Pointer(s.lim))-uintptr(unsafe.Pointer(s.cur)) < uintptr(n)
	}
	_ = YYRESTORE
	_ = YYLESSTHAN

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
number = [0-9]+;

hour24lz = [01][0-9] | "2"[0-4];
minutelz = [0-5][0-9];
monthlz = "0" [1-9] | "1" [0-2];
monthlzz = "0" [0-9] | "1" [0-2];
daylz   = "0" [1-9] | [1-2][0-9] | "3" [01];
daylzz  = "0" [0-9] | [1-2][0-9] | "3" [01];
secondlz = minutelz;
year4 = [0-9]{4};
weekofyear = "0"[1-9] | [1-4][0-9] | "5"[0-3];

space = [ \t]+;
datetimebasic  = year4 monthlz daylz "T" hour24lz minutelz secondlz "Z";
datetimeextended  = year4 "-" monthlz "-" daylz "T" hour24lz ':' minutelz ':' secondlz "Z";
period   = "P" (number "Y")? (number "M")? (number "W")? (number "D")? ("T" (number "H")? (number "M")? (number "S")?)?;
combinedrep = "P" year4 "-" monthlzz "-" daylzz "T" hour24lz ':' minutelz ':' secondlz;

recurrences = "R" number;

isoweekday       = year4 "-"? "W" weekofyear "-"? [0-7];
isoweek          = year4 "-"? "W" weekofyear;

recurrences
{
	str = timelibIsoString(s)
	ptr = str
	ptr = ptr[1:] // skip 'R'
	s.Recurrences = int(timelibGetUnsignedNr(&ptr, 9))
	s.HaveRecurrences = true
	return TIMELIB_PERIOD
}

datetimebasic | datetimeextended
{
	var current *Time

	if s.HaveDate || s.HavePeriod {
		current = s.End
		s.HaveEndDate = true
	} else {
		current = s.Begin
		s.HaveBeginDate = true
	}

	str = timelibIsoString(s)
	ptr = str
	current.Y = timelibGetUnsignedNr(&ptr, 4)
	// Skip '-' if present (extended format)
	if len(ptr) > 0 && ptr[0] == '-' {
		ptr = ptr[1:]
	}
	current.M = timelibGetUnsignedNr(&ptr, 2)
	// Skip '-' if present (extended format)
	if len(ptr) > 0 && ptr[0] == '-' {
		ptr = ptr[1:]
	}
	current.D = timelibGetUnsignedNr(&ptr, 2)
	// Skip 'T'
	if len(ptr) > 0 && ptr[0] == 'T' {
		ptr = ptr[1:]
	}
	current.H = timelibGetUnsignedNr(&ptr, 2)
	// Skip ':' if present (extended format)
	if len(ptr) > 0 && ptr[0] == ':' {
		ptr = ptr[1:]
	}
	current.I = timelibGetUnsignedNr(&ptr, 2)
	// Skip ':' if present (extended format)
	if len(ptr) > 0 && ptr[0] == ':' {
		ptr = ptr[1:]
	}
	current.S = timelibGetUnsignedNr(&ptr, 2)
	s.HaveDate = true
	return TIMELIB_ISO_DATE_INTRVL
}

period
{
	var nr int64
	inTime := false

	str = timelibIsoString(s)
	ptr = str
	ptr = ptr[1:] // skip 'P'

	for len(ptr) > 0 && s.errors.ErrorCount == 0 {
		if ptr[0] == 'T' {
			inTime = true
			ptr = ptr[1:]
		}
		if len(ptr) == 0 {
			addIsoError(s, "Missing expected time part")
			break
		}

		nr = timelibGetUnsignedNr(&ptr, 12)
		if len(ptr) == 0 {
			break
		}

		switch ptr[0] {
		case 'Y':
			s.Period.Y = nr
		case 'W':
			s.Period.D += nr * 7
		case 'D':
			s.Period.D += nr
		case 'H':
			s.Period.H = nr
		case 'S':
			s.Period.S = nr
		case 'M':
			if inTime {
				s.Period.I = nr
			} else {
				s.Period.M = nr
			}
		default:
			addIsoError(s, "Undefined period specifier")
		}
		ptr = ptr[1:]
	}
	s.HavePeriod = true
	return TIMELIB_PERIOD
}

combinedrep
{
	str = timelibIsoString(s)
	ptr = str
	ptr = ptr[1:] // skip 'P'
	s.Period.Y = timelibGetUnsignedNr(&ptr, 4)
	ptr = ptr[1:] // skip '-'
	s.Period.M = timelibGetUnsignedNr(&ptr, 2)
	ptr = ptr[1:] // skip '-'
	s.Period.D = timelibGetUnsignedNr(&ptr, 2)
	ptr = ptr[1:] // skip 'T'
	s.Period.H = timelibGetUnsignedNr(&ptr, 2)
	ptr = ptr[1:] // skip ':'
	s.Period.I = timelibGetUnsignedNr(&ptr, 2)
	ptr = ptr[1:] // skip ':'
	s.Period.S = timelibGetUnsignedNr(&ptr, 2)
	s.HavePeriod = true
	return TIMELIB_PERIOD
}

[ .,\t/]
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
	addIsoError(s, "Unexpected character")
	goto std
}
*/
}

/*!max:re2c */
