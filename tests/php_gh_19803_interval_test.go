package tests

import (
	"testing"
	"github.com/eutychus/timelib"
)

func TestPhpGh19803_Interval_01(t *testing.T) {
	_, _, _, _, errors := timelib.ParseIsoInterval("")
	if errors == nil || errors.ErrorCount != 1 {
		t.Fatalf("Expected 1 error, got %v", errors)
	}
	if errors.ErrorMessages[0].Message != "Empty string" {
		t.Errorf("Expected 'Empty string', got '%s'", errors.ErrorMessages[0].Message)
	}
}

func TestPhpGh19803_Interval_02(t *testing.T) {
	_, _, _, _, errors := timelib.ParseIsoInterval("  ")
	if errors == nil || errors.ErrorCount != 1 {
		t.Fatalf("Expected 1 error, got %v", errors)
	}
	if errors.ErrorMessages[0].Message != "Empty string" {
		t.Errorf("Expected 'Empty string', got '%s'", errors.ErrorMessages[0].Message)
	}
}
