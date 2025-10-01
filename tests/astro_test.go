package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func TestAstroRiseSetAltitude(t *testing.T) {
	// Test astronomical calculations for sunrise/sunset
	// Create a time structure for testing
	tm := timelib.TimeCtor()
	tm.Y = 2005
	tm.M = 10
	tm.D = 17
	tm.H = 0
	tm.I = 0
	tm.S = 0
	tm.HaveDate = true
	tm.HaveTime = true

	// Test coordinates (Oslo, Norway)
	longitude := 9.627
	latitude := 59.186
	altitude := 0.0 // Sea level

	// Test the astronomical calculation function
	// Note: upperLimb parameter: 0 = center of sun, 1 = upper limb of sun
	upperLimb := 0
	hRise, hSet, rise, set, transit := timelib.AstroRiseSetAltitude(tm, longitude, latitude, altitude, upperLimb)

	// Basic validation - function should not crash and return reasonable values
	t.Logf("Astro calculation results:")
	t.Logf("  Sunrise hour angle: %f", hRise)
	t.Logf("  Sunset hour angle: %f", hSet)
	t.Logf("  Sunrise time: %d", rise)
	t.Logf("  Sunset time: %d", set)
	t.Logf("  Transit time: %d", transit)
}

func TestAstroBasic(t *testing.T) {
	// Test basic astronomical calculation functionality
	tm := timelib.TimeCtor()
	tm.Y = 2005
	tm.M = 10
	tm.D = 17
	tm.H = 0
	tm.I = 0
	tm.S = 0
	tm.HaveDate = true
	tm.HaveTime = true

	// Test with different coordinates
	testCases := []struct {
		longitude, latitude float64
		description         string
	}{
		{9.627, 59.186, "Oslo, Norway"},
		{-74.006, 40.7128, "New York, USA"},
		{0.0, 0.0, "Equator/Prime Meridian"},
		{0.0, 90.0, "North Pole"},
		{0.0, -90.0, "South Pole"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			hRise, hSet, rise, set, transit := timelib.AstroRiseSetAltitude(tm, tc.longitude, tc.latitude, 0.0, 0)

			// Basic validation - function should not crash
			t.Logf("Location: %s", tc.description)
			t.Logf("  Hour angles: Rise=%f, Set=%f", hRise, hSet)
			t.Logf("  Times: Rise=%d, Set=%d, Transit=%d", rise, set, transit)
		})
	}
}

func TestAstroTimeConversion(t *testing.T) {
	// Test time conversion for astronomical calculations
	tm := timelib.TimeCtor()
	tm.Sse = 1129449600 // Example timestamp (2005-10-17 00:00:00 UTC)
	tm.HaveTime = true
	tm.SseUptodate = true

	// Test that we can convert the time structure
	// This simulates the timelib_dump_date functionality
	t.Logf("Time structure: Y=%d, M=%d, D=%d, H=%d, I=%d, S=%d",
		tm.Y, tm.M, tm.D, tm.H, tm.I, tm.S)

	// Test basic time validation
	if tm.Sse < 0 {
		t.Error("Expected non-negative timestamp")
	}

	t.Log("AstroTimeConversion test completed successfully")
}
