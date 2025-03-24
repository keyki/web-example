package user

import (
    "fmt"
    "web-example/types"
)

type User struct {
    id       int        `gorm:"primaryKey,autoIncrement"`
    UserName string     `gorm:"unique;not null;type:varchar(255)"`
    Password string     `gorm:"not null;type:varchar(255)"`
    Role     types.Role `gorm:"not null;type:varchar(100)"`
}

func (u *User) ToReponse() *Response {
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
