package product

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"web-example/database"
	"web-example/log"
	"web-example/types"
)

type Repository interface {
	ListAll(ctx context.Context) ([]*Product, error)
	Create(ctx context.Context, product *Product) error
	FindByName(ctx context.Context, name string) (*Product, error)
	FindAllByName(ctx context.Context, names []string) ([]*Product, error)
	FindAllByIds(ctx context.Context, ids []int) ([]*Product, error)
	UpdateQuantity(ctx context.Context, product *Product, tx *gorm.DB) error
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	store := &Store{db}
	store.init()
	return store
}

func (s *Store) ListAll(ctx context.Context) ([]*Product, error) {
	log.Logger(ctx).Info("Listing products")
	var products []*Product
	result := s.db.Find(&products)
	if result.Error != nil {
		return products, result.Error
	}
	return products, nil
}

func (s *Store) Create(ctx context.Context, product *Product) error {
	log.Logger(ctx).Infof("Creating product: %v", *product)
	result := s.db.Create(&product)
	err := result.Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				log.Logger(ctx).Infof("Product '%s' already exists", product.Name)
			} else {
				log.Logger(ctx).Info(pgErr)
			}
		} else {
			log.Logger(ctx).Infof("Cannot create product: %v", err)
		}
		return err
	}
	return nil
}

func (s *Store) FindByName(ctx context.Context, name string) (*Product, error) {
	log.Logger(ctx).Infof("Finding user by username: %v", name)
	var product Product
	result := s.db.Where("name = ?", name).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Logger(ctx).Infof("Found product: %v", product)
	return &product, nil
}

func (s *Store) FindAllByName(ctx context.Context, names []string) ([]*Product, error) {
	log.Logger(ctx).Infof("Finding %d products by names: %v", len(names), names)
	products := make([]*Product, 0)
	result := s.db.Where("name IN (?)", names).Find(&products)
	if result.Error != nil {
		return products, result.Error
	}
	log.Logger(ctx).Infof("Found %d products", len(products))
	return products, nil
}

func (s *Store) FindAllByIds(ctx context.Context, ids []int) ([]*Product, error) {
	log.Logger(ctx).Infof("Finding %d products by ids: %v", len(ids), ids)
	products := make([]*Product, 0)
	result := s.db.Where("id IN (?)", ids).Find(&products)
	if result.Error != nil {
		return products, result.Error
	}
	log.Logger(ctx).Infof("Found %d products", len(products))
	return products, nil
}

func (s *Store) UpdateQuantity(ctx context.Context, product *Product, tx *gorm.DB) error {
	log.Logger(ctx).Infof("Updating product: %v", *product)
	if tx != nil {
		return updateInTransaction(product, tx)
	}
	return updateInTransaction(product, s.db)
}

func updateInTransaction(product *Product, tx *gorm.DB) error {
	//result := tx.Model(&product).Select("Quantity").Save(&Product{Quantity: product.Quantity})
	result := tx.Model(&product).Update("quantity", product.Quantity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return database.TransactionError{Message: "row has been updated by someone else"}
	}
	return nil
}

func (s *Store) init() {
	s.Create(context.Background(), &Product{
		Name:        "pen",
		Description: "Blue pen",
		Price:       100,
		Currency:    types.HUF,
		Quantity:    50,
	})
	s.Create(context.Background(), &Product{
		Name:        "book",
		Description: "Harry Potter 1",
		Price:       5500,
		Currency:    types.HUF,
		Quantity:    12,
	})
}
