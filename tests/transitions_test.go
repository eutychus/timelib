package tests

import (
	"math"
	"testing"

	timelib "github.com/eutychus/timelib"
)

// UTC tests
func TestTransitionsUtc01(t *testing.T) {
	var err error
	var error int
	tzi, err := timelib.ParseTzfile("UTC", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(math.MinInt64/2, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", tto.Offset)
	}
	if tto.Abbr != "UTC" {
		t.Errorf("Expected abbr UTC, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsUtc02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("UTC", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(0, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", tto.Offset)
	}
	if tto.Abbr != "UTC" {
		t.Errorf("Expected abbr UTC, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsUtc03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("UTC", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(math.MaxInt64 / 2, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", tto.Offset)
	}
	if tto.Abbr != "UTC" {
		t.Errorf("Expected abbr UTC, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

// Tokyo tests
func TestTransitionsTokyo01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(math.MinInt64, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 33539 {
		t.Errorf("Expected offset 33539, got %d", tto.Offset)
	}
	if tto.Abbr != "LMT" {
		t.Errorf("Expected abbr LMT, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-2587712401, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 33539 {
		t.Errorf("Expected offset 33539, got %d", tto.Offset)
	}
	if tto.Abbr != "LMT" {
		t.Errorf("Expected abbr LMT, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-2587712400, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 32400 {
		t.Errorf("Expected offset 32400, got %d", tto.Offset)
	}
	if tto.Abbr != "JST" {
		t.Errorf("Expected abbr JST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo04(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-2587712399, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 32400 {
		t.Errorf("Expected offset 32400, got %d", tto.Offset)
	}
	if tto.Abbr != "JST" {
		t.Errorf("Expected abbr JST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo05(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-577962001, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "JDT" {
		t.Errorf("Expected abbr JDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo06(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-577962000, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 32400 {
		t.Errorf("Expected offset 32400, got %d", tto.Offset)
	}
	if tto.Abbr != "JST" {
		t.Errorf("Expected abbr JST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo07(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-577961999, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 32400 {
		t.Errorf("Expected offset 32400, got %d", tto.Offset)
	}
	if tto.Abbr != "JST" {
		t.Errorf("Expected abbr JST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo08(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(0, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 32400 {
		t.Errorf("Expected offset 32400, got %d", tto.Offset)
	}
	if tto.Abbr != "JST" {
		t.Errorf("Expected abbr JST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsTokyo09(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Asia/Tokyo", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(math.MaxInt64, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 32400 {
		t.Errorf("Expected offset 32400, got %d", tto.Offset)
	}
	if tto.Abbr != "JST" {
		t.Errorf("Expected abbr JST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

// Amsterdam tests
func TestTransitionsAms01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(math.MinInt64, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 1172 {
		t.Errorf("Expected offset 1172, got %d", tto.Offset)
	}
	if tto.Abbr != "LMT" {
		t.Errorf("Expected abbr LMT, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsAms02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-4260212372, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 1172 {
		t.Errorf("Expected offset 1172, got %d", tto.Offset)
	}
	if tto.Abbr != "AMT" {
		t.Errorf("Expected abbr AMT, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsAms03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-1025745573, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 4772 {
		t.Errorf("Expected offset 4772, got %d", tto.Offset)
	}
	if tto.Abbr != "NST" {
		t.Errorf("Expected abbr NST, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsAms04(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(-1025745572, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 4800 {
		t.Errorf("Expected offset 4800, got %d", tto.Offset)
	}
	if tto.Abbr != "+0120" {
		t.Errorf("Expected abbr +0120, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsAms05(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(811904399, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", tto.Offset)
	}
	if tto.Abbr != "CEST" {
		t.Errorf("Expected abbr CEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsAms06(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(811904440, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", tto.Offset)
	}
	if tto.Abbr != "CET" {
		t.Errorf("Expected abbr CET, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsAms07(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(828234000, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", tto.Offset)
	}
	if tto.Abbr != "CEST" {
		t.Errorf("Expected abbr CEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsAms08(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(846377999, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", tto.Offset)
	}
	if tto.Abbr != "CEST" {
		t.Errorf("Expected abbr CEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsAms09(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(846378000, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", tto.Offset)
	}
	if tto.Abbr != "CET" {
		t.Errorf("Expected abbr CET, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsAms10(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(846378001, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", tto.Offset)
	}
	if tto.Abbr != "CET" {
		t.Errorf("Expected abbr CET, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsAms11(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(859683599, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", tto.Offset)
	}
	if tto.Abbr != "CET" {
		t.Errorf("Expected abbr CET, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsAms12(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(859683600, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", tto.Offset)
	}
	if tto.Abbr != "CEST" {
		t.Errorf("Expected abbr CEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsAms13(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Amsterdam", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(859683600, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 7200 {
		t.Errorf("Expected offset 7200, got %d", tto.Offset)
	}
	if tto.Abbr != "CEST" {
		t.Errorf("Expected abbr CEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

// Canberra tests
func TestTransitionsCan01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1193500799, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "AEST" {
		t.Errorf("Expected abbr AEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsCan02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1193500800, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "AEDT" {
		t.Errorf("Expected abbr AEDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsCan03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1193500801, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "AEDT" {
		t.Errorf("Expected abbr AEDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsCan04(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1207411199, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "AEDT" {
		t.Errorf("Expected abbr AEDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsCan05(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1207411200, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "AEST" {
		t.Errorf("Expected abbr AEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsCan06(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1207411201, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "AEST" {
		t.Errorf("Expected abbr AEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsCan07(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1223135999, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "AEST" {
		t.Errorf("Expected abbr AEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsCan08(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1223136000, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "AEDT" {
		t.Errorf("Expected abbr AEDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsCan09(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1223136001, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "AEDT" {
		t.Errorf("Expected abbr AEDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsCan10(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1238860799, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "AEDT" {
		t.Errorf("Expected abbr AEDT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsCan11(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1238860800, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "AEST" {
		t.Errorf("Expected abbr AEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsCan12(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Canberra", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1238860801, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 36000 {
		t.Errorf("Expected offset 36000, got %d", tto.Offset)
	}
	if tto.Abbr != "AEST" {
		t.Errorf("Expected abbr AEST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

// Lord Howe tests
func TestTransitionsLH01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1207407599, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1207407600, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 37800 {
		t.Errorf("Expected offset 37800, got %d", tto.Offset)
	}
	if tto.Abbr != "+1030" {
		t.Errorf("Expected abbr +1030, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsLH03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1207407601, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 37800 {
		t.Errorf("Expected offset 37800, got %d", tto.Offset)
	}
	if tto.Abbr != "+1030" {
		t.Errorf("Expected abbr +1030, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsLH04(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1317482999, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 37800 {
		t.Errorf("Expected offset 37800, got %d", tto.Offset)
	}
	if tto.Abbr != "+1030" {
		t.Errorf("Expected abbr +1030, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsLH05(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1317483000, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH06(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1317483001, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH07(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1365260399, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH08(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1365260400, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 37800 {
		t.Errorf("Expected offset 37800, got %d", tto.Offset)
	}
	if tto.Abbr != "+1030" {
		t.Errorf("Expected abbr +1030, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsLH09(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1365260401, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 37800 {
		t.Errorf("Expected offset 37800, got %d", tto.Offset)
	}
	if tto.Abbr != "+1030" {
		t.Errorf("Expected abbr +1030, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsLH10(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1293800399, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH11(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1293800400, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH12(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1293800401, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH13(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1293839999, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH14(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1293840000, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsLH15(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Australia/Lord_Howe", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1293840001, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 39600 {
		t.Errorf("Expected offset 39600, got %d", tto.Offset)
	}
	if tto.Abbr != "+11" {
		t.Errorf("Expected abbr +11, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

// Fiji tests
func TestTransitionsFiji01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Pacific/Fiji", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1608386399, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 43200 {
		t.Errorf("Expected offset 43200, got %d", tto.Offset)
	}
	if tto.Abbr != "+12" {
		t.Errorf("Expected abbr +12, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsFiji02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Pacific/Fiji", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1608386400, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 46800 {
		t.Errorf("Expected offset 46800, got %d", tto.Offset)
	}
	if tto.Abbr != "+13" {
		t.Errorf("Expected abbr +13, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsFiji03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Pacific/Fiji", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1608386401, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 46800 {
		t.Errorf("Expected offset 46800, got %d", tto.Offset)
	}
	if tto.Abbr != "+13" {
		t.Errorf("Expected abbr +13, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

// Chatham tests
func TestTransitionsChat01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Pacific/Chatham", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1821880799, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 45900 {
		t.Errorf("Expected offset 45900, got %d", tto.Offset)
	}
	if tto.Abbr != "+1245" {
		t.Errorf("Expected abbr +1245, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsChat02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Pacific/Chatham", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1821880800, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 49500 {
		t.Errorf("Expected offset 49500, got %d", tto.Offset)
	}
	if tto.Abbr != "+1345" {
		t.Errorf("Expected abbr +1345, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsChat03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Pacific/Chatham", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1821880801, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 49500 {
		t.Errorf("Expected offset 49500, got %d", tto.Offset)
	}
	if tto.Abbr != "+1345" {
		t.Errorf("Expected abbr +1345, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

// Dublin tests (Ireland has negative DST)
func TestTransitionsEire01(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Dublin", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(859683599, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", tto.Offset)
	}
	if tto.Abbr != "GMT" {
		t.Errorf("Expected abbr GMT, got %s", tto.Abbr)
	}
	if tto.IsDst != 1 {
		t.Errorf("Expected is_dst 1, got %d", tto.IsDst)
	}
}

func TestTransitionsEire02(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Dublin", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(859683600, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", tto.Offset)
	}
	if tto.Abbr != "IST" {
		t.Errorf("Expected abbr IST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

func TestTransitionsEire03(t *testing.T) {
	var error int
	tzi, err := timelib.ParseTzfile("Europe/Dublin", timelib.BuiltinDB(), &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(859683601, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != 3600 {
		t.Errorf("Expected offset 3600, got %d", tto.Offset)
	}
	if tto.Abbr != "IST" {
		t.Errorf("Expected abbr IST, got %s", tto.Abbr)
	}
	if tto.IsDst != 0 {
		t.Errorf("Expected is_dst 0, got %d", tto.IsDst)
	}
}

// New York with modified full year DST
func TestTransitionsNY01(t *testing.T) {
	var error int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("New_York_mod_Full_Year_DST", testDirectory, &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1615483094, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != -14400 {
		t.Errorf("Expected offset -14400, got %d", tto.Offset)
	}
	if tto.Abbr != "EDT" {
		t.Errorf("Expected abbr EDT, got %s", tto.Abbr)
	}
}

func TestTransitionsNY02(t *testing.T) {
	var error int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("New_York_mod_Full_Year_DST", testDirectory, &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1609477199, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != -14400 {
		t.Errorf("Expected offset -14400, got %d", tto.Offset)
	}
	if tto.Abbr != "EDT" {
		t.Errorf("Expected abbr EDT, got %s", tto.Abbr)
	}
}

func TestTransitionsNY03(t *testing.T) {
	var error int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("New_York_mod_Full_Year_DST", testDirectory, &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1609477200, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != -14400 {
		t.Errorf("Expected offset -14400, got %d", tto.Offset)
	}
	if tto.Abbr != "EDT" {
		t.Errorf("Expected abbr EDT, got %s", tto.Abbr)
	}
}

func TestTransitionsNY04(t *testing.T) {
	var error int
	testDirectory, err2 := timelib.Zoneinfo("tests/c/files")
	if err2 != nil {
		t.Fatalf("Zoneinfo error: %v", err2)
	}

	tzi, err := timelib.ParseTzfile("New_York_mod_Full_Year_DST", testDirectory, &error)
	if err != nil || tzi == nil {
		t.Fatalf("Failed to parse timezone")
	}

	tto := timelib.GetTimeZoneInfo(1609477201, tzi)
	if tto == nil {
		t.Fatalf("tto is nil")
	}

	if tto.Offset != -14400 {
		t.Errorf("Expected offset -14400, got %d", tto.Offset)
	}
	if tto.Abbr != "EDT" {
		t.Errorf("Expected abbr EDT, got %s", tto.Abbr)
	}
}
