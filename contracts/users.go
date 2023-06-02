package contracts

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Bio       *string    `json:"bio,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
type GetUserByIdRequest struct {
	UserId int `param:"userId" validate:"nonzero"`
}
type GetUserByUserNameRequest struct {
	UserName string `param:"userName" validate:"nonzero"`
}
type UpdateUserRequest struct {
	UserId int     `param:"userId" validate:"nonzero"`
	Bio    *string `json:"bio"`
}
type DeleteUserRequest struct {
	UserId int `param:"userId" validate:"nonzero"`
}
type UpdateUserRoleRequest struct {
	UserId int    `param:"userId" validate:"nonzero"`
	Role   string `param:"role" validate:"role"`
}
