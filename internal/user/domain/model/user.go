package model

import (
	"time"

	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

type Users struct {
	Users              []User
	PaginationResponse commonModel.PaginationResponse
}

type User struct {
	UserID    string
	Name      string
	Email     string
	Password  string
	Role      string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreate struct {
	Name             string
	Email            string
	Password         string
	PasswordConfirm  string
	Role             string
	Verified         bool
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserUpdate struct {
	UserID    string
	Name      string
	UpdatedAt time.Time
}

type UserLogin struct {
	Email        string
	Password     string
	AccessToken  string
	RefreshToken string
}

type UserForgottenPassword struct {
	Email string
}

type UserResetPassword struct {
	Password        string
	PasswordConfirm string
}

type EmailData struct {
	URL          string
	TemplateName string
	TemplatePath string
	FirstName    string
	Subject      string
}

func NewEmailData(url, templateName, templatePath, firstName, subject string) EmailData {
	return EmailData{
		URL:          url,
		TemplateName: templateName,
		TemplatePath: templatePath,
		FirstName:    firstName,
		Subject:      subject,
	}
}
