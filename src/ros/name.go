package ros

import (
	"regexp"
)

const (
	Sep         = "/"
	GlobalNS    = "/"
	PrivateName = "~"
	Remap       = ":="
)

type Remapping map[string]string

type NameResolver struct {
	namespace string
	remapping Remapping
}

func newNameResolver(namespace string, remapping Remapping) *NameResolver {
	n := new(NameResolver)
	n.namespace = namespace
	n.remapping = remapping
	return n
}

func (n *NameResolver) resolve(name string) string {
	return name
}

func IsValidName(name string) bool {
	if len(name) == 0 {
		return true
	}
	if matched, _ := regexp.MatchString("^[~/]?([a-zA-Z]\\w*/)*[a-zA-Z]\\w*/?$", name); !matched {
		return false
	}
	return true
}
