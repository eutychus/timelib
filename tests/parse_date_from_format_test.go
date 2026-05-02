package tests

import (
	"testing"
	timelib "github.com/eutychus/timelib"
)

func TestParseDateFromFormat_Basic(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		withPrefix bool
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		z          int32
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		// naturalDateWithoutPrefix
		{"naturalDateWithoutPrefix", "Y/m/d", "2018/01/26", false, 2018, 1, 26, unset, unset, unset, 0, 0, "", 0, 0},
		// isoDateWithoutPrefix
		{"isoDateWithoutPrefix_1", "V.b.B", "53.7.2017", false, 2018, 1, 7, unset, unset, unset, 0, 0, "", 0, 0},
		{"isoDateWithoutPrefix_2", "V/B", "53/2017", false, 2018, 1, 1, unset, unset, unset, 0, 0, "", 0, 0},
		{"isoDateWithoutPrefix_3", "B", "2017", false, 2017, 1, 2, unset, unset, unset, 0, 0, "", 0, 0},
		// invalidISOWeek
		{"invalidISOWeek_1", "V/B", "55/2017", false, unset, unset, unset, unset, unset, unset, 0, 0, "", 1, 0},
		{"invalidISOWeek_2", "V", "52", false, unset, unset, unset, unset, unset, unset, 0, 0, "", 0, 1},
		// invalidISODayOfWeek
		{"invalidISODayOfWeek", "b", "8", false, unset, unset, unset, unset, unset, unset, 0, 0, "", 1, 0},
		// tzOffsetMinutes
		{"tzOffsetMinutes", "Y/m/d Z", "2018/01/26 +285", false, 2018, 1, 26, unset, unset, unset, 0, 285 * 60, "", 0, 0},
		// tzOffsetHours
		{"tzOffsetHours_1", "Y/m/d P", "2018/01/26 +02:00", false, 2018, 1, 26, unset, unset, unset, 0, 2 * 60 * 60, "", 0, 0},
		{"tzOffsetHours_2", "Y/m/d p", "2018/01/26 +02:00", false, 2018, 1, 26, unset, unset, unset, 0, 2 * 60 * 60, "", 0, 0},
		// tzUtc
		{"tzUtc", "Y/m/d H:i:sp", "2018/01/26 13:20:00Z", false, 2018, 1, 26, 13, 20, 0, 0, 0, "", 0, 0},
		// cannotMixISOWithNatural
		{"cannotMixISOWithNatural", "B/m/d", "2018/01/26", false, 2018, 1, 26, unset, unset, unset, 0, 0, "", 1, 0},
		// cannotHaveMeridianBeforeHour
		{"cannotHaveMeridianBeforeHour", "d M Y A h:i", "11 Mar 2013 PM 3:34", false, 2013, 3, 11, 3, 34, 0, 0, 0, "", 1, 0},
		// cannotHaveMeridianWithoutHour
		{"cannotHaveMeridianWithoutHour", "d M Y A", "11 Mar 2013 PM", false, 2013, 3, 11, 0, 0, 0, 0, 0, "", 1, 0},
		// cannotHaveDOYBeforeYear
		{"cannotHaveDOYBeforeYear", "z Y", "60 2020", false, 2020, unset, unset, unset, unset, unset, 0, 0, "", 1, 0},
		// DOYAfterLeapYear
		{"DOYAfterLeapYear", "Y z", "2020 60", false, 2020, 3, 1, unset, unset, unset, 0, 0, "", 0, 0},
		// DOYAfterYear
		{"DOYAfterYear", "Y z", "2021 60", false, 2021, 3, 2, unset, unset, unset, 0, 0, "", 0, 0},
		// naturalDateWithPrefix
		{"naturalDateWithPrefix", "Year %Y Month %m Day %d", "Year 2018 Month 01 Day 26", true, 2018, 1, 26, unset, unset, unset, 0, 0, "", 0, 0},
		// naturalDateWithTime
		{"naturalDateWithTime", "%Y-%m-%dT%H:%i:%sZ", "2018-01-26T11:56:02Z", true, 2018, 1, 26, 11, 56, 2, 0, 0, "", 0, 0},
		// isoDateWithPrefix
		{"isoDateWithPrefix_1", "%V.%b.%B", "53.7.2017", true, 2018, 1, 7, unset, unset, unset, 0, 0, "", 0, 0},
		{"isoDateWithPrefix_2", "%V/%B", "53/2017", true, 2018, 1, 1, unset, unset, unset, 0, 0, "", 0, 0},
		{"isoDateWithPrefix_3", "%B", "2017", true, 2017, 1, 2, unset, unset, unset, 0, 0, "", 0, 0},
		// missingPrefix
		{"missingPrefix", "Y/m/d", "2018/01/26", true, unset, unset, unset, unset, unset, unset, 0, 0, "", 5, 0},
		// prefixEscape
		{"prefixEscape_1", "%%", "%", true, unset, unset, unset, unset, unset, unset, 0, 0, "", 0, 0},
		{"prefixEscape_2", "%%Y", "%Y", true, unset, unset, unset, unset, unset, unset, 0, 0, "", 0, 0},
		// unmatchedPrefix
		{"unmatchedPrefix", "%Y%m", "2018", true, 2018, unset, unset, unset, unset, unset, 0, 0, "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result *timelib.Time
			var errContainer *timelib.ErrorContainer
			if tt.withPrefix {
				result, errContainer = timelib.ParseFromFormatWithPrefix(tt.format, tt.input)
			} else {
				result, errContainer = timelib.ParseFromFormat(tt.format, tt.input)
			}
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if result.TzAbbr != tt.tzAbbr {
				t.Errorf("TzAbbr: got %s, want %s", result.TzAbbr, tt.tzAbbr)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_American(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"american00", "m/d", "9/11", unset, 9, 11, unset, unset, unset, 0, 0},
		{"american01", "m/d", "09/11", unset, 9, 11, unset, unset, unset, 0, 0},
		{"american02", "m/d/y", "12/22/69", 2069, 12, 22, unset, unset, unset, 0, 0},
		{"american03", "m/d/y", "12/22/70", 1970, 12, 22, unset, unset, unset, 0, 0},
		{"american04", "m/d/y", "12/22/78", 1978, 12, 22, unset, unset, unset, 0, 0},
		{"american05", "m/d/Y", "12/22/1978", 1978, 12, 22, unset, unset, unset, 0, 0},
		{"american06", "m/d/Y", "12/22/2078", 2078, 12, 22, unset, unset, unset, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug41523(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"bug41523_00", "Y-m-d", "0000-00-00", 0, 0, 0, 0, 0},
		{"bug41523_01", "Y-m-d", "0001-00-00", 1, 0, 0, 0, 0},
		{"bug41523_02", "Y-m-d", "0002-00-00", 2, 0, 0, 0, 0},
		{"bug41523_03", "Y-m-d", "0003-00-00", 3, 0, 0, 0, 0},
		{"bug41523_04", "y-m-d", "00-00-00", 2000, 0, 0, 0, 0},
		{"bug41523_05", "y-m-d", "01-00-00", 2001, 0, 0, 0, 0},
		{"bug41523_06", "y-m-d", "02-00-00", 2002, 0, 0, 0, 0},
		{"bug41523_07", "y-m-d", "03-00-00", 2003, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug41842(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"bug41842_00", "-Y-m-d", "-0001-06-28", 1, 6, 28, 0, 0},
		{"bug41842_01", "-Y-m-d", "-2007-06-28", 2007, 6, 28, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug55240(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"bug55240_00", "d.m.Y Gi", "11.11.2009 800", 2009, 11, 11, 80, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_DateFull(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"datefull_00", "d F Y", "22 dec 1978", 1978, 12, 22, 0, 0},
		{"datefull_01", "d-F-y", "22-dec-78", 1978, 12, 22, 0, 0},
		{"datefull_02", "d F Y", "22 Dec 1978", 1978, 12, 22, 0, 0},
		{"datefull_03", "dFy  ", "22DEC78", 1978, 12, 22, 0, 0},
		{"datefull_04", "d F Y", "22 december 1978", 1978, 12, 22, 0, 0},
		{"datefull_05", "d-F-y", "22-december-78", 1978, 12, 22, 0, 0},
		{"datefull_06", "d F Y", "22 December 1978", 1978, 12, 22, 0, 0},
		{"datefull_07", "dFy  ", "22DECEMBER78", 1978, 12, 22, 0, 0},
		{"datefull_08", "d?F?Y", "22 dec\t1978", 1978, 12, 22, 0, 0},
		{"datefull_09", "d?F?Y", "22\tDec\t1978", 1978, 12, 22, 0, 0},
		{"datefull_10", "d?F?Y", "22\tdecember\t1978", 1978, 12, 22, 0, 0},
		{"datefull_11", "d?F?Y", "22\tDecember\t1978", 1978, 12, 22, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_DateNoColon(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"datenocolon_00", "Ymd", "19781222", 1978, 12, 22, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Date(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"date_00", "d.m.Y", "31.01.2006", 2006, 1, 31, 0, 0},
		{"date_01", "d.m.Y", "32.01.2006", 2006, 1, 32, 0, 0},
		{"date_02", "d.m.Y", "28.01.2006", 2006, 1, 28, 0, 0},
		{"date_03", "d.m.Y", "29.01.2006", 2006, 1, 29, 0, 0},
		{"date_04", "d.m.Y", "30.01.2006", 2006, 1, 30, 0, 0},
		{"date_05", "d.m.Y", "31.01.2006", 2006, 1, 31, 0, 0},
		{"date_06", "d.m.Y", "32.01.2006", 2006, 1, 32, 0, 0},
		{"date_07", "d-m-Y", "31-01-2006", 2006, 1, 31, 0, 0},
		{"date_08", "d-m-Y", "32-01-2006", 2006, 1, 32, 0, 0},
		{"date_09", "d-m-Y", "28-01-2006", 2006, 1, 28, 0, 0},
		{"date_10", "d-m-Y", "29-01-2006", 2006, 1, 29, 0, 0},
		{"date_11", "d-m-Y", "30-01-2006", 2006, 1, 30, 0, 0},
		{"date_12", "d-m-Y", "31-01-2006", 2006, 1, 31, 0, 0},
		{"date_13", "d-m-Y", "32-01-2006", 2006, 1, 32, 0, 0},
		{"date_14", "d-m-Y", "29-02-2006", 2006, 2, 29, 0, 0},
		{"date_15", "d-m-Y", "30-02-2006", 2006, 2, 30, 0, 0},
		{"date_16", "d-m-Y", "31-02-2006", 2006, 2, 31, 0, 0},
		{"date_17", "d-m-Y", "01-01-2006", 2006, 1, 1, 0, 0},
		{"date_18", "d-m-Y", "31-12-2006", 2006, 12, 31, 0, 0},
		{"date_19", "d.m.Y", "31.13.2006", 2006, 13, 31, 0, 0},
		{"date_20", "m/d/Y", "11/10/2006", 2006, 11, 10, 0, 0},
		{"date_21", "m/d/Y", "12/10/2006", 2006, 12, 10, 0, 0},
		{"date_22", "m/d/Y", "13/10/2006", 2006, 13, 10, 0, 0},
		{"date_23", "m/d/Y", "14/10/2006", 2006, 14, 10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_DateRoman(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"dateroman_00", "d F Y  ", "22 I 1978", 1978, 1, 22, 0, 0},
		{"dateroman_01", "d. F Y ", "22. II 1978", 1978, 2, 22, 0, 0},
		{"dateroman_02", "d F. Y ", "22 III. 1978", 1978, 3, 22, 0, 0},
		{"dateroman_03", "d- F- Y", "22- IV- 1978", 1978, 4, 22, 0, 0},
		{"dateroman_04", "d -F -Y", "22 -V -1978", 1978, 5, 22, 0, 0},
		{"dateroman_05", "d-F-Y  ", "22-VI-1978", 1978, 6, 22, 0, 0},
		{"dateroman_06", "d.F.Y  ", "22.VII.1978", 1978, 7, 22, 0, 0},
		{"dateroman_07", "d F Y  ", "22 VIII 1978", 1978, 8, 22, 0, 0},
		{"dateroman_08", "d F Y  ", "22 IX 1978", 1978, 9, 22, 0, 0},
		{"dateroman_09", "d F Y  ", "22 X 1978", 1978, 10, 22, 0, 0},
		{"dateroman_10", "d F Y  ", "22 XI 1978", 1978, 11, 22, 0, 0},
		{"dateroman_11", "d?F?Y  ", "22\tXII\t1978", 1978, 12, 22, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_DateSlash(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"dateslash_00", "Y/n/d", "2005/8/12", 2005, 8, 12, 0, 0},
		{"dateslash_01", "Y/m/d", "2005/01/02", 2005, 1, 2, 0, 0},
		{"dateslash_02", "Y/m/j", "2005/01/2", 2005, 1, 2, 0, 0},
		{"dateslash_03", "Y/n/d", "2005/1/02", 2005, 1, 2, 0, 0},
		{"dateslash_04", "Y/n/j", "2005/1/2", 2005, 1, 2, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Lenient(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"lenient_00", "Y-m-d", "2001-11-29 13:20:01", 2001, 11, 29, unset, unset, unset, 0, 0},
		{"lenient_01", "Y-m-d+", "2001-11-29 13:20:01", 2001, 11, 29, unset, unset, unset, 0, 0},
		{"lenient_02", "Y-m-d +", "2001-11-29 13:20:01", 2001, 11, 29, unset, unset, unset, 0, 0},
		{"lenient_03", "Y-m-d+", "2001-11-29", 2001, 11, 29, unset, unset, unset, 0, 0},
		{"lenient_04", "Y-m-d +", "2001-11-29", 2001, 11, 29, unset, unset, unset, 0, 0},
		{"lenient_05", "Y-m-d +", "2001-11-29", 2001, 11, 29, unset, unset, unset, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Mysql(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"mysql_00", "YmdHis", "19970523091528", 1997, 5, 23, 9, 15, 28, 0, 0},
		{"mysql_01", "YmdHis", "20001231185859", 2000, 12, 31, 18, 58, 59, 0, 0},
		{"mysql_02", "YmdHis", "20500410101010", 2050, 4, 10, 10, 10, 10, 0, 0},
		{"mysql_03", "YmdHis", "20050620091407", 2005, 6, 20, 9, 14, 7, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Pgsql(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"pgsql_00", "F d, Y", "January 8, 1999", 1999, 1, 8, 0, 0},
		{"pgsql_01", "F?d,?Y", "January\t8,\t1999", 1999, 1, 8, 0, 0},
		{"pgsql_02", "Y-m-d", "1999-01-08", 1999, 1, 8, 0, 0},
		{"pgsql_03", "m/d/Y", "1/8/1999", 1999, 1, 8, 0, 0},
		{"pgsql_04", "m/d/Y", "1/18/1999", 1999, 1, 18, 0, 0},
		{"pgsql_05", "m/d/y", "01/02/03", 2003, 1, 2, 0, 0},
		{"pgsql_06", "Y-M-d", "1999-Jan-08", 1999, 1, 8, 0, 0},
		{"pgsql_07", "M-d-Y", "Jan-08-1999", 1999, 1, 8, 0, 0},
		{"pgsql_08", "d-M-Y", "08-Jan-1999", 1999, 1, 8, 0, 0},
		{"pgsql_09", "y-M-d", "99-Jan-08", 1999, 1, 8, 0, 0},
		{"pgsql_10", "d-M-y", "08-Jan-99", 1999, 1, 8, 0, 0},
		{"pgsql_11", "M-d-y", "Jan-08-99", 1999, 1, 8, 0, 0},
		{"pgsql_12", "Ymd", "1999008", 1999, 1, 8, 0, 0},
		{"pgsql_13", "Y.z", "1999.008", 1999, 1, 9, 0, 0},
		{"pgsql_14", "Y.z", "1999.038", 1999, 2, 8, 0, 0},
		{"pgsql_15", "Y.z", "1999.238", 1999, 8, 27, 0, 0},
		{"pgsql_16", "Y.z", "1999.366", 2000, 1, 2, 0, 0},
		{"pgsql_17", "Yz", "1999008", 1999, 1, 9, 0, 0},
		{"pgsql_18", "Yz", "1999038", 1999, 2, 8, 0, 0},
		{"pgsql_19", "Yz", "1999238", 1999, 8, 27, 0, 0},
		{"pgsql_20", "Yz", "1999366", 2000, 1, 2, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_PointedDate(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		errorCount int
		warnCount  int
	}{
		{"pointeddate_00", "d.m.Y", "22.12.1978", 1978, 12, 22, 0, 0},
		{"pointeddate_01", "d.m.Y", "22.7.1978", 1978, 7, 22, 0, 0},
		{"pointeddate_02", "d.m.y", "22.12.78", 1978, 12, 22, 0, 0},
		{"pointeddate_03", "d.m.y", "22.7.78", 1978, 7, 22, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_TimeLong12(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"timelong12_00", "g:i:sa", "01:00:03am", unset, unset, unset, 1, 0, 3, 0, 0},
		{"timelong12_01", "g:i:sa", "01:03:12pm", unset, unset, unset, 13, 3, 12, 0, 0},
		{"timelong12_02", "g:i:s A", "12:31:13 A.M.", unset, unset, unset, 0, 31, 13, 0, 0},
		{"timelong12_03", "g:i:s A", "08:13:14 P.M.", unset, unset, unset, 20, 13, 14, 0, 0},
		{"timelong12_04", "g:i:s A", "11:59:15 AM", unset, unset, unset, 11, 59, 15, 0, 0},
		{"timelong12_05", "g:i:s A", "06:12:16 PM", unset, unset, unset, 18, 12, 16, 0, 0},
		{"timelong12_06", "g:i:s a", "07:08:17 am", unset, unset, unset, 7, 8, 17, 0, 0},
		{"timelong12_07", "g:i:s a", "08:09:18 p.m.", unset, unset, unset, 20, 9, 18, 0, 0},
		{"timelong12_08", "h:i:sa", "01.00.03am", unset, unset, unset, 1, 0, 3, 0, 0},
		{"timelong12_09", "h:i:sa", "01.03.12pm", unset, unset, unset, 13, 3, 12, 0, 0},
		{"timelong12_10", "h:i:s A", "12.31.13 A.M.", unset, unset, unset, 0, 31, 13, 0, 0},
		{"timelong12_11", "h:i:s A", "08.13.14 P.M.", unset, unset, unset, 20, 13, 14, 0, 0},
		{"timelong12_12", "h:i:s A", "11.59.15 AM", unset, unset, unset, 11, 59, 15, 0, 0},
		{"timelong12_13", "h:i:s A", "06.12.16 PM", unset, unset, unset, 18, 12, 16, 0, 0},
		{"timelong12_14", "h:i:s a", "07.08.17 am", unset, unset, unset, 7, 8, 17, 0, 0},
		{"timelong12_15", "h:i:s a", "08.09.18 p.m.", unset, unset, unset, 20, 9, 18, 0, 0},
		{"timelong12_16", "h:i:s a", "07.08.17\tam", unset, unset, unset, 7, 8, 17, 0, 0},
		{"timelong12_17", "h:i:s a", "08.09.18\tp.m.", unset, unset, unset, 20, 9, 18, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug44426(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"bug44426_00", "M d Y h:i:s:uA", "Aug 27 2007 12:00:00:000AM", 2007, 8, 27, 0, 0, 0, 0, 0},
		{"bug44426_01", "M d Y h:i:s.uA", "Aug 27 2007 12:00:00.000AM", 2007, 8, 27, 0, 0, 0, 0, 0},
		{"bug44426_02", "M d Y h:i:s:u", "Aug 27 2007 12:00:00:000", 2007, 8, 27, 12, 0, 0, 0, 0},
		{"bug44426_03", "M d Y h:i:s.u", "Aug 27 2007 12:00:00.000", 2007, 8, 27, 12, 0, 0, 0, 0},
		{"bug44426_04", "M d Y h:i:sA", "Aug 27 2007 12:00:00AM", 2007, 8, 27, 0, 0, 0, 0, 0},
		{"bug44426_05", "M d Y h:i:s:uA", "Aug 27 2007", 2007, 8, 27, unset, unset, unset, 0, 0},
		{"bug44426_06", "M d Y h:iA", "Aug 27 2007 12:00AM", 2007, 8, 27, 0, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug50392(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		errorCount int
		warnCount  int
	}{
		{"bug50392_00", "Y-m-d H:i:s", "2010-03-06 16:07:25", 2010, 3, 6, 16, 7, 25, 0, 0, 0},
		{"bug50392_01", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.", 2010, 3, 6, 16, 7, 25, 0, 0, 0},
		{"bug50392_02", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.1", 2010, 3, 6, 16, 7, 25, 100000, 0, 0},
		{"bug50392_03", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.12", 2010, 3, 6, 16, 7, 25, 120000, 0, 0},
		{"bug50392_04", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.123", 2010, 3, 6, 16, 7, 25, 123000, 0, 0},
		{"bug50392_05", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.1234", 2010, 3, 6, 16, 7, 25, 123400, 0, 0},
		{"bug50392_06", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.12345", 2010, 3, 6, 16, 7, 25, 123450, 0, 0},
		{"bug50392_07", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.123456", 2010, 3, 6, 16, 7, 25, 123456, 0, 0},
		{"bug50392_08", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.1234567", 2010, 3, 6, 16, 7, 25, 123456, 0, 0},
		{"bug50392_09", "Y-m-d H:i:s.u", "2010-03-06 16:07:25.12345678", 2010, 3, 6, 16, 7, 25, 123456, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug75577(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		errorCount int
		warnCount  int
	}{
		{"bug75577_00", "Y-m-d H:i:s.v", "2010-03-06 16:07:25.123", 2010, 3, 6, 16, 7, 25, 123000, 0, 0},
		{"bug75577_01", "Y-m-d H:i:s.v", "2010-03-06 16:07:25.1234", 2010, 3, 6, 16, 7, 25, 123000, 0, 0},
		{"bug75577_02", "Y-m-d H:i:s.v", "2010-03-06 16:07:25.12345", 2010, 3, 6, 16, 7, 25, 123000, 0, 0},
		{"bug75577_03", "Y-m-d H:i:s.v", "2010-03-06 16:07:25.123456", 2010, 3, 6, 16, 7, 25, 123000, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bugs(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		z          int32
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"bugs_00", "m/d/y Hi", "04/05/06 0045", 2006, 4, 5, 0, 45, 0, 0, 0, "", 0, 0},
		{"bugs_01", "H:i Y-m-d", "17:00 2004-01-03", 2004, 1, 3, 17, 0, 0, 0, 0, "", 0, 0},
		{"bugs_02", "Y-m-d H:i:s.ue", "2004-03-10 16:33:17.11403+01", 2004, 3, 10, 16, 33, 17, 114030, 3600, "", 0, 0},
		{"bugs_03", "Y-m-d H:i:se", "2004-03-10 16:33:17+01", 2004, 3, 10, 16, 33, 17, 0, 3600, "", 0, 0},
		{"bugs_04", "Y-m-d H:i:s e", "2003-11-19 08:00:00 T", 2003, 11, 19, 8, 0, 0, 0, -25200, "T", 0, 0},
		{"bugs_05", "d-M-Y H:i:s", "01-MAY-1982 00:00:00", 1982, 5, 1, 0, 0, 0, 0, 0, "", 0, 0},
		{"bugs_06", "Y-m-d?H:i:s", "2040-06-12T04:32:12", 2040, 6, 12, 4, 32, 12, 0, 0, "", 0, 0},
		{"bugs_07", "F dS", "july 14th", unset, 7, 14, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_08", "F d?e", "july 14tH", unset, 7, 14, unset, unset, unset, 0, 28800, "H", 0, 0},
		{"bugs_09", "dF", "11Oct", unset, 10, 11, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_10", "d F", "11 Oct", unset, 10, 11, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_11", "Fd, Y", "Jan14, 2004", 2004, 1, 14, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_12", "F d, Y", "Jan 14, 2004", 2004, 1, 14, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_13", "F.d, Y", "Jan.14, 2004", 2004, 1, 14, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_14", "Y-m-d", "1999-10-13", 1999, 10, 13, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_15", "F d  Y", "Oct 13  1999", 1999, 10, 13, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_16", "Y-m-d", "2000-01-19", 2000, 1, 19, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_17", "F d  Y", "Jan 19  2000", 2000, 1, 19, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_18", "Y-m-d", "2001-12-21", 2001, 12, 21, unset, unset, unset, 0, 0, "", 0, 0},
		{"bugs_19", "F d  H:is", "Dec 21  12:16", unset, 12, 21, 12, 16, 0, 0, 0, "", 0, 0},
		{"bugs_20", "Y-m-d H:i", "2001-12-21 12:16", 2001, 12, 21, 12, 16, 0, 0, 0, "", 0, 0},
		{"bugs_21", "F d Y H:i", "Dec 21 2001 12:16", 2001, 12, 21, 12, 16, 0, 0, 0, "", 0, 0},
		{"bugs_22", "F d  H:is", "Dec 21  12:16", unset, 12, 21, 12, 16, 0, 0, 0, "", 0, 0},
		{"bugs_23", "Y-m-d H:i:s", "2001-10-22 21:19:58", 2001, 10, 22, 21, 19, 58, 0, 0, "", 0, 0},
		{"bugs_24", "Y-m-d H:i:se", "2001-10-22 21:19:58-02", 2001, 10, 22, 21, 19, 58, 0, -7200, "", 0, 0},
		{"bugs_25", "Y-m-d H:i:se", "2001-10-22 21:19:58-0213", 2001, 10, 22, 21, 19, 58, 0, -7980, "", 0, 0},
		{"bugs_26", "Y-m-d H:i:se", "2001-10-22 21:19:58+02", 2001, 10, 22, 21, 19, 58, 0, 7200, "", 0, 0},
		{"bugs_27", "Y-m-d H:i:se", "2001-10-22 21:19:58+0213", 2001, 10, 22, 21, 19, 58, 0, 7980, "", 0, 0},
		{"bugs_28", "Y-m-d?H:i:se", "2001-10-22T21:20:58-03:40", 2001, 10, 22, 21, 20, 58, 0, -13200, "", 0, 0},
		{"bugs_29", "Ymd?Hise", "20011022T211958-2", 2001, 10, 22, 21, 19, 58, 0, -7200, "", 0, 0},
		{"bugs_30", "Ymd?Hise", "20011022T211958+0213", 2001, 10, 22, 21, 19, 58, 0, 7980, "", 0, 0},
		{"bugs_31", "Ymd?H:ie", "20011022T21:20+0215", 2001, 10, 22, 21, 20, 0, 0, 8100, "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if result.TzAbbr != tt.tzAbbr {
				t.Errorf("TzAbbr: got %s, want %s", result.TzAbbr, tt.tzAbbr)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_DateNoDay(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		z          int32
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"datenoday_00", "M Y", "Oct 2003", 2003, 10, unset, unset, unset, unset, 0, 0, "", 0, 0},
		{"datenoday_01", "d F Y", "20 October 2003", 2003, 10, 20, unset, unset, unset, 0, 0, "", 0, 0},
		{"datenoday_02", "M d", "Oct 03", unset, 10, 3, unset, unset, unset, 0, 0, "", 0, 0},
		{"datenoday_03", "M Y Hi", "Oct 2003 2045", 2003, 10, unset, 20, 45, 0, 0, 0, "", 0, 0},
		{"datenoday_04", "M Y H:i", "Oct 2003 20:45", 2003, 10, unset, 20, 45, 0, 0, 0, "", 0, 0},
		{"datenoday_05", "M Y H:i:s", "Oct 2003 20:45:37", 2003, 10, unset, 20, 45, 37, 0, 0, "", 0, 0},
		{"datenoday_06", "d F Y H:i e", "20 October 2003 00:00 CEST", 2003, 10, 20, 0, 0, 0, 0, 7200, "CEST", 0, 0},
		{"datenoday_07", "M d h:ie", "Oct 03 21:46m", unset, 10, 3, 21, 46, 0, 0, 43200, "M", 0, 0},
		{"datenoday_08", "M?Y?H:i", "Oct\t2003\t20:45", 2003, 10, unset, 20, 45, 0, 0, 0, "", 0, 0},
		{"datenoday_09", "M?d?H:ie", "Oct\t03\t21:46m", unset, 10, 3, 21, 46, 0, 0, 43200, "M", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if result.TzAbbr != tt.tzAbbr {
				t.Errorf("TzAbbr: got %s, want %s", result.TzAbbr, tt.tzAbbr)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_EpochSeconds(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		errorCount int
		warnCount  int
	}{
		{"epoch_seconds_00", "U", "-12219146756", 1582, 10, 16, 16, 34, 4, 0, 0, 0},
		{"epoch_seconds_01", "U", "12219146756", 2357, 3, 18, 7, 25, 56, 0, 0, 0},
		{"epoch_seconds_02", "U u", "-12219146756 123456", 1582, 10, 16, 16, 34, 4, 123456, 0, 0},
		{"epoch_seconds_03", "U u", "12219146756 123456", 2357, 3, 18, 7, 25, 56, 123456, 0, 0},
		{"epoch_seconds_04", "u U", "123456 -12219146756", 1582, 10, 16, 16, 34, 4, 123456, 0, 0},
		{"epoch_seconds_05", "u U", "123456 12219146756", 2357, 3, 18, 7, 25, 56, 123456, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Iso8601Long(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		errorCount int
		warnCount  int
	}{
		{"iso8601long_00", "H:i:s.u", "01:00:03.12345", unset, unset, unset, 1, 0, 3, 123450, 0, 0},
		{"iso8601long_01", "H:i:s.u", "13:03:12.45678", unset, unset, unset, 13, 3, 12, 456780, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Iso8601LongTz(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		z          int32
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"iso8601longtz_00", "H:i:s.u e", "01:00:03.12345 CET", unset, unset, unset, 1, 0, 3, 123450, 3600, "CET", 0, 0},
		{"iso8601longtz_01", "H:i:s.u e", "13:03:12.45678 CEST", unset, unset, unset, 13, 3, 12, 456780, 7200, "CEST", 0, 0},
		{"iso8601longtz_02", "H:i:s.ue", "15:57:41.0GMT", unset, unset, unset, 15, 57, 41, 0, 0, "GMT", 0, 0},
		{"iso8601longtz_03", "H:i:s.u e", "15:57:41.0 pdt", unset, unset, unset, 15, 57, 41, 0, -25200, "PDT", 0, 0},
		{"iso8601longtz_04", "H:i:s.ue", "23:41:00.0Z", unset, unset, unset, 23, 41, 0, 0, 0, "", 0, 0},
		{"iso8601longtz_05", "H:i:s.u e", "23:41:00.0 k", unset, unset, unset, 23, 41, 0, 0, 36000, "K", 0, 0},
		{"iso8601longtz_06", "H:i:s.ue", "04:05:07.789cast", unset, unset, unset, 4, 5, 7, 789000, 34200, "CAST", 0, 0},
		{"iso8601longtz_07", "H:i:s.u  e", "01:00:03.12345  +1", unset, unset, unset, 1, 0, 3, 123450, 3600, "", 0, 0},
		{"iso8601longtz_08", "H:i:s.u e", "13:03:12.45678 +0100", unset, unset, unset, 13, 3, 12, 456780, 3600, "", 0, 0},
		{"iso8601longtz_09", "H:i:s.ue", "15:57:41.0-0", unset, unset, unset, 15, 57, 41, 0, 0, "", 0, 0},
		{"iso8601longtz_10", "H:i:s.ue", "15:57:41.0-8", unset, unset, unset, 15, 57, 41, 0, -28800, "", 0, 0},
		{"iso8601longtz_11", "H:i:s.u e", "23:41:00.0 -0000", unset, unset, unset, 23, 41, 0, 0, 0, "", 0, 0},
		{"iso8601longtz_12", "H:i:s.u e", "04:05:07.789 +0930", unset, unset, unset, 4, 5, 7, 789000, 34200, "", 0, 0},
		{"iso8601longtz_13", "H:i:s.u (e)", "01:00:03.12345 (CET)", unset, unset, unset, 1, 0, 3, 123450, 3600, "CET", 0, 0},
		{"iso8601longtz_14", "H:i:s.u (e)", "13:03:12.45678 (CEST)", unset, unset, unset, 13, 3, 12, 456780, 7200, "CEST", 0, 0},
		{"iso8601longtz_15", "(e) H:i:s.u", "(CET) 01:00:03.12345", unset, unset, unset, 1, 0, 3, 123450, 3600, "CET", 0, 0},
		{"iso8601longtz_16", "(e) H:i:s.ue", "(CEST) 13:03:12.45678", unset, unset, unset, 13, 3, 12, 456780, 7200, "CEST", 0, 0},
		{"iso8601longtz_17", "H:i:s.u?(e)", "13:03:12.45678\t(CEST)", unset, unset, unset, 13, 3, 12, 456780, 7200, "CEST", 0, 0},
		{"iso8601longtz_18", "(e)?H:i:s.u", "(CEST)\t13:03:12.45678", unset, unset, unset, 13, 3, 12, 456780, 7200, "CEST", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if result.TzAbbr != tt.tzAbbr {
				t.Errorf("TzAbbr: got %s, want %s", result.TzAbbr, tt.tzAbbr)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Iso8601NoColon(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		z          int32
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"iso8601nocolon_00", "Hi", "2314", unset, unset, unset, 23, 14, 0, 0, "", 0, 0},
		{"iso8601nocolon_01", "Hi Y", "2314 2314", 2314, unset, unset, 23, 14, 0, 0, "", 0, 0},
		{"iso8601nocolon_02", "Hi e", "2314 PST", unset, unset, unset, 23, 14, 0, -28800, "PST", 0, 0},
		{"iso8601nocolon_03", "His e", "231431 CEST", unset, unset, unset, 23, 14, 31, 7200, "CEST", 0, 0},
		{"iso8601nocolon_04", "His e", "231431 CET", unset, unset, unset, 23, 14, 31, 3600, "CET", 0, 0},
		{"iso8601nocolon_05", "His", "231431", unset, unset, unset, 23, 14, 31, 0, "", 0, 0},
		{"iso8601nocolon_06", "His Y", "231431 2314", 2314, unset, unset, 23, 14, 31, 0, "", 0, 0},
		{"iso8601nocolon_07", "Hi?e", "2314\tPST", unset, unset, unset, 23, 14, 0, -28800, "PST", 0, 0},
		{"iso8601nocolon_08", "His?e", "231431\tCEST", unset, unset, unset, 23, 14, 31, 7200, "CEST", 0, 0},
		{"iso8601nocolon_09", "His?e", "231431\tCET", unset, unset, unset, 23, 14, 31, 3600, "CET", 0, 0},
		{"iso8601nocolon_10", "His?Y", "231431\t2314", 2314, unset, unset, 23, 14, 31, 0, "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if result.TzAbbr != tt.tzAbbr {
				t.Errorf("TzAbbr: got %s, want %s", result.TzAbbr, tt.tzAbbr)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Special(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		z          int32
		errorCount int
		warnCount  int
	}{
		{"special_00", "Y-m-d?H:i:se", "1998-9-15T09:05:32+4:0", 1998, 9, 15, 9, 5, 32, 0, 14400, 0, 0},
		{"special_01", "Y-m-d?H:i:se", "1998-09-15T09:05:32+04:00", 1998, 9, 15, 9, 5, 32, 0, 14400, 0, 0},
		{"special_02", "Y-m-d?H:i:s.ue", "1998-09-15T09:05:32.912+04:00", 1998, 9, 15, 9, 5, 32, 912000, 14400, 0, 0},
		{"special_03", "Y-m-d?H:i:s", "1998-09-15T09:05:32", 1998, 9, 15, 9, 5, 32, 0, 0, 0, 0},
		{"special_04", "Ymd?H:i:s", "19980915T09:05:32", 1998, 9, 15, 9, 5, 32, 0, 0, 0, 0},
		{"special_05", "Ymd?His", "19980915t090532", 1998, 9, 15, 9, 5, 32, 0, 0, 0, 0},
		{"special_06", "Y-m-d?H:i:se", "1998-09-15T09:05:32+4:9", 1998, 9, 15, 9, 5, 32, 0, 14940, 0, 0},
		{"special_07", "Y-m-d?H:i:se", "1998-9-15T09:05:32+4:30", 1998, 9, 15, 9, 5, 32, 0, 16200, 0, 0},
		{"special_08", "Y-m-d?H:i:se", "1998-09-15T09:05:32+04:9", 1998, 9, 15, 9, 5, 32, 0, 14940, 0, 0},
		{"special_09", "Y-m-d?H:i:se", "1998-9-15T09:05:32+04:30", 1998, 9, 15, 9, 5, 32, 0, 16200, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_TzCorrection(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		z          int32
		errorCount int
		warnCount  int
	}{
		{"tzcorrection_00", "e", "+4:30", 16200, 0, 0},
		{"tzcorrection_01", "e", "+4", 14400, 0, 0},
		{"tzcorrection_02", "e", "+1", 3600, 0, 0},
		{"tzcorrection_03", "P", "+14", 50400, 0, 0},
		{"tzcorrection_04", "P", "+42", 151200, 0, 0},
		{"tzcorrection_05", "P", "+4:0", 14400, 0, 0},
		{"tzcorrection_06", "P", "+4:01", 14460, 0, 0},
		{"tzcorrection_07", "T", "+4:30", 16200, 0, 0},
		{"tzcorrection_08", "T", "+401", 14460, 0, 0},
		{"tzcorrection_09", "T", "+402", 14520, 0, 0},
		{"tzcorrection_10", "T", "+430", 16200, 0, 0},
		{"tzcorrection_11", "O", "+0430", 16200, 0, 0},
		{"tzcorrection_12", "O", "+04:30", 16200, 0, 0},
		{"tzcorrection_13", "O", "+04:9", 14940, 0, 0},
		{"tzcorrection_14", "O", "+04:09", 14940, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_TzIdentifier(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"tz_identifier_00", "H:i:s.u e", "01:00:03.12345 Europe/Amsterdam", unset, unset, unset, 1, 0, 3, 123450, "", 0, 0},
		{"tz_identifier_01", "H:i:s.u e", "01:00:03.12345 America/Indiana/Knox", unset, unset, unset, 1, 0, 3, 123450, "", 0, 0},
		{"tz_identifier_02", "Y-m-d H:i:s e", "2005-07-14 22:30:41 America/Los_Angeles", 2005, 7, 14, 22, 30, 41, 0, "", 0, 0},
		{"tz_identifier_03", "Y-m-d?H:i:s?e", "2005-07-14\t22:30:41\tAmerica/Los_Angeles", 2005, 7, 14, 22, 30, 41, 0, "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug37017(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"bug37017_00", "Y-m-d H:i:s e", "2006-05-12 12:59:59 America/New_York", 2006, 5, 12, 12, 59, 59, "", 0, 0},
		{"bug37017_01", "Y-m-d H:i:s e", "2006-05-12 13:00:00 America/New_York", 2006, 5, 12, 13, 0, 0, "", 0, 0},
		{"bug37017_02", "Y-m-d H:i:s e", "2006-05-12 13:00:01 America/New_York", 2006, 5, 12, 13, 0, 1, "", 0, 0},
		{"bug37017_03", "Y-m-d H:i:s e", "2006-05-12 12:59:59 GMT", 2006, 5, 12, 12, 59, 59, "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Bug51393(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		z          int32
		errorCount int
		warnCount  int
	}{
		{"bug51393_00", "[d/M/Y:H:i:s O]", "[13/Mar/1969:23:40:00 +0100]", 1969, 3, 13, 23, 40, 0, 3600, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Combined(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		us         int64
		z          int32
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"combined_00", "D, d F Y H:i:s e", "Sat, 24 Apr 2004 21:48:40 +0200", 2004, 4, 24, 21, 48, 40, 0, 7200, "", 0, 0},
		{"combined_01", "D F d H:i:s e Y", "Sun Apr 25 01:05:41 CEST 2004", 2004, 4, 25, 1, 5, 41, 0, 7200, "CEST", 0, 0},
		{"combined_02", "D F d H:i:s Y", "Sun Apr 18 18:36:57 2004", 2004, 4, 18, 18, 36, 57, 0, 0, "", 0, 0},
		{"combined_03", "D, d F Y?H:i:s?e", "Sat, 24 Apr 2004\t21:48:40\t+0200", 2004, 4, 24, 21, 48, 40, 0, 7200, "", 0, 0},
		{"combined_04", "YmdHis e", "20040425010541 CEST", 2004, 4, 25, 1, 5, 41, 0, 7200, "CEST", 0, 0},
		{"combined_05", "YmdHis", "20040425010541", 2004, 4, 25, 1, 5, 41, 0, 0, "", 0, 0},
		{"combined_06", "Ymd?H:i:s", "19980717T14:08:55", 1998, 7, 17, 14, 8, 55, 0, 0, "", 0, 0},
		{"combined_07", "d/F/Y:H:i:s e", "10/Oct/2000:13:55:36 -0700", 2000, 10, 10, 13, 55, 36, 0, -25200, "", 0, 0},
		{"combined_08", "Y-m-d?H:i:s.u", "2001-11-29T13:20:01.123", 2001, 11, 29, 13, 20, 1, 123000, 0, "", 0, 0},
		{"combined_09", "Y-m-d?H:i:s.ue", "2001-11-29T13:20:01.123-05:00", 2001, 11, 29, 13, 20, 1, 123000, -18000, "", 0, 0},
		{"combined_10", "D F d H:i:s Y e", "Fri Aug 20 11:59:59 1993 GMT", 1993, 8, 20, 11, 59, 59, 0, 0, "GMT", 0, 0},
		{"combined_11", "D F d H:i:s Y e", "Fri Aug 20 11:59:59 1993 UTC", 1993, 8, 20, 11, 59, 59, 0, 0, "UTC", 0, 0},
		{"combined_12", "D?F?d?H:i:s?Y?e", "Fri\tAug\t20\t 11:59:59\t 1993\tUTC", unset, 8, 20, unset, 11, 59, 0, 0, "", 1, 0},
		{"combined_13", "F dS g:i e", "May 18th 5:05 UTC", unset, 5, 18, 5, 5, 0, 0, 0, "UTC", 0, 0},
		{"combined_14", "F dS g:ia e", "May 18th 5:05pm UTC", unset, 5, 18, 17, 5, 0, 0, 0, "UTC", 0, 0},
		{"combined_15", "F dS g:i a e", "May 18th 5:05 pm UTC", unset, 5, 18, 17, 5, 0, 0, 0, "UTC", 0, 0},
		{"combined_16", "F dS g:ia e", "May 18th 5:05am UTC", unset, 5, 18, 5, 5, 0, 0, 0, "UTC", 0, 0},
		{"combined_17", "F dS g:i a e", "May 18th 5:05 am UTC", unset, 5, 18, 5, 5, 0, 0, 0, "UTC", 0, 0},
		{"combined_18", "F dS Y g:ia e", "May 18th 2006 5:05pm UTC", 2006, 5, 18, 17, 5, 0, 0, 0, "UTC", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if result.US != tt.us {
				t.Errorf("US: got %d, want %d", result.US, tt.us)
			}
			if result.Z != tt.z {
				t.Errorf("Z: got %d, want %d", result.Z, tt.z)
			}
			if result.TzAbbr != tt.tzAbbr {
				t.Errorf("TzAbbr: got %s, want %s", result.TzAbbr, tt.tzAbbr)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_Day(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		errorCount int
		warnCount  int
	}{
		{"day_00", "D", "Monday", 0, 0},
		{"day_01", "D", "Wed", 0, 0},
		{"day_02", "D", "friday", 0, 0},
		{"day_03", "D", "SUNDAY", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_YearExpanded(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		errorCount int
		warnCount  int
	}{
		{"year_expanded_01", "x-m", "20-02", 20, 2, 0, 0},
		{"year_expanded_02", "x", "2002", 2002, unset, 0, 0},
		{"year_expanded_03", "x-m", "-2022-02", -2022, 2, 0, 0},
		{"year_expanded_04", "x-m", "-81120-02", -81120, 2, 0, 0},
		{"year_expanded_05", "x-m", "81120-02", 81120, 2, 0, 0},
		{"year_expanded_06", "x-m", "+82120-02", 82120, 2, 0, 0},
		{"year_expanded_07", "-x-m", "-81120-02", 81120, 2, 0, 0},
		{"year_expanded_08", "m-x", "02-81120", 81120, 2, 0, 0},
		{"year_expanded_09", "m-x", "02--81120", -81120, 2, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_YearExpandedCapitalX(t *testing.T) {
	unset := int64(-9999999)
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		errorCount int
		warnCount  int
	}{
		{"year_eXpanded_01", "X-m", "20-02", 20, 2, 0, 0},
		{"year_eXpanded_02", "X", "2002", 2002, unset, 0, 0},
		{"year_eXpanded_03", "X-m", "-2022-02", -2022, 2, 0, 0},
		{"year_eXpanded_04", "X-m", "-81120-02", -81120, 2, 0, 0},
		{"year_eXpanded_05", "X-m", "81120-02", 81120, 2, 0, 0},
		{"year_eXpanded_06", "X-m", "+82120-02", 82120, 2, 0, 0},
		{"year_eXpanded_07", "-X-m", "-81120-02", 81120, 2, 0, 0},
		{"year_eXpanded_08", "m-X", "02-81120", 81120, 2, 0, 0},
		{"year_eXpanded_09", "m-X", "02--81120", -81120, 2, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_BugGh9700(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		tzAbbr     string
		errorCount int
		warnCount  int
	}{
		{"bug_gh9700", "Y-m-d\\TH:i:sP[e]", "2022-02-18T00:00:00+01:00[Europe/Berlin]", "", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}

func TestParseDateFromFormat_BugGh11854(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		input      string
		y          int64
		m          int64
		d          int64
		h          int64
		i          int64
		s          int64
		errorCount int
		warnCount  int
	}{
		{"bug_gh11854", "D M d H:i:s Y", "Wed Aug  2 08:37:50 2023", 2023, 8, 2, 8, 37, 50, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errContainer := timelib.ParseFromFormat(tt.format, tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.Y != tt.y {
				t.Errorf("Y: got %d, want %d", result.Y, tt.y)
			}
			if result.M != tt.m {
				t.Errorf("M: got %d, want %d", result.M, tt.m)
			}
			if result.D != tt.d {
				t.Errorf("D: got %d, want %d", result.D, tt.d)
			}
			if result.H != tt.h {
				t.Errorf("H: got %d, want %d", result.H, tt.h)
			}
			if result.I != tt.i {
				t.Errorf("I: got %d, want %d", result.I, tt.i)
			}
			if result.S != tt.s {
				t.Errorf("S: got %d, want %d", result.S, tt.s)
			}
			if errContainer.ErrorCount != tt.errorCount {
				t.Errorf("ErrorCount: got %d, want %d", errContainer.ErrorCount, tt.errorCount)
			}
			if errContainer.WarningCount != tt.warnCount {
				t.Errorf("WarnCount: got %d, want %d", errContainer.WarningCount, tt.warnCount)
			}
		})
	}
}
