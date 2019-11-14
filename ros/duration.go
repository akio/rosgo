package ros

import (
	"time"
)

//Duration type which is a wrapper for a temporal value of {sec,nsec}
type Duration struct {
	temporal
}

//NewDuration instantiates a new Duration item with given sec and nsec integers
func NewDuration(sec uint32, nsec uint32) Duration {
	sec, nsec = normalizeTemporal(int64(sec), int64(nsec))
	return Duration{temporal{sec, nsec}}
}

//Add function for adding two durations together
func (d *Duration) Add(other Duration) Duration {
	sec, nsec := normalizeTemporal(int64(d.Sec)+int64(other.Sec),
		int64(d.NSec)+int64(other.NSec))
	return Duration{temporal{sec, nsec}}
}

//Sub function for subtracting a duration from another
func (d *Duration) Sub(other Duration) Duration {
	sec, nsec := normalizeTemporal(int64(d.Sec)-int64(other.Sec),
		int64(d.NSec)-int64(other.NSec))
	return Duration{temporal{sec, nsec}}
}

//Cmp function to compare two durations
func (d *Duration) Cmp(other Duration) int {
	return cmpUint64(d.ToNSec(), other.ToNSec())
}

//Sleep function pauses go routine for duration d
func (d *Duration) Sleep() error {
	if !d.IsZero() {
		time.Sleep(time.Duration(d.ToNSec()) * time.Nanosecond)
	}
	return nil
}
