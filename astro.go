package timelib

import (
	"math"
)

// Julian date constants
const (
	// JulianDayEpochOffset is the Julian Day number for Unix epoch (1970-01-01 00:00:00 UTC)
	JulianDayEpochOffset = 2440587.5
	// J2000Epoch is the Julian Day number for J2000.0 epoch (2000-01-01 12:00:00 TT)
	J2000Epoch = 2451545.0
	// SecondsPerDay is the number of seconds in a day
	SecondsPerDay = 86400.0
)

// TsToJulianDay converts a Unix timestamp to Julian Day number
// Julian Day is the continuous count of days since the beginning of
// the Julian Period (January 1, 4713 BC)
func TsToJulianDay(ts int64) float64 {
	tmp := float64(ts)
	tmp /= SecondsPerDay
	tmp += JulianDayEpochOffset
	return tmp
}

// Astronomical calculation constants
const (
	PI     = 3.1415926535897932384
	RADEG  = 180.0 / PI
	DEGRAD = PI / 180.0
	INV360 = 1.0 / 360.0
)

// Helper functions for degree-based trigonometry
func sind(x float64) float64 { return math.Sin(x * DEGRAD) }
func cosd(x float64) float64 { return math.Cos(x * DEGRAD) }

// func tand(x float64) float64  { return math.Tan(x * DEGRAD) }  // Unused - commented out
// func atand(x float64) float64 { return RADEG * math.Atan(x) }  // Unused - commented out
// func asind(x float64) float64 { return RADEG * math.Asin(x) }  // Unused - commented out
func acosd(x float64) float64     { return RADEG * math.Acos(x) }
func atan2d(y, x float64) float64 { return RADEG * math.Atan2(y, x) }

// astroRevolution reduces any angle to within 0..360 degrees
func astroRevolution(x float64) float64 {
	return x - 360.0*math.Floor(x*INV360)
}

// astroRev180 reduces angle to within -180..+180 degrees
func astroRev180(x float64) float64 {
	return x - 360.0*math.Floor(x*INV360+0.5)
}

// gmst0 computes the Greenwich Mean Sidereal Time at 0h UT
func gmst0(d float64) float64 {
	// Sidtime at 0h UT = L (Sun's mean longitude) + 180 degrees
	// L = M + w, as defined in sunpos(). Since I'm too lazy to
	// add these numbers, I'll let the C compiler do it for me.
	// 0.9856002585 is the number of degrees traversed by the Sun
	// each day. Thus we get GMST at 0h UT by computing
	// GMST = M + w + 180 at 0h UT
	sidtim0 := astroRevolution((180.0 + 356.0470 + 282.9404) + (0.9856002585+4.70935e-5)*d)
	return sidtim0
}

// sunpos computes the Sun's ecliptic longitude and distance
// at an instant given in d, number of days since 2000 Jan 0.0
// (Note: 2000 Jan 0.0 = 1999 Dec 31.0)
func sunpos(d float64) (lon, r float64) {
	// Compute Mean longitude of Sun
	M := astroRevolution(356.0470 + 0.9856002585*d)

	// Compute mean anomaly of the Sun
	w := 282.9404 + 4.70935e-5*d

	// Compute the Sun's ecliptic longitude and distance
	e := 0.016709 - 1.151e-9*d

	// Eccentric anomaly
	E := M + e*RADEG*sind(M)*(1.0+e*cosd(M))

	// Compute the Sun's distance r and true anomaly v
	x := cosd(E) - e
	y := math.Sqrt(1.0-e*e) * sind(E)

	r = math.Sqrt(x*x + y*y)
	v := atan2d(y, x)

	// Compute the Sun's true longitude
	lon = v + w

	// Normalize the longitude
	lon = astroRevolution(lon)

	return lon, r
}

