package model

import (
	"time"

	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
)

type Users struct {
	Users              []User
	PaginationResponse commonModel.PaginationResponse
}

type User struct {
	domainModel.BaseEntity
	Name     string
	Email    string
	Password string
	Role     string
	Verified bool
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
	ID        string
	Name      string
	UpdatedAt time.Time
}

type UserLogin struct {
	Email    string
	Password string
}

type UserToken struct {
	AccessToken  string
	RefreshToken string
}

type UserForgottenPassword struct {
	Email       string
	ResetToken  string
	ResetExpiry time.Time
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
