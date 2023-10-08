package validator

func IsSliceEmpty[T any](s []T) bool {
	return len(s) == 0
}

func IsSliceNotEmpty[T any](s []T) bool {
	return len(s) > 0
}