// sunRaDec computes the Sun's equatorial coordinates RA and Decl
// and also its distance, at an instant given in d,
// the number of days since 2000 Jan 0.0
func sunRaDec(d float64) (ra, dec, r float64) {
	// Compute Sun's ecliptical coordinates
	lon, r := sunpos(d)

	// Compute ecliptic rectangular coordinates (z=0)
	x := r * cosd(lon)
	y := r * sind(lon)

	// Compute obliquity of ecliptic (inclination of Earth's axis)
	obl_ecl := 23.4393 - 3.563e-7*d

	// Convert to equatorial rectangular coordinates - x is unchanged
	z := y * sind(obl_ecl)
	y = y * cosd(obl_ecl)

	// Convert to spherical coordinates
	ra = atan2d(y, x)
	dec = atan2d(z, math.Sqrt(x*x+y*y))

	return ra, dec, r
}

// daysSince2000Jan0 computes the number of days since 2000 Jan 0.0
// func daysSince2000Jan0(y, m, d int64) int64 {  // Unused - commented out
// 	return 367*y - (7*(y+((m+9)/12)))/4 + (275*m)/9 + d - 730530
// }

// AstroRiseSetAltitude computes rise, set and transit times for celestial objects
//
// Parameters:
//
//	t: Time structure with the date to compute for
//	lon: Longitude in degrees (East positive, West negative)
//	lat: Latitude in degrees (North positive, South negative)
//	altit: Altitude of object in degrees (negative for below horizon)
//	       Use -35.0/60.0 for Sun's upper limb touching horizon
//	upper_limb: 1 to compute rise/set for upper limb, 0 for center
//
// Returns:
//
//	hRise: Hour of rise (decimal hours, local time)
//	hSet: Hour of set (decimal hours, local time)
//	tsRise: Unix timestamp of rise
//	tsSet: Unix timestamp of set
//	tsTransit: Unix timestamp of transit (highest point)
func AstroRiseSetAltitude(t *Time, lon, lat, altit float64, upperLimb int) (hRise, hSet float64, tsRise, tsSet, tsTransit int64) {
	// Create UTC time at 00:00 of the given day
	tUtc := &Time{
		Y: t.Y, M: t.M, D: t.D,
		H: 0, I: 0, S: 0,
	}
	tUtc.UpdateTS(nil)

	// Compute d of 12h local mean solar time
	timestamp := tUtc.Sse
	d := TsToJ2000(timestamp) + 2 - lon/360.0

	// Compute local sidereal time of this moment
	sidtime := astroRevolution(gmst0(d) + 180.0 + lon)

	// Compute Sun's RA + Decl at this moment
	sRA, sdec, sr := sunRaDec(d)

	// Compute time when Sun is at south - in hours UT
	tsouth := 12.0 - astroRev180(sidtime-sRA)/15.0

	// Compute the Sun's apparent radius, degrees
	sradius := 0.2666 / sr

	// Do correction to upper limb, if necessary
	if upperLimb != 0 {
		altit -= sradius
	}

	// Compute the diurnal arc that the Sun traverses to reach the specified altitude
	cost := (sind(altit) - sind(lat)*sind(sdec)) / (cosd(lat) * cosd(sdec))
	tsTransit = tUtc.Sse + int64(tsouth*3600)

	if cost >= 1.0 {
		// Sun always below altit
		tsRise = tUtc.Sse + int64(tsouth*3600)
		tsSet = tUtc.Sse + int64(tsouth*3600)
		hRise = tsouth
		hSet = tsouth
	} else if cost <= -1.0 {
		// Sun always above altit
		tsRise = t.Sse - 12*3600
		tsSet = t.Sse + 12*3600
		hRise = tsouth - 12.0
		hSet = tsouth + 12.0
	} else {
		// The diurnal arc, hours
		ta := acosd(cost) / 15.0

		// Store rise and set times - as Unix Timestamp
		tsRise = tUtc.Sse + int64((tsouth-ta)*3600)
		tsSet = tUtc.Sse + int64((tsouth+ta)*3600)

		hRise = tsouth - ta
		hSet = tsouth + ta
	}

	return hRise, hSet, tsRise, tsSet, tsTransit
}
