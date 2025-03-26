package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"time"
	"web-example/cache"
	"web-example/log"
	"web-example/types"
	"web-example/util"
)

var userCache = cache.NewCache(30 * time.Second)

type Repository interface {
	ListAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	store := &Store{db}
	store.init()
	return store
}

func (s *Store) ListAll(ctx context.Context) ([]*User, error) {
	log.Logger(ctx).Info("Listing users")
	users := make([]*User, 0)
	result := s.db.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	log.Logger(ctx).Infof("Found %d users", len(users))
	return users, nil
}

func (s *Store) Create(ctx context.Context, user *User) error {
	log.Logger(ctx).Infof("Creating user: %v", *user)
	result := s.db.Create(&user)
	err := result.Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				log.Logger(ctx).Infof("User '%s' already exists", user.UserName)
			} else {
				log.Logger(ctx).Info(pgErr)
			}
		} else {
			log.Logger(ctx).Infof("Cannot create user: %v", err)
		}
		return err
	}
	return nil
}

func (s *Store) FindByUsername(ctx context.Context, username string) (*User, error) {
	log.Logger(ctx).Infof("Finding user by username: %v", username)
	var user User
	cachedUser, found := userCache.Get(username)
	if found {
		user, _ = cachedUser.(User)
		return &user, nil
	}
	result := s.db.Where("user_name = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userCache.Set(username, user, 10*time.Second)
	log.Logger(ctx).Infof("Found user: %v", user)
	return &user, nil
}

func (s *Store) init() {
	s.Create(context.Background(), &User{
		UserName: "admin",
		Password: util.HashPassword("admin"),
		Role:     types.ADMIN,
	})
	s.Create(context.Background(), &User{
		UserName: "alma",
		Password: util.HashPassword("alma"),
		Role:     types.USER,
	})
}
