package validator

func IsSliceEmpty[T any](s []T) bool {
	if len(s) == 0 {
		return true
	}
	return false
}

func IsSliceNotEmpty[T any](s []T) bool {
	if len(s) == 0 {
		return false
	}
	return true
}
