package main

import (
	"fmt"
	"log"
	"time"

	timelib "github.com/eutychus/timelib"
)

func main() {
	fmt.Println("=== timelib PHP Compatibility Functions Demo ===")
	fmt.Println()

	// Example 1: Parse absolute date to Unix timestamp
	fmt.Println("1. Parse absolute date to Unix timestamp:")
	ts := timelib.Strtotime("2024-01-15", 0)
	if ts == -1 {
		log.Fatal("Failed to parse date")
	}
	fmt.Printf("   Input: \"2024-01-15\"\n")
	fmt.Printf("   Unix Timestamp: %d\n", ts)
	fmt.Printf("   Formatted: %s\n\n", time.Unix(ts, 0).UTC().Format(time.RFC3339))

	// Example 2: Parse with time component
	fmt.Println("2. Parse datetime string:")
	ts2 := timelib.Strtotime("2024-06-15 10:30:45", 0)
	fmt.Printf("   Input: \"2024-06-15 10:30:45\"\n")
	fmt.Printf("   Unix Timestamp: %d\n", ts2)
	fmt.Printf("   Formatted: %s\n\n", time.Unix(ts2, 0).UTC().Format(time.RFC3339))

	// Example 3: Relative time - tomorrow
	fmt.Println("3. Calculate tomorrow from specific base time:")
	baseTime := int64(1704067200) // 2024-01-01 00:00:00 UTC
	tomorrow := timelib.Strtotime("tomorrow", baseTime)
	fmt.Printf("   Base time: %s\n", time.Unix(baseTime, 0).UTC().Format("2006-01-02"))
	fmt.Printf("   Tomorrow: %s\n\n", time.Unix(tomorrow, 0).UTC().Format("2006-01-02"))

	// Example 4: Add days to a date
	fmt.Println("4. Add 7 days to a specific date:")
	baseTime2 := int64(1704067200) // 2024-01-01
	plusWeek := timelib.Strtotime("+7 days", baseTime2)
	fmt.Printf("   Base: %s\n", time.Unix(baseTime2, 0).UTC().Format("2006-01-02"))
	fmt.Printf("   +7 days: %s\n\n", time.Unix(plusWeek, 0).UTC().Format("2006-01-02"))

	// Example 5: Get Go time.Time instead of Unix timestamp
	fmt.Println("5. Parse to Go time.Time:")
	goTime, err := timelib.StrtotimeToGoTime("2024-12-25 09:00:00", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Input: \"2024-12-25 09:00:00\"\n")
	fmt.Printf("   time.Time: %s\n", goTime.Format(time.RFC3339))
	fmt.Printf("   Formatted: %s\n\n", goTime.Format("Monday, January 2, 2006 at 3:04 PM"))

	// Example 6: Parse with base time as time.Time
	fmt.Println("6. Parse relative time with time.Time base:")
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nextWeek, err := timelib.StrtotimeToGoTime("+1 week", &base)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Base: %s\n", base.Format("2006-01-02"))
	fmt.Printf("   +1 week: %s\n\n", nextWeek.Format("2006-01-02"))

	// Example 7: Parse with specific timezone
	fmt.Println("7. Parse with timezone:")
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	nyTime, err := timelib.StrtotimeWithTimezone("2024-01-15 10:30:00", nil, loc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Input: \"2024-01-15 10:30:00\"\n")
	fmt.Printf("   In UTC: %s\n", nyTime.UTC().Format(time.RFC3339))
	fmt.Printf("   In New York: %s\n\n", nyTime.Format(time.RFC3339))

	// Example 8: Calculate expiration time
	fmt.Println("8. Calculate expiration time (30 days from now):")
	now := time.Now().Unix()
	expiration := timelib.Strtotime("+30 days", now)
	expirationTime := time.Unix(expiration, 0)
	fmt.Printf("   Current time: %s\n", time.Unix(now, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("   Expires at: %s\n", expirationTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Days until expiration: %.0f\n\n", expirationTime.Sub(time.Unix(now, 0)).Hours()/24)

	// Example 9: Day-of-week relative times
	fmt.Println("9. Day-of-week relative times:")
	base2 := int64(1704067200) // 2024-01-01 (Monday)
	nextMon := timelib.Strtotime("next monday", base2)
	fmt.Printf("   Base (2024-01-01): %s\n", time.Unix(base2, 0).UTC().Weekday())
	fmt.Printf("   Next Monday: %s (%s)\n\n",
		time.Unix(nextMon, 0).UTC().Format("2006-01-02"),
		time.Unix(nextMon, 0).UTC().Weekday())

	// Example 10: Error handling
	fmt.Println("10. Error handling:")
	invalidTs := timelib.Strtotime("not a valid date", 0)
	if invalidTs == -1 {
		fmt.Println("   Successfully caught invalid date string!")
	}

	_, err = timelib.StrtotimeToGoTime("also invalid", nil)
	if err != nil {
		fmt.Printf("   Error message: %v\n\n", err)
	}

	fmt.Println("=== Demo Complete ===")
}
