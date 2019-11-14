package ros

//Rate interface is a struct of Durations actual and expected cycle time, and start Time
type Rate struct {
	actualCycleTime   Duration
	expectedCycleTime Duration
	start             Time
}

//NewRate returns a new Rate object of given frequency
func NewRate(frequency float64) Rate {
	var actualCycleTime, expectedCycleTime Duration
	expectedCycleTime.FromSec(1.0 / frequency)
	start := Now()
	return Rate{actualCycleTime, expectedCycleTime, start}
}

//CycleTime returns a new Rate object of expected cycle time of duration given
func CycleTime(d Duration) Rate {
	var actualCycleTime Duration
	start := Now()
	return Rate{actualCycleTime, d, start}
}

//CycleTime returns duration of the Cycle time of a Rate object
func (r *Rate) CycleTime() Duration {
	return r.actualCycleTime
}

//ExpectedCycleTime returns duration of the Expected Cycle time of a Rate object
func (r *Rate) ExpectedCycleTime() Duration {
	return r.expectedCycleTime
}

//Reset sets actual Cycle time and start time of Rate to 0/now
func (r *Rate) Reset() {
	r.actualCycleTime = NewDuration(0, 0)
	r.start = Now()
}

//Sleep pauses go routine for time = expectedCycleTime - (Now - Rate start)
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
	r.start = now
	return nil
}
