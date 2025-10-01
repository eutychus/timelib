package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestDiffDaysAYear(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-01-01 GMT")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	t2, _ := timelib.Strtotime("2021-12-31 GMT")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 364 {
		t.Errorf("Expected 364 days, got %d", days)
	}
}

func TestDiffDaysPHP81458_1(t *testing.T) {
	t1, _ := timelib.Strtotime("2018-07-01 04:00 GMT+0000")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	t1.UpdateTS(nil)

	t2, _ := timelib.Strtotime("2018-07-02 00:00")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode int
	tzi2, _ := timelib.ParseTzfile("America/Toronto", timelib.BuiltinDB(), &errCode)
	t2.UpdateTS(tzi2)

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 1 {
		t.Errorf("Expected 1 day, got %d", days)
	}
}

func TestDiffDaysPHP81458_2(t *testing.T) {
	t1, _ := timelib.Strtotime("2018-12-01 00:00")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode1 int
	tzi1, _ := timelib.ParseTzfile("UTC", timelib.BuiltinDB(), &errCode1)
	t1.UpdateTS(tzi1)

	t2, _ := timelib.Strtotime("2018-12-02 00:01")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode2 int
	tzi2, _ := timelib.ParseTzfile("UTC", timelib.BuiltinDB(), &errCode2)
	t2.UpdateTS(tzi2)

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 1 {
		t.Errorf("Expected 1 day, got %d", days)
	}
}

func TestDiffDaysPHP78452(t *testing.T) {
	t1, _ := timelib.Strtotime("2019-09-24 11:47:24")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode1 int
	tzi1, _ := timelib.ParseTzfile("Asia/Tehran", timelib.BuiltinDB(), &errCode1)
	t1.UpdateTS(tzi1)

	t2, _ := timelib.Strtotime("2019-08-21 12:47:24")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode2 int
	tzi2, _ := timelib.ParseTzfile("Asia/Tehran", timelib.BuiltinDB(), &errCode2)
	t2.UpdateTS(tzi2)

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 33 {
		t.Errorf("Expected 33 days, got %d", days)
	}
}

func TestDiffDaysPHP74524(t *testing.T) {
	t1, _ := timelib.Strtotime("2017-11-17 22:05:26.000000")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode1 int
	tzi1, _ := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &errCode1)
	t1.UpdateTS(tzi1)

	t2, _ := timelib.Strtotime("2017-04-03 22:29:15.079459")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode2 int
	tzi2, _ := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &errCode2)
	t2.UpdateTS(tzi2)

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 227 {
		t.Errorf("Expected 227 days, got %d", days)
	}
}

func TestDiffDaysDateTimeFallType2Type2(t *testing.T) {
	t1, _ := timelib.Strtotime("2010-11-07 00:15:35 EDT")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	t2, _ := timelib.Strtotime("2010-11-07 00:10:20 EDT")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 0 {
		t.Errorf("Expected 0 days, got %d", days)
	}
}

func TestDiffDaysDateTimeFallType3Type3(t *testing.T) {
	t1, _ := timelib.Strtotime("2010-11-07 00:15:35")
	if t1 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode1 int
	tzi1, _ := timelib.ParseTzfile("America/New_York", timelib.BuiltinDB(), &errCode1)
	t1.UpdateTS(tzi1)

	t2, _ := timelib.Strtotime("2010-11-07 00:10:20")
	if t2 == nil {
		t.Skip("Strtotime failed to parse date")
	}
	var errCode2 int
	tzi2, _ := timelib.ParseTzfile("America/New_York", timelib.BuiltinDB(), &errCode2)
	t2.UpdateTS(tzi2)

	diff := t1.Diff(t2)
	days := diff.Days

	if days != 0 {
		t.Errorf("Expected 0 days, got %d", days)
	}
}
