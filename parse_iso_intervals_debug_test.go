package timelib

import (
	"fmt"
	"testing"
)

func TestIsoIntervalDebug(t *testing.T) {
	input := "2008-03-01T13:00:00Z"
	begin, end, period, recur, errors := ParseIsoInterval(input)

	fmt.Printf("Input: %q\n", input)
	if begin != nil {
		fmt.Printf("Begin: Y=%d M=%d D=%d H=%d I=%d S=%d\n", begin.Y, begin.M, begin.D, begin.H, begin.I, begin.S)
	} else {
		fmt.Println("Begin: nil")
	}
	if end != nil {
		fmt.Printf("End: Y=%d M=%d D=%d H=%d I=%d S=%d\n", end.Y, end.M, end.D, end.H, end.I, end.S)
	} else {
		fmt.Println("End: nil")
	}
	if period != nil {
		fmt.Printf("Period: Y=%d M=%d D=%d H=%d I=%d S=%d\n", period.Y, period.M, period.D, period.H, period.I, period.S)
	} else {
		fmt.Println("Period: nil")
	}
	fmt.Printf("Recurrences: %d\n", recur)
	if errors != nil && errors.ErrorCount > 0 {
		for _, e := range errors.ErrorMessages {
			fmt.Printf("Error: %s at position %d\n", e.Message, e.Position)
		}
	}
}
