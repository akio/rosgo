package ros

import (
	"testing"
)

func TestUnique(t *testing.T) {
	var data = []string{
		"a", "b", "c", "a",
	}
	result := unique(data)
	if len(result) != 3 {
		t.Fail()
	}
}

func TestUnion(t *testing.T) {
	var a = []string{
		"a", "b", "c", "a",
	}

	var b = []string{
		"a", "b", "d",
	}

	result := setUnion(a, b)
	if len(result) != 4 {
		t.Errorf("Expected 4 but %d", len(result))
	}

	for _, k := range []string{"a", "b", "c", "d"} {
		if !contains(result, k) {
			t.Fail()
		}
	}
}

func TestDifference(t *testing.T) {
	var a = []string{
		"a", "b", "c", "a",
	}

	var b = []string{
		"a", "b", "d", "e",
	}

	result := setDifference(a, b)
	if len(result) != 1 {
		t.Errorf("Expected 1 but %d", len(result))
	}
	for _, k := range []string{"c"} {
		if !contains(result, k) {
			t.Fail()
		}
	}

	result = setDifference(b, a)
	if len(result) != 2 {
		t.Errorf("Expected 2 but %d", len(result))
	}
	for _, k := range []string{"d", "e"} {
		if !contains(result, k) {
			t.Fail()
		}
	}
}
