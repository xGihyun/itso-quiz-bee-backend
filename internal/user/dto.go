package user

import "time"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
	Name     string `json:"name"`
}

type GetUserResponse struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
}
