package user

import "web-example/types"

type Response struct {
    UserName string     `json:"user_name"`
    Role     types.Role `json:"role"`
}
