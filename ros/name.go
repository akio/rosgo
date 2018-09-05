package ros

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	Sep       = "/"
	GlobalNS  = "/"
	PrivateNS = "~"
	Remap     = ":="
)

type NameMap map[string]string

func getNamespace(name string) string {
	if len(name) == 0 {
		return GlobalNS
	} else if name[len(name)-1] == '/' {
		name = name[:len(name)-1]
	}
	result := name[:strings.LastIndex(name, Sep)+1]
	if len(result) == 0 {
		return Sep
	} else {
		return result
	}
}

func qualifyNodeName(nodeName string) (string, string, error) {
	if nodeName == "" {
		return "", "", fmt.Errorf("Empty node name")
	}
	if nodeName[0] == "~" {
		return "", "", fmt.Errorf("Node name should not contain '~'")
	}
	canonName := canonicalizeName(nodeName)

	var components []string
	for _, c := range strings.Split(canonName, Sep) {
		if len(c) > 0 {
			components = append(components, c)
		}
	}
	if len(components) == 1 {
		return GlobalNS, components[0], nil
	} else {
		namespace := GlobalNS + strings.Join(components[:-1], Sep)
		return namespace, components[-1], nil
	}
}

func resolveName(name string, namespace string, mappings NameMap) string {
	var resolvedName string

	if len(name) == 0 {
		return getNamespace(namespace)
	}

	canonName := canonicalizeName(name)
	if isGlobalName(canonName) {
		resolvedName = canonName
	} else if isPrivateName(canonName) {
		resolvedName = canonicalizeName(namespace + Sep + canonName[1:])
	} else {
		resolvedName = getNamespace(namespace) + canonName
	}

	if mappings != nil {
		if remappedName, ok := mappings[resolvedName]; ok {
			return remappedName
		} else {
			return resolvedName
		}
	} else {
		return resolvedName
	}
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return true
	}
	if name == "/" || name == "~" {
		return true
	}
	if matched, _ := regexp.MatchString("^[~/]?([a-zA-Z]\\w*/)*[a-zA-Z]\\w*$", name); !matched {
		return false
	}
	return true
}

func isValidNamespace(name string) bool {
	if len(name) == 0 {
		return false
	}
	if matched, _ := regexp.MatchString("^/([a-zA-Z]\\w*/)*$", name); !matched {
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

// Remove sequential seperater
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

func processArguments(args []string) (NameMap, NameMap, NameMap, []string) {
	mapping := make(NameMap)
	params := make(NameMap)
	specials := make(NameMap)
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

type NameResolver struct {
	nodeName        string
	namespace       string
	mapping         NameMap
	resolvedMapping NameMap
}

func newNameResolver(nodeName string, remapping NameMap) *NameResolver {
	n := new(NameResolver)

	n.namespace = canonicalizeName(namespace)
	n.mapping = remapping
	n.resolvedMapping = make(NameMap)

	for k, v := range n.mapping {
		newKey := resolveName(k, namespace, nil)
		newValue := resolveName(v, namespace, nil)
		n.resolvedMapping[newKey] = newValue
	}

	for k, v := range n.resolvedMapping {
		fmt.Printf("%s: %s\n", k, v)
	}

	return n
}

func (n *NameResolver) resolve(name string) string {
	return resolveName(name, n.namespace, n.resolvedMapping)
}

func (n *NameResolver) remap(name string) string {
	r := resolveName(name, n.namespace, n.resolvedMapping)
	if remapped, ok := n.mapping[r]; ok {
		return resolveName(remapped, n.namespace, n.resolvedMapping)
	} else {
		return r
	}
}
