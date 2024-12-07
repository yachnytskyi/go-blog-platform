package validator

import (
	"time"
)

// IsTimeNotValid checks if the provided expiration time is still valid.
func IsTimeNotValid(expiryTime time.Time) bool {
	return expiryTime.Before(time.Now())
}
