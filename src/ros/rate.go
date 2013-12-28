package ros

// Impelement Rate interface
type rate struct {
    actualCycleTime Duration
    expectedCycleTime Duration
    start Time    
}


func newRate(frequency float64) *rate {
    r := new(rate)
    d := NewDuration()
    d.FromSec(1.0 / frequency)
    r.expectedCycleTime = d
    r.start = Now()
    return r
}


func newRateFromCycleTime(d Duration) *rate {
    r := new(rate)
    r.expectedCycleTime = d
    r.start = Now()
    return r
}


func (r *rate) CycleTime() Duration {
    return r.actualCycleTime
}


func (r *rate) ExpectedCycleTime() Duration {
    return r.expectedCycleTime
}


func (r *rate) Reset() {
   r.start = Now()
}

func (r *rate) Sleep() error {
    now := Now()
    r.start.Diff(now)
    d := now.Diff(r.start)
    d.Sleep()
    return nil
}


