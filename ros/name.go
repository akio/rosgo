package ros

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	//Sep is a namespace separator string
	Sep = "/"
	//GlobalNS is the global namespace initial separator string
	GlobalNS = "/"
	//PrivateNS is private namespace initial separator string
	PrivateNS = "~"
)

//NameMap is a string to string map of node names and resolved names
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
	}
	return result

}

func qualifyNodeName(nodeName string) (string, string, error) {
	if nodeName == "" {
		return "", "", fmt.Errorf("Empty node name")
	}
	if nodeName[:1] == PrivateNS {
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
	}
	namespace := GlobalNS + strings.Join(components[:len(components)-1], Sep)
	return namespace, components[len(components)-1], nil

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

// Remove sequential seperater
func canonicalizeName(name string) string {
	if name == GlobalNS {
		return name
	}
	components := []string{}
	for _, word := range strings.Split(name, Sep) {
		if len(word) > 0 {
			components = append(components, word)
		}
	}
	if name[0:1] == GlobalNS {
		return GlobalNS + strings.Join(components, Sep)
	}
	return strings.Join(components, Sep)

}

//NameResolver struct definition for the NameResolver object.
//Node name is the raw name of the node in question
//namespace is the directory/qualifier of the node, by default "/"
//mapping and resolvedMapping are NameMap maps
type NameResolver struct {
	nodeName        string
	namespace       string
	mapping         NameMap
	resolvedMapping NameMap
}

//
func newNameResolver(namespace string, nodeName string, remapping NameMap) *NameResolver {
	n := new(NameResolver)
	n.nodeName = nodeName
	n.namespace = canonicalizeName(namespace)
	n.mapping = remapping
	n.resolvedMapping = make(NameMap)

	for k, v := range n.mapping {
		newKey := n.resolve(k)
		newValue := n.resolve(v)
		n.resolvedMapping[newKey] = newValue
	}

	return n
}

// Resolve a ROS name to global name
func (n *NameResolver) resolve(name string) string {
	if len(name) == 0 {
		return n.namespace
	}

	var resolvedName string
	canonName := canonicalizeName(name)
	if isGlobalName(canonName) {
		resolvedName = canonName
	} else if isPrivateName(canonName) {
		resolvedName = canonicalizeName(n.namespace + Sep + n.nodeName + Sep + canonName[1:])
	} else {
		resolvedName = canonicalizeName(n.namespace + Sep + canonName)
	}

	return resolvedName
}

// Resolve a ROS name with remapping
func (n *NameResolver) remap(name string) string {
	key := n.resolve(name)
	if value, ok := n.resolvedMapping[key]; ok {
		return value
	}
	return key

}
