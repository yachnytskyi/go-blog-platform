package validator

// IsError checks if the provided error is not nil.
func IsError(err error) bool {
	return err != nil
}
