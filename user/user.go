package user

import (
    "web-example/types"
    "web-example/validator"
)

type User struct {
    id       int        `gorm:"primaryKey,autoIncrement"`
    UserName string     `json:"user_name" gorm:"unique;not null;type:varchar(255)" validate:"required;min=4,max=255"`
    Role     types.Role `json:"role" gorm:"not null;type:varchar(100)" validate:"required;oneOfRole=ADMIN,USER"`
}

func (User) TableName() string {
    return "users"
}

func (u *User) Validate() error {
    return validator.Validate(u)
}
