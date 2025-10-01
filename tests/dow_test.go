package tests

import (
	"testing"

	timelib "github.com/eutychus/timelib"
)

func testYMD(y, m, d int64) (iy, iw, id int64) {
	return timelib.IsoDateFromDate(y, m, d)
}

func checkYWD(t *testing.T, name string, iy, iw, id, e_iy, e_iw, e_id int64) {
	if iy != e_iy {
		t.Errorf("%s: ISO year mismatch: expected %d, got %d", name, e_iy, iy)
	}
	if iw != e_iw {
		t.Errorf("%s: ISO week mismatch: expected %d, got %d", name, e_iw, iw)
	}
	if id != e_id {
		t.Errorf("%s: ISO day mismatch: expected %d, got %d", name, e_id, id)
	}
}

func TestDOW_n0101_12_29(t *testing.T) {
	iy, iw, id := testYMD(-101, 12, 29)
	checkYWD(t, "t_n0101_12_29", iy, iw, id, -101, 52, 5)
}

func TestDOW_n0101_12_30(t *testing.T) {
	iy, iw, id := testYMD(-101, 12, 30)
	checkYWD(t, "t_n0101_12_30", iy, iw, id, -101, 52, 6)
}

func TestDOW_n0101_12_31(t *testing.T) {
	iy, iw, id := testYMD(-101, 12, 31)
	checkYWD(t, "t_n0101_12_31", iy, iw, id, -101, 52, 7)
}

func TestDOW_n0100_01_01(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 1)
	checkYWD(t, "t_n0100_01_01", iy, iw, id, -100, 1, 1)
}

func TestDOW_n0100_01_02(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 2)
	checkYWD(t, "t_n0100_01_02", iy, iw, id, -100, 1, 2)
}

func TestDOW_n0100_01_03(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 3)
	checkYWD(t, "t_n0100_01_03", iy, iw, id, -100, 1, 3)
}

func TestDOW_n0100_01_04(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 4)
	checkYWD(t, "t_n0100_01_04", iy, iw, id, -100, 1, 4)
}

func TestDOW_n0100_01_05(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 5)
	checkYWD(t, "t_n0100_01_05", iy, iw, id, -100, 1, 5)
}

func TestDOW_n0100_01_06(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 6)
	checkYWD(t, "t_n0100_01_06", iy, iw, id, -100, 1, 6)
}

func TestDOW_n0100_01_07(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 7)
	checkYWD(t, "t_n0100_01_07", iy, iw, id, -100, 1, 7)
}

func TestDOW_n0100_01_08(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 8)
	checkYWD(t, "t_n0100_01_08", iy, iw, id, -100, 2, 1)
}

func TestDOW_n0100_01_09(t *testing.T) {
	iy, iw, id := testYMD(-100, 1, 9)
	checkYWD(t, "t_n0100_01_09", iy, iw, id, -100, 2, 2)
}

func TestDOW_n0100_12_29(t *testing.T) {
	iy, iw, id := testYMD(-100, 12, 29)
	checkYWD(t, "t_n0100_12_29", iy, iw, id, -100, 52, 6)
}

func TestDOW_n0100_12_30(t *testing.T) {
	iy, iw, id := testYMD(-100, 12, 30)
	checkYWD(t, "t_n0100_12_30", iy, iw, id, -100, 52, 7)
}

func TestDOW_n0100_12_31(t *testing.T) {
	iy, iw, id := testYMD(-100, 12, 31)
	checkYWD(t, "t_n0100_12_31", iy, iw, id, -99, 1, 1)
}

func TestDOW_n0099_01_01(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 1)
	checkYWD(t, "t_n0099_01_01", iy, iw, id, -99, 1, 2)
}

func TestDOW_n0099_01_02(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 2)
	checkYWD(t, "t_n0099_01_02", iy, iw, id, -99, 1, 3)
}

func TestDOW_n0099_01_03(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 3)
	checkYWD(t, "t_n0099_01_03", iy, iw, id, -99, 1, 4)
}

func TestDOW_n0099_01_04(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 4)
	checkYWD(t, "t_n0099_01_04", iy, iw, id, -99, 1, 5)
}

