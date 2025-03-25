package order

import "gorm.io/gorm"

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ListAll(userId int) ([]*Order, error) {
	orders := make([]*Order, 0)
	result := s.db.Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return orders, result.Error
	}
	return orders, nil
}

func (s *Store) Create(order *Order) error {
	return nil
}
