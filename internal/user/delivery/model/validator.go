package model

import (
	"fmt"
	"net/mail"
	"strings"
)

func (userCreateView *UserCreateView) UserCreateViewValidator() error {
	var message string
	var err error

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
		message = message + "key: `UserCreateView.PasswordConfirm` error: field validation for `password_confirm` failed, `password_confirm` cannot be empty "
		err = fmt.Errorf(message)

	}

	if userCreateView.Email != "" {
		_, ok := mail.ParseAddress(userCreateView.Email)
		if ok != nil {
			message = message + "key: `UserCreateView.Email` error: field validation for `email` failed, invalid email address "
			err = fmt.Errorf(message)

		}
	}

	if userCreateView.Password != "" && userCreateView.PasswordConfirm != "" && userCreateView.Password != userCreateView.PasswordConfirm {
		message = message + "key: `UserCreateView.PasswordConfirm` error: field validation for `password_confirm` failed, passwords do not match "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userUpdateView *UserUpdateView) UserUpdateViewValidator() error {
	var message string
	var err error

	if userUpdateView.Name == "" {
		message = "key: `UserUpdateView.Name` error: field validation for `name` failed, `name` cannot be empty "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userLoginView *UserLoginView) UserSignInViewValidator() error {
	var message string
	var err error

	if userLoginView.Email == "" {
		message = "key: `UserLogInView.Email` error: field validation for `email` failed, `email` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userLoginView.Password == "" {
		message = message + "key: `UserLogInView.Password` error: field validation for `password` failed, `password` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userLoginView.Email != "" {
		_, ok := mail.ParseAddress(userLoginView.Email)
		if ok != nil {
			message = message + "key: `UserCreateView.Email` error: field validation for `email` failed, invalid email address "
			err = fmt.Errorf(message)

		}
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userForgottenPasswordView *UserForgottenPasswordView) UserForgottenPasswordViewValidator() error {
	var message string
	var err error

	if userForgottenPasswordView.Email == "" {
		message = "key: `UserForgottenPasswordView.Email` error: field validation for `email` failed, `email` cannot be empty "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userResetPasswordView *UserResetPasswordView) UserResetPasswordViewValidator() error {
	var message string
	var err error

	if userResetPasswordView.Password == "" {
		message = "key: `UserResetPasswordView.Password` error: field validation for `password` failed, `password` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userResetPasswordView.PasswordConfirm == "" {
		message = "key: `UserResetPasswordView.PasswordConfirm` error: field validation for `password_confirm` failed, `password_confirm` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userResetPasswordView.Password != "" && userResetPasswordView.PasswordConfirm != "" && userResetPasswordView.Password != userResetPasswordView.PasswordConfirm {
		message = message + "key: `UserResetPasswordView.PasswordConfirm` error: field validation for `password_confirm` failed, passwords do not match "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}
