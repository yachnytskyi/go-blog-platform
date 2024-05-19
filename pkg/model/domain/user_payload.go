package domain

// UserTokenPayload represents the payload of a user token,
// which includes the user ID and the user's role.
type UserTokenPayload struct {
	UserID string // The unique identifier of the user.
	Role   string // The role of the user (e.g., "admin", "user").
}

// NewUserTokenPayload creates a new UserTokenPayload with the provided user ID and role.
// Parameters:
// - userID: A string representing the unique identifier of the user.
// - role: A string representing the role of the user.
// Returns:
// - A UserTokenPayload struct populated with the given user ID and role.
func NewUserTokenPayload(userID, role string) UserTokenPayload {
	return UserTokenPayload{
		UserID: userID,
		Role:   role,
	}
}
