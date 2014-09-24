package ros

import (
	"testing"
)

func TestNameValidation(t *testing.T) {
	// Positive testing
	positives := [...]string{
		"",
		"/",
		"~",
		"foo",
		"foo/",
		"foo/bar",
		"foo/bar/",
		"foo_0/bar1_/",
		"/foo",
		"/foo/",
		"/foo/bar",
		"/foo/bar/",
		"~foo",
		"~foo/",
		"~foo/bar",
		"~foo/bar/",
	}
	for _, p := range positives {
		if !isValidName(p) {
			t.Error(p)
		}
	}

	// Negative testing
	negatives := [...]string{
		"foo//bar",
		"^foo//bar",
		"//foo",
		"0foo",
		"_0foo",
		"foo/0bar",
		"foo/_bar",
		"foo/~bar",
		"foo bar",
	}
	for _, n := range negatives {
		if isValidName(n) {
			t.Error(n)
		}
	}
}

func TestCanonicalize(t *testing.T) {
	if canonicalizeName("/") != "/" {
		t.Fail()
	}

	if canonicalizeName("/foo//bar/") != "/foo/bar" {
		t.Fail()
	}

	if canonicalizeName("foo//bar///baz/") != "foo/bar/baz" {
		t.Fail()
	}
}

func TestSpecialNamespace(t *testing.T) {
	if !isGlobalName("/foo") {
		t.Fail()
	}
	if isGlobalName("~foo") {
		t.Fail()
	}
	if isGlobalName("foo") {
		t.Fail()
	}

	if isPrivateName("/foo") {
		t.Fail()
	}
	if !isPrivateName("~foo") {
		t.Fail()
	}
	if isPrivateName("foo") {
		t.Fail()
	}
}

func TestResolution1(t *testing.T) {
	remapping := map[string]string{}
	resolver := newNameResolver("/node1", remapping)
	var result string

	result = resolver.resolve("bar")
	if result != "/bar" {
		t.Fail()
	}

	result = resolver.resolve("/bar")
	if result != "/bar" {
		t.Fail()
	}

	result = resolver.resolve("~bar")
	if result != "/node1/bar" {
		t.Fail()
	}
}

func TestResolution2(t *testing.T) {
	remapping := map[string]string{}
	resolver := newNameResolver("/go/node2", remapping)
	var result string

	result = resolver.resolve("bar")
	if result != "/go/bar" {
		t.Fail()
	}

	result = resolver.resolve("/bar")
	if result != "/bar" {
		t.Fail()
	}

	result = resolver.resolve("~bar")
	if result != "/go/node2/bar" {
		t.Fail()
	}
}

func TestResolution3(t *testing.T) {
	remapping := map[string]string{}
	resolver := newNameResolver("/go/node3", remapping)
	var result string

	result = resolver.resolve("foo/bar")
	if result != "/go/foo/bar" {
		t.Fail()
	}

	result = resolver.resolve("/foo/bar")
	if result != "/foo/bar" {
		t.Fail()
	}

	result = resolver.resolve("~foo/bar")
	if result != "/go/node3/foo/bar" {
		t.Fail()
	}
}

func TestRemapping1(t *testing.T) {
	remapping := map[string]string{
		"foo": "bar",
	}

	resolver := newNameResolver("/", remapping)
	var result string

	result = resolver.resolve("foo")
	if result != "/bar" {
		t.Fail()
	}

	result = resolver.resolve("/foo")
	if result != "/bar" {
		t.Fail()
	}
}

func TestRemapping2(t *testing.T) {
	remapping := map[string]string{
		"foo": "bar",
	}

	resolver := newNameResolver("/baz", remapping)
	var result string

	result = resolver.resolve("foo")
	if result != "/baz/bar" {
		t.Fail()
	}

	result = resolver.resolve("/baz/foo")
	if result != "/baz/bar" {
		t.Fail()
	}
}

func TestRemapping3(t *testing.T) {
	remapping := map[string]string{
		"/foo": "bar",
	}

	resolver := newNameResolver("/", remapping)
	var result string

	result = resolver.resolve("foo")
	if result != "/bar" {
		t.Fail()
	}

	result = resolver.resolve("/foo")
	if result != "/bar" {
		t.Fail()
	}
}

func TestRemapping4(t *testing.T) {
	remapping := map[string]string{
		"/foo": "bar",
	}

	resolver := newNameResolver("/baz", remapping)
	var result string

	result = resolver.resolve("/foo")
	if result != "/baz/bar" {
		t.Fail()
	}
}

func TestRemapping5(t *testing.T) {
	remapping := map[string]string{
		"/foo": "/a/b/c/bar",
	}

	resolver := newNameResolver("/baz", remapping)
	var result string

	result = resolver.resolve("/foo")
	if result != "/a/b/c/bar" {
		t.Fail()
	}
}

func TestRemapping(t *testing.T) {
	t.Fail()
}
