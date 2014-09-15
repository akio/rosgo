package ros

import (
	"regexp"
)

const (
	Sep       = '/'
	GlobalNS  = '/'
	PrivateNS = '~'
	Remap     = ":="
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

func isValidName(name string) bool {
	if len(name) == 0 {
		return true
	}
	if matched, _ := regexp.MatchString("^[~/]?([a-zA-Z]\\w*/)*[a-zA-Z]\\w*/?$", name); !matched {
		return false
	}
	return true
}

func isGlobalName(name string) bool {
	return len(name) > 0 && name[0] == GlobalNS
}

func isPrivateName(name string) bool {
	return len(name) > 0 && name[0] == PrivateNS
}
