package tests

import (
	"math"
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestConvertPositiveHMSToDecimal(t *testing.T) {
	hour, min, sec := 2, 19, 48
	d := timelib.HMSToDecimalHour(hour, min, sec)

	if math.Abs(d-2.33) > 0.00001 {
		t.Errorf("Expected 2.33, got %f", d)
	}
}

func TestConvertZeroHMSToDecimal(t *testing.T) {
	hour, min, sec := int(0), int(0), int(0)
	d := timelib.HMSToDecimalHour(hour, min, sec)

	if math.Abs(d-0) > 0.00001 {
		t.Errorf("Expected 0, got %f", d)
	}
}

func TestConvertNegativeHMSToDecimal(t *testing.T) {
	hour, min, sec := int(-2), int(20), int(0)
	d := timelib.HMSToDecimalHour(hour, min, sec)

	if math.Abs(d-(-2.333333)) > 0.000001 {
		t.Errorf("Expected -2.333333, got %f", d)
	}
}

func TestConvertNegativeZeroHMSToDecimal(t *testing.T) {
	hour, min, sec := int(0), int(0), int(0)
	d := timelib.HMSToDecimalHour(hour, min, sec)

	if math.Abs(d-0) > 0.00001 {
		t.Errorf("Expected 0, got %f", d)
	}
}

func TestConvertPositiveDecimalToHMS(t *testing.T) {
	d := 2.33
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != 2 || min != 19 || sec != 48 {
		t.Errorf("Expected 2:19:48, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertZeroDecimalToHMS(t *testing.T) {
	d := 0.0
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != 0 || min != 0 || sec != 0 {
		t.Errorf("Expected 0:0:0, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertNegativeDecimalToHMS(t *testing.T) {
	d := -2.33
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != -2 || min != 19 || sec != 48 {
		t.Errorf("Expected -2:19:48, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertNegativeZeroDecimalToHMS(t *testing.T) {
	d := 0.0
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != 0 || min != 0 || sec != 0 {
		t.Errorf("Expected 0:0:0, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertPositiveHMSFToDecimal(t *testing.T) {
	hour, min, sec, usec := int(2), int(19), int(48), int(250000)
	d := timelib.HMSFToDecimalHour(hour, min, sec, usec)

	if math.Abs(d-2.330069) > 0.000001 {
		t.Errorf("Expected 2.330069, got %f", d)
	}
}

func TestConvertPositiveHMSFToDecimalFullSec(t *testing.T) {
	hour, min, sec, usec := int(2), int(19), int(47), int(1000000)
	d := timelib.HMSFToDecimalHour(hour, min, sec, usec)

	if math.Abs(d-2.33) > 0.000001 {
		t.Errorf("Expected 2.33, got %f", d)
	}
}

func TestConvertZeroHMSFToDecimal(t *testing.T) {
	hour, min, sec, usec := int(0), int(0), int(0), int(0)
	d := timelib.HMSFToDecimalHour(hour, min, sec, usec)

	if math.Abs(d-0) > 0.00001 {
		t.Errorf("Expected 0, got %f", d)
	}
}

func TestConvertNegativeHMSFToDecimal(t *testing.T) {
	hour, min, sec, usec := int(-2), int(20), int(0), int(50000)
	d := timelib.HMSFToDecimalHour(hour, min, sec, usec)

	if math.Abs(d-(-2.333347)) > 0.000001 {
		t.Errorf("Expected -2.333347, got %f", d)
	}
}

func TestConvertNegativeHMSFToDecimalSmall(t *testing.T) {
	hour, min, sec, usec := int(-2), int(20), int(0), int(5)
	base := timelib.HMSToDecimalHour(hour, min, sec)
	d := timelib.HMSFToDecimalHour(hour, min, sec, usec)

	if math.Abs((d-base)-(-1.388889e-09)) > 0.00001 {
		t.Errorf("Expected difference -1.388889e-09, got %e", d-base)
	}
}

func TestConvertNegativeZeroHMSFToDecimal(t *testing.T) {
	hour, min, sec, usec := int(0), int(0), int(0), int(10)
	d := timelib.HMSFToDecimalHour(hour, min, sec, usec)

	if math.Abs(d-0) > 0.00001 {
		t.Errorf("Expected 0, got %f", d)
	}
}

func TestConvertPositiveDecimalToHMSOverflow15(t *testing.T) {
	d := 9.333333333333333
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != 9 || min != 19 || sec != 59 {
		t.Errorf("Expected 9:19:59, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertPositiveDecimalToHMSOverflow16(t *testing.T) {
	d := 9.3333333333333333
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != 9 || min != 20 || sec != 0 {
		t.Errorf("Expected 9:20:0, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertNegativeDecimalToHMSOverflow15(t *testing.T) {
	d := -9.333333333333333
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != -9 || min != 19 || sec != 59 {
		t.Errorf("Expected -9:19:59, got %d:%d:%d", hour, min, sec)
	}
}

func TestConvertNegativeDecimalToHMSOverflow16(t *testing.T) {
	d := -9.3333333333333333
	hour, min, sec := timelib.DecimalHourToHMS(d)

	if hour != -9 || min != 20 || sec != 0 {
		t.Errorf("Expected -9:20:0, got %d:%d:%d", hour, min, sec)
	}
}

// Additional tests from timelib_hmsf_to_decimal_hour.cpp

func TestHMSFZero(t *testing.T) {
	d := timelib.HMSFToDecimalHour(0, 0, 0, 0)
	if math.Abs(d-0) > 0.00000000001 {
		t.Errorf("Expected 0, got %f", d)
	}
}

func TestHMSFSmallestPositive(t *testing.T) {
	d := timelib.HMSFToDecimalHour(0, 0, 0, 1)
	if math.Abs(d-2.777778e-10) > 0.00000000001 {
		t.Errorf("Expected 2.777778e-10, got %e", d)
	}
}

func TestHMSFOneSecondPositive(t *testing.T) {
	d := timelib.HMSFToDecimalHour(0, 0, 1, 0)
	if math.Abs(d-0.00027777778) > 0.00000000001 {
		t.Errorf("Expected 0.00027777778, got %f", d)
	}
}

func TestHMSFSixMinutePositive(t *testing.T) {
	d := timelib.HMSFToDecimalHour(0, 6, 0, 0)
	if math.Abs(d-0.1) > 0.00000000001 {
		t.Errorf("Expected 0.1, got %f", d)
	}
}

func TestHMSFThreeHoursPositive(t *testing.T) {
	d := timelib.HMSFToDecimalHour(3, 0, 0, 0)
	if math.Abs(d-3) > 0.00000000001 {
		t.Errorf("Expected 3, got %f", d)
	}
}

func TestHMSFThreeHoursFifteenMinutesPositive(t *testing.T) {
	d := timelib.HMSFToDecimalHour(3, 15, 0, 0)
	if math.Abs(d-3.25) > 0.00000000001 {
		t.Errorf("Expected 3.25, got %f", d)
	}
}

func TestHMSFOneSecondNegative(t *testing.T) {
	d := timelib.HMSFToDecimalHour(0, 0, -1, 0)
	if math.Abs(d-(-0.00027777778)) > 0.00000000001 {
		t.Errorf("Expected -0.00027777778, got %f", d)
	}
}

func TestHMSFSixMinuteNegative(t *testing.T) {
	d := timelib.HMSFToDecimalHour(0, -6, 0, 0)
	if math.Abs(d-(-0.1)) > 0.00000000001 {
		t.Errorf("Expected -0.1, got %f", d)
	}
}

func TestHMSFThreeHoursNegative(t *testing.T) {
	d := timelib.HMSFToDecimalHour(-3, 0, 0, 0)
	if math.Abs(d-(-3)) > 0.00000000001 {
		t.Errorf("Expected -3, got %f", d)
	}
}

func TestHMSFThreeHoursFifteenMinutesNegative(t *testing.T) {
	d := timelib.HMSFToDecimalHour(-3, 15, 0, 0)
	if math.Abs(d-(-3.25)) > 0.00000000001 {
		t.Errorf("Expected -3.25, got %f", d)
	}
}

func TestHMSFThreeHoursNegativeFifteenMinutesNegative(t *testing.T) {
	d := timelib.HMSFToDecimalHour(-3, -15, 0, 0)
	if math.Abs(d-(-2.75)) > 0.00000000001 {
		t.Errorf("Expected -2.75, got %f", d)
	}
}
