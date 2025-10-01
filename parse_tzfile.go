package timelib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// readUint32BE reads a big-endian uint32
func readUint32BE(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

// readInt32BE reads a big-endian int32
func readInt32BE(data []byte) int32 {
	return int32(binary.BigEndian.Uint32(data))
}

// readInt64BE reads a big-endian int64
func readInt64BE(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
}

// detectFormat detects the timezone file format
func detectFormat(data []byte) (string, int, error) {
	if len(data) < 4 {
		return "", 0, errors.New("file too small")
	}

	// Check for PHP format
	if data[0] == 'P' && data[1] == 'H' && data[2] == 'P' {
		version := int(data[3] - '0')
		return "PHP", version, nil
	}

	// Check for TZif format
	if data[0] == 'T' && data[1] == 'Z' && data[2] == 'i' && data[3] == 'f' {
		version := 0
		if len(data) > 4 {
			switch data[4] {
			case '\x00':
				version = 1
			case '2':
				version = 2
			case '3':
				version = 3
			case '4':
				version = 4
			default:
				return "", 0, fmt.Errorf("unsupported TZif version: %c", data[4])
			}
		}
		return "TZif", version, nil
	}

	return "", 0, errors.New("unknown timezone file format")
}

// readPHPPreamble reads PHP-format timezone file preamble
func readPHPPreamble(data []byte, tz *TzInfo) (int, int, error) {
	if len(data) < 20 {
		return 0, 0, errors.New("PHP preamble too short")
	}

	version := int(data[3] - '0')
	pos := 4

	// Read BC flag
	tz.Bc = data[pos]
	pos++

	// Read country code
	copy(tz.Location.CountryCode[:], data[pos:pos+2])
	pos += 2

	// Skip rest of preamble (13 bytes)
	pos += 13

	return version, pos, nil
}

// readTZifPreamble reads TZif-format timezone file preamble
func readTZifPreamble(data []byte, tz *TzInfo) (int, int, error) {
	if len(data) < 20 {
		return 0, 0, errors.New("TZif preamble too short")
	}

	version := 0
	switch data[4] {
	case '\x00':
		version = 1
	case '2':
		version = 2
	case '3':
		version = 3
	case '4':
		version = 4
	default:
		return 0, 0, fmt.Errorf("unsupported TZif version")
	}

	pos := 5

	// Set BC flag and country code to default
	tz.Bc = 0
	tz.Location.CountryCode[0] = '?'
	tz.Location.CountryCode[1] = '?'
	tz.Location.CountryCode[2] = '\x00'

	// Skip rest of preamble (15 bytes)
	pos += 15

	return version, pos, nil
}

// readHeader reads TZif header
func readHeader(data []byte, pos int) (header struct {
	ttisgmtcnt uint32
	ttisstdcnt uint32
	leapcnt    uint32
	timecnt    uint32
	typecnt    uint32
	charcnt    uint32
}, newPos int, err error) {
	if len(data) < pos+24 {
		return header, pos, errors.New("header too short")
	}

	header.ttisgmtcnt = readUint32BE(data[pos : pos+4])
	header.ttisstdcnt = readUint32BE(data[pos+4 : pos+8])
	header.leapcnt = readUint32BE(data[pos+8 : pos+12])
	header.timecnt = readUint32BE(data[pos+12 : pos+16])
	header.typecnt = readUint32BE(data[pos+16 : pos+20])
	header.charcnt = readUint32BE(data[pos+20 : pos+24])

	return header, pos + 24, nil
}

// readTransitions reads transition times and indices
func readTransitions(data []byte, pos int, timecnt uint32, use64bit bool) (trans []int64, transIdx []uint8, newPos int, err error) {
	trans = make([]int64, timecnt)
	transIdx = make([]uint8, timecnt)

	// Read transition times
	if use64bit {
		if len(data) < pos+int(timecnt)*8 {
			return nil, nil, pos, errors.New("not enough data for 64-bit transitions")
		}
		for i := uint32(0); i < timecnt; i++ {
			trans[i] = readInt64BE(data[pos : pos+8])
			pos += 8
		}
	} else {
		if len(data) < pos+int(timecnt)*4 {
			return nil, nil, pos, errors.New("not enough data for 32-bit transitions")
		}
		for i := uint32(0); i < timecnt; i++ {
			trans[i] = int64(readInt32BE(data[pos : pos+4]))
			pos += 4
		}
	}

	// Read transition type indices
	if len(data) < pos+int(timecnt) {
		return nil, nil, pos, errors.New("not enough data for transition indices")
	}
	for i := uint32(0); i < timecnt; i++ {
		transIdx[i] = data[pos]
		pos++
	}

	return trans, transIdx, pos, nil
}

// readTypes reads timezone type information
func readTypes(data []byte, pos int, typecnt uint32) (types []TTInfo, newPos int, err error) {
	types = make([]TTInfo, typecnt)

	if len(data) < pos+int(typecnt)*6 {
		return nil, pos, errors.New("not enough data for types")
	}

	for i := uint32(0); i < typecnt; i++ {
		types[i].Offset = readInt32BE(data[pos : pos+4])
		types[i].IsDst = int(data[pos+4])
		types[i].AbbrIdx = int(data[pos+5])
		pos += 6
	}

	return types, pos, nil
}

// readAbbreviations reads abbreviation strings
func readAbbreviations(data []byte, pos int, charcnt uint32) (abbr string, newPos int, err error) {
	if len(data) < pos+int(charcnt) {
		return "", pos, errors.New("not enough data for abbreviations")
	}

	abbr = string(data[pos : pos+int(charcnt)])
	pos += int(charcnt)

	return abbr, pos, nil
}

// readLeapSeconds reads leap second records
func readLeapSeconds(data []byte, pos int, leapcnt uint32, use64bit bool) (leaps []TLInfo, newPos int, err error) {
	leaps = make([]TLInfo, leapcnt)

	if use64bit {
		if len(data) < pos+int(leapcnt)*12 {
			return nil, pos, errors.New("not enough data for 64-bit leap seconds")
		}
		for i := uint32(0); i < leapcnt; i++ {
			leaps[i].Trans = readInt64BE(data[pos : pos+8])
			leaps[i].Corr = int64(readInt32BE(data[pos+8 : pos+12]))
			pos += 12
		}
	} else {
		if len(data) < pos+int(leapcnt)*8 {
			return nil, pos, errors.New("not enough data for 32-bit leap seconds")
		}
		for i := uint32(0); i < leapcnt; i++ {
			leaps[i].Trans = int64(readInt32BE(data[pos : pos+4]))
			leaps[i].Corr = int64(readInt32BE(data[pos+4 : pos+8]))
			pos += 8
		}
	}

	return leaps, pos, nil
}

// readIndicators reads standard/wall and UT/local indicators
func readIndicators(data []byte, pos int, ttisstdcnt, ttisgmtcnt uint32, types []TTInfo) (newPos int, err error) {
	// Read isstd indicators
	if len(data) < pos+int(ttisstdcnt) {
		return pos, errors.New("not enough data for isstd indicators")
	}
	for i := uint32(0); i < ttisstdcnt && int(i) < len(types); i++ {
		types[i].IsStd = int(data[pos])
		pos++
	}

	// Read isut indicators
	if len(data) < pos+int(ttisgmtcnt) {
		return pos, errors.New("not enough data for isut indicators")
	}
	for i := uint32(0); i < ttisgmtcnt && int(i) < len(types); i++ {
		types[i].IsUtc = int(data[pos])
		pos++
	}

	return pos, nil
}

// readPosixString reads POSIX TZ string from v2+ format
func readPosixString(data []byte, pos int) (string, error) {
	if pos >= len(data) {
		return "", nil
	}

	// Find the newline that precedes the POSIX string
	for pos < len(data) && data[pos] == '\n' {
		pos++
	}

	if pos >= len(data) {
		return "", nil
	}

	// Read until newline or end of data
	end := pos
	for end < len(data) && data[end] != '\n' && data[end] != '\x00' {
		end++
	}

	if end > pos {
		return string(data[pos:end]), nil
	}

	return "", nil
}

// validateTransitions ensures transitions are monotonically increasing
func validateTransitions(trans []int64) error {
	for i := 1; i < len(trans); i++ {
		if trans[i] <= trans[i-1] {
			return errors.New("transitions do not increase monotonically")
		}
	}
	return nil
}

// ParseTzfileData parses timezone data from raw bytes
func ParseTzfileData(tzName string, data []byte, errorCode *int) (*TzInfo, error) {
	if errorCode != nil {
		*errorCode = TIMELIB_ERROR_NO_ERROR
	}

	if len(data) < 44 {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_CANNOT_OPEN_FILE
		}
		return nil, errors.New("file too small to be a valid timezone file")
	}

	tz := &TzInfo{
		Name: tzName,
	}

	// Detect format
	format, version, err := detectFormat(data)
	if err != nil {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_UNSUPPORTED_VERSION
		}
		return nil, err
	}

	var pos int

	// Read preamble based on format
	if format == "PHP" {
		version, pos, err = readPHPPreamble(data, tz)
		if err != nil {
			return nil, err
		}
	} else {
		version, pos, err = readTZifPreamble(data, tz)
		if err != nil {
			if errorCode != nil {
				*errorCode = TIMELIB_ERROR_UNSUPPORTED_VERSION
			}
			return nil, err
		}
	}

	// For version 2+, we need to skip the 32-bit section and read the 64-bit section
	if version >= 2 {
		// Read 32-bit header to know how much to skip
		header32, newPos, err := readHeader(data, pos)
		if err != nil {
			return nil, err
		}
		pos = newPos

		// Skip 32-bit data section
		// transition times (4 bytes each) + transition indices (1 byte each)
		pos += int(header32.timecnt) * 5
		// type info (6 bytes each)
		pos += int(header32.typecnt) * 6
		// abbreviation strings
		pos += int(header32.charcnt)
		// leap seconds (8 bytes each in v1)
		pos += int(header32.leapcnt) * 8
		// isstd and isut indicators
		pos += int(header32.ttisstdcnt)
		pos += int(header32.ttisgmtcnt)

		// Now read the TZif2+ preamble for 64-bit section
		if pos+20 > len(data) {
			if errorCode != nil {
				*errorCode = TIMELIB_ERROR_CORRUPT_NO_64BIT_PREAMBLE
			}
			return nil, errors.New("no 64-bit section found in v2+ file")
		}

		// Verify TZif magic again
		if data[pos] != 'T' || data[pos+1] != 'Z' || data[pos+2] != 'i' || data[pos+3] != 'f' {
			if errorCode != nil {
				*errorCode = TIMELIB_ERROR_CORRUPT_NO_64BIT_PREAMBLE
			}
			return nil, errors.New("invalid 64-bit section header")
		}

		pos += 20 // Skip second header
	}

	// Read header
	header, newPos, err := readHeader(data, pos)
	if err != nil {
		return nil, err
	}
	pos = newPos

	// Store header in appropriate bit field
	if version >= 2 {
		tz.Bit64.Ttisgmtcnt = uint64(header.ttisgmtcnt)
		tz.Bit64.Ttisstdcnt = uint64(header.ttisstdcnt)
		tz.Bit64.Leapcnt = uint64(header.leapcnt)
		tz.Bit64.Timecnt = uint64(header.timecnt)
		tz.Bit64.Typecnt = uint64(header.typecnt)
		tz.Bit64.Charcnt = uint64(header.charcnt)
	} else {
		tz.Bit32.Ttisgmtcnt = header.ttisgmtcnt
		tz.Bit32.Ttisstdcnt = header.ttisstdcnt
		tz.Bit32.Leapcnt = header.leapcnt
		tz.Bit32.Timecnt = header.timecnt
		tz.Bit32.Typecnt = header.typecnt
		tz.Bit32.Charcnt = header.charcnt
	}

	// Read transitions
	tz.Trans, tz.TransIdx, pos, err = readTransitions(data, pos, header.timecnt, version >= 2)
	if err != nil {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE
		}
		return nil, err
	}

	// Validate transitions
	if err := validateTransitions(tz.Trans); err != nil {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_CORRUPT_TRANSITIONS_DONT_INCREASE
		}
		return nil, err
	}

	// Read types
	tz.Type, pos, err = readTypes(data, pos, header.typecnt)
	if err != nil {
		return nil, err
	}

	// Read abbreviations
	tz.TimezoneAbbr, pos, err = readAbbreviations(data, pos, header.charcnt)
	if err != nil {
		return nil, err
	}

	// Read leap seconds
	tz.LeapTimes, pos, err = readLeapSeconds(data, pos, header.leapcnt, version >= 2)
	if err != nil {
		return nil, err
	}

	// Read indicators
	pos, err = readIndicators(data, pos, header.ttisstdcnt, header.ttisgmtcnt, tz.Type)
	if err != nil {
		return nil, err
	}

	// Read POSIX TZ string for v2+
	if version >= 2 {
		tz.PosixString, err = readPosixString(data, pos)
		if err != nil {
			return nil, err
		}

		// Parse POSIX string if present
		if tz.PosixString != "" {
			tz.PosixInfo, _ = ParsePosixString(tz.PosixString, tz)
		}
	}

	return tz, nil
}

