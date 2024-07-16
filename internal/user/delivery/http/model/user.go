package model

import (
	"time"

	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

// TokenView represents a view for a token response.
type TokenView struct {
	Token string `json:"token"`
}

// UsersView represents a view for multiple users with pagination information.
type UsersView struct {
	UsersView              []UserView                       `json:"users"`
	HTTPPaginationResponse httpModel.HTTPPaginationResponse `json:"pagination_response"`
}

// UserView represents a view for a single user.
type UserView struct {
	httpModel.BaseEntity        // Embedding BaseEntity for common fields
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
}

// UserCreateView represents a view for creating a new user.
type UserCreateView struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// UserUpdateView represents a view for updating a user.
type UserUpdateView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserLoginView represents a view for user login credentials.
type UserLoginView struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserTokenView represents a view for access and refresh tokens.
type UserTokenView struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserForgottenPasswordView represents a view for initiating forgotten password flow.
type UserForgottenPasswordView struct {
	Email string `json:"email"`
}

// UserResetPasswordView represents a view for resetting user password.
type UserResetPasswordView struct {
	ResetToken      string `json:"reset_token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// UserWelcomeMessageView represents a view for a welcome message notification.
type UserWelcomeMessageView struct {
	Notification string `json:"notification"`
}

// NewWelcomeMessageView creates a new instance of UserWelcomeMessageView with a notification message.
func NewWelcomeMessageView(notification string) UserWelcomeMessageView {
	return UserWelcomeMessageView{
		Notification: notification,
	}
}

// NewUsersView creates a new instance of UsersView with provided user views and pagination response.
func NewUsersView(users []UserView, paginationResponse httpModel.HTTPPaginationResponse) UsersView {
	return UsersView{
		UsersView:              users,
		HTTPPaginationResponse: paginationResponse,
	}
}

// NewUserView creates a new instance of UserView with the provided attributes.
// Returns:
// - UserView: The initialized UserView struct with the specified attributes.
func NewUserView(id string, createdAt, updatedAt time.Time, name, email, role string) UserView {
	return UserView{
		BaseEntity: httpModel.NewBaseEntity(id, createdAt, updatedAt),
		Name:       name,
		Email:      email,
		Role:       role,
	}
}

// NewUserCreateView creates a new instance of UserCreateView with provided name, email, password, and password confirmation.
func NewUserCreateView(name, email, password, passwordConfirm string) UserCreateView {
	return UserCreateView{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

// NewUserUpdateView creates a new instance of UserUpdateView with provided ID and name.
func NewUserUpdateView(id, name string) UserUpdateView {
	return UserUpdateView{
		ID:   id,
		Name: name,
	}
}

// NewUserLoginView creates a new instance of UserLoginView with provided email and password.
func NewUserLoginView(email, password string) UserLoginView {
	return UserLoginView{
		Email:    email,
		Password: password,
	}
}

// NewUserForgottenPasswordView creates a new instance of UserForgottenPasswordView with provided email.
func NewUserForgottenPasswordView(email string) UserForgottenPasswordView {
	return UserForgottenPasswordView{
		Email: email,
	}
}

// NewUserResetPasswordView creates a new instance of UserResetPasswordView with provided reset token, password, and password confirmation.
func NewUserResetPasswordView(resetToken, password, passwordConfirm string) UserResetPasswordView {
	return UserResetPasswordView{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

// NewUserTokenView creates a new instance of UserTokenView with provided access token and refresh token.
func NewUserTokenView(accessToken, refreshToken string) UserTokenView {
	return UserTokenView{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
