package ros

import (
	"testing"
)

func TestNewTime(t *testing.T) {
	t1 := NewTime(1, 2)
	if t1.Sec != 1 {
		t.Fail()
	}
	if t1.NSec != 2 {
		t.Fail()
	}
}

func TestTimeAdd(t *testing.T) {
	var t1 Time
	t1.FromNSec(500000000)

	var d Duration
	d.FromNSec(800000000)

	t2 := t1.Add(d)
	if t2.Sec != 1 {
		t.Error(t2.Sec)
	}
	if t2.NSec != 300000000 {
		t.Error(t2.NSec)
	}
}

func TestTimeSub(t *testing.T) {
	var t1 Time
	t1.FromNSec(1300000000)

	var d Duration
	d.FromNSec(500000000)

	t2 := t1.Sub(d)
	if t2.Sec != 0 {
		t.Error(t2.Sec)
	}
	if t2.NSec != 800000000 {
		t.Error(t2.NSec)
	}
}

func TestTimeDiff(t *testing.T) {
	var t1, t2 Time
	t1.FromNSec(1300000000)
	t2.FromNSec(500000000)

	d := t1.Diff(t2)
	if d.Sec != 0 {
		t.Error(d.Sec)
	}
	if d.NSec != 800000000 {
		t.Error(d.NSec)
	}
}
