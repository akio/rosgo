package ros

import (
	"time"
)

type Duration struct {
	temporal
}

func NewDuration(sec uint32, nsec uint32) Duration {
	sec, nsec = normalizeTemporal(int64(sec), int64(nsec))
	return Duration{temporal{sec, nsec}}
}

func (d *Duration) Add(other Duration) Duration {
	sec, nsec := normalizeTemporal(int64(d.Sec)+int64(other.Sec),
		int64(d.NSec)+int64(other.NSec))
	return Duration{temporal{sec, nsec}}
}

func (d *Duration) Sub(other Duration) Duration {
	sec, nsec := normalizeTemporal(int64(d.Sec)-int64(other.Sec),
		int64(d.NSec)-int64(other.NSec))
	return Duration{temporal{sec, nsec}}
}

func (d *Duration) Cmp(other Duration) int {
	return cmpUint64(d.ToNSec(), other.ToNSec())
}

func (d *Duration) Sleep() error {
	if !d.IsZero() {
		time.Sleep(time.Duration(d.ToNSec()) * time.Nanosecond)
	}
	return nil
}
