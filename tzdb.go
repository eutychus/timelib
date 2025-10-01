package timelib

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ZoneinfoDir loads timezone database from a directory
func ZoneinfoDir(directory string) (*TzDB, error) {
	if directory == "" {
		return nil, fmt.Errorf("directory cannot be empty")
	}

	// Check if directory exists
	info, err := os.Stat(directory)
	if err != nil {
		return nil, fmt.Errorf("cannot access directory %s: %v", directory, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", directory)
	}

	// Get absolute path for the directory
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path for %s: %v", directory, err)
	}

	tzdb := &TzDB{
		Version:   "custom",
		IndexSize: 0,
		Index:     []TzDBIndexEntry{},
		Data:      []byte{},
		BaseDir:   absDir,
	}

	// Walk the directory tree
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		// Skip directories
		if info.IsDir() {
			// Skip certain directories
			name := info.Name()
			if name == "posix" || name == "right" {
				return filepath.SkipDir
			}
			return nil
		}

		// Filter out non-timezone files
		name := info.Name()
		if name == "." || name == ".." ||
			name == "posixrules" || name == "localtime" ||
			strings.Contains(name, ".list") ||
			strings.Contains(name, ".tab") {
			return nil
		}

		// Check if this looks like a timezone file
		if !isValidTzFile(path) {
			return nil
		}

		// Extract relative path from directory
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return nil // Skip this file
		}

		// Normalize path separators
		relPath = filepath.ToSlash(relPath)

		// Add to index
		entry := TzDBIndexEntry{
			ID:  relPath,
			Pos: 0, // File-based entries use Pos=0
		}
		tzdb.Index = append(tzdb.Index, entry)
		tzdb.IndexSize++

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error scanning directory: %v", err)
	}

	// Sort index by ID
	sort.Slice(tzdb.Index, func(i, j int) bool {
		return strings.ToLower(tzdb.Index[i].ID) < strings.ToLower(tzdb.Index[j].ID)
	})

	return tzdb, nil
}

// isValidTzFile checks if a file is likely a valid timezone file
func isValidTzFile(path string) bool {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Check file size
	info, err := file.Stat()
	if err != nil {
		return false
	}

	if !info.Mode().IsRegular() || info.Size() < 20 {
		return false
	}

	// Read first 20 bytes to check magic number
	buf := make([]byte, 20)
	n, err := file.Read(buf)
	if err != nil || n < 20 {
		return false
	}

	// Check for TZif or PHP magic
	if buf[0] == 'T' && buf[1] == 'Z' && buf[2] == 'i' && buf[3] == 'f' {
		return true
	}

	if buf[0] == 'P' && buf[1] == 'H' && buf[2] == 'P' {
		return true
	}

	return false
}

// LoadTzFileFromDB loads timezone file from database
func LoadTzFileFromDB(tzName string, tzdb *TzDB) ([]byte, string, error) {
	if tzdb == nil {
		return nil, "", fmt.Errorf("timezone database is nil")
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
		return nil, "", fmt.Errorf("timezone '%s' not found in database", tzName)
	}

	// If database has embedded data
	if len(tzdb.Data) > 0 && entry.Pos < len(tzdb.Data) {
		// Find the end of this entry by finding the next entry's position
		endPos := len(tzdb.Data)
		for i := range tzdb.Index {
			if tzdb.Index[i].Pos > entry.Pos && tzdb.Index[i].Pos < endPos {
				endPos = tzdb.Index[i].Pos
			}
		}

		// Extract data for this timezone
		if entry.Pos < endPos && endPos <= len(tzdb.Data) {
			data := tzdb.Data[entry.Pos:endPos]
			return data, tzName, nil
		}

		return nil, "", fmt.Errorf("invalid timezone data boundaries")
	}

	// Load from file - ID contains the relative path
	// First try BaseDir if available
	if tzdb.BaseDir != "" {
		fullPath := filepath.Join(tzdb.BaseDir, filepath.FromSlash(entry.ID))
		data, err := os.ReadFile(fullPath)
		if err == nil {
			return data, fullPath, nil
		}
	}

	// Fall back to system paths and direct path
	possiblePaths := []string{
		entry.ID,
		filepath.Join("/usr/share/zoneinfo", entry.ID),
		filepath.Join("/var/db/timezone/zoneinfo", entry.ID),
	}

	for _, path := range possiblePaths {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, path, nil
		}
	}

	return nil, "", fmt.Errorf("cannot find timezone file for '%s'", tzName)
}

// UpdateParseTzfile updates the ParseTzfile function to use the new parser
func UpdateParseTzfile(timezone string, tzdb *TzDB, errorCode *int) (*TzInfo, error) {
	// Try to load from database first
	if tzdb != nil {
		// Search in index
		var found *TzDBIndexEntry
		for i := range tzdb.Index {
			if tzdb.Index[i].ID == timezone {
				found = &tzdb.Index[i]
				break
			}
		}

		if found != nil {
			// Try to load timezone file data
			data, _, err := LoadTzFileFromDB(timezone, tzdb)
			if err == nil {
				return ParseTzfileData(timezone, data, errorCode)
			}
		}
	}

	// Fall back to system paths
	systemPaths := []string{
		filepath.Join("/usr/share/zoneinfo", timezone),
		filepath.Join("/var/db/timezone/zoneinfo", timezone),
		timezone, // Try as direct path
	}

	for _, path := range systemPaths {
		tz, err := ParseTzfileFromFile(path, errorCode)
		if err == nil {
			tz.Name = timezone // Use requested name, not filename
			return tz, nil
		}
	}

	if errorCode != nil {
		*errorCode = TIMELIB_ERROR_NO_SUCH_TIMEZONE
	}
	return nil, fmt.Errorf("timezone '%s' not found", timezone)
}
