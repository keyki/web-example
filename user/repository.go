package user

type UserRepository interface {
    ListAll() ([]*User, error)
    Create(user *User) error
    FindByUsername(username string) (*User, error)
}
