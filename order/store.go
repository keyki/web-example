package order

import (
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	ListAll(userId int) ([]*Order, error)
	Create(order *Order, tx *gorm.DB) (int, error)
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ListAll(userId int) ([]*Order, error) {
	orders := make([]*Order, 0)
	result := s.db.Preload("Products.Product").Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

func (s *Store) Create(order *Order, tx *gorm.DB) (int, error) {
	if tx != nil {
		return createInTransaction(order, tx)
	}
	return createInTransaction(order, s.db)
}

func createInTransaction(order *Order, tx *gorm.DB) (int, error) {
	result := tx.Create(&order)
	err := result.Error
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return 0, err
	}
	return order.ID, nil
}
