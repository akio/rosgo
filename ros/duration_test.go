package ros

import (
	"testing"
	"time"
)

func TestNewDuration(t *testing.T) {
	d := NewDuration(1, 2)
	if d.Sec != 1 {
		t.Fail()
	}
	if d.NSec != 2 {
		t.Fail()
	}
}

func TestDurationAdd(t *testing.T) {
	var d1, d2 Duration
	d1.FromNSec(500000000)
	d2.FromNSec(800000000)

	d3 := d1.Add(d2)
	if d3.Sec != 1 {
		t.Error(d3.Sec)
	}
	if d3.NSec != 300000000 {
		t.Error(d3.NSec)
	}
}

func TestDurationSub(t *testing.T) {
	var d1, d2 Duration
	d1.FromNSec(1300000000)
	d2.FromNSec(500000000)

	d3 := d1.Sub(d2)
	if d3.Sec != 0 {
		t.Error(d3.Sec)
	}
	if d3.NSec != 800000000 {
		t.Error(d3.NSec)
	}
}

func TestDurationSleep(t *testing.T) {
	d := NewDuration(1, 100000000)
	start := time.Now().UnixNano()
	d.Sleep()
	end := time.Now().UnixNano()
	// The jitter tolerance (5msec) doesn't have strong basis.
	const Tolerance int64 = 5000000
	elapsed := end - start
	delta := elapsed - int64(d.ToNSec())
	if delta < 0 {
		delta = -delta
	}
	if delta > Tolerance {
		t.Errorf("expected: %d  actual0: %d  delta: %d", d.ToNSec(), elapsed, delta)
	}
}
