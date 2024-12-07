package model

type UserTokenPayload struct {
	UserID string
	Role   string
}

func NewUserTokenPayload(userID, role string) UserTokenPayload {
	return UserTokenPayload{
		UserID: userID,
		Role:   role,
	}
}
