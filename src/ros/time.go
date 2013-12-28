package ros

import (
    "time"
)


func canonicalize(sec int64, nsec int64) (uint32, uint32) {
    for nsec > 10e9 {
        sec += 1
        nsec -= 10e9
    }
    for nsec < 0 {
        sec -= 1
        nsec += 10e9
    }
    return uint32(sec), uint32(nsec)
}

type timeBase struct {
    sec uint32
    nsec uint32
}


func (t *timeBase) Sec() uint32 {
    return t.sec
}

func (t *timeBase) SetSec(sec uint32) {
    t.sec = sec
}

func (t *timeBase) NSec() uint32 {
    return t.nsec
}

func (t *timeBase) SetNSec(nsec uint32) {
    t.nsec = nsec
}

func (t *timeBase) Set(sec uint32, nsec uint32) {
    t.sec = sec
    t.nsec = nsec
}


/// 
func (t *timeBase) IsZero() bool {
    return t.sec == 0 && t.nsec == 0
}


func (t *timeBase) ToSec() float64 {
    return float64(t.sec) + float64(t.nsec) * 10e9
}


func (t *timeBase) ToNSec() uint64 {
    return uint64(t.sec) * 10e9 + uint64(t.nsec)
}


func (t *timeBase) FromSec(sec float64) {
    nsec := uint64(sec * 10e9)
    t.FromNSec(nsec)
}


func (t *timeBase) FromNSec(nsec uint64) {
    s, ns := canonicalize(0, int64(nsec))
    t.sec = s
    t.nsec = ns
}


func (t *timeBase) Add(d Duration) Time {
    sec, nsec := canonicalize(int64(t.sec) + int64(d.Sec()),
                              int64(t.nsec) + int64(d.NSec()))
    t.Set(sec, nsec)
    return t
}


func (t *timeBase) Sub(d Duration) Time {
    sec, nsec := canonicalize(int64(t.sec) - int64(d.Sec()),
                              int64(t.nsec) - int64(d.NSec()))
    t.Set(sec, nsec)
    return t
}


func (t *timeBase) Diff(from Time) Duration {
    sec, nsec := canonicalize(int64(t.sec) - int64(from.Sec()),
                              int64(t.nsec) - int64(from.NSec()))
    return &duration { timeBase{ sec, nsec } }
}


func (t *timeBase) Cmp(other Time) (result int) {
    lhs := t.ToNSec()
    rhs := other.ToNSec()
    if lhs > rhs {
        result = 1
    } else if lhs < rhs {
        result = -1
    } else {
        result = 0
    }
    return result
}



type _Time struct {
    timeBase
}


func now() Time {
    t := new(_Time)
    t.FromNSec(uint64(time.Now().UnixNano()))
    return t
}


type duration struct {
    timeBase
}


func (d *duration) Sleep() error {
    time.Sleep(time.Duration(d.ToNSec()))
    return nil
}

