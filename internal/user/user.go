package user

import "time"

type Role string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

type User struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"` // WARN: Should this be here?
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
}
