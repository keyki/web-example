package order

import (
	"context"
	"gorm.io/gorm"
	"web-example/log"
)

type Repository interface {
	ListAll(ctx context.Context, userId int) ([]*Order, error)
	Create(ctx context.Context, order *Order, tx *gorm.DB) (int, error)
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ListAll(_ context.Context, userId int) ([]*Order, error) {
	orders := make([]*Order, 0)
	result := s.db.Preload("Products.Product").Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

func (s *Store) Create(ctx context.Context, order *Order, tx *gorm.DB) (int, error) {
	if tx != nil {
		return createInTransaction(ctx, order, tx)
	}
	return createInTransaction(ctx, order, s.db)
}

func createInTransaction(ctx context.Context, order *Order, tx *gorm.DB) (int, error) {
	result := tx.Create(&order)
	err := result.Error
	if err != nil {
		log.Logger(ctx).Infof("Error creating order: %v", err)
		return 0, err
	}
	return order.ID, nil
}
