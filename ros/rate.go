package ros

// Impelement Rate interface
type Rate struct {
	actualCycleTime   Duration
	expectedCycleTime Duration
	start             Time
}

func NewRate(frequency float64) Rate {
	var actualCycleTime, expectedCycleTime Duration
	expectedCycleTime.FromSec(1.0 / frequency)
	start := Now()
	return Rate{actualCycleTime, expectedCycleTime, start}
}

func CycleTime(d Duration) Rate {
	var actualCycleTime Duration
	start := Now()
	return Rate{actualCycleTime, d, start}
}

func (r *Rate) CycleTime() Duration {
	return r.actualCycleTime
}

func (r *Rate) ExpectedCycleTime() Duration {
	return r.expectedCycleTime
}

func (r *Rate) Reset() {
	r.actualCycleTime = NewDuration(0, 0)
	r.start = Now()
}

func (r *Rate) Sleep() error {
	end := Now()
	diff := end.Diff(r.start)
	var remaining Duration
	if r.expectedCycleTime.Cmp(diff) >= 0 {
		remaining = r.expectedCycleTime.Sub(diff)
	}
	remaining.Sleep()
	now := Now()
	r.actualCycleTime = now.Diff(r.start)
	r.start = r.start.Add(r.expectedCycleTime)
	return nil
}
