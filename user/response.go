package user

import "web-example/types"

type UserResponse struct {
    UserName string     `json:"user_name"`
    Role     types.Role `json:"role"`
}
