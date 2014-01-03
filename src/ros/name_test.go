package ros

import (
    "testing"
)


func TestRelativeResolution(t *testing.T) {
    resolver := newNameResolver("/node1", nil)
    if resolver.resolve("bar") != "/bar" {
        t.Fail()
    }

    resolver := newNameResolver("/go/node1", nil)
    if resolver.resolve("bar") != "/go/bar" {
        t.Fail()
    }

    if resolver.resolve("foo/bar") != "/bar" {
        t.Fail()
    }
}


func TestGlobalResolution(t *testing.T) {
    resolver := newNameResolver("/node1", nil)
    if resolver.resolve("/bar") != "/bar" {
        t.Fail()
    }

    resolver := newNameResolver("/go/node1", nil)
    if resolver.resolve("/bar") != "/bar" {
        t.Fail()
    }

    if resolver.resolve("/foo/bar") != "/bar" {
        t.Fail()
    }
}


func TestPrivateResolution(t *testing.T) {
    resolver := newNameResolver("/node1", nil)
    if resolver.resolve("~bar") != "/node1/bar" {
        t.Fail()
    }

    resolver := newNameResolver("/go/node1", nil)
    if resolver.resolve("~bar") != "/go/node1/bar" {
        t.Fail()
    }

    if resolver.resolve("~foo/bar") != "/go/node1/foo/bar" {
        t.Fail()
    }
}


func TestRemappingRelativeToRelative(t *testing.T) {
    remapping := map[string]string {
        "foo": "bar",
    }

    resolver := newNameResolver("/", remapping)
    var result string

    if resolver.resolve("foo") != "/bar" {
        t.Fail()
    }

    if resolver.resolve("/foo") != "/bar" {
        t.Fail()
    }

    resolver = newNameResolver("/baz", remapping)

    if resolver.resolve("foo") != "/baz/bar" {
        t.Fail()
    }

    if resolver.resolve("/baz/foo") != "/baz/bar" {
        t.Fail()
    }
}


func TestRemappingGlobalToRelative(t *testing.T) {
    remapping := map[string]string {
        "/foo": "bar",
    }

    resolver := newNameResolver("/", remapping)
    var result string

    if resolver.resolve("foo") != "/bar" {
        t.Fail()
    }

    if resolver.resolve("/foo") != "/bar" {
        t.Fail()
    }

    resolver := newNameResolver("/baz", remapping)

    if resolver.resolve("/foo") != "/baz/bar" {
        t.Fail()
    }
}


func TestRemappingGlobalToGlobal(t *testing.T) {
    remapping := map[string]string {
        "/foo": "/a/b/c/bar",
    }

    resolver := newNameResolver("/baz", remapping)
    var result string

    if resolver.resolve("/foo") != "/a/b/c/bar" {
        t.Fail()
    }
}


