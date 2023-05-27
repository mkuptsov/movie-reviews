package users

import "time"

const (
	AdminRole  = "admin"
	EditorRole = "editor"
	UserRole   = "user"
)

type User struct {
	Id        int        `json:"user_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Bio       *string    `json:"bio,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

type UserWithPassword struct {
	*User
	PasswordHash string
}
