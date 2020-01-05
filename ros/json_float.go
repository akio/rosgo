package ros

// IMPORT REQUIRED PACKAGES.

import (
	"encoding/json"
	"math"
	"strconv"
)

// DEFINE PUBLIC STRUCTURES.

// DEFINE PRIVATE STRUCTURES.

type JsonFloat32 struct {
	F float32
}

type JsonFloat64 struct {
	F float64
}

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

// DEFINE PUBLIC STATIC FUNCTIONS.

// DEFINE PRIVATE RECEIVER FUNCTIONS.

func (f JsonFloat32) String() string {
	return strconv.FormatFloat(float64(f.F), 'f', 5, 32)
}

func (f JsonFloat32) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(f.F)) {
		return json.Marshal("nan")
	} else if math.IsInf(float64(f.F), 1) {
		return json.Marshal("+inf")
	} else if math.IsInf(float64(f.F), -1) {
		return json.Marshal("-inf")
	}
	return json.Marshal(f.F)
}

func (f JsonFloat64) String() string {
	return strconv.FormatFloat(f.F, 'f', 5, 64)
}

func (f JsonFloat64) MarshalJSON() ([]byte, error) {
	if math.IsNaN(f.F) {
		return json.Marshal("nan")
	} else if math.IsInf(f.F, 1) {
		return json.Marshal("+inf")
	} else if math.IsInf(f.F, -1) {
		return json.Marshal("-inf")
	}
	return json.Marshal(f.F)
}

// ALL DONE.
