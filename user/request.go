package user

import (
    "web-example/types"
    "web-example/validator"
)

type UserRequest struct {
    UserName string     `json:"user_name" validate:"required;min=4,max=255"`
    Password string     `json:"password" validate:"required;min=1,max=255"`
    Role     types.Role `json:"role" validate:"required;oneOfRole=ADMIN,USER"`
}

func (u *UserRequest) Validate() error {
    return validator.Validate(u)
}

func (u *UserRequest) ToUser() *User {
    return &User{
        UserName: u.UserName,
        Password: u.Password,
        Role:     u.Role,
    }
}
