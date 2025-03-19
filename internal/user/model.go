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
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url"`
}
