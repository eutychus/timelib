# timelib-go

Go port of the PHP timelib library for date/time parsing and timezone handling.

## Features

- **Complete date/time parser** supporting multiple formats (ISO 8601, relative times, natural language, etc.)
- **Full timezone support** with builtin timezone database (598 timezones, version 2025.2)
- **Timezone abbreviation resolution** (1,127 abbreviations including military timezones)
- **Interval parsing** for ISO 8601 durations and recurring intervals
- **Timezone transitions** with historical DST change tracking
- **Arithmetic operations** on dates and times

## Installation

```bash
go get github.com/eutychus/timelib
```

## Basic Usage

```go
package main

import (
    "fmt"
    timelib "github.com/eutychus/timelib"
)

func main() {
    // Parse a date string
    t, err := timelib.StrToTime("2025-10-04 15:30:00 EST", timelib.BuiltinDB())
    if err != nil {
        panic(err)
    }

    fmt.Printf("Year: %d, Month: %d, Day: %d\n", t.Y, t.M, t.D)
    fmt.Printf("Hour: %d, Minute: %d, Second: %d\n", t.H, t.I, t.S)

    // Clean up
    timelib.TimeDtor(t)
}
```

## Updating the Timezone Database

The timezone database is embedded in the library as binary files (`tzdata_index.bin` and `tzdata_data.bin`). To update the timezone database when a new version is released:

### Prerequisites

- Access to the C timelib source code (specifically `timezonedb.h`)
- GCC compiler

### Quick Reference

```bash
# 1. Navigate to the C timelib directory
cd /path/to/timelib

# 2. Update timezonedb.h from upstream (if needed)
# wget https://github.com/derickr/timelib/raw/master/timezonedb.h

# 3. Verify extract_tzdata.c includes (should be "timelib.h", not "../timelib.h")
grep "include" extract_tzdata.c

# 4. Compile and run the extraction tool
gcc -o extract_tzdata extract_tzdata.c
./extract_tzdata

# 5. Copy generated files to Go package
cp tzdata_index.bin tzdata_data.bin /path/to/timelib-go/

# 6. Verify in Go
cd /path/to/timelib-go
go test ./tests
```

### Detailed Steps

1. **Update the C timezone database header**

   Obtain the latest `timezonedb.h` from the upstream timelib project and place it in the parent timelib directory:
   ```bash
   cd /path/to/timelib
   # Download from: https://github.com/derickr/timelib/blob/master/timezonedb.h
   ```

2. **Verify the extract tool includes**

   **Important**: Check that `extract_tzdata.c` has the correct includes:
   ```bash
   grep "#include" extract_tzdata.c
   ```

   Should show:
   ```c
   #include "timelib.h"
   #include "timezonedb.h"
   ```

   If it shows `#include "../timelib.h"`, fix it:
   ```bash
   sed -i 's|#include "../timelib.h"|#include "timelib.h"|' extract_tzdata.c
   sed -i 's|#include "../timezonedb.h"|#include "timezonedb.h"|' extract_tzdata.c
   ```

3. **Compile the extraction tool**

   ```bash
   cd /path/to/timelib
   gcc -o extract_tzdata extract_tzdata.c
   ```

4. **Extract the timezone database**

   Run the extraction tool to generate the binary files:
   ```bash
   ./extract_tzdata
   ```

   Expected output:
   ```
   Extracted timezone database:
     Index: 598 entries
     Data: 355607 bytes
     Version: 2025.2
   ```

   This creates:
   - `tzdata_index.bin` - Index of 598 timezone IDs and their positions
   - `tzdata_data.bin` - Binary timezone transition data (355,607 bytes)

5. **Copy the files to the Go package**

   ```bash
   cp tzdata_index.bin tzdata_data.bin /path/to/timelib-go/
   ```

6. **Verify the update**

   Run the tests to ensure the new database works correctly:
   ```bash
   cd /path/to/timelib-go
   go test ./tests
   ```

   Should output: `ok github.com/eutychus/timelib/tests`

7. **Test the database programmatically**

   Create and run a test program:
   ```go
   package main

   import (
       "fmt"
       timelib "github.com/eutychus/timelib"
   )

   func main() {
       db := timelib.BuiltinDB()
       fmt.Printf("Database version: %s\n", db.Version)
       fmt.Printf("Total timezones: %d\n", db.IndexSize)

       // Test critical timezone lookups
       testIDs := []string{"UTC", "GMT", "EST", "America/New_York", "Europe/London"}
       for _, id := range testIDs {
           entry := timelib.FindTimezone(id)
           if entry != nil {
               fmt.Printf("✓ Found %s\n", id)
           } else {
               fmt.Printf("✗ Missing %s\n", id)
           }
       }
   }
   ```

   All test timezones should be found.

### Database Structure

The timezone database consists of two embedded binary files:

- **`tzdata_index.bin`**: Contains 598 timezone index entries. Each entry has:
  - Timezone ID (null-terminated string, e.g., "America/New_York")
  - Position offset (4-byte big-endian integer pointing into tzdata_data.bin)

- **`tzdata_data.bin`**: Contains the raw timezone transition data for all timezones in the compiled TZif format.

The index is parsed at runtime by `parseBuiltinIndex()` in `tzdata_builtin.go`, and timezone data is accessed on-demand from the embedded data.

### Timezone Abbreviations

Timezone abbreviations (like EST, PST, CET) are defined in `timezone_abbr.go` and are ported from the C timelib `timezonemap.h` file. These generally don't need updating as frequently as the timezone database itself, but if new abbreviations are added:

1. Update the abbreviation table by regenerating from the C `timezonemap.h`
2. Ensure military timezone abbreviations (A-Z) remain included
3. Run tests to verify abbreviation resolution works correctly

### Detailed Documentation

For more detailed information about the timezone database, see:
- **TIMEZONE_UPDATE_NOTES.md** - Comprehensive update notes, troubleshooting, and verification procedures

## Testing

```bash
# Run all tests
go test ./tests

# Run specific test suites
go test ./tests -run TestTimezone
go test ./tests -run TestParse

# Run with verbose output
go test ./tests -v
```

### Test Status

- **All tests**: ✅ 667 tests passing, 0 failures, 0 skipped
- **Success rate**: 100%
- **Test coverage**: Comprehensive coverage of date/time parsing, timezone handling, arithmetic operations, and interval parsing

The project is fully functional and ready for production use.

## Documentation

- **Parser architecture**: See source code comments and test files for implementation details
- **Timezone implementation**: Built-in timezone database with 598 timezones (version 2025.2)

## License

This is a port of the PHP timelib library. Please refer to the original timelib license.

## Credits

Original C library: https://github.com/derickr/timelib
