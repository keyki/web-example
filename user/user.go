package user

import (
    "fmt"
    "web-example/types"
    "web-example/validator"
)

type User struct {
    id       int        `gorm:"primaryKey,autoIncrement"`
    UserName string     `json:"user_name" gorm:"unique;not null;type:varchar(255)" validate:"required;min=4,max=255"`
    Password string     `json:"password" gorm:"not null;type:varchar(255)" validate:"required;min=1,max=255"`
    Role     types.Role `json:"role" gorm:"not null;type:varchar(100)" validate:"required;oneOfRole=ADMIN,USER"`
}

type UserResponse struct {
    UserName string     `json:"user_name"`
    Role     types.Role `json:"role"`
}

func (u *User) ToReponse() *UserResponse {
    return &UserResponse{
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

func (u *User) Validate() error {
    return validator.Validate(u)
}
