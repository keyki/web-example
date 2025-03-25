package database

import "gorm.io/gorm"

type TransactionError struct {
	Message string
}

func (e TransactionError) Error() string {
	return e.Message
}

type Transactional interface {
	BeginTransaction() *gorm.DB
}

type TransactionService struct {
	db *gorm.DB
}

func NewDbTransaction(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (t *TransactionService) BeginTransaction() *gorm.DB {
	return t.db.Begin()
}
