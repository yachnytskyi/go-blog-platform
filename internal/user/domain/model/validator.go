package model

import (
	"fmt"
	"strings"
)

func (userCreate *UserCreate) UserCreateValidator() error {
	var err error
	var message string

	if len(userCreate.Name) < 4 || len(userCreate.Name) > 40 {
		message = "key: `UserCreate.Name` error: field validation for `name` failed, `name` can be between 4 and 40 characters "
		err = fmt.Errorf(message)
	}

	if len(userCreate.Email) > 40 {
		message = message + "key: `UserCreate.Email` error: field validation for `email` failed, `email` cannot be more that 40 characters long "
		err = fmt.Errorf(message)
	}

	if len(userCreate.Password) < 8 || len(userCreate.Password) > 40 {
		message = message + "key: `UserCreate.Password` error: field validation for `password` failed, `password` can be between 8 and 40 characters long "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil

}

func (userUpdate *UserUpdate) UserUpdateValidator() error {
	var message string
	var err error

	if len(userUpdate.Name) < 4 || len(userUpdate.Name) > 40 {
		message = "key: `UserUpdate.Name` error: field validation for `name` failed, `name` can be between 4 and 40 characters "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userLogin *UserLogin) UserLoginValidator() error {
	var message string
	var err error

	if len(userLogin.Email) > 40 {
		message = "key: `UserLogin.Email` error: field validation for `email` failed, `email` cannot be more that 40 characters long "
		err = fmt.Errorf(message)
	}

	if len(userLogin.Password) < 4 || len(userLogin.Email) > 40 {
		message = message + "key: `UserLogin.Password` error: field validation for `password` failed, `password` can be between 4 and 40 characters "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userForgottenPassword *UserForgottenPassword) UserForgottenPasswordValidator() error {
	var err error
	var message string

	if len(userForgottenPassword.Email) > 40 {
		message = "key: `UserForgottenPassword.Email` error: field validation for `email` failed, `email` cannot be more that 40 characters long "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil

}

func (userResetPassword *UserResetPassword) UserResetPasswordValidator() error {
	var err error
	var message string

	if len(userResetPassword.Password) < 4 || len(userResetPassword.Password) > 40 {
		message = "key: `UserResetPassword.Password` error: field validation for `password` failed, `password` can be between 8 and 40 characters "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}
