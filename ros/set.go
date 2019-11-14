package ros

func contains(array []string, key string) bool {
	for _, item := range array {
		if item == key {
			return true
		}
	}
	return false
}

func unique(array []string) []string {
	set := map[string]bool{}
	for _, item := range array {
		set[item] = true
	}
	var result []string
	for k := range set {
		result = append(result, k)
	}
	return result
}

func setUnion(lhs []string, rhs []string) []string {
	set := map[string]bool{}
	for _, item := range lhs {
		set[item] = true
	}
	for _, item := range rhs {
		set[item] = true
	}
	var result []string
	for k := range set {
		result = append(result, k)
	}
	return result
}

func setDifference(lhs []string, rhs []string) []string {
	left := map[string]bool{}
	for _, item := range lhs {
		left[item] = true
	}
	right := map[string]bool{}
	for _, item := range rhs {
		right[item] = true
	}
	for k := range right {
		delete(left, k)
	}
	var result []string
	for k := range left {
		result = append(result, k)
	}
	return result
}
