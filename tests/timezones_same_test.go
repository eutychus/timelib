package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestTimezoneSameType1Type1Same1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 GMT+0100")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 GMT+0100")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t2.ZoneType)
	}
}

func TestTimezoneSameType1Type1NotSame1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 GMT+0200")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 GMT+0100")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t2.ZoneType)
	}
}

func TestTimezoneSameType1Type1NotSame2(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 GMT+0100")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 GMT+0200")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2Same1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 CET")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 CET")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2Same2(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 BST")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 CET")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2Same3(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 CDT")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 EST")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2Same4(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 EST")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 CDT")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2NotSame1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 EDT")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 CDT")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2NotSame2(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 CET")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 CEST")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType3Type3Same1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	var dummyError int
	tzi1, _ := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &dummyError)
	tzi2, _ := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &dummyError)

	t1.UpdateTS(tzi1)
	t2.UpdateTS(tzi2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t2.ZoneType)
	}
}

func TestTimezoneSameType3Type3NotSame1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	var dummyError int
	tzi1, _ := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &dummyError)
	tzi2, _ := timelib.ParseTzfile("America/Chicago", timelib.BuiltinDB(), &dummyError)

	t1.UpdateTS(tzi1)
	t2.UpdateTS(tzi2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t2.ZoneType)
	}
}

func TestTimezoneSameType1Type2(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 GMT+0100")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 BST")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType1Type3(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 GMT+0100")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	var dummyError int
	tzi2, _ := timelib.ParseTzfile("Europe/Berlin", timelib.BuiltinDB(), &dummyError)
	t2.UpdateTS(tzi2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 CEST")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 GMT+0200")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type3(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39 CET")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	var dummyError int
	tzi2, _ := timelib.ParseTzfile("Europe/Berlin", timelib.BuiltinDB(), &dummyError)
	t2.UpdateTS(tzi2)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t2.ZoneType)
	}
}

func TestTimezoneSameType3Type1(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 GMT+0200")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	var dummyError int
	tzi1, _ := timelib.ParseTzfile("Europe/Berlin", timelib.BuiltinDB(), &dummyError)
	t1.UpdateTS(tzi1)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_OFFSET {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_OFFSET, t2.ZoneType)
	}
}

func TestTimezoneSameType3Type2(t *testing.T) {
	t1, _ := timelib.Strtotime("2021-11-05 11:23:39")
	t2, _ := timelib.Strtotime("2021-11-05 11:24:07 CET")
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	var dummyError int
	tzi1, _ := timelib.ParseTzfile("Europe/Berlin", timelib.BuiltinDB(), &dummyError)
	t1.UpdateTS(tzi1)

	same := timelib.SameTimezone(t1, t2)

	if same {
		t.Errorf("Expected different timezone comparison result, got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ID {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ID, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}
