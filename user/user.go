package user

import (
	"fmt"
	"web-example/types"
	"web-example/validator"
)

type User struct {
	ID       int        `gorm:"primaryKey,autoIncrement"`
	UserName string     `gorm:"unique;not null;type:varchar(255)"`
	Password string     `gorm:"not null;type:varchar(255)"`
	Role     types.Role `gorm:"not null;type:varchar(100)"`
}

type Request struct {
	UserName string     `json:"user_name" validate:"required;min=4,max=255"`
	Password string     `json:"password" validate:"required;min=1,max=255"`
	Role     types.Role `json:"role" validate:"required;oneOfRole=ADMIN,USER"`
}

type Response struct {
	UserName string     `json:"user_name"`
	Role     types.Role `json:"role"`
}

func (u *User) ToResponse() *Response {
	return &Response{
		UserName: u.UserName,
		Role:     u.Role,
	}
}

func (u User) String() string {
	return fmt.Sprintf("username: %s role: %s", u.UserName, u.Role)
}

func (User) TableName() string {
	return "users"
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
