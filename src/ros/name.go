package ros


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


