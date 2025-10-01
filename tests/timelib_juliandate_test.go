package tests

import (
	"math"
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestJ2000Epoch(t *testing.T) {
	d := int64(946728000)
	js := timelib.TsToJ2000(d)

	if math.Abs(js-0) > 0.0000001 {
		t.Errorf("J2000Epoch: expected 0, got %f", js)
	}
}

func TestAugust2017_J2000(t *testing.T) {
	d := int64(1502755200)
	js := timelib.TsToJ2000(d)

	if math.Abs(js-6435.5) > 0.0000001 {
		t.Errorf("August2017: expected 6435.5, got %f", js)
	}
}

func TestJulianDayEpoch(t *testing.T) {
	d := int64(-210866760000)
	js := timelib.TsToJulianDay(d)

	if math.Abs(js-0) > 0.0000001 {
		t.Errorf("JulianDayEpoch: expected 0, got %f", js)
	}
}

func TestJulianDateExampleFromWikipedia(t *testing.T) {
	d := int64(1357000200)
	js := timelib.TsToJulianDay(d)

	if math.Abs(js-2456293.520833) > 0.000001 {
		t.Errorf("JulianDateExampleFromWikipedia: expected 2456293.520833, got %f", js)
	}
}

func TestAugust2017_JulianDay(t *testing.T) {
	d := int64(1502755200)
	js := timelib.TsToJulianDay(d)

	if math.Abs(js-2457980.5) > 0.0000001 {
		t.Errorf("August2017: expected 2457980.5, got %f", js)
	}
}
