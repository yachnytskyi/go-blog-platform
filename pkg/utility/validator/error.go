package validator

// IsError checks if the provided error is not nil.
// It returns true if the error is not nil, indicating that an error has occurred.
//
// Parameters:
// - err: The error to be checked.
//
// Returns:
// - A boolean indicating whether the error is not nil.
//
// Example:
//
//	err := someFunction()
//	if validator.IsError(err) {
//	  // Handle the error
//	}
func IsError(err error) bool {
	return err != nil
}
