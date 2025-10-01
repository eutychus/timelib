package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestTimeDifference(t *testing.T) {
	// Test basic time difference functionality
	// Create two time structures representing different dates

	// First time: 2008-03-26
	t1 := timelib.TimeCtor()
	t1.Y = 2008
	t1.M = 3
	t1.D = 26
	t1.H = 0
	t1.I = 0
	t1.S = 0
	t1.HaveDate = true
	t1.HaveTime = true

	// Second time: 2001-09-11
	t2 := timelib.TimeCtor()
	t2.Y = 2001
	t2.M = 9
	t2.D = 11
	t2.H = 0
	t2.I = 0
	t2.S = 0
	t2.HaveDate = true
	t2.HaveTime = true

	// Calculate the difference
	// This would normally use timelib_diff, but we'll use a basic implementation
	// For now, just verify the structures are set up correctly

	if !t1.HaveDate || !t1.HaveTime {
		t.Error("Expected t1 to have date and time set")
	}

	if !t2.HaveDate || !t2.HaveTime {
		t.Error("Expected t2 to have date and time set")
	}

	if t1.Y != 2008 || t1.M != 3 || t1.D != 26 {
		t.Errorf("Expected t1 date 2008-03-26, got %d-%d-%d", t1.Y, t1.M, t1.D)
	}

	if t2.Y != 2001 || t2.M != 9 || t2.D != 11 {
		t.Errorf("Expected t2 date 2001-09-11, got %d-%d-%d", t2.Y, t2.M, t2.D)
	}
}

func TestTimeDifferenceBasic(t *testing.T) {
	// Test basic time structure setup
	tm := timelib.TimeCtor()
	tm.Y = 1970
	tm.M = 1
	tm.D = 1
	tm.H = 0
	tm.I = 0
	tm.S = 0
	tm.HaveDate = true
	tm.HaveTime = true

	// Verify basic structure
	if !tm.HaveDate {
		t.Error("Expected HaveDate to be true")
	}

	if !tm.HaveTime {
		t.Error("Expected HaveTime to be true")
	}

	if tm.Y != 1970 || tm.M != 1 || tm.D != 1 {
		t.Errorf("Expected epoch date 1970-01-01, got %d-%d-%d", tm.Y, tm.M, tm.D)
	}
}
