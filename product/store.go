package product

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
	"web-example/types"
)

type Repository interface {
	ListAll() ([]*Product, error)
	Create(product *Product) error
	FindByName(name string) (*Product, error)
	FindAllByName(names []string) ([]*Product, error)
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	store := &Store{db}
	store.init()
	return store
}

func (s *Store) ListAll() ([]*Product, error) {
	log.Printf("Listing products")
	var products []*Product
	result := s.db.Find(&products)
	if result.Error != nil {
		return products, result.Error
	}
	return products, nil
}

func (s *Store) Create(product *Product) error {
	log.Printf("Creating product: %v", *product)
	result := s.db.Create(&product)
	err := result.Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				log.Printf("Product '%s' already exists", product.Name)
			} else {
				log.Println(pgErr)
			}
		} else {
			log.Printf("Cannot create product: %v", err)
		}
		return err
	}
	return nil
}

func (s *Store) FindByName(name string) (*Product, error) {
	log.Printf("Finding user by username: %v", name)
	var product Product
	result := s.db.Where("name = ?", name).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Printf("Found product: %v", product)
	return &product, nil
}

func (s *Store) FindAllByName(names []string) ([]*Product, error) {
	log.Printf("Finding %d products by names: %v", len(names), names)
	products := make([]*Product, 0)
	result := s.db.Where("name IN (?)", names).Find(&products)
	if result.Error != nil {
		return products, result.Error
	}
	log.Printf("Found %d products", len(products))
	return products, nil
}

func (s *Store) init() {
	s.Create(&Product{
		Name:        "pen",
		Description: "Blue pen",
		Price:       100,
		Currency:    types.HUF,
		Quantity:    50,
	})
	s.Create(&Product{
		Name:        "book",
		Description: "Harry Potter 1",
		Price:       5500,
		Currency:    types.HUF,
		Quantity:    12,
	})
}
