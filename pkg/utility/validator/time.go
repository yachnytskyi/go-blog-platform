package validator

import (
	"time"
)

// IsTimeNotValid checks if the provided expiration time is still valid.
func IsTimeNotValid(expiryTime time.Time) bool {
	currentTime := time.Now()
	if expiryTime.Before(currentTime) {
		return true
	}

	return false
}
