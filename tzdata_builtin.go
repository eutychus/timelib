package timelib

import (
	_ "embed"
	"encoding/binary"
	"strings"
)

// Embedded timezone database files
//
//go:embed tzdata_data.bin
var builtinTzData []byte

//go:embed tzdata_index.bin
var builtinTzIndexRaw []byte

// parseBuiltinIndex parses the timezone index from the embedded data
func parseBuiltinIndex() []TzDBIndexEntry {
	entries := make([]TzDBIndexEntry, 0, 598)
	data := builtinTzIndexRaw
	pos := 0

	for pos < len(data) {
		// Find timezone name (ends with newline)
		nameEnd := pos
		for nameEnd < len(data) && data[nameEnd] != '\n' {
			nameEnd++
		}
		if nameEnd >= len(data) {
			break
		}

		name := string(data[pos:nameEnd])
		pos = nameEnd + 1

		// Read 4-byte position (big-endian)
		if pos+4 > len(data) {
			break
		}
		offset := binary.BigEndian.Uint32(data[pos : pos+4])
		pos += 4

		// Skip newline after position
		if pos < len(data) && data[pos] == '\n' {
			pos++
		}

		entries = append(entries, TzDBIndexEntry{
			ID:  name,
			Pos: int(offset),
		})
	}

	return entries
}

// builtinTzDB is the singleton built-in timezone database
var builtinTzDB *TzDB

// initBuiltinTzDB initializes the built-in timezone database
func initBuiltinTzDB() {
	if builtinTzDB != nil {
		return
	}

	index := parseBuiltinIndex()

	builtinTzDB = &TzDB{
		Version:   "2025.2",
		IndexSize: len(index),
		Index:     index,
		Data:      builtinTzData,
	}
}

// BuiltinDB returns the built-in timezone database
func BuiltinDB() *TzDB {
	if builtinTzDB == nil {
		initBuiltinTzDB()
	}
	return builtinTzDB
}

// FindTimezone looks up a timezone by name in the builtin database
func FindTimezone(name string) *TzDBIndexEntry {
	db := BuiltinDB()
	if db == nil {
		return nil
	}

	// Try exact match first
	for i := range db.Index {
		if db.Index[i].ID == name {
			return &db.Index[i]
		}
	}

	// Try case-insensitive match
	lowerName := strings.ToLower(name)
	for i := range db.Index {
		if strings.ToLower(db.Index[i].ID) == lowerName {
			return &db.Index[i]
		}
	}

	return nil
}
