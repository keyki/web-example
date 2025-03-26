package user

import (
	"web-example/types"
	"web-example/validator"
)

type Request struct {
	UserName string     `json:"user_name" validate:"required;min=4,max=255"`
	Password string     `json:"password" validate:"required;min=1,max=255"`
	Role     types.Role `json:"role" validate:"required;oneOfRole=ADMIN,USER"`
}

type Response struct {
	UserName string     `json:"user_name"`
	Role     types.Role `json:"role"`
}

func (u *Request) Validate() error {
	return validator.Validate(u)
}

func (u *Request) ToUser() *User {
	return &User{
		UserName: u.UserName,
		Password: u.Password,
		Role:     u.Role,
	}
}