// ParseTzfileFromFile parses timezone file from filesystem
func ParseTzfileFromFile(filename string, errorCode *int) (*TzInfo, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_CANNOT_OPEN_FILE
		}
		return nil, err
	}

	tzName := filepath.Base(filename)
	return ParseTzfileData(tzName, data, errorCode)
}

// ParseTzfileFromDB parses timezone file from database
func ParseTzfileFromDB(tzName string, tzdb *TzDB, errorCode *int) (*TzInfo, error) {
	if errorCode != nil {
		*errorCode = TIMELIB_ERROR_NO_ERROR
	}

	if tzdb == nil {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_NO_SUCH_TIMEZONE
		}
		return nil, errors.New("timezone database is nil")
	}

	// Find timezone in index
	var entry *TzDBIndexEntry
	for i := range tzdb.Index {
		if tzdb.Index[i].ID == tzName {
			entry = &tzdb.Index[i]
			break
		}
	}

	if entry == nil {
		if errorCode != nil {
			*errorCode = TIMELIB_ERROR_NO_SUCH_TIMEZONE
		}
		return nil, fmt.Errorf("timezone '%s' not found in database", tzName)
	}

	// If database has embedded data, parse from it
	if len(tzdb.Data) > 0 && entry.Pos < len(tzdb.Data) {
		// For embedded database, we'd need to know the size of each entry
		// For now, try to parse from the position
		return ParseTzfileData(tzName, tzdb.Data[entry.Pos:], errorCode)
	}

	// Otherwise, try to load from file
	// The index might store file paths
	if entry.Pos == 0 && entry.ID != "" {
		// Try as filename
		return ParseTzfileFromFile(entry.ID, errorCode)
	}

	if errorCode != nil {
		*errorCode = TIMELIB_ERROR_NO_SUCH_TIMEZONE
	}
	return nil, fmt.Errorf("cannot load timezone data for '%s'", tzName)
}
