package user

import (
    "gorm.io/gorm"
    "log"
)

type Store struct {
    db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
    return &Store{db}
}

func (s *Store) ListAll() ([]*User, error) {
    log.Printf("Listing users")
    var users []*User
    result := s.db.Find(&users)
    if result.Error != nil {
        return users, result.Error
    }
    return users, nil
}

func (s *Store) Create(user *User) error {
    log.Printf("Creating user: %v", *user)
    result := s.db.Create(&user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (s *Store) FindByUsername(username string) (*User, error) {
    log.Printf("Finding user by username: %v", username)
    var user User
    result := s.db.Where("user_name = ?", username).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    log.Printf("Found user: %v", user)
    return &user, nil
}
