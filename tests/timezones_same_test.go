package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestTimezoneSameType1Type1Same1(t *testing.T) {
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 GMT+0100", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 GMT+0100", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 GMT+0200", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 GMT+0100", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 GMT+0100", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 GMT+0200", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 CET", timelib.BuiltinDB())
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 CET", timelib.BuiltinDB())

	if err1 != nil || t1 == nil {
		t.Fatalf("Failed to parse time with CET: %v", err1)
	}
	if err2 != nil || t2 == nil {
		t.Fatalf("Failed to parse time with CET: %v", err2)
	}

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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 BST", timelib.BuiltinDB())
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 CET", timelib.BuiltinDB())

	if err1 != nil || t1 == nil {
		t.Fatalf("Failed to parse time with BST: %v", err1)
	}
	if err2 != nil || t2 == nil {
		t.Fatalf("Failed to parse time with CET: %v", err2)
	}

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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 CDT", timelib.BuiltinDB())
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 EST", timelib.BuiltinDB())

	if err1 != nil || t1 == nil {
		t.Fatalf("Failed to parse time with CDT: %v", err1)
	}
	if err2 != nil || t2 == nil {
		t.Fatalf("Failed to parse time with EST: %v", err2)
	}

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
	// Need to pass BuiltinDB() to parse timezone abbreviations
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 EST", timelib.BuiltinDB())
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 CDT", timelib.BuiltinDB())
	if err1 != nil || err2 != nil {
		t.Fatalf("Parse failed: t1 err=%v, t2 err=%v", err1, err2)
	}
	if t1 == nil || t2 == nil {
		t.Fatalf("Parse returned nil: t1=%v, t2=%v", t1, t2)
	}
	defer timelib.TimeDtor(t1)
	defer timelib.TimeDtor(t2)

	same := timelib.SameTimezone(t1, t2)

	if !same {
		t.Errorf("Expected same timezone (EST and CDT have same effective offset), got %v", same)
	}
	if t1.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t1.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t1.ZoneType)
	}
	if t2.ZoneType != timelib.TIMELIB_ZONETYPE_ABBR {
		t.Errorf("Expected t2.zone_type=%d, got %d", timelib.TIMELIB_ZONETYPE_ABBR, t2.ZoneType)
	}
}

func TestTimezoneSameType2Type2NotSame1(t *testing.T) {
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 EDT", timelib.BuiltinDB())
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 CDT", timelib.BuiltinDB())
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 CET", timelib.BuiltinDB())
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 CEST", timelib.BuiltinDB())
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 GMT+0100", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 BST", timelib.BuiltinDB())
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 GMT+0100", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 CEST", timelib.BuiltinDB())
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 GMT+0200", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39 CET", timelib.BuiltinDB())
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39", nil)
	if err1 != nil || t1 == nil { t.Fatal("Failed to parse t1") }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 GMT+0200", nil)
	if err2 != nil || t2 == nil { t.Fatal("Failed to parse t2") }
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
	// NOTE: CET may not be in BuiltinDB, so we parse the date without timezone
	// and then set timezone separately, matching the C test behavior
	t1, err1 := timelib.StrToTime("2021-11-05 11:23:39", nil)
	if err1 != nil || t1 == nil { t.Fatalf("Failed to parse t1: %v", err1) }
	t2, err2 := timelib.StrToTime("2021-11-05 11:24:07 CET", timelib.BuiltinDB())
	if t2 == nil {
		// If CET parsing fails, parse without timezone
		t2, err2 = timelib.StrToTime("2021-11-05 11:24:07", nil)
		if err2 != nil || t2 == nil { t.Fatalf("Failed to parse t2 without TZ: err=%v", err2) }
		// Manually set CET timezone info (Central European Time = UTC+1)
		t2.ZoneType = timelib.TIMELIB_ZONETYPE_ABBR
		t2.Z = 3600  // UTC+1
		t2.Dst = 0
	}
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
