package ros

import (
    "testing"
)

func TestResolution1(t *testing.T) {
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
    remapping := map[string]string {
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
    remapping := map[string]string {
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
    remapping := map[string]string {
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
    remapping := map[string]string {
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
    remapping := map[string]string {
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

}
