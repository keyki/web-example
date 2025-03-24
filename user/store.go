package user

import (
    "errors"
    "github.com/jackc/pgx/v5/pgconn"
    "gorm.io/gorm"
    "log"
    "web-example/types"
    "web-example/util"
)

type Store struct {
    db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
    store := &Store{db}
    store.init()
    return store
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

func (s *Store) init() {
    err := s.Create(&User{
        UserName: "admin",
        Password: util.HashPassword("admin"),
        Role:     types.ADMIN,
    })
    if err == nil {
        return
    }
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        if pgErr.Code == "23505" {
            log.Printf("Admin user already exists")
        } else {
            log.Println(err)
        }
    } else {
        log.Printf("Cannot create admin user: %v", err)
    }
}
