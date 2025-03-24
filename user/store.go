package user

import (
    "gorm.io/gorm"
    "log"
)

type UserStore struct {
    db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
    return &UserStore{db}
}

func (s *UserStore) ListAll() ([]*User, error) {
    log.Printf("Listing users")
    var users []*User
    result := s.db.Find(&users)
    if result.Error != nil {
        return users, result.Error
    }
    return users, nil
}

func (s *UserStore) Create(user *User) error {
    log.Printf("Creating user: %v", *user)
    result := s.db.Create(&user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
