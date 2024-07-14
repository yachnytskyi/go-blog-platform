package validator

import (
	"time"
)

// IsTimeNotValid checks if the provided expiration time is still valid.
//
// Parameters:
// - expiryTime: The time to check.
//
// Returns:
// - A boolean indicating whether the time is valid.
func IsTimeNotValid(expiryTime time.Time) bool {
	// Get the current time
	currentTime := time.Now()

	// Compare the provided expiry time with the current time
	if expiryTime.Before(currentTime) {
		// If the expiry time is before the current time, the time has expired
		return true
	}

	// If the expiry time is after or equal to the current time, the time is still valid
	return false
}
