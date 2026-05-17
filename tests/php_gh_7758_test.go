package tests

import (
	"testing"
	"github.com/eutychus/timelib"
)

func TestPhpGh7758(t *testing.T) {
	time, _, err := timelib.ParseDateString("@-0.4", timelib.BuiltinDB(), timelib.ParseTzfile)
	if err != nil {
		t.Fatalf("ParseDateString returned error: %v", err)
	}
	if time == nil {
		t.Fatalf("ParseDateString returned nil time")
	}

	if time.Y != 1970 { t.Errorf("Expected Y=1970, got %d", time.Y) }
	if time.M != 1 { t.Errorf("Expected M=1, got %d", time.M) }
	if time.D != 1 { t.Errorf("Expected D=1, got %d", time.D) }
	if time.H != 0 { t.Errorf("Expected H=0, got %d", time.H) }
	if time.I != 0 { t.Errorf("Expected I=0, got %d", time.I) }
	if time.S != 0 { t.Errorf("Expected S=0, got %d", time.S) }
	if time.Relative.S != 0 { t.Errorf("Expected Relative.S=0, got %d", time.Relative.S) }
	if time.Relative.US != -400000 { t.Errorf("Expected Relative.US=-400000, got %d", time.Relative.US) }
}
