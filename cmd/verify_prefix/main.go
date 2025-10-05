package main

import (
	"fmt"

	timelib "github.com/eutychus/timelib"
)

func main() {
	// Test with prefix (matching C version's test_parse_with_prefix)
	fmt.Println("Testing with % prefix (matching C version):")
	result, errors := timelib.ParseFromFormatWithPrefix("%Y-%m-%dT%H:%i:%sZ", "2018-01-26T11:56:02Z")
	if errors != nil && errors.ErrorCount > 0 {
		fmt.Printf("  Error: %v\n", errors)
	} else {
		fmt.Printf("  Year: %d, Month: %d, Day: %d\n", result.Y, result.M, result.D)
		fmt.Printf("  Hour: %d, Minute: %d, Second: %d\n", result.H, result.I, result.S)
	}

	fmt.Println("\nTesting without prefix (original Go style):")
	result2, errors2 := timelib.ParseFromFormat("Y/m/d", "2018/01/26")
	if errors2 != nil && errors2.ErrorCount > 0 {
		fmt.Printf("  Error: %v\n", errors2)
	} else {
		fmt.Printf("  Year: %d, Month: %d, Day: %d\n", result2.Y, result2.M, result2.D)
	}

	fmt.Println("\nTesting prefix with ISO week format:")
	result3, errors3 := timelib.ParseFromFormatWithPrefix("%V.%b.%B", "53.7.2017")
	if errors3 != nil && errors3.ErrorCount > 0 {
		fmt.Printf("  Error: %v\n", errors3)
	} else {
		fmt.Printf("  Year: %d, Month: %d, Day: %d\n", result3.Y, result3.M, result3.D)
	}
}
