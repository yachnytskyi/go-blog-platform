package model

import (
	"time"

	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
)

// Users represents a collection of users with pagination information.
type Users struct {
	Users              []User                         // Slice of User entities
	PaginationResponse commonModel.PaginationResponse // Pagination metadata
}

// User represents a user entity with basic details.
type User struct {
	domainModel.BaseEntity        // Embedded BaseEntity for common fields
	Name                   string // User's full name
	Email                  string // User's email address
	Password               string // User's hashed password
	Role                   string // User's role (e.g., admin, user)
	Verified               bool   // Flag indicating if the user's email is verified
}

// UserCreate represents data for creating a new user.
type UserCreate struct {
	Name             string    // User's full name
	Email            string    // User's email address
	Password         string    // User's plaintext password (to be hashed)
	PasswordConfirm  string    // Confirmation of the plaintext password
	Role             string    // User's role (e.g., admin, user)
	Verified         bool      // Flag indicating if the user's email is verified
	VerificationCode string    // Verification code for email verification
	CreatedAt        time.Time // Timestamp of user creation
	UpdatedAt        time.Time // Timestamp of last update
}

// UserUpdate represents data for updating a user.
type UserUpdate struct {
	ID        string    // ID of the user to be updated
	Name      string    // Updated user's full name
	UpdatedAt time.Time // Timestamp of update
}

// UserLogin represents data for user login credentials.
type UserLogin struct {
	Email    string // User's email address for login
	Password string // User's password for login
}

// UserToken represents data for access and refresh tokens.
type UserToken struct {
	AccessToken  string // JWT access token for user authentication
	RefreshToken string // Refresh token for generating new access tokens
}

// UserForgottenPassword represents data for initiating forgotten password flow.
type UserForgottenPassword struct {
	Email       string    // User's email address requesting password reset
	ResetToken  string    // Unique token for resetting the password
	ResetExpiry time.Time // Expiry time for the password reset token
}

// UserResetPassword represents data for resetting a user's password.
type UserResetPassword struct {
	ResetToken      string // Unique token for resetting the password
	Password        string // New password for the user
	PasswordConfirm string // Confirmation of the new password
}

// UserResetToken represents data for storing reset token expiry information.
type UserResetToken struct {
	ResetExpiry time.Time // Expiry time for the password reset token
}

// EmailData represents data required for sending an email.
type EmailData struct {
	URL          string // URL to include in the email
	TemplateName string // Name of the email template
	TemplatePath string // Path to the email template file
	FirstName    string // First name of the recipient
	Subject      string // Subject of the email
}

// NewUsers creates a new instance of Users with provided users and pagination response.
func NewUsers(users []User, paginationResponse commonModel.PaginationResponse) Users {
	return Users{
		Users:              users,
		PaginationResponse: paginationResponse,
	}
}

// NewUser creates a new instance of User with provided details.
func NewUser(id string, createdAt, updatedAt time.Time, name, email, role string) User {
	return User{
		BaseEntity: domainModel.NewBaseEntity(id, createdAt, updatedAt),
		Name:       name,
		Email:      email,
		Role:       role,
	}
}

// NewUserCreate creates a new instance of UserCreate with provided details.
func NewUserCreate(name, email, password, passwordConfirm string) UserCreate {
	return UserCreate{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

// NewUserUpdate creates a new instance of UserUpdate with provided details.
func NewUserUpdate(id, name string) UserUpdate {
	return UserUpdate{
		ID:   id,
		Name: name,
	}
}

// NewUserLogin creates a new instance of UserLogin with provided email and password.
func NewUserLogin(email, password string) UserLogin {
	return UserLogin{
		Email:    email,
		Password: password,
	}
}

// NewUserToken creates a new instance of UserToken with provided access token and refresh token.
func NewUserToken(accessToken, refreshToken string) UserToken {
	return UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// NewUserForgottenPassword creates a new instance of UserForgottenPassword with provided details.
func NewUserForgottenPassword(email string) UserForgottenPassword {
	return UserForgottenPassword{
		Email: email,
	}
}

// NewUserResetPassword creates a new instance of UserResetPassword with provided details.
func NewUserResetPassword(resetToken, password, passwordConfirm string) UserResetPassword {
	return UserResetPassword{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

// NewUserResetToken creates a new instance of UserResetToken with provided reset token expiry.
func NewUserResetToken(resetExpiry time.Time) UserResetToken {
	return UserResetToken{
		ResetExpiry: resetExpiry,
	}
}

// NewEmailData creates a new instance of EmailData with provided details.
func NewEmailData(url, templateName, templatePath, firstName, subject string) EmailData {
	return EmailData{
		URL:          url,
		TemplateName: templateName,
		TemplatePath: templatePath,
		FirstName:    firstName,
		Subject:      subject,
	}
}
