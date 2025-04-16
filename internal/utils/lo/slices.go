package lo

// Uniq creates new array containing only unique elements. Order is not defined
func Uniq[T comparable](elems []T) []T {
	var uniq = make(map[T]struct{}, len(elems))
	for _, e := range elems {
		uniq[e] = struct{}{}
	}

	var result []T
	for e := range uniq {
		result = append(result, e)
	}

	return result
}
