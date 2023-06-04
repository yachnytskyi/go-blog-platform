package model

import "fmt"

func (userCreate *UserCreate) UserCreateValidator() error {
	var err error
	var message string

	if len(userCreate.Name) < 4 || len(userCreate.Name) > 40 {
		message = "something"
		err = fmt.Errorf(message)
	}

	return err

}
