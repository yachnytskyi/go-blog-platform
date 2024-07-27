package validator

func IsSliceContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

func IsSliceNotContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return false
		}
	}

	return true
}
