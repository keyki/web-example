package user

var (
    ADMIN = Role("ADMIN")
    USER  = Role("USER")
)

type Role string

type User struct {
    ID       int    `json:"id" gorm:"primaryKey,autoIncrement"`
    UserName string `json:"user_name" gorm:"unique;not null;type:varchar(255)"`
    Role     Role   `json:"role" gorm:"not null;type:varchar(100)"`
}

func (User) TableName() string {
    return "users"
}
