package ros

import (
	"testing"
)

func TestNormalizeTemporal(t *testing.T) {
	var sec, nsec uint32
	sec, nsec = normalizeTemporal(1, 2)
	if sec != 1 || nsec != 2 {
		t.Error(sec, nsec)
	}

	sec, nsec = normalizeTemporal(1, 2000000001)
	if sec != 3 || nsec != 1 {
		t.Error(sec, nsec)
	}

	sec, nsec = normalizeTemporal(3, -2000000001)
	if sec != 0 || nsec != 999999999 {
		t.Error(sec, nsec)
	}
}

func TestTemporalIsZero(t *testing.T) {
	var t1 temporal
	if !t1.IsZero() {
		t.Fail()
	}

	t1.Sec = 1
	t1.NSec = 0
	if t1.IsZero() {
		t.Fail()
	}

	t1.Sec = 0
	t1.NSec = 1
	if t1.IsZero() {
		t.Fail()
	}

	t1.Sec = 1
	t1.NSec = 1
	if t1.IsZero() {
		t.Fail()
	}
}

func TestTemporalSet(t *testing.T) {
	t1 := temporal{1, 2}
	if t1.Sec != 1 {
		t.Fail()
	}
	if t1.NSec != 2 {
		t.Fail()
	}
}

func TestTemporalToSec(t *testing.T) {
	t1 := temporal{1, 500000000}
	if t1.ToSec() != 1.5 {
		t.Error(t1.ToSec())
	}
	t1.Sec, t1.NSec = 0, 1500000000
	if t1.ToSec() != 1.5 {
		t.Error(t1.ToSec())
	}
}

func TestTemporalToNSec(t *testing.T) {
	t1 := temporal{1, 500000000}
	if t1.ToNSec() != 1500000000 {
		t.Fail()
	}
	t1.Sec, t1.NSec = 0, 1500000000
	if t1.ToNSec() != 1500000000 {
		t.Fail()
	}
}

func TestTemporalFromSec(t *testing.T) {
	var t1 temporal
	t1.FromSec(1.5)
	if t1.Sec != 1 || t1.NSec != 500000000 {
		t.Fail()
	}
}

func TestTemporalFromNSec(t *testing.T) {
	var t1 temporal
	t1.FromNSec(1500000000)
	if t1.Sec != 1 || t1.NSec != 500000000 {
		t.Fail()
	}
}
