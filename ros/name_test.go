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

func TestCanonicalizeName(t *testing.T) {
	if canonicalizeName("/") != "/" {
		t.Fail()
	}

	if canonicalizeName("/foo//bar/") != "/foo/bar" {
		t.Fail()
	}

	if canonicalizeName("foo//bar///baz/") != "foo/bar/baz" {
		t.Fail()
	}

	if canonicalizeName("~foo//bar///baz/") != "~foo/bar/baz" {
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
	remapping := NameMap{}
	resolver := newNameResolver("/", "node1", remapping)
	var result string

	result = resolver.resolve("bar")
	if result != "/bar" {
		t.Error(result)
	}

	result = resolver.resolve("/bar")
	if result != "/bar" {
		t.Error(result)
	}

	result = resolver.resolve("~bar")
	if result != "/node1/bar" {
		t.Error(result)
	}
}

func TestResolution2(t *testing.T) {
	remapping := NameMap{}
	resolver := newNameResolver("/go", "node2", remapping)
	var result string

	result = resolver.resolve("bar")
	if result != "/go/bar" {
		t.Error(result)
	}

	result = resolver.resolve("/bar")
	if result != "/bar" {
		t.Error(result)
	}

	result = resolver.resolve("~bar")
	if result != "/go/node2/bar" {
		t.Error(result)
	}
}

func TestResolution3(t *testing.T) {
	remapping := NameMap{}
	resolver := newNameResolver("/go", "node3", remapping)
	var result string

	result = resolver.resolve("foo/bar")
	if result != "/go/foo/bar" {
		t.Error(result)
	}

	result = resolver.resolve("/foo/bar")
	if result != "/foo/bar" {
		t.Error(result)
	}

	result = resolver.resolve("~foo/bar")
	if result != "/go/node3/foo/bar" {
		t.Error(result)
	}
}

func TestNameMap1(t *testing.T) {
	remapping := NameMap{
		"foo": "bar",
	}

	resolver := newNameResolver("/", "mynode", remapping)
	var result string

	result = resolver.remap("foo")
	if result != "/bar" {
		t.Error(result)
	}

	result = resolver.remap("/foo")
	if result != "/bar" {
		t.Error(result)
	}
}

func TestNameMap2(t *testing.T) {
	remapping := NameMap{
		"foo": "bar",
	}

	resolver := newNameResolver("/baz", "mynode", remapping)
	var result string

	result = resolver.remap("foo")
	if result != "/baz/bar" {
		t.Error(result)
		t.Error(resolver.mapping)
	}

	result = resolver.remap("/baz/foo")
	if result != "/baz/bar" {
		t.Error(result)
	}
}

func TestNameMap3(t *testing.T) {
	remapping := NameMap{
		"/foo": "bar",
	}

	resolver := newNameResolver("/", "mynode", remapping)
	var result string

	result = resolver.remap("foo")
	if result != "/bar" {
		t.Error(result)
	}

	result = resolver.remap("/foo")
	if result != "/bar" {
		t.Error(result)
	}
}

func TestNameMap4(t *testing.T) {
	remapping := NameMap{
		"/foo": "bar",
	}

	resolver := newNameResolver("/baz", "mynode", remapping)
	var result string

	result = resolver.remap("/foo")
	if result != "/baz/bar" {
		t.Error(resolver.mapping)
		t.Error(result)
	}
}

func TestNameMap5(t *testing.T) {
	remapping := NameMap{
		"/foo": "/a/b/c/bar",
	}

	resolver := newNameResolver("/baz", "mynode", remapping)
	var result string

	result = resolver.remap("/foo")
	if result != "/a/b/c/bar" {
		t.Error(result)
	}
}

func TestGetNamespace(t *testing.T) {
	var ns string
	ns = getNamespace("")
	if ns != "/" {
		t.Error(ns)
	}

	ns = getNamespace("/")
	if ns != "/" {
		t.Error(ns)
	}

	ns = getNamespace("/foo")
	if ns != "/" {
		t.Error(ns)
	}

	ns = getNamespace("/foo/")
	if ns != "/" {
		t.Error(ns)
	}

	ns = getNamespace("/foo/bar")
	if ns != "/foo/" {
		t.Error(ns)
	}

	ns = getNamespace("/foo/bar/baz")
	if ns != "/foo/bar/" {
		t.Error(ns)
	}
}

func TestProcessArguments(t *testing.T) {
	args := []string{
		"foo:=bar",
		"_param:=value",
		"__master:=http://localhost:11311",
		"foo",
		"42",
	}

	mapping, params, specials, rest := processArguments(args)
	if mapping["foo"] != "bar" {
		t.Fail()
	}
	if params["param"] != "value" {
		t.Fail()
	}
	if specials["__master"] != "http://localhost:11311" {
		t.Fail()
	}
	if len(rest) != 2 {
		t.Fail()
	}
	if rest[0] != "foo" || rest[1] != "42" {
		t.Fail()
	}
}
