package user

import "time"

type Role string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
	Name     string `json:"name"`
}

type UserResponse struct {
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username"`
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatarUrl"`
}
