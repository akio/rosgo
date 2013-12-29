package ros

import (
    "testing"
)


func TestTimeBaseIsZero(t *testing.T) {
    t := timeBase { 0, 0 }
    if !t.IsZero() {
        t.Fail()
    }

    t.sec = 1
    t.nsec = 0
    if t.IsZero() {
        t.Fail()
    }

    t.sec = 0
    t.nsec = 1
    if t.IsZero() {
        t.Fail()
    }

    t.sec = 1
    t.nsec = 1
    if t.IsZero() {
        t.Fail()
    }
}


func TestTimeBaseSet(t *testing.T) {
    var t timeBase
    t.Set(1, 2)
}


func TestTimeBaseToSec(t *testing.T) {
}


func TestTimeBaseToNSec(t *testing.T) {

}


func TestTimeBaseFromSec(t *testing.T) {
}


func TestTimeBaseFromNSec(t *testing.T) {
}


func TestTimeBaseAdd(t *testing.T) {
}


func TestTimeBaseSub(t *testing.T) {
}



func TestTimeBaseDiff(t *testing.T) { 

}

func TestTimeBaseCmp(t *testing.T) {
}

func TestNewTime(t *testing.T) {

}


func TestNewDuration(t *testing.T) {
    d := newDuration()

}

func TestDurationSleep(t *testing.T) {

}


