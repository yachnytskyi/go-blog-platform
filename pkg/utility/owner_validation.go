package utility

import (
	"bytes"
	"errors"
)

func IsUserOwner(currentUserID string, userID string) error {
	convertedCurrentUserID := []byte(currentUserID)
	convertedUserID := []byte(userID)
	comparisson := bytes.Compare(convertedCurrentUserID, convertedUserID)

	if comparisson != 0 {
		return errors.New("sorry, but you do not have permissions to do that")
	}
	return nil
}
