package tests

import (
	"math"
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestPHPSunInfo001(t *testing.T) {
	time := &timelib.Time{
		Y: 2006,
		M: 12,
		D: 12,
		H: 0,
		I: 0,
		S: 0,
	}
	time.UpdateTS(nil)

	hRise, hSet, tsRise, tsSet, tsTransit := timelib.AstroRiseSetAltitude(
		time, 31.7667, 35.2333, -35.0/60.0, 1,
	)

	if math.Abs(hRise-4.86) > 0.01 {
		t.Errorf("h_rise: expected 4.86, got %f", hRise)
	}
	if math.Abs(hSet-14.69) > 0.01 {
		t.Errorf("h_set: expected 14.69, got %f", hSet)
	}
	if tsRise != 1165899111 {
		t.Errorf("ts_rise: expected 1165899111, got %d", tsRise)
	}
	if tsSet != 1165934475 {
		t.Errorf("ts_set: expected 1165934475, got %d", tsSet)
	}
	expectedTransit := int64((1165899111 + 1165934475) / 2)
	if tsTransit != expectedTransit {
		t.Errorf("ts_transit: expected %d, got %d", expectedTransit, tsTransit)
	}
}

func TestPHPSunInfo002(t *testing.T) {
	time := &timelib.Time{
		Y: 2007,
		M: 4,
		D: 13,
		H: 11,
		I: 10,
		S: 54,
	}
	time.UpdateTS(nil)

	hRise, hSet, tsRise, tsSet, tsTransit := timelib.AstroRiseSetAltitude(
		time, 9.61, 59.21, -35.0/60.0, 1,
	)

	if math.Abs(hRise-4.23) > 0.01 {
		t.Errorf("h_rise: expected 4.23, got %f", hRise)
	}
	if math.Abs(hSet-18.51) > 0.01 {
		t.Errorf("h_set: expected 18.51, got %f", hSet)
	}
	if tsRise != 1176437611 {
		t.Errorf("ts_rise: expected 1176437611, got %d", tsRise)
	}
	if tsSet != 1176489051 {
		t.Errorf("ts_set: expected 1176489051, got %d", tsSet)
	}
	expectedTransit := int64((1176489051 + 1176437611) / 2)
	if tsTransit != expectedTransit {
		t.Errorf("ts_transit: expected %d, got %d", expectedTransit, tsTransit)
	}
}
