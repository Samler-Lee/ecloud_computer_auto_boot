package util

func InArray[T comparable](slice []T, search T) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}

	return false
}