func TestDOW_n0099_01_05(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 5)
	checkYWD(t, "t_n0099_01_06", iy, iw, id, -99, 1, 6)
}

func TestDOW_n0099_01_06(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 6)
	checkYWD(t, "t_n0099_01_06", iy, iw, id, -99, 1, 7)
}

func TestDOW_n0099_01_07(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 7)
	checkYWD(t, "t_n0099_01_07", iy, iw, id, -99, 2, 1)
}

func TestDOW_n0099_01_08(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 8)
	checkYWD(t, "t_n0099_01_08", iy, iw, id, -99, 2, 2)
}

func TestDOW_n0099_01_09(t *testing.T) {
	iy, iw, id := testYMD(-99, 1, 9)
	checkYWD(t, "t_n0099_01_09", iy, iw, id, -99, 2, 3)
}

func TestDOW_n0001_12_26(t *testing.T) {
	iy, iw, id := testYMD(-1, 12, 26)
	checkYWD(t, "t_n0001_12_26", iy, iw, id, -1, 51, 7)
}

func TestDOW_n0001_12_27(t *testing.T) {
	iy, iw, id := testYMD(-1, 12, 27)
	checkYWD(t, "t_n0001_12_27", iy, iw, id, -1, 52, 1)
}

func TestDOW_n0001_12_28(t *testing.T) {
	iy, iw, id := testYMD(-1, 12, 28)
	checkYWD(t, "t_n0001_12_28", iy, iw, id, -1, 52, 2)
}

func TestDOW_n0001_12_29(t *testing.T) {
	iy, iw, id := testYMD(-1, 12, 29)
	checkYWD(t, "t_n0001_12_29", iy, iw, id, -1, 52, 3)
}

func TestDOW_n0001_12_30(t *testing.T) {
	iy, iw, id := testYMD(-1, 12, 30)
	checkYWD(t, "t_n0001_12_30", iy, iw, id, -1, 52, 4)
}

func TestDOW_n0001_12_31(t *testing.T) {
	iy, iw, id := testYMD(-1, 12, 31)
	checkYWD(t, "t_n0001_12_31", iy, iw, id, -1, 52, 5)
}

func TestDOW_0000_01_01(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 1)
	checkYWD(t, "t_0000_01_01", iy, iw, id, -1, 52, 6)
}

func TestDOW_0000_01_02(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 2)
	checkYWD(t, "t_0000_01_02", iy, iw, id, -1, 52, 7)
}

func TestDOW_0000_01_03(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 3)
	checkYWD(t, "t_0000_01_03", iy, iw, id, 0, 1, 1)
}

func TestDOW_0000_01_04(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 4)
	checkYWD(t, "t_0000_01_04", iy, iw, id, 0, 1, 2)
}

func TestDOW_0000_01_05(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 5)
	checkYWD(t, "t_0000_01_05", iy, iw, id, 0, 1, 3)
}

func TestDOW_0000_01_06(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 6)
	checkYWD(t, "t_0000_01_06", iy, iw, id, 0, 1, 4)
}

func TestDOW_0000_01_07(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 7)
	checkYWD(t, "t_0000_01_07", iy, iw, id, 0, 1, 5)
}

func TestDOW_0000_01_08(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 8)
	checkYWD(t, "t_0000_01_08", iy, iw, id, 0, 1, 6)
}

func TestDOW_0000_01_09(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 9)
	checkYWD(t, "t_0000_01_09", iy, iw, id, 0, 1, 7)
}

func TestDOW_0000_01_10(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 10)
	checkYWD(t, "t_0000_01_10", iy, iw, id, 0, 2, 1)
}

func TestDOW_0000_01_11(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 11)
	checkYWD(t, "t_0000_01_11", iy, iw, id, 0, 2, 2)
}

func TestDOW_0000_01_12(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 12)
	checkYWD(t, "t_0000_01_12", iy, iw, id, 0, 2, 3)
}

func TestDOW_0000_01_13(t *testing.T) {
	iy, iw, id := testYMD(0, 1, 13)
	checkYWD(t, "t_0000_01_13", iy, iw, id, 0, 2, 4)
}
