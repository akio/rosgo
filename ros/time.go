package ros

import (
	gotime "time"
)

//Time struct contains a temporal value {sec,nsec}
type Time struct {
	temporal
}

//NewTime creates a Time object of given integers {sec,nsec}
func NewTime(sec uint32, nsec uint32) Time {
	sec, nsec = normalizeTemporal(int64(sec), int64(nsec))
	return Time{temporal{sec, nsec}}
}

//Now creates a Time object of value Now
func Now() Time {
	var t Time
	t.FromNSec(uint64(gotime.Now().UnixNano()))
	return t
}

//Diff returns difference of two Time objects as a Duration
func (t *Time) Diff(from Time) Duration {
	sec, nsec := normalizeTemporal(int64(t.Sec)-int64(from.Sec),
		int64(t.NSec)-int64(from.NSec))
	return Duration{temporal{sec, nsec}}
}

//Add returns sum of Time and Duration given
func (t *Time) Add(d Duration) Time {
	sec, nsec := normalizeTemporal(int64(t.Sec)+int64(d.Sec),
		int64(t.NSec)+int64(d.NSec))
	return Time{temporal{sec, nsec}}
}

//Sub returns subtraction of Time and Duration given
func (t *Time) Sub(d Duration) Time {
	sec, nsec := normalizeTemporal(int64(t.Sec)-int64(d.Sec),
		int64(t.NSec)-int64(d.NSec))
	return Time{temporal{sec, nsec}}
}

//Cmp returns int comparison of two Time objects
func (t *Time) Cmp(other Time) int {
	return cmpUint64(t.ToNSec(), other.ToNSec())
}
