package ros

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewRate(t *testing.T) {
	r := NewRate(100)
	if !r.actualCycleTime.IsZero() {
		t.Fail()
	}
	if r.expectedCycleTime.ToSec() != 0.01 {
		t.Fail()
	}
}

func TestCycleTime(t *testing.T) {
	const MeasureTolerance int64 = 1000000

	var d Duration
	d.FromSec(0.01)
	r := CycleTime(d)
	if !r.actualCycleTime.IsZero() {
		t.Fail()
	}
	if r.expectedCycleTime.ToSec() != 0.01 {
		t.Fail()
	}

	start := time.Now().UnixNano()
	r.Sleep()
	end := time.Now().UnixNano()

	actual := r.CycleTime()
	elapsed := end - start
	delta := int64(actual.ToNSec()) - elapsed
	if delta < 0 {
		delta = -delta
	}
	if delta > MeasureTolerance {
		t.Error(delta)
	}
}

func TestRateReset(t *testing.T) {
	r := NewRate(100)
	r.Sleep()

	if r.actualCycleTime.IsZero() {
		t.Fail()
	}
	r.Reset()
	if !r.actualCycleTime.IsZero() {
		t.Fail()
	}
}

func TestRateSleep(t *testing.T) {
	// The jitter tolerance (5msec) doesn't have strong basis.
	const JitterTolerance int64 = 5000000
	ct := NewDuration(0, 100000000) // 10msec
	r := CycleTime(ct)
	if ct.Cmp(r.ExpectedCycleTime()) != 0 {
		t.Fail()
	}
	for i := 0; i < 10; i++ {
		start := time.Now().UnixNano()
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		r.Sleep()
		end := time.Now().UnixNano()

		elapsed := end - start
		delta := elapsed - int64(ct.ToNSec())
		if delta < 0 {
			delta = -delta
		}
		if delta > JitterTolerance {
			actual := r.CycleTime()
			t.Errorf("expected: %d  actual: %d  measured: %d  delta: %d",
				ct.ToNSec(), actual.ToNSec(), elapsed, delta)
		}
	}
}
