package ros

import (
	"testing"
)

func TestLoadJsonFromString(t *testing.T) {
	value, err := loadParamFromString("42")
	if err != nil {
		t.Error(err)
	}
	i, ok := value.(float64)
	if !ok {
		t.Fail()
	}
	if i != 42.0 {
		t.Error(i)
	}
}
