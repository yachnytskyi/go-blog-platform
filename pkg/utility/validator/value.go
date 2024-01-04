package validator

func IsValueEmpty[T comparable](data T) bool {
	var t T
	return data == t
}

func IsValueNotEmpty[T comparable](data T) bool {
	var t T
	return data != t
}
