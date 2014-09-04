package ros

const maxUint32 = int64(^uint32(0))

func normalizeTemporal(sec int64, nsec int64) (uint32, uint32) {
	const SecondInNanosecond = 1000000000
	if nsec > SecondInNanosecond {
		sec += nsec / SecondInNanosecond
		nsec = nsec % SecondInNanosecond
	} else if nsec < 0 {
		sec += nsec/SecondInNanosecond - 1
		nsec = nsec%SecondInNanosecond + SecondInNanosecond
	}

	if sec < 0 || sec > maxUint32 {
		panic("Time is out of range")
	}

	return uint32(sec), uint32(nsec)
}

func cmpUint64(lhs, rhs uint64) int {
	var result int
	if lhs > rhs {
		result = 1
	} else if lhs < rhs {
		result = -1
	} else {
		result = 0
	}
	return result
}

type temporal struct {
	Sec  uint32
	NSec uint32
}

func (t *temporal) IsZero() bool {
	return t.Sec == 0 && t.NSec == 0
}

func (t *temporal) ToSec() float64 {
	return float64(t.Sec) + float64(t.NSec)*1e-9
}

func (t *temporal) ToNSec() uint64 {
	return uint64(t.Sec)*1000000000 + uint64(t.NSec)
}

func (t *temporal) FromSec(sec float64) {
	nsec := uint64(sec * 1e9)
	t.FromNSec(nsec)
}

func (t *temporal) FromNSec(nsec uint64) {
	t.Sec, t.NSec = normalizeTemporal(0, int64(nsec))
}

func (t *temporal) Normalize() {
	t.Sec, t.NSec = normalizeTemporal(int64(t.Sec), int64(t.NSec))
}
