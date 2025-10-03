package timelib

import (
	"fmt"
	"testing"
	"unsafe"
)

// TestDebugScanner helps debug the scanner behavior
func TestDebugScanner(t *testing.T) {
	str := "2025"

	// Initialize scanner with null terminator
	strBytes := make([]byte, len(str)+1)
	copy(strBytes, str)
	strBytes[len(str)] = 0

	s := &Scanner{
		str:    strBytes,
		errors: &ErrorContainer{},
		time: &Time{
			Y:  TIMELIB_UNSET,
			M:  TIMELIB_UNSET,
			D:  TIMELIB_UNSET,
			H:  TIMELIB_UNSET,
			I:  TIMELIB_UNSET,
			S:  TIMELIB_UNSET,
			US: TIMELIB_UNSET,
			Z:  TIMELIB_UNSET,
			Dst: TIMELIB_UNSET,
		},
		tzdb: nil,
	}

	// Set up pointers
	s.cur = &s.str[0]
	s.lim = &s.str[len(str)]

	fmt.Printf("Input string: %q\n", str)
	fmt.Printf("Input bytes: %v\n", s.str)
	fmt.Printf("Scanner initialized:\n")
	fmt.Printf("  str: %p (len=%d)\n", &s.str[0], len(s.str))
	fmt.Printf("  cur: %p (value=%d '%c')\n", s.cur, *s.cur, *s.cur)
	fmt.Printf("  lim: %p\n", s.lim)
	fmt.Printf("  Distance to limit: %d bytes\n", uintptr(unsafe.Pointer(s.lim))-uintptr(unsafe.Pointer(s.cur)))

	// Run scanner with debug
	maxIter := 10
	for i := 0; i < maxIter; i++ {
		fmt.Printf("\n=== Iteration %d ===\n", i+1)
		fmt.Printf("Before scan:\n")
		fmt.Printf("  cur: %p", s.cur)
		if s.cur != nil && uintptr(unsafe.Pointer(s.cur)) < uintptr(unsafe.Pointer(s.lim)) {
			fmt.Printf(" (value=%d '%c')\n", *s.cur, *s.cur)
		} else {
			fmt.Printf(" (past limit)\n")
		}

		token := scan(s, nil)

		fmt.Printf("After scan:\n")
		fmt.Printf("  token: %d\n", token)
		fmt.Printf("  cur: %p", s.cur)
		if s.cur != nil && uintptr(unsafe.Pointer(s.cur)) < uintptr(unsafe.Pointer(s.lim)) {
			fmt.Printf(" (value=%d '%c')\n", *s.cur, *s.cur)
		} else {
			fmt.Printf(" (past limit)\n")
		}
		fmt.Printf("  Time.Y: %d\n", s.time.Y)

		if token == EOI {
			fmt.Printf("\nReached EOI\n")
			break
		}
	}

	fmt.Printf("\n=== Final State ===\n")
	fmt.Printf("Time.Y: %d\n", s.time.Y)
	fmt.Printf("Time.M: %d\n", s.time.M)
	fmt.Printf("Time.D: %d\n", s.time.D)
	fmt.Printf("Time.H: %d\n", s.time.H)
	fmt.Printf("Time.I: %d\n", s.time.I)
	fmt.Printf("Time.S: %d\n", s.time.S)
	fmt.Printf("HaveDate: %v\n", s.time.HaveDate)
	fmt.Printf("HaveTime: %v\n", s.time.HaveTime)
	fmt.Printf("HaveRelative: %v\n", s.time.HaveRelative)
	fmt.Printf("Errors: %d\n", s.errors.ErrorCount)
	if s.errors.ErrorCount > 0 {
		for i, e := range s.errors.ErrorMessages {
			fmt.Printf("  Error %d: %s at position %d\n", i, e.Message, e.Position)
		}
	}
}

// TestParseSimple tests simple parsing
func TestParseSimple(t *testing.T) {
	testCases := []string{"2025", "now", "yesterday"}
	for _, input := range testCases {
		time, err := StrToTime(input, nil)
		if err != nil {
			t.Fatalf("Error parsing %q: %v", input, err)
		}
		t.Logf("%q => Y=%d M=%d D=%d H=%d I=%d S=%d", input, time.Y, time.M, time.D, time.H, time.I, time.S)
		t.Logf("  HaveDate=%v HaveTime=%v HaveRelative=%v", time.HaveDate, time.HaveTime, time.HaveRelative)
		if time.HaveRelative {
			t.Logf("  Relative: Y=%d M=%d D=%d H=%d I=%d S=%d",
				time.Relative.Y, time.Relative.M, time.Relative.D,
				time.Relative.H, time.Relative.I, time.Relative.S)
		}
	}
}
