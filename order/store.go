package order

import (
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	ListAll(userId int) ([]*Order, error)
	Create(order *Order) (int, error)
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ListAll(userId int) ([]*Order, error) {
	orders := make([]*Order, 0)
	result := s.db.Preload("Products").Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

func (s *Store) Create(order *Order) (int, error) {
	log.Printf("Creating order: %v", *order)
	result := s.db.Create(&order)
	err := result.Error
	if err != nil {
		return 0, err
	}
	return order.ID, nil
}
