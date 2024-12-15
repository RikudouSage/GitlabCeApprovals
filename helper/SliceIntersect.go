package helper

import "slices"

func SliceIntersect[TType comparable](slice1 []TType, slice2 []TType) []TType {
	result := make([]TType, 0)
	for _, value := range slice2 {
		if slices.Contains(slice1, value) {
			result = append(result, value)
		}
	}

	return result
}
