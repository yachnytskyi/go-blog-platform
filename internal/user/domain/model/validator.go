package model

import (
	"fmt"
	"strings"
)

func (userCreate *UserCreate) UserCreateValidator() error {
	var err error
	var message string

	if len(userCreate.Name) < 4 || len(userCreate.Name) > 40 {
		message = "key: `UserCreateView.Name` error: field validation for `name` failed, `name` can be between 4 and 40 characters "
		err = fmt.Errorf(message)
	}

	if len(userCreate.Email) > 40 {
		message = message + "key: `UserCreateView.Email` error: field validation for `email` failed, `email` cannot be more that 40 characters long "
		err = fmt.Errorf(message)
	}

	if len(userCreate.Password) < 8 || len(userCreate.Password) > 40 {
		message = message + "key: `UserCreateView.Password` error: field validation for `password` failed, `password` can be between 8 and 40 characters long "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil

}
