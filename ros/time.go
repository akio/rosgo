package ros

import (
	gotime "time"
)

type Time struct {
	temporal
}

func NewTime(sec uint32, nsec uint32) Time {
	sec, nsec = normalizeTemporal(int64(sec), int64(nsec))
	return Time{temporal{sec, nsec}}
}

func Now() Time {
	var t Time
	t.FromNSec(uint64(gotime.Now().UnixNano()))
	return t
}

func (t *Time) Diff(from Time) Duration {
	sec, nsec := normalizeTemporal(int64(t.Sec)-int64(from.Sec),
		int64(t.NSec)-int64(from.NSec))
	return Duration{temporal{sec, nsec}}
}

func (t *Time) Add(d Duration) Time {
	sec, nsec := normalizeTemporal(int64(t.Sec)+int64(d.Sec),
		int64(t.NSec)+int64(d.NSec))
	return Time{temporal{sec, nsec}}
}

func (t *Time) Sub(d Duration) Time {
	sec, nsec := normalizeTemporal(int64(t.Sec)-int64(d.Sec),
		int64(t.NSec)-int64(d.NSec))
	return Time{temporal{sec, nsec}}
}

func (t *Time) Cmp(other Time) int {
	return cmpUint64(t.ToNSec(), other.ToNSec())
}
