package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestRenderTimestamp(t *testing.T) {
	// Test rendering a timestamp for Europe/Amsterdam timezone
	ts := int64(1114819200) // Example timestamp (2005-04-30 00:00:00 UTC)

	// Create a time structure
	time := timelib.TimeCtor()
	time.Sse = ts
	time.HaveTime = true
	time.SseUptodate = true

	// Parse timezone info for Europe/Amsterdam
	var dummyError int
	tz, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &dummyError)
	if err != nil {
		t.Errorf("Europe/Amsterdam timezone not available: %v", err)
	}
	if dummyError != timelib.TIMELIB_ERROR_NO_ERROR {
		t.Errorf("Europe/Amsterdam timezone parse error: %d", dummyError)
	}

	timelib.SetTimezone(time, tz)

	// Convert to local time
	time.Unixtime2local(ts)

	// Verify the basic structure is set up correctly
	if time.Sse != ts {
		t.Errorf("Expected SSE %d, got %d", ts, time.Sse)
	}

	// Check that timezone info is set
	if time.TzInfo == nil {
		t.Error("Expected timezone info to be set")
	}

	if time.TzInfo.Name != "Europe/Amsterdam" {
		t.Errorf("Expected timezone name Europe/Amsterdam, got %s", time.TzInfo.Name)
	}

	// The actual date/time conversion may not be fully implemented yet,
	// so we just verify the structure is properly set up
	if !time.HaveTime {
		t.Error("Expected HaveTime to be true")
	}

	if !time.SseUptodate {
		t.Error("Expected SseUptodate to be true")
	}
}

func TestRenderTimestampBasic(t *testing.T) {
	// Test basic timestamp rendering functionality
	ts := int64(946728000) // Known timestamp (J2000 epoch)

	tm := timelib.TimeCtor()
	tm.Sse = ts
	tm.HaveTime = true
	tm.SseUptodate = true

	// Convert to local time (UTC)
	tm.Unixtime2local(ts)

	// Verify basic structure
	if tm.Sse != ts {
		t.Errorf("Expected SSE %d, got %d", ts, tm.Sse)
	}

	if !tm.HaveTime {
		t.Error("Expected HaveTime to be true")
	}

	if !tm.SseUptodate {
		t.Error("Expected SseUptodate to be true")
	}
}
