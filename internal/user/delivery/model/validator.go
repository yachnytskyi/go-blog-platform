package model

import (
	"fmt"
	"net/mail"
)

func (userCreateView *UserCreateView) UserCreateViewValidator() error {
	var err error
	var message string

	if userCreateView.Name == "" {
		message = "key: `UserCreateView.Name` error: field validation for `name` failed, `name` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userCreateView.Email == "" {
		message = message + "key: `UserCreateView.Email` error: field validation for `email` failed, `email` cannot be empty "
		err = fmt.Errorf(message)

	}

	if userCreateView.Password == "" {
		message = message + "key: `UserCreateView.Password` error: field validation for `password` failed, `password` cannot be empty "
		err = fmt.Errorf(message)

	}

	if userCreateView.PasswordConfirm == "" {
		message = message + "key: `UserCreateView.PasswordConfirm` error: field validation for `password_confirm` failed, `password_confirm` cannot be empty"
		err = fmt.Errorf(message)

	}

	if userCreateView.Email != "" {
		_, ok := mail.ParseAddress(userCreateView.Email)
		if ok != nil {
			message = message + "key: `UserCreateView.Email` error: field validation for `email` failed, invalid email address"
			err = fmt.Errorf(message)

		}
	}

	if userCreateView.Password != "" && userCreateView.PasswordConfirm != "" && userCreateView.Password != userCreateView.PasswordConfirm {
		message = message + "key: `UserCreateView.PasswordConfirm` error: field validation for `password_confirm` failed, passwords do not match"
		err = fmt.Errorf(message)

	}

	return err
}
