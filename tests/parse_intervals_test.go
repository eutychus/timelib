package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestParseIntervalsWeeksOnly(t *testing.T) {
	errors := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval("P2W", errors)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if errors.WarningCount != 0 || errors.ErrorCount != 0 {
		t.Errorf("Expected no errors, got %d warnings, %d errors", errors.WarningCount, errors.ErrorCount)
	}
	if p.Y != 0 || p.M != 0 || p.D != 14 {
		t.Errorf("Expected 0Y 0M 14D, got %dY %dM %dD", p.Y, p.M, p.D)
	}
}

func TestParseIntervalsCombiningWeeksAndDays(t *testing.T) {
	errors := &timelib.ErrorContainer{}
	_, _, p, _, err := timelib.Strtointerval("P2W3D", errors)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if errors.WarningCount != 0 || errors.ErrorCount != 0 {
		t.Errorf("Expected no errors, got %d warnings, %d errors", errors.WarningCount, errors.ErrorCount)
	}
	if p.Y != 0 || p.M != 0 || p.D != 17 {
		t.Errorf("Expected 0Y 0M 17D, got %dY %dM %dD", p.Y, p.M, p.D)
	}

	errors = &timelib.ErrorContainer{}
	_, _, p, _, err = timelib.Strtointerval("P1Y3M1W5DT2H", errors)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if errors.WarningCount != 0 || errors.ErrorCount != 0 {
		t.Errorf("Expected no errors, got %d warnings, %d errors", errors.WarningCount, errors.ErrorCount)
	}
	if p.Y != 1 || p.M != 3 || p.D != 12 || p.H != 2 {
		t.Errorf("Expected 1Y 3M 12D 2H, got %dY %dM %dD %dH", p.Y, p.M, p.D, p.H)
	}
}
