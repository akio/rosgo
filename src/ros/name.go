package ros

import (
	"regexp"
	"strings"
)

const (
	Sep       = "/"
	GlobalNS  = "/"
	PrivateNS = "~"
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
	if name == "/" || name == "~" {
		return true
	}
	if matched, _ := regexp.MatchString("^[~/]?([a-zA-Z]\\w*/)*[a-zA-Z]\\w*/?$", name); !matched {
		return false
	}
	return true
}

func isGlobalName(name string) bool {
	return len(name) > 0 && name[0:1] == GlobalNS
}

func isPrivateName(name string) bool {
	return len(name) > 0 && name[0:1] == PrivateNS
}

func canonicalizeName(name string) string {
	if name == GlobalNS {
		return name
	} else {
		components := []string{}
		for _, word := range strings.Split(name, Sep) {
			if len(word) > 0 {
				components = append(components, word)
			}
		}
		if name[0:1] == GlobalNS {
			return GlobalNS + strings.Join(components, Sep)
		} else {
			return strings.Join(components, Sep)
		}
	}
}

type nodeArguments struct {
	remapping   map[string]string
	params      map[string]string
	specialKeys map[string]string
}

func processArguments(args []string) (map[string]string, map[string]string, map[string]string, []string) {
	mapping := make(map[string]string)
	params := make(map[string]string)
	specials := make(map[string]string)
	rest := make([]string, 0)
	for _, arg := range args {
		components := strings.Split(arg, Remap)
		if len(components) == 2 {
			key := components[0]
			value := components[1]
			if strings.HasPrefix(key, "__") {
				specials[key] = value
			} else if strings.HasPrefix(key, "_") {
				params[key] = value
			} else {
				mapping[key] = value
			}
		} else {
			rest = append(rest, arg)
		}
	}
	return mapping, params, specials, rest
}
