package order

import (
	"context"
	"gorm.io/gorm"
	"web-example/log"
)

type Repository interface {
	ListAll(ctx context.Context, userId int) ([]*Order, error)
	Create(ctx context.Context, order *Order, tx *gorm.DB) (int, error)
	Find(ctx context.Context, orderId int, userId int) (*Order, error)
	Delete(ctx context.Context, order *Order) error
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

func (s *Store) Find(ctx context.Context, orderId int, userId int) (*Order, error) {
	order := &Order{}
	logger := log.Logger(ctx)
	logger.Infof("Finding order by id: %d", orderId)
	result := s.db.Where("id = ? AND user_id = ?", orderId, userId).Find(order)
	if result.Error != nil {
		return order, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	logger.Infof("Found order: %v", order)
	return order, nil
}

func (s *Store) Delete(ctx context.Context, order *Order) error {
	logger := log.Logger(ctx)
	logger.Infof("Deleting order: %d", order.ID)
	result := s.db.Delete(order)
	if result.Error != nil {
		return result.Error
	}
	logger.Infof("Deleted order: %d", order.ID)
	return nil
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
