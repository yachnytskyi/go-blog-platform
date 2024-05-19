package validator

// IsError checks if the provided error is not nil.
// Parameters:
// - err: The error to be checked.
// Returns:
// - A boolean indicating whether the error is not nil.
func IsError(err error) bool {
	return err != nil
}
